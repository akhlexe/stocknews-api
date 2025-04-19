package cache

import (
	"context"
	"sync"
	"time"

	"github.com/akhlexe/stocknews-api/internal/models"
	"github.com/akhlexe/stocknews-api/internal/storage"
	"github.com/rs/zerolog/log"
)

type PersistentCache struct {
	storage         storage.Storage
	memory          map[string]CacheItem
	mu              sync.RWMutex
	ttl             time.Duration
	cleanupInterval time.Duration
	stopCleanup     chan struct{}
}

func NewPersistentCache(storage storage.Storage, ttl time.Duration) *PersistentCache {
	cache := &PersistentCache{
		storage:         storage,
		memory:          make(map[string]CacheItem),
		ttl:             ttl,
		cleanupInterval: ttl / 2,
		stopCleanup:     make(chan struct{}),
	}

	go cache.startCleanupTimer()

	return cache
}

func (c *PersistentCache) startCleanupTimer() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopCleanup:
			return
		}
	}
}

func (c *PersistentCache) cleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Clean up memory cache
	c.mu.RLock()
	now := time.Now()
	for k, v := range c.memory {
		if now.After(v.Expiration) {
			delete(c.memory, k)
		}
	}

	c.mu.RUnlock()
	c.mu.Unlock()

	if err := c.storage.DeleteExpired(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to delete expired items from storage")
	}
}

func (c *PersistentCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.memory = make(map[string]CacheItem)

	log.Info().Msg("Memory cache cleared")
}

func (c *PersistentCache) GetArticles(ctx context.Context, ticker string) ([]models.Article, bool) {
	cacheKey := "news_" + ticker

	// Try memory cache first (Fast path)
	c.mu.RLock()
	item, foundInMemory := c.memory[cacheKey]
	c.mu.RUnlock()

	if foundInMemory && time.Now().Before(item.Expiration) {
		articles, ok := item.Value.([]models.Article)
		if ok {
			log.Debug().Str("ticker", ticker).Msg("Cache hit (memory)")
			return articles, true
		}
	}

	// Try persistent storage (Slow path)
	articlesData, expiration, found, err := c.storage.GetArticles(ctx, ticker)
	if err != nil {
		log.Error().Err(err).Str("ticker", ticker).Msg("Error retrieving articles from storage")
		return nil, false
	}

	if found {
		var articles []models.Article
		articles, err = models.UnmarshalArticles(articlesData)
		if err != nil {
			log.Error().Err(err).Str("ticker", ticker).Msg("Error unmarshalling articles from storage")
			return nil, false
		}

		c.mu.Lock()
		c.memory[cacheKey] = CacheItem{
			Value:      articles,
			Expiration: expiration,
		}
		c.mu.Unlock()

		log.Debug().Str("ticker", ticker).Msg("Cache hit (storage)")
		return articles, true
	}

	log.Debug().Str("ticker", ticker).Msg("Cache miss: not found in memory or storage")
	return nil, false
}

func (c *PersistentCache) SetArticles(ctx context.Context, ticker string, articles []models.Article) {
	cacheKey := "news_" + ticker
	expiration := time.Now().Add(c.ttl)

	c.mu.Lock()
	c.memory[cacheKey] = CacheItem{
		Value:      articles,
		Expiration: expiration,
	}
	c.mu.Unlock()

	// Serializa articles for storage
	articlesData, err := models.MarshalArticles(articles)

	if err != nil {
		log.Error().Err(err).Str("ticker", ticker).Msg("Error marshalling articles for storage")
		return
	}

	// Update persistent storage
	if err := c.storage.SaveArticles(ctx, ticker, articlesData, expiration); err != nil {
		log.Error().Err(err).Str("ticker", ticker).Msg("Error saving articles to storage")
		return
	} else {
		log.Debug().Str("ticker", ticker).Msg("Articles saved to storage")
	}
}

func (c *PersistentCache) Stop() {
	close(c.stopCleanup)
}
