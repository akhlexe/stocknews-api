package news

import (
	"context"
	"fmt"

	"github.com/akhlexe/stocknews-api/internal/cache"
	"github.com/akhlexe/stocknews-api/internal/models"
	"github.com/rs/zerolog/log"
)

type AlphaVantageFetcher struct {
	apiKey string
	cache  *cache.PersistentCache
}

func NewAlphaVantageFetcher(apiKey string, cache *cache.PersistentCache) *AlphaVantageFetcher {
	return &AlphaVantageFetcher{
		apiKey: apiKey,
		cache:  cache,
	}
}

func (f *AlphaVantageFetcher) GetNewsByTicker(ctx context.Context, ticker string) ([]models.Article, error) {
	if cached, ok := f.cache.GetArticles(ctx, ticker); ok {
		return cached, nil
	}

	log.Info().
		Str("ticker", ticker).
		Str("method", "GetNewsByTicker").
		Msgf("🌍 Fetching news from API for %s", ticker)

	resp, err := GetNewsByTicker(ctx, f.apiKey, ticker)

	if err != nil {
		return nil, fmt.Errorf("error fetching news: %w", err)
	}

	f.cache.SetArticles(ctx, ticker, resp)

	return resp, nil
}
