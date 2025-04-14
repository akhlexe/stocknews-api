package filter

import (
	"testing"

	"github.com/akhlexe/stocknews-api/internal/news"
)

func TestFilterByQuery(t *testing.T) {
	articles := []news.Article{
		{Title: "Apple Inc. Announces New Product", Summary: "Apple is launching a new iPhone."},
		{Title: "Google's Stock Surges", Summary: "Google's stock price has increased."},
		{Title: "Microsoft's Cloud Services", Summary: "Microsoft is expanding its cloud offerings."},
		{Title: "Tesla's New Car", Summary: "Tesla releases a new electric car."},
	}

	testCases := []struct {
		name     string
		query    string
		expected int
	}{
		{"Filter by Apple", "apple", 1},
		{"Filter by Google", "google", 1},
		{"Filter by Microsoft", "microsoft", 1},
		{"Filter by Tesla", "tesla", 1},
		{"Filter by iPhone", "iphone", 1},
		{"Filter by Stock", "stock", 1},
		{"Filter by Cloud", "cloud", 1},
		{"Filter by Car", "car", 1},
		{"Filter by Nonexistent", "nonexistent", 0},
		{"Filter by Empty", "", 4},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filtered := FilterByQuery(articles, tc.query)
			if len(filtered) != tc.expected {
				t.Errorf("Expected %d articles, got %d", tc.expected, len(filtered))
			}
		})
	}

}
