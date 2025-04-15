package news

import "context"

type Provider interface {
	GetNewsByTicker(ctx context.Context, ticker string) ([]Article, error)
}
