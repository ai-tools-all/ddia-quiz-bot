package evaluator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/abhishek/quiz-evaluator/internal/models"
)

// Evaluator orchestrates the evaluation process
type Evaluator struct {
	aiClient      AIClient
	questionIndex models.QuestionIndex
	config        Config
}

// Config holds evaluator configuration
type Config struct {
	ParallelWorkers int
	Timeout         time.Duration
	Verbose         bool
}

// NewEvaluator creates a new evaluator instance
func NewEvaluator(aiClient AIClient, questionIndex models.QuestionIndex, config Config) *Evaluator {
	if config.ParallelWorkers <= 0 {
		config.ParallelWorkers = 5
	}
	if config.Timeout <= 0 {
		config.Timeout = 30 * time.Second
	}

	return &Evaluator{
		aiClient:      aiClient,
		questionIndex: questionIndex,
		config:        config,
	}
}

// EvaluateResponses processes multiple user responses
func (e *Evaluator) EvaluateResponses(ctx context.Context, responses []models.UserResponse) ([]models.Evaluation, error) {
	if len(responses) == 0 {
		return []models.Evaluation{}, nil
	}

	// Channel for jobs and results
	jobs := make(chan evaluationJob, len(responses))
	results := make(chan models.Evaluation, len(responses))
	errors := make(chan error, len(responses))

	// Start worker pool
	var wg sync.WaitGroup
	for i := 0; i < e.config.ParallelWorkers; i++ {
		wg.Add(1)
		go e.evaluationWorker(ctx, jobs, results, errors, &wg)
	}

	// Queue jobs
	for _, response := range responses {
		jobs <- evaluationJob{
			response: response,
			index:    0, // Could be used for ordering if needed
		}
	}
	close(jobs)

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Collect results
	var evaluations []models.Evaluation
	var evalErrors []error

	done := false
	for !done {
		select {
		case eval, ok := <-results:
			if ok {
				evaluations = append(evaluations, eval)
			} else {
				done = true
			}
		case err, ok := <-errors:
			if ok && err != nil {
				evalErrors = append(evalErrors, err)
			}
		}
	}

	// Report errors if any
	if len(evalErrors) > 0 {
		fmt.Printf("Warning: %d evaluation(s) failed\n", len(evalErrors))
		if e.config.Verbose {
			for _, err := range evalErrors {
				fmt.Printf("  - %v\n", err)
			}
		}
	}

	return evaluations, nil
}

// evaluationJob represents a job for the worker pool
type evaluationJob struct {
	response models.UserResponse
	index    int
}

// evaluationWorker processes evaluation jobs
func (e *Evaluator) evaluationWorker(
	ctx context.Context,
	jobs <-chan evaluationJob,
	results chan<- models.Evaluation,
	errors chan<- error,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for job := range jobs {
		eval, err := e.evaluateSingleResponse(ctx, job.response)
		if err != nil {
			errors <- fmt.Errorf("failed to evaluate question %s: %w", job.response.QuestionID, err)
			
			// Send a failed evaluation with error details
			results <- models.Evaluation{
				QuestionID:   job.response.QuestionID,
				UserResponse: job.response.UserResponse,
				Score:        0,
				Feedback:     fmt.Sprintf("Evaluation failed: %v", err),
				Timestamp:    time.Now(),
			}
		} else {
			results <- eval
		}
	}
}

// evaluateSingleResponse evaluates a single user response
func (e *Evaluator) evaluateSingleResponse(ctx context.Context, response models.UserResponse) (models.Evaluation, error) {
	// Find the question
	question, exists := e.questionIndex[response.QuestionID]
	if !exists {
		return models.Evaluation{}, fmt.Errorf("question not found: %s", response.QuestionID)
	}

	// Create evaluation request
	evalReq := EvaluationRequest{
		Question:     question,
		UserResponse: response.UserResponse,
		MaxTokens:    1500,
		Temperature:  0.3,
	}

	// Apply timeout
	evalCtx, cancel := context.WithTimeout(ctx, e.config.Timeout)
	defer cancel()

	// Call AI for evaluation
	aiResponse, err := e.aiClient.EvaluateResponse(evalCtx, evalReq)
	if err != nil {
		return models.Evaluation{}, fmt.Errorf("AI evaluation failed: %w", err)
	}

	// Build evaluation result
	evaluation := models.Evaluation{
		QuestionID:    response.QuestionID,
		UserResponse:  response.UserResponse,
		Score:         aiResponse.Score,
		Feedback:      aiResponse.Feedback,
		QuestionTitle: question.Title,
		Level:         question.Level,
		Strengths:     aiResponse.Strengths,
		Improvements:  aiResponse.Improvements,
		Timestamp:     aiResponse.Timestamp,
	}

	return evaluation, nil
}

// GetStatistics returns evaluation statistics
func (e *Evaluator) GetStatistics(evaluations []models.Evaluation) Statistics {
	stats := Statistics{
		TotalEvaluated: len(evaluations),
		ScoreDistribution: make(map[string]int),
		LevelDistribution: make(map[string]int),
	}

	if len(evaluations) == 0 {
		return stats
	}

	var totalScore float64
	minScore := 100.0
	maxScore := 0.0

	for _, eval := range evaluations {
		// Update score statistics
		totalScore += eval.Score
		if eval.Score < minScore {
			minScore = eval.Score
		}
		if eval.Score > maxScore {
			maxScore = eval.Score
		}

		// Update score distribution
		bucket := getScoreBucket(eval.Score)
		stats.ScoreDistribution[bucket]++

		// Update level distribution
		if eval.Level != "" {
			stats.LevelDistribution[eval.Level]++
		}
	}

	stats.AverageScore = totalScore / float64(len(evaluations))
	stats.MinScore = minScore
	stats.MaxScore = maxScore

	return stats
}

// Statistics holds evaluation statistics
type Statistics struct {
	TotalEvaluated    int
	AverageScore      float64
	MinScore          float64
	MaxScore          float64
	ScoreDistribution map[string]int
	LevelDistribution map[string]int
}

// getScoreBucket returns the score bucket for distribution
func getScoreBucket(score float64) string {
	switch {
	case score >= 90:
		return "Excellent (90-100)"
	case score >= 80:
		return "Good (80-89)"
	case score >= 70:
		return "Satisfactory (70-79)"
	case score >= 60:
		return "Needs Improvement (60-69)"
	default:
		return "Below Expectations (<60)"
	}
}
