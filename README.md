# stocknews-api

# Stock News API

A RESTful API service for retrieving financial news about stocks with filtering and summarization capabilities.

## Features

- **Ticker-based News Retrieval**: Get the latest news for specific stock tickers
- **News Filtering**: Filter news articles by query terms
- **AI Summarization**: Generate concise summaries of stock news (when enabled)
- **Caching**: Built-in caching to reduce external API calls and improve performance
- **Multiple Data Sources**: Extensible architecture for multiple news providers

## Technologies

- Go (Golang) 1.20+
- Gin Web Framework
- AlphaVantage API for financial data
- Ollama for AI summarization (optional)
- Zerolog for structured logging

## Getting Started

### Prerequisites

- Go 1.20 or higher
- AlphaVantage API key (get one at [alphavantage.co](https://www.alphavantage.co/))
- Ollama (optional, for AI summarization)

### Installation

1. Clone the repository

```bash
git clone https://github.com/yourusername/stocknews-api.git
cd stocknews-api
```

2. Install dependencies

```bash
make deps
```

3. Create a `.env` file based on the example

```bash
cp .env.example .env
```

4. Edit the `.env` file and add your API keys

```
ALPHAVANTAGE_API_KEY=your_key_here
OLLAMA_URL=http://localhost:11434
```

5. Build and run the application

```bash
make run
```

## Usage

### API Endpoints

- **GET /health**: Check API health
- **GET /news/{ticker}**: Get news for a specific ticker
  - Query Parameters:
    - `q`: Filter news by text search
    - `summarize`: Set to "true" to get an AI-generated summary

### Examples

Retrieve news for Apple Inc:

```
GET /news/AAPL
```

Filter news containing "revenue":

```
GET /news/MSFT?q=revenue
```

Get an AI-generated summary of all news for Tesla:

```
GET /news/TSLA?summarize=true
```

## Development

This project includes a Makefile to simplify common development tasks:

| Command                            | Description                      |
| ---------------------------------- | -------------------------------- |
| `make run`                         | Run the application              |
| `make build`                       | Build the binary                 |
| `make test`                        | Run all tests                    |
| `make test-verbose`                | Run tests with verbose output    |
| `make test-file FILE=path/to/file` | Run tests in a specific file     |
| `make fmt`                         | Format code                      |
| `make tidy`                        | Run go mod tidy                  |
| `make coverage`                    | Generate test coverage report    |
| `make clean`                       | Clean up binaries and test cache |
| `make deps`                        | Install dependencies             |

## Architecture

The application follows a clean architecture approach:

- `cmd/server`: Application entry point
- `internal/api`: HTTP handlers and server configuration
- `internal/news`: News providers and article models
- `internal/cache`: In-memory caching functionality
- `internal/filter`: News article filtering logic
- `internal/ai`: AI summarization capabilities
- `internal/apperrors`: Application-specific error types

## License

[MIT License](LICENSE)
