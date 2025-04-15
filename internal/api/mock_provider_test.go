package api

import (
	"context"

	"github.com/akhlexe/stocknews-api/internal/news"
	"github.com/stretchr/testify/mock"
)

type MockNewsProvider struct {
	mock.Mock
}

func (m *MockNewsProvider) GetNewsByTicker(ctx context.Context, ticker string) ([]news.Article, error) {
	args := m.Called(ctx, ticker)

	var articles []news.Article
	if args.Get(0) != nil {
		val, ok := args.Get(0).([]news.Article)
		if ok {
			articles = val
		}
	}

	return articles, args.Error(1)
}
