package screens

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/abhishek/ddia-clicker/internal/config"
	"github.com/abhishek/ddia-clicker/internal/markdown"
	"github.com/abhishek/ddia-clicker/internal/models"
	"github.com/abhishek/ddia-clicker/internal/tui/components"
	"github.com/abhishek/ddia-clicker/internal/tui/session"
)

// ScreenState represents the current screen being displayed
type ScreenState int

const (
	StateWelcome ScreenState = iota
	StateModeSelect
	StateTopicSelect
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
	availableTopics   []markdown.TopicInfo
	selectedTopic     *markdown.TopicInfo
	selectedMode      string // "mcq", "subjective", or "mixed"
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

	// MCQ specific fields
	mcqComponent    *components.MCQ
	currentQType    string // "subjective" or "mcq"
	showExplanation bool   // Show MCQ explanation after answer

	// Topic selection UI state
	topicCursor             int
	topicPageStart          int
	topicsPerPage           int
	loadingQuestions        bool
	topicLoadError          error
	awaitingSessionDecision bool
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

	// Set default mode from config
	defaultMode := cfg.DefaultMode
	if defaultMode == "" {
		defaultMode = "mixed"
	}
	if cfg.ChaptersRootPath == "" {
		defaultMode = "subjective"
	}

	return ImprovedAppModel{
		user:           user,
		config:         cfg,
		state:          StateWelcome,
		sessionManager: session.NewManager(cfg.SessionsDir),
		selectedMode:   defaultMode,
		currentIndex:   0,
		textarea:       ta,
		renderer:       renderer,
		topicsPerPage:  8,
	}
}

// Init initializes the model
func (m ImprovedAppModel) Init() tea.Cmd {
	// Check if using topic selection mode
	if m.config.ChaptersRootPath != "" {
		return tea.Batch(
			m.discoverTopicsCmd(),
			textarea.Blink,
		)
	}
	// Fallback to legacy single-topic mode
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
				if m.config.ChaptersRootPath != "" {
					if len(m.availableTopics) > 0 {
						m.state = StateModeSelect
					}
				} else if len(m.existingSessions) > 0 {
					m.state = StateSessionSelect
				} else {
					return m, m.createNewSessionCmd()
				}
			} else if msg.String() == "q" {
				m.quitting = true
				return m, tea.Quit
			}

		case StateModeSelect:
			// Handle mode selection keys (using letters to avoid conflict with topic numbers)
			switch msg.String() {
			case "m":
				m.selectedMode = "mcq"
				m.resetTopicSelectionState()
				m.state = StateTopicSelect
			case "s":
				m.selectedMode = "subjective"
				m.resetTopicSelectionState()
				m.state = StateTopicSelect
			case "b", "x":
				m.selectedMode = "mixed"
				m.resetTopicSelectionState()
				m.state = StateTopicSelect
			case "enter":
				if len(m.availableTopics) > 0 {
					m.resetTopicSelectionState()
					m.state = StateTopicSelect
				}
			case "esc", "backspace":
				m.state = StateWelcome
			case "q":
				m.quitting = true
				return m, tea.Quit
			}

		case StateTopicSelect:
			if cmd := m.handleTopicSelectKeys(msg); cmd != nil {
				return m, cmd
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
			if m.currentQType == "mcq" && m.mcqComponent != nil {
				// MCQ-specific handling
				if msg.String() == "enter" || msg.String() == " " {
					if !m.mcqComponent.Submitted {
						// Submit MCQ answer
						m.mcqComponent.Submit()
						return m, m.saveMCQAnswerCmd()
					} else {
						// Move to next question after viewing feedback
						return m, m.moveToNextQuestionCmd()
					}
				} else if msg.String() == "e" && m.mcqComponent.Submitted {
					// Toggle explanation
					m.mcqComponent.ToggleExplanation()
				} else if msg.String() == "n" && m.mcqComponent.Submitted {
					// Next question
					return m, m.moveToNextQuestionCmd()
				} else {
					// Pass to MCQ component for navigation
					var cmd tea.Cmd
					m.mcqComponent, cmd = m.mcqComponent.Update(msg)
					cmds = append(cmds, cmd)
				}
			} else {
				// Subjective question handling
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
			}

		case StateComplete:
			if msg.String() == "q" {
				m.quitting = true
				return m, tea.Quit
			}
		}

	case topicsDiscoveredMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		m.availableTopics = msg.topics
		m.resetTopicSelectionState()
		if len(m.availableTopics) == 0 && m.config.ChaptersRootPath != "" {
			m.topicLoadError = fmt.Errorf("no topics found in %s", m.config.ChaptersRootPath)
		}

	case questionsLoadedMsg:
		if msg.topicName != "" {
			if m.selectedTopic == nil || m.selectedTopic.Name != msg.topicName {
				break
			}
		}
		if msg.mode != "" && msg.mode != m.selectedMode && m.config.ChaptersRootPath != "" {
			break
		}
		m.loadingQuestions = false
		if msg.err != nil {
			if msg.topicName != "" {
				m.topicLoadError = msg.err
				m.awaitingSessionDecision = false
				m.questions = nil
				return m, nil
			}
			m.err = msg.err
			return m, tea.Quit
		}
		m.topicLoadError = nil
		m.questions = msg.questions
		m.currentIndex = 0
		if msg.topicName != "" {
			if cmd := m.decideNextStep(); cmd != nil {
				return m, cmd
			}
		} else {
			if m.existingSessions == nil {
				return m, nil
			}
			if len(m.existingSessions) > 0 {
				m.state = StateSessionSelect
			} else {
				return m, m.createNewSessionCmd()
			}
		}

	case existingSessionsMsg:
		if msg.mode != "" && msg.mode != m.selectedMode && !(m.config.ChaptersRootPath == "" && msg.mode == "subjective") {
			break
		}
		if msg.topic != "" {
			if m.selectedTopic == nil || m.selectedTopic.Name != msg.topic {
				break
			}
		} else if m.selectedTopic != nil && m.config.ChaptersRootPath != "" {
			break
		}
		if msg.err != nil {
			if msg.topic != "" {
				m.topicLoadError = msg.err
				m.awaitingSessionDecision = false
				return m, nil
			}
			m.err = msg.err
			return m, nil
		}
		if msg.sessions == nil {
			m.existingSessions = []*session.Session{}
		} else {
			m.existingSessions = msg.sessions
		}
		if cmd := m.decideNextStep(); cmd != nil {
			return m, cmd
		}
		if m.config.ChaptersRootPath == "" && m.questions != nil {
			if len(m.existingSessions) > 0 {
				m.state = StateSessionSelect
			} else if m.state == StateWelcome {
				return m, m.createNewSessionCmd()
			}
		}

	case sessionCreatedMsg:
		m.currentSession = msg.session
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		m.state = StateQuestion
		m.questionStartTime = time.Now()
		m.initializeQuestionComponent()
		if m.currentQType == "mcq" {
			m.textarea.Blur()
		} else {
			m.textarea.Focus()
		}
		return m, tea.Batch(
			m.startAutoSaveCmd(),
			textarea.Blink,
		)

	case autoSaveTickMsg:
		if m.state == StateQuestion && m.currentQType != "mcq" && m.textarea.Value() != "" {
			// Check if enough time has passed since last save
			if time.Since(m.lastSaveTime) >= m.config.AutoSaveInterval {
				return m, m.saveCurrentAnswerCmd()
			}
		}
		return m, m.waitForAutoSaveTick()

	case answerSavedMsg:
		if m.currentQType != "mcq" {
			m.lastSaveTime = time.Now()
		}
		return m, nil

	case nextQuestionMsg:
		// Update UI for next question
		m.currentIndex = msg.newIndex
		m.textarea.Reset()
		m.questionStartTime = time.Now()

		if m.currentIndex >= len(m.questions) {
			m.state = StateComplete
		} else {
			// Initialize component for new question
			m.initializeQuestionComponent()
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
	case StateModeSelect:
		return m.renderModeSelect()
	case StateTopicSelect:
		return m.renderTopicSelect()
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

	title := titleStyle.Render("📚 Quiz TUI")

	if m.config.ChaptersRootPath != "" && m.availableTopics == nil {
		return fmt.Sprintf("%s\n\n%s\n", title, "Discovering topics...")
	}

	if m.config.ChaptersRootPath == "" && m.questions == nil {
		return fmt.Sprintf("%s\n\n%s\n", title, "Loading questions...")
	}

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(1)

	var info string
	if m.config.ChaptersRootPath != "" {
		if len(m.availableTopics) > 0 {
			info = fmt.Sprintf(
				"User: %s\nAvailable topics: %d\nDefault mode: %s",
				m.user,
				len(m.availableTopics),
				formatModeDisplay(m.selectedMode),
			)
		} else {
			info = fmt.Sprintf("User: %s\nNo topics discovered yet", m.user)
		}
	} else if m.questions != nil {
		info = fmt.Sprintf(
			"User: %s\nLoaded questions: %d\nMode: %s",
			m.user,
			len(m.questions),
			formatModeDisplay(m.selectedMode),
		)
	} else {
		info = fmt.Sprintf("User: %s\nInitializing...", m.user)
	}

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		MarginTop(2)

	help := "Press Enter to choose a quiz mode • Press q to quit"

	return fmt.Sprintf("%s\n\n%s\n\n%s\n", title, infoStyle.Render(info), helpStyle.Render(help))
}

// renderModeSelect renders the mode selection screen
func (m ImprovedAppModel) renderModeSelect() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginTop(2).
		MarginBottom(2)

	title := titleStyle.Render("🎯 Select Question Mode")

	// Mode descriptions
	modeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		MarginLeft(2)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginLeft(6)

	var modeList strings.Builder
	modeList.WriteString("Choose your quiz mode:\n\n")

	// MCQ mode
	modeList.WriteString(modeStyle.Render("  [M] MCQ Questions\n"))
	modeList.WriteString(descStyle.Render("      Multiple choice questions with instant feedback\n\n"))

	// Subjective mode
	modeList.WriteString(modeStyle.Render("  [S] Subjective Questions\n"))
	modeList.WriteString(descStyle.Render("      Open-ended questions requiring written answers\n\n"))

	// Mixed mode
	modeList.WriteString(modeStyle.Render("  [B] Both (Mixed Mode)\n"))
	modeList.WriteString(descStyle.Render("      Both MCQ and subjective questions\n\n"))

	// Current selection indicator
	currentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true).
		MarginTop(1)

	current := currentStyle.Render(fmt.Sprintf("Current: %s Mode", formatModeDisplay(m.selectedMode)))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(2)

	help := helpStyle.Render("M/S/B to select • Enter to continue • Esc to go back • q to quit")

	return fmt.Sprintf("%s\n\n%s\n%s\n\n%s\n", title, modeList.String(), current, help)
}

// renderTopicSelect renders the topic selection screen
func (m ImprovedAppModel) renderTopicSelect() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginTop(2).
		MarginBottom(2)

	title := titleStyle.Render("📚 Select Your Topic")

	pageSize := m.topicsPerPage
	if pageSize <= 0 {
		pageSize = 8
	}

	modeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)
	modeText := modeStyle.Render(fmt.Sprintf("Mode: %s", formatModeDisplay(m.selectedMode)))

	if m.availableTopics == nil {
		loadingStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			MarginTop(1)
		helpStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			MarginTop(2)
		help := helpStyle.Render("Discovering topics... Please wait or press M to change mode.")
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			modeText,
			"",
			loadingStyle.Render("Loading topics..."),
			help,
		)
	}

	if len(m.availableTopics) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			MarginTop(1)
		helpStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			MarginTop(2)
		help := helpStyle.Render("No topics available. Press M to change mode or q to quit.")
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			modeText,
			"",
			emptyStyle.Render("No topics discovered for this mode."),
			help,
		)
	}

	start := m.topicPageStart
	if start < 0 {
		start = 0
	}
	if start >= len(m.availableTopics) {
		start = (len(m.availableTopics) / pageSize) * pageSize
		if start >= len(m.availableTopics) {
			start = 0
		}
	}
	end := start + pageSize
	if end > len(m.availableTopics) {
		end = len(m.availableTopics)
	}

	var listBuilder strings.Builder
	cursorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true)
	nameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))
	selectedNameStyle := nameStyle.Copy().
		Foreground(lipgloss.Color("212")).
		Bold(true)
	countStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
	lineStyle := lipgloss.NewStyle().
		MarginLeft(2)

	for i := start; i < end; i++ {
		topic := m.availableTopics[i]
		isCursor := i == m.topicCursor
		indicator := "  "
		if isCursor {
			indicator = cursorStyle.Render("▸ ")
		}

		var countText string
		switch m.selectedMode {
		case "mixed":
			countText = fmt.Sprintf("(MCQ: %d | Subjective: %d)", topic.MCQCount, topic.SubjectiveCount)
		case "mcq":
			countText = fmt.Sprintf("(%d questions)", topic.MCQCount)
		default:
			countText = fmt.Sprintf("(%d questions)", topic.SubjectiveCount)
		}

		name := nameStyle.Render(topic.DisplayName)
		if isCursor {
			name = selectedNameStyle.Render(topic.DisplayName)
		}

		line := fmt.Sprintf("%s%s %s", indicator, name, countStyle.Render(countText))
		listBuilder.WriteString(lineStyle.Render(line))
		listBuilder.WriteRune('\n')
	}

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(1)

	statusLines := []string{}
	if m.loadingQuestions && m.selectedTopic != nil {
		loadingStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			MarginTop(1)
		statusLines = append(statusLines, loadingStyle.Render(fmt.Sprintf("Loading %s questions...", m.selectedTopic.DisplayName)))
	}
	if m.topicLoadError != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			MarginTop(1)
		statusLines = append(statusLines, errorStyle.Render(m.topicLoadError.Error()))
	}

	pageInfo := ""
	if len(m.availableTopics) > pageSize {
		currentPage := (m.topicCursor / pageSize) + 1
		totalPages := (len(m.availableTopics) + pageSize - 1) / pageSize
		pageInfo = helpStyle.Render(fmt.Sprintf("Page %d of %d", currentPage, totalPages))
	}

	help := helpStyle.Render("↑/↓ move • PgUp/PgDn page • Enter select • M change mode • q quit")

	sections := []string{title, modeText, "", strings.TrimRight(listBuilder.String(), "\n")}
	if pageInfo != "" {
		sections = append(sections, pageInfo)
	}
	sections = append(sections, statusLines...)
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
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

	var infoLines []string
	if m.selectedTopic != nil {
		infoLines = append(infoLines, fmt.Sprintf("Topic: %s", m.selectedTopic.DisplayName))
	}
	infoLines = append(infoLines, fmt.Sprintf("Mode: %s", formatModeDisplay(m.selectedMode)))
	infoLines = append(infoLines, fmt.Sprintf("Found %d incomplete session(s)", len(m.existingSessions)))
	info := strings.Join(infoLines, "\n")

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

	// Progress bar with topic and level info
	progressStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	progressText := fmt.Sprintf("Question %d of %d", m.currentIndex+1, len(m.questions))
	if m.selectedTopic != nil {
		progressText = fmt.Sprintf("%s | Level: %s | %s", m.selectedTopic.DisplayName, question.Level, progressText)
	}
	progress := progressStyle.Render(progressText)

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

	// Answer section - different for MCQ vs subjective
	var answerSection string
	var help string

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(1)

	if m.currentQType == "mcq" && m.mcqComponent != nil {
		// MCQ answer section
		answerLabel := lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			MarginTop(1).
			Render("Select Answer:")

		mcqView := m.mcqComponent.View()
		answerSection = fmt.Sprintf("%s\n%s", answerLabel, mcqView)

		// MCQ-specific help text
		if !m.mcqComponent.Submitted {
			help = helpStyle.Render("↑↓: Navigate • Enter/Space: Submit • Ctrl+C: Quit")
		} else {
			help = helpStyle.Render("N: Next Question • E: Toggle Explanation • Ctrl+C: Quit")
		}
	} else {
		// Subjective answer section
		answerLabel := lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			MarginTop(1).
			Render("Your Answer:")

		textareaView := m.textarea.View()
		answerSection = fmt.Sprintf("%s\n%s", answerLabel, textareaView)

		// Subjective-specific help text
		help = helpStyle.Render("Ctrl+N: Next • Ctrl+S: Save • Ctrl+C: Quit")
	}

	// Auto-save indicator (only for subjective)
	saveIndicator := ""
	if m.currentQType == "subjective" && time.Since(m.lastSaveTime) < 2*time.Second {
		saveIndicator = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Render(" ✓ Saved")
	}

	return fmt.Sprintf(
		"%s%s\n\n%s\n\n%s\n\n%s\n",
		progress,
		saveIndicator,
		questionBox,
		answerSection,
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

func formatModeDisplay(mode string) string {
	switch mode {
	case "mcq":
		return "MCQ"
	case "subjective":
		return "Subjective"
	case "mixed":
		return "Mixed"
	}
	if mode == "" {
		return "Mixed"
	}
	runes := []rune(mode)
	return strings.ToUpper(string(runes[0])) + string(runes[1:])
}

func (m *ImprovedAppModel) resetTopicSelectionState() {
	if len(m.availableTopics) == 0 {
		m.topicCursor = 0
		m.topicPageStart = 0
	} else {
		if m.topicCursor < 0 || m.topicCursor >= len(m.availableTopics) {
			m.topicCursor = 0
		}
		if m.topicsPerPage <= 0 {
			m.topicsPerPage = 8
		}
		m.topicPageStart = (m.topicCursor / m.topicsPerPage) * m.topicsPerPage
	}
	m.selectedTopic = nil
	m.questions = nil
	m.existingSessions = nil
	m.currentSession = nil
	m.currentIndex = 0
	m.loadingQuestions = false
	m.topicLoadError = nil
	m.awaitingSessionDecision = false
	m.mcqComponent = nil
	m.textarea.Reset()
	m.textarea.Blur()
}

func (m *ImprovedAppModel) ensureTopicCursorVisible() {
	if len(m.availableTopics) == 0 {
		m.topicPageStart = 0
		return
	}
	if m.topicsPerPage <= 0 {
		m.topicsPerPage = 8
	}
	if m.topicCursor < 0 {
		m.topicCursor = 0
	}
	if m.topicCursor >= len(m.availableTopics) {
		m.topicCursor = len(m.availableTopics) - 1
	}
	if m.topicCursor < m.topicPageStart {
		m.topicPageStart = m.topicCursor
	}
	if m.topicCursor >= m.topicPageStart+m.topicsPerPage {
		m.topicPageStart = m.topicCursor - (m.topicsPerPage - 1)
	}
	maxStart := len(m.availableTopics) - m.topicsPerPage
	if maxStart < 0 {
		maxStart = 0
	}
	if m.topicPageStart > maxStart {
		m.topicPageStart = maxStart
	}
	if m.topicPageStart < 0 {
		m.topicPageStart = 0
	}
}

func (m *ImprovedAppModel) handleTopicSelectKeys(msg tea.KeyMsg) tea.Cmd {
	key := msg.String()
	switch key {
	case "up", "k":
		if len(m.availableTopics) == 0 {
			return nil
		}
		if m.topicCursor > 0 {
			m.topicCursor--
			m.ensureTopicCursorVisible()
		}
	case "down", "j":
		if len(m.availableTopics) == 0 {
			return nil
		}
		if m.topicCursor < len(m.availableTopics)-1 {
			m.topicCursor++
			m.ensureTopicCursorVisible()
		}
	case "pgdown", "right", "l":
		if len(m.availableTopics) == 0 {
			return nil
		}
		if m.topicsPerPage <= 0 {
			m.topicsPerPage = 8
		}
		next := m.topicCursor + m.topicsPerPage
		if next >= len(m.availableTopics) {
			next = len(m.availableTopics) - 1
		}
		if next != m.topicCursor {
			m.topicCursor = next
			m.topicPageStart = (m.topicCursor / m.topicsPerPage) * m.topicsPerPage
			m.ensureTopicCursorVisible()
		}
	case "pgup", "left", "h":
		if len(m.availableTopics) == 0 {
			return nil
		}
		if m.topicsPerPage <= 0 {
			m.topicsPerPage = 8
		}
		prev := m.topicCursor - m.topicsPerPage
		if prev < 0 {
			prev = 0
		}
		if prev != m.topicCursor {
			m.topicCursor = prev
			m.topicPageStart = (m.topicCursor / m.topicsPerPage) * m.topicsPerPage
			m.ensureTopicCursorVisible()
		}
	case "home", "g":
		if len(m.availableTopics) == 0 {
			return nil
		}
		m.topicCursor = 0
		m.topicPageStart = 0
		m.ensureTopicCursorVisible()
	case "end", "G":
		if len(m.availableTopics) == 0 {
			return nil
		}
		m.topicCursor = len(m.availableTopics) - 1
		if m.topicsPerPage <= 0 {
			m.topicsPerPage = 8
		}
		m.topicPageStart = (m.topicCursor / m.topicsPerPage) * m.topicsPerPage
		m.ensureTopicCursorVisible()
	case "enter":
		if m.loadingQuestions || len(m.availableTopics) == 0 {
			return nil
		}
		m.selectedTopic = &m.availableTopics[m.topicCursor]
		m.loadingQuestions = true
		m.topicLoadError = nil
		m.awaitingSessionDecision = true
		m.questions = nil
		m.existingSessions = nil
		m.currentSession = nil
		m.currentIndex = 0
		return tea.Batch(
			m.loadTopicQuestionsCmd(),
			m.checkTopicSessionsCmd(),
		)
	case "m", "backspace", "esc":
		m.state = StateModeSelect
		m.loadingQuestions = false
		m.awaitingSessionDecision = false
	case "q":
		m.quitting = true
		return tea.Quit
	}
	return nil
}

func (m *ImprovedAppModel) decideNextStep() tea.Cmd {
	if !m.awaitingSessionDecision {
		return nil
	}
	if m.loadingQuestions || m.selectedTopic == nil {
		return nil
	}
	if m.questions == nil || len(m.questions) == 0 {
		return nil
	}
	if m.existingSessions == nil {
		return nil
	}
	m.awaitingSessionDecision = false
	if len(m.existingSessions) > 0 {
		m.state = StateSessionSelect
		return nil
	}
	m.currentIndex = 0
	return m.createNewSessionCmd()
}

func (m *ImprovedAppModel) saveBeforeQuit() {
	if m.state == StateQuestion && m.currentSession != nil {
		timeSpent := int(time.Since(m.questionStartTime).Seconds())
		question := m.questions[m.currentIndex]
		questionID := question.ID
		answer := strings.TrimSpace(m.textarea.Value())

		if answer != "" {
			qType := question.Type
			if qType == "" {
				qType = "subjective"
			}
			m.sessionManager.UpdateResponse(m.currentSession, questionID, qType, answer, timeSpent)
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

type topicsDiscoveredMsg struct {
	topics []markdown.TopicInfo
	err    error
}

func (m ImprovedAppModel) discoverTopicsCmd() tea.Cmd {
	return func() tea.Msg {
		scanner := markdown.NewScanner("")
		topics, err := scanner.DiscoverTopics(m.config.ChaptersRootPath)
		return topicsDiscoveredMsg{topics: topics, err: err}
	}
}

type questionsLoadedMsg struct {
	questions []*models.Question
	err       error
	mode      string
	topicName string
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

		return questionsLoadedMsg{questions: questions, mode: "subjective"}
	}
}

func (m ImprovedAppModel) loadTopicQuestionsCmd() tea.Cmd {
	return func() tea.Msg {
		if m.selectedTopic == nil {
			return questionsLoadedMsg{err: fmt.Errorf("no topic selected"), mode: m.selectedMode}
		}

		mode := m.selectedMode
		topicName := m.selectedTopic.Name
		topicDisplay := m.selectedTopic.DisplayName

		var subjectiveQuestions []*models.Question
		var mcqQuestions []*models.Question

		// Load subjective questions if mode is subjective or mixed
		if mode == "subjective" || mode == "mixed" {
			scanner := markdown.NewScanner(m.selectedTopic.Path)
			if index, err := scanner.ScanQuestions(); err == nil {
				subjectiveQuestions = scanner.GetProgressiveQuestions(index)
				for _, q := range subjectiveQuestions {
					if q != nil && q.Type == "" {
						q.Type = "subjective"
					}
				}
			}
		}

		// Load MCQ questions if mode is mcq or mixed
		if mode == "mcq" || mode == "mixed" {
			topicPath := filepath.Dir(m.selectedTopic.Path)
			mcqPath := filepath.Join(topicPath, "mcq")
			if info, err := os.Stat(mcqPath); err == nil && info.IsDir() {
				scanner := markdown.NewScanner(mcqPath)
				if index, err := scanner.ScanQuestions(); err == nil {
					ids := make([]string, 0, len(index))
					for id := range index {
						ids = append(ids, id)
					}
					sort.Strings(ids)
					for _, id := range ids {
						q := index[id]
						if q != nil && q.Type == "" {
							q.Type = "mcq"
						}
						mcqQuestions = append(mcqQuestions, q)
					}
				}
			}
		}

		var allQuestions []*models.Question
		switch mode {
		case "subjective":
			if len(subjectiveQuestions) == 0 {
				return questionsLoadedMsg{
					err:       fmt.Errorf("no subjective questions available for %s", topicDisplay),
					mode:      mode,
					topicName: topicName,
				}
			}
			allQuestions = append(allQuestions, subjectiveQuestions...)
		case "mcq":
			if len(mcqQuestions) == 0 {
				return questionsLoadedMsg{
					err:       fmt.Errorf("no mcq questions available for %s", topicDisplay),
					mode:      mode,
					topicName: topicName,
				}
			}
			allQuestions = append(allQuestions, mcqQuestions...)
		case "mixed":
			if len(subjectiveQuestions)+len(mcqQuestions) == 0 {
				return questionsLoadedMsg{
					err:       fmt.Errorf("no questions available for %s", topicDisplay),
					mode:      mode,
					topicName: topicName,
				}
			}
			allQuestions = append(allQuestions, subjectiveQuestions...)
			allQuestions = append(allQuestions, mcqQuestions...)
		default:
			return questionsLoadedMsg{
				err:       fmt.Errorf("unsupported mode: %s", mode),
				mode:      mode,
				topicName: topicName,
			}
		}

		return questionsLoadedMsg{
			questions: allQuestions,
			mode:      mode,
			topicName: topicName,
		}
	}
}

type existingSessionsMsg struct {
	sessions []*session.Session
	err      error
	mode     string
	topic    string
}

func (m ImprovedAppModel) checkExistingSessionsCmd() tea.Cmd {
	return func() tea.Msg {
		mode := m.selectedMode
		if mode == "" {
			mode = "subjective"
		}
		sessions, err := m.sessionManager.ListIncompleteSessions(m.user, mode)
		return existingSessionsMsg{sessions: sessions, err: err, mode: mode}
	}
}

func (m ImprovedAppModel) checkTopicSessionsCmd() tea.Cmd {
	return func() tea.Msg {
		if m.selectedTopic == nil {
			return existingSessionsMsg{sessions: []*session.Session{}, err: nil, mode: m.selectedMode}
		}
		mode := m.selectedMode
		if mode == "" {
			mode = "subjective"
		}
		sessions, err := m.sessionManager.ListIncompleteSessionsForTopic(m.user, mode, m.selectedTopic.Name)
		return existingSessionsMsg{sessions: sessions, err: err, mode: mode, topic: m.selectedTopic.Name}
	}
}

type sessionCreatedMsg struct {
	session *session.Session
	err     error
}

func (m ImprovedAppModel) createNewSessionCmd() tea.Cmd {
	return func() tea.Msg {
		var sess *session.Session
		var err error
		mode := m.selectedMode
		if mode == "" {
			mode = "subjective"
		}

		if m.selectedTopic != nil {
			// Create session with topic information
			sess, err = m.sessionManager.CreateSessionWithTopic(
				m.user,
				mode,
				m.selectedTopic.Name,
				m.selectedTopic.DisplayName,
				m.questions,
			)
		} else {
			// Legacy mode without topic
			sess, err = m.sessionManager.CreateSession(m.user, mode, m.questions)
		}

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

// initializeQuestionComponent initializes the appropriate input component based on question type
func (m *ImprovedAppModel) initializeQuestionComponent() {
	if m.currentIndex >= len(m.questions) {
		return
	}

	question := m.questions[m.currentIndex]
	m.currentQType = question.Type

	// Default to subjective if type not specified
	if m.currentQType == "" {
		m.currentQType = "subjective"
	}

	if m.currentQType == "mcq" {
		// Initialize MCQ component
		m.mcqComponent = components.NewMCQ(question.Options, question.Answer)
		m.mcqComponent.SetExplanation(question.Explanation)
		m.showExplanation = false
		m.textarea.Blur()
	} else {
		// Clear MCQ component for subjective questions
		m.mcqComponent = nil

		// Load existing answer if resuming
		if m.currentSession != nil {
			existingResponse := m.sessionManager.GetResponse(m.currentSession, question.ID)
			if existingResponse != nil {
				m.textarea.SetValue(existingResponse.Answer)
			} else {
				m.textarea.Reset()
			}
		}
		m.textarea.Focus()
	}
}

func (m ImprovedAppModel) saveCurrentAnswerCmd() tea.Cmd {
	return func() tea.Msg {
		if m.currentSession == nil || m.currentIndex >= len(m.questions) {
			return answerSavedMsg{}
		}
		question := m.questions[m.currentIndex]
		if question.Type == "mcq" {
			return answerSavedMsg{}
		}
		timeSpent := int(time.Since(m.questionStartTime).Seconds())
		questionID := question.ID
		answer := strings.TrimSpace(m.textarea.Value())

		qType := question.Type
		if qType == "" {
			qType = "subjective"
		}
		m.sessionManager.UpdateResponse(m.currentSession, questionID, qType, answer, timeSpent)
		err := m.sessionManager.SaveSession(m.currentSession)

		return answerSavedMsg{err: err}
	}
}

func (m ImprovedAppModel) saveMCQAnswerCmd() tea.Cmd {
	return func() tea.Msg {
		timeSpent := int(time.Since(m.questionStartTime).Seconds())
		questionID := m.questions[m.currentIndex].ID

		if m.currentSession == nil {
			return answerSavedMsg{err: fmt.Errorf("session not initialized")}
		}
		if m.mcqComponent == nil {
			return answerSavedMsg{err: fmt.Errorf("MCQ component not initialized")}
		}

		selectedLetter := m.mcqComponent.GetSelectedLetter()
		isCorrect := m.mcqComponent.IsCorrect()
		correct := isCorrect

		existing := m.sessionManager.GetResponse(m.currentSession, questionID)
		if existing != nil {
			existing.QuestionType = "mcq"
			existing.Answer = selectedLetter
			existing.IsCorrect = &correct
			existing.SelectedOption = selectedLetter
			existing.TimeSpentSeconds += timeSpent
			existing.UpdatedAt = time.Now()
		} else {
			response := session.Response{
				QuestionID:       questionID,
				QuestionType:     "mcq",
				Answer:           selectedLetter,
				IsCorrect:        &correct,
				SelectedOption:   selectedLetter,
				UpdatedAt:        time.Now(),
				TimeSpentSeconds: timeSpent,
			}
			m.currentSession.Responses = append(m.currentSession.Responses, response)
		}
		m.currentSession.Session.Answered = len(m.currentSession.Responses)

		err := m.sessionManager.SaveSession(m.currentSession)
		return answerSavedMsg{err: err}
	}
}

type nextQuestionMsg struct {
	newIndex int
}

func (m ImprovedAppModel) moveToNextQuestionCmd() tea.Cmd {
	return func() tea.Msg {
		if m.currentSession != nil && m.currentIndex < len(m.questions) {
			question := m.questions[m.currentIndex]
			if question.Type != "mcq" {
				timeSpent := int(time.Since(m.questionStartTime).Seconds())
				answer := strings.TrimSpace(m.textarea.Value())
				if answer != "" {
					qType := question.Type
					if qType == "" {
						qType = "subjective"
					}
					m.sessionManager.UpdateResponse(m.currentSession, question.ID, qType, answer, timeSpent)
					m.sessionManager.SaveSession(m.currentSession)
				}
			}
		}

		// Move to next question
		newIndex := m.currentIndex + 1

		// Check if we've completed all questions
		if newIndex >= len(m.questions) && m.currentSession != nil {
			m.sessionManager.CompleteSession(m.currentSession)
		}

		return nextQuestionMsg{newIndex: newIndex}
	}
}
