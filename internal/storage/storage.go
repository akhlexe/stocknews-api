package storage

import (
	"context"
	"time"

	"github.com/akhlexe/stocknews-api/internal/news"
)

// Storage defines the interface for persistent storage operations
type Storage interface {
	// SaveArticles stores news articles for a ticker with expiration time
	SaveArticles(ctx context.Context, ticker string, articles []news.Article, expiration time.Time) error

	// GetArticles retrieves news articles for a ticker if not expired
	GetArticles(ctx context.Context, ticker string) ([]news.Article, time.Time, bool, error)

	// DeleteExpired removes expired articles from storage
	DeleteExpired(ctx context.Context) error

	// Close closes the storage connection
	Close() error
}
