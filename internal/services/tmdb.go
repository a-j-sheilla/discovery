package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"movie-discovery-app/configs"
	"movie-discovery-app/internal/models"
)

// TMDBClient handles TMDB API interactions
type TMDBClient struct {
	config     *configs.TMDBConfig
	httpClient *http.Client
	cache      *Cache
	rateLimiter *RateLimiter
}

// Cache represents a simple in-memory cache
type Cache struct {
	data map[string]CacheItem
	mu   sync.RWMutex
}

// CacheItem represents a cached item
type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

// RateLimiter implements rate limiting
type RateLimiter struct {
	requests chan time.Time
	limit    int
	window   time.Duration
}

// NewTMDBClient creates a new TMDB API client
func NewTMDBClient(config *configs.TMDBConfig, cacheConfig *configs.CacheConfig, rateConfig *configs.RateLimitConfig) *TMDBClient {
	return &TMDBClient{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cache: &Cache{
			data: make(map[string]CacheItem),
		},
		rateLimiter: &RateLimiter{
			requests: make(chan time.Time, rateConfig.RequestsPerMinute),
			limit:    rateConfig.RequestsPerMinute,
			window:   time.Minute,
		},
	}
}

// SearchMovies searches for movies using TMDB API
func (c *TMDBClient) SearchMovies(query string, page int) (*models.SearchResult, error) {
	cacheKey := fmt.Sprintf("search_movies_%s_%d", query, page)
	
	// Check cache first
	if cached := c.cache.Get(cacheKey); cached != nil {
		if result, ok := cached.(*models.SearchResult); ok {
			return result, nil
		}
	}

	// Rate limiting
	if err := c.rateLimiter.Wait(); err != nil {
		return nil, err
	}

	// Build URL
	params := url.Values{}
	params.Add("api_key", c.config.APIKey)
	params.Add("query", query)
	params.Add("page", strconv.Itoa(page))

	url := fmt.Sprintf("%s/search/movie?%s", c.config.BaseURL, params.Encode())

	// Make request
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result models.SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Cache the result
	c.cache.Set(cacheKey, &result, 30*time.Minute)

	return &result, nil
}

// SearchTVShows searches for TV shows using TMDB API
func (c *TMDBClient) SearchTVShows(query string, page int) (*models.SearchResult, error) {
	cacheKey := fmt.Sprintf("search_tv_%s_%d", query, page)
	
	// Check cache first
	if cached := c.cache.Get(cacheKey); cached != nil {
		if result, ok := cached.(*models.SearchResult); ok {
			return result, nil
		}
	}

	// Rate limiting
	if err := c.rateLimiter.Wait(); err != nil {
		return nil, err
	}

	// Build URL
	params := url.Values{}
	params.Add("api_key", c.config.APIKey)
	params.Add("query", query)
	params.Add("page", strconv.Itoa(page))

	url := fmt.Sprintf("%s/search/tv?%s", c.config.BaseURL, params.Encode())

	// Make request
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to search TV shows: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result models.SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Cache the result
	c.cache.Set(cacheKey, &result, 30*time.Minute)

	return &result, nil
}

// GetMovieDetails gets detailed movie information
func (c *TMDBClient) GetMovieDetails(movieID int) (*models.Movie, error) {
	cacheKey := fmt.Sprintf("movie_details_%d", movieID)
	
	// Check cache first
	if cached := c.cache.Get(cacheKey); cached != nil {
		if movie, ok := cached.(*models.Movie); ok {
			return movie, nil
		}
	}

	// Rate limiting
	if err := c.rateLimiter.Wait(); err != nil {
		return nil, err
	}

	// Build URL
	params := url.Values{}
	params.Add("api_key", c.config.APIKey)

	url := fmt.Sprintf("%s/movie/%d?%s", c.config.BaseURL, movieID, params.Encode())

	// Make request
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var movie models.Movie
	if err := json.Unmarshal(body, &movie); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Cache the result
	c.cache.Set(cacheKey, &movie, 30*time.Minute)

	return &movie, nil
}

// GetTrendingMovies gets trending movies
func (c *TMDBClient) GetTrendingMovies(timeWindow string, page int) (*models.TrendingResponse, error) {
	cacheKey := fmt.Sprintf("trending_movies_%s_%d", timeWindow, page)
	
	// Check cache first
	if cached := c.cache.Get(cacheKey); cached != nil {
		if result, ok := cached.(*models.TrendingResponse); ok {
			return result, nil
		}
	}

	// Rate limiting
	if err := c.rateLimiter.Wait(); err != nil {
		return nil, err
	}

	// Build URL
	params := url.Values{}
	params.Add("api_key", c.config.APIKey)
	params.Add("page", strconv.Itoa(page))

	url := fmt.Sprintf("%s/trending/movie/%s?%s", c.config.BaseURL, timeWindow, params.Encode())

	// Make request
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending movies: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result models.TrendingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Cache the result
	c.cache.Set(cacheKey, &result, 30*time.Minute)

	return &result, nil
}

// Cache methods
func (c *Cache) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	item, exists := c.data[key]
	if !exists || time.Now().After(item.ExpiresAt) {
		return nil
	}
	
	return item.Data
}

func (c *Cache) Set(key string, data interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.data[key] = CacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(duration),
	}
}

// RateLimiter methods
func (rl *RateLimiter) Wait() error {
	now := time.Now()
	
	// Remove old requests outside the window
	for {
		select {
		case t := <-rl.requests:
			if now.Sub(t) < rl.window {
				// Put it back if it's still within the window
				select {
				case rl.requests <- t:
				default:
					return fmt.Errorf("rate limit exceeded")
				}
				goto checkLimit
			}
		default:
			goto checkLimit
		}
	}
	
checkLimit:
	// Check if we can make a new request
	if len(rl.requests) >= rl.limit {
		return fmt.Errorf("rate limit exceeded")
	}
	
	// Add current request
	select {
	case rl.requests <- now:
		return nil
	default:
		return fmt.Errorf("rate limit exceeded")
	}
}
