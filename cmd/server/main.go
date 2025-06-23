package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"movie-discovery-app/configs"
	"movie-discovery-app/internal/api"
	"movie-discovery-app/internal/services"
)

func main() {
	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate required API keys
	if config.TMDB.APIKey == "" {
		log.Fatal("TMDB_API_KEY environment variable is required")
	}
	if config.OMDB.APIKey == "" {
		log.Fatal("OMDB_API_KEY environment variable is required")
	}

	// Initialize services
	discoveryService := services.NewDiscoveryService(config)
	watchlistService := services.NewWatchlistService()
	recommendationService := services.NewRecommendationService(discoveryService, watchlistService)
	genreService := services.NewGenreService(config)

	// Initialize handlers
	handlers := api.NewHandlers(discoveryService, watchlistService, recommendationService, genreService)

	// Setup router
	router := api.SetupRouter(handlers)

	// Start server
	addr := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		if err := server.Close(); err != nil {
			log.Printf("Error during server shutdown: %v", err)
		}
	}()

	log.Printf("Starting Movie Discovery App server on %s", addr)
	log.Printf("API endpoints available at http://%s/api/v1", addr)
	log.Printf("Web interface available at http://%s", addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}

	log.Println("Server stopped")
}
