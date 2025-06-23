package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"movie-discovery-app/configs"
	"movie-discovery-app/internal/models"
	"movie-discovery-app/internal/services"
)

func setupTestHandlers() *Handlers {
	config := &configs.Config{
		TMDB: configs.TMDBConfig{
			APIKey:  "test_key",
			BaseURL: "https://api.themoviedb.org/3",
		},
		OMDB: configs.OMDBConfig{
			APIKey:  "test_key",
			BaseURL: "http://www.omdbapi.com",
		},
		Cache: configs.CacheConfig{
			Duration: 30 * time.Minute,
		},
		Rate: configs.RateLimitConfig{
			RequestsPerMinute: 60,
		},
	}

	discoveryService := services.NewDiscoveryService(config)
	watchlistService := services.NewWatchlistService()
	recommendationService := services.NewRecommendationService(discoveryService, watchlistService)
	genreService := services.NewGenreService(config)

	return NewHandlers(discoveryService, watchlistService, recommendationService, genreService)
}

func TestHandlers_HealthCheck(t *testing.T) {
	handlers := setupTestHandlers()

	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HealthCheck)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response["status"])
	}

	if response["service"] != "movie-discovery-app" {
		t.Errorf("Expected service 'movie-discovery-app', got '%s'", response["service"])
	}
}

func TestHandlers_SearchMovies_InvalidQuery(t *testing.T) {
	handlers := setupTestHandlers()

	tests := []struct {
		name           string
		query          string
		expectedStatus int
	}{
		{
			name:           "Empty query",
			query:          "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "No query parameter",
			query:          "no_q_param",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var url string
			if tt.query == "no_q_param" {
				url = "/api/v1/search/movies"
			} else {
				url = "/api/v1/search/movies?q=" + tt.query
			}

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.SearchMovies)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestHandlers_AddToWatchlist(t *testing.T) {
	handlers := setupTestHandlers()

	watchlistItem := models.WatchlistItem{
		ID:         "123",
		Type:       "movie",
		Title:      "Test Movie",
		PosterPath: "/test.jpg",
		Watched:    false,
		Rating:     0,
	}

	jsonData, err := json.Marshal(watchlistItem)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/v1/watchlist", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.AddToWatchlist)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if response["status"] != "success" {
		t.Errorf("Expected status 'success', got '%s'", response["status"])
	}
}

func TestHandlers_AddToWatchlist_InvalidJSON(t *testing.T) {
	handlers := setupTestHandlers()

	invalidJSON := []byte(`{"invalid": json}`)

	req, err := http.NewRequest("POST", "/api/v1/watchlist", bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.AddToWatchlist)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestHandlers_AddToWatchlist_WrongMethod(t *testing.T) {
	handlers := setupTestHandlers()

	req, err := http.NewRequest("GET", "/api/v1/watchlist", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.AddToWatchlist)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestHandlers_GetWatchlist(t *testing.T) {
	handlers := setupTestHandlers()

	// First add an item to the watchlist
	watchlistItem := models.WatchlistItem{
		ID:         "123",
		Type:       "movie",
		Title:      "Test Movie",
		PosterPath: "/test.jpg",
		Watched:    false,
		Rating:     0,
	}

	// Add item first
	jsonData, _ := json.Marshal(watchlistItem)
	addReq, _ := http.NewRequest("POST", "/api/v1/watchlist", bytes.NewBuffer(jsonData))
	addReq.Header.Set("Content-Type", "application/json")
	addRR := httptest.NewRecorder()
	handlers.AddToWatchlist(addRR, addReq)

	// Now get the watchlist
	req, err := http.NewRequest("GET", "/api/v1/watchlist", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GetWatchlist)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var watchlist []models.WatchlistItem
	err = json.Unmarshal(rr.Body.Bytes(), &watchlist)
	if err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if len(watchlist) != 1 {
		t.Errorf("Expected 1 item in watchlist, got %d", len(watchlist))
	}

	if watchlist[0].ID != "123" {
		t.Errorf("Expected item ID '123', got '%s'", watchlist[0].ID)
	}
}

func TestHandlers_GetWatchlistStats(t *testing.T) {
	handlers := setupTestHandlers()

	req, err := http.NewRequest("GET", "/api/v1/watchlist/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GetWatchlistStats)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var stats map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &stats)
	if err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	// Check that all expected fields are present
	expectedFields := []string{"total_items", "watched_items", "unwatched_items", "movies", "tv_shows", "average_rating"}
	for _, field := range expectedFields {
		if _, exists := stats[field]; !exists {
			t.Errorf("Expected field '%s' in stats response", field)
		}
	}
}

func TestHandlers_CORS(t *testing.T) {
	handlers := setupTestHandlers()

	req, err := http.NewRequest("OPTIONS", "/api/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	
	// Create a test handler to wrap with CORS
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	
	corsHandler := handlers.EnableCORS(testHandler)
	corsHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check CORS headers
	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type, Authorization",
	}

	for header, expectedValue := range expectedHeaders {
		if value := rr.Header().Get(header); value != expectedValue {
			t.Errorf("Expected header %s to be '%s', got '%s'", header, expectedValue, value)
		}
	}
}

// Benchmark tests
func BenchmarkHandlers_HealthCheck(b *testing.B) {
	handlers := setupTestHandlers()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/api/v1/health", nil)
		rr := httptest.NewRecorder()
		handlers.HealthCheck(rr, req)
	}
}

func BenchmarkHandlers_GetWatchlist(b *testing.B) {
	handlers := setupTestHandlers()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/api/v1/watchlist", nil)
		rr := httptest.NewRecorder()
		handlers.GetWatchlist(rr, req)
	}
}
