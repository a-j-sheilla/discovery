package models

import "time"

// Movie represents a movie with combined data from TMDB and OMDB
type Movie struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	Overview     string  `json:"overview"`
	ReleaseDate  string  `json:"release_date"`
	PosterPath   string  `json:"poster_path"`
	BackdropPath string  `json:"backdrop_path"`
	VoteAverage  float64 `json:"vote_average"`
	VoteCount    int     `json:"vote_count"`
	Popularity   float64 `json:"popularity"`
	GenreIDs     []int   `json:"genre_ids"`
	Genres       []Genre `json:"genres"`
	Runtime      int     `json:"runtime"`

	// OMDB specific fields
	IMDBRating     string `json:"imdb_rating"`
	RottenTomatoes string `json:"rotten_tomatoes"`
	Plot           string `json:"plot"`
	Director       string `json:"director"`
	Writer         string `json:"writer"`
	Actors         string `json:"actors"`
	Language       string `json:"language"`
	Country        string `json:"country"`
	Awards         string `json:"awards"`
	IMDBId         string `json:"imdb_id"`
}

// TVShow represents a TV show with combined data from TMDB and OMDB
type TVShow struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Overview         string  `json:"overview"`
	FirstAirDate     string  `json:"first_air_date"`
	LastAirDate      string  `json:"last_air_date"`
	PosterPath       string  `json:"poster_path"`
	BackdropPath     string  `json:"backdrop_path"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
	Popularity       float64 `json:"popularity"`
	GenreIDs         []int   `json:"genre_ids"`
	Genres           []Genre `json:"genres"`
	NumberOfSeasons  int     `json:"number_of_seasons"`
	NumberOfEpisodes int     `json:"number_of_episodes"`

	// OMDB specific fields
	IMDBRating     string `json:"imdb_rating"`
	RottenTomatoes string `json:"rotten_tomatoes"`
	Plot           string `json:"plot"`
	Director       string `json:"director"`
	Writer         string `json:"writer"`
	Actors         string `json:"actors"`
	Language       string `json:"language"`
	Country        string `json:"country"`
	Awards         string `json:"awards"`
	IMDBId         string `json:"imdb_id"`
}

// Genre represents a movie/TV show genre
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// SearchResult represents search results
type SearchResult struct {
	Page         int           `json:"page"`
	Results      []interface{} `json:"results"`
	TotalPages   int           `json:"total_pages"`
	TotalResults int           `json:"total_results"`
}

// WatchlistItem represents an item in user's watchlist
type WatchlistItem struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"` // "movie" or "tv"
	Title      string    `json:"title"`
	PosterPath string    `json:"poster_path"`
	AddedAt    time.Time `json:"added_at"`
	Watched    bool      `json:"watched"`
	Rating     float64   `json:"rating,omitempty"`
}

// TrendingResponse represents trending content response
type TrendingResponse struct {
	Page         int           `json:"page"`
	Results      []interface{} `json:"results"`
	TotalPages   int           `json:"total_pages"`
	TotalResults int           `json:"total_results"`
}

// APIError represents an API error response
type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Success    bool   `json:"success"`
}

// WatchProvider represents a streaming service provider
type WatchProvider struct {
	ProviderID      int    `json:"provider_id"`
	ProviderName    string `json:"provider_name"`
	LogoPath        string `json:"logo_path"`
	DisplayPriority int    `json:"display_priority"`
}

// WatchProviders represents watch provider information for a movie/TV show
type WatchProviders struct {
	ID      int                            `json:"id"`
	Results map[string]WatchProviderRegion `json:"results"`
}

// WatchProviderRegion represents watch providers for a specific region
type WatchProviderRegion struct {
	Link     string          `json:"link"`
	Flatrate []WatchProvider `json:"flatrate,omitempty"` // Subscription services
	Buy      []WatchProvider `json:"buy,omitempty"`      // Purchase options
	Rent     []WatchProvider `json:"rent,omitempty"`     // Rental options
}

// Trailer represents a movie/TV show trailer
type Trailer struct {
	ID          string `json:"id"`
	ISO639_1    string `json:"iso_639_1"`
	ISO3166_1   string `json:"iso_3166_1"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Site        string `json:"site"`
	Size        int    `json:"size"`
	Type        string `json:"type"`
	Official    bool   `json:"official"`
	PublishedAt string `json:"published_at"`
}

// TrailerResponse represents the response from trailer API
type TrailerResponse struct {
	ID      int       `json:"id"`
	Results []Trailer `json:"results"`
}

// YouTubeVideo represents a YouTube video from search
type YouTubeVideo struct {
	VideoID      string `json:"video_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Thumbnail    string `json:"thumbnail"`
	ChannelTitle string `json:"channel_title"`
	PublishedAt  string `json:"published_at"`
	Duration     string `json:"duration,omitempty"`
}
