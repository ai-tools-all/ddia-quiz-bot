package quiz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// ValidationReport represents the complete validation report
type ValidationReport struct {
	TotalFiles        int              `json:"total_files"`
	ValidFiles        int              `json:"valid_files"`
	InvalidFiles      int              `json:"invalid_files"`
	FilesWithWarnings int              `json:"files_with_warnings"`
	FileValidations   []FileValidation `json:"results"`
}

// Reporter generates validation reports
type Reporter struct {
	verbose bool
	quiet   bool
}

// NewReporter creates a new reporter
func NewReporter(verbose, quiet bool) *Reporter {
	return &Reporter{
		verbose: verbose,
		quiet:   quiet,
	}
}

// Generate creates a validation report from file validation results
func (r *Reporter) Generate(validations []FileValidation) ValidationReport {
	report := ValidationReport{
		TotalFiles:      len(validations),
		FileValidations: validations,
	}

	for _, validation := range validations {
		if validation.Valid {
			report.ValidFiles++
		} else {
			report.InvalidFiles++
		}

		// Count warnings
		hasWarnings := false
		for _, issue := range validation.Issues {
			if issue.Type == "warning" {
				hasWarnings = true
				break
			}
		}
		if hasWarnings {
			report.FilesWithWarnings++
		}
	}

	return report
}

// OutputText generates a human-readable text report
func (r *Reporter) OutputText(w io.Writer, report ValidationReport) error {
	var buf bytes.Buffer

	// Summary section
	if !r.quiet {
		buf.WriteString("Validation Report\n")
		buf.WriteString("=================\n")
		fmt.Fprintf(&buf, "Total files:         %d\n", report.TotalFiles)
		fmt.Fprintf(&buf, "Valid files:         %d\n", report.ValidFiles)
		fmt.Fprintf(&buf, "Invalid files:       %d\n", report.InvalidFiles)
		fmt.Fprintf(&buf, "Files with warnings: %d\n", report.FilesWithWarnings)
		buf.WriteString("\n")
	}

	// Detailed results
	for _, validation := range report.FileValidations {
		// Skip valid files without warnings in quiet mode
		if r.quiet && validation.Valid && !r.hasWarnings(validation) {
			continue
		}

		// File header
		if validation.Valid {
			if r.hasWarnings(validation) {
				fmt.Fprintf(&buf, "⚠ %s (VALID with warnings)\n", validation.FilePath)
			} else if !r.quiet {
				fmt.Fprintf(&buf, "✓ %s (VALID)\n", validation.FilePath)
			}
		} else {
			fmt.Fprintf(&buf, "✗ %s (INVALID)\n", validation.FilePath)
		}

		// Show issues
		for _, issue := range validation.Issues {
			// Skip info messages unless verbose
			if issue.Type == "info" && !r.verbose {
				continue
			}

			icon := r.getIssueIcon(issue.Type)
			label := r.getIssueLabel(issue.Type)

			if r.verbose || !validation.Valid || issue.Type == "error" {
				fmt.Fprintf(&buf, "  %s %s: %s\n", icon, label, issue.Message)
			} else if issue.Type == "warning" {
				fmt.Fprintf(&buf, "  %s %s: %s\n", icon, label, issue.Message)
			}
		}

		// Show question details in verbose mode
		if r.verbose && validation.Question != nil {
			q := validation.Question
			buf.WriteString("  Details:\n")
			fmt.Fprintf(&buf, "    ID: %s\n", q.ID)
			fmt.Fprintf(&buf, "    Title: %s\n", q.Title)
			fmt.Fprintf(&buf, "    Level: %s\n", q.Level)
			fmt.Fprintf(&buf, "    Category: %s\n", q.Category)
			fmt.Fprintf(&buf, "    Core Concepts: %d items\n", len(q.CoreConcepts))
			fmt.Fprintf(&buf, "    Peripheral Concepts: %d items\n", len(q.PeripheralConcepts))
			fmt.Fprintf(&buf, "    Has Sample Excellent: %v\n", q.SampleExcellent != "")
			fmt.Fprintf(&buf, "    Has Sample Acceptable: %v\n", q.SampleAcceptable != "")
			fmt.Fprintf(&buf, "    Evaluation Rubric Items: %d\n", len(q.EvaluationRubric))
		}

		// Add spacing
		if !r.quiet || !validation.Valid || r.hasWarnings(validation) {
			buf.WriteString("\n")
		}
	}

	_, err := io.Copy(w, &buf)
	return err
}

// OutputJSON generates a JSON report
func (r *Reporter) OutputJSON(w io.Writer, report ValidationReport) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

// hasWarnings checks if a validation result has any warnings
func (r *Reporter) hasWarnings(validation FileValidation) bool {
	for _, issue := range validation.Issues {
		if issue.Type == "warning" {
			return true
		}
	}
	return false
}

// getIssueIcon returns the appropriate icon for an issue type
func (r *Reporter) getIssueIcon(issueType string) string {
	switch issueType {
	case "error":
		return "✗"
	case "warning":
		return "⚠"
	case "info":
		return "ℹ"
	default:
		return "•"
	}
}

// getIssueLabel returns the label for an issue type
func (r *Reporter) getIssueLabel(issueType string) string {
	switch issueType {
	case "error":
		return "ERROR"
	case "warning":
		return "WARNING"
	case "info":
		return "INFO"
	default:
		return issueType
	}
}
