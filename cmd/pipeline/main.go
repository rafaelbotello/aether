package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/rafaelbotello/aether/internal/worker"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	urls := []string{
		"https://tcbonepiecechapters.com/mangas/5/one-piece",
		"https://api.mangadex.org/manga/a1c7c817-4e59-43b7-9365-09675a149a6f/feed?translatedLanguage[]=en",
	}

	pipeline := worker.NewPipeline(logger, client, 5)

	logger.Info("starting ingestion pipeline")
	if err := pipeline.Run(context.Background(), urls); err != nil {
		logger.Error("pipeline encountered a fatal error", "error", err)
		os.Exit(1)
	}

	logger.Info("pipeline finished checking for updates successfully")
}
