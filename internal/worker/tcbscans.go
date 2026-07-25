package worker

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rafaelbotello/aether/internal/models"
)

var chapterNumberRegex = regexp.MustCompile(`\d+`)

type TCBScans struct {
	client *http.Client
}

func (t *TCBScans) FetchUpdates(ctx context.Context, targetURL string) ([]models.ScrapeResult, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create tcb request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (platform; rv:gecko-version) Gecko/gecko-trail Firefox/firefox-version")

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute tcb fetchUpdate req: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create new document: %w", err)
	}

	var results []models.ScrapeResult

	doc.Find("a.bg-card").Each(func(i int, s *goquery.Selection) {

		chapterID, _ := s.Attr("href")

		rawChapterNumber := s.Find("div.font-bold").Text()
		match := chapterNumberRegex.FindString(rawChapterNumber)

		rawTitle := s.Find("div.text-gray-500").Text()
		cleanTitle := strings.TrimSpace(rawTitle)

		result := models.ScrapeResult{
			SourceName:    "TCBScans",
			SourceURL:     "https://tcbonepiecechapters.com" + chapterID,
			ChapterNumber: match,
			Title:         cleanTitle,
			Language:      "en",
		}
		results = append(results, result)

	})

	return results, nil

}
