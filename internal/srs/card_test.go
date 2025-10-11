package srs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCard(t *testing.T) {
	card := NewCard("test-001", "zookeeper", "L3", "baseline")

	assert.Equal(t, "test-001", card.QuestionID)
	assert.Equal(t, "zookeeper", card.Topic)
	assert.Equal(t, "L3", card.Level)
	assert.Equal(t, "baseline", card.Category)
	assert.Equal(t, 0, card.Interval)
	assert.Equal(t, 0, card.Repetitions)
	assert.Equal(t, 2.5, card.EaseFactor)
	assert.Equal(t, CardStateNew, card.State)
	assert.True(t, card.IsDue())
}

func TestCardIsNew(t *testing.T) {
	card := NewCard("test-001", "test", "L3", "baseline")
	assert.True(t, card.IsNew())

	card.State = CardStateLearning
	assert.False(t, card.IsNew())
}

func TestCardIsDue(t *testing.T) {
	card := NewCard("test-001", "test", "L3", "baseline")

	// New cards are due immediately
	assert.True(t, card.IsDue())

	// Set due date in past
	card.DueDate = time.Now().Add(-24 * time.Hour)
	assert.True(t, card.IsDue())

	// Set due date in future
	card.DueDate = time.Now().Add(24 * time.Hour)
	assert.False(t, card.IsDue())
}

func TestCardIsOverdue(t *testing.T) {
	card := NewCard("test-001", "test", "L3", "baseline")
	card.State = CardStateReview
	card.LastReviewed = time.Now().Add(-10 * 24 * time.Hour)

	// Not overdue if new
	card.State = CardStateNew
	assert.False(t, card.IsOverdue())

	// Overdue if past due date by more than 24 hours
	card.State = CardStateReview
	card.DueDate = time.Now().Add(-25 * time.Hour)
	assert.True(t, card.IsOverdue())

	// Not overdue if within 24 hours
	card.DueDate = time.Now().Add(-23 * time.Hour)
	assert.False(t, card.IsOverdue())
}

func TestCardDaysOverdue(t *testing.T) {
	card := NewCard("test-001", "test", "L3", "baseline")
	card.State = CardStateReview

	// 3 days overdue
	card.DueDate = time.Now().Add(-3 * 24 * time.Hour)
	days := card.DaysOverdue()
	assert.GreaterOrEqual(t, days, 2) // Allow for timing variance
	assert.LessOrEqual(t, days, 3)
}

func TestCardUpdateState(t *testing.T) {
	card := NewCard("test-001", "test", "L3", "baseline")

	// New state
	card.Repetitions = 0
	card.UpdateState()
	assert.Equal(t, CardStateNew, card.State)

	// Learning state
	card.Repetitions = 2
	card.Interval = 6
	card.UpdateState()
	assert.Equal(t, CardStateLearning, card.State)

	// Mature state
	card.Repetitions = 5
	card.Interval = 30
	card.UpdateState()
	assert.Equal(t, CardStateMature, card.State)
}

func TestCardRecordReview(t *testing.T) {
	card := NewCard("test-001", "test", "L3", "baseline")

	card.RecordReview(QualityGood, 120, 0)

	assert.Equal(t, 1, card.TotalReviews)
	assert.Equal(t, 1, card.SuccessCount)
	assert.Equal(t, 120, card.AverageTime)
	assert.False(t, card.LastReviewed.IsZero())

	// Second review
	card.RecordReview(QualityEasy, 90, 1)

	assert.Equal(t, 2, card.TotalReviews)
	assert.Equal(t, 2, card.SuccessCount)
	assert.Equal(t, 105, card.AverageTime) // Average of 120 and 90
	assert.Equal(t, 1, card.HintsUsed)
}

func TestCardRetentionRate(t *testing.T) {
	card := NewCard("test-001", "test", "L3", "baseline")

	// No reviews yet
	assert.Equal(t, 0.0, card.RetentionRate())

	// 3 out of 5 successful
	card.TotalReviews = 5
	card.SuccessCount = 3
	assert.InDelta(t, 0.6, card.RetentionRate(), 0.01)

	// Perfect record
	card.TotalReviews = 10
	card.SuccessCount = 10
	assert.Equal(t, 1.0, card.RetentionRate())
}

func TestReviewQuality(t *testing.T) {
	qualities := []ReviewQuality{
		QualityBlackout,
		QualityWrong,
		QualityHard,
		QualityGood,
		QualityEasy,
		QualityPerfect,
	}

	for i, q := range qualities {
		assert.Equal(t, i, int(q))
	}
}
