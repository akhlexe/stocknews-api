package news

type Provider interface {
	GetNewsByTicker(ticker string) ([]Article, error)
}
