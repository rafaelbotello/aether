package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rafaelbotello/aether/internal/models"
)

type MangaDex struct {
	client *http.Client
}

type mdResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			Chapter string `json:"chapter"`
			Title   string `json:"title"`
		} `json:"attributes"`
	} `json:"data"`
}

func (m *MangaDex) FetchUpdates(ctx context.Context, targetURL string) ([]models.ScrapeResult, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create md request: %w", err)
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute fetchUpdates request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response mdResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode fetchUpdate response: %w", err)
	}

	var results []models.ScrapeResult
	for _, item := range response.Data {
		result := models.ScrapeResult{
			SourceName:    "MangaDex",
			SourceURL:     "https://mangadex.org/chapter/" + item.ID,
			ChapterNumber: item.Attributes.Chapter,
			Title:         item.Attributes.Title,
			Language:      "en",
		}
		results = append(results, result)
	}

	return results, nil
}
