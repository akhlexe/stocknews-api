package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akhlexe/stocknews-api/internal/apperrors"
	"github.com/akhlexe/stocknews-api/internal/news"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestRouter(fetcher news.Provider) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/news/:ticker", func(c *gin.Context) {
		handleNews(c, fetcher)
	})

	return router
}

func TestHandleNews(t *testing.T) {
	samplerArticles := []news.Article{
		{Title: "Test Stock Up", Summary: "Good news for TEST", Tickers: []string{"TEST"}},
		{Title: "TEST Results", Summary: "Quarterly results analysis", Tickers: []string{"TEST"}},
	}

	// The handler returns tickers as []interface{} not []string, so adjust accordingly
	expectedSampleArticleBody := []interface{}{
		map[string]interface{}{
			"title":                   "Test Stock Up",
			"url":                     "",
			"summary":                 "Good news for TEST",
			"banner_image":            "",
			"time_published":          "",
			"source":                  "",
			"overall_sentiment_label": "",
			"tickers":                 []interface{}{"TEST"},
		},
		map[string]interface{}{
			"title":                   "TEST Results",
			"url":                     "",
			"summary":                 "Quarterly results analysis",
			"banner_image":            "",
			"time_published":          "",
			"source":                  "",
			"overall_sentiment_label": "",
			"tickers":                 []interface{}{"TEST"},
		},
	}

	testCases := []struct {
		name           string
		tickerParam    string
		queryParams    map[string]string
		mockSetup      func(*MockNewsProvider)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:        "Success - Valid Ticker No Filter",
			tickerParam: "TEST",
			queryParams: nil,
			mockSetup: func(mf *MockNewsProvider) {
				mf.On("GetNewsByTicker",
					mock.AnythingOfType("*context.timerCtx"),
					"TEST",
				).Return(samplerArticles, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"ticker": "TEST",
				"news":   expectedSampleArticleBody,
			},
		},
		{
			name:        "Success - Valid Ticker With Filter",
			tickerParam: "TEST",
			queryParams: map[string]string{"q": "Stock Up"},
			mockSetup: func(mf *MockNewsProvider) {
				mf.On("GetNewsByTicker",
					mock.AnythingOfType("*context.timerCtx"),
					"TEST",
				).Return(samplerArticles, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"ticker": "TEST",
				"news":   []interface{}{expectedSampleArticleBody[0]},
			},
		},
		{
			name:        "Error - Invalid Ticker Format",
			tickerParam: "TEST!",
			queryParams: nil,
			mockSetup: func(mf *MockNewsProvider) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid ticker format.",
			},
		},
		{
			name:        "Error - Ticker Not Found (Fetcher error)",
			tickerParam: "UNKNOWN",
			queryParams: nil,
			mockSetup: func(mf *MockNewsProvider) {
				mf.On("GetNewsByTicker",
					mock.AnythingOfType("*context.timerCtx"),
					"UNKNOWN",
				).Return(nil, apperrors.ErrNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "No news found for the specified ticker.",
			},
		},
		{
			name:        "Error - Service Unavailable (Fetcher error)",
			tickerParam: "TEST",
			queryParams: nil,
			mockSetup: func(mf *MockNewsProvider) {
				mf.On("GetNewsByTicker",
					mock.AnythingOfType("*context.timerCtx"),
					"TEST",
				).Return(nil, apperrors.ErrServiceUnavailable).Once()
			},
			expectedStatus: http.StatusServiceUnavailable,
			expectedBody: map[string]interface{}{
				"error": "External service unavailable.",
			},
		},
		{
			name:        "Error - Generic Fetcher Error",
			tickerParam: "TEST",
			queryParams: nil,
			mockSetup: func(mf *MockNewsProvider) {
				// Configure mock to return a generic error
				mf.On("GetNewsByTicker",
					mock.AnythingOfType("*context.timerCtx"),
					"TEST",
				).Return(nil, errors.New("something unexpected happened")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Unknown error.",
			},
		},
	}

	// -- Run Test Cases --
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Create Mock & Setup Expectations
			mockFetcher := new(MockNewsProvider)
			if tc.mockSetup != nil {
				tc.mockSetup(mockFetcher)
			}

			// 2. Setup Router with Mock
			router := setupTestRouter(mockFetcher)

			// 3. Create Request & Response Recorder
			w := httptest.NewRecorder()

			// Build URL with path param and optional query params
			urlPath := "/news/" + tc.tickerParam
			req, _ := http.NewRequest(http.MethodGet, urlPath, nil)
			q := req.URL.Query()
			if tc.queryParams != nil {
				for k, v := range tc.queryParams {
					q.Add(k, v)
				}
			}

			req.URL.RawQuery = q.Encode()

			// 4. Perform Request through Router
			router.ServeHTTP(w, req)

			// 5. Assert Status Code
			assert.Equal(t, tc.expectedStatus, w.Code, "HTTP status code mismatch")

			// 6. Assert Response Body (by decoding JSON)
			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			if !assert.NoError(t, err, "Failed to unmarshal response body JSON") {
				t.Logf("Response body: %s", w.Body.String())
				return
			}

			// For debugging, print the actual response vs expected
			if !assert.Equal(t, tc.expectedBody, responseBody, "Response body mismatch") {
				t.Logf("Expected: %#v", tc.expectedBody)
				t.Logf("Actual: %#v", responseBody)
			}

			// 7. Verify Mock Expectations were met
			mockFetcher.AssertExpectations(t)
		})
	}
}
