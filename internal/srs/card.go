package srs

import (
	"time"
)

// CardState represents the lifecycle state of an SRS card
type CardState string

const (
	CardStateNew       CardState = "new"       // Never reviewed
	CardStateLearning  CardState = "learning"  // In initial learning phase
	CardStateReview    CardState = "review"    // Regular reviews
	CardStateMature    CardState = "mature"    // Long intervals (21+ days)
	CardStateSuspended CardState = "suspended" // User paused
)

// Card represents an SRS flashcard for a question
type Card struct {
	QuestionID string    `json:"question_id"`
	Topic      string    `json:"topic"`  // e.g., "zookeeper", "gfs"
	Level      string    `json:"level"`  // L3, L4, L5, etc.
	Category   string    `json:"category"` // baseline, bar-raiser

	// SRS State
	Interval     int       `json:"interval"`      // Days until next review
	Repetitions  int       `json:"repetitions"`   // Number of successful reviews
	EaseFactor   float64   `json:"ease_factor"`   // Difficulty multiplier (default 2.5)
	DueDate      time.Time `json:"due_date"`      // Next review date
	LastReviewed time.Time `json:"last_reviewed"` // Last review timestamp

	// Performance Tracking
	TotalReviews int `json:"total_reviews"` // All review attempts
	SuccessCount int `json:"success_count"` // Good/excellent reviews
	AverageTime  int `json:"average_time"`  // Average answer time (seconds)
	HintsUsed    int `json:"hints_used"`    // Total hints across reviews

	// Lifecycle
	State     CardState `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ReviewQuality represents how well the user answered
type ReviewQuality int

const (
	QualityBlackout ReviewQuality = 0 // Complete failure
	QualityWrong    ReviewQuality = 1 // Incorrect answer
	QualityHard     ReviewQuality = 2 // Correct with difficulty
	QualityGood     ReviewQuality = 3 // Correct with some thought
	QualityEasy     ReviewQuality = 4 // Correct, felt easy
	QualityPerfect  ReviewQuality = 5 // Instant correct recall
)

// ReviewResult captures user's performance on a question
type ReviewResult struct {
	QuestionID       string        `json:"question_id"`
	Timestamp        time.Time     `json:"timestamp"`
	TimeSpent        int           `json:"time_spent"` // seconds
	Quality          ReviewQuality `json:"quality"`
	HintsUsed        int           `json:"hints_used"`
	PreviousInterval int           `json:"previous_interval"`
	NewInterval      int           `json:"new_interval"`
	WasOverdue       bool          `json:"was_overdue"`
}

// NewCard creates a new SRS card for a question
func NewCard(questionID, topic, level, category string) *Card {
	now := time.Now()
	return &Card{
		QuestionID:   questionID,
		Topic:        topic,
		Level:        level,
		Category:     category,
		Interval:     0,
		Repetitions:  0,
		EaseFactor:   2.5, // Default starting ease
		DueDate:      now, // Due immediately for first review
		LastReviewed: time.Time{},
		TotalReviews: 0,
		SuccessCount: 0,
		AverageTime:  0,
		HintsUsed:    0,
		State:        CardStateNew,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// IsNew returns true if the card has never been reviewed
func (c *Card) IsNew() bool {
	return c.State == CardStateNew
}

// IsDue returns true if the card is due for review
func (c *Card) IsDue() bool {
	return time.Now().After(c.DueDate) || time.Now().Equal(c.DueDate)
}

// IsOverdue returns true if the card is past its due date
func (c *Card) IsOverdue() bool {
	if c.IsNew() {
		return false
	}
	return time.Now().After(c.DueDate.Add(24 * time.Hour))
}

// DaysSinceLastReview returns the number of days since last review
func (c *Card) DaysSinceLastReview() int {
	if c.LastReviewed.IsZero() {
		return 0
	}
	return int(time.Since(c.LastReviewed).Hours() / 24)
}

// DaysOverdue returns how many days past due the card is
func (c *Card) DaysOverdue() int {
	if !c.IsOverdue() {
		return 0
	}
	return int(time.Since(c.DueDate).Hours() / 24)
}

// UpdateState updates the card's state based on interval
func (c *Card) UpdateState() {
	if c.Repetitions == 0 {
		c.State = CardStateNew
	} else if c.Interval < 21 {
		c.State = CardStateLearning
	} else {
		c.State = CardStateMature
	}
}

// UpdateAverageTime updates the running average time
// Note: This should be called AFTER incrementing TotalReviews
func (c *Card) UpdateAverageTime(timeSpent int) {
	if c.TotalReviews == 1 {
		// First review - just set it
		c.AverageTime = timeSpent
	} else {
		// Running average: ((old_avg * (n-1)) + new_value) / n
		c.AverageTime = (c.AverageTime*(c.TotalReviews-1) + timeSpent) / c.TotalReviews
	}
}

// RecordReview updates card statistics after a review
func (c *Card) RecordReview(quality ReviewQuality, timeSpent int, hintsUsed int) {
	c.TotalReviews++
	c.UpdateAverageTime(timeSpent)
	c.HintsUsed += hintsUsed
	
	if quality >= QualityGood {
		c.SuccessCount++
	}
	
	c.LastReviewed = time.Now()
	c.UpdatedAt = time.Now()
}

// RetentionRate returns the success rate (0.0 to 1.0)
func (c *Card) RetentionRate() float64 {
	if c.TotalReviews == 0 {
		return 0.0
	}
	return float64(c.SuccessCount) / float64(c.TotalReviews)
}
