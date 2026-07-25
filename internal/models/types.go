package models

import "context"

type ScrapeResult struct {
	SourceName    string
	SourceURL     string
	ChapterNumber string
	Title         string
	Language      string
}

type Provider interface {
	FetchUpdates(ctx context.Context, url string) ([]ScrapeResult, error)
}
