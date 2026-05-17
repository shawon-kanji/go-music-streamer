package song

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"go-music-streamer/internal/domain/apperror"
	"go-music-streamer/internal/domain/dto"
)

const defaultSemantiTagSearchURL = "http://localhost:8000/search"
const defaultSemantiTagTenantID = "1234"
const defaultSemantiTagSearchTimeoutSec = 20

type semantiTagSearchItem struct {
	ID       string                 `json:"id"`
	Distance float64                `json:"distance"`
	Metadata map[string]interface{} `json:"metadata"`
}

type semantiTagSearchResponse struct {
	Query string                 `json:"query"`
	Count int                    `json:"count"`
	Items []semantiTagSearchItem `json:"items"`
}

func (useCase *songUseCase) SearchSongs(req *dto.SearchSongsRequest) (*dto.SemanticSearchResponse, error) {
	if req == nil {
		return nil, apperror.New(apperror.BadRequest, "search request is required")
	}

	query := strings.TrimSpace(req.Query)
	if query == "" {
		return nil, apperror.New(apperror.BadRequest, "query is required")
	}

	payload := map[string]interface{}{
		"query": query,
	}

	if req.NResults != nil {
		payload["nResults"] = *req.NResults
	}

	if req.ItemType != nil {
		itemType := strings.TrimSpace(*req.ItemType)
		if itemType != "" {
			payload["itemType"] = itemType
		}
	}

	if req.IncludeTagVectors != nil {
		payload["includeTagVectors"] = *req.IncludeTagVectors
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, apperror.Newf(apperror.InternalError, "failed to marshal search payload: %v", err)
	}

	searchURL := os.Getenv("SEMANTITAG_SEARCH_URL")
	if searchURL == "" {
		searchURL = defaultSemantiTagSearchURL
	}

	tenantID := os.Getenv("SEMANTITAG_TENANT_ID")
	if tenantID == "" {
		tenantID = defaultSemantiTagTenantID
	}

	timeoutSec := defaultSemantiTagSearchTimeoutSec
	if val := os.Getenv("SEMANTITAG_SEARCH_TIMEOUT_SEC"); val != "" {
		if parsed, parseErr := strconv.Atoi(val); parseErr == nil && parsed > 0 {
			timeoutSec = parsed
		}
	}

	client := &http.Client{Timeout: time.Duration(timeoutSec) * time.Second}

	httpReq, err := http.NewRequest(http.MethodPost, searchURL, bytes.NewReader(body))
	if err != nil {
		return nil, apperror.Newf(apperror.InternalError, "failed to create search request: %v", err)
	}

	httpReq.Header.Set("tenantId", tenantID)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, apperror.Newf(apperror.InternalError, "failed to call semantic search service: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		trimmed := strings.TrimSpace(string(respBody))
		if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError {
			return nil, apperror.Newf(apperror.BadRequest, "semantic search request failed: %s", trimmed)
		}

		return nil, apperror.Newf(apperror.InternalError, "semantic search service failed with status %d: %s", resp.StatusCode, trimmed)
	}

	var semantiResp semantiTagSearchResponse
	if err := json.Unmarshal(respBody, &semantiResp); err != nil {
		return nil, apperror.Newf(apperror.InternalError, "failed to parse semantic search response: %v", err)
	}

	items := make([]dto.SemanticSearchItem, 0, len(semantiResp.Items))
	for _, item := range semantiResp.Items {
		itemID, extractErr := extractItemID(item)
		if extractErr != nil {
			continue
		}

		songEntity, repoErr := useCase.repo.GetSongByID(itemID)
		if repoErr != nil {
			var appErr *apperror.AppError
			if errors.As(repoErr, &appErr) && appErr.Code == apperror.NotFound {
				continue
			}

			return nil, apperror.Newf(apperror.InternalError, "failed to fetch song %d from database: %v", itemID, repoErr)
		}

		items = append(items, dto.SemanticSearchItem{
			ItemID:   itemID,
			Distance: item.Distance,
			Song: &dto.SongResponse{
				ID:        songEntity.ID,
				Title:     songEntity.Title,
				Artist:    songEntity.Artist,
				Album:     songEntity.Album,
				Duration:  songEntity.Duration,
				URL:       songEntity.URL,
				LikeCount: songEntity.LikeCount,
				Genre:     songEntity.Genre,
				Thumbnail: songEntity.Thumbnail,
			},
		})
	}

	searchResp := &dto.SemanticSearchResponse{
		Query: semantiResp.Query,
		Count: len(items),
		Items: items,
	}

	return searchResp, nil
}

func extractItemID(item semantiTagSearchItem) (uint, error) {
	if rawItemID, ok := item.Metadata["itemId"]; ok {
		switch v := rawItemID.(type) {
		case string:
			parsed, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64)
			if err == nil {
				return uint(parsed), nil
			}
		case float64:
			if v >= 0 {
				return uint(v), nil
			}
		}
	}

	parts := strings.Split(item.ID, ":")
	if len(parts) >= 3 {
		parsed, err := strconv.ParseUint(parts[2], 10, 64)
		if err == nil {
			return uint(parsed), nil
		}
	}

	return 0, apperror.Newf(apperror.BadRequest, "unable to extract itemId from semantic result id=%s", item.ID)
}
