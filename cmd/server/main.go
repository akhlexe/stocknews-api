package main

import (
	server "github.com/akhlexe/stocknews-api/internal/api"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	server.Run()
}
