package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"movie-discovery-app/configs"
)

// OMDBClient handles OMDB API interactions
type OMDBClient struct {
	config      *configs.OMDBConfig
	httpClient  *http.Client
	cache       *Cache
	rateLimiter *RateLimiter
}

// OMDBResponse represents the response from OMDB API
type OMDBResponse struct {
	Title      string `json:"Title"`
	Year       string `json:"Year"`
	Rated      string `json:"Rated"`
	Released   string `json:"Released"`
	Runtime    string `json:"Runtime"`
	Genre      string `json:"Genre"`
	Director   string `json:"Director"`
	Writer     string `json:"Writer"`
	Actors     string `json:"Actors"`
	Plot       string `json:"Plot"`
	Language   string `json:"Language"`
	Country    string `json:"Country"`
	Awards     string `json:"Awards"`
	Poster     string `json:"Poster"`
	Ratings    []struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	} `json:"Ratings"`
	Metascore    string `json:"Metascore"`
	IMDBRating   string `json:"imdbRating"`
	IMDBVotes    string `json:"imdbVotes"`
	IMDBID       string `json:"imdbID"`
	Type         string `json:"Type"`
	TotalSeasons string `json:"totalSeasons,omitempty"`
	Response     string `json:"Response"`
	Error        string `json:"Error,omitempty"`
}

// NewOMDBClient creates a new OMDB API client
func NewOMDBClient(config *configs.OMDBConfig, cacheConfig *configs.CacheConfig, rateConfig *configs.RateLimitConfig) *OMDBClient {
	return &OMDBClient{
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

// GetMovieByTitle gets movie details from OMDB by title
func (c *OMDBClient) GetMovieByTitle(title string, year string) (*OMDBResponse, error) {
	cacheKey := fmt.Sprintf("omdb_title_%s_%s", title, year)
	
	// Check cache first
	if cached := c.cache.Get(cacheKey); cached != nil {
		if response, ok := cached.(*OMDBResponse); ok {
			return response, nil
		}
	}

	// Rate limiting
	if err := c.rateLimiter.Wait(); err != nil {
		return nil, err
	}

	// Build URL
	params := url.Values{}
	params.Add("apikey", c.config.APIKey)
	params.Add("t", title)
	params.Add("type", "movie")
	params.Add("plot", "full")
	if year != "" {
		params.Add("y", year)
	}

	requestURL := fmt.Sprintf("%s?%s", c.config.BaseURL, params.Encode())

	// Make request
	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie from OMDB: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var omdbResp OMDBResponse
	if err := json.Unmarshal(body, &omdbResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check if the response contains an error
	if omdbResp.Response == "False" {
		return nil, fmt.Errorf("OMDB error: %s", omdbResp.Error)
	}

	// Cache the result
	c.cache.Set(cacheKey, &omdbResp, 30*time.Minute)

	return &omdbResp, nil
}

// GetMovieByIMDBID gets movie details from OMDB by IMDB ID
func (c *OMDBClient) GetMovieByIMDBID(imdbID string) (*OMDBResponse, error) {
	cacheKey := fmt.Sprintf("omdb_imdb_%s", imdbID)
	
	// Check cache first
	if cached := c.cache.Get(cacheKey); cached != nil {
		if response, ok := cached.(*OMDBResponse); ok {
			return response, nil
		}
	}

	// Rate limiting
	if err := c.rateLimiter.Wait(); err != nil {
		return nil, err
	}

	// Build URL
	params := url.Values{}
	params.Add("apikey", c.config.APIKey)
	params.Add("i", imdbID)
	params.Add("plot", "full")

	requestURL := fmt.Sprintf("%s?%s", c.config.BaseURL, params.Encode())

	// Make request
	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie from OMDB: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var omdbResp OMDBResponse
	if err := json.Unmarshal(body, &omdbResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check if the response contains an error
	if omdbResp.Response == "False" {
		return nil, fmt.Errorf("OMDB error: %s", omdbResp.Error)
	}

	// Cache the result
	c.cache.Set(cacheKey, &omdbResp, 30*time.Minute)

	return &omdbResp, nil
}

// GetTVShowByTitle gets TV show details from OMDB by title
func (c *OMDBClient) GetTVShowByTitle(title string, year string) (*OMDBResponse, error) {
	cacheKey := fmt.Sprintf("omdb_tv_%s_%s", title, year)
	
	// Check cache first
	if cached := c.cache.Get(cacheKey); cached != nil {
		if response, ok := cached.(*OMDBResponse); ok {
			return response, nil
		}
	}

	// Rate limiting
	if err := c.rateLimiter.Wait(); err != nil {
		return nil, err
	}

	// Build URL
	params := url.Values{}
	params.Add("apikey", c.config.APIKey)
	params.Add("t", title)
	params.Add("type", "series")
	params.Add("plot", "full")
	if year != "" {
		params.Add("y", year)
	}

	requestURL := fmt.Sprintf("%s?%s", c.config.BaseURL, params.Encode())

	// Make request
	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get TV show from OMDB: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var omdbResp OMDBResponse
	if err := json.Unmarshal(body, &omdbResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check if the response contains an error
	if omdbResp.Response == "False" {
		return nil, fmt.Errorf("OMDB error: %s", omdbResp.Error)
	}

	// Cache the result
	c.cache.Set(cacheKey, &omdbResp, 30*time.Minute)

	return &omdbResp, nil
}

// SearchMovies searches for movies by title
func (c *OMDBClient) SearchMovies(title string, page int) (*OMDBSearchResponse, error) {
	cacheKey := fmt.Sprintf("omdb_search_%s_%d", title, page)
	
	// Check cache first
	if cached := c.cache.Get(cacheKey); cached != nil {
		if response, ok := cached.(*OMDBSearchResponse); ok {
			return response, nil
		}
	}

	// Rate limiting
	if err := c.rateLimiter.Wait(); err != nil {
		return nil, err
	}

	// Build URL
	params := url.Values{}
	params.Add("apikey", c.config.APIKey)
	params.Add("s", title)
	params.Add("type", "movie")
	if page > 1 {
		params.Add("page", fmt.Sprintf("%d", page))
	}

	requestURL := fmt.Sprintf("%s?%s", c.config.BaseURL, params.Encode())

	// Make request
	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies in OMDB: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var searchResp OMDBSearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check if the response contains an error
	if searchResp.Response == "False" {
		return nil, fmt.Errorf("OMDB error: %s", searchResp.Error)
	}

	// Cache the result
	c.cache.Set(cacheKey, &searchResp, 30*time.Minute)

	return &searchResp, nil
}

// OMDBSearchResponse represents search results from OMDB
type OMDBSearchResponse struct {
	Search []struct {
		Title  string `json:"Title"`
		Year   string `json:"Year"`
		IMDBID string `json:"imdbID"`
		Type   string `json:"Type"`
		Poster string `json:"Poster"`
	} `json:"Search"`
	TotalResults string `json:"totalResults"`
	Response     string `json:"Response"`
	Error        string `json:"Error,omitempty"`
}

// GetRottenTomatoesRating extracts Rotten Tomatoes rating from ratings array
func (r *OMDBResponse) GetRottenTomatoesRating() string {
	for _, rating := range r.Ratings {
		if rating.Source == "Rotten Tomatoes" {
			return rating.Value
		}
	}
	return ""
}
