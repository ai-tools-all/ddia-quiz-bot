package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MCQ represents a multiple choice question component
type MCQ struct {
	Options      []string
	SelectedIdx  int
	Submitted    bool
	CorrectIdx   int
	ShowFeedback bool
	Explanation  string
}

// NewMCQ creates a new MCQ component
func NewMCQ(options []string, correctAnswer string) *MCQ {
	correctIdx := -1
	// Find correct answer index (e.g., "A)" -> 0, "B)" -> 1)
	if len(correctAnswer) > 0 {
		letter := strings.ToUpper(string(correctAnswer[0]))
		correctIdx = int(letter[0] - 'A')
	}

	return &MCQ{
		Options:      options,
		SelectedIdx:  0,
		Submitted:    false,
		CorrectIdx:   correctIdx,
		ShowFeedback: false,
		Explanation:  "",
	}
}

// SetExplanation sets the MCQ explanation text
func (m *MCQ) SetExplanation(explanation string) {
	m.Explanation = explanation
}

// Update handles messages for the MCQ component
func (m *MCQ) Update(msg tea.Msg) (*MCQ, tea.Cmd) {
	if m.Submitted {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.MoveUp()
		case "down", "j":
			m.MoveDown()
		}
	}

	return m, nil
}

// MoveUp moves selection up
func (m *MCQ) MoveUp() {
	if m.SelectedIdx > 0 {
		m.SelectedIdx--
	}
}

// MoveDown moves selection down
func (m *MCQ) MoveDown() {
	if m.SelectedIdx < len(m.Options)-1 {
		m.SelectedIdx++
	}
}

// Submit submits the current selection
func (m *MCQ) Submit() bool {
	m.Submitted = true
	m.ShowFeedback = true
	return m.IsCorrect()
}

// IsCorrect returns true if the selected answer is correct
func (m *MCQ) IsCorrect() bool {
	return m.SelectedIdx == m.CorrectIdx
}

// GetSelectedLetter returns the selected option letter (A, B, C, D)
func (m *MCQ) GetSelectedLetter() string {
	if m.SelectedIdx < 0 || m.SelectedIdx >= 26 {
		return ""
	}
	return string(rune('A' + m.SelectedIdx))
}

// ToggleExplanation toggles the explanation visibility
func (m *MCQ) ToggleExplanation() {
	m.ShowFeedback = !m.ShowFeedback
}

// View renders the MCQ component
func (m *MCQ) View() string {
	var b strings.Builder

	// Style definitions
	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	correctStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)

	incorrectStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	optionStyle := lipgloss.NewStyle().
		Padding(0, 2)

	// Render options
	for i, option := range m.Options {
		var prefix string
		var style lipgloss.Style

		if m.Submitted {
			// After submission, show correct/incorrect
			if i == m.CorrectIdx {
				prefix = "✓ "
				style = correctStyle
			} else if i == m.SelectedIdx {
				prefix = "✗ "
				style = incorrectStyle
			} else {
				prefix = "  "
				style = lipgloss.NewStyle()
			}
		} else {
			// Before submission, show selection
			if i == m.SelectedIdx {
				prefix = "▸ "
				style = selectedStyle
			} else {
				prefix = "  "
				style = lipgloss.NewStyle()
			}
		}

		optionText := style.Render(fmt.Sprintf("%s%s", prefix, option))
		b.WriteString(optionStyle.Render(optionText))
		b.WriteString("\n")
	}

	// Show feedback after submission
	if m.Submitted && m.ShowFeedback {
		b.WriteString("\n")
		
		var feedbackStyle lipgloss.Style
		var feedbackText string
		
		if m.IsCorrect() {
			feedbackStyle = correctStyle.Copy().Padding(1, 2)
			feedbackText = "✓ Correct!"
		} else {
			feedbackStyle = incorrectStyle.Copy().Padding(1, 2)
			feedbackText = "✗ Incorrect"
		}
		
		b.WriteString(feedbackStyle.Render(feedbackText))
		b.WriteString("\n")

		// Show explanation if available
		if m.Explanation != "" {
			explanationStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				Padding(1, 2).
				Width(80)
			
			b.WriteString("\n")
			b.WriteString(explanationStyle.Render("Explanation:\n" + m.Explanation))
		}
	}

	return b.String()
}
