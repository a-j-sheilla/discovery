package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"movie-discovery-app/internal/models"
)

// YouTubeService handles YouTube API interactions
type YouTubeService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewYouTubeService creates a new YouTube service
func NewYouTubeService() *YouTubeService {
	return &YouTubeService{
		apiKey:  os.Getenv("YOUTUBE_API_KEY"),
		baseURL: os.Getenv("YOUTUBE_BASE_URL"),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// YouTubeSearchResponse represents YouTube search API response
type YouTubeSearchResponse struct {
	Items []YouTubeSearchItem `json:"items"`
}

// YouTubeSearchItem represents a single YouTube search result
type YouTubeSearchItem struct {
	ID      YouTubeVideoID `json:"id"`
	Snippet YouTubeSnippet `json:"snippet"`
}

// YouTubeVideoID represents YouTube video ID
type YouTubeVideoID struct {
	VideoID string `json:"videoId"`
}

// YouTubeSnippet represents YouTube video snippet
type YouTubeSnippet struct {
	Title        string                 `json:"title"`
	Description  string                 `json:"description"`
	ChannelTitle string                 `json:"channelTitle"`
	PublishedAt  string                 `json:"publishedAt"`
	Thumbnails   YouTubeThumbnails      `json:"thumbnails"`
}

// YouTubeThumbnails represents YouTube video thumbnails
type YouTubeThumbnails struct {
	Default YouTubeThumbnail `json:"default"`
	Medium  YouTubeThumbnail `json:"medium"`
	High    YouTubeThumbnail `json:"high"`
}

// YouTubeThumbnail represents a single thumbnail
type YouTubeThumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// SearchTrailers searches for movie/TV show trailers on YouTube
func (s *YouTubeService) SearchTrailers(title string, year string, mediaType string) ([]models.YouTubeVideo, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("YouTube API key not configured")
	}

	// Build search query
	query := fmt.Sprintf("%s trailer", title)
	if year != "" {
		query += fmt.Sprintf(" %s", year)
	}
	if mediaType == "tv" {
		query += " tv series"
	}

	// Prepare API request
	params := url.Values{}
	params.Set("part", "snippet")
	params.Set("q", query)
	params.Set("type", "video")
	params.Set("maxResults", "5")
	params.Set("order", "relevance")
	params.Set("key", s.apiKey)

	requestURL := fmt.Sprintf("%s/search?%s", s.baseURL, params.Encode())

	// Make API request
	resp, err := s.client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search YouTube: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("YouTube API error: %d", resp.StatusCode)
	}

	// Parse response
	var searchResponse YouTubeSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("failed to decode YouTube response: %w", err)
	}

	// Convert to our model
	var trailers []models.YouTubeVideo
	for _, item := range searchResponse.Items {
		if item.ID.VideoID == "" {
			continue
		}

		thumbnail := item.Snippet.Thumbnails.High.URL
		if thumbnail == "" {
			thumbnail = item.Snippet.Thumbnails.Medium.URL
		}
		if thumbnail == "" {
			thumbnail = item.Snippet.Thumbnails.Default.URL
		}

		trailer := models.YouTubeVideo{
			VideoID:      item.ID.VideoID,
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			Thumbnail:    thumbnail,
			ChannelTitle: item.Snippet.ChannelTitle,
			PublishedAt:  item.Snippet.PublishedAt,
		}

		trailers = append(trailers, trailer)
	}

	return trailers, nil
}

// GetOfficialTrailer searches for the most relevant official trailer
func (s *YouTubeService) GetOfficialTrailer(title string, year string, mediaType string) (*models.YouTubeVideo, error) {
	trailers, err := s.SearchTrailers(title, year, mediaType)
	if err != nil {
		return nil, err
	}

	if len(trailers) == 0 {
		return nil, fmt.Errorf("no trailers found")
	}

	// Find the most likely official trailer
	for _, trailer := range trailers {
		titleLower := strings.ToLower(trailer.Title)
		if strings.Contains(titleLower, "official") && 
		   strings.Contains(titleLower, "trailer") {
			return &trailer, nil
		}
	}

	// If no official trailer found, return the first one
	return &trailers[0], nil
}

// IsConfigured checks if YouTube service is properly configured
func (s *YouTubeService) IsConfigured() bool {
	return s.apiKey != "" && s.baseURL != ""
}
