package screens

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/abhishek/ddia-clicker/internal/config"
	"github.com/abhishek/ddia-clicker/internal/markdown"
	"github.com/abhishek/ddia-clicker/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuestionLoadingFromSubjectiveFolder(t *testing.T) {
	// Test loading questions from the subjective folder as configured
	tests := []struct {
		name           string
		configPath     string
		expectQuestions bool
		minQuestions   int
		description    string
	}{
		{
			name:           "Load questions from GFS subjective folder",
			configPath:     "ddia-quiz-bot/content/chapters/09-distributed-systems-gfs/subjective",
			expectQuestions: true,
			minQuestions:   1,
			description:    "Should successfully load questions from GFS subjective folder",
		},
		{
			name:           "Fail on non-existent chapter folder",
			configPath:     "ddia-quiz-bot/content/chapters/03-storage-and-retrieval/subjective",
			expectQuestions: false,
			minQuestions:   0,
			description:    "Should handle missing subjective folder gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get absolute path
			absPath := filepath.Join("/home/abhishek/Downloads/experiments/ai-tools/ddia-clicker", tt.configPath)
			
			// Create scanner
			scanner := markdown.NewScanner(absPath)
			
			// Try to scan questions
			index, err := scanner.ScanQuestions()
			
			if tt.expectQuestions {
				require.NoError(t, err, "Should not error when loading questions from %s", tt.configPath)
				require.NotNil(t, index, "Question index should not be nil")
				assert.GreaterOrEqual(t, len(index), tt.minQuestions, 
					"Should have at least %d questions, got %d", tt.minQuestions, len(index))
				
				// Verify questions have required fields
				for id, question := range index {
					assert.NotEmpty(t, question.ID, "Question ID should not be empty")
					assert.NotEmpty(t, question.MainQuestion, "Question text should not be empty for ID: %s", id)
					assert.NotEmpty(t, question.Level, "Question level should not be empty for ID: %s", id)
					assert.Contains(t, []string{"L3", "L4", "L5", "L6", "L7"}, question.Level, 
						"Question level should be valid for ID: %s", id)
				}
				
				// Test progressive question ordering
				questions := scanner.GetProgressiveQuestions(index)
				assert.Equal(t, len(index), len(questions), "Progressive questions should include all questions")
				
				// Verify level ordering
				var lastLevel string
				levelOrder := map[string]int{"L3": 1, "L4": 2, "L5": 3, "L6": 4, "L7": 5}
				for i, q := range questions {
					if lastLevel != "" {
						assert.GreaterOrEqual(t, levelOrder[q.Level], levelOrder[lastLevel],
							"Question %d: Level %s should come after or equal to %s", i, q.Level, lastLevel)
					}
					lastLevel = q.Level
				}
			} else {
				// Should error or return empty when folder doesn't exist
				if err == nil {
					assert.Empty(t, index, "Should return empty index for non-existent folder")
				}
			}
		})
	}
}

func TestTopicDiscoveryWithChaptersRoot(t *testing.T) {
	// Test topic discovery functionality
	chaptersRoot := "/home/abhishek/Downloads/experiments/ai-tools/ddia-clicker/ddia-quiz-bot/content/chapters"
	
	scanner := markdown.NewScanner("")
	topics, err := scanner.DiscoverTopics(chaptersRoot)
	
	require.NoError(t, err, "Topic discovery should not error")
	require.NotEmpty(t, topics, "Should discover at least one topic")
	
	// Find GFS topic
	var gfsTopic *markdown.TopicInfo
	for _, topic := range topics {
		if topic.Name == "09-distributed-systems-gfs" {
			gfsTopic = &topic
			break
		}
	}
	
	require.NotNil(t, gfsTopic, "Should find GFS topic")
	assert.Equal(t, "Distributed Systems Gfs", gfsTopic.DisplayName)
	assert.Greater(t, gfsTopic.TotalCount, 0, "GFS topic should have questions")
	assert.NotEmpty(t, gfsTopic.LevelCounts, "Should have level counts")
	
	// Verify only GFS has questions (based on current state)
	hasQuestions := 0
	for _, topic := range topics {
		if topic.TotalCount > 0 {
			hasQuestions++
			t.Logf("Topic %s has %d questions", topic.Name, topic.TotalCount)
		}
	}
	
	assert.Equal(t, 1, hasQuestions, "Currently only one chapter (GFS) should have subjective questions")
}

func TestAppModelQuestionLoading(t *testing.T) {
	// Test the actual app model's question loading flow
	cfg := &config.TUIConfig{
		ChaptersRootPath: "ddia-quiz-bot/content/chapters",
		SessionsDir:      "test_sessions",
		AutoSaveInterval: 30 * time.Second,
	}
	
	app := NewImprovedAppModel("testuser", cfg)
	
	// Initialize the app
	cmd := app.Init()
	require.NotNil(t, cmd, "Init should return a command")
	
	// Simulate topic discovery message
	scanner := markdown.NewScanner("")
	topics, err := scanner.DiscoverTopics(
		"/home/abhishek/Downloads/experiments/ai-tools/ddia-clicker/" + cfg.ChaptersRootPath)
	
	if err == nil && len(topics) > 0 {
		// Send topics discovered message
		msg := topicsDiscoveredMsg{topics: topics, err: nil}
		model, _ := app.Update(msg)
		updatedApp := model.(ImprovedAppModel)
		
		assert.Equal(t, len(topics), len(updatedApp.availableTopics), 
			"Should have discovered topics")
		assert.NotEqual(t, StateWelcome, updatedApp.state, 
			"Should not be stuck in welcome state after topic discovery")
	}
}

func TestWelcomeScreenLoadingState(t *testing.T) {
	// Test that welcome screen properly shows loading state
	
	t.Run("Topic mode", func(t *testing.T) {
		cfg := &config.TUIConfig{
			ChaptersRootPath: "ddia-quiz-bot/content/chapters",
			SessionsDir:      "test_sessions",
			AutoSaveInterval: 30 * time.Second,
		}
		
		app := NewImprovedAppModel("testuser", cfg)
		
		// Before topics are discovered
		view := app.renderWelcome()
		assert.Contains(t, view, "Discovering topics...", 
			"Should show discovering topics message in topic mode")
		
		// After topics are discovered
		app.availableTopics = []markdown.TopicInfo{
			{Name: "test-topic", DisplayName: "Test Topic", TotalCount: 10},
		}
		view = app.renderWelcome()
		assert.NotContains(t, view, "Discovering topics...", 
			"Should not show discovering message after topics are loaded")
		assert.Contains(t, view, "Found: 1 topics", 
			"Should show topic count when loaded")
	})
	
	t.Run("Single topic mode", func(t *testing.T) {
		cfg := &config.TUIConfig{
			ContentPath:      "ddia-quiz-bot/content/chapters/09-distributed-systems-gfs/subjective",
			ChaptersRootPath: "", // Empty for single topic mode
			SessionsDir:      "test_sessions",
			AutoSaveInterval: 30 * time.Second,
		}
		
		app := NewImprovedAppModel("testuser", cfg)
		
		// Before questions are loaded
		view := app.renderWelcome()
		assert.Contains(t, view, "Loading questions...", 
			"Should show loading questions message in single topic mode")
		
		// After questions are loaded
		app.questions = make([]*models.Question, 5)
		view = app.renderWelcome()
		assert.NotContains(t, view, "Loading questions...", 
			"Should not show loading message after questions are loaded")
		assert.Contains(t, view, "Loaded: 5 questions", 
			"Should show question count when loaded")
	})
}

func TestSingleTopicModeFallback(t *testing.T) {
	// Test fallback to single topic mode when chapters_root_path is empty
	cfg := &config.TUIConfig{
		ContentPath:      "ddia-quiz-bot/content/chapters/09-distributed-systems-gfs/subjective",
		ChaptersRootPath: "", // Empty to trigger single topic mode
		SessionsDir:      "test_sessions",
		AutoSaveInterval: 30 * time.Second,
	}
	
	app := NewImprovedAppModel("testuser", cfg)
	
	// Init should trigger loadQuestionsCmd instead of discoverTopicsCmd
	cmd := app.Init()
	require.NotNil(t, cmd, "Init should return a command")
	
	// In single topic mode, it should try to load questions directly
	scanner := markdown.NewScanner(
		"/home/abhishek/Downloads/experiments/ai-tools/ddia-clicker/" + cfg.ContentPath)
	index, err := scanner.ScanQuestions()
	
	if err == nil {
		assert.NotEmpty(t, index, "Should load questions in single topic mode")
	}
}
