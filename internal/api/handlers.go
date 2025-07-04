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

// GetTVShowDetails handles TV show details requests
func (h *Handlers) GetTVShowDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tvIDStr := vars["id"]

	tvID, err := strconv.Atoi(tvIDStr)
	if err != nil {
		http.Error(w, "Invalid TV show ID", http.StatusBadRequest)
		return
	}

	tvShow, err := h.discoveryService.GetTVShowDetails(tvID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get TV show details: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tvShow)
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

// MarkAsUnwatched handles marking items as unwatched
func (h *Handlers) MarkAsUnwatched(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := "default_user"
	vars := mux.Vars(r)
	itemID := vars["id"]
	itemType := vars["type"]

	if err := h.watchlistService.MarkAsUnwatched(userID, itemID, itemType); err != nil {
		http.Error(w, fmt.Sprintf("Failed to mark as unwatched: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// ExportWatchlistAsJSON handles exporting watchlist as JSON
func (h *Handlers) ExportWatchlistAsJSON(w http.ResponseWriter, r *http.Request) {
	userID := "default_user"

	data, err := h.watchlistService.ExportWatchlist(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to export watchlist: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=watchlist.json")
	w.Write(data)
}

// ExportWatchlistAsCSV handles exporting watchlist as CSV
func (h *Handlers) ExportWatchlistAsCSV(w http.ResponseWriter, r *http.Request) {
	userID := "default_user"

	data, err := h.watchlistService.ExportWatchlistAsCSV(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to export watchlist as CSV: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=watchlist.csv")
	w.Write(data)
}

// ExportWatchlistAsPDF handles exporting watchlist as PDF
func (h *Handlers) ExportWatchlistAsPDF(w http.ResponseWriter, r *http.Request) {
	userID := "default_user"

	data, err := h.watchlistService.ExportWatchlistAsPDF(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to export watchlist as PDF: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=watchlist.pdf")
	w.Write(data)
}

// GetTrailers handles getting trailers for a movie or TV show
func (h *Handlers) GetTrailers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mediaType := vars["type"]
	mediaIDStr := vars["id"]

	mediaID, err := strconv.Atoi(mediaIDStr)
	if err != nil {
		http.Error(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	var trailers []models.YouTubeVideo
	switch mediaType {
	case "movie":
		trailers, err = h.discoveryService.GetMovieTrailers(mediaID)
	case "tv":
		trailers, err = h.discoveryService.GetTVTrailers(mediaID)
	default:
		http.Error(w, "Invalid media type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get trailers: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trailers)
}

// GetOfficialTrailer handles getting the official trailer for a movie or TV show
func (h *Handlers) GetOfficialTrailer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mediaType := vars["type"]
	mediaIDStr := vars["id"]

	mediaID, err := strconv.Atoi(mediaIDStr)
	if err != nil {
		http.Error(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	trailer, err := h.discoveryService.GetOfficialTrailer(mediaID, mediaType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get official trailer: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trailer)
}

// GetWatchProviders handles getting watch providers for a movie or TV show
func (h *Handlers) GetWatchProviders(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mediaType := vars["type"]
	mediaIDStr := vars["id"]

	mediaID, err := strconv.Atoi(mediaIDStr)
	if err != nil {
		http.Error(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	providers, err := h.discoveryService.GetWatchProviders(mediaID, mediaType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get watch providers: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(providers)
}

// GetStreamingServices handles getting streaming services for a specific region
func (h *Handlers) GetStreamingServices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mediaType := vars["type"]
	mediaIDStr := vars["id"]
	region := r.URL.Query().Get("region")

	if region == "" {
		region = "US" // Default to US
	}

	mediaID, err := strconv.Atoi(mediaIDStr)
	if err != nil {
		http.Error(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	services, err := h.discoveryService.GetStreamingServices(mediaID, mediaType, region)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get streaming services: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
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

	// Parse content type parameter (movies or tv)
	contentType := r.URL.Query().Get("type")
	if contentType == "" {
		contentType = "movies" // Default to movies
	}

	// Validate content type
	if contentType != "movies" && contentType != "tv" {
		http.Error(w, "Invalid content type. Must be 'movies' or 'tv'", http.StatusBadRequest)
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

	// Call appropriate discovery method based on content type
	var results *models.SearchResult
	switch contentType {
	case "movies":
		results, err = h.genreService.DiscoverMoviesByGenre(genreID, page, filters)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to discover movies: %v", err), http.StatusInternalServerError)
			return
		}
	case "tv":
		results, err = h.genreService.DiscoverTVShowsByGenre(genreID, page, filters)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to discover TV shows: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
