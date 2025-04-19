package main

import (
	"fmt"
	"os"
	"time"

	"github.com/akhlexe/stocknews-api/internal/api"
	"github.com/akhlexe/stocknews-api/internal/cache"
	"github.com/akhlexe/stocknews-api/internal/news"
	"github.com/akhlexe/stocknews-api/internal/storage"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
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

	// Initializa the Postgres Storage
	log.Info().Msg("Initializing storage")
	postgresStorage, err := CreatePostgresStorage()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create storage")
	}
	defer postgresStorage.Close()

	cache := cache.NewPersistentCache(postgresStorage, 10*time.Minute)
	apiKey := os.Getenv("ALPHAVANTAGE_API_KEY")
	fetcher := news.NewAlphaVantageFetcher(apiKey, cache)
	multiFetcher := news.NewMultiFetcher(fetcher)

	server := api.NewServer(multiFetcher)
	server.Run()
}

func CreatePostgresStorage() (*storage.PostgresStorage, error) {
	// Get PostgreSQL connection details from environment variables
	host := getEnvOrDefault("POSTGRES_HOST", "localhost")
	port := getEnvOrDefault("POSTGRES_PORT", "5434")
	user := getEnvOrDefault("POSTGRES_USER", "postgres")
	password := getEnvOrDefault("POSTGRES_PASSWORD", "postgres")
	dbName := getEnvOrDefault("POSTGRES_DB", "stocknews")

	log.Info().Msgf("%s, %s, %s, %s, %s", host, port, user, password, dbName)

	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	log.Info().Msgf("Connecting to PostgreSQL database: %s", connStr)

	log.Debug().
		Str("host", host).
		Str("port", port).
		Str("user", user).
		Str("password", password).
		Str("dbName", dbName).
		Msg("Connecting to PostgreSQL database")

	return storage.NewPostgresStorage(connStr)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
