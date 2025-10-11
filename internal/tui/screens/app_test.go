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

	// Verify multiple topics have questions
	hasQuestions := 0
	for _, topic := range topics {
		if topic.TotalCount > 0 {
			hasQuestions++
			t.Logf("Topic %s has %d questions", topic.Name, topic.TotalCount)
		}
	}

	assert.GreaterOrEqual(t, hasQuestions, 1, "At least one chapter should have subjective questions")
	t.Logf("Found %d topics with questions", hasQuestions)
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
		assert.Equal(t, StateWelcome, updatedApp.state,
			"Should remain on welcome screen until user input")

		model, _ = updatedApp.Update(tea.KeyMsg{Type: tea.KeyEnter})
		afterInput := model.(ImprovedAppModel)
		assert.Equal(t, StateModeSelect, afterInput.state,
			"Enter should advance to mode selection after topics are loaded")
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
		assert.Contains(t, view, "Available topics: 1",
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
		assert.Contains(t, view, "Loaded questions: 5",
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

	require.Equal(t, StateModeSelect, updated.state, "Enter should advance to mode selection when topics exist")

	model, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	modeSelected := model.(ImprovedAppModel)
	require.Equal(t, StateTopicSelect, modeSelected.state, "Selecting a mode should enter topic selection")
	require.Equal(t, "subjective", modeSelected.selectedMode)

	model, cmd := modeSelected.Update(tea.KeyMsg{Type: tea.KeyEnter})
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
	manager.UpdateResponse(older, "q1", "subjective", "old answer", 10)
	require.NoError(t, manager.SaveSession(older))

	newer, err := manager.CreateSessionWithTopic("alice", "subjective", "09-topic", "Topic", app.questions)
	require.NoError(t, err)
	newer.Session.SessionID = "newer-session"
	manager.UpdateResponse(newer, "q1", "subjective", "latest answer", 5)
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

// TestInitializeQuestionComponentMCQDetection tests MCQ vs subjective question type detection
func TestInitializeQuestionComponentMCQDetection(t *testing.T) {
	tests := []struct {
		name          string
		question      *models.Question
		wantQType     string
		wantMCQNotNil bool
	}{
		{
			name: "MCQ question with type set",
			question: &models.Question{
				ID:           "mcq-001",
				Type:         "mcq",
				MainQuestion: "Test MCQ?",
				Options:      []string{"A) Option 1", "B) Option 2", "C) Option 3"},
				Answer:       "B",
				Explanation:  "Test explanation",
			},
			wantQType:     "mcq",
			wantMCQNotNil: true,
		},
		{
			name: "Subjective question with type set",
			question: &models.Question{
				ID:           "subj-001",
				Type:         "subjective",
				MainQuestion: "Explain the concept...",
			},
			wantQType:     "subjective",
			wantMCQNotNil: false,
		},
		{
			name: "Question with no type defaults to subjective",
			question: &models.Question{
				ID:           "no-type-001",
				Type:         "",
				MainQuestion: "Some question...",
			},
			wantQType:     "subjective",
			wantMCQNotNil: false,
		},
		{
			name: "Question with empty type string defaults to subjective",
			question: &models.Question{
				ID:           "empty-type-001",
				MainQuestion: "Another question...",
			},
			wantQType:     "subjective",
			wantMCQNotNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.TUIConfig{
				SessionsDir:      t.TempDir(),
				AutoSaveInterval: time.Second,
			}

			app := NewImprovedAppModel("testuser", cfg)
			app.questions = []*models.Question{tt.question}
			app.currentIndex = 0

			// Initialize component
			app.initializeQuestionComponent()

			// Verify question type
			assert.Equal(t, tt.wantQType, app.currentQType,
				"Question type should be %s", tt.wantQType)

			// Verify MCQ component state
			if tt.wantMCQNotNil {
				assert.NotNil(t, app.mcqComponent, "MCQ component should be initialized for MCQ questions")
				assert.False(t, app.showExplanation, "Explanation should be hidden initially")
			} else {
				assert.Nil(t, app.mcqComponent, "MCQ component should be nil for subjective questions")
			}
		})
	}
}

// TestSwitchingBetweenQuestionTypes tests state management when switching between MCQ and subjective
func TestSwitchingBetweenQuestionTypes(t *testing.T) {
	cfg := &config.TUIConfig{
		SessionsDir:      t.TempDir(),
		AutoSaveInterval: time.Second,
	}

	app := NewImprovedAppModel("testuser", cfg)

	// Create a mix of questions
	app.questions = []*models.Question{
		{
			ID:           "subj-001",
			Type:         "subjective",
			MainQuestion: "Subjective question 1",
		},
		{
			ID:           "mcq-001",
			Type:         "mcq",
			MainQuestion: "MCQ question 1?",
			Options:      []string{"A) Option 1", "B) Option 2"},
			Answer:       "A",
			Explanation:  "Explanation for MCQ",
		},
		{
			ID:           "subj-002",
			Type:         "subjective",
			MainQuestion: "Subjective question 2",
		},
	}

	// Start with first subjective question
	app.currentIndex = 0
	app.initializeQuestionComponent()

	assert.Equal(t, "subjective", app.currentQType)
	assert.Nil(t, app.mcqComponent)

	// Switch to MCQ
	app.currentIndex = 1
	app.initializeQuestionComponent()

	assert.Equal(t, "mcq", app.currentQType)
	assert.NotNil(t, app.mcqComponent)

	// Switch back to subjective
	app.currentIndex = 2
	app.initializeQuestionComponent()

	assert.Equal(t, "subjective", app.currentQType)
	assert.Nil(t, app.mcqComponent, "MCQ component should be cleared when switching to subjective")

	// Switch back to MCQ again
	app.currentIndex = 1
	app.initializeQuestionComponent()

	assert.Equal(t, "mcq", app.currentQType)
	assert.NotNil(t, app.mcqComponent, "MCQ component should be re-initialized")
	assert.Equal(t, 0, app.mcqComponent.SelectedIdx, "Should start fresh with first option selected")
	assert.False(t, app.mcqComponent.Submitted, "Should not be submitted on fresh initialization")
}

// TestTextareaResetAndAnswerLoading tests textarea behavior for subjective questions
func TestTextareaResetAndAnswerLoading(t *testing.T) {
	cfg := &config.TUIConfig{
		SessionsDir:      t.TempDir(),
		AutoSaveInterval: time.Second,
	}

	app := NewImprovedAppModel("testuser", cfg)

	questions := []*models.Question{
		{
			ID:           "subj-001",
			Type:         "subjective",
			MainQuestion: "Question 1",
		},
		{
			ID:           "subj-002",
			Type:         "subjective",
			MainQuestion: "Question 2",
		},
	}

	app.questions = questions

	// Create a session with an existing answer
	session, err := app.sessionManager.CreateSession("testuser", "subjective", questions)
	require.NoError(t, err)

	// Add existing answer for first question
	app.sessionManager.UpdateResponse(session, "subj-001", "subjective", "existing answer for q1", 30)
	require.NoError(t, app.sessionManager.SaveSession(session))

	app.currentSession = session

	// Initialize first question with existing answer
	app.currentIndex = 0
	app.initializeQuestionComponent()

	assert.Equal(t, "existing answer for q1", app.textarea.Value(),
		"Should load existing answer from session")

	// Switch to second question (no existing answer)
	app.currentIndex = 1
	app.initializeQuestionComponent()

	assert.Empty(t, app.textarea.Value(),
		"Should reset textarea when no existing answer")
}

// TestMCQComponentInitialization tests MCQ component is properly configured
func TestMCQComponentInitialization(t *testing.T) {
	cfg := &config.TUIConfig{
		SessionsDir:      t.TempDir(),
		AutoSaveInterval: time.Second,
	}

	app := NewImprovedAppModel("testuser", cfg)

	question := &models.Question{
		ID:           "mcq-001",
		Type:         "mcq",
		MainQuestion: "What is 2+2?",
		Options:      []string{"A) 3", "B) 4", "C) 5", "D) 6"},
		Answer:       "B",
		Explanation:  "Two plus two equals four",
		Hook:         "Basic math",
	}

	app.questions = []*models.Question{question}
	app.currentIndex = 0

	app.initializeQuestionComponent()

	require.NotNil(t, app.mcqComponent)

	// Verify MCQ component properties
	assert.Equal(t, 4, len(app.mcqComponent.Options), "Should have 4 options")
	assert.Equal(t, 0, app.mcqComponent.SelectedIdx, "Should start with first option selected")
	assert.False(t, app.mcqComponent.Submitted, "Should not be submitted initially")
	assert.Equal(t, 1, app.mcqComponent.CorrectIdx, "Correct answer B should be index 1")
	assert.Equal(t, "Two plus two equals four", app.mcqComponent.Explanation,
		"Explanation should be set")
}

// TestInitializeOutOfBoundsHandling tests boundary conditions
func TestInitializeOutOfBoundsHandling(t *testing.T) {
	cfg := &config.TUIConfig{
		SessionsDir:      t.TempDir(),
		AutoSaveInterval: time.Second,
	}

	app := NewImprovedAppModel("testuser", cfg)

	app.questions = []*models.Question{
		{ID: "q1", Type: "subjective", MainQuestion: "Question 1"},
	}

	// Test with index out of bounds
	app.currentIndex = 999
	app.initializeQuestionComponent()

	// Should not panic and should not set currentQType
	assert.Empty(t, app.currentQType, "Should not set type when out of bounds")

	// Test with empty questions slice
	app.questions = []*models.Question{}
	app.currentIndex = 0
	app.initializeQuestionComponent()

	assert.Empty(t, app.currentQType, "Should not set type when questions is empty")
}

// TestMCQStateTransitions tests the state machine for MCQ question flow
func TestMCQStateTransitions(t *testing.T) {
	cfg := &config.TUIConfig{
		SessionsDir:      t.TempDir(),
		AutoSaveInterval: time.Second,
	}

	app := NewImprovedAppModel("testuser", cfg)

	mcqQuestion := &models.Question{
		ID:           "mcq-001",
		Type:         "mcq",
		MainQuestion: "Test MCQ?",
		Options:      []string{"A) Option 1", "B) Option 2", "C) Option 3"},
		Answer:       "B",
		Explanation:  "B is correct",
	}

	app.questions = []*models.Question{mcqQuestion}
	app.currentIndex = 0
	app.state = StateQuestion

	// Create session
	session, err := app.sessionManager.CreateSession("testuser", "mcq", app.questions)
	require.NoError(t, err)
	app.currentSession = session
	app.questionStartTime = time.Now()

	// Initialize MCQ component
	app.initializeQuestionComponent()
	require.NotNil(t, app.mcqComponent)

	// Test navigation (move down)
	model, _ := app.Update(tea.KeyMsg{Type: tea.KeyDown})
	updated := model.(ImprovedAppModel)
	assert.Equal(t, 1, updated.mcqComponent.SelectedIdx, "Should move to second option")

	// Test navigation (move up)
	model, _ = updated.Update(tea.KeyMsg{Type: tea.KeyUp})
	updated = model.(ImprovedAppModel)
	assert.Equal(t, 0, updated.mcqComponent.SelectedIdx, "Should move back to first option")

	// Move to correct answer (B = index 1)
	model, _ = updated.Update(tea.KeyMsg{Type: tea.KeyDown})
	updated = model.(ImprovedAppModel)

	// Submit answer
	model, cmd := updated.Update(tea.KeyMsg{Type: tea.KeyEnter})
	submitted := model.(ImprovedAppModel)

	assert.True(t, submitted.mcqComponent.Submitted, "MCQ should be marked as submitted")
	assert.NotNil(t, cmd, "Should trigger save command")

	// Verify save command was issued
	msg := cmd()
	_, ok := msg.(answerSavedMsg)
	assert.True(t, ok, "Should return answerSavedMsg")
}

// TestResumeMCQSession tests resuming a session with MCQ answers
func TestResumeMCQSession(t *testing.T) {
	cfg := &config.TUIConfig{
		SessionsDir:      t.TempDir(),
		AutoSaveInterval: time.Second,
	}

	app := NewImprovedAppModel("testuser", cfg)

	mcqQuestion := &models.Question{
		ID:           "mcq-001",
		Type:         "mcq",
		MainQuestion: "Test?",
		Options:      []string{"A) Opt 1", "B) Opt 2"},
		Answer:       "A",
	}

	app.questions = []*models.Question{mcqQuestion}

	// Create session with MCQ answer
	session, err := app.sessionManager.CreateSession("testuser", "mcq", app.questions)
	require.NoError(t, err)

	isCorrect := true
	app.sessionManager.UpdateResponse(session, "mcq-001", "mcq", "A", 10)
	// Update the response to include MCQ-specific fields
	session.Responses[0].QuestionType = "mcq"
	session.Responses[0].IsCorrect = &isCorrect
	session.Responses[0].SelectedOption = "A"
	require.NoError(t, app.sessionManager.SaveSession(session))

	// Load session
	loaded, err := app.sessionManager.LoadSession("testuser", "mcq", session.Session.SessionID)
	require.NoError(t, err)

	assert.Len(t, loaded.Responses, 1)
	assert.Equal(t, "mcq", loaded.Responses[0].QuestionType)
	assert.Equal(t, "A", loaded.Responses[0].Answer)
	assert.NotNil(t, loaded.Responses[0].IsCorrect)
	assert.True(t, *loaded.Responses[0].IsCorrect)
}
