package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(connectionString string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Error().Err(err).Msg("Failed to open Postgres database")
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Error().Err(err).Msg("Failed to ping Postgres database")
		return nil, err
	}

	storage := &PostgresStorage{db: db}

	if err := storage.initialize(); err != nil {
		db.Close()
		return nil, err
	}

	return storage, nil
}

func (s *PostgresStorage) initialize() error {
	// Create the articles table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS articles (
		ticker TEXT NOT NULL,
		data BYTEA NOT NULL,
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

func (s *PostgresStorage) SaveArticles(ctx context.Context, ticker string, articles []byte, expiration time.Time) error {
	query := `
	INSERT INTO articles (ticker, data, expiration, created_at)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT(ticker)
	DO UPDATE SET data = $2, expiration = $3, created_at = $4
	`
	_, err := s.db.ExecContext(ctx, query, ticker, articles, expiration, time.Now())
	if err != nil {
		log.Error().Err(err).Msg("Failed to save articles")
		return err
	}

	log.Debug().Str("ticker", ticker).Msg("Saved articles to SQLite storage")
	return nil
}

func (s *PostgresStorage) GetArticles(ctx context.Context, ticker string) ([]byte, time.Time, bool, error) {
	query := `
	SELECT data, expiration FROM articles 
	WHERE ticker = $1 AND expiration > $2
	`
	row := s.db.QueryRowContext(ctx, query, ticker)

	var data []byte
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

	log.Debug().Str("ticker", ticker).Msg("Retrieved articles from SQLite storage")
	return data, expiration, true, nil
}

func (s *PostgresStorage) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM articles WHERE expiration <= $1`

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
		log.Info().Int64("count", count).Msg("Deleted expired articles from Postgres storage")
	}

	return nil
}

func (s *PostgresStorage) Close() error {
	log.Info().Msg("Closing PostgreSQL database connection")
	return s.db.Close()
}
