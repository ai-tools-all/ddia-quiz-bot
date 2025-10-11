package srs

import (
	"time"
)

// ReviewSession manages a study session
type ReviewSession struct {
	ID        string    `json:"id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`

	// Session Config
	MaxNewCards int      `json:"max_new_cards"` // New questions per session
	MaxReviews  int      `json:"max_reviews"`   // Review limit
	Topics      []string `json:"topics"`        // Filter by topics
	Levels      []string `json:"levels"`        // Filter by levels

	// Session State
	CardsReviewed []string `json:"cards_reviewed"`
	NewCardsSeen  int      `json:"new_cards_seen"`
	ReviewsDone   int      `json:"reviews_done"`

	// Statistics
	TotalTime    int     `json:"total_time"` // seconds
	AvgQuality   float64 `json:"avg_quality"`
	CardsCorrect int     `json:"cards_correct"`
	CardsFailed  int     `json:"cards_failed"`
}

// Statistics tracks overall SRS progress
type Statistics struct {
	// Overall Stats
	TotalCards    int `json:"total_cards"`
	NewCards      int `json:"new_cards"`
	LearningCards int `json:"learning_cards"`
	ReviewCards   int `json:"review_cards"`
	MatureCards   int `json:"mature_cards"`

	// Today
	DueToday       int `json:"due_today"`
	CompletedToday int `json:"completed_today"`
	NewSeenToday   int `json:"new_seen_today"`

	// Upcoming
	DueTomorrow int `json:"due_tomorrow"`
	DueThisWeek int `json:"due_this_week"`

	// Performance
	RetentionRate float64 `json:"retention_rate"` // % of reviews answered well
	CurrentStreak int     `json:"current_streak"` // Days with reviews
	LongestStreak int     `json:"longest_streak"`

	// Per-Topic Stats
	TopicStats map[string]*TopicStatistics `json:"topic_stats"`

	// Metadata
	LastUpdated time.Time `json:"last_updated"`
}

// TopicStatistics tracks stats for a specific topic
type TopicStatistics struct {
	Topic         string  `json:"topic"`
	TotalCards    int     `json:"total_cards"`
	MasteredCards int     `json:"mastered_cards"` // Mature cards
	RetentionRate float64 `json:"retention_rate"`
	AvgInterval   int     `json:"avg_interval"` // Average days between reviews
}

// NewReviewSession creates a new review session
func NewReviewSession(maxNew, maxReviews int, topics, levels []string) *ReviewSession {
	return &ReviewSession{
		ID:            generateSessionID(),
		StartTime:     time.Now(),
		MaxNewCards:   maxNew,
		MaxReviews:    maxReviews,
		Topics:        topics,
		Levels:        levels,
		CardsReviewed: []string{},
		NewCardsSeen:  0,
		ReviewsDone:   0,
		TotalTime:     0,
		AvgQuality:    0,
		CardsCorrect:  0,
		CardsFailed:   0,
	}
}

// RecordCardReview records a card review in the session
func (s *ReviewSession) RecordCardReview(cardID string, quality ReviewQuality, timeSpent int, isNew bool) {
	s.CardsReviewed = append(s.CardsReviewed, cardID)
	s.TotalTime += timeSpent

	if isNew {
		s.NewCardsSeen++
	} else {
		s.ReviewsDone++
	}

	// Update quality average
	count := len(s.CardsReviewed)
	s.AvgQuality = (s.AvgQuality*float64(count-1) + float64(quality)) / float64(count)

	// Track correct/failed
	if quality >= QualityGood {
		s.CardsCorrect++
	} else {
		s.CardsFailed++
	}
}

// IsComplete returns true if session limits are reached
func (s *ReviewSession) IsComplete() bool {
	totalCards := s.NewCardsSeen + s.ReviewsDone
	return totalCards >= s.MaxReviews
}

// CanAddNewCard returns true if we can show another new card
func (s *ReviewSession) CanAddNewCard() bool {
	return s.NewCardsSeen < s.MaxNewCards
}

// End marks the session as complete
func (s *ReviewSession) End() {
	s.EndTime = time.Now()
}

// Duration returns the session duration in seconds
func (s *ReviewSession) Duration() int {
	if s.EndTime.IsZero() {
		return int(time.Since(s.StartTime).Seconds())
	}
	return int(s.EndTime.Sub(s.StartTime).Seconds())
}

// generateSessionID creates a unique session ID
func generateSessionID() string {
	return time.Now().Format("20060102-150405")
}
