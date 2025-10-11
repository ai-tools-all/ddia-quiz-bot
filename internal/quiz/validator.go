package quiz

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/abhishek/ddia-clicker/internal/markdown"
	"github.com/abhishek/ddia-clicker/internal/models"
)

// ValidationIssue represents a single validation problem
type ValidationIssue struct {
	Type    string `json:"type"`    // "error" or "warning"
	Field   string `json:"field"`   // Which field has the issue
	Message string `json:"message"` // Human-readable message
}

// FileValidation represents the validation result for a single file
type FileValidation struct {
	FilePath string            `json:"filepath"`
	Valid    bool              `json:"valid"`
	Issues   []ValidationIssue `json:"issues"`
	Question *models.Question  `json:"question,omitempty"`
}

// Validator validates quiz question files
type Validator struct {
	parser      *markdown.Parser
	strictMode  bool
	skipSpecial bool
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		parser:      markdown.NewParser(),
		strictMode:  false,
		skipSpecial: true,
	}
}

// SetStrictMode enables strict mode where warnings also fail validation
func (v *Validator) SetStrictMode(strict bool) {
	v.strictMode = strict
}

// ValidateFile validates a single markdown file
func (v *Validator) ValidateFile(filepath string) FileValidation {
	result := FileValidation{
		FilePath: filepath,
		Valid:    true,
		Issues:   []ValidationIssue{},
	}

	// Check if it's a markdown file
	if !strings.HasSuffix(strings.ToLower(filepath), ".md") {
		result.Valid = false
		result.Issues = append(result.Issues, ValidationIssue{
			Type:    "error",
			Field:   "file",
			Message: "not a markdown file",
		})
		return result
	}

	// Skip special files
	if v.skipSpecial && v.isSpecialFile(filepath) {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:    "info",
			Field:   "file",
			Message: "skipped special file (readme, index, or guidelines)",
		})
		return result
	}

	// Parse the file
	question, err := v.parser.ParseQuestionFile(filepath)
	if err != nil {
		result.Valid = false
		result.Issues = append(result.Issues, ValidationIssue{
			Type:    "error",
			Field:   "parse",
			Message: fmt.Sprintf("parse error: %v", err),
		})
		return result
	}

	result.Question = question

	// Validate the question
	issues := v.ValidateQuestion(question)
	result.Issues = append(result.Issues, issues...)

	// Check if there are any errors
	for _, issue := range result.Issues {
		if issue.Type == "error" {
			result.Valid = false
			break
		}
	}

	// In strict mode, warnings also make the file invalid
	if v.strictMode {
		for _, issue := range result.Issues {
			if issue.Type == "warning" {
				result.Valid = false
				break
			}
		}
	}

	return result
}

// ValidateQuestion validates a parsed question
func (v *Validator) ValidateQuestion(q *models.Question) []ValidationIssue {
	var issues []ValidationIssue

	// Required fields - errors
	if q.ID == "" {
		issues = append(issues, ValidationIssue{
			Type:    "error",
			Field:   "id",
			Message: "question ID is required",
		})
	}

	if q.MainQuestion == "" {
		issues = append(issues, ValidationIssue{
			Type:    "error",
			Field:   "question",
			Message: "main question content is required",
		})
	}

	// Recommended fields - warnings
	if q.Title == "" {
		issues = append(issues, ValidationIssue{
			Type:    "warning",
			Field:   "title",
			Message: "question title is recommended",
		})
	}

	if q.Level == "" {
		issues = append(issues, ValidationIssue{
			Type:    "warning",
			Field:   "level",
			Message: "difficulty level (L3-L7) is recommended",
		})
	} else if !v.isValidLevel(q.Level) {
		issues = append(issues, ValidationIssue{
			Type:    "warning",
			Field:   "level",
			Message: fmt.Sprintf("unusual level '%s', expected L3-L7", q.Level),
		})
	}

	if q.Category == "" {
		issues = append(issues, ValidationIssue{
			Type:    "warning",
			Field:   "category",
			Message: "category (baseline/bar-raiser) is recommended",
		})
	} else if !v.isValidCategory(q.Category) {
		issues = append(issues, ValidationIssue{
			Type:    "warning",
			Field:   "category",
			Message: fmt.Sprintf("unusual category '%s', expected 'baseline' or 'bar-raiser'", q.Category),
		})
	}

	// Content completeness - warnings
	if len(q.CoreConcepts) == 0 {
		issues = append(issues, ValidationIssue{
			Type:    "warning",
			Field:   "core_concepts",
			Message: "no core concepts defined",
		})
	}

	if len(q.PeripheralConcepts) == 0 {
		issues = append(issues, ValidationIssue{
			Type:    "warning",
			Field:   "peripheral_concepts",
			Message: "no peripheral concepts defined",
		})
	}

	if q.SampleExcellent == "" {
		issues = append(issues, ValidationIssue{
			Type:    "warning",
			Field:   "sample_excellent",
			Message: "sample excellent answer is recommended",
		})
	}

	if q.SampleAcceptable == "" {
		issues = append(issues, ValidationIssue{
			Type:    "warning",
			Field:   "sample_acceptable",
			Message: "sample acceptable answer is recommended",
		})
	}

	if len(q.EvaluationRubric) == 0 {
		issues = append(issues, ValidationIssue{
			Type:    "warning",
			Field:   "evaluation_rubric",
			Message: "evaluation rubric is recommended",
		})
	}

	return issues
}

// ValidateDirectory validates all markdown files in a directory
func (v *Validator) ValidateDirectory(dirPath string, recursive bool) ([]FileValidation, error) {
	var results []FileValidation

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			if !recursive && path != dirPath {
				return filepath.SkipDir
			}
			return nil
		}

		// Validate markdown files
		if strings.HasSuffix(strings.ToLower(path), ".md") {
			result := v.ValidateFile(path)
			results = append(results, result)
		}

		return nil
	}

	err := filepath.Walk(dirPath, walkFunc)
	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	return results, nil
}

// isValidLevel checks if the level format is valid
func (v *Validator) isValidLevel(level string) bool {
	validLevels := []string{"L3", "L4", "L5", "L6", "L7"}
	for _, valid := range validLevels {
		if level == valid {
			return true
		}
	}
	return false
}

// isValidCategory checks if the category is valid
func (v *Validator) isValidCategory(category string) bool {
	validCategories := []string{"baseline", "bar-raiser", "bar_raiser"}
	categoryLower := strings.ToLower(category)
	for _, valid := range validCategories {
		if categoryLower == valid {
			return true
		}
	}
	return false
}

// isSpecialFile checks if a file should be skipped (readme, index, guidelines)
func (v *Validator) isSpecialFile(filepath string) bool {
	filename := strings.ToLower(filepath)
	return strings.HasSuffix(filename, "readme.md") ||
		strings.HasSuffix(filename, "index.md") ||
		strings.HasSuffix(filename, "guidelines.md")
}
