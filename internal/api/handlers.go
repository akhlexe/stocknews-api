package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/akhlexe/stocknews-api/internal/ai"
	"github.com/akhlexe/stocknews-api/internal/apperrors"
	"github.com/akhlexe/stocknews-api/internal/filter"
	"github.com/akhlexe/stocknews-api/internal/news"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var validTickerRegex = regexp.MustCompile(`^[A-Z]{1,10}$`)

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

	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("status", fmt.Sprintf("%d", c.Writer.Status())).
			Dur("latency", latency).
			Msg("Request handled")

	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/news/:ticker", func(c *gin.Context) {
		handleNews(c, s.MultiFetcher)
	})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
		log.Info().Msgf("Defaulting to port %s", port)
	} else {
		log.Info().Msgf("Listening on port %s", port)
	}

	address := fmt.Sprintf(":%s", port)

	log.Info().Str("address", address).Msg("Starting server")

	err := router.Run(address)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("address", address).
			Msg("Failed to start server")
	}
}

func handleNews(c *gin.Context, fetcher news.Provider) {
	ticker := c.Param("ticker")
	query := c.Query("q")
	summarize := c.DefaultQuery("summarize", "false") == "true"

	requestLog := log.With().Str("ticker", ticker).Logger()

	// Input validation.
	if !validTickerRegex.MatchString(ticker) {
		requestLog.Warn().Msg("Invalid ticker format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticker format."})
		return
	}

	// add a timeout to the context for the fetcher call
	fetchCtx, cancelFetch := context.WithTimeout(c, 10*time.Second)
	defer cancelFetch() // Important: ensure cancel is called to release resources

	articles, err := fetcher.GetNewsByTicker(fetchCtx, ticker)

	if err != nil {
		requestLog.Error().Err(err).Msg("Error processing news request")

		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out."})
		} else if errors.Is(err, context.Canceled) {
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "Request canceled."})
		} else if errors.Is(err, apperrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No news found for the specified ticker."})
		} else if errors.Is(err, apperrors.ErrServiceUnavailable) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "External service unavailable."})
		} else if errors.Is(err, apperrors.ErrInternal) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown error."})
		}
		return
	}

	if summarize {
		var allArticles string
		for _, a := range articles {
			allArticles += a.Title + ":" + a.Summary + "\n"
		}

		if allArticles == "" {
			requestLog.Warn().Msg("No article content to summarize.")
			c.JSON(http.StatusOK, gin.H{"ticker": ticker, "summary": ""})
			return
		}

		aiCtx, cancelAI := context.WithTimeout(c, 10*time.Second)
		defer cancelAI() // Important: ensure cancel is called to release resources

		summary, err := ai.SummarizeArticles(aiCtx, allArticles)
		if err != nil {
			requestLog.Error().Err(err).Msg("Failed to generate AI summary")

			if errors.Is(err, context.DeadlineExceeded) {
				c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out."})
			} else if errors.Is(err, context.Canceled) {
				c.JSON(http.StatusRequestTimeout, gin.H{"error": "Request canceled."})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}

			return
		}

		requestLog.Info().Msg("AI summary generated successfully")
		c.JSON(http.StatusOK, gin.H{"ticker": ticker, "summary": summary})
		return
	}

	if query != "" {
		articles = filter.FilterByQuery(articles, query)
		requestLog.Debug().Str("query", query).Int("result_count", len(articles)).Msg("Filtered articles by query")
	}

	requestLog.Info().Int("article_count", len(articles)).Msg("Successfully retrieved news articles")
	c.JSON(http.StatusOK, gin.H{"ticker": ticker, "news": articles})
}
