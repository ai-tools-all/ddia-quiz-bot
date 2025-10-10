# Feature: Interactive Quiz TUI with Bubbletea

## Date: 2025-10-10
## Category: Feature
## Status: Planning

## Overview
Design and implement an interactive Terminal User Interface (TUI) for taking quizzes that reads questions from markdown files, collects user responses interactively, stores them in CSV format, and optionally evaluates them using the existing quiz-eval CLI tool.

## Objectives

### Primary Goals
1. **Interactive Question Presentation**: Display quiz questions one-by-one in a beautiful, readable format
2. **User Response Collection**: Capture typed answers with a smooth input experience
3. **Progress Tracking**: Show quiz progress and allow navigation between questions
4. **CSV Export**: Save responses in the same format expected by quiz-eval
5. **Markdown Rendering**: Render questions beautifully using Glow
6. **Seamless Integration**: Optionally chain to quiz-eval for immediate evaluation

### User Experience Goals
- Intuitive keyboard navigation
- Beautiful, styled terminal interface
- Responsive and performant
- Progress persistence (resume interrupted quizzes)
- Review mode before submission

## What's Already Implemented (Reusable)

### Existing Components in quiz-evaluator

✅ **internal/models/**
- `Question` struct with ID, Title, MainQuestion, Level, Category
- `QuestionIndex` for fast lookups
- `UserResponse` struct (needs enhancement)
- JSON tags already present

✅ **internal/markdown/**
- `Scanner.ScanQuestions()` - recursively finds all markdown files
- `Parser.ParseQuestionFile()` - extracts YAML frontmatter and sections
- Supports section parsing (Core Concepts, Rubrics, etc.)

✅ **internal/config/**
- Complete configuration system with Viper
- Environment variable binding
- AI provider configs (for evaluation)

⚠️ **Needs Enhancement:**
1. Parser to handle MCQ format (options, answer, explanation)
2. Question model to include Options, CorrectAnswer, Explanation fields
3. New JSON writer module (currently only has CSV)

### Question Types Found in ddia-quiz-bot/content

1. **MCQ Questions** (~171 files):
   - `## question` - the question text
   - `## options` - A) B) C) D) choices
   - `## answer` - correct answer (e.g., "A")
   - `## explanation` - why the answer is correct

2. **Subjective Questions** (sample files):
   - `## Question` - open-ended question
   - `## Core Concepts` - required concepts
   - `## Sample Excellent Answer` - reference answer
   - `## Evaluation Rubric` - grading criteria

## Architecture Overview

### High-Level Components

```
┌─────────────────────────────────────────────────────────┐
│           Interactive Quiz TUI (quiz-tui)               │
│                                                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │   REUSE:     │  │   Response   │  │    JSON      │ │
│  │   markdown   │→ │  Collector   │→ │   Writer     │ │
│  │   Scanner    │  │              │  │   (NEW)      │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
│         ↓                                               │
│  ┌──────────────┐  ┌──────────────────────────────┐   │
│  │  Bubbletea   │  │    Charm Components          │   │
│  │     TUI      │  │  (Huh, Lipgloss, Glamour)   │   │
│  └──────────────┘  └──────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
                    ↓
        ┌──────────────────────┐
        │   quiz-eval CLI      │
        │ (Optional Evaluation)│
        └──────────────────────┘
```

### Tech Stack

#### Charm Ecosystem Libraries
1. **Bubbletea** (`github.com/charmbracelet/bubbletea`)
   - Core TUI framework using The Elm Architecture
   - State management and event handling
   - Main application loop

2. **Huh** (`github.com/charmbracelet/huh`)
   - Terminal forms and prompts
   - Input fields for responses
   - Multi-line text areas for long answers
   - Confirmation dialogs

3. **Lipgloss** (`github.com/charmbracelet/lipgloss`)
   - Styling and layout
   - Color schemes and themes
   - Borders, padding, margins
   - Responsive layouts

4. **Glow** (`github.com/charmbracelet/glow`)
   - Markdown rendering
   - Syntax highlighting
   - Beautiful question display
   - Code block rendering

5. **Bubbles** (`github.com/charmbracelet/bubbles`)
   - Pre-built UI components
   - Progress bars
   - Spinners
   - Viewports for scrolling
   - Text input components

6. **Glamour** (`github.com/charmbracelet/glamour`)
   - Markdown to ANSI rendering (used by Glow)
   - Custom styles for question rendering
   - Code highlighting

#### Additional Libraries
- **Cobra** (existing) - CLI command structure
- **Viper** (existing) - Configuration management
- **Existing internal packages** - Question parsing, CSV I/O

## Detailed Component Design

### 1. Enhanced Question Model (ENHANCE EXISTING)

**Location**: `internal/models/question.go`

**Add to existing Question struct**:
```go
type Question struct {
    // Existing fields (keep as is)
    ID                 string            `yaml:"id" json:"id"`
    Title              string            `yaml:"title" json:"title"`
    MainQuestion       string            `yaml:"question" json:"question"`
    Level              string            `yaml:"level" json:"level"`
    Category           string            `yaml:"category" json:"category"`
    
    // NEW: MCQ support
    QuestionType       string            `yaml:"-" json:"question_type"` // "mcq" or "subjective"
    Options            []string          `yaml:"options" json:"options,omitempty"`
    CorrectAnswer      string            `yaml:"answer" json:"correct_answer,omitempty"`
    Explanation        string            `yaml:"explanation" json:"explanation,omitempty"`
    
    // Existing subjective fields
    CoreConcepts       []string          `yaml:"core_concepts" json:"core_concepts,omitempty"`
    PeripheralConcepts []string          `yaml:"peripheral_concepts" json:"peripheral_concepts,omitempty"`
    SampleExcellent    string            `yaml:"sample_excellent" json:"sample_excellent,omitempty"`
    SampleAcceptable   string            `yaml:"sample_acceptable" json:"sample_acceptable,omitempty"`
    EvaluationRubric   map[string]string `yaml:"evaluation_rubric" json:"evaluation_rubric,omitempty"`
    
    // Metadata
    Day                int               `yaml:"day" json:"day,omitempty"`
    Tags               []string          `yaml:"tags" json:"tags,omitempty"`
    FilePath           string            `yaml:"-" json:"-"`
}

// Helper method
func (q *Question) IsMCQ() bool {
    return len(q.Options) > 0
}
```

### 2. Enhanced Parser (ENHANCE EXISTING)

**Location**: `internal/markdown/parser.go`

**Add to Parser.parseBody()** to handle MCQ sections:
```go
// Add cases in saveSection():
case "options":
    question.Options = p.parseListSection(content)
    question.QuestionType = "mcq"
case "answer":
    question.CorrectAnswer = strings.TrimSpace(content)
case "explanation":
    question.Explanation = content
```

**QuestionType auto-detection**:
- If `Options` field is populated → "mcq"
- Otherwise → "subjective"

### 3. Question Loader Module (NEW - wraps existing Scanner)

**Location**: `internal/quiz/loader.go`

**Responsibilities**:
- Wrapper around existing markdown.Scanner
- Filtering and selection logic
- Randomization and shuffling
- Question type filtering

**Interface**:
```go
type QuizLoader struct {
    scanner    *markdown.Scanner
    allQuestions models.QuestionIndex
}

func NewQuizLoader(questionsDir string) (*QuizLoader, error) {
    scanner := markdown.NewScanner(questionsDir)
    index, err := scanner.ScanQuestions()
    if err != nil {
        return nil, err
    }
    return &QuizLoader{
        scanner: scanner,
        allQuestions: index,
    }, nil
}

func (l *QuizLoader) FilterByOptions(opts LoadOptions) []*models.Question
func (l *QuizLoader) Shuffle(questions []*models.Question, seed int64) []*models.Question
func (l *QuizLoader) GetQuestionCount() int

type LoadOptions struct {
    Chapters  []string
    Levels    []string
    Tags      []string
    Type      string  // "mcq", "subjective", "all"
    Limit     int
    Random    bool
}
```

**Key Features**:
- Reuses existing Scanner and Parser
- Adds filtering logic on top
- Supports both MCQ and subjective questions

### 2. TUI Application (Bubbletea)

**Location**: `internal/tui/app.go`

**State Management**:
```go
type Model struct {
    // Quiz state
    questions       []Question
    currentIndex    int
    responses       map[string]Response
    startTime       time.Time
    
    // UI state
    state           AppState
    width           int
    height          int
    
    // Components
    viewport        viewport.Model
    textInput       textarea.Model
    progressBar     progress.Model
    
    // Styling
    styles          Styles
    
    // Configuration
    config          QuizConfig
}

type AppState int
const (
    StateWelcome AppState = iota
    StateQuizSetup
    StateQuestion
    StateReview
    StateSubmit
    StateComplete
)
```

**The Elm Architecture Implementation**:
```go
// Init: Initialize the application
func (m Model) Init() tea.Cmd

// Update: Handle messages and update state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)

// View: Render the UI
func (m Model) View() string
```

### 3. UI Screens/States

#### A. Welcome Screen
**Purpose**: Introduction and setup

**Features**:
- Project title and branding
- Quick instructions
- Configuration options (chapters, levels, count)
- Start button

**Implementation**: Use `huh.Form` for setup options

**Layout**:
```
╭──────────────────────────────────────────╮
│     📚 DDIA Quiz Interactive Session     │
├──────────────────────────────────────────┤
│                                          │
│  Welcome! Ready to test your knowledge?  │
│                                          │
│  Select chapters to include:             │
│  [x] Chapter 3: Storage and Retrieval    │
│  [x] Chapter 5: Replication              │
│  [ ] Chapter 7: Transactions             │
│                                          │
│  Difficulty level: [All / L3 / L5 / L7] │
│  Question count: [10]                    │
│  Randomize order: [Yes]                  │
│                                          │
│       [Start Quiz] [Quit]                │
╰──────────────────────────────────────────╯
```

#### B. Question Screen
**Purpose**: Display question and collect answer

**Features**:
- Rendered markdown question (Glow/Glamour)
- Response input (Huh textarea for subjective, select for MCQ)
- Navigation buttons
- Progress indicator
- Timer (optional)
- Keyboard shortcuts

**Layout**:
```
╭──────────────────────────────────────────────────╮
│ Question 3/10  [████████░░░░░░░░░░]  30%  ⏱ 5min │
├──────────────────────────────────────────────────┤
│ Chapter 7: Transactions | Level: L5              │
│                                                  │
│ Read Committed Isolation                         │
│                                                  │
│ Under read committed isolation level, which      │
│ phenomenon is prevented?                         │
│                                                  │
│ Options:                                         │
│ • A) Dirty reads and dirty writes                │
│ • B) Lost updates                                │
│ • C) Write skew                                  │
│ • D) Phantom reads                               │
│                                                  │
├──────────────────────────────────────────────────┤
│ Your Answer:                                     │
│ ┌──────────────────────────────────────────────┐ │
│ │ A                                          █ │ │
│ └──────────────────────────────────────────────┘ │
├──────────────────────────────────────────────────┤
│ [p] Previous [n] Next [r] Review [s] Submit      │
╰──────────────────────────────────────────────────╯
```

**For Subjective Questions**:
```
╭──────────────────────────────────────────────────╮
│ Question 5/10  [████████████░░░░░░░]  50%  ⏱ 8min│
├──────────────────────────────────────────────────┤
│ Chapter 9: Distributed Systems | Level: L7       │
│                                                  │
│ Raft Consensus Algorithm                         │
│                                                  │
│ Explain how the Raft consensus algorithm handles│
│ leader election and log replication. Discuss the│
│ safety properties that Raft guarantees.          │
│                                                  │
├──────────────────────────────────────────────────┤
│ Your Answer: (Ctrl+S to save, ESC to cancel)    │
│ ┌──────────────────────────────────────────────┐ │
│ │ Raft uses randomized timeouts for leader    █│ │
│ │ election. When a follower doesn't hear from   │ │
│ │ the leader, it increments its term and        │ │
│ │ requests votes...                             │ │
│ │                                               │ │
│ │ [Lines: 4/unlimited | Words: 28]              │ │
│ └──────────────────────────────────────────────┘ │
├──────────────────────────────────────────────────┤
│ [p] Previous [n] Next [r] Review [s] Submit      │
╰──────────────────────────────────────────────────╯
```

#### C. Review Screen
**Purpose**: Review all answers before submission

**Features**:
- List of all questions
- Answer status (answered/skipped)
- Jump to specific question
- Answer preview

**Layout**:
```
╭──────────────────────────────────────────────────╮
│              📋 Review Your Answers              │
├──────────────────────────────────────────────────┤
│                                                  │
│  ✓ Q1: Read Committed Isolation         [Edit]  │
│     Answer: A                                    │
│                                                  │
│  ✓ Q2: ACID Properties                  [Edit]  │
│     Answer: Atomicity ensures...                │
│                                                  │
│  ⚠ Q3: Write Skew                       [Edit]  │
│     No answer provided                           │
│                                                  │
│  ✓ Q4: Snapshot Isolation               [Edit]  │
│     Answer: C                                    │
│                                                  │
│  Progress: 9/10 answered (1 skipped)             │
│                                                  │
│  [Continue Editing] [Submit Anyway]              │
╰──────────────────────────────────────────────────╯
```

#### D. Submit Confirmation
**Purpose**: Final confirmation before saving

**Features**:
- Summary statistics
- Warning for unanswered questions
- Confirm/Cancel options
- Option to evaluate immediately

**Layout**:
```
╭──────────────────────────────────────────────────╮
│              ✅ Ready to Submit?                 │
├──────────────────────────────────────────────────┤
│                                                  │
│  Quiz Summary:                                   │
│  • Questions attempted: 9/10                     │
│  • Time taken: 15 minutes                        │
│  • Chapters covered: 3, 5, 7                     │
│                                                  │
│  Your responses will be saved to:               │
│  responses_abhishek_2025-10-10.csv              │
│                                                  │
│  [ ] Evaluate responses immediately with AI      │
│                                                  │
│  [Submit] [Review Again] [Cancel]                │
╰──────────────────────────────────────────────────╯
```

#### E. Completion Screen
**Purpose**: Show completion status and next steps

**Features**:
- Success message
- File location
- Evaluation option
- Statistics

**Layout**:
```
╭──────────────────────────────────────────────────╮
│              🎉 Quiz Complete!                   │
├──────────────────────────────────────────────────┤
│                                                  │
│  Your responses have been saved successfully!    │
│                                                  │
│  📁 File: responses_abhishek_2025-10-10.csv      │
│  📊 Questions: 10                                │
│  ⏱  Time: 15 minutes 32 seconds                 │
│                                                  │
│  What would you like to do next?                │
│                                                  │
│  [1] Evaluate responses with AI                  │
│  [2] Take another quiz                           │
│  [3] View response file                          │
│  [4] Exit                                        │
│                                                  │
╰──────────────────────────────────────────────────╯
```

### 4. Response Collector

**Location**: `internal/quiz/collector.go`

**Responsibilities**:
- Track user responses
- Handle answer validation
- Support answer editing
- Auto-save functionality
- Resume capability

**Interface**:
```go
type ResponseCollector interface {
    SetResponse(questionID string, answer string) error
    GetResponse(questionID string) (string, bool)
    DeleteResponse(questionID string) error
    GetAllResponses() map[string]Response
    GetProgress() Progress
    IsComplete() bool
    SaveDraft() error
    LoadDraft() error
}

type Response struct {
    QuestionID  string
    Question    string
    UserAnswer  string
    Timestamp   time.Time
    TimeTaken   time.Duration
    IsAnswered  bool
}

type Progress struct {
    TotalQuestions int
    Answered       int
    Skipped        int
    CurrentIndex   int
    TimeElapsed    time.Duration
}
```

### 5. JSON Response Writer (NEW)

**Location**: `internal/json/writer.go`

**Responsibilities**:
- Save quiz session and responses to JSON
- Include metadata (user, timestamp, session info)
- Support for both MCQ and subjective responses

**JSON Format**:
```json
{
  "session": {
    "session_id": "quiz-20251010-233440",
    "user": "abhishek",
    "started_at": "2025-10-10T23:34:40Z",
    "completed_at": "2025-10-10T23:48:15Z",
    "duration_seconds": 815,
    "total_questions": 10,
    "answered": 9,
    "skipped": 1
  },
  "responses": [
    {
      "question_id": "ch07-read-committed",
      "question_title": "Read Committed Isolation",
      "question_type": "mcq",
      "level": "L3",
      "chapter": "Chapter 7",
      "user_response": "A",
      "timestamp": "2025-10-10T23:35:12Z",
      "time_taken_seconds": 32
    },
    {
      "question_id": "ch09-raft-consensus",
      "question_title": "Raft Consensus",
      "question_type": "subjective",
      "level": "L7",
      "chapter": "Chapter 9",
      "user_response": "Raft uses randomized timeouts for leader election...",
      "timestamp": "2025-10-10T23:38:45Z",
      "time_taken_seconds": 213
    }
  ]
}
```

**Interface**:
```go
type JSONWriter struct {
    outputPath string
}

func (w *JSONWriter) WriteSession(session *QuizSession) error
func (w *JSONWriter) LoadSession(sessionID string) (*QuizSession, error)

type QuizSession struct {
    SessionID        string            `json:"session_id"`
    User             string            `json:"user"`
    StartedAt        time.Time         `json:"started_at"`
    CompletedAt      time.Time         `json:"completed_at,omitempty"`
    DurationSeconds  int               `json:"duration_seconds,omitempty"`
    TotalQuestions   int               `json:"total_questions"`
    Answered         int               `json:"answered"`
    Skipped          int               `json:"skipped"`
    Responses        []ResponseRecord  `json:"responses"`
}

type ResponseRecord struct {
    QuestionID        string    `json:"question_id"`
    QuestionTitle     string    `json:"question_title"`
    QuestionType      string    `json:"question_type"`
    Level             string    `json:"level"`
    Chapter           string    `json:"chapter"`
    UserResponse      string    `json:"user_response"`
    Timestamp         time.Time `json:"timestamp"`
    TimeTakenSeconds  int       `json:"time_taken_seconds"`
}
```

### 6. Styling System

**Location**: `internal/tui/styles.go`

**Color Schemes**:
```go
type Styles struct {
    // Colors
    Primary     lipgloss.Color
    Secondary   lipgloss.Color
    Success     lipgloss.Color
    Warning     lipgloss.Color
    Error       lipgloss.Color
    Muted       lipgloss.Color
    
    // Styles
    TitleStyle      lipgloss.Style
    HeaderStyle     lipgloss.Style
    QuestionStyle   lipgloss.Style
    InputStyle      lipgloss.Style
    ButtonStyle     lipgloss.Style
    ProgressStyle   lipgloss.Style
    HelpStyle       lipgloss.Style
    
    // Layout
    ContentWidth    int
    BorderType      lipgloss.Border
}
```

**Theme Options**:
- Light mode
- Dark mode
- High contrast
- Custom (from config)

## Implementation Plan

### Summary: Building on Existing Foundation

**What We're Reusing** (✅ ~60% of infrastructure exists):
- Complete markdown parsing and scanning system
- Question model with metadata support
- Configuration management with Viper
- AI evaluation system (for later integration)
- CSV reader/writer (for evaluate command compatibility)

**What We're Building** (🆕 ~40% new code):
- Bubbletea TUI application and screens
- Interactive response collection
- JSON session storage
- Quiz session management
- Charm ecosystem integration

**What We're Enhancing** (⚠️ minor additions):
- Question model: Add MCQ fields (Options, CorrectAnswer, Explanation)
- Parser: Handle MCQ sections
- Config: Add TUI-specific settings

**Estimated LOC**: ~2,500 new lines of Go code

### Phase 1: Foundation & Model Enhancement (Days 1-2)
**Goal**: Enhance existing models and set up basic TUI structure

**Tasks**:
1. Enhance Question model (internal/models/question.go)
   - Add MCQ fields: Options, CorrectAnswer, Explanation, QuestionType
   - Add Day and Tags fields for metadata
   - Add IsMCQ() helper method
   - Update JSON tags

2. Enhance Parser (internal/markdown/parser.go)
   - Add parsing for "options", "answer", "explanation" sections
   - Auto-detect question type based on presence of options
   - Parse tags from frontmatter

3. Create QuizLoader wrapper (internal/quiz/loader.go)
   - Wrap existing Scanner
   - Add filtering by chapters, levels, tags, type
   - Add shuffling functionality
   - Question selection logic

4. Set up Bubbletea application structure (internal/tui/)
   - Create main model and state machine
   - Implement Init, Update, View functions
   - Set up message handling
   - Define AppState enum

5. Install Charm dependencies
   - Add bubbletea, bubbles, lipgloss, huh, glamour to go.mod
   - Test basic TUI render

**Deliverables**:
- Enhanced models supporting MCQ questions
- Parser that handles both question types
- Basic Bubbletea app that compiles
- Dependency setup complete

### Phase 2: Question Display & Input (Week 1-2)
**Goal**: Rich question rendering and response collection

**Tasks**:
1. Integrate Glow/Glamour for markdown rendering
   - Render question text beautifully
   - Handle code blocks and formatting
   - Syntax highlighting

2. Implement response input components
   - Multiple choice selector (Huh select)
   - Multi-line text area (Huh textarea)
   - Input validation
   - Character/word count

3. Progress tracking UI
   - Progress bar with Bubbles
   - Question counter
   - Timer display (optional)
   - Status indicators

4. Keyboard shortcuts
   - Navigation (n, p, g)
   - Quick actions (s, r, q)
   - Help overlay (?)

**Deliverables**:
- Beautiful question display
- Functional input collection
- Progress tracking

### Phase 3: Response Management (Week 2)
**Goal**: Save, review, and edit responses

**Tasks**:
1. Response collector implementation
   - In-memory storage
   - Answer editing
   - Validation

2. Review screen
   - List all questions
   - Show answer status
   - Jump to question
   - Answer preview

3. Draft/Auto-save functionality
   - Periodic auto-save
   - Resume interrupted sessions
   - Draft storage location

4. CSV export
   - Format responses
   - Include metadata
   - File naming (user, timestamp)

**Deliverables**:
- Working response collection
- Review functionality
- CSV export

### Phase 4: Integration & Polish (Week 2-3)
**Goal**: Integration with quiz-eval and UX refinements

**Tasks**:
1. Integration with quiz-eval
   - Chain to evaluation command
   - Pass CSV file
   - Display evaluation results

2. Configuration system
   - Default settings
   - User preferences
   - Theme selection

3. Error handling
   - Graceful failures
   - User-friendly messages
   - Recovery options

4. Help system
   - Keyboard shortcuts help
   - Inline tips
   - Documentation

**Deliverables**:
- End-to-end workflow
- Polished UX
- Documentation

### Phase 5: Testing & Optimization (Week 3)
**Goal**: Ensure reliability and performance

**Tasks**:
1. Unit tests
   - Question loader
   - Response collector
   - CSV writer

2. Integration tests
   - Full quiz workflow
   - CSV format validation
   - State management

3. Performance optimization
   - Efficient rendering
   - Memory management
   - Fast markdown parsing

4. User acceptance testing
   - Real quiz sessions
   - Feedback collection
   - Bug fixes

**Deliverables**:
- Test coverage
- Performance benchmarks
- Bug-free release

## CLI Command Structure

### New Command: `quiz-tui`

```bash
# Main command - interactive mode
quiz-tui start [questions-dir]

# With options
quiz-tui start ./content --chapters 3,5,7 --level L5 --count 10

# Resume saved session
quiz-tui resume ./drafts/quiz-session-123.json

# Non-interactive mode (same as quiz-eval)
quiz-tui evaluate input.csv questions-dir
```

### Flags

```
Global Flags:
  --config string         Config file (default: ./config.yaml)
  -v, --verbose          Verbose output
  --user string          User identifier (default: $USER)

Start Command Flags:
  --chapters strings     Chapters to include (e.g., 3,5,7)
  --levels strings       Difficulty levels (e.g., L3,L5,L7)
  --tags strings         Filter by tags
  --count int           Number of questions (default: 10)
  --random              Randomize question order
  --timer               Show timer
  --auto-save int       Auto-save interval in seconds (default: 60)
  --output string       Output CSV file
  --evaluate            Evaluate immediately after completion

Resume Command Flags:
  --session string      Session file to resume

Evaluate Command Flags:
  (Same as existing quiz-eval flags)
```

## Configuration

**File**: `config.yaml`

```yaml
# TUI Configuration
tui:
  theme: dark           # light, dark, high-contrast
  show_timer: true
  auto_save_interval: 60  # seconds
  enable_markdown: true
  
  # Colors (optional override)
  colors:
    primary: "#7C3AED"
    success: "#10B981"
    warning: "#F59E0B"
    error: "#EF4444"

# Quiz Defaults
quiz:
  default_chapters: [3, 5, 7]
  default_level: "all"  # all, L3, L5, L7
  default_count: 10
  randomize: true
  
  # Directories
  questions_dir: "./content/chapters"
  drafts_dir: "./drafts"
  output_dir: "./responses"

# User Settings
user:
  name: "abhishek"
  auto_identify: true  # Use $USER if not specified

# Integration
integration:
  auto_evaluate: false
  eval_command: "./quiz-eval"
  eval_flags: "--workers 5 --verbose"

# Existing AI configuration (for evaluation)
ai:
  provider: "openai"
  api_key: "${OPENAI_API_KEY}"
  model: "gpt-4-turbo-preview"
```

## File Structure

```
quiz-evaluator/
├── cmd/
│   ├── root.go                    # ✅ EXISTING (reuse)
│   ├── evaluate.go                # ✅ EXISTING (reuse)
│   └── start.go                   # 🆕 NEW: Interactive quiz command
│
├── internal/
│   ├── tui/                       # 🆕 NEW: TUI components
│   │   ├── app.go                # Main Bubbletea app
│   │   ├── models.go             # TUI models
│   │   ├── states.go             # State definitions
│   │   ├── styles.go             # Lipgloss styles
│   │   ├── screens/              # Screen implementations
│   │   │   ├── welcome.go
│   │   │   ├── question.go
│   │   │   ├── review.go
│   │   │   ├── submit.go
│   │   │   └── complete.go
│   │   ├── components/           # Reusable components
│   │   │   ├── progress.go
│   │   │   ├── question_viewer.go
│   │   │   ├── answer_input.go
│   │   │   └── help.go
│   │   └── renderer/             # Markdown rendering
│   │       └── question_renderer.go
│   │
│   ├── quiz/                      # 🆕 NEW: Quiz logic
│   │   ├── loader.go             # Wraps markdown.Scanner
│   │   ├── collector.go          # Response collector
│   │   ├── session.go            # Session management
│   │   └── filter.go             # Question filtering
│   │
│   ├── json/                      # 🆕 NEW: JSON I/O
│   │   ├── writer.go             # JSON session writer
│   │   └── reader.go             # JSON session reader
│   │
│   ├── csv/                       # ✅ EXISTING (reuse)
│   │   ├── reader.go             # Used by evaluate command
│   │   └── writer.go             # Used by evaluate command
│   │
│   ├── markdown/                  # ⚠️ ENHANCE EXISTING
│   │   ├── parser.go             # Add MCQ parsing (options, answer, explanation)
│   │   └── scanner.go            # Reuse as-is
│   │
│   ├── models/                    # ⚠️ ENHANCE EXISTING
│   │   ├── question.go           # Add Options, CorrectAnswer, Explanation, QuestionType
│   │   ├── response.go           # Reuse as-is
│   │   └── evaluation.go         # Reuse as-is
│   │
│   ├── evaluator/                 # ✅ EXISTING (reuse for evaluation)
│   │   ├── ai_client.go
│   │   ├── openai_client.go
│   │   ├── prompt.go
│   │   └── evaluator.go
│   │
│   └── config/                    # ⚠️ ENHANCE EXISTING
│       └── config.go              # Add TUI config section
│
├── main.go                        # ✅ EXISTING (reuse)
├── go.mod                         # ⚠️ ENHANCE (add Charm deps)
└── go.sum                         # ⚠️ ENHANCE (auto-updated)
```

### Legend:
- ✅ **EXISTING (reuse)**: Use as-is, no changes needed
- ⚠️ **ENHANCE EXISTING**: Extend with new functionality
- 🆕 **NEW**: Create from scratch

## User Experience Flow

### Happy Path

```
1. User runs: quiz-tui start ./content --user abhishek

2. Welcome screen appears with:
   - Greeting
   - Chapter selection
   - Level selection
   - Question count

3. User configures quiz:
   - Selects chapters: 3, 5, 7
   - Selects level: All
   - Sets count: 10
   - Enables randomize

4. Quiz starts:
   - Question 1/10 displays
   - Beautiful markdown rendering
   - User types answer

5. Navigation:
   - User presses 'n' for next
   - Answer auto-saved
   - Question 2/10 loads

6. Continues answering questions

7. User presses 'r' to review:
   - Sees all 10 questions
   - 8 answered, 2 skipped
   - Reviews answer to Q5
   - Edits answer

8. User presses 's' to submit:
   - Confirmation dialog
   - Warning about 2 skipped
   - Option to evaluate
   - User confirms

9. Completion screen:
   - Success message
   - JSON file location: `responses_abhishek_2025-10-10.json`
   - Option to evaluate
   - User selects evaluate

10. Evaluation runs:
    - Converts JSON to CSV format for evaluator
    - Progress shown
    - Results displayed
    - Output file shown

11. User exits
```

### Error Scenarios

1. **No questions found**
   - Show friendly error
   - Suggest checking directory
   - Exit gracefully

2. **Invalid configuration**
   - Highlight errors
   - Suggest fixes
   - Allow retry

3. **Session interrupted**
   - Draft auto-saved
   - Resume option on restart
   - No data loss

4. **Evaluation fails**
   - CSV still saved
   - Error message shown
   - Option to retry

## Testing Strategy

### Unit Tests
- Question loader filtering
- Response collector state
- CSV formatting
- Markdown parsing
- Style rendering

### Integration Tests
- Full quiz workflow
- State transitions
- CSV export/import
- Configuration loading

### Manual Testing
- Multiple quiz sessions
- Different terminal sizes
- Various color schemes
- Keyboard navigation
- Edge cases (0 questions, huge responses)

## Performance Considerations

### Optimization Areas
1. **Markdown Rendering**: Cache rendered questions
2. **State Updates**: Minimize re-renders
3. **CSV Writing**: Buffered I/O
4. **Memory**: Stream large question sets

### Targets
- Render time: < 16ms (60fps)
- Response time: < 100ms
- Memory: < 50MB for 1000 questions
- Startup: < 1 second

## Dependencies

### New Dependencies
```go
require (
    github.com/charmbracelet/bubbletea v0.27.1
    github.com/charmbracelet/bubbles v0.20.0
    github.com/charmbracelet/lipgloss v0.13.1
    github.com/charmbracelet/huh v0.6.0
    github.com/charmbracelet/glamour v0.8.0
)
```

### Version Compatibility
- Go 1.21+
- Terminal with ANSI color support
- UTF-8 encoding

## Security Considerations

1. **Input Validation**: Sanitize all user inputs
2. **Path Traversal**: Validate file paths
3. **Resource Limits**: Cap response length
4. **API Keys**: Never log or display
5. **File Permissions**: Secure draft files

## Accessibility

1. **Keyboard Only**: Full functionality without mouse
2. **Screen Readers**: Semantic markup where possible
3. **High Contrast**: Support high-contrast themes
4. **Clear Labels**: Descriptive text for all actions

## Documentation

### User Documentation
1. README with screenshots
2. Command reference
3. Configuration guide
4. Keyboard shortcuts
5. Troubleshooting

### Developer Documentation
1. Architecture overview
2. Component API docs
3. State machine diagram
4. Contributing guide

## Success Metrics

### Functionality
- ✅ Users can take quizzes interactively
- ✅ Responses saved in correct CSV format
- ✅ Integration with quiz-eval works
- ✅ Sessions can be resumed

### Quality
- ✅ No crashes during normal use
- ✅ Responsive on various terminals
- ✅ Beautiful, readable display
- ✅ Intuitive navigation

### Performance
- ✅ Handles 100+ questions smoothly
- ✅ Fast response to keypresses
- ✅ Low memory footprint

## Future Enhancements

### Phase 2 Features
1. **Spaced Repetition**: Track quiz history, suggest retakes
2. **Statistics Dashboard**: Show progress over time
3. **Collaborative Mode**: Multiple users, leaderboards
4. **Export Formats**: JSON, PDF, HTML
5. **Question Editor**: Create questions in TUI
6. **Offline Mode**: No internet required
7. **Mobile Support**: SSH-friendly interface
8. **Plugin System**: Custom question sources

### Integration Options
1. **Web Interface**: Browser-based alternative
2. **API Server**: REST API for programmatic access
3. **Database Backend**: Persistent storage
4. **Analytics**: Track performance metrics
5. **Notification System**: Quiz reminders

## Risks and Mitigations

### Technical Risks
1. **Terminal Compatibility**
   - Mitigation: Test on major terminals, graceful fallbacks

2. **Performance with Large Datasets**
   - Mitigation: Lazy loading, pagination, caching

3. **State Management Complexity**
   - Mitigation: Clear state machine, comprehensive tests

### User Experience Risks
1. **Learning Curve**
   - Mitigation: Built-in help, clear shortcuts, tutorial mode

2. **Data Loss**
   - Mitigation: Aggressive auto-save, session recovery

## Timeline

### Week 1: Foundation & Display
- Days 1-2: Bubbletea setup, welcome screen
- Days 3-4: Question navigation, markdown rendering
- Days 5-7: Styling system, basic input

### Week 2: Interaction & Storage
- Days 1-2: Response collection, answer input
- Days 3-4: Review screen, editing
- Days 5-7: CSV export, session management

### Week 3: Integration & Polish
- Days 1-2: quiz-eval integration, configuration
- Days 3-4: Error handling, help system
- Days 5-7: Testing, bug fixes, documentation

## Conclusion

This interactive quiz TUI will provide a delightful terminal-based experience for taking quizzes, leveraging the beautiful Charm ecosystem for a polished, professional interface. By reusing existing infrastructure (markdown parsing, CSV I/O, evaluation) and focusing on the interactive layer, we can deliver a powerful tool that enhances the quiz-taking workflow while maintaining seamless integration with the evaluation pipeline.

The phased implementation approach ensures steady progress with usable milestones, while the modular architecture allows for future enhancements and customization. The result will be a production-ready TUI that feels native to the terminal environment and makes quiz-taking efficient and enjoyable.
