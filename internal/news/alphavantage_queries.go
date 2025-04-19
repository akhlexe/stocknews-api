package news

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akhlexe/stocknews-api/internal/apperrors"
	"github.com/akhlexe/stocknews-api/internal/models"
	"github.com/rs/zerolog/log"
)

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

func GetNewsByTicker(ctx context.Context, apiKey string, ticker string) ([]models.Article, error) {
	if apiKey == "" {
		log.Error().Msg("Missing ALPHAVANTAGE_API_KEY environment variable")
		return nil, fmt.Errorf("%w: missing ALPHAVANTAGE_API_KEY environment variable", apperrors.ErrConfiguration)
	}

	endpoint := "https://www.alphavantage.co/query"
	params := url.Values{}
	params.Set("function", "NEWS_SENTIMENT")
	params.Set("tickers", ticker)
	params.Set("apikey", apiKey)

	fullUrl := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullUrl, nil)
	if err != nil {
		log.Error().Err(err).Str("url", fullUrl).Msg("Error creating AlphaVantage request")
		return nil, fmt.Errorf("%w: error creating AlphaVantage request: %v", apperrors.ErrServiceUnavailable, err)
	}

	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {

		if ctx.Err() != nil {
			log.Warn().Err(ctx.Err()).Str("url", fullUrl).Msg("Context error requesting AlphaVantage")
			return nil, ctx.Err()
		}

		log.Error().Err(err).Str("url", fullUrl).Msg("Error requesting AlphaVantage")
		return nil, fmt.Errorf("%w: error resqueting AlphaVantage: %v", apperrors.ErrServiceUnavailable, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().
			Int("status_code", resp.StatusCode).
			Str("ticker", ticker).
			Msg("Error requesting AlphaVantage")

		return nil, fmt.Errorf("%w: Alphavantage API error: status code %d", apperrors.ErrServiceUnavailable, resp.StatusCode)
	}

	var result apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error().Err(err).Str("ticker", ticker).Msg("Error decoding AlphaVantage API response")
		return nil, fmt.Errorf("%w: error decoding API response: %v", apperrors.ErrInternal, err)
	}

	if len(result.Feed) == 0 {
		log.Warn().Str("ticker", ticker).Msg("No news articles found for the given ticker")
		return nil, apperrors.ErrNotFound
	}

	var articles []models.Article

	for _, item := range result.Feed {
		var tickers []string
		for _, t := range item.TickerData {
			tickers = append(tickers, t.Ticker)
		}

		articles = append(articles, models.Article{
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
