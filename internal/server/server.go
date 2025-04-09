package server

import (
	"net/http"

	"github.com/akhlexe/stocknews-api/internal/stocks"
	"github.com/gin-gonic/gin"
)

func Run() {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/news/:ticker", handleNews)

	router.Run(":8080")
}

func handleNews(c *gin.Context) {
	ticker := c.Param("ticker")
	query := c.Query("q")

	articles, err := stocks.GetNewsByTicker(ticker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if query != "" {
		articles = stocks.FilterByQuery(articles, query)
	}

	c.JSON(http.StatusOK, gin.H{"ticker": ticker, "news": articles})
}
