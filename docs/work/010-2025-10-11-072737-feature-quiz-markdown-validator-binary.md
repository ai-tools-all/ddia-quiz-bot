# Quiz Markdown Validator Binary

**Date**: 2025-10-11  
**Type**: Feature  
**Status**: In Progress

## 1. Overview

Create a standalone binary `validate-quiz` that validates markdown quiz files for TUI compatibility, sharing the parsing logic from the existing TUI application. This tool will help content creators verify their quiz questions are properly formatted before integration.

## 2. Requirements

### 2.1 Functional Requirements
- Validate single markdown files or entire directories
- Check YAML frontmatter syntax and required fields
- Verify markdown structure and section formatting
- Support multiple output formats (text, JSON)
- Provide detailed error messages and warnings
- Exit with appropriate status codes for CI/CD integration

### 2.2 Non-Functional Requirements
- Share parsing logic with TUI to ensure consistency
- Fast performance for large directory scans
- Clear, actionable error messages
- Support for recursive directory traversal

## 3. Architecture Design

### 3.1 Component Structure

```
cmd/validate-quiz/
├── main.go                 # CLI entry point and command handling

internal/markdown/          # Shared parsing logic (already exists)
├── parser.go              # Core parsing functionality
├── scanner.go             # Directory scanning and indexing

internal/quiz/             # New shared quiz validation logic
├── validator.go           # Validation rules and checks
├── reporter.go            # Report generation and formatting
```

### 3.2 Shared Components

The validator will reuse existing components from the TUI:
- `internal/markdown/Parser` - Markdown parsing with frontmatter
- `internal/markdown/Scanner` - Directory scanning
- `internal/models/Question` - Question data model

## 4. Detailed Implementation Plan

### 4.1 Command-Line Interface

```go
// cmd/validate-quiz/main.go
var rootCmd = &cobra.Command{
    Use:   "validate-quiz [file/directory]",
    Short: "Validate markdown quiz files for TUI compatibility",
    Args:  cobra.ExactArgs(1),
    RunE:  runValidation,
}

// CLI flags
func init() {
    rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "text", 
        "Output format (text, json)")
    rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, 
        "Show detailed validation information")
    rootCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, 
        "Recursively validate directories")
    rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, 
        "Only show files with errors")
}
```

### 4.2 Validation Logic

```go
// internal/quiz/validator.go
package quiz

import (
    "github.com/abhishek/ddia-clicker/internal/models"
)

type Validator struct {
    strictMode bool  // Fail on warnings
}

type ValidationIssue struct {
    Type     string // "error" or "warning"
    Field    string // Which field has the issue
    Message  string // Human-readable message
}

func (v *Validator) ValidateQuestion(q *models.Question) []ValidationIssue {
    var issues []ValidationIssue
    
    // Required fields
    if q.ID == "" {
        issues = append(issues, ValidationIssue{
            Type:    "error",
            Field:   "id",
            Message: "Question ID is required",
        })
    }
    
    if q.MainQuestion == "" {
        issues = append(issues, ValidationIssue{
            Type:    "error",
            Field:   "question",
            Message: "Main question content is required",
        })
    }
    
    // Recommended fields
    if q.Level == "" {
        issues = append(issues, ValidationIssue{
            Type:    "warning",
            Field:   "level",
            Message: "Difficulty level (L3-L7) is recommended",
        })
    } else if !isValidLevel(q.Level) {
        issues = append(issues, ValidationIssue{
            Type:    "warning",
            Field:   "level",
            Message: fmt.Sprintf("Unusual level '%s', expected L3-L7", q.Level),
        })
    }
    
    // Validate category
    if q.Category != "" && !isValidCategory(q.Category) {
        issues = append(issues, ValidationIssue{
            Type:    "warning",
            Field:   "category",
            Message: fmt.Sprintf("Unusual category '%s', expected 'baseline' or 'bar-raiser'", q.Category),
        })
    }
    
    // Check content completeness
    if len(q.CoreConcepts) == 0 {
        issues = append(issues, ValidationIssue{
            Type:    "warning",
            Field:   "core_concepts",
            Message: "No core concepts defined",
        })
    }
    
    return issues
}
```

### 4.3 Report Generation

```go
// internal/quiz/reporter.go
package quiz

type FileValidation struct {
    FilePath string
    Valid    bool
    Issues   []ValidationIssue
    Question *models.Question
}

type ValidationReport struct {
    TotalFiles         int
    ValidFiles         int
    InvalidFiles       int
    FilesWithWarnings  int
    FileValidations    []FileValidation
}

type Reporter interface {
    Generate(validations []FileValidation) ValidationReport
    OutputText(report ValidationReport, verbose bool) string
    OutputJSON(report ValidationReport) ([]byte, error)
}

type DefaultReporter struct{}

func (r *DefaultReporter) OutputText(report ValidationReport, verbose bool) string {
    var buf bytes.Buffer
    
    // Summary
    fmt.Fprintf(&buf, "Validation Report\n")
    fmt.Fprintf(&buf, "=================\n")
    fmt.Fprintf(&buf, "Total files:         %d\n", report.TotalFiles)
    fmt.Fprintf(&buf, "Valid files:         %d\n", report.ValidFiles)
    fmt.Fprintf(&buf, "Invalid files:       %d\n", report.InvalidFiles)
    fmt.Fprintf(&buf, "Files with warnings: %d\n", report.FilesWithWarnings)
    
    // Details for each file
    for _, fv := range report.FileValidations {
        if fv.Valid {
            fmt.Fprintf(&buf, "✓ %s\n", fv.FilePath)
        } else {
            fmt.Fprintf(&buf, "✗ %s\n", fv.FilePath)
        }
        
        for _, issue := range fv.Issues {
            icon := "⚠"
            if issue.Type == "error" {
                icon = "✗"
            }
            fmt.Fprintf(&buf, "  %s %s: %s\n", icon, issue.Field, issue.Message)
        }
    }
    
    return buf.String()
}
```

### 4.4 Integration with Existing Parser

The validator will use the existing parser with enhanced error handling:

```go
// Enhance internal/markdown/parser.go with validation support
func (p *Parser) ParseQuestionFileWithValidation(filepath string) (*models.Question, []ValidationIssue, error) {
    content, err := os.ReadFile(filepath)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to read file: %w", err)
    }
    
    var issues []ValidationIssue
    
    // Extract frontmatter
    frontmatter, body := p.extractFrontmatter(string(content))
    if frontmatter == "" {
        issues = append(issues, ValidationIssue{
            Type:    "warning",
            Field:   "frontmatter",
            Message: "No YAML frontmatter found",
        })
    }
    
    // Continue parsing...
    question := &models.Question{}
    // ... existing parsing logic ...
    
    return question, issues, nil
}
```

## 5. Validation Rules

### 5.1 Required Fields (Errors)
- `id` or `question_id` - Unique identifier
- `question` or `main_question` - The actual question text

### 5.2 Recommended Fields (Warnings)
- `title` - Question title
- `level` - Difficulty level (L3-L7)
- `category` - Type (baseline/bar-raiser)
- `core_concepts` - List of core concepts
- `peripheral_concepts` - List of peripheral concepts
- `sample_excellent` - Example excellent answer
- `sample_acceptable` - Example acceptable answer
- `evaluation_rubric` - Scoring criteria

### 5.3 Format Validation
- Level must match pattern: L3, L4, L5, L6, L7
- Category must be: baseline, bar-raiser, or bar_raiser
- Lists must use proper markdown formatting (-, *, or numbered)
- Sections should use ## headers

## 6. Usage Examples

### 6.1 Basic Usage

```bash
# Validate single file
./validate-quiz question.md

# Validate directory
./validate-quiz ./questions

# Recursive validation
./validate-quiz -r ./content

# JSON output for CI/CD
./validate-quiz -f json ./content > validation-report.json

# Quiet mode - only errors
./validate-quiz -q ./content

# Verbose mode - full details
./validate-quiz -v ./content
```

### 6.2 Expected Output

#### Text Output (Default)
```
Validation Report
=================
Total files:         17
Valid files:         15
Invalid files:       1
Files with warnings: 3

✓ L3-baseline/01-linearizability-basics.md
✓ L3-baseline/02-fifo-guarantees.md
⚠ L3-baseline/03-watch-mechanism.md
  ⚠ level: Difficulty level (L3-L7) is recommended
✗ L4-baseline/01-configuration-management.md
  ✗ id: Question ID is required
  ⚠ core_concepts: No core concepts defined
```

#### JSON Output
```json
{
  "total_files": 17,
  "valid_files": 15,
  "invalid_files": 1,
  "files_with_warnings": 3,
  "results": [
    {
      "filepath": "L3-baseline/01-linearizability-basics.md",
      "valid": true,
      "errors": [],
      "warnings": [],
      "question": {
        "id": "zk-l3-linearizability",
        "title": "Understanding Linearizability",
        "level": "L3",
        "category": "baseline"
      }
    }
  ]
}
```

## 7. Testing Strategy

### 7.1 Unit Tests

```go
// cmd/validate-quiz/main_test.go
func TestValidateFile(t *testing.T) {
    tests := []struct {
        name     string
        content  string
        wantValid bool
        wantErrors int
        wantWarnings int
    }{
        {
            name: "valid question",
            content: `---
id: test-001
level: L3
category: baseline
---
## Question
What is linearizability?`,
            wantValid: true,
            wantErrors: 0,
            wantWarnings: 0,
        },
        {
            name: "missing id",
            content: `## Question
What is linearizability?`,
            wantValid: false,
            wantErrors: 1,
            wantWarnings: 0,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create temp file
            tmpfile := createTempMarkdown(t, tt.content)
            defer os.Remove(tmpfile)
            
            result := validateFile(tmpfile)
            assert.Equal(t, tt.wantValid, result.Valid)
            assert.Len(t, result.Errors, tt.wantErrors)
            assert.Len(t, result.Warnings, tt.wantWarnings)
        })
    }
}
```

### 7.2 Integration Tests

Test against actual quiz files:
```bash
# Test on existing Zookeeper questions
./validate-quiz -r ddia-quiz-bot/content/chapters/10-mit-6824-zookeeper/

# Test on GFS questions
./validate-quiz -r ddia-quiz-bot/content/chapters/09-distributed-systems-gfs/
```

## 8. Build Integration

### 8.1 Update build script

```bash
# scripts/build.sh modification
echo "Building validate-quiz..."
go build -o build/validate-quiz ./cmd/validate-quiz
```

### 8.2 Add to Makefile (if exists)

```makefile
validate-quiz:
	go build -o build/validate-quiz ./cmd/validate-quiz

validate: validate-quiz
	./build/validate-quiz -r ddia-quiz-bot/content/

.PHONY: validate-quiz validate
```

## 9. CI/CD Integration

### 9.1 GitHub Actions Example

```yaml
- name: Validate Quiz Files
  run: |
    go build -o validate-quiz ./cmd/validate-quiz
    ./validate-quiz -f json -r ddia-quiz-bot/content/ > validation.json
    
    # Check exit code
    if [ $? -ne 0 ]; then
      echo "Quiz validation failed!"
      cat validation.json | jq '.'
      exit 1
    fi
```

## 10. Implementation Steps

1. **Extract shared validation logic** (30 min)
   - Create `internal/quiz` package
   - Move validation rules from main.go
   - Create validator and reporter interfaces

2. **Enhance parser with validation hooks** (20 min)
   - Add validation issue tracking to parser
   - Return parsing issues alongside errors

3. **Implement comprehensive validation rules** (40 min)
   - Required field validation
   - Format validation (levels, categories)
   - Content completeness checks
   - Cross-reference validation

4. **Add test coverage** (30 min)
   - Unit tests for validator
   - Integration tests with sample files
   - Edge case testing

5. **Update build system** (10 min)
   - Add to build.sh
   - Create make target
   - Document in README

6. **Test on existing quiz content** (20 min)
   - Run on all existing chapters
   - Fix any discovered issues
   - Document common problems

## 11. Future Enhancements

1. **Schema validation** - Validate against JSON Schema
2. **Content analysis** - Check question quality, complexity
3. **Cross-reference checking** - Validate follow-up question references
4. **Auto-fix mode** - Automatically fix common issues
5. **Watch mode** - Continuous validation during editing
6. **VSCode extension** - Real-time validation in editor
7. **Difficulty analysis** - Verify appropriate difficulty progression

## 12. Status

- [x] Create initial validator binary structure
- [x] Extract shared validation logic to internal/quiz
- [x] Enhance parser with validation support
- [x] Add comprehensive test coverage
- [x] Test on existing quiz files
- [x] Update build script
- [ ] Document usage in README (optional)

## 13. Implementation Summary

**Date Completed**: 2025-10-11

### What Was Built

1. **internal/quiz/validator.go**
   - Comprehensive validation logic for quiz questions
   - Support for required field validation (ID, question)
   - Support for recommended field warnings (title, level, category, concepts, samples, rubric)
   - Level format validation (L3-L7)
   - Category validation (baseline, bar-raiser, bar_raiser)
   - Directory validation with recursive support
   - Strict mode option (treat warnings as errors)

2. **internal/quiz/reporter.go**
   - Text and JSON report generation
   - Verbose mode for detailed output
   - Quiet mode for showing only errors
   - Clean, formatted output with icons (✓, ✗, ⚠)
   - Summary statistics (total, valid, invalid, warnings)

3. **cmd/validate-quiz/main.go**
   - Refactored to use extracted validator and reporter packages
   - Cobra-based CLI with comprehensive flags
   - Supports single file or directory validation
   - Exit code 1 on validation failures for CI/CD integration

4. **Comprehensive Test Coverage**
   - internal/quiz/validator_test.go - 14 test functions covering:
     - Question validation rules
     - File validation
     - Directory validation (recursive and non-recursive)
     - Strict mode behavior
     - Level and category validation
   - internal/quiz/reporter_test.go - 8 test functions covering:
     - Report generation
     - Text output formatting
     - JSON output
     - Verbose and quiet modes
     - Issue icons and labels
   - All tests passing with 100% coverage of key functionality

5. **Build Integration**
   - Created scripts/build-all.sh for building all project binaries
   - Successfully builds quiz-tui (21M) and validate-quiz (4.2M)

### Testing Results

Tested on actual quiz content:
- **Zookeeper chapter**: 17 files validated, all valid with some warnings
- **Transactions chapter**: Files validated successfully
- JSON output format verified
- Quiet mode verified
- Verbose mode with detailed question information verified

### Usage Examples

```bash
# Validate single file
./build/validate-quiz question.md

# Validate directory recursively
./build/validate-quiz -r ./content

# JSON output for CI/CD
./build/validate-quiz -f json -r ./content > report.json

# Quiet mode - only show errors
./build/validate-quiz -q -r ./content

# Verbose mode with details
./build/validate-quiz -v question.md

# Strict mode - warnings fail validation
./build/validate-quiz -s -r ./content
```

### Key Features Implemented

1. ✅ Shares parsing logic with TUI (internal/markdown package)
2. ✅ Validates required and recommended fields
3. ✅ Multiple output formats (text, JSON)
4. ✅ Recursive directory scanning
5. ✅ Special file filtering (README, index, guidelines)
6. ✅ Proper exit codes for CI/CD integration
7. ✅ Comprehensive test coverage
8. ✅ Clean, maintainable code structure
9. ✅ Detailed error and warning messages
10. ✅ Strict mode for enforcing all recommendations

## 13. Notes

- The validator shares parsing logic with TUI to ensure consistency
- Exit codes: 0 = all valid, 1 = validation errors found
- JSON output is designed for easy parsing by CI/CD tools
- Warnings don't fail validation by default (use --strict for that)
