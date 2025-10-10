package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/abhishek/quiz-evaluator/internal/models"
)

// Writer handles CSV output operations
type Writer struct {
	filepath string
}

// NewWriter creates a new CSV writer
func NewWriter(filepath string) *Writer {
	return &Writer{filepath: filepath}
}

// WriteEvaluations writes evaluation results to CSV file
func (w *Writer) WriteEvaluations(evaluations []models.Evaluation) error {
	file, err := os.Create(w.filepath)
	if err != nil {
		return fmt.Errorf("failed to create output CSV file: %w", err)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	// Write header
	header := []string{
		"question_id",
		"question_title",
		"level",
		"user_response",
		"score",
		"strengths",
		"improvements",
		"feedback",
	}

	if err := csvWriter.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, eval := range evaluations {
		record := []string{
			eval.QuestionID,
			eval.QuestionTitle,
			eval.Level,
			eval.UserResponse,
			fmt.Sprintf("%.1f", eval.Score),
			strings.Join(eval.Strengths, "; "),
			strings.Join(eval.Improvements, "; "),
			eval.Feedback,
		}

		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("failed to write evaluation record: %w", err)
		}
	}

	return nil
}

// WriteSummary writes a summary report
func (w *Writer) WriteSummary(evaluations []models.Evaluation) error {
	if len(evaluations) == 0 {
		return nil
	}

	summaryPath := strings.TrimSuffix(w.filepath, ".csv") + "_summary.txt"
	file, err := os.Create(summaryPath)
	if err != nil {
		return fmt.Errorf("failed to create summary file: %w", err)
	}
	defer file.Close()

	// Calculate statistics
	var totalScore float64
	scoreDistribution := make(map[string]int)
	for _, eval := range evaluations {
		totalScore += eval.Score
		
		// Bucket scores
		if eval.Score >= 90 {
			scoreDistribution["Excellent (90-100)"]++
		} else if eval.Score >= 80 {
			scoreDistribution["Good (80-89)"]++
		} else if eval.Score >= 70 {
			scoreDistribution["Satisfactory (70-79)"]++
		} else if eval.Score >= 60 {
			scoreDistribution["Needs Improvement (60-69)"]++
		} else {
			scoreDistribution["Below Expectations (<60)"]++
		}
	}

	avgScore := totalScore / float64(len(evaluations))

	// Write summary
	fmt.Fprintf(file, "Quiz Evaluation Summary\n")
	fmt.Fprintf(file, "========================\n\n")
	fmt.Fprintf(file, "Total Responses Evaluated: %d\n", len(evaluations))
	fmt.Fprintf(file, "Average Score: %.1f\n\n", avgScore)
	fmt.Fprintf(file, "Score Distribution:\n")
	
	for category, count := range scoreDistribution {
		percentage := float64(count) / float64(len(evaluations)) * 100
		fmt.Fprintf(file, "  %s: %d (%.1f%%)\n", category, count, percentage)
	}

	return nil
}
