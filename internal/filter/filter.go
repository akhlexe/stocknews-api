package filter

import (
	"strings"

	"github.com/akhlexe/stocknews-api/internal/models"
)

func FilterByQuery(articles []models.Article, query string) []models.Article {
	var result []models.Article

	query = strings.ToLower(query)
	for _, a := range articles {
		if strings.Contains(strings.ToLower(a.Title), query) ||
			strings.Contains(strings.ToLower(a.Summary), query) {
			result = append(result, a)
		}
	}

	return result
}
