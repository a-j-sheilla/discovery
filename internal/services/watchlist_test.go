package services

import (
	"testing"

	"movie-discovery-app/internal/models"
)

func TestWatchlistService_AddToWatchlist(t *testing.T) {
	service := NewWatchlistService()
	userID := "test_user"

	item := models.WatchlistItem{
		ID:         "123",
		Type:       "movie",
		Title:      "Test Movie",
		PosterPath: "/test.jpg",
		Watched:    false,
		Rating:     0,
	}

	// Test adding valid item
	err := service.AddToWatchlist(userID, item)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test adding duplicate item
	err = service.AddToWatchlist(userID, item)
	if err == nil {
		t.Error("Expected error for duplicate item")
	}

	// Test adding invalid item (empty ID)
	invalidItem := models.WatchlistItem{
		ID:    "",
		Type:  "movie",
		Title: "Invalid Movie",
	}
	err = service.AddToWatchlist(userID, invalidItem)
	if err == nil {
		t.Error("Expected error for invalid item")
	}
}

func TestWatchlistService_RemoveFromWatchlist(t *testing.T) {
	service := NewWatchlistService()
	userID := "test_user"

	item := models.WatchlistItem{
		ID:         "123",
		Type:       "movie",
		Title:      "Test Movie",
		PosterPath: "/test.jpg",
		Watched:    false,
		Rating:     0,
	}

	// Add item first
	service.AddToWatchlist(userID, item)

	// Test removing existing item
	err := service.RemoveFromWatchlist(userID, "123", "movie")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test removing non-existent item
	err = service.RemoveFromWatchlist(userID, "999", "movie")
	if err == nil {
		t.Error("Expected error for non-existent item")
	}

	// Test removing from non-existent user
	err = service.RemoveFromWatchlist("non_existent_user", "123", "movie")
	if err == nil {
		t.Error("Expected error for non-existent user")
	}
}

func TestWatchlistService_GetWatchlist(t *testing.T) {
	service := NewWatchlistService()
	userID := "test_user"

	// Test getting empty watchlist
	watchlist, err := service.GetWatchlist(userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(watchlist) != 0 {
		t.Errorf("Expected empty watchlist, got %d items", len(watchlist))
	}

	// Add some items
	items := []models.WatchlistItem{
		{
			ID:         "123",
			Type:       "movie",
			Title:      "Movie 1",
			PosterPath: "/movie1.jpg",
			Watched:    false,
			Rating:     0,
		},
		{
			ID:         "456",
			Type:       "tv",
			Title:      "TV Show 1",
			PosterPath: "/tv1.jpg",
			Watched:    true,
			Rating:     8.5,
		},
	}

	for _, item := range items {
		service.AddToWatchlist(userID, item)
	}

	// Test getting populated watchlist
	watchlist, err = service.GetWatchlist(userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(watchlist) != 2 {
		t.Errorf("Expected 2 items, got %d", len(watchlist))
	}
}

func TestWatchlistService_MarkAsWatched(t *testing.T) {
	service := NewWatchlistService()
	userID := "test_user"

	item := models.WatchlistItem{
		ID:         "123",
		Type:       "movie",
		Title:      "Test Movie",
		PosterPath: "/test.jpg",
		Watched:    false,
		Rating:     0,
	}

	// Add item first
	service.AddToWatchlist(userID, item)

	// Test marking as watched with rating
	err := service.MarkAsWatched(userID, "123", "movie", 8.5)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the item is marked as watched
	watchlist, _ := service.GetWatchlist(userID)
	if len(watchlist) != 1 {
		t.Fatal("Expected 1 item in watchlist")
	}
	if !watchlist[0].Watched {
		t.Error("Expected item to be marked as watched")
	}
	if watchlist[0].Rating != 8.5 {
		t.Errorf("Expected rating 8.5, got %f", watchlist[0].Rating)
	}

	// Test marking non-existent item
	err = service.MarkAsWatched(userID, "999", "movie", 7.0)
	if err == nil {
		t.Error("Expected error for non-existent item")
	}
}

func TestWatchlistService_GetWatchedItems(t *testing.T) {
	service := NewWatchlistService()
	userID := "test_user"

	items := []models.WatchlistItem{
		{
			ID:         "123",
			Type:       "movie",
			Title:      "Watched Movie",
			PosterPath: "/watched.jpg",
			Watched:    true,
			Rating:     8.0,
		},
		{
			ID:         "456",
			Type:       "movie",
			Title:      "Unwatched Movie",
			PosterPath: "/unwatched.jpg",
			Watched:    false,
			Rating:     0,
		},
	}

	for _, item := range items {
		service.AddToWatchlist(userID, item)
	}

	watchedItems, err := service.GetWatchedItems(userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(watchedItems) != 1 {
		t.Errorf("Expected 1 watched item, got %d", len(watchedItems))
	}
	if watchedItems[0].ID != "123" {
		t.Errorf("Expected watched item ID 123, got %s", watchedItems[0].ID)
	}
}

func TestWatchlistService_GetUnwatchedItems(t *testing.T) {
	service := NewWatchlistService()
	userID := "test_user"

	items := []models.WatchlistItem{
		{
			ID:         "123",
			Type:       "movie",
			Title:      "Watched Movie",
			PosterPath: "/watched.jpg",
			Watched:    true,
			Rating:     8.0,
		},
		{
			ID:         "456",
			Type:       "movie",
			Title:      "Unwatched Movie",
			PosterPath: "/unwatched.jpg",
			Watched:    false,
			Rating:     0,
		},
	}

	for _, item := range items {
		service.AddToWatchlist(userID, item)
	}

	unwatchedItems, err := service.GetUnwatchedItems(userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(unwatchedItems) != 1 {
		t.Errorf("Expected 1 unwatched item, got %d", len(unwatchedItems))
	}
	if unwatchedItems[0].ID != "456" {
		t.Errorf("Expected unwatched item ID 456, got %s", unwatchedItems[0].ID)
	}
}

func TestWatchlistService_IsInWatchlist(t *testing.T) {
	service := NewWatchlistService()
	userID := "test_user"

	item := models.WatchlistItem{
		ID:         "123",
		Type:       "movie",
		Title:      "Test Movie",
		PosterPath: "/test.jpg",
		Watched:    false,
		Rating:     0,
	}

	// Test item not in watchlist
	exists := service.IsInWatchlist(userID, "123", "movie")
	if exists {
		t.Error("Expected item not to be in watchlist")
	}

	// Add item and test again
	service.AddToWatchlist(userID, item)
	exists = service.IsInWatchlist(userID, "123", "movie")
	if !exists {
		t.Error("Expected item to be in watchlist")
	}

	// Test different type
	exists = service.IsInWatchlist(userID, "123", "tv")
	if exists {
		t.Error("Expected item with different type not to be in watchlist")
	}
}

func TestWatchlistService_GetWatchlistStats(t *testing.T) {
	service := NewWatchlistService()
	userID := "test_user"

	items := []models.WatchlistItem{
		{
			ID:         "123",
			Type:       "movie",
			Title:      "Movie 1",
			PosterPath: "/movie1.jpg",
			Watched:    true,
			Rating:     8.0,
		},
		{
			ID:         "456",
			Type:       "movie",
			Title:      "Movie 2",
			PosterPath: "/movie2.jpg",
			Watched:    true,
			Rating:     9.0,
		},
		{
			ID:         "789",
			Type:       "tv",
			Title:      "TV Show 1",
			PosterPath: "/tv1.jpg",
			Watched:    false,
			Rating:     0,
		},
	}

	for _, item := range items {
		service.AddToWatchlist(userID, item)
	}

	stats, err := service.GetWatchlistStats(userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedStats := map[string]interface{}{
		"total_items":     3,
		"watched_items":   2,
		"unwatched_items": 1,
		"movies":          2,
		"tv_shows":        1,
		"average_rating":  8.5,
	}

	for key, expected := range expectedStats {
		if stats[key] != expected {
			t.Errorf("Expected %s to be %v, got %v", key, expected, stats[key])
		}
	}
}

func TestWatchlistService_validateWatchlistItem(t *testing.T) {
	service := NewWatchlistService()

	tests := []struct {
		name    string
		item    models.WatchlistItem
		wantErr bool
	}{
		{
			name: "Valid movie item",
			item: models.WatchlistItem{
				ID:         "123",
				Type:       "movie",
				Title:      "Test Movie",
				PosterPath: "/test.jpg",
				Watched:    false,
				Rating:     8.5,
			},
			wantErr: false,
		},
		{
			name: "Valid TV item",
			item: models.WatchlistItem{
				ID:         "456",
				Type:       "tv",
				Title:      "Test TV Show",
				PosterPath: "/test.jpg",
				Watched:    true,
				Rating:     0,
			},
			wantErr: false,
		},
		{
			name: "Empty ID",
			item: models.WatchlistItem{
				ID:    "",
				Type:  "movie",
				Title: "Test Movie",
			},
			wantErr: true,
		},
		{
			name: "Invalid type",
			item: models.WatchlistItem{
				ID:    "123",
				Type:  "invalid",
				Title: "Test Movie",
			},
			wantErr: true,
		},
		{
			name: "Empty title",
			item: models.WatchlistItem{
				ID:    "123",
				Type:  "movie",
				Title: "",
			},
			wantErr: true,
		},
		{
			name: "Invalid rating (negative)",
			item: models.WatchlistItem{
				ID:     "123",
				Type:   "movie",
				Title:  "Test Movie",
				Rating: -1,
			},
			wantErr: true,
		},
		{
			name: "Invalid rating (too high)",
			item: models.WatchlistItem{
				ID:     "123",
				Type:   "movie",
				Title:  "Test Movie",
				Rating: 11,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateWatchlistItem(tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateWatchlistItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
