package services

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"movie-discovery-app/internal/models"

	"github.com/jung-kurt/gofpdf/v2"
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
		"total_items":      len(watchlist),
		"watched_items":    0,
		"unwatched_items":  0,
		"movies":           0,
		"tv_shows":         0,
		"average_rating":   0.0,
		"highest_rated":    0.0,
		"total_watch_time": 0, // This could be enhanced with actual runtime data
	}

	var totalRating float64
	var ratedItems int
	var highestRating float64

	for _, item := range watchlist {
		if item.Watched {
			stats["watched_items"] = stats["watched_items"].(int) + 1
			if item.Rating > 0 {
				totalRating += item.Rating
				ratedItems++
				if item.Rating > highestRating {
					highestRating = item.Rating
				}
			}
		} else {
			stats["unwatched_items"] = stats["unwatched_items"].(int) + 1
		}

		switch item.Type {
		case "movie":
			stats["movies"] = stats["movies"].(int) + 1
		case "tv":
			stats["tv_shows"] = stats["tv_shows"].(int) + 1
		}
	}

	if ratedItems > 0 {
		stats["average_rating"] = totalRating / float64(ratedItems)
		stats["highest_rated"] = highestRating
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

// ExportWatchlistAsCSV exports user's watchlist as CSV
func (s *WatchlistService) ExportWatchlistAsCSV(userID string) ([]byte, error) {
	watchlist, err := s.GetWatchlist(userID)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write CSV header
	header := []string{"ID", "Type", "Title", "Poster Path", "Added At", "Watched", "Rating"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write watchlist items
	for _, item := range watchlist {
		record := []string{
			item.ID,
			item.Type,
			item.Title,
			item.PosterPath,
			item.AddedAt.Format("2006-01-02 15:04:05"),
			strconv.FormatBool(item.Watched),
			strconv.FormatFloat(item.Rating, 'f', 1, 64),
		}
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	return buf.Bytes(), nil
}

// ExportWatchlistAsPDF exports user's watchlist as PDF
func (s *WatchlistService) ExportWatchlistAsPDF(userID string) ([]byte, error) {
	watchlist, err := s.GetWatchlist(userID)
	if err != nil {
		return nil, err
	}

	// Create new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "My Watchlist")
	pdf.Ln(15)

	// Add stats
	stats, _ := s.GetWatchlistStats(userID)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, fmt.Sprintf("Total Items: %d | Movies: %d | TV Shows: %d | Watched: %d | To Watch: %d",
		stats["total_items"], stats["movies"], stats["tv_shows"], stats["watched_items"], stats["unwatched_items"]))
	pdf.Ln(10)

	// Table headers
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(15, 8, "Type")
	pdf.Cell(60, 8, "Title")
	pdf.Cell(25, 8, "Added")
	pdf.Cell(20, 8, "Status")
	pdf.Cell(15, 8, "Rating")
	pdf.Ln(8)

	// Table content
	pdf.SetFont("Arial", "", 8)
	for _, item := range watchlist {
		// Check if we need a new page
		if pdf.GetY() > 270 {
			pdf.AddPage()
			// Re-add headers on new page
			pdf.SetFont("Arial", "B", 9)
			pdf.Cell(15, 8, "Type")
			pdf.Cell(60, 8, "Title")
			pdf.Cell(25, 8, "Added")
			pdf.Cell(20, 8, "Status")
			pdf.Cell(15, 8, "Rating")
			pdf.Ln(8)
			pdf.SetFont("Arial", "", 8)
		}

		// Truncate title if too long
		title := item.Title
		if len(title) > 35 {
			title = title[:32] + "..."
		}

		status := "To Watch"
		if item.Watched {
			status = "Watched"
		}

		rating := "-"
		if item.Rating > 0 {
			rating = strconv.FormatFloat(item.Rating, 'f', 1, 64)
		}

		pdf.Cell(15, 6, item.Type)
		pdf.Cell(60, 6, title)
		pdf.Cell(25, 6, item.AddedAt.Format("2006-01-02"))
		pdf.Cell(20, 6, status)
		pdf.Cell(15, 6, rating)
		pdf.Ln(6)
	}

	// Generate PDF bytes
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
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
