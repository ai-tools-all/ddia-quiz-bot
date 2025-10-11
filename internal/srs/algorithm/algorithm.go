package algorithm

// Algorithm defines the interface for spaced repetition algorithms
type Algorithm interface {
	// CalculateInterval computes the next review interval based on performance
	// Returns the new interval in days
	CalculateInterval(
		interval int,           // Current interval
		repetitions int,        // Number of successful reviews
		easeFactor float64,     // Current ease factor
		quality int,            // Review quality (0-5)
		timeSpent int,          // Time spent in seconds
		hintsUsed int,          // Number of hints used
		isOverdue bool,         // Whether card is overdue
		daysOverdue int,        // How many days overdue
	) (newInterval int, newRepetitions int, newEaseFactor float64)
	
	// Name returns the algorithm name
	Name() string
}

// Config holds algorithm configuration
type Config struct {
	// Starting values
	StartingEase float64 `json:"starting_ease"` // Default 2.5
	MinEase      float64 `json:"min_ease"`      // Floor at 1.3
	
	// Intervals
	MinInterval int `json:"min_interval"` // Minimum 1 day
	MaxInterval int `json:"max_interval"` // Maximum 365 days
	
	// New card graduation intervals
	NewCardIntervals []int `json:"new_card_intervals"` // [1, 3, 6]
	
	// Adjustments
	HintPenalty           int     `json:"hint_penalty"`             // Quality reduction per hint
	FastAnswerThreshold   int     `json:"fast_answer_threshold"`    // Seconds for "easy" bonus
	SlowAnswerThreshold   int     `json:"slow_answer_threshold"`    // Seconds for difficulty penalty
	OverdueReductionMax   float64 `json:"overdue_reduction_max"`    // Max interval reduction (0.5 = 50%)
	
	// Decay for mature cards
	DecayEnabled   bool    `json:"decay_enabled"`
	DecayThreshold int     `json:"decay_threshold"` // Days before applying decay
	DecayFactor    float64 `json:"decay_factor"`    // Reduction factor (0.95 = 5% reduction)
	
	// Mature threshold
	MatureThreshold int `json:"mature_threshold"` // Days to be considered "mature"
}

// DefaultConfig returns the default algorithm configuration
func DefaultConfig() *Config {
	return &Config{
		StartingEase:        2.5,
		MinEase:             1.3,
		MinInterval:         1,
		MaxInterval:         365,
		NewCardIntervals:    []int{1, 3, 6},
		HintPenalty:         1,
		FastAnswerThreshold: 10,
		SlowAnswerThreshold: 300,
		OverdueReductionMax: 0.5,
		DecayEnabled:        true,
		DecayThreshold:      30,
		DecayFactor:         0.95,
		MatureThreshold:     21,
	}
}
