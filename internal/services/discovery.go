package services

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"movie-discovery-app/configs"
	"movie-discovery-app/internal/models"
)

// DiscoveryService combines TMDB and OMDB data for comprehensive movie/TV discovery
type DiscoveryService struct {
	tmdbClient *TMDBClient
	omdbClient *OMDBClient
}

// NewDiscoveryService creates a new discovery service
func NewDiscoveryService(config *configs.Config) *DiscoveryService {
	tmdbClient := NewTMDBClient(&config.TMDB, &config.Cache, &config.Rate)
	omdbClient := NewOMDBClient(&config.OMDB, &config.Cache, &config.Rate)

	return &DiscoveryService{
		tmdbClient: tmdbClient,
		omdbClient: omdbClient,
	}
}

// SearchMovies searches for movies using both TMDB and OMDB
func (s *DiscoveryService) SearchMovies(query string, page int) (*models.SearchResult, error) {
	// Get results from TMDB first (primary source)
	tmdbResults, err := s.tmdbClient.SearchMovies(query, page)
	if err != nil {
		log.Printf("TMDB search error: %v", err)
		// Fallback to OMDB if TMDB fails
		return s.searchMoviesOMDBFallback(query, page)
	}

	// Enhance TMDB results with OMDB data
	enhancedResults := make([]interface{}, 0, len(tmdbResults.Results))
	for _, result := range tmdbResults.Results {
		if movieData, ok := result.(map[string]interface{}); ok {
			enhanced := s.enhanceMovieWithOMDB(movieData)
			enhancedResults = append(enhancedResults, enhanced)
		}
	}

	tmdbResults.Results = enhancedResults
	return tmdbResults, nil
}

// SearchTVShows searches for TV shows using both TMDB and OMDB
func (s *DiscoveryService) SearchTVShows(query string, page int) (*models.SearchResult, error) {
	// Get results from TMDB first (primary source)
	tmdbResults, err := s.tmdbClient.SearchTVShows(query, page)
	if err != nil {
		log.Printf("TMDB TV search error: %v", err)
		// Return error for TV shows as OMDB has limited TV support
		return nil, fmt.Errorf("failed to search TV shows: %w", err)
	}

	// Enhance TMDB results with OMDB data where possible
	enhancedResults := make([]interface{}, 0, len(tmdbResults.Results))
	for _, result := range tmdbResults.Results {
		if tvData, ok := result.(map[string]interface{}); ok {
			enhanced := s.enhanceTVShowWithOMDB(tvData)
			enhancedResults = append(enhancedResults, enhanced)
		}
	}

	tmdbResults.Results = enhancedResults
	return tmdbResults, nil
}

// GetMovieDetails gets comprehensive movie details from both APIs
func (s *DiscoveryService) GetMovieDetails(movieID int) (*models.Movie, error) {
	// Get basic details from TMDB
	tmdbMovie, err := s.tmdbClient.GetMovieDetails(movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie details from TMDB: %w", err)
	}

	// Try to enhance with OMDB data
	if tmdbMovie.Title != "" {
		year := ""
		if len(tmdbMovie.ReleaseDate) >= 4 {
			year = tmdbMovie.ReleaseDate[:4]
		}

		omdbMovie, err := s.omdbClient.GetMovieByTitle(tmdbMovie.Title, year)
		if err != nil {
			log.Printf("Failed to get OMDB data for movie %s: %v", tmdbMovie.Title, err)
		} else {
			// Merge OMDB data into TMDB movie
			s.mergeOMDBIntoMovie(tmdbMovie, omdbMovie)
		}
	}

	return tmdbMovie, nil
}

// GetTrendingMovies gets trending movies from TMDB
func (s *DiscoveryService) GetTrendingMovies(timeWindow string, page int) (*models.TrendingResponse, error) {
	// Validate time window
	if timeWindow != "day" && timeWindow != "week" {
		timeWindow = "week" // default
	}

	return s.tmdbClient.GetTrendingMovies(timeWindow, page)
}

// GetTrendingTVShows gets trending TV shows from TMDB
func (s *DiscoveryService) GetTrendingTVShows(timeWindow string, page int) (*models.TrendingResponse, error) {
	// This would be implemented similar to GetTrendingMovies
	// For now, return an empty result
	return &models.TrendingResponse{
		Page:         page,
		Results:      []interface{}{},
		TotalPages:   0,
		TotalResults: 0,
	}, nil
}

// enhanceMovieWithOMDB enhances TMDB movie data with OMDB information
func (s *DiscoveryService) enhanceMovieWithOMDB(movieData map[string]interface{}) map[string]interface{} {
	title, ok := movieData["title"].(string)
	if !ok {
		return movieData
	}

	releaseDate, _ := movieData["release_date"].(string)
	year := ""
	if len(releaseDate) >= 4 {
		year = releaseDate[:4]
	}

	// Try to get OMDB data
	omdbData, err := s.omdbClient.GetMovieByTitle(title, year)
	if err != nil {
		log.Printf("Failed to enhance movie %s with OMDB data: %v", title, err)
		return movieData
	}

	// Add OMDB fields to the movie data
	movieData["imdb_rating"] = omdbData.IMDBRating
	movieData["rotten_tomatoes"] = omdbData.GetRottenTomatoesRating()
	movieData["plot"] = omdbData.Plot
	movieData["director"] = omdbData.Director
	movieData["writer"] = omdbData.Writer
	movieData["actors"] = omdbData.Actors
	movieData["language"] = omdbData.Language
	movieData["country"] = omdbData.Country
	movieData["awards"] = omdbData.Awards
	movieData["imdb_id"] = omdbData.IMDBID

	return movieData
}

// enhanceTVShowWithOMDB enhances TMDB TV show data with OMDB information
func (s *DiscoveryService) enhanceTVShowWithOMDB(tvData map[string]interface{}) map[string]interface{} {
	name, ok := tvData["name"].(string)
	if !ok {
		return tvData
	}

	firstAirDate, _ := tvData["first_air_date"].(string)
	year := ""
	if len(firstAirDate) >= 4 {
		year = firstAirDate[:4]
	}

	// Try to get OMDB data
	omdbData, err := s.omdbClient.GetTVShowByTitle(name, year)
	if err != nil {
		log.Printf("Failed to enhance TV show %s with OMDB data: %v", name, err)
		return tvData
	}

	// Add OMDB fields to the TV show data
	tvData["imdb_rating"] = omdbData.IMDBRating
	tvData["rotten_tomatoes"] = omdbData.GetRottenTomatoesRating()
	tvData["plot"] = omdbData.Plot
	tvData["director"] = omdbData.Director
	tvData["writer"] = omdbData.Writer
	tvData["actors"] = omdbData.Actors
	tvData["language"] = omdbData.Language
	tvData["country"] = omdbData.Country
	tvData["awards"] = omdbData.Awards
	tvData["imdb_id"] = omdbData.IMDBID

	return tvData
}

// mergeOMDBIntoMovie merges OMDB data into a TMDB movie struct
func (s *DiscoveryService) mergeOMDBIntoMovie(movie *models.Movie, omdbData *OMDBResponse) {
	movie.IMDBRating = omdbData.IMDBRating
	movie.RottenTomatoes = omdbData.GetRottenTomatoesRating()
	movie.Plot = omdbData.Plot
	movie.Director = omdbData.Director
	movie.Writer = omdbData.Writer
	movie.Actors = omdbData.Actors
	movie.Language = omdbData.Language
	movie.Country = omdbData.Country
	movie.Awards = omdbData.Awards
	movie.IMDBId = omdbData.IMDBID

	// Parse runtime if available
	if omdbData.Runtime != "" && omdbData.Runtime != "N/A" {
		runtimeStr := strings.Replace(omdbData.Runtime, " min", "", -1)
		if runtime, err := strconv.Atoi(runtimeStr); err == nil {
			movie.Runtime = runtime
		}
	}
}

// searchMoviesOMDBFallback provides fallback search using OMDB when TMDB fails
func (s *DiscoveryService) searchMoviesOMDBFallback(query string, page int) (*models.SearchResult, error) {
	omdbResults, err := s.omdbClient.SearchMovies(query, page)
	if err != nil {
		return nil, fmt.Errorf("both TMDB and OMDB search failed: %w", err)
	}

	// Convert OMDB results to our standard format
	results := make([]interface{}, 0, len(omdbResults.Search))
	for _, movie := range omdbResults.Search {
		movieData := map[string]interface{}{
			"title":        movie.Title,
			"release_date": movie.Year,
			"poster_path":  movie.Poster,
			"imdb_id":      movie.IMDBID,
			"overview":     "", // OMDB search doesn't provide overview
		}
		results = append(results, movieData)
	}

	totalResults, _ := strconv.Atoi(omdbResults.TotalResults)
	totalPages := (totalResults + 9) / 10 // OMDB returns 10 results per page

	return &models.SearchResult{
		Page:         page,
		Results:      results,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	}, nil
}

// ValidateSearchQuery validates and sanitizes search queries
func (s *DiscoveryService) ValidateSearchQuery(query string) error {
	if strings.TrimSpace(query) == "" {
		return fmt.Errorf("search query cannot be empty")
	}
	
	if len(query) > 100 {
		return fmt.Errorf("search query too long (max 100 characters)")
	}
	
	return nil
}

// ValidatePage validates page numbers
func (s *DiscoveryService) ValidatePage(page int) error {
	if page < 1 {
		return fmt.Errorf("page number must be positive")
	}
	
	if page > 1000 {
		return fmt.Errorf("page number too high (max 1000)")
	}
	
	return nil
}
