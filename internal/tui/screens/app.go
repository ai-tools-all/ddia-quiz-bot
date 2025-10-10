package screens

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/abhishek/ddia-clicker/internal/config"
	"github.com/abhishek/ddia-clicker/internal/markdown"
	"github.com/abhishek/ddia-clicker/internal/models"
	"github.com/abhishek/ddia-clicker/internal/tui/session"
)

// ScreenState represents the current screen being displayed
type ScreenState int

const (
	StateWelcome ScreenState = iota
	StateSessionSelect
	StateQuestion
	StateComplete
)

// ImprovedAppModel represents the improved main application state
type ImprovedAppModel struct {
	user              string
	config            *config.TUIConfig
	state             ScreenState
	sessionManager    *session.Manager
	currentSession    *session.Session
	existingSessions  []*session.Session
	questions         []*models.Question
	currentIndex      int
	textarea          textarea.Model
	questionStartTime time.Time
	lastSaveTime      time.Time
	err               error
	quitting          bool
	width             int
	height            int
	renderer          *glamour.TermRenderer
}

// NewImprovedAppModel creates a new improved application model
func NewImprovedAppModel(user string, cfg *config.TUIConfig) ImprovedAppModel {
	ta := textarea.New()
	ta.Placeholder = "Type your answer here..."
	ta.SetWidth(80)
	ta.SetHeight(12)
	ta.CharLimit = 0
	ta.ShowLineNumbers = false

	// Create glamour renderer for markdown
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	return ImprovedAppModel{
		user:           user,
		config:         cfg,
		state:          StateWelcome,
		sessionManager: session.NewManager(cfg.SessionsDir),
		currentIndex:   0,
		textarea:       ta,
		renderer:       renderer,
	}
}

// Init initializes the model
func (m ImprovedAppModel) Init() tea.Cmd {
	return tea.Batch(
		m.loadQuestionsCmd(),
		m.checkExistingSessionsCmd(),
		textarea.Blink,
	)
}

// Update handles messages and updates the model
func (m ImprovedAppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		// Handle global keys
		if msg.String() == "ctrl+c" {
			m.saveBeforeQuit()
			m.quitting = true
			return m, tea.Quit
		}

		// Handle state-specific keys
		switch m.state {
		case StateWelcome:
			if msg.String() == "enter" {
				if len(m.existingSessions) > 0 {
					m.state = StateSessionSelect
				} else {
					return m, m.createNewSessionCmd()
				}
			} else if msg.String() == "q" {
				m.quitting = true
				return m, tea.Quit
			}

		case StateSessionSelect:
			if msg.String() == "n" {
				// Create new session
				return m, m.createNewSessionCmd()
			} else if msg.String() == "r" && len(m.existingSessions) > 0 {
				// Resume most recent session
				m.currentSession = m.existingSessions[0]
				m.state = StateQuestion
				m.restoreSession()
				m.textarea.Focus()
				return m, m.startAutoSaveCmd()
			} else if msg.String() == "q" {
				m.quitting = true
				return m, tea.Quit
			}

		case StateQuestion:
			if msg.String() == "ctrl+s" {
				// Manual save
				return m, m.saveCurrentAnswerCmd()
			} else if msg.String() == "ctrl+n" || msg.String() == "ctrl+enter" {
				// Save and move to next question
				return m, m.moveToNextQuestionCmd()
			} else {
				// Pass to textarea
				var cmd tea.Cmd
				m.textarea, cmd = m.textarea.Update(msg)
				cmds = append(cmds, cmd)
			}

		case StateComplete:
			if msg.String() == "q" {
				m.quitting = true
				return m, tea.Quit
			}
		}

	case questionsLoadedMsg:
		m.questions = msg.questions
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}

	case existingSessionsMsg:
		m.existingSessions = msg.sessions
		if msg.err != nil {
			m.err = msg.err
		}

	case sessionCreatedMsg:
		m.currentSession = msg.session
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		m.state = StateQuestion
		m.questionStartTime = time.Now()
		m.textarea.Focus()
		return m, tea.Batch(
			m.startAutoSaveCmd(),
			textarea.Blink,
		)

	case autoSaveTickMsg:
		if m.state == StateQuestion && m.textarea.Value() != "" {
			// Check if enough time has passed since last save
			if time.Since(m.lastSaveTime) >= m.config.AutoSaveInterval {
				return m, m.saveCurrentAnswerCmd()
			}
		}
		return m, m.waitForAutoSaveTick()

	case answerSavedMsg:
		m.lastSaveTime = time.Now()
		return m, nil

	case nextQuestionMsg:
		// Update UI for next question
		m.currentIndex = msg.newIndex
		m.textarea.Reset()
		m.questionStartTime = time.Now()
		
		if m.currentIndex >= len(m.questions) {
			m.state = StateComplete
		}
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m ImprovedAppModel) View() string {
	if m.quitting {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true).
			Render("Goodbye!\n")
	}

	if m.err != nil {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true).
			Render(fmt.Sprintf("Error: %v\n", m.err))
	}

	switch m.state {
	case StateWelcome:
		return m.renderWelcome()
	case StateSessionSelect:
		return m.renderSessionSelect()
	case StateQuestion:
		return m.renderQuestion()
	case StateComplete:
		return m.renderComplete()
	default:
		return "Unknown state\n"
	}
}

// renderWelcome renders the welcome screen
func (m ImprovedAppModel) renderWelcome() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginTop(2).
		MarginBottom(2).
		Padding(1, 2)

	title := titleStyle.Render("📚 Quiz TUI - Subjective Questions")

	if m.questions == nil {
		return fmt.Sprintf("%s\n\n%s\n", title, "Loading questions...")
	}

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(1)

	info := fmt.Sprintf(
		"Loaded: %d questions\nUser: %s\nMode: Subjective",
		len(m.questions),
		m.user,
	)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		MarginTop(2)

	help := "Press Enter to continue • Press q to quit"

	return fmt.Sprintf("%s\n\n%s\n\n%s\n", title, infoStyle.Render(info), helpStyle.Render(help))
}

// renderSessionSelect renders the session selection screen
func (m ImprovedAppModel) renderSessionSelect() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(2)

	title := titleStyle.Render("Resume or Start New?")

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	info := fmt.Sprintf("Found %d incomplete session(s)", len(m.existingSessions))

	if len(m.existingSessions) > 0 {
		session := m.existingSessions[0]
		sessionInfo := fmt.Sprintf(
			"\nMost recent session:\n  Created: %s\n  Progress: %d/%d questions answered",
			session.Session.CreatedAt.Format("2006-01-02 15:04"),
			session.Session.Answered,
			session.Session.QuestionCount,
		)
		info += sessionInfo
	}

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		MarginTop(2)

	help := "Press r to resume • Press n for new session • Press q to quit"

	return fmt.Sprintf("%s\n\n%s\n\n%s\n", title, infoStyle.Render(info), helpStyle.Render(help))
}

// renderQuestion renders the current question
func (m ImprovedAppModel) renderQuestion() string {
	if m.currentSession == nil || m.currentIndex >= len(m.questions) {
		return "No question to display\n"
	}

	question := m.questions[m.currentIndex]

	// Progress bar
	progressStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	progress := progressStyle.Render(fmt.Sprintf("Question %d of %d", m.currentIndex+1, len(m.questions)))

	// Get question text
	questionText := question.MainQuestion
	
	// Debug: check if question text is empty
	if questionText == "" {
		questionText = "[ERROR: Question text is empty. ID: " + question.ID + "]"
	}

	// Try to render with markdown, but use plain text as fallback
	rendered := questionText
	if m.renderer != nil && questionText != "" {
		if r, err := m.renderer.Render(questionText); err == nil && strings.TrimSpace(r) != "" {
			rendered = r
		}
	}

	// Question display with styling
	questionStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1).
		Width(80)

	questionBox := questionStyle.Render(rendered)

	// Answer section
	answerLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		MarginTop(1).
		Render("Your Answer:")

	textareaView := m.textarea.View()

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(1)

	help := helpStyle.Render("Ctrl+N: Next • Ctrl+S: Save • Ctrl+C: Quit")

	// Auto-save indicator
	saveIndicator := ""
	if time.Since(m.lastSaveTime) < 2*time.Second {
		saveIndicator = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Render(" ✓ Saved")
	}

	return fmt.Sprintf(
		"%s%s\n\n%s\n\n%s\n%s\n\n%s\n",
		progress,
		saveIndicator,
		questionBox,
		answerLabel,
		textareaView,
		help,
	)
}

// renderComplete renders the completion screen
func (m ImprovedAppModel) renderComplete() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("42")).
		MarginTop(2)

	title := titleStyle.Render("🎉 Quiz Complete!")

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(2)

	info := fmt.Sprintf(
		"You answered %d questions.\nSession saved successfully.",
		m.currentSession.Session.Answered,
	)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		MarginTop(2)

	help := helpStyle.Render("Press q to quit")

	return fmt.Sprintf("%s\n\n%s\n\n%s\n", title, infoStyle.Render(info), help)
}

// Helper methods

func (m *ImprovedAppModel) saveBeforeQuit() {
	if m.state == StateQuestion && m.currentSession != nil {
		timeSpent := int(time.Since(m.questionStartTime).Seconds())
		questionID := m.questions[m.currentIndex].ID
		answer := strings.TrimSpace(m.textarea.Value())
		
		if answer != "" {
			m.sessionManager.UpdateResponse(m.currentSession, questionID, answer, timeSpent)
			m.sessionManager.SaveSession(m.currentSession)
		}
	}
}

func (m *ImprovedAppModel) restoreSession() {
	// Find the current question index
	answeredCount := len(m.currentSession.Responses)
	if answeredCount < len(m.questions) {
		m.currentIndex = answeredCount
	} else {
		m.currentIndex = len(m.questions) - 1
	}

	// Load existing answer if any
	if m.currentIndex < len(m.questions) {
		questionID := m.questions[m.currentIndex].ID
		if resp := m.sessionManager.GetResponse(m.currentSession, questionID); resp != nil {
			m.textarea.SetValue(resp.Answer)
		}
	}

	m.questionStartTime = time.Now()
}

// Commands

type questionsLoadedMsg struct {
	questions []*models.Question
	err       error
}

func (m ImprovedAppModel) loadQuestionsCmd() tea.Cmd {
	return func() tea.Msg {
		scanner := markdown.NewScanner(m.config.ContentPath)
		index, err := scanner.ScanQuestions()
		if err != nil {
			return questionsLoadedMsg{err: err}
		}

		byLevel := scanner.GetQuestionsByLevel(index)
		var questions []*models.Question

		levels := []string{"L3", "L4", "L5", "L6", "L7"}
		for _, level := range levels {
			if qs, ok := byLevel[level]; ok {
				questions = append(questions, qs...)
			}
		}

		return questionsLoadedMsg{questions: questions}
	}
}

type existingSessionsMsg struct {
	sessions []*session.Session
	err      error
}

func (m ImprovedAppModel) checkExistingSessionsCmd() tea.Cmd {
	return func() tea.Msg {
		sessions, err := m.sessionManager.ListIncompleteSessions(m.user, "subjective")
		return existingSessionsMsg{sessions: sessions, err: err}
	}
}

type sessionCreatedMsg struct {
	session *session.Session
	err     error
}

func (m ImprovedAppModel) createNewSessionCmd() tea.Cmd {
	return func() tea.Msg {
		sess, err := m.sessionManager.CreateSession(m.user, "subjective", m.questions)
		return sessionCreatedMsg{session: sess, err: err}
	}
}

type autoSaveTickMsg struct{}

func (m ImprovedAppModel) startAutoSaveCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(m.config.AutoSaveInterval)
		return autoSaveTickMsg{}
	}
}

func (m ImprovedAppModel) waitForAutoSaveTick() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(m.config.AutoSaveInterval)
		return autoSaveTickMsg{}
	}
}

type answerSavedMsg struct {
	err error
}

func (m ImprovedAppModel) saveCurrentAnswerCmd() tea.Cmd {
	return func() tea.Msg {
		timeSpent := int(time.Since(m.questionStartTime).Seconds())
		questionID := m.questions[m.currentIndex].ID
		answer := strings.TrimSpace(m.textarea.Value())

		m.sessionManager.UpdateResponse(m.currentSession, questionID, answer, timeSpent)
		err := m.sessionManager.SaveSession(m.currentSession)

		return answerSavedMsg{err: err}
	}
}

type nextQuestionMsg struct {
	newIndex int
}

func (m ImprovedAppModel) moveToNextQuestionCmd() tea.Cmd {
	return func() tea.Msg {
		// Save current answer first
		timeSpent := int(time.Since(m.questionStartTime).Seconds())
		questionID := m.questions[m.currentIndex].ID
		answer := strings.TrimSpace(m.textarea.Value())

		if answer != "" {
			m.sessionManager.UpdateResponse(m.currentSession, questionID, answer, timeSpent)
			m.sessionManager.SaveSession(m.currentSession)
		}

		// Move to next question
		newIndex := m.currentIndex + 1

		// Check if we've completed all questions
		if newIndex >= len(m.questions) {
			m.sessionManager.CompleteSession(m.currentSession)
		}

		return nextQuestionMsg{newIndex: newIndex}
	}
}
