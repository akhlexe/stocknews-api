package news

import (
	"fmt"

	"github.com/akhlexe/stocknews-api/internal/cache"
	"github.com/rs/zerolog/log"
)

type AlphaVantageFetcher struct {
	APIKey string
	Cache  *cache.Cache
}

func NewAlphaVantageFetcher(apiKey string, cache *cache.Cache) *AlphaVantageFetcher {
	return &AlphaVantageFetcher{
		APIKey: apiKey,
		Cache:  cache,
	}
}

func (f *AlphaVantageFetcher) GetNewsByTicker(ticker string) ([]Article, error) {
	cacheKey := fmt.Sprintf("news_%s", ticker)

	if cached, ok := f.Cache.Get(cacheKey); ok {
		return cached.([]Article), nil
	}

	log.Info().
		Str("ticker", ticker).
		Str("method", "GetNewsByTicker").
		Msgf("🌍 Fetching news from API for %s", ticker)

	resp, err := GetNewsByTicker(f.APIKey, ticker)

	if err != nil {
		return nil, fmt.Errorf("error fetching news: %w", err)
	}

	f.Cache.Set(cacheKey, resp)

	return resp, nil
}
