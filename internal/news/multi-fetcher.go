package news

import (
	"context"

	"github.com/akhlexe/stocknews-api/internal/models"
)

type MultiFetcher struct {
	Providers []Provider
}

func NewMultiFetcher(providers ...Provider) *MultiFetcher {
	return &MultiFetcher{
		Providers: providers,
	}
}

func (m *MultiFetcher) GetNewsByTicker(ctx context.Context, ticker string) ([]models.Article, error) {
	var allArticles []models.Article

	for _, provider := range m.Providers {
		articles, err := provider.GetNewsByTicker(ctx, ticker)
		if err != nil {
			return nil, err
		}
		allArticles = append(allArticles, articles...)
	}

	return allArticles, nil
}
