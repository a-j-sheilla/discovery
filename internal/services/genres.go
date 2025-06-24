package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"movie-discovery-app/configs"
	"movie-discovery-app/internal/models"
)

// GenreService handles genre-related operations
type GenreService struct {
	tmdbClient *TMDBClient
}

// NewGenreService creates a new genre service
func NewGenreService(config *configs.Config) *GenreService {
	tmdbClient := NewTMDBClient(&config.TMDB, &config.Cache, &config.Rate)
	return &GenreService{
		tmdbClient: tmdbClient,
	}
}

// GetMovieGenres gets all available movie genres
func (s *GenreService) GetMovieGenres() ([]models.Genre, error) {
	cacheKey := "movie_genres"

	// Check cache first
	if cached := s.tmdbClient.cache.Get(cacheKey); cached != nil {
		if genres, ok := cached.([]models.Genre); ok {
			return genres, nil
		}
	}

	// Rate limiting
	if err := s.tmdbClient.rateLimiter.Wait(); err != nil {
		return nil, err
	}

	// Build URL
	params := url.Values{}
	params.Add("api_key", s.tmdbClient.config.APIKey)

	requestURL := fmt.Sprintf("%s/genre/movie/list?%s", s.tmdbClient.config.BaseURL, params.Encode())

	// Make request
	resp, err := s.tmdbClient.httpClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie genres: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response struct {
		Genres []models.Genre `json:"genres"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Cache the result
	s.tmdbClient.cache.Set(cacheKey, response.Genres, 24*60*60*1000) // Cache for 24 hours

	return response.Genres, nil
}

// GetTVGenres gets all available TV show genres
func (s *GenreService) GetTVGenres() ([]models.Genre, error) {
	cacheKey := "tv_genres"

	// Check cache first
	if cached := s.tmdbClient.cache.Get(cacheKey); cached != nil {
		if genres, ok := cached.([]models.Genre); ok {
			return genres, nil
		}
	}

	// Rate limiting
	if err := s.tmdbClient.rateLimiter.Wait(); err != nil {
		return nil, err
	}

	// Build URL
	params := url.Values{}
	params.Add("api_key", s.tmdbClient.config.APIKey)

	requestURL := fmt.Sprintf("%s/genre/tv/list?%s", s.tmdbClient.config.BaseURL, params.Encode())

	// Make request
	resp, err := s.tmdbClient.httpClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get TV genres: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response struct {
		Genres []models.Genre `json:"genres"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Cache the result
	s.tmdbClient.cache.Set(cacheKey, response.Genres, 24*60*60*1000) // Cache for 24 hours

	return response.Genres, nil
}

// DiscoverMoviesByGenre discovers movies by genre with additional filters
func (s *GenreService) DiscoverMoviesByGenre(genreID int, page int, filters DiscoveryFilters) (*models.SearchResult, error) {
	cacheKey := fmt.Sprintf("discover_movies_genre_%d_page_%d", genreID, page)

	// Check cache first
	if cached := s.tmdbClient.cache.Get(cacheKey); cached != nil {
		if result, ok := cached.(*models.SearchResult); ok {
			return result, nil
		}
	}

	// Rate limiting
	if err := s.tmdbClient.rateLimiter.Wait(); err != nil {
		return nil, err
	}

	// Build URL with filters
	params := url.Values{}
	params.Add("api_key", s.tmdbClient.config.APIKey)
	params.Add("with_genres", strconv.Itoa(genreID))
	params.Add("page", strconv.Itoa(page))
	params.Add("sort_by", filters.SortBy)

	if filters.MinRating > 0 {
		params.Add("vote_average.gte", fmt.Sprintf("%.1f", filters.MinRating))
	}
	if filters.MaxRating > 0 {
		params.Add("vote_average.lte", fmt.Sprintf("%.1f", filters.MaxRating))
	}
	if filters.MinYear > 0 {
		params.Add("primary_release_date.gte", fmt.Sprintf("%d-01-01", filters.MinYear))
	}
	if filters.MaxYear > 0 {
		params.Add("primary_release_date.lte", fmt.Sprintf("%d-12-31", filters.MaxYear))
	}

	requestURL := fmt.Sprintf("%s/discover/movie?%s", s.tmdbClient.config.BaseURL, params.Encode())

	// Make request
	resp, err := s.tmdbClient.httpClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to discover movies: %w", err)
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
	s.tmdbClient.cache.Set(cacheKey, &result, 30*60*1000) // Cache for 30 minutes

	return &result, nil
}

// DiscoverTVShowsByGenre discovers TV shows by genre with additional filters
func (s *GenreService) DiscoverTVShowsByGenre(genreID int, page int, filters DiscoveryFilters) (*models.SearchResult, error) {
	cacheKey := fmt.Sprintf("discover_tv_genre_%d_page_%d", genreID, page)

	// Check cache first
	if cached := s.tmdbClient.cache.Get(cacheKey); cached != nil {
		if result, ok := cached.(*models.SearchResult); ok {
			return result, nil
		}
	}

	// Rate limiting
	if err := s.tmdbClient.rateLimiter.Wait(); err != nil {
		return nil, err
	}

	// Build URL with filters
	params := url.Values{}
	params.Add("api_key", s.tmdbClient.config.APIKey)
	params.Add("with_genres", strconv.Itoa(genreID))
	params.Add("page", strconv.Itoa(page))
	params.Add("sort_by", filters.SortBy)

	if filters.MinRating > 0 {
		params.Add("vote_average.gte", fmt.Sprintf("%.1f", filters.MinRating))
	}
	if filters.MaxRating > 0 {
		params.Add("vote_average.lte", fmt.Sprintf("%.1f", filters.MaxRating))
	}
	if filters.MinYear > 0 {
		params.Add("first_air_date.gte", fmt.Sprintf("%d-01-01", filters.MinYear))
	}
	if filters.MaxYear > 0 {
		params.Add("first_air_date.lte", fmt.Sprintf("%d-12-31", filters.MaxYear))
	}

	requestURL := fmt.Sprintf("%s/discover/tv?%s", s.tmdbClient.config.BaseURL, params.Encode())

	// Make request
	resp, err := s.tmdbClient.httpClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to discover TV shows: %w", err)
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
	s.tmdbClient.cache.Set(cacheKey, &result, 30*60*1000) // Cache for 30 minutes

	return &result, nil
}

// DiscoveryFilters represents filters for movie/TV discovery
type DiscoveryFilters struct {
	SortBy    string  `json:"sort_by"`    // popularity.desc, vote_average.desc, release_date.desc, etc.
	MinRating float64 `json:"min_rating"` // Minimum vote average
	MaxRating float64 `json:"max_rating"` // Maximum vote average
	MinYear   int     `json:"min_year"`   // Minimum release year
	MaxYear   int     `json:"max_year"`   // Maximum release year
}

// GetDefaultFilters returns default discovery filters
func GetDefaultFilters() DiscoveryFilters {
	return DiscoveryFilters{
		SortBy:    "popularity.desc",
		MinRating: 0,
		MaxRating: 10,
		MinYear:   1900,
		MaxYear:   2030,
	}
}

// ValidateFilters validates discovery filters
func (f *DiscoveryFilters) Validate() error {
	validSortOptions := map[string]bool{
		"popularity.desc":           true,
		"popularity.asc":            true,
		"vote_average.desc":         true,
		"vote_average.asc":          true,
		"release_date.desc":         true,
		"release_date.asc":          true,
		"revenue.desc":              true,
		"revenue.asc":               true,
		"primary_release_date.desc": true,
		"primary_release_date.asc":  true,
	}

	if !validSortOptions[f.SortBy] {
		f.SortBy = "popularity.desc" // Default fallback
	}

	if f.MinRating < 0 {
		f.MinRating = 0
	}
	if f.MaxRating > 10 {
		f.MaxRating = 10
	}
	if f.MinRating > f.MaxRating {
		f.MinRating, f.MaxRating = f.MaxRating, f.MinRating
	}

	if f.MinYear < 1900 {
		f.MinYear = 1900
	}
	if f.MaxYear > 2030 {
		f.MaxYear = 2030
	}
	if f.MinYear > f.MaxYear {
		f.MinYear, f.MaxYear = f.MaxYear, f.MinYear
	}

	return nil
}

// GetPopularGenres returns the most popular genres based on current trending content
func (s *GenreService) GetPopularGenres(contentType string) ([]models.Genre, error) {
	// This is a simplified implementation
	// In a real system, you'd analyze trending content to determine popular genres

	if contentType == "movie" {
		return s.GetMovieGenres()
	} else {
		return s.GetTVGenres()
	}
}

// SearchByGenreAndKeyword searches for content by both genre and keyword
func (s *GenreService) SearchByGenreAndKeyword(genreID int, keyword string, contentType string, page int) (*models.SearchResult, error) {
	// This would combine genre filtering with keyword search
	// For now, we'll use the basic discovery endpoint
	filters := GetDefaultFilters()
	return s.DiscoverMoviesByGenre(genreID, page, filters)
}
