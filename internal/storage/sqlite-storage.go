package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/akhlexe/stocknews-api/internal/news"
	"github.com/rs/zerolog/log"
)

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		log.Error().Err(err).Msg("Failed to open SQLite database")
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Error().Err(err).Msg("Failed to ping SQLite database")
		return nil, err
	}

	storage := &SQLiteStorage{db: db}

	if err := storage.initialize(); err != nil {
		db.Close()
		return nil, err
	}

	return storage, nil
}

func (s *SQLiteStorage) initialize() error {
	// Create the articles table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS articles (
		ticker TEXT NOT NULL,
		data TEXT NOT NULL,
		expiration TIMESTAMP NOT NULL,
		created_at TIMESTAMP NOT NULL,
		PRIMARY KEY (ticker)
	);
	CREATE INDEX IF NOT EXISTS idx_expiration ON articles(expiration);
	`

	if _, err := s.db.Exec(query); err != nil {
		log.Error().Err(err).Msg("Failed to create articles table")
		return err
	}

	return nil
}

func (s *SQLiteStorage) SaveArticles(ctx context.Context, ticker string, articles []news.Article, expiration time.Time) error {

	data, err := json.Marshal(articles)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal articles")
		return err
	}

	query := `
	INSERT OR REPLACE INTO articles (ticker, data, expiration, created_at)
	VALUES (?, ?, ?, ?)
	`
	_, err = s.db.ExecContext(ctx, query, ticker, data, expiration, time.Now())
	if err != nil {
		log.Error().Err(err).Msg("Failed to save articles")
		return err
	}

	log.Debug().Str("ticker", ticker).Msg("Saved articles to SQLite storage")
	return nil
}

func (s *SQLiteStorage) GetArticles(ctx context.Context, ticker string) ([]news.Article, time.Time, bool, error) {
	query := `
	SELECT data, expiration 
	FROM articles 
	WHERE ticker = ?
	`
	row := s.db.QueryRowContext(ctx, query, ticker)

	var data string
	var expiration time.Time

	err := row.Scan(&data, &expiration)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug().Str("ticker", ticker).Msg("No articles found in SQLite storage")
			return nil, time.Time{}, false, nil
		}
		log.Error().Err(err).Msg("Failed to retrieve articles")
		return nil, time.Time{}, false, err
	}

	var articles []news.Article
	err = json.Unmarshal([]byte(data), &articles)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal articles")
		return nil, time.Time{}, false, err
	}

	log.Debug().Str("ticker", ticker).Msg("Retrieved articles from SQLite storage")
	return articles, expiration, true, nil
}

func (s *SQLiteStorage) DeleteExpired(ctx context.Context) error {
	query := `
	DELETE FROM articles 
	WHERE expiration < ?
	`

	result, err := s.db.ExecContext(ctx, query, time.Now())
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete expired articles")
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get number of deleted rows")
		return nil
	}

	if count > 0 {
		log.Info().Int64("count", count).Msg("Deleted expired articles from SQLite storage")
	}

	return nil
}

func (s *SQLiteStorage) Close() error {
	log.Info().Msg("Closing SQLite storage")
	return s.db.Close()
}
