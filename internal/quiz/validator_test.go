package quiz

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/abhishek/ddia-clicker/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateQuestion(t *testing.T) {
	tests := []struct {
		name         string
		question     *models.Question
		wantErrors   int
		wantWarnings int
		checkIssues  func(t *testing.T, issues []ValidationIssue)
	}{
		{
			name: "valid complete question",
			question: &models.Question{
				ID:                 "test-001",
				Title:              "Test Question",
				MainQuestion:       "What is linearizability?",
				Level:              "L3",
				Category:           "baseline",
				CoreConcepts:       []string{"linearizability", "consistency"},
				PeripheralConcepts: []string{"distributed systems"},
				SampleExcellent:    "Excellent answer here",
				SampleAcceptable:   "Acceptable answer here",
				EvaluationRubric:   map[string]string{"criterion1": "description"},
			},
			wantErrors:   0,
			wantWarnings: 0,
		},
		{
			name: "missing required ID",
			question: &models.Question{
				MainQuestion: "What is linearizability?",
			},
			wantErrors:   1,
			wantWarnings: 0,
			checkIssues: func(t *testing.T, issues []ValidationIssue) {
				found := false
				for _, issue := range issues {
					if issue.Type == "error" && issue.Field == "id" {
						found = true
						assert.Contains(t, issue.Message, "required")
					}
				}
				assert.True(t, found, "Expected error about missing ID")
			},
		},
		{
			name: "missing required question",
			question: &models.Question{
				ID: "test-001",
			},
			wantErrors:   1,
			wantWarnings: 0,
			checkIssues: func(t *testing.T, issues []ValidationIssue) {
				found := false
				for _, issue := range issues {
					if issue.Type == "error" && issue.Field == "question" {
						found = true
					}
				}
				assert.True(t, found, "Expected error about missing question")
			},
		},
		{
			name: "missing recommended fields",
			question: &models.Question{
				ID:           "test-001",
				MainQuestion: "What is linearizability?",
			},
			wantErrors:   0,
			wantWarnings: 8, // title, level, category, core/peripheral concepts, samples, rubric
		},
		{
			name: "invalid level format",
			question: &models.Question{
				ID:           "test-001",
				MainQuestion: "What is linearizability?",
				Level:        "L99",
			},
			wantErrors:   0,
			wantWarnings: 1, // at least the invalid level warning
			checkIssues: func(t *testing.T, issues []ValidationIssue) {
				found := false
				for _, issue := range issues {
					if issue.Type == "warning" && issue.Field == "level" {
						found = true
						assert.Contains(t, issue.Message, "unusual")
					}
				}
				assert.True(t, found, "Expected warning about invalid level")
			},
		},
		{
			name: "invalid category",
			question: &models.Question{
				ID:           "test-001",
				MainQuestion: "What is linearizability?",
				Category:     "invalid-category",
			},
			wantErrors:   0,
			wantWarnings: 1, // at least the invalid category warning
			checkIssues: func(t *testing.T, issues []ValidationIssue) {
				found := false
				for _, issue := range issues {
					if issue.Type == "warning" && issue.Field == "category" {
						found = true
						assert.Contains(t, issue.Message, "unusual")
					}
				}
				assert.True(t, found, "Expected warning about invalid category")
			},
		},
		{
			name: "valid with bar_raiser category",
			question: &models.Question{
				ID:           "test-001",
				MainQuestion: "What is linearizability?",
				Category:     "bar_raiser",
			},
			wantErrors: 0,
		},
		{
			name: "valid with case insensitive category",
			question: &models.Question{
				ID:           "test-001",
				MainQuestion: "What is linearizability?",
				Category:     "BASELINE",
			},
			wantErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewValidator()
			issues := validator.ValidateQuestion(tt.question)

			// Count errors and warnings
			errorCount := 0
			warningCount := 0
			for _, issue := range issues {
				if issue.Type == "error" {
					errorCount++
				} else if issue.Type == "warning" {
					warningCount++
				}
			}

			assert.Equal(t, tt.wantErrors, errorCount, "Unexpected error count")
			if tt.wantWarnings > 0 {
				assert.GreaterOrEqual(t, warningCount, tt.wantWarnings, "Expected at least %d warnings", tt.wantWarnings)
			}

			if tt.checkIssues != nil {
				tt.checkIssues(t, issues)
			}
		})
	}
}

func TestValidateFile(t *testing.T) {
	// Create temp directory for test files
	tmpDir := t.TempDir()

	tests := []struct {
		name         string
		filename     string
		content      string
		wantValid    bool
		wantErrors   int
		wantWarnings int
	}{
		{
			name:     "valid question file",
			filename: "test-valid.md",
			content: `---
id: test-001
title: Test Question
level: L3
category: baseline
---

## Question
What is linearizability?

## Core Concepts
- linearizability
- consistency

## Sample Excellent Answer
This is an excellent answer.

## Sample Acceptable Answer
This is an acceptable answer.
`,
			wantValid:    true,
			wantErrors:   0,
			wantWarnings: 0,
		},
		{
			name:     "missing ID",
			filename: "test-no-id.md",
			content: `---
title: Test Question
---

## Question
What is linearizability?
`,
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name:     "non-markdown file",
			filename: "test.txt",
			content:  "This is not a markdown file",
			wantValid: false,
			wantErrors: 1,
		},
		{
			name:     "readme file should be skipped",
			filename: "README.md",
			content:  "# README",
			wantValid: true, // skipped files are considered valid
		},
		{
			name:     "minimal valid question",
			filename: "test-minimal.md",
			content: `---
id: test-minimal
---

## Question
What is the answer?
`,
			wantValid:    true,
			wantErrors:   0,
			wantWarnings: 0, // will have warnings but still valid
		},
	}

	validator := NewValidator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			filepath := filepath.Join(tmpDir, tt.filename)
			err := os.WriteFile(filepath, []byte(tt.content), 0644)
			require.NoError(t, err)

			// Validate
			result := validator.ValidateFile(filepath)

			assert.Equal(t, tt.wantValid, result.Valid, "Unexpected validation result")

			// Count errors and warnings
			errorCount := 0
			warningCount := 0
			for _, issue := range result.Issues {
				if issue.Type == "error" {
					errorCount++
				} else if issue.Type == "warning" {
					warningCount++
				}
			}

			if tt.wantErrors > 0 {
				assert.GreaterOrEqual(t, errorCount, tt.wantErrors, "Expected at least %d errors", tt.wantErrors)
			}
			if tt.wantWarnings > 0 {
				assert.GreaterOrEqual(t, warningCount, tt.wantWarnings, "Expected at least %d warnings", tt.wantWarnings)
			}
		})
	}
}

func TestValidateDirectory(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()

	// Create test files
	files := map[string]string{
		"valid1.md": `---
id: test-001
---
## Question
Test question 1?
`,
		"valid2.md": `---
id: test-002
---
## Question
Test question 2?
`,
		"invalid.md": `---
title: No ID
---
## Question
Invalid question?
`,
		"subdir/nested.md": `---
id: test-003
---
## Question
Nested question?
`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)
		err := os.MkdirAll(dir, 0755)
		require.NoError(t, err)
		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err)
	}

	validator := NewValidator()

	t.Run("non-recursive", func(t *testing.T) {
		results, err := validator.ValidateDirectory(tmpDir, false)
		require.NoError(t, err)

		// Should only find files in root directory (not subdir)
		assert.Equal(t, 3, len(results))

		validCount := 0
		invalidCount := 0
		for _, result := range results {
			if result.Valid {
				validCount++
			} else {
				invalidCount++
			}
		}

		assert.Equal(t, 2, validCount)
		assert.Equal(t, 1, invalidCount)
	})

	t.Run("recursive", func(t *testing.T) {
		results, err := validator.ValidateDirectory(tmpDir, true)
		require.NoError(t, err)

		// Should find all files including nested
		assert.Equal(t, 4, len(results))

		validCount := 0
		invalidCount := 0
		for _, result := range results {
			if result.Valid {
				validCount++
			} else {
				invalidCount++
			}
		}

		assert.Equal(t, 3, validCount)
		assert.Equal(t, 1, invalidCount)
	})
}

func TestStrictMode(t *testing.T) {
	t.Run("normal mode - warnings don't fail", func(t *testing.T) {
		validator := NewValidator()
		validator.SetStrictMode(false)

		filepath := filepath.Join(t.TempDir(), "test.md")
		content := `---
id: test-001
---
## Question
What is the answer?
`
		err := os.WriteFile(filepath, []byte(content), 0644)
		require.NoError(t, err)

		result := validator.ValidateFile(filepath)
		assert.True(t, result.Valid, "Should be valid in normal mode despite warnings")

		// Should have warnings
		hasWarnings := false
		for _, issue := range result.Issues {
			if issue.Type == "warning" {
				hasWarnings = true
				break
			}
		}
		assert.True(t, hasWarnings, "Should have warnings")
	})

	t.Run("strict mode - warnings fail validation", func(t *testing.T) {
		validator := NewValidator()
		validator.SetStrictMode(true)

		filepath := filepath.Join(t.TempDir(), "test.md")
		content := `---
id: test-001
---
## Question
What is the answer?
`
		err := os.WriteFile(filepath, []byte(content), 0644)
		require.NoError(t, err)

		result := validator.ValidateFile(filepath)
		assert.False(t, result.Valid, "Should be invalid in strict mode with warnings")
	})
}

func TestIsValidLevel(t *testing.T) {
	validator := NewValidator()

	validLevels := []string{"L3", "L4", "L5", "L6", "L7"}
	for _, level := range validLevels {
		assert.True(t, validator.isValidLevel(level), "Expected %s to be valid", level)
	}

	invalidLevels := []string{"L1", "L2", "L8", "L99", "l3", "3", ""}
	for _, level := range invalidLevels {
		assert.False(t, validator.isValidLevel(level), "Expected %s to be invalid", level)
	}
}

func TestIsValidCategory(t *testing.T) {
	validator := NewValidator()

	validCategories := []string{"baseline", "bar-raiser", "bar_raiser", "BASELINE", "Bar-Raiser"}
	for _, category := range validCategories {
		assert.True(t, validator.isValidCategory(category), "Expected %s to be valid", category)
	}

	invalidCategories := []string{"invalid", "test", ""}
	for _, category := range invalidCategories {
		assert.False(t, validator.isValidCategory(category), "Expected %s to be invalid", category)
	}
}
