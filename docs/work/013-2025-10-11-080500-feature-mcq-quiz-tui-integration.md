# MCQ Integration with Quiz TUI - Detailed Code Plan

## Date: 2025-10-11 08:05:00

## Task
Add MCQ (Multiple Choice Question) support to the existing quiz-tui system

## Current State Analysis

### Existing Infrastructure
1. **Spaced Repetition System (SRS)**
   - `internal/srs/` - Complete SRS implementation with SM-2+ algorithm
   - Cards track questions with intervals, ease factors, and performance
   - Already handles question IDs, topics, levels, categories

2. **TUI Architecture**
   - `internal/tui/screens/app.go` - Main app model with state machine
   - `internal/tui/session/` - Session management for quiz progress
   - States: Welcome → TopicSelect → SessionSelect → Question → Complete
   - Currently only supports subjective questions with textarea

3. **Question Model**
   - `internal/models/question.go` - Base Question struct
   - Currently lacks MCQ-specific fields (options, answer)

4. **Markdown Parser**
   - `internal/markdown/parser.go` - Parses question markdown files
   - Needs extension for MCQ format parsing

## Detailed Implementation Plan

### Phase 1: Extend Data Models

#### 1.1 Update Question Model
```go
// internal/models/question.go
type Question struct {
    // ... existing fields ...
    
    // MCQ specific fields
    Type        string   `yaml:"type" json:"type"`         // "subjective" or "mcq"
    Options     []string `yaml:"options" json:"options"`   // MCQ options
    Answer      string   `yaml:"answer" json:"answer"`     // Correct answer (A, B, C, etc.)
    Explanation string   `yaml:"explanation" json:"explanation"` // MCQ explanation
    Hook        string   `yaml:"hook" json:"hook"`         // Engagement hook
}
```

#### 1.2 Extend SRS Card for MCQ Performance
```go
// internal/srs/card.go
type Card struct {
    // ... existing fields ...
    
    // MCQ specific metrics
    QuestionType    string  `json:"question_type"`    // "subjective" or "mcq"
    MCQAccuracy     float64 `json:"mcq_accuracy"`     // MCQ-specific accuracy
    LastMCQChoice   string  `json:"last_mcq_choice"`  // Track last selected option
    IncorrectChoices []string `json:"incorrect_choices"` // Track wrong patterns
}
```

### Phase 2: Parser Extensions

#### 2.1 Update Markdown Parser
```go
// internal/markdown/parser.go
// Add MCQ parsing logic to parseBody() and saveSection()

func (p *Parser) saveSection(section, content string, question *models.Question) {
    // ... existing sections ...
    
    // MCQ sections
    if section == "question" {
        question.MainQuestion = content
    } else if section == "options" {
        question.Options = p.parseOptions(content)
    } else if section == "answer" {
        question.Answer = strings.TrimSpace(content)
    } else if section == "explanation" {
        question.Explanation = content
    } else if section == "hook" {
        question.Hook = content
    }
}

func (p *Parser) parseOptions(content string) []string {
    // Parse "- A) option text" format
}
```

### Phase 3: TUI Components

#### 3.1 Create MCQ Display Component
```go
// internal/tui/components/mcq.go
package components

type MCQComponent struct {
    Options      []string
    SelectedIdx  int
    Submitted    bool
    CorrectIdx   int
    ShowFeedback bool
}

func (m *MCQComponent) View() string {
    // Render MCQ options with selection state
    // Show feedback after submission
}

func (m *MCQComponent) Update(msg tea.Msg) (*MCQComponent, tea.Cmd) {
    // Handle arrow keys for selection
    // Handle enter for submission
    // Handle space for toggle selection
}
```

#### 3.2 Update App Model
```go
// internal/tui/screens/app.go
type ImprovedAppModel struct {
    // ... existing fields ...
    
    // MCQ specific fields
    mcqComponent    *components.MCQComponent
    currentQType    string  // "subjective" or "mcq"
    mcqAnswer       string  // Selected MCQ answer
    showExplanation bool    // Show MCQ explanation after answer
}
```

### Phase 4: State Machine Updates

#### 4.1 Modify Question State Handler
```go
// internal/tui/screens/app.go

func (m ImprovedAppModel) updateQuestionState(msg tea.Msg) (tea.Model, tea.Cmd) {
    question := m.questions[m.currentIndex]
    
    if question.Type == "mcq" {
        return m.handleMCQInput(msg)
    } else {
        return m.handleSubjectiveInput(msg)
    }
}

func (m ImprovedAppModel) handleMCQInput(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up", "k":
            m.mcqComponent.MoveUp()
        case "down", "j":
            m.mcqComponent.MoveDown()
        case "enter", " ":
            return m.submitMCQAnswer()
        case "e":
            if m.mcqComponent.Submitted {
                m.showExplanation = !m.showExplanation
            }
        }
    }
}
```

### Phase 5: SRS Integration

#### 5.1 Update Review Recording
```go
// internal/srs/scheduler.go

func (s *Scheduler) RecordMCQReview(questionID string, correct bool, timeSpent int) (*ReviewResult, error) {
    quality := QualityWrong
    if correct {
        quality = QualityGood  // Can be adjusted based on speed
    }
    
    return s.RecordReview(questionID, quality, timeSpent, 0)
}
```

#### 5.2 MCQ-specific Metrics
```go
// internal/srs/card.go

func (c *Card) RecordMCQAttempt(selected string, correct bool, timeSpent int) {
    c.TotalReviews++
    if correct {
        c.SuccessCount++
    }
    
    c.LastMCQChoice = selected
    if !correct && !contains(c.IncorrectChoices, selected) {
        c.IncorrectChoices = append(c.IncorrectChoices, selected)
    }
    
    // Update MCQ accuracy
    c.MCQAccuracy = float64(c.SuccessCount) / float64(c.TotalReviews)
}
```

### Phase 6: UI/UX Design

#### 6.1 MCQ Display Format
```
┌─────────────────────────────────────────────────────┐
│ Question 3/10 [MCQ]                                 │
├─────────────────────────────────────────────────────┤
│                                                     │
│ What type of failures can primary-backup           │
│ replication typically handle?                      │
│                                                     │
│   [ ] A) Software bugs that cause incorrect        │
│          calculations                               │
│                                                     │
│   [•] B) Hardware failures that cause the server   │
│          to stop executing                          │
│                                                     │
├─────────────────────────────────────────────────────┤
│ [↑↓] Navigate  [Enter] Submit  [E] Explanation     │
└─────────────────────────────────────────────────────┘

After submission:
├─────────────────────────────────────────────────────┤
│ ✓ Correct!                                          │
│                                                     │
│ Explanation:                                        │
│ Primary-backup replication handles fail-stop       │
│ failures where servers stop cleanly, not bugs      │
│ that would affect both replicas identically.       │
├─────────────────────────────────────────────────────┤
│ [N] Next  [E] Toggle Explanation  [Q] Quit         │
└─────────────────────────────────────────────────────┘
```

#### 6.2 Mixed Quiz Flow
- Support both MCQ and subjective in same session
- Clear visual indicators for question type
- Different interaction patterns (selection vs typing)

### Phase 7: Session Management

#### 7.1 Update Session Format
```go
// internal/tui/session/session.go

type Answer struct {
    QuestionID   string    `json:"question_id"`
    Type         string    `json:"type"`        // "mcq" or "subjective"
    Content      string    `json:"content"`     // Text answer or selected option
    IsCorrect    *bool     `json:"is_correct,omitempty"` // For MCQ
    TimeSpent    int       `json:"time_spent"`
    Timestamp    time.Time `json:"timestamp"`
}
```

### Phase 8: Testing Strategy

#### 8.1 Unit Tests
- MCQ parser tests
- MCQ component state transitions
- SRS integration with MCQ metrics

#### 8.2 Integration Tests
- Mixed quiz sessions (MCQ + subjective)
- Session persistence with MCQ answers
- SRS scheduling with MCQ performance

### Phase 9: Configuration

#### 9.1 Config Updates
```go
// internal/config/config.go

type TUIConfig struct {
    // ... existing fields ...
    
    // MCQ settings
    EnableMCQ           bool   `yaml:"enable_mcq"`
    MCQFeedbackDelay    int    `yaml:"mcq_feedback_delay"`    // ms
    ShowExplanations    bool   `yaml:"show_explanations"`
    MCQSessionRatio     float64 `yaml:"mcq_session_ratio"`     // % of MCQs in mixed sessions
}
```

### Implementation Order

1. **Week 1: Foundation**
   - Extend Question model (1.1)
   - Update markdown parser (2.1)
   - Basic MCQ component (3.1)

2. **Week 2: Integration**
   - App model updates (3.2)
   - State machine modifications (4.1)
   - Session management (7.1)

3. **Week 3: SRS & Polish**
   - SRS integration (5.1, 5.2)
   - UI/UX refinements (6.1, 6.2)
   - Testing and configuration (8.1, 8.2, 9.1)

### Migration Path

1. **Backward Compatibility**
   - Default `Type: "subjective"` for existing questions
   - Graceful handling of missing MCQ fields
   - Session format versioning

2. **Content Migration**
   - Script to validate existing MCQ markdown files
   - Batch import of MCQ questions into SRS
   - Topic mapping for MCQ questions

### Success Metrics

- MCQ questions load and display correctly
- Users can select and submit answers
- Feedback (correct/incorrect + explanation) displays
- SRS tracks MCQ performance separately
- Sessions persist MCQ progress
- Mixed sessions work seamlessly

### Risk Mitigation

1. **Parser Complexity**: Start with simple MCQ format, extend later
2. **UI Complexity**: Use existing lipgloss styles, minimal custom rendering
3. **State Management**: Reuse existing session/state patterns
4. **Performance**: Lazy load MCQ questions, cache parsed results

### Future Enhancements

1. **Advanced MCQ Types**
   - Multiple correct answers
   - Matching questions
   - Fill-in-the-blank

2. **Analytics**
   - MCQ-specific performance metrics
   - Common wrong answer patterns
   - Time-to-answer analysis

3. **Adaptive Learning**
   - Adjust MCQ difficulty based on performance
   - Mix MCQ/subjective based on topic mastery
   - Personalized review schedules

## Summary

This plan adds comprehensive MCQ support to the quiz-tui while:
- Leveraging existing SRS infrastructure
- Maintaining backward compatibility
- Providing clear UI/UX for different question types
- Integrating seamlessly with session management
- Supporting mixed quiz sessions

---

## Implementation Status

### Completed (2025-10-11)

#### Phase 1: Data Models ✅
- Extended `Question` model with MCQ fields (Type, Options, Answer, Explanation, Hook)
- Updated tags from `yaml` to `toml` for frontmatter
- Extended `Card` model with MCQ performance tracking (MCQAccuracy, LastMCQChoice, IncorrectChoices)
- Updated `Response` model to support MCQ data (IsCorrect, SelectedOption, QuestionType)

#### Phase 2: Parser Extensions ✅
- Updated markdown parser to support both TOML (+++...+++) and YAML (---...---) frontmatter
- Added MCQ section parsing (options, answer, explanation, hook)
- Implemented `parseMCQOptions()` function to parse MCQ options
- Maintained backward compatibility with existing YAML format

#### Phase 3: TUI Components ✅
- Created `internal/tui/components/mcq.go` component
- Implemented navigation (up/down arrows)
- Implemented selection and submission
- Added visual feedback (correct/incorrect with colors)
- Added explanation display toggle

#### Phase 4: State Management ✅
- Updated `ImprovedAppModel` with MCQ-specific fields
- Implemented `initializeQuestionComponent()` to detect and initialize correct component
- Modified `StateQuestion` handler to support both MCQ and subjective inputs
- Added `saveMCQAnswerCmd()` for MCQ-specific answer saving
- Updated `nextQuestionMsg` handler to reinitialize components

#### Phase 5: SRS Integration ✅
- Added `RecordMCQAttempt()` method to Card for MCQ-specific tracking
- Extended session persistence to include MCQ metadata

#### Phase 6: UI/UX ✅
- Updated `renderQuestion()` to display MCQ vs subjective differently
- Added context-aware help text
- Implemented visual indicators for question type

#### Testing ✅
- Fixed duplicate function issues
- Updated parser tests for new signature
- All markdown parser tests passing
- Build successful

### Sample MCQ File Created
Created test MCQ at: `ddia-quiz-bot/content/chapters/11-mit-6824-fault-tolerance/mcq/test-mcq-001.md`

### Pending Items

- [ ] Phase 8.1: Write comprehensive unit tests for MCQ components
- [ ] Phase 9.1: Add MCQ configuration options to TUIConfig
- [ ] Integration testing with real quiz sessions
- [ ] Performance testing with mixed MCQ/subjective sessions
- [ ] SRS scheduler MCQ-specific review logic
- [ ] MCQ analytics and reporting
