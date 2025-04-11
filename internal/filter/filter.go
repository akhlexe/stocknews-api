package filter

import (
	"strings"

	"github.com/akhlexe/stocknews-api/internal/news"
)

func FilterByQuery(articles []news.Article, query string) []news.Article {
	var result []news.Article

	query = strings.ToLower(query)
	for _, a := range articles {
		if strings.Contains(strings.ToLower(a.Title), query) ||
			strings.Contains(strings.ToLower(a.Summary), query) {
			result = append(result, a)
		}
	}

	return result
}
