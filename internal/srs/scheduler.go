package srs

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/abhishek/ddia-clicker/internal/models"
	"github.com/abhishek/ddia-clicker/internal/srs/algorithm"
)

// Scheduler manages the SRS system
type Scheduler struct {
	storage   *Storage
	algorithm algorithm.Algorithm
	cards     map[string]*Card // questionID -> Card
}

// NewScheduler creates a new SRS scheduler
func NewScheduler(dataDir string) (*Scheduler, error) {
	storage, err := NewStorage(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	cards, err := storage.LoadCards()
	if err != nil {
		return nil, fmt.Errorf("failed to load cards: %w", err)
	}

	// Use SM-2+ algorithm by default
	algo := algorithm.NewSM2Plus(algorithm.DefaultConfig())

	return &Scheduler{
		storage:   storage,
		algorithm: algo,
		cards:     cards,
	}, nil
}

// AddQuestion adds a new question to the SRS system
func (s *Scheduler) AddQuestion(question *models.Question) error {
	// Check if card already exists
	if _, exists := s.cards[question.ID]; exists {
		return nil // Already added
	}

	// Create new card
	card := NewCard(question.ID, extractTopic(question), question.Level, question.Category)
	s.cards[question.ID] = card

	// Save immediately
	return s.storage.SaveCards(s.cards)
}

// AddQuestions adds multiple questions in bulk
func (s *Scheduler) AddQuestions(questions []*models.Question) error {
	for _, q := range questions {
		if _, exists := s.cards[q.ID]; !exists {
			card := NewCard(q.ID, extractTopic(q), q.Level, q.Category)
			s.cards[q.ID] = card
		}
	}
	
	return s.storage.SaveCards(s.cards)
}

// RecordReview records a review result and updates the card
func (s *Scheduler) RecordReview(questionID string, quality ReviewQuality, timeSpent int, hintsUsed int) (*ReviewResult, error) {
	card, exists := s.cards[questionID]
	if !exists {
		return nil, fmt.Errorf("card not found: %s", questionID)
	}

	// Record previous state
	previousInterval := card.Interval
	wasOverdue := card.IsOverdue()

	// Calculate new interval using algorithm
	newInterval, newRepetitions, newEaseFactor := s.algorithm.CalculateInterval(
		card.Interval,
		card.Repetitions,
		card.EaseFactor,
		int(quality),
		timeSpent,
		hintsUsed,
		card.IsOverdue(),
		card.DaysOverdue(),
	)
	
	// Update card
	card.Interval = newInterval
	card.Repetitions = newRepetitions
	card.EaseFactor = newEaseFactor
	card.DueDate = time.Now().Add(time.Duration(newInterval) * 24 * time.Hour)
	card.RecordReview(quality, timeSpent, hintsUsed)
	card.UpdateState()

	// Create review result
	result := &ReviewResult{
		QuestionID:       questionID,
		Timestamp:        time.Now(),
		TimeSpent:        timeSpent,
		Quality:          quality,
		HintsUsed:        hintsUsed,
		PreviousInterval: previousInterval,
		NewInterval:      newInterval,
		WasOverdue:       wasOverdue,
	}

	// Save to history
	if err := s.storage.AppendReview(result); err != nil {
		return nil, fmt.Errorf("failed to save review: %w", err)
	}

	// Save updated cards
	if err := s.storage.SaveCards(s.cards); err != nil {
		return nil, fmt.Errorf("failed to save cards: %w", err)
	}

	return result, nil
}

// GetDueCards returns cards that are due for review
func (s *Scheduler) GetDueCards(maxNew, maxReviews int, topics, levels []string) []*Card {
	
	var overdue []*Card
	var dueToday []*Card
	var newCards []*Card

	for _, card := range s.cards {
		// Filter by topics if specified
		if len(topics) > 0 && !contains(topics, card.Topic) {
			continue
		}
		
		// Filter by levels if specified
		if len(levels) > 0 && !contains(levels, card.Level) {
			continue
		}

		// Skip suspended cards
		if card.State == CardStateSuspended {
			continue
		}

		// Categorize cards
		if card.IsNew() {
			newCards = append(newCards, card)
		} else if card.IsOverdue() {
			overdue = append(overdue, card)
		} else if card.IsDue() {
			dueToday = append(dueToday, card)
		}
	}

	// Sort overdue by how many days overdue (most overdue first)
	sort.Slice(overdue, func(i, j int) bool {
		return overdue[i].DaysOverdue() > overdue[j].DaysOverdue()
	})

	// Sort due today by due date (earliest first)
	sort.Slice(dueToday, func(i, j int) bool {
		return dueToday[i].DueDate.Before(dueToday[j].DueDate)
	})

	// Shuffle new cards to mix topics
	shuffleCards(newCards)

	// Build result with priority: overdue > due today > new
	result := []*Card{}

	// Add all overdue (critical)
	result = append(result, overdue...)

	// Add due today up to review limit
	remaining := maxReviews - len(overdue)
	if remaining > 0 && len(dueToday) > 0 {
		toAdd := min(remaining, len(dueToday))
		result = append(result, dueToday[:toAdd]...)
	}

	// Add new cards if under total limit
	if len(result) < maxReviews {
		newToAdd := min(maxNew, maxReviews-len(result))
		if newToAdd > 0 && len(newCards) > 0 {
			toAdd := min(newToAdd, len(newCards))
			result = append(result, newCards[:toAdd]...)
		}
	}

	return result
}

// GetCard returns a card by question ID
func (s *Scheduler) GetCard(questionID string) (*Card, bool) {
	card, exists := s.cards[questionID]
	return card, exists
}

// GetStatistics calculates current statistics
func (s *Scheduler) GetStatistics() (*Statistics, error) {
	sessions, err := s.storage.LoadTodaysSessions()
	if err != nil {
		return nil, fmt.Errorf("failed to load sessions: %w", err)
	}

	stats := CalculateStatistics(s.cards, sessions)
	return stats, nil
}

// CreateSession creates a new review session
func (s *Scheduler) CreateSession(maxNew, maxReviews int, topics, levels []string) *ReviewSession {
	return NewReviewSession(maxNew, maxReviews, topics, levels)
}

// SaveSession saves a review session
func (s *Scheduler) SaveSession(session *ReviewSession) error {
	return s.storage.SaveSession(session)
}

// Backup creates a backup of all data
func (s *Scheduler) Backup() error {
	return s.storage.Backup()
}

// GetAllCards returns all cards
func (s *Scheduler) GetAllCards() map[string]*Card {
	return s.cards
}

// GetCardsByTopic returns cards filtered by topic
func (s *Scheduler) GetCardsByTopic(topic string) []*Card {
	var result []*Card
	for _, card := range s.cards {
		if card.Topic == topic {
			result = append(result, card)
		}
	}
	return result
}

// GetCardsByLevel returns cards filtered by level
func (s *Scheduler) GetCardsByLevel(level string) []*Card {
	var result []*Card
	for _, card := range s.cards {
		if card.Level == level {
			result = append(result, card)
		}
	}
	return result
}

// ResetCard resets a card to new state
func (s *Scheduler) ResetCard(questionID string) error {
	card, exists := s.cards[questionID]
	if !exists {
		return fmt.Errorf("card not found: %s", questionID)
	}

	card.Interval = 0
	card.Repetitions = 0
	card.EaseFactor = 2.5
	card.DueDate = time.Now()
	card.State = CardStateNew
	card.UpdatedAt = time.Now()

	return s.storage.SaveCards(s.cards)
}

// SuspendCard suspends a card (won't show in reviews)
func (s *Scheduler) SuspendCard(questionID string) error {
	card, exists := s.cards[questionID]
	if !exists {
		return fmt.Errorf("card not found: %s", questionID)
	}

	card.State = CardStateSuspended
	card.UpdatedAt = time.Now()

	return s.storage.SaveCards(s.cards)
}

// UnsuspendCard reactivates a suspended card
func (s *Scheduler) UnsuspendCard(questionID string) error {
	card, exists := s.cards[questionID]
	if !exists {
		return fmt.Errorf("card not found: %s", questionID)
	}

	if card.State == CardStateSuspended {
		card.UpdateState() // Restore appropriate state
		card.UpdatedAt = time.Now()
	}

	return s.storage.SaveCards(s.cards)
}

// Helper functions

// extractTopic extracts topic from question
func extractTopic(question *models.Question) string {
	// Question ID format: topic-subjective-L3-001
	// or simple: topic-001
	// FilePath might contain: chapters/10-mit-6824-zookeeper/...
	
	// Try to extract from FilePath
	if question.FilePath != "" {
		// Look for pattern like "10-mit-6824-zookeeper"
		parts := splitPath(question.FilePath)
		for _, part := range parts {
			if len(part) > 3 && part[0] >= '0' && part[0] <= '9' {
				// Remove leading number prefix like "10-"
				if idx := indexByte(part, '-'); idx > 0 && idx < len(part)-1 {
					topic := part[idx+1:]
					return topic
				}
			}
		}
	}
	
	// Fallback to extracting from ID
	if question.ID != "" {
		parts := splitString(question.ID, '-')
		if len(parts) > 0 {
			return parts[0]
		}
	}
	
	return "unknown"
}

// Helper to split path
func splitPath(path string) []string {
	var parts []string
	current := ""
	for _, r := range path {
		if r == '/' || r == '\\' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(r)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// Helper to split string by delimiter
func splitString(s string, delim rune) []string {
	var parts []string
	current := ""
	for _, r := range s {
		if r == delim {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(r)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// Helper to find byte in string
func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}

// contains checks if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// min returns minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// shuffleCards randomly shuffles a slice of cards
func shuffleCards(cards []*Card) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
}
