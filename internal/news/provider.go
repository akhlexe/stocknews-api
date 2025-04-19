package news

import (
	"context"

	"github.com/akhlexe/stocknews-api/internal/models"
)

type Provider interface {
	GetNewsByTicker(ctx context.Context, ticker string) ([]models.Article, error)
}
