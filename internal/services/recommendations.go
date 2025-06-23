package services

import (
	"fmt"
	"math"
	"sort"

	"movie-discovery-app/internal/models"
)

// RecommendationService provides movie/TV show recommendations
type RecommendationService struct {
	discoveryService *DiscoveryService
	watchlistService *WatchlistService
}

// NewRecommendationService creates a new recommendation service
func NewRecommendationService(discoveryService *DiscoveryService, watchlistService *WatchlistService) *RecommendationService {
	return &RecommendationService{
		discoveryService: discoveryService,
		watchlistService: watchlistService,
	}
}

// RecommendationScore represents a recommendation with its score
type RecommendationScore struct {
	Item  interface{} `json:"item"`
	Score float64     `json:"score"`
	Type  string      `json:"type"`
}

// GetRecommendations gets personalized recommendations for a user
func (s *RecommendationService) GetRecommendations(userID string, limit int) ([]RecommendationScore, error) {
	// Get user's watchlist to understand preferences
	watchlist, err := s.watchlistService.GetWatchlist(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user watchlist: %w", err)
	}

	if len(watchlist) == 0 {
		// No watchlist data, return trending content as fallback
		return s.getTrendingRecommendations(limit)
	}

	// Analyze user preferences
	preferences := s.analyzeUserPreferences(watchlist)
	
	// Get recommendations based on preferences
	recommendations := s.generateRecommendations(preferences, limit)
	
	return recommendations, nil
}

// UserPreferences represents analyzed user preferences
type UserPreferences struct {
	FavoriteGenres    map[int]float64 // genre_id -> preference_score
	PreferredRatings  []float64       // ratings of liked movies
	PreferredYears    []int           // years of liked movies
	MovieVsTVRatio    float64         // preference for movies vs TV shows
	AverageRating     float64         // user's average rating
}

// analyzeUserPreferences analyzes user's watchlist to determine preferences
func (s *RecommendationService) analyzeUserPreferences(watchlist []models.WatchlistItem) *UserPreferences {
	preferences := &UserPreferences{
		FavoriteGenres: make(map[int]float64),
	}

	var totalRating float64
	var ratedItems int
	var movieCount, tvCount int

	for _, item := range watchlist {
		// Count movie vs TV preference
		if item.Type == "movie" {
			movieCount++
		} else {
			tvCount++
		}

		// Analyze ratings
		if item.Rating > 0 {
			preferences.PreferredRatings = append(preferences.PreferredRatings, item.Rating)
			totalRating += item.Rating
			ratedItems++
		}

		// Note: In a real implementation, you'd fetch full movie details
		// to get genre information. For this demo, we'll use simplified logic.
	}

	// Calculate movie vs TV ratio
	total := movieCount + tvCount
	if total > 0 {
		preferences.MovieVsTVRatio = float64(movieCount) / float64(total)
	}

	// Calculate average rating
	if ratedItems > 0 {
		preferences.AverageRating = totalRating / float64(ratedItems)
	}

	return preferences
}

// generateRecommendations generates recommendations based on user preferences
func (s *RecommendationService) generateRecommendations(preferences *UserPreferences, limit int) []RecommendationScore {
	var recommendations []RecommendationScore

	// For this demo, we'll use trending content and apply preference-based scoring
	// In a real system, you'd use collaborative filtering, content-based filtering, etc.

	// Get trending movies
	trendingMovies, err := s.discoveryService.GetTrendingMovies("week", 1)
	if err == nil && trendingMovies.Results != nil {
		for _, item := range trendingMovies.Results {
			if movieData, ok := item.(map[string]interface{}); ok {
				score := s.calculateRecommendationScore(movieData, preferences, "movie")
				recommendations = append(recommendations, RecommendationScore{
					Item:  movieData,
					Score: score,
					Type:  "movie",
				})
			}
		}
	}

	// Sort by score (highest first)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Return top recommendations
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations
}

// calculateRecommendationScore calculates a recommendation score for an item
func (s *RecommendationService) calculateRecommendationScore(item map[string]interface{}, preferences *UserPreferences, itemType string) float64 {
	score := 0.0

	// Base score from popularity/rating
	if popularity, ok := item["popularity"].(float64); ok {
		score += math.Log(popularity+1) * 0.1 // Log scale to prevent dominance
	}

	if voteAverage, ok := item["vote_average"].(float64); ok {
		score += voteAverage * 0.3
	}

	// Preference for movie vs TV
	if itemType == "movie" {
		score += preferences.MovieVsTVRatio * 2.0
	} else {
		score += (1.0 - preferences.MovieVsTVRatio) * 2.0
	}

	// Recency bonus (newer content gets slight boost)
	if releaseDate, ok := item["release_date"].(string); ok && len(releaseDate) >= 4 {
		// Simple recency calculation - in real system you'd parse the date properly
		if releaseDate >= "2020" {
			score += 0.5
		}
	}

	// Genre matching would go here if we had genre data
	// For now, we'll use a simplified approach

	return score
}

// getTrendingRecommendations returns trending content as fallback recommendations
func (s *RecommendationService) getTrendingRecommendations(limit int) ([]RecommendationScore, error) {
	var recommendations []RecommendationScore

	// Get trending movies
	trendingMovies, err := s.discoveryService.GetTrendingMovies("week", 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending movies: %w", err)
	}

	if trendingMovies.Results != nil {
		for i, item := range trendingMovies.Results {
			if i >= limit {
				break
			}
			
			if movieData, ok := item.(map[string]interface{}); ok {
				// Simple scoring based on popularity and rating
				score := 5.0 // Base score
				if popularity, ok := movieData["popularity"].(float64); ok {
					score += math.Log(popularity+1) * 0.1
				}
				if voteAverage, ok := movieData["vote_average"].(float64); ok {
					score += voteAverage * 0.3
				}

				recommendations = append(recommendations, RecommendationScore{
					Item:  movieData,
					Score: score,
					Type:  "movie",
				})
			}
		}
	}

	return recommendations, nil
}

// GetSimilarMovies gets movies similar to a given movie
func (s *RecommendationService) GetSimilarMovies(movieID int, limit int) ([]RecommendationScore, error) {
	// In a real implementation, you'd use the TMDB "similar movies" endpoint
	// For this demo, we'll return trending movies as similar content
	return s.getTrendingRecommendations(limit)
}

// GetRecommendationsByGenre gets recommendations for a specific genre
func (s *RecommendationService) GetRecommendationsByGenre(genreID int, limit int) ([]RecommendationScore, error) {
	// In a real implementation, you'd filter by genre
	// For this demo, we'll return trending content
	return s.getTrendingRecommendations(limit)
}

// RecommendationExplanation provides explanation for why an item was recommended
type RecommendationExplanation struct {
	Reason string  `json:"reason"`
	Score  float64 `json:"score"`
}

// ExplainRecommendation explains why an item was recommended
func (s *RecommendationService) ExplainRecommendation(item interface{}, preferences *UserPreferences) RecommendationExplanation {
	// Simple explanation logic
	if movieData, ok := item.(map[string]interface{}); ok {
		if voteAverage, ok := movieData["vote_average"].(float64); ok && voteAverage >= 7.0 {
			return RecommendationExplanation{
				Reason: "Highly rated movie",
				Score:  voteAverage,
			}
		}
		
		if popularity, ok := movieData["popularity"].(float64); ok && popularity >= 100 {
			return RecommendationExplanation{
				Reason: "Popular trending movie",
				Score:  popularity,
			}
		}
	}

	return RecommendationExplanation{
		Reason: "Based on your preferences",
		Score:  5.0,
	}
}
