package worker

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/rafaelbotello/aether/internal/models"
	"golang.org/x/sync/errgroup"
)

type Pipeline struct {
	client      *http.Client
	logger      *slog.Logger
	workerCount int
}

func NewPipeline(logger *slog.Logger, client *http.Client, workerCount int) *Pipeline {
	return &Pipeline{
		logger:      logger,
		client:      client,
		workerCount: workerCount,
	}
}

func (p *Pipeline) Run(ctx context.Context, urls []string) error {
	jobs := make(chan string, 100)
	results := make(chan models.ScrapeResult, 100)

	g, ctx := errgroup.WithContext(context.TODO())

	g.Go(func() error {
		defer close(jobs)

		for _, url := range urls {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case jobs <- url:
			}
		}

		return nil
	})

	for range p.workerCount {
		g.Go(func() error {
			return p.worker(ctx, jobs, results)
		})
	}

	collectorDone := make(chan struct{})
	go func() {
		for result := range results {
			p.logger.Info("logging", "source", result.SourceName, "chapter", result.ChapterNumber, "title", result.Title, "source url", result.SourceURL)
		}
		close(collectorDone)
	}()

	err := g.Wait()
	close(results)
	<-collectorDone

	return err
}

func (p *Pipeline) worker(ctx context.Context, jobs <-chan string, results chan<- models.ScrapeResult) error {

	providers := map[string]models.Provider{
		"api.mangadex.org":        &MangaDex{client: p.client},
		"tcbonepiecechapters.com": &TCBScans{client: p.client},
	}

	for jobURL := range jobs {

		parsedURL, err := url.Parse(jobURL)
		if err != nil {
			p.logger.Error("invalid URL in jobs channel", "url", jobURL, "error", err)
			continue
		}

		provider, exists := providers[parsedURL.Host]
		if !exists {
			p.logger.Warn("no provider built for this domain", "domain", parsedURL.Host)
			continue
		}

		scrapeResults, err := provider.FetchUpdates(ctx, jobURL)
		if err != nil {
			p.logger.Error("failed to fetch updates from", "source", jobURL, "error", err)
			continue
		}

		for _, res := range scrapeResults {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case results <- res:

			}
		}
	}
	return nil
}

// func (p *Pipeline) fetchUpdate(ctx context.Context, job string) (*GetUpdateResponse, error) {

// 	req, err := http.NewRequestWithContext(ctx, http.MethodGet, job, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create request: %w", err)
// 	}

// 	resp, err := p.client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to execute fetchUpdate request: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// 	}

// 	var response GetUpdateResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {

// 		return nil, fmt.Errorf("failed to decode fetchUpdate response: %w", err)
// 	}

// 	return &response, nil
// }
