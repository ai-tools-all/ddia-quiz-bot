package quiz

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/abhishek/ddia-clicker/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	validations := []FileValidation{
		{
			FilePath: "test1.md",
			Valid:    true,
			Issues:   []ValidationIssue{},
		},
		{
			FilePath: "test2.md",
			Valid:    false,
			Issues: []ValidationIssue{
				{Type: "error", Field: "id", Message: "missing ID"},
			},
		},
		{
			FilePath: "test3.md",
			Valid:    true,
			Issues: []ValidationIssue{
				{Type: "warning", Field: "title", Message: "missing title"},
			},
		},
	}

	reporter := NewReporter(false, false)
	report := reporter.Generate(validations)

	assert.Equal(t, 3, report.TotalFiles)
	assert.Equal(t, 2, report.ValidFiles)
	assert.Equal(t, 1, report.InvalidFiles)
	assert.Equal(t, 1, report.FilesWithWarnings)
}

func TestOutputText(t *testing.T) {
	question := &models.Question{
		ID:                 "test-001",
		Title:              "Test Question",
		Level:              "L3",
		Category:           "baseline",
		CoreConcepts:       []string{"concept1"},
		PeripheralConcepts: []string{"concept2"},
		SampleExcellent:    "Excellent",
		SampleAcceptable:   "Acceptable",
		EvaluationRubric:   map[string]string{"criterion": "description"},
	}

	tests := []struct {
		name        string
		verbose     bool
		quiet       bool
		validations []FileValidation
		checkOutput func(t *testing.T, output string)
	}{
		{
			name:    "normal output with summary",
			verbose: false,
			quiet:   false,
			validations: []FileValidation{
				{FilePath: "test1.md", Valid: true, Issues: []ValidationIssue{}},
				{FilePath: "test2.md", Valid: false, Issues: []ValidationIssue{
					{Type: "error", Field: "id", Message: "missing ID"},
				}},
			},
			checkOutput: func(t *testing.T, output string) {
				assert.Contains(t, output, "Validation Report")
				assert.Contains(t, output, "Total files:")
				assert.Contains(t, output, "✓ test1.md")
				assert.Contains(t, output, "✗ test2.md")
				assert.Contains(t, output, "ERROR")
			},
		},
		{
			name:    "quiet mode",
			verbose: false,
			quiet:   true,
			validations: []FileValidation{
				{FilePath: "test1.md", Valid: true, Issues: []ValidationIssue{}},
				{FilePath: "test2.md", Valid: false, Issues: []ValidationIssue{
					{Type: "error", Field: "id", Message: "missing ID"},
				}},
			},
			checkOutput: func(t *testing.T, output string) {
				assert.NotContains(t, output, "Validation Report")
				assert.NotContains(t, output, "✓ test1.md")
				assert.Contains(t, output, "✗ test2.md")
			},
		},
		{
			name:    "verbose mode with details",
			verbose: true,
			quiet:   false,
			validations: []FileValidation{
				{FilePath: "test1.md", Valid: true, Issues: []ValidationIssue{}, Question: question},
			},
			checkOutput: func(t *testing.T, output string) {
				assert.Contains(t, output, "Details:")
				assert.Contains(t, output, "ID: test-001")
				assert.Contains(t, output, "Level: L3")
			},
		},
		{
			name:    "warnings display",
			verbose: false,
			quiet:   false,
			validations: []FileValidation{
				{FilePath: "test1.md", Valid: true, Issues: []ValidationIssue{
					{Type: "warning", Field: "title", Message: "missing title"},
				}},
			},
			checkOutput: func(t *testing.T, output string) {
				assert.Contains(t, output, "⚠")
				assert.Contains(t, output, "WARNING")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reporter := NewReporter(tt.verbose, tt.quiet)
			report := reporter.Generate(tt.validations)

			var buf bytes.Buffer
			err := reporter.OutputText(&buf, report)
			require.NoError(t, err)

			output := buf.String()
			tt.checkOutput(t, output)
		})
	}
}

func TestOutputJSON(t *testing.T) {
	validations := []FileValidation{
		{
			FilePath: "test1.md",
			Valid:    true,
			Issues:   []ValidationIssue{},
			Question: &models.Question{
				ID:    "test-001",
				Title: "Test Question",
			},
		},
		{
			FilePath: "test2.md",
			Valid:    false,
			Issues: []ValidationIssue{
				{Type: "error", Field: "id", Message: "missing ID"},
				{Type: "warning", Field: "title", Message: "missing title"},
			},
		},
	}

	reporter := NewReporter(false, false)
	report := reporter.Generate(validations)

	var buf bytes.Buffer
	err := reporter.OutputJSON(&buf, report)
	require.NoError(t, err)

	// Parse JSON to verify structure
	var decoded ValidationReport
	err = json.Unmarshal(buf.Bytes(), &decoded)
	require.NoError(t, err)

	assert.Equal(t, 2, decoded.TotalFiles)
	assert.Equal(t, 1, decoded.ValidFiles)
	assert.Equal(t, 1, decoded.InvalidFiles)
	assert.Equal(t, 1, decoded.FilesWithWarnings)
	assert.Len(t, decoded.FileValidations, 2)

	// Check first validation
	assert.Equal(t, "test1.md", decoded.FileValidations[0].FilePath)
	assert.True(t, decoded.FileValidations[0].Valid)
	assert.NotNil(t, decoded.FileValidations[0].Question)

	// Check second validation
	assert.Equal(t, "test2.md", decoded.FileValidations[1].FilePath)
	assert.False(t, decoded.FileValidations[1].Valid)
	assert.Len(t, decoded.FileValidations[1].Issues, 2)
}

func TestGetIssueIcon(t *testing.T) {
	reporter := NewReporter(false, false)

	tests := []struct {
		issueType string
		expected  string
	}{
		{"error", "✗"},
		{"warning", "⚠"},
		{"info", "ℹ"},
		{"unknown", "•"},
	}

	for _, tt := range tests {
		t.Run(tt.issueType, func(t *testing.T) {
			icon := reporter.getIssueIcon(tt.issueType)
			assert.Equal(t, tt.expected, icon)
		})
	}
}

func TestGetIssueLabel(t *testing.T) {
	reporter := NewReporter(false, false)

	tests := []struct {
		issueType string
		expected  string
	}{
		{"error", "ERROR"},
		{"warning", "WARNING"},
		{"info", "INFO"},
		{"custom", "custom"},
	}

	for _, tt := range tests {
		t.Run(tt.issueType, func(t *testing.T) {
			label := reporter.getIssueLabel(tt.issueType)
			assert.Equal(t, tt.expected, label)
		})
	}
}

func TestHasWarnings(t *testing.T) {
	reporter := NewReporter(false, false)

	tests := []struct {
		name       string
		validation FileValidation
		expected   bool
	}{
		{
			name: "has warnings",
			validation: FileValidation{
				Issues: []ValidationIssue{
					{Type: "warning", Field: "test", Message: "test"},
				},
			},
			expected: true,
		},
		{
			name: "only errors",
			validation: FileValidation{
				Issues: []ValidationIssue{
					{Type: "error", Field: "test", Message: "test"},
				},
			},
			expected: false,
		},
		{
			name: "no issues",
			validation: FileValidation{
				Issues: []ValidationIssue{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := reporter.hasWarnings(tt.validation)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestQuietModeFiltering(t *testing.T) {
	validations := []FileValidation{
		{FilePath: "valid-no-warnings.md", Valid: true, Issues: []ValidationIssue{}},
		{FilePath: "valid-with-warnings.md", Valid: true, Issues: []ValidationIssue{
			{Type: "warning", Field: "test", Message: "test warning"},
		}},
		{FilePath: "invalid.md", Valid: false, Issues: []ValidationIssue{
			{Type: "error", Field: "test", Message: "test error"},
		}},
	}

	reporter := NewReporter(false, true) // quiet mode
	report := reporter.Generate(validations)

	var buf bytes.Buffer
	err := reporter.OutputText(&buf, report)
	require.NoError(t, err)

	output := buf.String()

	// In quiet mode, valid file without warnings should not appear
	assert.NotContains(t, output, "valid-no-warnings.md")

	// Valid file with warnings should appear
	assert.Contains(t, output, "valid-with-warnings.md")

	// Invalid file should appear
	assert.Contains(t, output, "invalid.md")
}

func TestReportSummarySection(t *testing.T) {
	validations := []FileValidation{
		{FilePath: "test1.md", Valid: true, Issues: []ValidationIssue{}},
		{FilePath: "test2.md", Valid: false, Issues: []ValidationIssue{{Type: "error"}}},
		{FilePath: "test3.md", Valid: true, Issues: []ValidationIssue{{Type: "warning"}}},
	}

	reporter := NewReporter(false, false)
	report := reporter.Generate(validations)

	var buf bytes.Buffer
	err := reporter.OutputText(&buf, report)
	require.NoError(t, err)

	output := buf.String()
	lines := strings.Split(output, "\n")

	// Find summary section
	var summaryLines []string
	inSummary := false
	for _, line := range lines {
		if strings.Contains(line, "Validation Report") {
			inSummary = true
		}
		if inSummary {
			summaryLines = append(summaryLines, line)
			if line == "" && len(summaryLines) > 1 {
				break
			}
		}
	}

	summaryText := strings.Join(summaryLines, "\n")
	assert.Contains(t, summaryText, "Total files:         3")
	assert.Contains(t, summaryText, "Valid files:         2")
	assert.Contains(t, summaryText, "Invalid files:       1")
	assert.Contains(t, summaryText, "Files with warnings: 1")
}
