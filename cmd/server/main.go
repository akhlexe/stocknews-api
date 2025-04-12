package main

import (
	"os"
	"time"

	"github.com/akhlexe/stocknews-api/internal/api"
	"github.com/akhlexe/stocknews-api/internal/cache"
	"github.com/akhlexe/stocknews-api/internal/news"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	apiKey := os.Getenv("ALPHAVANTAGE_API_KEY")
	cache := cache.NewCache(10 * time.Minute)
	fetcher := news.NewAlphaVantageFetcher(apiKey, cache)
	multiFetcher := news.NewMultiFetcher(fetcher)

	server := api.NewServer(multiFetcher)
	server.Run()
}
