package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"movie-discovery-app/internal/models"
	"movie-discovery-app/internal/services"

	"github.com/gorilla/mux"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	discoveryService      *services.DiscoveryService
	watchlistService      *services.WatchlistService
	recommendationService *services.RecommendationService
	genreService          *services.GenreService
}

// NewHandlers creates a new handlers instance
func NewHandlers(discoveryService *services.DiscoveryService, watchlistService *services.WatchlistService, recommendationService *services.RecommendationService, genreService *services.GenreService) *Handlers {
	return &Handlers{
		discoveryService:      discoveryService,
		watchlistService:      watchlistService,
		recommendationService: recommendationService,
		genreService:          genreService,
	}
}

// SearchMovies handles movie search requests
func (h *Handlers) SearchMovies(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	// Validate query
	if err := h.discoveryService.ValidateSearchQuery(query); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse page parameter
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	// Validate page
	if err := h.discoveryService.ValidatePage(page); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Search movies
	results, err := h.discoveryService.SearchMovies(query, page)
	if err != nil {
		http.Error(w, fmt.Sprintf("Search failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// SearchTVShows handles TV show search requests
func (h *Handlers) SearchTVShows(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	// Validate query
	if err := h.discoveryService.ValidateSearchQuery(query); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse page parameter
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	// Validate page
	if err := h.discoveryService.ValidatePage(page); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Search TV shows
	results, err := h.discoveryService.SearchTVShows(query, page)
	if err != nil {
		http.Error(w, fmt.Sprintf("Search failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// GetMovieDetails handles movie details requests
func (h *Handlers) GetMovieDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieIDStr := vars["id"]

	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	movie, err := h.discoveryService.GetMovieDetails(movieID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get movie details: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

// GetTrendingMovies handles trending movies requests
func (h *Handlers) GetTrendingMovies(w http.ResponseWriter, r *http.Request) {
	timeWindow := r.URL.Query().Get("time_window")
	if timeWindow == "" {
		timeWindow = "week"
	}

	// Parse page parameter
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	// Validate page
	if err := h.discoveryService.ValidatePage(page); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results, err := h.discoveryService.GetTrendingMovies(timeWindow, page)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get trending movies: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// AddToWatchlist handles adding items to watchlist
func (h *Handlers) AddToWatchlist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For demo purposes, we'll use a default user ID
	// In a real app, this would come from authentication
	userID := "default_user"

	var item models.WatchlistItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.watchlistService.AddToWatchlist(userID, item); err != nil {
		log.Printf("Watchlist validation error: %v, Item: %+v", err, item)
		http.Error(w, fmt.Sprintf("Failed to add to watchlist: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// RemoveFromWatchlist handles removing items from watchlist
func (h *Handlers) RemoveFromWatchlist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := "default_user"
	vars := mux.Vars(r)
	itemID := vars["id"]
	itemType := vars["type"]

	if err := h.watchlistService.RemoveFromWatchlist(userID, itemID, itemType); err != nil {
		http.Error(w, fmt.Sprintf("Failed to remove from watchlist: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// GetWatchlist handles getting user's watchlist
func (h *Handlers) GetWatchlist(w http.ResponseWriter, r *http.Request) {
	userID := "default_user"

	watchlist, err := h.watchlistService.GetWatchlist(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get watchlist: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(watchlist)
}

// MarkAsWatched handles marking items as watched
func (h *Handlers) MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := "default_user"
	vars := mux.Vars(r)
	itemID := vars["id"]
	itemType := vars["type"]

	// Parse rating from request body
	var requestBody struct {
		Rating float64 `json:"rating"`
	}
	json.NewDecoder(r.Body).Decode(&requestBody)

	if err := h.watchlistService.MarkAsWatched(userID, itemID, itemType, requestBody.Rating); err != nil {
		http.Error(w, fmt.Sprintf("Failed to mark as watched: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// GetWatchlistStats handles getting watchlist statistics
func (h *Handlers) GetWatchlistStats(w http.ResponseWriter, r *http.Request) {
	userID := "default_user"

	stats, err := h.watchlistService.GetWatchlistStats(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get watchlist stats: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// HealthCheck handles health check requests
func (h *Handlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "movie-discovery-app",
	})
}

// CORS middleware
func (h *Handlers) EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs HTTP requests
func (h *Handlers) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// GetRecommendations handles recommendation requests
func (h *Handlers) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	userID := "default_user"

	// Parse limit parameter
	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	recommendations, err := h.recommendationService.GetRecommendations(userID, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get recommendations: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendations)
}

// GetMovieGenres handles movie genres requests
func (h *Handlers) GetMovieGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := h.genreService.GetMovieGenres()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get movie genres: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
}

// GetTVGenres handles TV show genres requests
func (h *Handlers) GetTVGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := h.genreService.GetTVGenres()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get TV genres: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
}

// DiscoverByGenre handles genre-based discovery requests
func (h *Handlers) DiscoverByGenre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	genreIDStr := vars["genreId"]

	genreID, err := strconv.Atoi(genreIDStr)
	if err != nil {
		http.Error(w, "Invalid genre ID", http.StatusBadRequest)
		return
	}

	// Parse page parameter
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	// Parse filters
	filters := services.GetDefaultFilters()
	if sortBy := r.URL.Query().Get("sort_by"); sortBy != "" {
		filters.SortBy = sortBy
	}
	if minRatingStr := r.URL.Query().Get("min_rating"); minRatingStr != "" {
		if minRating, err := strconv.ParseFloat(minRatingStr, 64); err == nil {
			filters.MinRating = minRating
		}
	}
	if maxRatingStr := r.URL.Query().Get("max_rating"); maxRatingStr != "" {
		if maxRating, err := strconv.ParseFloat(maxRatingStr, 64); err == nil {
			filters.MaxRating = maxRating
		}
	}
	if minYearStr := r.URL.Query().Get("min_year"); minYearStr != "" {
		if minYear, err := strconv.Atoi(minYearStr); err == nil {
			filters.MinYear = minYear
		}
	}
	if maxYearStr := r.URL.Query().Get("max_year"); maxYearStr != "" {
		if maxYear, err := strconv.Atoi(maxYearStr); err == nil {
			filters.MaxYear = maxYear
		}
	}

	// Validate filters
	filters.Validate()

	results, err := h.genreService.DiscoverMoviesByGenre(genreID, page, filters)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to discover movies: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
