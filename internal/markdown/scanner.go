package markdown

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/abhishek/ddia-clicker/internal/models"
)

// TopicInfo represents a discovered topic with metadata
type TopicInfo struct {
	Name            string         // e.g., "09-distributed-systems-gfs"
	DisplayName     string         // e.g., "GFS & Distributed Systems"
	Path            string         // absolute path to topic directory (subjective)
	LevelCounts     map[string]int // question count per level (L3, L4, etc.)
	TotalCount      int            // total question count
	MCQCount        int            // MCQ question count
	SubjectiveCount int            // Subjective question count
}

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

// DiscoverTopics scans the chapters directory and returns available topics
func (s *Scanner) DiscoverTopics(chaptersPath string) ([]TopicInfo, error) {
	entries, err := os.ReadDir(chaptersPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read chapters directory: %w", err)
	}

	var topics []TopicInfo

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		topicName := entry.Name()
		topicPath := filepath.Join(chaptersPath, topicName)
		subjectivePath := filepath.Join(topicPath, "subjective")
		mcqPath := filepath.Join(topicPath, "mcq")

		// Count questions by level
		levelCounts := make(map[string]int)
		totalCount := 0
		subjectiveCount := 0
		mcqCount := 0

		// Scan the subjective directory if it exists
		if _, err := os.Stat(subjectivePath); err == nil {
			levelDirs, err := os.ReadDir(subjectivePath)
			if err == nil {
				for _, levelDir := range levelDirs {
					if !levelDir.IsDir() {
						continue
					}

					levelName := levelDir.Name()
					// Extract level (L3, L4, etc.)
					var level string
					if strings.HasPrefix(levelName, "L") {
						// Extract L3, L4, etc. from "L3-baseline" or "L3-bar-raiser"
						parts := strings.Split(levelName, "-")
						if len(parts) > 0 {
							level = parts[0]
						}
					}

					if level == "" {
						continue
					}

					// Count markdown files in this level directory
					levelPath := filepath.Join(subjectivePath, levelName)
					files, err := os.ReadDir(levelPath)
					if err != nil {
						continue
					}

					count := 0
					for _, file := range files {
						if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".md") {
							filename := strings.ToLower(file.Name())
							if filename != "readme.md" && filename != "index.md" && filename != "guidelines.md" {
								count++
							}
						}
					}

					levelCounts[level] += count
					totalCount += count
					subjectiveCount += count
				}
			}
		}

		// Scan the MCQ directory if it exists
		if _, err := os.Stat(mcqPath); err == nil {
			files, err := os.ReadDir(mcqPath)
			if err == nil {
				for _, file := range files {
					if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".md") {
						filename := strings.ToLower(file.Name())
						if filename != "readme.md" && filename != "index.md" && filename != "guidelines.md" {
							mcqCount++
							totalCount++
						}
					}
				}
			}
		}

		if totalCount == 0 {
			continue
		}

		// Create display name from topic name
		displayName := formatTopicName(topicName)

		topics = append(topics, TopicInfo{
			Name:            topicName,
			DisplayName:     displayName,
			Path:            subjectivePath,
			LevelCounts:     levelCounts,
			TotalCount:      totalCount,
			MCQCount:        mcqCount,
			SubjectiveCount: subjectiveCount,
		})
	}

	if len(topics) == 0 {
		return nil, fmt.Errorf("no valid topics found in %s", chaptersPath)
	}

	// Sort topics by name for consistent ordering
	sort.Slice(topics, func(i, j int) bool {
		return topics[i].Name < topics[j].Name
	})

	return topics, nil
}

// GetProgressiveQuestions returns questions in progressive order: L3 -> L4 -> L5 -> L6 -> L7
// Within each level: baseline before bar-raiser, then alphabetically
func (s *Scanner) GetProgressiveQuestions(index models.QuestionIndex) []*models.Question {
	// Define level order
	levelOrder := []string{"L3", "L4", "L5", "L6", "L7"}

	var result []*models.Question

	for _, level := range levelOrder {
		// Get all questions for this level
		var levelQuestions []*models.Question
		for _, question := range index {
			if question.Level == level {
				levelQuestions = append(levelQuestions, question)
			}
		}

		// Sort by category (baseline before bar-raiser) and then by file path
		sort.Slice(levelQuestions, func(i, j int) bool {
			qi, qj := levelQuestions[i], levelQuestions[j]

			// First sort by category: baseline < bar-raiser
			if qi.Category != qj.Category {
				if qi.Category == "baseline" {
					return true
				}
				if qj.Category == "baseline" {
					return false
				}
				return qi.Category < qj.Category
			}

			// Then sort alphabetically by file path
			return qi.FilePath < qj.FilePath
		})

		result = append(result, levelQuestions...)
	}

	return result
}

// formatTopicName converts a topic directory name to a display name
func formatTopicName(name string) string {
	// Remove leading numbers and dashes (e.g., "03-storage-and-retrieval" -> "storage-and-retrieval")
	parts := strings.SplitN(name, "-", 2)
	if len(parts) < 2 {
		return name
	}

	// Convert to title case and replace dashes with spaces
	displayName := strings.ReplaceAll(parts[1], "-", " ")
	words := strings.Fields(displayName)
	for i, word := range words {
		// Capitalize first letter, keep rest as is for acronyms
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}

	return strings.Join(words, " ")
}
