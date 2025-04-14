package main

import (
	"os"
	"time"

	"github.com/akhlexe/stocknews-api/internal/api"
	"github.com/akhlexe/stocknews-api/internal/cache"
	"github.com/akhlexe/stocknews-api/internal/news"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	apiKey := os.Getenv("ALPHAVANTAGE_API_KEY")
	cache := cache.NewCache(10 * time.Minute)
	fetcher := news.NewAlphaVantageFetcher(apiKey, cache)
	multiFetcher := news.NewMultiFetcher(fetcher)

	server := api.NewServer(multiFetcher)
	server.Run()
}
