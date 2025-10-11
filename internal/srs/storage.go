package srs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// Storage handles persistence of SRS data
type Storage struct {
	dataDir string
	mu      sync.RWMutex
}

// StorageData represents the persisted card data
type StorageData struct {
	Version   string           `json:"version"`
	UpdatedAt time.Time        `json:"updated_at"`
	Cards     map[string]*Card `json:"cards"` // questionID -> Card
}

// HistoryData represents the persisted review history
type HistoryData struct {
	Version string          `json:"version"`
	Reviews []*ReviewResult `json:"reviews"`
}

// NewStorage creates a new storage instance
func NewStorage(dataDir string) (*Storage, error) {
	if dataDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		dataDir = filepath.Join(homeDir, ".ddia-clicker", "srs")
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	return &Storage{
		dataDir: dataDir,
	}, nil
}

// SaveCards saves all cards to disk
func (s *Storage) SaveCards(cards map[string]*Card) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data := StorageData{
		Version:   "1.0",
		UpdatedAt: time.Now(),
		Cards:     cards,
	}

	return s.writeJSON("cards.json", data)
}

// LoadCards loads all cards from disk
func (s *Storage) LoadCards() (map[string]*Card, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var data StorageData
	if err := s.readJSON("cards.json", &data); err != nil {
		if os.IsNotExist(err) {
			// No cards file yet, return empty map
			return make(map[string]*Card), nil
		}
		return nil, err
	}

	return data.Cards, nil
}

// AppendReview appends a review to the history log
func (s *Storage) AppendReview(review *ReviewResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Load existing history
	var data HistoryData
	if err := s.readJSON("history.json", &data); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// New history file
		data = HistoryData{
			Version: "1.0",
			Reviews: []*ReviewResult{},
		}
	}

	// Append new review
	data.Reviews = append(data.Reviews, review)

	return s.writeJSON("history.json", data)
}

// LoadHistory loads all review history
func (s *Storage) LoadHistory() ([]*ReviewResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var data HistoryData
	if err := s.readJSON("history.json", &data); err != nil {
		if os.IsNotExist(err) {
			return []*ReviewResult{}, nil
		}
		return nil, err
	}

	return data.Reviews, nil
}

// SaveSession saves a review session
func (s *Storage) SaveSession(session *ReviewSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionsDir := filepath.Join(s.dataDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		return fmt.Errorf("failed to create sessions directory: %w", err)
	}

	filename := fmt.Sprintf("%s.json", session.ID)
	filepath := filepath.Join(sessionsDir, filename)

	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	return os.WriteFile(filepath, data, 0644)
}

// LoadSession loads a specific review session
func (s *Storage) LoadSession(sessionID string) (*ReviewSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filename := fmt.Sprintf("%s.json", sessionID)
	filepath := filepath.Join(s.dataDir, "sessions", filename)

	var session ReviewSession
	if err := s.readJSON(filepath, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// LoadTodaysSessions loads all sessions from today
func (s *Storage) LoadTodaysSessions() ([]*ReviewSession, error) {
	today := time.Now().Format("20060102")
	
	sessionsDir := filepath.Join(s.dataDir, "sessions")
	entries, err := os.ReadDir(sessionsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*ReviewSession{}, nil
		}
		return nil, err
	}

	var sessions []*ReviewSession
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		// Check if session is from today
		if len(entry.Name()) >= 8 && entry.Name()[:8] == today {
			session, err := s.LoadSession(entry.Name()[:len(entry.Name())-5]) // remove .json
			if err != nil {
				continue // Skip corrupted sessions
			}
			sessions = append(sessions, session)
		}
	}

	return sessions, nil
}

// Backup creates a backup of the cards file
func (s *Storage) Backup() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	backupDir := filepath.Join(s.dataDir, "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	timestamp := time.Now().Format("20060102-150405")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("cards-%s.json.bak", timestamp))
	
	cardsFile := filepath.Join(s.dataDir, "cards.json")
	
	// Check if cards file exists
	if _, err := os.Stat(cardsFile); os.IsNotExist(err) {
		return nil // Nothing to backup
	}

	// Copy file
	data, err := os.ReadFile(cardsFile)
	if err != nil {
		return fmt.Errorf("failed to read cards file: %w", err)
	}

	if err := os.WriteFile(backupFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	// Clean old backups (keep last 7 days)
	if err := s.cleanOldBackups(7); err != nil {
		// Log but don't fail
		fmt.Printf("Warning: failed to clean old backups: %v\n", err)
	}

	return nil
}

// cleanOldBackups removes backups older than keepDays
func (s *Storage) cleanOldBackups(keepDays int) error {
	backupDir := filepath.Join(s.dataDir, "backups")
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return err
	}

	cutoff := time.Now().AddDate(0, 0, -keepDays)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		info, err := entry.Info()
		if err != nil {
			continue
		}
		
		if info.ModTime().Before(cutoff) {
			filepath := filepath.Join(backupDir, entry.Name())
			os.Remove(filepath) // Ignore errors
		}
	}

	return nil
}

// GetDataDir returns the data directory path
func (s *Storage) GetDataDir() string {
	return s.dataDir
}

// writeJSON writes data as JSON to a file
func (s *Storage) writeJSON(filename string, data interface{}) error {
	filepath := filepath.Join(s.dataDir, filename)
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write atomically using temp file + rename
	tempFile := filepath + ".tmp"
	if err := os.WriteFile(tempFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := os.Rename(tempFile, filepath); err != nil {
		os.Remove(tempFile) // Clean up temp file
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// readJSON reads and unmarshals JSON from a file
func (s *Storage) readJSON(filename string, data interface{}) error {
	var filePath string
	
	// If filename is already an absolute path, use it
	if !filepath.IsAbs(filename) {
		filePath = filepath.Join(s.dataDir, filename)
	} else {
		filePath = filename
	}
	
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, data)
}

// topicData is a helper struct for aggregating topic statistics
type topicData struct {
	totalCards    int
	matureCards   int
	totalReviews  int
	successCount  int
	totalInterval int
}

// CalculateStatistics computes statistics from cards and history
func CalculateStatistics(cards map[string]*Card, sessions []*ReviewSession) *Statistics {
	stats := &Statistics{
		TopicStats:  make(map[string]*TopicStatistics),
		LastUpdated: time.Now(),
	}

	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	weekEnd := today.Add(7 * 24 * time.Hour)

	// Track per-topic data
	topicDataMap := make(map[string]*topicData)

	for _, card := range cards {
		stats.TotalCards++

		// Count by state
		switch card.State {
		case CardStateNew:
			stats.NewCards++
		case CardStateLearning:
			stats.LearningCards++
		case CardStateReview:
			stats.ReviewCards++
		case CardStateMature:
			stats.MatureCards++
		}

		// Count due cards
		if card.IsDue() {
			if card.DueDate.Before(tomorrow) {
				stats.DueToday++
			} else if card.DueDate.Before(tomorrow.Add(24 * time.Hour)) {
				stats.DueTomorrow++
			}
			
			if card.DueDate.Before(weekEnd) {
				stats.DueThisWeek++
			}
		}

		// Per-topic statistics
		if _, exists := topicDataMap[card.Topic]; !exists {
			topicDataMap[card.Topic] = &topicData{
				totalCards:    0,
				matureCards:   0,
				totalReviews:  0,
				successCount:  0,
				totalInterval: 0,
			}
		}
		
		td := topicDataMap[card.Topic]
		td.totalCards++
		if card.State == CardStateMature {
			td.matureCards++
		}
		td.totalReviews += card.TotalReviews
		td.successCount += card.SuccessCount
		td.totalInterval += card.Interval
	}

	// Calculate retention rate
	totalReviews := 0
	totalSuccess := 0
	for _, card := range cards {
		totalReviews += card.TotalReviews
		totalSuccess += card.SuccessCount
	}
	if totalReviews > 0 {
		stats.RetentionRate = float64(totalSuccess) / float64(totalReviews)
	}

	// Count today's completed reviews from sessions
	for _, session := range sessions {
		stats.CompletedToday += len(session.CardsReviewed)
		stats.NewSeenToday += session.NewCardsSeen
	}

	// Calculate streaks
	stats.CurrentStreak = calculateStreak(sessions)
	stats.LongestStreak = calculateLongestStreak(sessions)

	// Build topic statistics
	for topic, td := range topicDataMap {
		topicStat := &TopicStatistics{
			Topic:         topic,
			TotalCards:    td.totalCards,
			MasteredCards: td.matureCards,
		}
		
		if td.totalReviews > 0 {
			topicStat.RetentionRate = float64(td.successCount) / float64(td.totalReviews)
		}
		
		if td.totalCards > 0 {
			topicStat.AvgInterval = td.totalInterval / td.totalCards
		}
		
		stats.TopicStats[topic] = topicStat
	}

	return stats
}

// calculateStreak calculates current review streak in days
func calculateStreak(sessions []*ReviewSession) int {
	if len(sessions) == 0 {
		return 0
	}

	// Sort sessions by date (newest first)
	sortedDates := make([]string, 0)
	dateMap := make(map[string]bool)
	
	for _, session := range sessions {
		date := session.StartTime.Format("20060102")
		dateMap[date] = true
	}
	
	for date := range dateMap {
		sortedDates = append(sortedDates, date)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(sortedDates)))

	// Count consecutive days from today
	streak := 0
	
	for i, date := range sortedDates {
		expectedDate := time.Now().AddDate(0, 0, -i).Format("20060102")
		if date != expectedDate {
			break
		}
		streak++
	}

	return streak
}

// calculateLongestStreak finds the longest review streak
func calculateLongestStreak(sessions []*ReviewSession) int {
	if len(sessions) == 0 {
		return 0
	}

	// Get unique dates
	dateMap := make(map[string]bool)
	for _, session := range sessions {
		date := session.StartTime.Format("20060102")
		dateMap[date] = true
	}
	
	// Convert to sorted slice
	dates := make([]string, 0, len(dateMap))
	for date := range dateMap {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	// Find longest consecutive sequence
	maxStreak := 1
	currentStreak := 1
	
	for i := 1; i < len(dates); i++ {
		prevDate, _ := time.Parse("20060102", dates[i-1])
		currDate, _ := time.Parse("20060102", dates[i])
		
		daysDiff := int(currDate.Sub(prevDate).Hours() / 24)
		
		if daysDiff == 1 {
			currentStreak++
			if currentStreak > maxStreak {
				maxStreak = currentStreak
			}
		} else {
			currentStreak = 1
		}
	}

	return maxStreak
}
