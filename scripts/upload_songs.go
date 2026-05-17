package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	APIURL             string `json:"apiUrl"`
	AuthToken          string `json:"authToken"`
	LoginURL           string `json:"loginUrl"`
	LoginEmail         string `json:"loginEmail"`
	LoginPassword      string `json:"loginPassword"`
	FilePath           string `json:"filePath"`
	CSVPath            string `json:"csvPath"`
	TotalUploads       int    `json:"totalUploads"`
	Concurrency        int    `json:"concurrency"`
	RequestTimeoutSec  int    `json:"requestTimeoutSec"`
	AppendIndexToTitle bool   `json:"appendIndexToTitle"`
}

type SongRow struct {
	Title    string
	Artist   string
	Album    string
	Genre    string
	Duration uint
}

type uploadResult struct {
	Index  int
	Status int
	Body   string
	Err    error
	Title  string
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginData struct {
	Token string `json:"token"`
}

type apiResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func main() {
	defaultConfigPath := filepath.Join("scripts", "upload_songs.config.json")
	configPath := flag.String("config", defaultConfigPath, "Path to uploader JSON config")
	flag.Parse()

	cfg, err := loadConfig(*configPath)
	if err != nil {
		exitErr("config error: %v", err)
	}

	rows, err := loadSongRows(cfg.CSVPath)
	if err != nil {
		exitErr("csv error: %v", err)
	}

	client := &http.Client{Timeout: time.Duration(cfg.RequestTimeoutSec) * time.Second}

	token, err := resolveAuthToken(client, cfg)
	if err != nil {
		exitErr("auth error: %v", err)
	}
	cfg.AuthToken = token

	fmt.Println("Starting uploads")
	fmt.Printf("  API_URL=%s\n", cfg.APIURL)
	fmt.Printf("  FILE_PATH=%s\n", cfg.FilePath)
	fmt.Printf("  CSV_PATH=%s\n", cfg.CSVPath)
	fmt.Printf("  TOTAL_UPLOADS=%d\n", cfg.TotalUploads)
	fmt.Printf("  CONCURRENCY=%d\n", cfg.Concurrency)

	jobs := make(chan int)
	results := make(chan uploadResult)

	var wg sync.WaitGroup
	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				row := rows[(idx-1)%len(rows)]
				if cfg.AppendIndexToTitle {
					row.Title = fmt.Sprintf("%s-%d", row.Title, idx)
				}

				res := uploadSong(client, cfg, row, idx)
				results <- res
			}
		}()
	}

	go func() {
		for i := 1; i <= cfg.TotalUploads; i++ {
			jobs <- i
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	var successCount int64
	var failCount int64

	for res := range results {
		if res.Err == nil && res.Status >= http.StatusOK && res.Status < http.StatusMultipleChoices {
			atomic.AddInt64(&successCount, 1)
			fmt.Printf("OK index=%d status=%d title=%s\n", res.Index, res.Status, res.Title)
			continue
		}

		atomic.AddInt64(&failCount, 1)
		if res.Err != nil {
			fmt.Printf("FAIL index=%d err=%v title=%s\n", res.Index, res.Err, res.Title)
		} else {
			fmt.Printf("FAIL index=%d status=%d body=%s title=%s\n", res.Index, res.Status, res.Body, res.Title)
		}
	}

	fmt.Printf("Upload summary: success=%d fail=%d\n", successCount, failCount)
	if failCount > 0 {
		os.Exit(1)
	}
}

func resolveAuthToken(client *http.Client, cfg Config) (string, error) {
	if strings.TrimSpace(cfg.AuthToken) != "" {
		return strings.TrimSpace(cfg.AuthToken), nil
	}

	if strings.TrimSpace(cfg.LoginEmail) == "" || strings.TrimSpace(cfg.LoginPassword) == "" {
		return "", fmt.Errorf("authToken is empty and loginEmail/loginPassword are not configured")
	}

	loginURL := strings.TrimSpace(cfg.LoginURL)
	if loginURL == "" {
		derivedLoginURL, err := deriveLoginURL(cfg.APIURL)
		if err != nil {
			return "", err
		}
		loginURL = derivedLoginURL
	}

	payload := loginRequest{
		Email:    strings.TrimSpace(cfg.LoginEmail),
		Password: cfg.LoginPassword,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal login request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, loginURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("login failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	var wrapper apiResponse
	if err := json.Unmarshal(respBody, &wrapper); err != nil {
		return "", fmt.Errorf("decode login response: %w", err)
	}

	var data loginData
	if err := json.Unmarshal(wrapper.Data, &data); err != nil {
		return "", fmt.Errorf("decode login data: %w", err)
	}

	if strings.TrimSpace(data.Token) == "" {
		return "", fmt.Errorf("login succeeded but token is empty")
	}

	fmt.Printf("Login successful via %s\n", loginURL)
	return data.Token, nil
}

func deriveLoginURL(apiURL string) (string, error) {
	parsed, err := url.Parse(apiURL)
	if err != nil {
		return "", fmt.Errorf("invalid apiUrl: %w", err)
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("invalid apiUrl: must include scheme and host")
	}

	parsed.Path = "/auth/login"
	parsed.RawQuery = ""
	parsed.Fragment = ""

	return parsed.String(), nil
}

func uploadSong(client *http.Client, cfg Config, row SongRow, index int) uploadResult {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := addFilePart(writer, "file", cfg.FilePath); err != nil {
		_ = writer.Close()
		return uploadResult{Index: index, Err: err, Title: row.Title}
	}

	_ = writer.WriteField("title", row.Title)
	_ = writer.WriteField("artist", row.Artist)
	_ = writer.WriteField("album", row.Album)
	_ = writer.WriteField("genre", row.Genre)
	_ = writer.WriteField("duration", strconv.FormatUint(uint64(row.Duration), 10))

	if err := writer.Close(); err != nil {
		return uploadResult{Index: index, Err: err, Title: row.Title}
	}

	req, err := http.NewRequest(http.MethodPost, cfg.APIURL, &body)
	if err != nil {
		return uploadResult{Index: index, Err: err, Title: row.Title}
	}

	req.Header.Set("Authorization", "Bearer "+cfg.AuthToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return uploadResult{Index: index, Err: err, Title: row.Title}
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	return uploadResult{
		Index:  index,
		Status: resp.StatusCode,
		Body:   strings.TrimSpace(string(respBody)),
		Title:  row.Title,
	}
}

func addFilePart(writer *multipart.Writer, fieldName, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	part, err := writer.CreateFormFile(fieldName, filepath.Base(path))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	return err
}

func loadSongRows(csvPath string) ([]SongRow, error) {
	f, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true
	r.FieldsPerRecord = -1
	r.LazyQuotes = true

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("csv must include header and at least 1 row")
	}

	headerMap := make(map[string]int)
	for i, col := range records[0] {
		headerMap[strings.ToLower(strings.TrimSpace(col))] = i
	}

	titleIdx, err := resolveHeaderIndex(headerMap, []string{"title", "track_name"}, "title")
	if err != nil {
		return nil, err
	}

	artistIdx, err := resolveHeaderIndex(headerMap, []string{"artist", "artists"}, "artist")
	if err != nil {
		return nil, err
	}

	albumIdx, err := resolveHeaderIndex(headerMap, []string{"album", "album_name"}, "album")
	if err != nil {
		return nil, err
	}

	genreIdx, err := resolveHeaderIndex(headerMap, []string{"genre", "track_genre"}, "genre")
	if err != nil {
		return nil, err
	}

	durationIdx, durationInMS, err := resolveDurationIndex(headerMap)
	if err != nil {
		return nil, err
	}

	rows := make([]SongRow, 0, len(records)-1)
	skipped := 0
	for _, rec := range records[1:] {
		if len(rec) == 0 {
			skipped++
			continue
		}

		if titleIdx >= len(rec) || artistIdx >= len(rec) || albumIdx >= len(rec) || genreIdx >= len(rec) || durationIdx >= len(rec) {
			skipped++
			continue
		}

		durationStr := strings.TrimSpace(rec[durationIdx])
		duration, err := strconv.ParseUint(durationStr, 10, 64)
		if err != nil {
			skipped++
			continue
		}
		if durationInMS {
			duration = duration / 1000
			if duration == 0 {
				duration = 1
			}
		}

		row := SongRow{
			Title:    strings.TrimSpace(rec[titleIdx]),
			Artist:   strings.TrimSpace(rec[artistIdx]),
			Album:    strings.TrimSpace(rec[albumIdx]),
			Genre:    strings.TrimSpace(rec[genreIdx]),
			Duration: uint(duration),
		}

		if row.Title == "" || row.Artist == "" || row.Album == "" || row.Genre == "" {
			skipped++
			continue
		}

		rows = append(rows, row)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("csv has no usable rows after validation")
	}

	fmt.Printf("Loaded %d rows from %s (skipped %d invalid rows)\n", len(rows), csvPath, skipped)

	return rows, nil
}

func resolveHeaderIndex(headerMap map[string]int, candidates []string, fieldName string) (int, error) {
	for _, candidate := range candidates {
		if idx, ok := headerMap[candidate]; ok {
			return idx, nil
		}
	}

	return -1, fmt.Errorf("csv missing required column for %s (accepted: %s)", fieldName, strings.Join(candidates, ", "))
}

func resolveDurationIndex(headerMap map[string]int) (int, bool, error) {
	if idx, ok := headerMap["duration"]; ok {
		return idx, false, nil
	}
	if idx, ok := headerMap["duration_seconds"]; ok {
		return idx, false, nil
	}
	if idx, ok := headerMap["duration_sec"]; ok {
		return idx, false, nil
	}
	if idx, ok := headerMap["duration_ms"]; ok {
		return idx, true, nil
	}

	return -1, false, fmt.Errorf("csv missing required column for duration (accepted: duration, duration_seconds, duration_sec, duration_ms)")
}

func loadConfig(path string) (Config, error) {
	var cfg Config

	b, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}

	if cfg.APIURL == "" {
		return cfg, fmt.Errorf("apiUrl is required")
	}
	if strings.TrimSpace(cfg.AuthToken) == "" {
		if strings.TrimSpace(cfg.LoginEmail) == "" || strings.TrimSpace(cfg.LoginPassword) == "" {
			return cfg, fmt.Errorf("either authToken or loginEmail/loginPassword is required")
		}
	}
	if cfg.FilePath == "" {
		return cfg, fmt.Errorf("filePath is required")
	}
	if cfg.CSVPath == "" {
		return cfg, fmt.Errorf("csvPath is required")
	}
	if cfg.TotalUploads < 1 {
		return cfg, fmt.Errorf("totalUploads must be >= 1")
	}
	if cfg.Concurrency < 1 {
		return cfg, fmt.Errorf("concurrency must be >= 1")
	}
	if cfg.RequestTimeoutSec < 1 {
		cfg.RequestTimeoutSec = 30
	}

	if _, err := os.Stat(cfg.FilePath); err != nil {
		return cfg, fmt.Errorf("filePath not found: %w", err)
	}

	return cfg, nil
}

func exitErr(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
