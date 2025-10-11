package srs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMCQCardTracking tests MCQ-specific card behavior
func TestMCQCardTracking(t *testing.T) {
	tests := []struct {
		name         string
		attempts     []mcqAttempt
		finalAccuracy float64
		description   string
	}{
		{
			name: "Perfect accuracy",
			attempts: []mcqAttempt{
				{selected: "A", correct: true, quality: QualityGood},
				{selected: "B", correct: true, quality: QualityGood},
				{selected: "C", correct: true, quality: QualityEasy},
			},
			finalAccuracy: 1.0,
			description: "All answers correct - 100% accuracy",
		},
		{
			name: "Mixed accuracy",
			attempts: []mcqAttempt{
				{selected: "A", correct: true, quality: QualityGood},
				{selected: "B", correct: false, quality: QualityWrong},
				{selected: "C", correct: true, quality: QualityGood},
				{selected: "A", correct: false, quality: QualityWrong},
			},
			finalAccuracy: 0.5,
			description: "50% accuracy (2 correct out of 4)",
		},
		{
			name: "Poor accuracy",
			attempts: []mcqAttempt{
				{selected: "A", correct: false, quality: QualityWrong},
				{selected: "B", correct: false, quality: QualityWrong},
				{selected: "C", correct: true, quality: QualityGood},
			},
			finalAccuracy: 0.33,
			description: "33% accuracy (1 correct out of 3)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create MCQ card
			card := NewCard("mcq-test-001", "gfs", "L3", "baseline")
			card.QuestionType = "mcq"

			// Simulate MCQ attempts
			correctCount := 0
			for _, attempt := range tt.attempts {
				if attempt.correct {
					correctCount++
				}
				// Record MCQ attempt first (selects and checks)
				card.RecordMCQAttempt(attempt.selected, attempt.correct)
				// Then record the review outcome
				card.RecordReview(attempt.quality, 30, 0) // 30 seconds, no hints
			}

			// Calculate expected accuracy manually
			expectedAccuracy := float64(correctCount) / float64(len(tt.attempts))

			// Verify tracking
			assert.Equal(t, len(tt.attempts), card.TotalReviews, "Should track total reviews")
			assert.Equal(t, card.QuestionType, "mcq", "Should maintain MCQ type")
			assert.InDelta(t, expectedAccuracy, card.MCQAccuracy, 0.01, tt.description)
			assert.NotEmpty(t, card.LastMCQChoice, "Should track last MCQ choice")

			// Verify incorrect choices tracking
			incorrectAttempts := 0
			for _, attempt := range tt.attempts {
				if !attempt.correct {
					incorrectAttempts++
				}
			}
			if incorrectAttempts > 0 {
				assert.Len(t, card.IncorrectChoices, incorrectAttempts, "Should track unique incorrect choices")
			}
		})
	}
}

// TestMCQIncorrectChoicePattern tests that incorrect answer patterns are tracked
func TestMCQIncorrectChoicePattern(t *testing.T) {
	card := NewCard("pattern-test-001", "raft", "L5", "baseline")
	card.QuestionType = "mcq"

	// Simulate common incorrect answer patterns
	attempts := []mcqAttempt{
		{selected: "A", correct: false},
		{selected: "A", correct: false}, // Same wrong answer
		{selected: "B", correct: false}, // Different wrong answer
		{selected: "A", correct: false}, // Repeated wrong answer
		{selected: "C", correct: true},  // Finally correct
	}

	for _, attempt := range attempts {
		card.RecordMCQAttempt(attempt.selected, attempt.correct)
		card.RecordReview(QualityGood, 25, 0)
	}

	// Verify incorrect choice tracking
	assert.Contains(t, card.IncorrectChoices, "A", "Should track A as incorrect")
	assert.Contains(t, card.IncorrectChoices, "B", "Should track B as incorrect")
	assert.Len(t, card.IncorrectChoices, 2, "Should track unique incorrect choices only")

	// Verify final state - calculate expected accuracy manually
	expectedAccuracy := 1.0 / 5.0 // 1 correct out of 5 attempts
	assert.Equal(t, "C", card.LastMCQChoice, "Last choice should be correct answer")
	assert.InDelta(t, expectedAccuracy, card.MCQAccuracy, 0.01, "Should calculate 20% accuracy (1/5)")
}

// mcqAttempt represents a single MCQ attempt for testing
type mcqAttempt struct {
	selected string
	correct  bool
	quality  ReviewQuality
}

// TestMCQSRSSchedulerIntegration tests MCQ with SRS scheduler
func TestMCQSRSSchedulerIntegration(t *testing.T) {
	// Create temporary directory for test storage
	tmpDir := t.TempDir()
	
	scheduler, err := NewScheduler(tmpDir)
	require.NoError(t, err)

	// Create MCQ card
	card := NewCard("scheduler-test-001", "zookeeper", "L4", "baseline")
	card.QuestionType = "mcq"

	// Add card to scheduler (need to add to cards map directly since AddCard expects models.Question)
	scheduler.cards[card.QuestionID] = card
	err = scheduler.storage.SaveCards(scheduler.cards)
	require.NoError(t, err)

	// Test correct MCQ answer  
	updatedCard, exists := scheduler.GetCard(card.QuestionID)
	require.True(t, exists)
	initialReviews := updatedCard.TotalReviews

	result, err := scheduler.RecordReview(card.QuestionID, QualityGood, 20, 0)
	require.NoError(t, err)

	// Verify SRS update
	assert.Greater(t, result.NewInterval, result.PreviousInterval, "Correct answer should increase interval")
	assert.Equal(t, QualityGood, result.Quality)

	// Verify MCQ-specific tracking
	finalCard, exists := scheduler.GetCard(card.QuestionID)
	require.True(t, exists)
	assert.Equal(t, "mcq", finalCard.QuestionType, "Should maintain MCQ type")
	assert.Greater(t, finalCard.TotalReviews, initialReviews, "Should increment total reviews")
}

// TestMCQQualityMapping tests how different qualities map to MCQ performance
func TestMCQQualityMapping(t *testing.T) {
	tests := []struct {
		name          string
		mcqScenario   string
		expectedQuality ReviewQuality
	}{
		{
			name:        "Quick correct MCQ",
			mcqScenario: "correct in <10 seconds",
			expectedQuality: QualityEasy,
		},
		{
			name:        "Normal correct MCQ",  
			mcqScenario: "correct in 30 seconds",
			expectedQuality: QualityGood,
		},
		{
			name:        "Slow correct MCQ",
			mcqScenario: "correct in 60+ seconds",
			expectedQuality: QualityHard,
		},
		{
			name:        "Incorrect MCQ",
			mcqScenario: "wrong answer selected",
			expectedQuality: QualityWrong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card := NewCard("quality-test-001", "gfs", "L3", "baseline")
			card.QuestionType = "mcq"

			// Simulate quality-based review
			timeSpent := simulateTimeSpent(tt.mcqScenario)
			correct := !containsString([]string{"wrong answer selected"}, tt.mcqScenario)

			// Record review with mapped quality
			card.RecordReview(tt.expectedQuality, timeSpent, 0)
			card.RecordMCQAttempt("A", correct)

			// Verify quality was recorded
			assert.Equal(t, 1, card.TotalReviews)
			if tt.expectedQuality >= QualityGood {
				assert.Equal(t, 1, card.SuccessCount, "Correct answers should count as success")
			}
		})
	}
}

// simulateTimeSpent converts scenario description to time in seconds
func simulateTimeSpent(scenario string) int {
	switch {
	case containsString([]string{"<10 seconds", "fast", "quick"}, scenario):
		return 8
	case containsString([]string{"30 seconds", "normal"}, scenario):
		return 30
	case containsString([]string{"60+ seconds", "slow", "timed out"}, scenario):
		return 75
	default:
		return 25
	}
}



// TestMCQLearningCurve tests how MCQ performance affects learning progression
func TestMCQLearningCurve(t *testing.T) {
	scenarios := []struct {
		name           string
		attempts       consecutiveAttempts
		expectedState  CardState
		expectedMinInterval int
	}{
		{
			name: "Correct streak advances beyond new",
			attempts: consecutiveAttempts{
				count:       6,
				correctRate: 1.0, // 100% correct
				avgTime:     20,  // Quick answers
			},
			expectedState: CardStateLearning, // Might not reach mature in 6 attempts
			expectedMinInterval: 1,  // Should have some interval growth
		},
		{
			name: "Mixed performance stays in learning",
			attempts: consecutiveAttempts{
				count:       6,
				correctRate: 0.6, // 60% correct
				avgTime:     35,  // Slower answers
			},
			expectedState: CardStateLearning,
			expectedMinInterval: 0,  // Learning cards might have 0 interval if struggling
		},
		{
			name: "Poor performance keeps in learning",
			attempts: consecutiveAttempts{
				count:       6,
				correctRate: 0.3, // 30% correct
				avgTime:     45,  // Slow answers
			},
			expectedState: CardStateLearning,
			expectedMinInterval: 0,  // Still in learning phase with poor performance
		},
	}

	for _, tc := range scenarios {
		t.Run(tc.name, func(t *testing.T) {
			// Create temporary directory for SRS
			tmpDir := t.TempDir()
			
			card := NewCard("learning-test-001", "gfs", "L3", "baseline")
			card.QuestionType = "mcq"

			scheduler, err := NewScheduler(tmpDir)
			require.NoError(t, err)
			scheduler.cards[card.QuestionID] = card

			// Simulate the learning attempts using scheduler
			for i := 0; i < tc.attempts.count; i++ {
				isCorrect := i < int(float64(tc.attempts.count)*tc.attempts.correctRate)
				quality := QualityHard
				if isCorrect {
					if tc.attempts.avgTime < 25 {
						quality = QualityGood
					} else {
						quality = QualityHard
					}
				} else {
					quality = QualityWrong
				}

				// Use scheduler to apply proper SRS algorithm
				_, err := scheduler.RecordReview(card.QuestionID, quality, tc.attempts.avgTime, 0)
				require.NoError(t, err)
			}

			// Get final card state
			finalCard, exists := scheduler.GetCard(card.QuestionID)
			require.True(t, exists)

			// Verify learning progression
			assert.Equal(t, tc.expectedState, finalCard.State, tc.name)
			assert.GreaterOrEqual(t, finalCard.Interval, tc.expectedMinInterval, 
				"Interval should be at least %d for state %s", tc.expectedMinInterval, tc.expectedState)
		})
	}
}

// consecutiveAttempts represents a series of MCQ attempts for learning curve testing
type consecutiveAttempts struct {
	count       int
	correctRate float64
	avgTime     int
}
