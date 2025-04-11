package api

import (
	"net/http"
	"time"

	"github.com/akhlexe/stocknews-api/internal/cache"
	"github.com/akhlexe/stocknews-api/internal/filter"
	"github.com/akhlexe/stocknews-api/internal/news"
	"github.com/gin-gonic/gin"
)

func Run() {
	router := gin.Default()

	apiKey := "YOUR_API_KEY"
	cache := cache.NewCache(10 * time.Minute)
	fetcher := news.NewFetcher(apiKey, cache)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/news/:ticker", func(c *gin.Context) {
		handleNews(c, fetcher)
	})

	router.Run(":8080")
}

func handleNews(c *gin.Context, fetcher *news.Fetcher) {
	ticker := c.Param("ticker")
	query := c.Query("q")

	articles, err := fetcher.FetchNews(ticker)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if query != "" {
		articles = filter.FilterByQuery(articles, query)
	}

	c.JSON(http.StatusOK, gin.H{"ticker": ticker, "news": articles})
}
