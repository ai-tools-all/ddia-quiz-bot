package screens

import (
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/abhishek/ddia-clicker/internal/config"
	"github.com/abhishek/ddia-clicker/internal/markdown"
	"github.com/abhishek/ddia-clicker/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuestionLoadingFromSubjectiveFolder(t *testing.T) {
	// Test loading questions from the subjective folder as configured
	tests := []struct {
		name            string
		configPath      string
		expectQuestions bool
		minQuestions    int
		description     string
	}{
		{
			name:            "Load questions from GFS subjective folder",
			configPath:      "ddia-quiz-bot/content/chapters/09-distributed-systems-gfs/subjective",
			expectQuestions: true,
			minQuestions:    1,
			description:     "Should successfully load questions from GFS subjective folder",
		},
		{
			name:            "Fail on non-existent chapter folder",
			configPath:      "ddia-quiz-bot/content/chapters/03-storage-and-retrieval/subjective",
			expectQuestions: false,
			minQuestions:    0,
			description:     "Should handle missing subjective folder gracefully",
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

func TestTopicSelectionTransition(t *testing.T) {
	cfg := &config.TUIConfig{
		ChaptersRootPath: "unused",
		SessionsDir:      t.TempDir(),
		AutoSaveInterval: 100 * time.Millisecond,
	}

	app := NewImprovedAppModel("tester", cfg)
	app.availableTopics = []markdown.TopicInfo{
		{Name: "01-topic-alpha", DisplayName: "Topic Alpha", TotalCount: 5, Path: "/tmp/topic-alpha"},
		{Name: "02-topic-beta", DisplayName: "Topic Beta", TotalCount: 3, Path: "/tmp/topic-beta"},
	}

	model, _ := app.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updated := model.(ImprovedAppModel)

	require.Equal(t, StateTopicSelect, updated.state, "Enter should advance to topic selection when topics exist")

	model, cmd := updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	selected := model.(ImprovedAppModel)

	require.NotNil(t, cmd, "Selecting a topic should trigger follow-up commands")
	require.NotNil(t, selected.selectedTopic)
	assert.Equal(t, "01-topic-alpha", selected.selectedTopic.Name)
	assert.Equal(t, StateTopicSelect, selected.state, "State stays in topic select until questions arrive")
}

func TestGetProgressiveQuestionsOrderingSynthetic(t *testing.T) {
	index := models.QuestionIndex{
		"q1": {ID: "q1", Level: "L4", Category: "bar-raiser", FilePath: "b.md"},
		"q2": {ID: "q2", Level: "L3", Category: "baseline", FilePath: "a.md"},
		"q3": {ID: "q3", Level: "L3", Category: "bar-raiser", FilePath: "c.md"},
		"q4": {ID: "q4", Level: "L5", Category: "baseline", FilePath: "a.md"},
		"q5": {ID: "q5", Level: "L4", Category: "baseline", FilePath: "a.md"},
	}

	scanner := markdown.NewScanner("")
	ordered := scanner.GetProgressiveQuestions(index)

	require.Len(t, ordered, 5)
	ids := make([]string, 0, len(ordered))
	for _, q := range ordered {
		ids = append(ids, q.ID)
	}

	assert.Equal(t, []string{"q2", "q3", "q5", "q1", "q4"}, ids,
		"Questions should follow L3→L4→L5 with baseline before bar-raiser and alphabetical fallback")
}

func TestAutoSaveCyclePersistsResponse(t *testing.T) {
	tempDir := t.TempDir()
	cfg := &config.TUIConfig{
		SessionsDir:      tempDir,
		AutoSaveInterval: 10 * time.Millisecond,
	}

	app := NewImprovedAppModel("autosaver", cfg)
	question := &models.Question{ID: "q1", MainQuestion: "Q1", Level: "L3", Category: "baseline"}
	app.questions = []*models.Question{question}
	app.currentIndex = 0
	app.state = StateQuestion
	app.questionStartTime = time.Now().Add(-time.Second)
	app.lastSaveTime = time.Now().Add(-cfg.AutoSaveInterval * 2)
	app.textarea.SetValue(" autosave answer ")

	sessionObj, err := app.sessionManager.CreateSession("autosaver", "subjective", app.questions)
	require.NoError(t, err)
	require.NoError(t, app.sessionManager.SaveSession(sessionObj))
	app.currentSession = sessionObj

	model, cmd := app.Update(autoSaveTickMsg{})
	updated := model.(ImprovedAppModel)
	require.NotNil(t, cmd, "Auto-save tick should issue a save command when interval elapsed")

	msg := cmd()
	answerSaved, ok := msg.(answerSavedMsg)
	require.True(t, ok, "Expected answerSavedMsg from auto-save command")

	model, _ = updated.Update(answerSaved)
	final := model.(ImprovedAppModel)
	assert.WithinDuration(t, time.Now(), final.lastSaveTime, 50*time.Millisecond)

	saved, err := final.sessionManager.LoadSession("autosaver", "subjective", sessionObj.Session.SessionID)
	require.NoError(t, err)
	require.Len(t, saved.Responses, 1)
	assert.Equal(t, "autosave answer", saved.Responses[0].Answer)
	assert.Equal(t, 1, saved.Session.Answered)
}

func TestResumeMostRecentSessionLoadsState(t *testing.T) {
	tempDir := t.TempDir()
	cfg := &config.TUIConfig{
		SessionsDir:      tempDir,
		AutoSaveInterval: time.Second,
	}

	app := NewImprovedAppModel("alice", cfg)
	app.state = StateSessionSelect
	app.selectedTopic = &markdown.TopicInfo{Name: "09-topic", DisplayName: "Topic", TotalCount: 1, Path: "/tmp/topic"}

	question := &models.Question{ID: "q1", MainQuestion: "Q1", Level: "L3", Category: "baseline"}
	app.questions = []*models.Question{question}

	manager := app.sessionManager

	older, err := manager.CreateSessionWithTopic("alice", "subjective", "09-topic", "Topic", app.questions)
	require.NoError(t, err)
	older.Session.SessionID = "older-session"
	older.Session.CreatedAt = time.Now().Add(-2 * time.Hour)
	manager.UpdateResponse(older, "q1", "old answer", 10)
	require.NoError(t, manager.SaveSession(older))

	newer, err := manager.CreateSessionWithTopic("alice", "subjective", "09-topic", "Topic", app.questions)
	require.NoError(t, err)
	newer.Session.SessionID = "newer-session"
	manager.UpdateResponse(newer, "q1", "latest answer", 5)
	require.NoError(t, manager.SaveSession(newer))

	cmd := app.checkTopicSessionsCmd()
	msg := cmd()
	sessionsMsg, ok := msg.(existingSessionsMsg)
	require.True(t, ok)

	model, _ := app.Update(sessionsMsg)
	updated := model.(ImprovedAppModel)
	require.Len(t, updated.existingSessions, 2)
	assert.Equal(t, newer.Session.SessionID, updated.existingSessions[0].Session.SessionID,
		"Newest session should be first")

	model, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	resumed := model.(ImprovedAppModel)

	assert.Equal(t, StateQuestion, resumed.state)
	require.NotNil(t, resumed.currentSession)
	assert.Equal(t, newer.Session.SessionID, resumed.currentSession.Session.SessionID)
	assert.Equal(t, "latest answer", resumed.textarea.Value(), "Existing answer should be loaded into textarea")
}

func TestSaveBeforeQuitPersistsAnswer(t *testing.T) {
	tempDir := t.TempDir()
	cfg := &config.TUIConfig{
		SessionsDir:      tempDir,
		AutoSaveInterval: time.Second,
	}

	app := NewImprovedAppModel("quituser", cfg)
	question := &models.Question{ID: "q1", MainQuestion: "Q1", Level: "L3", Category: "baseline"}
	app.questions = []*models.Question{question}
	app.state = StateQuestion
	app.currentIndex = 0
	app.questionStartTime = time.Now().Add(-2 * time.Second)
	app.textarea.SetValue("pending answer")

	sessionObj, err := app.sessionManager.CreateSession("quituser", "subjective", app.questions)
	require.NoError(t, err)
	app.currentSession = sessionObj

	app.saveBeforeQuit()

	saved, err := app.sessionManager.LoadSession("quituser", "subjective", sessionObj.Session.SessionID)
	require.NoError(t, err)
	require.Len(t, saved.Responses, 1)
	assert.Equal(t, "pending answer", saved.Responses[0].Answer)
	assert.Equal(t, 1, saved.Session.Answered)
	assert.Greater(t, saved.Responses[0].TimeSpentSeconds, 0)
}
