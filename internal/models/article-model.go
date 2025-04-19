package models

import "encoding/json"

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

func MarshalArticles(articles []Article) ([]byte, error) {
	return json.Marshal(articles)
}

func UnmarshalArticles(data []byte) ([]Article, error) {
	var articles []Article
	err := json.Unmarshal(data, &articles)
	return articles, err
}
