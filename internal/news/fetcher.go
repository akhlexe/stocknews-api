package news

import (
	"context"
	"fmt"
	"log"

	"github.com/akhlexe/stocknews-api/internal/cache"
	"github.com/akhlexe/stocknews-api/internal/models"
)

type Fetcher struct {
	APIKey string
	Cache  *cache.Cache
}

func NewFetcher(apiKey string, cache *cache.Cache) *Fetcher {
	return &Fetcher{
		APIKey: apiKey,
		Cache:  cache,
	}
}

func (f *Fetcher) FetchNews(ctx context.Context, ticker string) ([]models.Article, error) {
	cacheKey := fmt.Sprintf("news_%s", ticker)

	if cached, ok := f.Cache.Get(cacheKey); ok {
		return cached.([]models.Article), nil
	}

	log.Printf("🌍 Fetching news from API for %s", ticker)

	resp, err := GetNewsByTicker(ctx, f.APIKey, ticker)

	if err != nil {
		return nil, fmt.Errorf("error fetching news: %w", err)
	}

	f.Cache.Set(cacheKey, resp)

	return resp, nil
}
