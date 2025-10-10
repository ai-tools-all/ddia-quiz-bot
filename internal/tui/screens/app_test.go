package screens

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/abhishek/ddia-clicker/internal/config"
	"github.com/abhishek/ddia-clicker/internal/markdown"
	"github.com/abhishek/ddia-clicker/internal/models"
	"github.com/abhishek/ddia-clicker/internal/tui/session"
)

// TestQuestionLoading tests that questions are loaded correctly from markdown files
func TestQuestionLoading(t *testing.T) {
	// Use actual content path
	contentPath := filepath.Join("..", "..", "..", "ddia-quiz-bot", "content", "chapters",
		"09-distributed-systems-gfs", "subjective")

	if _, err := os.Stat(contentPath); os.IsNotExist(err) {
		t.Skip("Skipping test - content path not found")
	}

	scanner := markdown.NewScanner(contentPath)
	index, err := scanner.ScanQuestions()
	if err != nil {
		t.Fatalf("Failed to scan questions: %v", err)
	}

	if len(index) == 0 {
		t.Fatal("No questions found in content directory")
	}

	t.Logf("Loaded %d questions", len(index))

	// Verify each question has essential fields
	emptyQuestions := 0
	for id, q := range index {
		if q.MainQuestion == "" {
			t.Errorf("Question %s has empty MainQuestion field", id)
			t.Logf("Question details: ID=%s, Level=%s, Category=%s, FilePath=%s",
				q.ID, q.Level, q.Category, q.FilePath)
			emptyQuestions++
		}
		
		if q.ID == "" {
			t.Errorf("Question has empty ID")
		}
		
		if q.Level == "" {
			t.Logf("Warning: Question %s has no level", id)
		}
	}

	if emptyQuestions > 0 {
		t.Fatalf("Found %d questions with empty MainQuestion field", emptyQuestions)
	}
}

// TestQuestionsByLevel tests that questions can be organized by level
func TestQuestionsByLevel(t *testing.T) {
	contentPath := filepath.Join("..", "..", "..", "ddia-quiz-bot", "content", "chapters",
		"09-distributed-systems-gfs", "subjective")

	if _, err := os.Stat(contentPath); os.IsNotExist(err) {
		t.Skip("Skipping test - content path not found")
	}

	scanner := markdown.NewScanner(contentPath)
	index, err := scanner.ScanQuestions()
	if err != nil {
		t.Fatalf("Failed to scan questions: %v", err)
	}

	byLevel := scanner.GetQuestionsByLevel(index)
	
	levels := []string{"L3", "L4", "L5", "L6", "L7"}
	var allQuestions []*models.Question
	
	for _, level := range levels {
		if qs, ok := byLevel[level]; ok {
			t.Logf("Level %s: %d questions", level, len(qs))
			allQuestions = append(allQuestions, qs...)
			
			// Verify each question in this level
			for _, q := range qs {
				if q.MainQuestion == "" {
					t.Errorf("Question %s (Level %s) has empty MainQuestion", q.ID, level)
				}
			}
		}
	}

	if len(allQuestions) == 0 {
		t.Fatal("No questions organized by level")
	}

	t.Logf("Total questions across all levels: %d", len(allQuestions))
}

// TestRenderQuestionWithEmptyText tests fallback when question text is empty
func TestRenderQuestionWithEmptyText(t *testing.T) {
	cfg := &config.TUIConfig{
		SessionsDir: t.TempDir(),
		ContentPath: "test",
	}

	model := NewImprovedAppModel("testuser", cfg)
	
	// Create a question with empty MainQuestion
	model.questions = []*models.Question{
		{
			ID:           "test-empty",
			MainQuestion: "",
			Level:        "L3",
		},
	}
	
	model.state = StateQuestion
	model.currentIndex = 0
	
	// Create a dummy session
	model.currentSession = &session.Session{}

	// Render and check output
	output := model.renderQuestion()
	
	if !strings.Contains(output, "ERROR: Question text is empty") {
		t.Error("Expected error message for empty question text")
	}
	
	if !strings.Contains(output, "test-empty") {
		t.Error("Expected question ID in error message")
	}
}

// TestRenderQuestionWithValidText tests normal rendering
func TestRenderQuestionWithValidText(t *testing.T) {
	cfg := &config.TUIConfig{
		SessionsDir: t.TempDir(),
		ContentPath: "test",
	}

	model := NewImprovedAppModel("testuser", cfg)
	
	questionText := "What is the purpose of replication in GFS?"
	
	model.questions = []*models.Question{
		{
			ID:           "test-valid",
			MainQuestion: questionText,
			Level:        "L3",
		},
	}
	
	model.state = StateQuestion
	model.currentIndex = 0
	model.currentSession = &session.Session{}

	output := model.renderQuestion()
	
	// Should contain the question text
	if !strings.Contains(output, questionText) {
		t.Errorf("Expected question text %q in output", questionText)
		t.Logf("Output: %s", output)
	}
	
	// Should contain progress indicator
	if !strings.Contains(output, "Question 1 of 1") {
		t.Error("Expected progress indicator")
	}
	
	// Should contain help text
	if !strings.Contains(output, "Ctrl+N") {
		t.Error("Expected keyboard shortcut help")
	}
}

// TestParseRealQuestionsE2E is an end-to-end test parsing real questions
func TestParseRealQuestionsE2E(t *testing.T) {
	contentPath := filepath.Join("..", "..", "..", "ddia-quiz-bot", "content", "chapters",
		"09-distributed-systems-gfs", "subjective")

	if _, err := os.Stat(contentPath); os.IsNotExist(err) {
		t.Skip("Skipping E2E test - content path not found")
	}

	// Find all markdown files
	var mdFiles []string
	err := filepath.Walk(contentPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".md") {
			filename := strings.ToLower(info.Name())
			if filename != "readme.md" && filename != "index.md" && filename != "guidelines.md" {
				mdFiles = append(mdFiles, path)
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to walk content directory: %v", err)
	}

	t.Logf("Found %d markdown files", len(mdFiles))

	parser := markdown.NewParser()
	failedFiles := 0
	emptyQuestions := 0

	for _, path := range mdFiles {
		question, err := parser.ParseQuestionFile(path)
		if err != nil {
			t.Logf("Failed to parse %s: %v", filepath.Base(path), err)
			failedFiles++
			continue
		}

		// Check critical fields
		if question.ID == "" {
			t.Errorf("File %s: Question has no ID", filepath.Base(path))
		}

		if question.MainQuestion == "" {
			t.Errorf("File %s: Question has empty MainQuestion (ID: %s)", filepath.Base(path), question.ID)
			emptyQuestions++
		} else {
			t.Logf("✓ %s: MainQuestion length = %d chars", filepath.Base(path), len(question.MainQuestion))
		}
	}

	if failedFiles > 0 {
		t.Errorf("Failed to parse %d out of %d files", failedFiles, len(mdFiles))
	}

	if emptyQuestions > 0 {
		t.Fatalf("Found %d questions with empty MainQuestion", emptyQuestions)
	}
}

// TestQuestionTextNotEmpty verifies all loaded questions have non-empty text
func TestQuestionTextNotEmpty(t *testing.T) {
	contentPath := filepath.Join("..", "..", "..", "ddia-quiz-bot", "content", "chapters",
		"09-distributed-systems-gfs", "subjective")

	if _, err := os.Stat(contentPath); os.IsNotExist(err) {
		t.Skip("Skipping test - content path not found")
	}

	scanner := markdown.NewScanner(contentPath)
	index, err := scanner.ScanQuestions()
	if err != nil {
		t.Fatalf("Failed to scan questions: %v", err)
	}

	var problematicQuestions []string

	for id, q := range index {
		if q.MainQuestion == "" {
			problematicQuestions = append(problematicQuestions, id)
			t.Errorf("Question %s has empty MainQuestion", id)
			
			// Read the file to debug
			if q.FilePath != "" {
				content, err := os.ReadFile(q.FilePath)
				if err == nil {
					t.Logf("File content preview for %s:\n%s\n", id, string(content[:min(500, len(content))]))
				}
			}
		} else {
			// Log successful parsing for verification
			preview := q.MainQuestion
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			t.Logf("✓ %s: %s", id, preview)
		}
	}

	if len(problematicQuestions) > 0 {
		t.Fatalf("Found %d questions with empty MainQuestion: %v",
			len(problematicQuestions), problematicQuestions)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
