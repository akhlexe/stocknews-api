package stocks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/akhlexe/stocknews-api/internal/ai"
)

type NewsArticle struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Summary     string   `json:"summary"`
	Image       string   `json:"banner_image"`
	PublishedAt string   `json:"time_published"`
	Source      string   `json:"source"`
	Sentiment   string   `json:"overall_sentiment_label"`
	Tickers     []string `json:"tickers"`
}

type apiResponse struct {
	Items string `json:"items"`
	Feed  []struct {
		Title       string `json:"title"`
		URL         string `json:"url"`
		Summary     string `json:"summary"`
		BannerImage string `json:"banner_image"`
		Time        string `json:"time_published"`
		Source      string `json:"source"`
		Sentiment   string `json:"overall_sentiment_label"`
		TickerData  []struct {
			Ticker string `json:"ticker"`
		} `json:"ticker_sentiment"`
	} `json:"feed"`
}

func GetNewsByTicker(ticker string) ([]NewsArticle, error) {
	apiKey := "YOUR_API_KEY"

	if apiKey == "" {
		return nil, fmt.Errorf("missing ALPHAVANTAGE_API_KEY environment variable")
	}

	endpoint := "https://www.alphavantage.co/query"
	params := url.Values{}
	params.Set("function", "NEWS_SENTIMENT")
	params.Set("tickers", ticker)
	params.Set("apikey", apiKey)

	fullUrl := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	resp, err := http.Get(fullUrl)
	if err != nil {
		return nil, fmt.Errorf("error resqueting AlphaVantage: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status code %d", resp.StatusCode)
	}

	var result apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding API response: %w", err)
	}

	var articles []NewsArticle

	for _, item := range result.Feed {
		var tickers []string
		for _, t := range item.TickerData {
			tickers = append(tickers, t.Ticker)
		}

		articles = append(articles, NewsArticle{
			Title:       item.Title,
			URL:         item.URL,
			Summary:     item.Summary,
			Image:       item.BannerImage,
			PublishedAt: item.Time,
			Source:      item.Source,
			Sentiment:   item.Sentiment,
			Tickers:     tickers,
		})
	}

	if len(articles) > 0 {
		first := articles[0]
		summaryPrompt := fmt.Sprintf("Summarize this news article in 1-2 lines: \n\nTitle: %s\n\n%s", first.Title, first.Summary)

		summary, err := ai.GenerateSummary(summaryPrompt)
		if err != nil {
			log.Println("❌ Failed to generate AI summary:", err)
			return nil, fmt.Errorf("error generating summary: %w", err)
		} else {
			log.Println("✅ AI summary generated successfully")
			articles[0].Summary = summary
		}
	}

	return articles, nil
}

func FilterByQuery(articles []NewsArticle, query string) []NewsArticle {
	var result []NewsArticle

	query = strings.ToLower(query)
	for _, a := range articles {
		if strings.Contains(strings.ToLower(a.Title), query) ||
			strings.Contains(strings.ToLower(a.Summary), query) {
			result = append(result, a)
		}
	}

	return result
}
