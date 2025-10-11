package algorithm

import (
	"math"
)

// SM2Plus implements the enhanced SuperMemo 2 algorithm
type SM2Plus struct {
	config *Config
}

// NewSM2Plus creates a new SM-2+ algorithm instance
func NewSM2Plus(config *Config) *SM2Plus {
	if config == nil {
		config = DefaultConfig()
	}
	return &SM2Plus{
		config: config,
	}
}

// Name returns the algorithm name
func (s *SM2Plus) Name() string {
	return "SM-2+"
}

// CalculateInterval computes the next review interval
func (s *SM2Plus) CalculateInterval(
	interval int,
	repetitions int,
	easeFactor float64,
	quality int,
	timeSpent int,
	hintsUsed int,
	isOverdue bool,
	daysOverdue int,
) (newInterval int, newRepetitions int, newEaseFactor float64) {
	// Adjust quality based on context
	adjustedQuality := s.adjustQuality(quality, timeSpent, hintsUsed)

	// Failed review - reset to learning
	if adjustedQuality < 3 { // Less than "Good"
		newRepetitions = 0
		newEaseFactor = math.Max(s.config.MinEase, easeFactor-0.2)
		newInterval = s.config.MinInterval
		return
	}

	// Successful review
	newRepetitions = repetitions + 1
	newEaseFactor = s.updateEaseFactor(easeFactor, adjustedQuality)

	// Calculate interval
	if newRepetitions <= len(s.config.NewCardIntervals) {
		// Use graduated intervals for new cards
		newInterval = s.config.NewCardIntervals[newRepetitions-1]
	} else {
		// Standard SM-2 formula
		newInterval = int(float64(interval) * newEaseFactor)

		// Apply decay for mature cards
		if s.config.DecayEnabled && newInterval >= s.config.DecayThreshold {
			newInterval = int(float64(newInterval) * s.config.DecayFactor)
		}
	}

	// Apply overdue adjustment
	if isOverdue {
		newInterval = s.applyOverdueAdjustment(interval, daysOverdue, newInterval)
	}

	// Clamp to limits
	if newInterval < s.config.MinInterval {
		newInterval = s.config.MinInterval
	}
	if newInterval > s.config.MaxInterval {
		newInterval = s.config.MaxInterval
	}

	return
}

// adjustQuality modifies quality based on hints and timing
func (s *SM2Plus) adjustQuality(quality int, timeSpent int, hintsUsed int) int {
	adjusted := quality

	// Hint penalty - each hint reduces quality by 1
	if hintsUsed > 0 {
		adjusted -= s.config.HintPenalty * hintsUsed
	}

	// Time-based adjustment (only for successful answers)
	if quality >= 3 { // Quality Good or better
		// Very fast answer = easier than thought
		if timeSpent < s.config.FastAnswerThreshold {
			adjusted++
		}
		// Very slow answer = harder than thought
		if timeSpent > s.config.SlowAnswerThreshold {
			adjusted--
		}
	}

	// Clamp to valid range (0-5)
	if adjusted < 0 {
		adjusted = 0
	}
	if adjusted > 5 {
		adjusted = 5
	}

	return adjusted
}

// updateEaseFactor updates the ease factor based on review quality
func (s *SM2Plus) updateEaseFactor(currentEase float64, quality int) float64 {
	// SM-2 formula: EF' = EF + (0.1 - (5-q) * (0.08 + (5-q) * 0.02))
	q := float64(quality)
	newEase := currentEase + (0.1 - (5-q)*(0.08+(5-q)*0.02))

	// Floor at minimum ease
	if newEase < s.config.MinEase {
		newEase = s.config.MinEase
	}

	return newEase
}

// applyOverdueAdjustment reduces interval for overdue cards
func (s *SM2Plus) applyOverdueAdjustment(previousInterval int, daysOverdue int, newInterval int) int {
	if daysOverdue == 0 {
		return newInterval
	}

	// Calculate reduction factor based on how overdue
	// More overdue = bigger reduction (up to max)
	reductionFactor := math.Min(
		s.config.OverdueReductionMax,
		float64(daysOverdue)/float64(previousInterval)*0.5,
	)

	reducedInterval := int(float64(newInterval) * (1.0 - reductionFactor))

	// Ensure minimum interval
	if reducedInterval < s.config.MinInterval {
		reducedInterval = s.config.MinInterval
	}

	return reducedInterval
}
