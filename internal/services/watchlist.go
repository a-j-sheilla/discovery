package services

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"movie-discovery-app/internal/models"
)

// WatchlistService manages user watchlists
// Note: In a real application, this would be backed by a database
// For this demo, we're using in-memory storage that simulates localStorage
type WatchlistService struct {
	watchlists map[string][]models.WatchlistItem // userID -> watchlist items
	mu         sync.RWMutex
}

// NewWatchlistService creates a new watchlist service
func NewWatchlistService() *WatchlistService {
	return &WatchlistService{
		watchlists: make(map[string][]models.WatchlistItem),
	}
}

// AddToWatchlist adds an item to user's watchlist
func (s *WatchlistService) AddToWatchlist(userID string, item models.WatchlistItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate item
	if err := s.validateWatchlistItem(item); err != nil {
		return err
	}

	// Set added timestamp
	item.AddedAt = time.Now()

	// Get user's watchlist
	watchlist, exists := s.watchlists[userID]
	if !exists {
		watchlist = []models.WatchlistItem{}
	}

	// Check if item already exists
	for _, existing := range watchlist {
		if existing.ID == item.ID && existing.Type == item.Type {
			return fmt.Errorf("item already in watchlist")
		}
	}

	// Add item to watchlist
	watchlist = append(watchlist, item)
	s.watchlists[userID] = watchlist

	return nil
}

// RemoveFromWatchlist removes an item from user's watchlist
func (s *WatchlistService) RemoveFromWatchlist(userID, itemID, itemType string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	watchlist, exists := s.watchlists[userID]
	if !exists {
		return fmt.Errorf("watchlist not found")
	}

	// Find and remove item
	for i, item := range watchlist {
		if item.ID == itemID && item.Type == itemType {
			// Remove item from slice
			watchlist = append(watchlist[:i], watchlist[i+1:]...)
			s.watchlists[userID] = watchlist
			return nil
		}
	}

	return fmt.Errorf("item not found in watchlist")
}

// GetWatchlist gets user's complete watchlist
func (s *WatchlistService) GetWatchlist(userID string) ([]models.WatchlistItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	watchlist, exists := s.watchlists[userID]
	if !exists {
		return []models.WatchlistItem{}, nil
	}

	// Return a copy to prevent external modification
	result := make([]models.WatchlistItem, len(watchlist))
	copy(result, watchlist)

	return result, nil
}

// MarkAsWatched marks an item as watched in user's watchlist
func (s *WatchlistService) MarkAsWatched(userID, itemID, itemType string, rating float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	watchlist, exists := s.watchlists[userID]
	if !exists {
		return fmt.Errorf("watchlist not found")
	}

	// Find and update item
	for i, item := range watchlist {
		if item.ID == itemID && item.Type == itemType {
			watchlist[i].Watched = true
			if rating > 0 && rating <= 10 {
				watchlist[i].Rating = rating
			}
			s.watchlists[userID] = watchlist
			return nil
		}
	}

	return fmt.Errorf("item not found in watchlist")
}

// MarkAsUnwatched marks an item as unwatched in user's watchlist
func (s *WatchlistService) MarkAsUnwatched(userID, itemID, itemType string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	watchlist, exists := s.watchlists[userID]
	if !exists {
		return fmt.Errorf("watchlist not found")
	}

	// Find and update item
	for i, item := range watchlist {
		if item.ID == itemID && item.Type == itemType {
			watchlist[i].Watched = false
			watchlist[i].Rating = 0
			s.watchlists[userID] = watchlist
			return nil
		}
	}

	return fmt.Errorf("item not found in watchlist")
}

// GetWatchedItems gets only watched items from user's watchlist
func (s *WatchlistService) GetWatchedItems(userID string) ([]models.WatchlistItem, error) {
	watchlist, err := s.GetWatchlist(userID)
	if err != nil {
		return nil, err
	}

	var watched []models.WatchlistItem
	for _, item := range watchlist {
		if item.Watched {
			watched = append(watched, item)
		}
	}

	return watched, nil
}

// GetUnwatchedItems gets only unwatched items from user's watchlist
func (s *WatchlistService) GetUnwatchedItems(userID string) ([]models.WatchlistItem, error) {
	watchlist, err := s.GetWatchlist(userID)
	if err != nil {
		return nil, err
	}

	var unwatched []models.WatchlistItem
	for _, item := range watchlist {
		if !item.Watched {
			unwatched = append(unwatched, item)
		}
	}

	return unwatched, nil
}

// IsInWatchlist checks if an item is in user's watchlist
func (s *WatchlistService) IsInWatchlist(userID, itemID, itemType string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	watchlist, exists := s.watchlists[userID]
	if !exists {
		return false
	}

	for _, item := range watchlist {
		if item.ID == itemID && item.Type == itemType {
			return true
		}
	}

	return false
}

// GetWatchlistStats gets statistics about user's watchlist
func (s *WatchlistService) GetWatchlistStats(userID string) (map[string]interface{}, error) {
	watchlist, err := s.GetWatchlist(userID)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_items":    len(watchlist),
		"watched_items":  0,
		"unwatched_items": 0,
		"movies":         0,
		"tv_shows":       0,
		"average_rating": 0.0,
	}

	var totalRating float64
	var ratedItems int

	for _, item := range watchlist {
		if item.Watched {
			stats["watched_items"] = stats["watched_items"].(int) + 1
			if item.Rating > 0 {
				totalRating += item.Rating
				ratedItems++
			}
		} else {
			stats["unwatched_items"] = stats["unwatched_items"].(int) + 1
		}

		if item.Type == "movie" {
			stats["movies"] = stats["movies"].(int) + 1
		} else if item.Type == "tv" {
			stats["tv_shows"] = stats["tv_shows"].(int) + 1
		}
	}

	if ratedItems > 0 {
		stats["average_rating"] = totalRating / float64(ratedItems)
	}

	return stats, nil
}

// ExportWatchlist exports user's watchlist as JSON
func (s *WatchlistService) ExportWatchlist(userID string) ([]byte, error) {
	watchlist, err := s.GetWatchlist(userID)
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(watchlist, "", "  ")
}

// ImportWatchlist imports a watchlist from JSON
func (s *WatchlistService) ImportWatchlist(userID string, data []byte, merge bool) error {
	var importedWatchlist []models.WatchlistItem
	if err := json.Unmarshal(data, &importedWatchlist); err != nil {
		return fmt.Errorf("invalid watchlist format: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if !merge {
		// Replace existing watchlist
		s.watchlists[userID] = importedWatchlist
		return nil
	}

	// Merge with existing watchlist
	existingWatchlist, exists := s.watchlists[userID]
	if !exists {
		existingWatchlist = []models.WatchlistItem{}
	}

	// Create a map for quick lookup of existing items
	existingItems := make(map[string]bool)
	for _, item := range existingWatchlist {
		key := fmt.Sprintf("%s_%s", item.ID, item.Type)
		existingItems[key] = true
	}

	// Add new items that don't already exist
	for _, item := range importedWatchlist {
		key := fmt.Sprintf("%s_%s", item.ID, item.Type)
		if !existingItems[key] {
			if err := s.validateWatchlistItem(item); err == nil {
				existingWatchlist = append(existingWatchlist, item)
			}
		}
	}

	s.watchlists[userID] = existingWatchlist
	return nil
}

// validateWatchlistItem validates a watchlist item
func (s *WatchlistService) validateWatchlistItem(item models.WatchlistItem) error {
	if item.ID == "" {
		return fmt.Errorf("item ID cannot be empty")
	}

	if item.Type != "movie" && item.Type != "tv" {
		return fmt.Errorf("item type must be 'movie' or 'tv'")
	}

	if item.Title == "" {
		return fmt.Errorf("item title cannot be empty")
	}

	if item.Rating < 0 || item.Rating > 10 {
		return fmt.Errorf("rating must be between 0 and 10")
	}

	return nil
}
