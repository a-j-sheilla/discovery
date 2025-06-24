package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRouter sets up all routes for the application
func SetupRouter(handlers *Handlers) *mux.Router {
	r := mux.NewRouter()

	// Apply middleware
	r.Use(handlers.EnableCORS)
	r.Use(handlers.LoggingMiddleware)

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Health check
	api.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	// Search endpoints
	api.HandleFunc("/search/movies", handlers.SearchMovies).Methods("GET")
	api.HandleFunc("/search/tv", handlers.SearchTVShows).Methods("GET")

	// Movie details
	api.HandleFunc("/movies/{id:[0-9]+}", handlers.GetMovieDetails).Methods("GET")

	// TV show details
	api.HandleFunc("/tv/{id:[0-9]+}", handlers.GetTVShowDetails).Methods("GET")

	// Trending content
	api.HandleFunc("/trending/movies", handlers.GetTrendingMovies).Methods("GET")

	// Recommendations
	api.HandleFunc("/recommendations", handlers.GetRecommendations).Methods("GET")

	// Genres
	api.HandleFunc("/genres/movies", handlers.GetMovieGenres).Methods("GET")
	api.HandleFunc("/genres/tv", handlers.GetTVGenres).Methods("GET")
	api.HandleFunc("/discover/genre/{genreId:[0-9]+}", handlers.DiscoverByGenre).Methods("GET")

	// Watchlist endpoints
	api.HandleFunc("/watchlist", handlers.GetWatchlist).Methods("GET")
	api.HandleFunc("/watchlist", handlers.AddToWatchlist).Methods("POST")
	api.HandleFunc("/watchlist/{type}/{id}", handlers.RemoveFromWatchlist).Methods("DELETE")
	api.HandleFunc("/watchlist/{type}/{id}/watched", handlers.MarkAsWatched).Methods("PUT")
	api.HandleFunc("/watchlist/stats", handlers.GetWatchlistStats).Methods("GET")

	// Static file serving
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	// Serve the main page
	r.HandleFunc("/", serveIndex).Methods("GET")
	r.HandleFunc("/movie/{id:[0-9]+}", serveIndex).Methods("GET")
	r.HandleFunc("/search", serveIndex).Methods("GET")
	r.HandleFunc("/watchlist", serveIndex).Methods("GET")
	r.HandleFunc("/trending", serveIndex).Methods("GET")

	return r
}

// serveIndex serves the main HTML page
func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/templates/index.html")
}
