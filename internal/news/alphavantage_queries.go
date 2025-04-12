package news

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Article struct {
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

func GetNewsByTicker(apiKey string, ticker string) ([]Article, error) {
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

	var articles []Article

	for _, item := range result.Feed {
		var tickers []string
		for _, t := range item.TickerData {
			tickers = append(tickers, t.Ticker)
		}

		articles = append(articles, Article{
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

	return articles, nil
}
