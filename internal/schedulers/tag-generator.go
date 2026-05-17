package schedulers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-music-streamer/internal/domain/entity"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const tagGeneratorWorkerCount = 10

type TagGenerator struct {
	tagGeneratorQueue chan *entity.Song
	httpClient        *http.Client
	generateURL       string
	tenantID          string
}

type generateTagPayload struct {
	ItemID   string                 `json:"itemId"`
	Type     string                 `json:"type"`
	Metadata map[string]interface{} `json:"metadata"`
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

	return &TagGenerator{
		tagGeneratorQueue: make(chan *entity.Song),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		generateURL: endpoint,
		tenantID:    tenantID,
	}
}

func (s *TagGenerator) Start() {
	for i := 0; i < tagGeneratorWorkerCount; i++ {
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
		return fmt.Errorf("unexpected status %d from semantitag: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	log.Printf("tag generated for song id %d", song.ID)

	return nil
}

func (s *TagGenerator) Stop() {
	close(s.tagGeneratorQueue)
}

func (s *TagGenerator) AddTask(task *entity.Song) {
	s.tagGeneratorQueue <- task
}
