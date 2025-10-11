package algorithm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSM2PlusName(t *testing.T) {
	algo := NewSM2Plus(nil)
	assert.Equal(t, "SM-2+", algo.Name())
}

func TestSM2PlusNewCard(t *testing.T) {
	algo := NewSM2Plus(nil)
	
	// New card state
	interval := 0
	repetitions := 0
	easeFactor := 2.5

	// First review - quality good (3)
	newInterval, newReps, newEase := algo.CalculateInterval(
		interval, repetitions, easeFactor, 3, 60, 0, false, 0,
	)
	
	// Should get first graduation interval (1 day)
	assert.Equal(t, 1, newInterval)
	assert.Equal(t, 1, newReps)
	assert.Greater(t, newEase, 2.4) // Ease should increase slightly
}

func TestSM2PlusGraduationIntervals(t *testing.T) {
	algo := NewSM2Plus(nil)
	
	interval := 0
	repetitions := 0
	easeFactor := 2.5

	// First successful review -> 1 day
	interval1, reps1, ease1 := algo.CalculateInterval(
		interval, repetitions, easeFactor, 3, 60, 0, false, 0,
	)
	assert.Equal(t, 1, interval1)
	assert.Equal(t, 1, reps1)

	// Second successful review -> 3 days
	interval2, reps2, ease2 := algo.CalculateInterval(
		interval1, reps1, ease1, 3, 60, 0, false, 0,
	)
	assert.Equal(t, 3, interval2)
	assert.Equal(t, 2, reps2)

	// Third successful review -> 6 days
	interval3, reps3, ease3 := algo.CalculateInterval(
		interval2, reps2, ease2, 3, 60, 0, false, 0,
	)
	assert.Equal(t, 6, interval3)
	assert.Equal(t, 3, reps3)

	// Fourth review -> start using ease factor
	interval4, reps4, _ := algo.CalculateInterval(
		interval3, reps3, ease3, 3, 60, 0, false, 0,
	)
	assert.Greater(t, interval4, 6) // Should be 6 * ease_factor
	assert.Equal(t, 4, reps4)
}

func TestSM2PlusFailedReview(t *testing.T) {
	algo := NewSM2Plus(nil)
	card := srs.NewCard("test-001", "test", "L3", "baseline")
	
	// Set card to review state
	card.Repetitions = 5
	card.Interval = 30
	card.EaseFactor = 2.5
	originalEase := card.EaseFactor

	// Failed review
	interval := algo.CalculateInterval(card, srs.QualityHard, 60, 0)

	assert.Equal(t, 0, card.Repetitions) // Reset
	assert.Equal(t, 1, interval)         // Back to 1 day
	assert.Less(t, card.EaseFactor, originalEase) // Ease reduced
}

func TestSM2PlusEaseFactorAdjustment(t *testing.T) {
	algo := NewSM2Plus(nil)
	card := srs.NewCard("test-001", "test", "L3", "baseline")
	card.Repetitions = 3
	card.Interval = 6
	card.EaseFactor = 2.5

	// Quality Easy (4) should increase ease
	algo.CalculateInterval(card, srs.QualityEasy, 60, 0)
	assert.Greater(t, card.EaseFactor, 2.5)

	// Reset
	card.EaseFactor = 2.5

	// Quality Good (3) should slightly increase ease
	algo.CalculateInterval(card, srs.QualityGood, 60, 0)
	assert.GreaterOrEqual(t, card.EaseFactor, 2.5)

	// Reset
	card.EaseFactor = 2.5
	card.Repetitions = 3

	// Quality Hard (2) should decrease ease
	algo.CalculateInterval(card, srs.QualityHard, 60, 0)
	assert.Less(t, card.EaseFactor, 2.5)
}

func TestSM2PlusHintPenalty(t *testing.T) {
	algo := NewSM2Plus(nil)
	card1 := srs.NewCard("test-001", "test", "L3", "baseline")
	card2 := srs.NewCard("test-002", "test", "L3", "baseline")
	
	card1.Repetitions = 3
	card1.Interval = 6
	card1.EaseFactor = 2.5
	
	card2.Repetitions = 3
	card2.Interval = 6
	card2.EaseFactor = 2.5

	// Review without hints
	interval1 := algo.CalculateInterval(card1, srs.QualityGood, 60, 0)

	// Review with hints (should be treated as lower quality)
	interval2 := algo.CalculateInterval(card2, srs.QualityGood, 60, 2)

	// Interval with hints should be less
	assert.Less(t, interval2, interval1)
}

func TestSM2PlusTimingAdjustment(t *testing.T) {
	algo := NewSM2Plus(nil)
	
	// Very fast answer (< 10s) should be treated as easier
	card1 := srs.NewCard("test-001", "test", "L3", "baseline")
	card1.Repetitions = 3
	card1.Interval = 6
	card1.EaseFactor = 2.5
	interval1 := algo.CalculateInterval(card1, srs.QualityGood, 5, 0)
	
	// Normal speed
	card2 := srs.NewCard("test-002", "test", "L3", "baseline")
	card2.Repetitions = 3
	card2.Interval = 6
	card2.EaseFactor = 2.5
	interval2 := algo.CalculateInterval(card2, srs.QualityGood, 60, 0)
	
	// Fast answer should have longer interval
	assert.GreaterOrEqual(t, interval1, interval2)
}

func TestSM2PlusMinMaxIntervals(t *testing.T) {
	config := DefaultConfig()
	config.MinInterval = 1
	config.MaxInterval = 100
	algo := NewSM2Plus(config)
	
	card := srs.NewCard("test-001", "test", "L3", "baseline")
	card.Repetitions = 10
	card.Interval = 90
	card.EaseFactor = 3.0 // High ease to push beyond max

	// Should be clamped to max
	interval := algo.CalculateInterval(card, srs.QualityPerfect, 60, 0)
	assert.LessOrEqual(t, interval, config.MaxInterval)
	
	// Failed review should respect minimum
	interval = algo.CalculateInterval(card, srs.QualityWrong, 60, 0)
	assert.GreaterOrEqual(t, interval, config.MinInterval)
}

func TestSM2PlusDecayForMatureCards(t *testing.T) {
	config := DefaultConfig()
	config.DecayEnabled = true
	config.DecayThreshold = 30
	config.DecayFactor = 0.95
	algo := NewSM2Plus(config)

	card := srs.NewCard("test-001", "test", "L3", "baseline")
	card.Repetitions = 10
	card.Interval = 40 // Above decay threshold
	card.EaseFactor = 2.0

	interval := algo.CalculateInterval(card, srs.QualityGood, 60, 0)
	
	// Should be less than straight multiplication due to decay
	expectedWithoutDecay := int(float64(40) * 2.0)
	assert.Less(t, interval, expectedWithoutDecay)
}

func TestSM2PlusMinEaseFloor(t *testing.T) {
	config := DefaultConfig()
	config.MinEase = 1.3
	algo := NewSM2Plus(config)

	card := srs.NewCard("test-001", "test", "L3", "baseline")
	card.EaseFactor = 1.4 // Just above minimum

	// Multiple failed reviews should not push below minimum
	for i := 0; i < 5; i++ {
		algo.CalculateInterval(card, srs.QualityWrong, 60, 0)
		card.Repetitions = 1 // Simulate some progress
	}

	assert.GreaterOrEqual(t, card.EaseFactor, config.MinEase)
}
