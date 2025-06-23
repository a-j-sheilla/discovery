package services

import (
	"testing"
	"time"

	"movie-discovery-app/configs"
)

func TestDiscoveryService_ValidateSearchQuery(t *testing.T) {
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

	service := NewDiscoveryService(config)

	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			name:    "Valid query",
			query:   "Inception",
			wantErr: false,
		},
		{
			name:    "Empty query",
			query:   "",
			wantErr: true,
		},
		{
			name:    "Whitespace only query",
			query:   "   ",
			wantErr: true,
		},
		{
			name:    "Very long query",
			query:   string(make([]byte, 101)), // 101 characters
			wantErr: true,
		},
		{
			name:    "Valid short query",
			query:   "a",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateSearchQuery(tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSearchQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDiscoveryService_ValidatePage(t *testing.T) {
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

	service := NewDiscoveryService(config)

	tests := []struct {
		name    string
		page    int
		wantErr bool
	}{
		{
			name:    "Valid page 1",
			page:    1,
			wantErr: false,
		},
		{
			name:    "Valid page 100",
			page:    100,
			wantErr: false,
		},
		{
			name:    "Valid page 1000",
			page:    1000,
			wantErr: false,
		},
		{
			name:    "Invalid page 0",
			page:    0,
			wantErr: true,
		},
		{
			name:    "Invalid negative page",
			page:    -1,
			wantErr: true,
		},
		{
			name:    "Invalid page too high",
			page:    1001,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidatePage(tt.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCache_SetAndGet(t *testing.T) {
	cache := &Cache{
		data: make(map[string]CacheItem),
	}

	// Test setting and getting data
	testData := "test value"
	cache.Set("test_key", testData, 1*time.Hour)

	retrieved := cache.Get("test_key")
	if retrieved != testData {
		t.Errorf("Expected %v, got %v", testData, retrieved)
	}

	// Test getting non-existent key
	nonExistent := cache.Get("non_existent")
	if nonExistent != nil {
		t.Errorf("Expected nil for non-existent key, got %v", nonExistent)
	}

	// Test expired data
	cache.Set("expired_key", "expired_value", -1*time.Hour) // Already expired
	expired := cache.Get("expired_key")
	if expired != nil {
		t.Errorf("Expected nil for expired key, got %v", expired)
	}
}

func TestRateLimiter_Wait(t *testing.T) {
	rateLimiter := &RateLimiter{
		requests: make(chan time.Time, 2), // Small limit for testing
		limit:    2,
		window:   time.Minute,
	}

	// First request should succeed
	err := rateLimiter.Wait()
	if err != nil {
		t.Errorf("First request should succeed, got error: %v", err)
	}

	// Second request should succeed
	err = rateLimiter.Wait()
	if err != nil {
		t.Errorf("Second request should succeed, got error: %v", err)
	}

	// Third request should fail (rate limit exceeded)
	err = rateLimiter.Wait()
	if err == nil {
		t.Error("Third request should fail due to rate limit")
	}
}

// Benchmark tests
func BenchmarkCache_Set(b *testing.B) {
	cache := &Cache{
		data: make(map[string]CacheItem),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("test_key", "test_value", 1*time.Hour)
	}
}

func BenchmarkCache_Get(b *testing.B) {
	cache := &Cache{
		data: make(map[string]CacheItem),
	}
	cache.Set("test_key", "test_value", 1*time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("test_key")
	}
}

// Example test for integration testing (would require actual API keys)
func ExampleDiscoveryService_SearchMovies() {
	// This is an example of how to use the SearchMovies method
	// In a real test, you would mock the HTTP client
	
	config := &configs.Config{
		TMDB: configs.TMDBConfig{
			APIKey:  "your_api_key",
			BaseURL: "https://api.themoviedb.org/3",
		},
		OMDB: configs.OMDBConfig{
			APIKey:  "your_api_key",
			BaseURL: "http://www.omdbapi.com",
		},
		Cache: configs.CacheConfig{
			Duration: 30 * time.Minute,
		},
		Rate: configs.RateLimitConfig{
			RequestsPerMinute: 60,
		},
	}

	service := NewDiscoveryService(config)
	
	// This would make actual API calls in a real scenario
	// results, err := service.SearchMovies("Inception", 1)
	// if err != nil {
	//     log.Fatal(err)
	// }
	// fmt.Printf("Found %d results\n", len(results.Results))
	
	_ = service // Use the service to avoid unused variable error
}
