package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TextArea wraps a textarea for answer input
type TextArea struct {
	textarea textarea.Model
	focused  bool
}

// NewTextArea creates a new text area component
func NewTextArea() TextArea {
	ta := textarea.New()
	ta.Placeholder = "Type your answer here..."
	ta.SetWidth(80)
	ta.SetHeight(10)
	ta.Focus()

	return TextArea{
		textarea: ta,
		focused:  true,
	}
}

// Init initializes the text area
func (t TextArea) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles messages for the text area
func (t TextArea) Update(msg tea.Msg) (TextArea, tea.Cmd) {
	var cmd tea.Cmd
	t.textarea, cmd = t.textarea.Update(msg)
	return t, cmd
}

// View renders the text area
func (t TextArea) View() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1)

	return style.Render(t.textarea.View())
}

// SetValue sets the text area value
func (t *TextArea) SetValue(value string) {
	t.textarea.SetValue(value)
}

// Value returns the current value
func (t TextArea) Value() string {
	return strings.TrimSpace(t.textarea.Value())
}

// Focus focuses the text area
func (t *TextArea) Focus() tea.Cmd {
	t.focused = true
	return t.textarea.Focus()
}

// Blur unfocuses the text area
func (t *TextArea) Blur() {
	t.focused = false
	t.textarea.Blur()
}
