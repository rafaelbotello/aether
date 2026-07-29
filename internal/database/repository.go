package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rafaelbotello/aether/internal/models"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) SaveScrapeResult(ctx context.Context, result models.ScrapeResult) error {

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	mangaInsert := `INSERT INTO mangas (canonical_id, title) VALUES($1, $2) ON CONFLICT (canonical_id) DO NOTHING;`

	if _, err := tx.Exec(ctx, mangaInsert, result.MangaID, result.MangaTitle); err != nil {
		return fmt.Errorf("failed to upsert manga: %w", err)
	}

	chapterInsert := `INSERT INTO chapters (manga_id, chapter_number, language, title) VALUES($1,$2,$3,$4) ON CONFLICT(manga_id, chapter_number, language) DO UPDATE SET
	title = COALESCE(NULLIF(EXCLUDED.title, ''), chapters.title) RETURNING id;`
	var chapterID int
	err = tx.QueryRow(ctx, chapterInsert, result.MangaID, result.ChapterNumber, result.Language, result.Title).Scan(&chapterID)
	if err != nil {
		return fmt.Errorf("failed to upsert chapter: %w", err)
	}

	sourceInsert := `INSERT INTO chapter_sources (chapter_id, source_name, source_url) 
	VALUES($1,$2,$3) ON CONFLICT (chapter_id, source_name) DO UPDATE SET source_url = EXCLUDED.source_url;`

	if _, err := tx.Exec(ctx, sourceInsert, chapterID, result.SourceName, result.SourceURL); err != nil {
		return fmt.Errorf("failed to upsert source: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
