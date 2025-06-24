package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"movie-discovery-app/internal/models"
)

// ProvidersService handles watch providers integration
type ProvidersService struct {
	tmdbAPIKey string
	tmdbBaseURL string
	client     *http.Client
}

// NewProvidersService creates a new providers service
func NewProvidersService() *ProvidersService {
	return &ProvidersService{
		tmdbAPIKey:  os.Getenv("TMDB_API_KEY"),
		tmdbBaseURL: os.Getenv("TMDB_BASE_URL"),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetMovieWatchProviders gets watch providers for a movie
func (s *ProvidersService) GetMovieWatchProviders(movieID int) (*models.WatchProviders, error) {
	if s.tmdbAPIKey == "" {
		return nil, fmt.Errorf("TMDB API key not configured")
	}

	url := fmt.Sprintf("%s/movie/%d/watch/providers?api_key=%s", s.tmdbBaseURL, movieID, s.tmdbAPIKey)
	
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie watch providers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	var providers models.WatchProviders
	if err := json.NewDecoder(resp.Body).Decode(&providers); err != nil {
		return nil, fmt.Errorf("failed to decode watch providers response: %w", err)
	}

	return &providers, nil
}

// GetTVWatchProviders gets watch providers for a TV show
func (s *ProvidersService) GetTVWatchProviders(tvID int) (*models.WatchProviders, error) {
	if s.tmdbAPIKey == "" {
		return nil, fmt.Errorf("TMDB API key not configured")
	}

	url := fmt.Sprintf("%s/tv/%d/watch/providers?api_key=%s", s.tmdbBaseURL, tvID, s.tmdbAPIKey)
	
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get TV watch providers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	var providers models.WatchProviders
	if err := json.NewDecoder(resp.Body).Decode(&providers); err != nil {
		return nil, fmt.Errorf("failed to decode watch providers response: %w", err)
	}

	return &providers, nil
}

// GetWatchProviders gets watch providers for any media type
func (s *ProvidersService) GetWatchProviders(mediaID int, mediaType string) (*models.WatchProviders, error) {
	switch mediaType {
	case "movie":
		return s.GetMovieWatchProviders(mediaID)
	case "tv":
		return s.GetTVWatchProviders(mediaID)
	default:
		return nil, fmt.Errorf("unsupported media type: %s", mediaType)
	}
}

// GetProvidersForRegion gets watch providers for a specific region
func (s *ProvidersService) GetProvidersForRegion(mediaID int, mediaType string, region string) (*models.WatchProviderRegion, error) {
	providers, err := s.GetWatchProviders(mediaID, mediaType)
	if err != nil {
		return nil, err
	}

	if regionProviders, exists := providers.Results[region]; exists {
		return &regionProviders, nil
	}

	return nil, fmt.Errorf("no providers found for region: %s", region)
}

// GetAvailableRegions gets all available regions for watch providers
func (s *ProvidersService) GetAvailableRegions(mediaID int, mediaType string) ([]string, error) {
	providers, err := s.GetWatchProviders(mediaID, mediaType)
	if err != nil {
		return nil, err
	}

	var regions []string
	for region := range providers.Results {
		regions = append(regions, region)
	}

	return regions, nil
}

// GetStreamingServices gets only streaming/subscription services for a region
func (s *ProvidersService) GetStreamingServices(mediaID int, mediaType string, region string) ([]models.WatchProvider, error) {
	regionProviders, err := s.GetProvidersForRegion(mediaID, mediaType, region)
	if err != nil {
		return nil, err
	}

	return regionProviders.Flatrate, nil
}

// GetPurchaseOptions gets purchase and rental options for a region
func (s *ProvidersService) GetPurchaseOptions(mediaID int, mediaType string, region string) (map[string][]models.WatchProvider, error) {
	regionProviders, err := s.GetProvidersForRegion(mediaID, mediaType, region)
	if err != nil {
		return nil, err
	}

	options := make(map[string][]models.WatchProvider)
	if len(regionProviders.Buy) > 0 {
		options["buy"] = regionProviders.Buy
	}
	if len(regionProviders.Rent) > 0 {
		options["rent"] = regionProviders.Rent
	}

	return options, nil
}

// IsConfigured checks if providers service is properly configured
func (s *ProvidersService) IsConfigured() bool {
	return s.tmdbAPIKey != "" && s.tmdbBaseURL != ""
}

// GetProviderLogoURL constructs the full URL for a provider logo
func (s *ProvidersService) GetProviderLogoURL(logoPath string) string {
	if logoPath == "" {
		return ""
	}
	return fmt.Sprintf("https://image.tmdb.org/t/p/original%s", logoPath)
}

// GetPopularProviders returns a list of popular streaming providers
func (s *ProvidersService) GetPopularProviders() []models.WatchProvider {
	return []models.WatchProvider{
		{ProviderID: 8, ProviderName: "Netflix", LogoPath: "/t2yyOv40HZeVlLjYsCsPHnWLk4W.jpg", DisplayPriority: 1},
		{ProviderID: 15, ProviderName: "Hulu", LogoPath: "/giwM8XX4V2AQb9vsoN7yti82tKK.jpg", DisplayPriority: 2},
		{ProviderID: 337, ProviderName: "Disney Plus", LogoPath: "/7rwgEs15tFwyR9NPQ5vpzxTj19Q.jpg", DisplayPriority: 3},
		{ProviderID: 384, ProviderName: "HBO Max", LogoPath: "/Ajqyt5aNxNGjmF9uOfxArGrdf3X.jpg", DisplayPriority: 4},
		{ProviderID: 9, ProviderName: "Amazon Prime Video", LogoPath: "/emthp39XA2YScoYL1p0sdbAH2WA.jpg", DisplayPriority: 5},
		{ProviderID: 350, ProviderName: "Apple TV Plus", LogoPath: "/6uhKBfmtzFqOcLousHwZuzcrScK.jpg", DisplayPriority: 6},
		{ProviderID: 531, ProviderName: "Paramount Plus", LogoPath: "/xbhHHa1YgtpwhC8lb1NQ3ACVcLd.jpg", DisplayPriority: 7},
		{ProviderID: 386, ProviderName: "Peacock", LogoPath: "/xTVM8uXT9QocigQ01LKPkBpmpnx.jpg", DisplayPriority: 8},
	}
}
