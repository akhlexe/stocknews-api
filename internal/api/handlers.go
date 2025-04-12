package api

import (
	"log"
	"net/http"

	"github.com/akhlexe/stocknews-api/internal/ai"
	"github.com/akhlexe/stocknews-api/internal/filter"
	"github.com/akhlexe/stocknews-api/internal/news"
	"github.com/gin-gonic/gin"
)

type Server struct {
	MultiFetcher *news.MultiFetcher
}

func NewServer(multiFetcher *news.MultiFetcher) *Server {
	return &Server{
		MultiFetcher: multiFetcher,
	}
}

func (s *Server) Run() {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/news/:ticker", func(c *gin.Context) {
		handleNews(c, s.MultiFetcher)
	})

	router.Run(":8080")
}

func handleNews(c *gin.Context, fetcher *news.MultiFetcher) {
	ticker := c.Param("ticker")
	query := c.Query("q")
	summarize := c.DefaultQuery("summarize", "false") == "true"

	articles, err := fetcher.GetNewsByTicker(ticker)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if summarize {
		var allArticles string
		for _, a := range articles {
			allArticles += a.Title + ":" + a.Summary + "\n"
		}

		summary, err := ai.SummarizeArticles(allArticles)
		if err != nil {
			log.Println("❌ Failed to generate AI summary:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Println("✅ AI summary generated successfully")
		c.JSON(http.StatusOK, gin.H{"ticker": ticker, "summary": summary})
		return
	}

	if query != "" {
		articles = filter.FilterByQuery(articles, query)
	}

	c.JSON(http.StatusOK, gin.H{"ticker": ticker, "news": articles})
}
