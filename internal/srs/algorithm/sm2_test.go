package algorithm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Quality constants (from srs package to avoid import cycle)
const (
	QualityBlackout = 0
	QualityWrong    = 1
	QualityHard     = 2
	QualityGood     = 3
	QualityEasy     = 4
	QualityPerfect  = 5
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
	assert.GreaterOrEqual(t, newEase, 2.3) // Ease should be reasonable
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

	// Set card to review state
	interval := 30
	repetitions := 5
	easeFactor := 2.5
	originalEase := easeFactor

	// Failed review (quality < 3 means failed)
	newInterval, newReps, newEase := algo.CalculateInterval(
		interval, repetitions, easeFactor,
		int(QualityHard), // QualityHard = 2
		60, 0, false, 0,
	)

	assert.Equal(t, 0, newReps)           // Reset
	assert.Equal(t, 1, newInterval)       // Back to 1 day
	assert.Less(t, newEase, originalEase) // Ease reduced
}

func TestSM2PlusEaseFactorAdjustment(t *testing.T) {
	algo := NewSM2Plus(nil)

	// Quality Perfect (5) should increase ease most
	_, _, ease0 := algo.CalculateInterval(6, 3, 2.5, int(QualityPerfect), 10, 0, false, 0)
	assert.Greater(t, ease0, 2.5)

	// Quality Easy (4) should maintain or increase ease
	_, _, ease1 := algo.CalculateInterval(6, 3, 2.5, int(QualityEasy), 60, 0, false, 0)
	assert.GreaterOrEqual(t, ease1, 2.4)

	// Quality Good (3) - ease behavior varies
	_, _, ease2 := algo.CalculateInterval(6, 3, 2.5, int(QualityGood), 60, 0, false, 0)
	assert.GreaterOrEqual(t, ease2, 2.0) // Should be reasonable

	// Quality Hard (2) should decrease ease
	_, _, ease3 := algo.CalculateInterval(6, 3, 2.5, int(QualityHard), 60, 0, false, 0)
	assert.Less(t, ease3, 2.5)
}

func TestSM2PlusHintPenalty(t *testing.T) {
	algo := NewSM2Plus(nil)

	// Review without hints
	interval1, _, _ := algo.CalculateInterval(6, 3, 2.5, int(QualityGood), 60, 0, false, 0)

	// Review with hints (should be treated as lower quality)
	interval2, _, _ := algo.CalculateInterval(6, 3, 2.5, int(QualityGood), 60, 2, false, 0)

	// Interval with hints should be less or equal
	assert.LessOrEqual(t, interval2, interval1)
}

func TestSM2PlusTimingAdjustment(t *testing.T) {
	algo := NewSM2Plus(nil)

	// Very fast answer (< 10s) should be treated as easier
	interval1, _, _ := algo.CalculateInterval(6, 3, 2.5, int(QualityGood), 5, 0, false, 0)

	// Normal speed
	interval2, _, _ := algo.CalculateInterval(6, 3, 2.5, int(QualityGood), 60, 0, false, 0)

	// Fast answer should have longer or equal interval
	assert.GreaterOrEqual(t, interval1, interval2)
}

func TestSM2PlusMinMaxIntervals(t *testing.T) {
	config := DefaultConfig()
	config.MinInterval = 1
	config.MaxInterval = 100
	algo := NewSM2Plus(config)

	// High ease to push beyond max
	interval, _, _ := algo.CalculateInterval(90, 10, 3.0, int(QualityPerfect), 60, 0, false, 0)

	// Should be clamped to max
	assert.LessOrEqual(t, interval, config.MaxInterval)

	// Failed review should respect minimum
	interval2, _, _ := algo.CalculateInterval(90, 10, 3.0, int(QualityWrong), 60, 0, false, 0)
	assert.GreaterOrEqual(t, interval2, config.MinInterval)
}

func TestSM2PlusDecayForMatureCards(t *testing.T) {
	config := DefaultConfig()
	config.DecayEnabled = true
	config.DecayThreshold = 30
	config.DecayFactor = 0.95
	algo := NewSM2Plus(config)

	interval, _, _ := algo.CalculateInterval(40, 10, 2.0, int(QualityGood), 60, 0, false, 0)

	// Should be less than straight multiplication due to decay
	expectedWithoutDecay := int(float64(40) * 2.0)
	assert.Less(t, interval, expectedWithoutDecay)
}

func TestSM2PlusMinEaseFloor(t *testing.T) {
	config := DefaultConfig()
	config.MinEase = 1.3
	algo := NewSM2Plus(config)

	easeFactor := 1.4 // Just above minimum

	// Multiple failed reviews should not push below minimum
	for i := 0; i < 5; i++ {
		_, _, easeFactor = algo.CalculateInterval(1, 1, easeFactor, int(QualityWrong), 60, 0, false, 0)
	}

	assert.GreaterOrEqual(t, easeFactor, config.MinEase)
}
