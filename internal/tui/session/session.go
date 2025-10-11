package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/abhishek/ddia-clicker/internal/models"
)

// Status represents the status of a quiz session
type Status string

const (
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusAborted    Status = "aborted"
)

// Session represents a quiz session with all its metadata and responses
type Session struct {
	Session   SessionMetadata   `json:"session"`
	Questions []QuestionSummary `json:"questions"`
	Responses []Response        `json:"responses"`
}

// SessionMetadata holds session-level information
type SessionMetadata struct {
	SessionID     string    `json:"session_id"`
	User          string    `json:"user"`
	Mode          string    `json:"mode"`
	Topic         string    `json:"topic"`         // Topic identifier (e.g., "09-distributed-systems-gfs")
	TopicDisplay  string    `json:"topic_display"` // Human-readable topic name
	Status        Status    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	QuestionCount int       `json:"question_count"`
	Answered      int       `json:"answered"`
}

// QuestionSummary holds minimal question info for the session
type QuestionSummary struct {
	ID       string            `json:"id"`
	Title    string            `json:"title"`
	Level    string            `json:"level"`
	Chapter  string            `json:"chapter"`
	Metadata map[string]string `json:"metadata"`
}

// Response represents a user's response to a question
type Response struct {
	QuestionID       string    `json:"question_id"`
	QuestionType     string    `json:"question_type"` // "subjective" or "mcq"
	Answer           string    `json:"answer"`
	IsCorrect        *bool     `json:"is_correct,omitempty"`      // For MCQ, nil for subjective
	SelectedOption   string    `json:"selected_option,omitempty"` // MCQ: A, B, C, D
	UpdatedAt        time.Time `json:"updated_at"`
	TimeSpentSeconds int       `json:"time_spent_seconds"`
}

// Manager handles session persistence and retrieval
type Manager struct {
	baseDir string
}

// NewManager creates a new session manager
func NewManager(baseDir string) *Manager {
	return &Manager{
		baseDir: baseDir,
	}
}

// CreateSession creates a new session for a user
func (m *Manager) CreateSession(user string, mode string, questions []*models.Question) (*Session, error) {
	return m.CreateSessionWithTopic(user, mode, "", "", questions)
}

// CreateSessionWithTopic creates a new session for a user with topic information
func (m *Manager) CreateSessionWithTopic(user string, mode string, topic string, topicDisplay string, questions []*models.Question) (*Session, error) {
	timestamp := time.Now().Format("20060102-150405")
	sessionID := fmt.Sprintf("%s-%s-%s", timestamp, user, mode)
	if topic != "" {
		sessionID = fmt.Sprintf("%s-%s-%s-%s", timestamp, user, mode, topic)
	}

	session := &Session{
		Session: SessionMetadata{
			SessionID:     sessionID,
			User:          user,
			Mode:          mode,
			Topic:         topic,
			TopicDisplay:  topicDisplay,
			Status:        StatusInProgress,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			QuestionCount: len(questions),
			Answered:      0,
		},
		Questions: make([]QuestionSummary, 0, len(questions)),
		Responses: make([]Response, 0),
	}

	// Convert questions to summaries
	for _, q := range questions {
		session.Questions = append(session.Questions, QuestionSummary{
			ID:      q.ID,
			Title:   q.Title,
			Level:   q.Level,
			Chapter: q.Category,
			Metadata: map[string]string{
				"category": q.Category,
			},
		})
	}

	return session, nil
}

// SaveSession saves a session to disk
func (m *Manager) SaveSession(session *Session) error {
	// Update timestamp
	session.Session.UpdatedAt = time.Now()

	// Create directory structure
	sessionDir := filepath.Join(m.baseDir, session.Session.User, session.Session.Mode)
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return fmt.Errorf("failed to create session directory: %w", err)
	}

	// Write to temp file first, then rename for atomicity
	filename := fmt.Sprintf("%s.json", session.Session.SessionID)
	filepath := filepath.Join(sessionDir, filename)
	tempPath := filepath + ".tmp"

	// Marshal to JSON
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Write to temp file
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempPath, filepath); err != nil {
		os.Remove(tempPath) // Clean up temp file
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

// LoadSession loads a session from disk by session ID
func (m *Manager) LoadSession(user string, mode string, sessionID string) (*Session, error) {
	filename := fmt.Sprintf("%s.json", sessionID)
	filepath := filepath.Join(m.baseDir, user, mode, filename)

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

// ListIncompleteSessions lists all incomplete sessions for a user and mode
func (m *Manager) ListIncompleteSessions(user string, mode string) ([]*Session, error) {
	return m.ListIncompleteSessionsForTopic(user, mode, "")
}

// ListIncompleteSessionsForTopic lists incomplete sessions filtered by topic
func (m *Manager) ListIncompleteSessionsForTopic(user string, mode string, topic string) ([]*Session, error) {
	sessionDir := filepath.Join(m.baseDir, user, mode)

	// Check if directory exists
	if _, err := os.Stat(sessionDir); os.IsNotExist(err) {
		return []*Session{}, nil
	}

	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read session directory: %w", err)
	}

	var sessions []*Session
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		sessionID := entry.Name()[:len(entry.Name())-5] // Remove .json extension
		session, err := m.LoadSession(user, mode, sessionID)
		if err != nil {
			// Log but continue
			fmt.Fprintf(os.Stderr, "Warning: Failed to load session %s: %v\n", sessionID, err)
			continue
		}

		// Filter by status and optionally by topic
		if session.Session.Status == StatusInProgress {
			if topic == "" || session.Session.Topic == topic {
				sessions = append(sessions, session)
			}
		}
	}

	// Sort by creation time (newest first)
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Session.CreatedAt.After(sessions[j].Session.CreatedAt)
	})

	return sessions, nil
}

// UpdateResponse updates or adds a response to the session
func (m *Manager) UpdateResponse(session *Session, questionID string, answer string, timeSpent int) {
	// Check if response already exists
	found := false
	for i := range session.Responses {
		if session.Responses[i].QuestionID == questionID {
			session.Responses[i].Answer = answer
			session.Responses[i].UpdatedAt = time.Now()
			session.Responses[i].TimeSpentSeconds += timeSpent
			found = true
			break
		}
	}

	// If not found, add new response
	if !found {
		session.Responses = append(session.Responses, Response{
			QuestionID:       questionID,
			Answer:           answer,
			UpdatedAt:        time.Now(),
			TimeSpentSeconds: timeSpent,
		})
	}

	// Update answered count
	session.Session.Answered = len(session.Responses)
}

// CompleteSession marks a session as completed
func (m *Manager) CompleteSession(session *Session) error {
	session.Session.Status = StatusCompleted
	session.Session.UpdatedAt = time.Now()
	return m.SaveSession(session)
}

// GetResponse retrieves a response for a specific question
func (m *Manager) GetResponse(session *Session, questionID string) *Response {
	for i := range session.Responses {
		if session.Responses[i].QuestionID == questionID {
			return &session.Responses[i]
		}
	}
	return nil
}
