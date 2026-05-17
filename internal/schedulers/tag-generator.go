package schedulers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-music-streamer/internal/domain/entity"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultTagGeneratorWorkerCount = 10
const defaultTagGeneratorQueueSize = 1000
const defaultTagGeneratorHTTPTimeoutSec = 30
const defaultTagGeneratorMaxRetries = 3
const defaultTagGeneratorRetryBackoffMs = 500

type TagGenerator struct {
	tagGeneratorQueue chan *entity.Song
	httpClient        *http.Client
	generateURL       string
	tenantID          string
	workerCount       int
	maxRetries        int
	retryBackoff      time.Duration
}

type generateTagPayload struct {
	ItemID   string                 `json:"itemId"`
	Type     string                 `json:"type"`
	Metadata map[string]interface{} `json:"metadata"`
}

type httpStatusError struct {
	statusCode int
	body       string
}

func (e *httpStatusError) Error() string {
	return fmt.Sprintf("unexpected status %d from semantitag: %s", e.statusCode, strings.TrimSpace(e.body))
}

func NewTagGenerator() *TagGenerator {
	endpoint := os.Getenv("SEMANTITAG_GENERATE_URL")
	if endpoint == "" {
		endpoint = "http://localhost:8000/generate"
	}

	tenantID := os.Getenv("SEMANTITAG_TENANT_ID")
	if tenantID == "" {
		tenantID = "1234"
	}

	queueSize := defaultTagGeneratorQueueSize
	if val := os.Getenv("TAG_GENERATOR_QUEUE_SIZE"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
			queueSize = parsed
		}
	}

	workerCount := defaultTagGeneratorWorkerCount
	if val := os.Getenv("TAG_GENERATOR_WORKERS"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
			workerCount = parsed
		}
	}

	timeoutSec := defaultTagGeneratorHTTPTimeoutSec
	if val := os.Getenv("TAG_GENERATOR_HTTP_TIMEOUT_SEC"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
			timeoutSec = parsed
		}
	}

	maxRetries := defaultTagGeneratorMaxRetries
	if val := os.Getenv("TAG_GENERATOR_MAX_RETRIES"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed >= 0 {
			maxRetries = parsed
		}
	}

	backoffMs := defaultTagGeneratorRetryBackoffMs
	if val := os.Getenv("TAG_GENERATOR_RETRY_BACKOFF_MS"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
			backoffMs = parsed
		}
	}

	return &TagGenerator{
		tagGeneratorQueue: make(chan *entity.Song, queueSize),
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSec) * time.Second,
		},
		generateURL: endpoint,
		tenantID:    tenantID,
		workerCount: workerCount,
		maxRetries:  maxRetries,
		retryBackoff: time.Duration(backoffMs) *
			time.Millisecond,
	}
}

func (s *TagGenerator) Start() {
	for i := 0; i < s.workerCount; i++ {
		go func() {
			for song := range s.tagGeneratorQueue {
				if song == nil {
					continue
				}

				if err := s.generateTags(song); err != nil {
					log.Printf("tag generation failed for song id %d: %v", song.ID, err)
				}
			}
		}()
	}
}

func (s *TagGenerator) generateTags(song *entity.Song) error {
	payload := generateTagPayload{
		ItemID: strconv.FormatUint(uint64(song.ID), 10),
		Type:   "song",
		Metadata: map[string]interface{}{
			"artist":   song.Artist,
			"genre":    song.Genre,
			"title":    song.Title,
			"album":    song.Album,
			"duration": song.Duration,
		},
	}

	maxAttempts := s.maxRetries + 1
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := s.sendGenerateRequest(payload)
		if err == nil {
			log.Printf("tag generated for song id %d", song.ID)
			return nil
		}

		lastErr = err
		if !isRetryableError(err) || attempt == maxAttempts {
			break
		}

		delay := s.retryBackoff * time.Duration(1<<(attempt-1))
		if delay > 8*time.Second {
			delay = 8 * time.Second
		}

		log.Printf("tag generation retry for song id %d (attempt %d/%d): %v", song.ID, attempt+1, maxAttempts, err)
		time.Sleep(delay)
	}

	return lastErr
}

func (s *TagGenerator) sendGenerateRequest(payload generateTagPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.generateURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("tenantId", s.tenantID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)

	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return &httpStatusError{statusCode: resp.StatusCode, body: string(respBody)}
	}

	return nil
}

func isRetryableError(err error) bool {
	var statusErr *httpStatusError
	if errors.As(err, &statusErr) {
		return statusErr.statusCode == http.StatusTooManyRequests || statusErr.statusCode == http.StatusRequestTimeout || statusErr.statusCode >= http.StatusInternalServerError
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "connection refused") || strings.Contains(errStr, "connection reset") || strings.Contains(errStr, "context deadline exceeded") {
		return true
	}

	return false
}

func (s *TagGenerator) Stop() {
	close(s.tagGeneratorQueue)
}

func (s *TagGenerator) AddTask(task *entity.Song) {
	if task == nil {
		return
	}

	select {
	case s.tagGeneratorQueue <- task:
	default:
		log.Printf("tag generator queue full, dropping song id %d", task.ID)
	}
}
