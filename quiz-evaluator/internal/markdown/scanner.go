package markdown

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/abhishek/quiz-evaluator/internal/models"
)

// Scanner recursively scans directories for markdown files
type Scanner struct {
	basePath string
	parser   *Parser
}

// NewScanner creates a new markdown scanner
func NewScanner(basePath string) *Scanner {
	return &Scanner{
		basePath: basePath,
		parser:   NewParser(),
	}
}

// ScanQuestions scans the directory for markdown question files and builds an index
func (s *Scanner) ScanQuestions() (models.QuestionIndex, error) {
	index := make(models.QuestionIndex)
	
	err := filepath.Walk(s.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Process only markdown files
		if !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}

		// Skip index/readme files
		filename := strings.ToLower(info.Name())
		if filename == "readme.md" || filename == "index.md" {
			return nil
		}

		// Parse the markdown file
		question, err := s.parser.ParseQuestionFile(path)
		if err != nil {
			// Log error but continue scanning
			fmt.Fprintf(os.Stderr, "Warning: Failed to parse %s: %v\n", path, err)
			return nil
		}

		// Skip if not a valid question file
		if question == nil || question.ID == "" {
			return nil
		}

		// Add to index
		if existing, exists := index[question.ID]; exists {
			fmt.Fprintf(os.Stderr, "Warning: Duplicate question ID '%s' found in:\n  - %s\n  - %s\n", 
				question.ID, existing.FilePath, path)
			return nil
		}

		question.FilePath = path
		index[question.ID] = question
		
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error scanning questions directory: %w", err)
	}

	if len(index) == 0 {
		return nil, fmt.Errorf("no valid question files found in %s", s.basePath)
	}

	return index, nil
}

// CountQuestions returns the number of questions found
func (s *Scanner) CountQuestions(index models.QuestionIndex) int {
	return len(index)
}

// GetQuestionsByLevel groups questions by difficulty level
func (s *Scanner) GetQuestionsByLevel(index models.QuestionIndex) map[string][]*models.Question {
	byLevel := make(map[string][]*models.Question)
	
	for _, question := range index {
		level := question.Level
		if level == "" {
			level = "Unknown"
		}
		byLevel[level] = append(byLevel[level], question)
	}
	
	return byLevel
}
