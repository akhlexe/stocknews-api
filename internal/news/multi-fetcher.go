package news

type MultiFetcher struct {
	Providers []Provider
}

func NewMultiFetcher(providers ...Provider) *MultiFetcher {
	return &MultiFetcher{
		Providers: providers,
	}
}

func (m *MultiFetcher) GetNewsByTicker(ticker string) ([]Article, error) {
	var allArticles []Article

	for _, provider := range m.Providers {
		articles, err := provider.GetNewsByTicker(ticker)
		if err != nil {
			return nil, err
		}
		allArticles = append(allArticles, articles...)
	}

	return allArticles, nil
}
