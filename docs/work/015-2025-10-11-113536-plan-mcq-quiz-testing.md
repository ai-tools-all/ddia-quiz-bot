# MCQ Quiz Testing Implementation Plan

## Date: 2025-10-11 11:35:36
## Category: Testing Plan
## Priority: High

## Overview
Based on analysis of the existing codebase, MCQ format, and parser implementation, this plan outlines comprehensive tests for ensuring MCQ quiz functionality.

## MCQ Format Analysis (from test-mcq-001.md)

### Frontmatter Format
```toml
+++
id = "fault-tolerance-mcq-L3-001"
title = "Primary-Backup Failure Types"
level = "L3"
category = "baseline"
type = "mcq"
+++
```

### Question Structure
```
## Question
What type of failures can primary-backup replication typically handle?

## Options
- A) Software bugs that cause incorrect calculations
- B) Hardware failures that cause the server to stop executing
- C) Network partitions that split the system
- D) All of the above

## Answer
B

## Explanation
Primary-backup replication handles...

## Hook
Understanding the limitations...

## Core Concepts
- Fail-stop failures
- Replication limitations
- Fault models
```

## Phase 1: Mock Data Setup

### 1.1 Create Test Data Directory
```bash
mkdir -p test-data/mcq/
mkdir -p test-data/mixed/
```

### 1.2 Create Mock MCQ Samples
- Copy `test-mcq-001.md` as base template
- Create variations:
  - `basic-mcq-001.md` - Simple MCQ format
  - `mcq-with-explanation.md` - Has detailed explanation
  - `mcq-with-hook.md` - Has engagement hook
  - `yaml-mcq.md` - YAML frontmatter for compatibility
  - `edge-case-mcq.md` - Malformed options for error testing

### 1.3 Create Mixed Content Samples
- Combine MCQ and subjective questions for integration testing
- Test session management with mixed types

## Phase 2: Core Testing Components

### 2.1 Parser Tests (extend `internal/markdown/parser_test.go`)

#### Test Cases
```go
func TestParseMCQFile(t *testing.T) {
    tests := []struct {
        name           string
        filename       string
        expectType     string
        expectOptions  int
        expectAnswer   string
        expectErr      bool
    }{
        {
            name:          "Basic MCQ format",
            filename:      "basic-mcq-001.md",
            expectType:    "mcq",
            expectOptions: 4,
            expectAnswer:  "B",
            expectErr:     false,
        },
        {
            name:          "MCQ with explanation",
            filename:      "mcq-with-explanation.md",
            expectType:    "mcq",
            expectOptions: 4,
            expectAnswer:  "C",
            expectErr:     false,
        },
        {
            name:          "YAML frontmatter compatibility",
            filename:      "yaml-mcq.md",
            expectType:    "mcq",
            expectOptions: 4,
            expectAnswer:  "A",
            expectErr:     false,
        },
        {
            name:          "Edge case - malformed options",
            filename:      "edge-case-mcq.md",
            expectErr:     true,
        },
    }
}
```

#### Key Validation Points
- **Options Parsing**: Test `parseMCQOptions()` function with various formats
  - Standard: `- A) Option text`
  - Variations: `A) text`, `* B) text`, `• C) text`
- **Frontmatter Support**: Both TOML (+++) and YAML (---)
- **Auto-detection**: When `type` field is missing
- **Required Fields**: id, level, options, answer
- **Optional Fields**: explanation, hook, core_concepts

### 2.2 SRS Integration Tests (new file: `internal/srs/mcq_test.go`)

#### Core SRS Tests
```go
func TestMCQCardTracking(t *testing.T) {
    card := &Card{
        QuestionID: "test-mcq-001",
        QuestionType: "mcq",
    }
    
    // Test correct answer
    card.RecordMCQAttempt("B", true, 30)
    assert.Equal(t, 1.0, card.MCQAccuracy)
    assert.Equal(t, "B", card.LastMCQChoice)
    assert.Equal(t, 0, len(card.IncorrectChoices))
    
    // Test incorrect answer
    card.RecordMCQAttempt("A", false, 45)
    assert.Equal(t, 0.5, card.MCQAccuracy) // 1/2 success rate
    assert.Contains(t, card.IncorrectChoices, "A")
}
```

#### SRS Scheduler Tests
```go
func TestSRSMCQIntegration(t *testing.T) {
    scheduler := NewScheduler()
    
    // Test MCQ review recording
    result, err := scheduler.RecordMCQReview("test-mcq-001", true, 30)
    assert.NoError(t, err)
    assert.Greater(t, result.Card.Interval, float64(0))
    
    // Test performance-based quality adjustment
    quickCorrectResult, _ := scheduler.RecordMCQReview("test-mcq-002", true, 10)
    slowCorrectResult, _ := scheduler.RecordMCQReview("test-mcq-003", true, 60)
    
    // Quick correct should get better quality than slow correct
    assert.Greater(t, quickCorrectResult.Quality, slowCorrectResult.Quality)
}
```

### 2.3 Parser + SRS Integration Tests (new file: `integration_test.go`)

#### End-to-End Workflow
```go
func TestMCQEndToEndWorkflow(t *testing.T) {
    // 1. Parse MCQ file
    parser := NewParser()
    question, err := parser.ParseQuestionFile("test-data/mcq/basic-mcq-001.md")
    require.NoError(t, err)
    assert.Equal(t, "mcq", question.Type)
    
    // 2. Create SRS card
    scheduler := NewScheduler()
    card, err := scheduler.CreateCard(question)
    require.NoError(t, err)
    
    // 3. Simulate MCQ attempt
    selectedAnswer := "B"
    isCorrect := selectedAnswer == question.Answer
    
    result, err := scheduler.RecordMCQReview(question.ID, isCorrect, 30)
    require.NoError(t, err)
    
    // 4. Verify SRS update
    updatedCard, err := scheduler.GetCard(question.ID)
    require.NoError(t, err)
    assert.Equal(t, selectedAnswer, updatedCard.LastMCQChoice)
    assert.Equal(t, question.Type, updatedCard.QuestionType)
}
```

#### Mixed Session Testing
```go
func TestMixedMCQSubjectiveSession(t *testing.T) {
    // Load mix of MCQ and subjective questions
    questions := loadTestQuestions("test-data/mixed/")
    
    assert.Greater(t, len(questions.MCQ), 0)
    assert.Greater(t, len(questions.Subjective), 0)
    
    // Simulate session handling
    session := &Session{Questions: questions}
    
    for _, q := range session.Questions {
        if q.Type == "mcq" {
            // Test MCQ path
            answer := simulateMCQAnswer(q)
            session.RecordMCQAnswer(q.ID, answer)
        } else {
            // Test subjective path
            answer := simulateSubjectiveAnswer(q)
            session.RecordSubjectiveAnswer(q.ID, answer)
        }
    }
    
    // Verify session persistence includes all metadata
    assert.Equal(t, len(session.Questions), len(session.Responses))
    
    for _, resp := range session.Responses {
        assert.NotNil(t, resp.QuestionType)
        if resp.QuestionType == "mcq" {
            assert.NotNil(t, resp.IsCorrect)
            assert.NotNil(t, resp.SelectedOption)
        }
    }
}
```

## Phase 3: Component Tests

### 3.1 MCQ UI Component Tests (extend `internal/tui/components/`)

#### MCQ Component State Management
```go
func TestMCQComponentStateTransitions(t *testing.T) {
    component := &MCQComponent{
        Options: []string{"A) Option 1", "B) Option 2", "C) Option 3", "D) Option 4"},
    }
    
    // Test navigation
    msg := tea.KeyMsg{Type: tea.KeyDown}
    model, cmd := component.Update(msg)
    assert.Equal(t, 1, model.SelectedIdx)
    
    // Test option selection
    msg = tea.KeyMsg{Type: tea.KeyEnter}
    model, cmd = component.Update(msg)
    assert.True(t, model.Submitted)
    assert.NotNil(t, cmd) // Should return saveMCQAnswerCmd
}
```

#### Visual Feedback Tests
```go
func TestMCQVisualFeedback(t *testing.T) {
    component := &MCQComponent{
        Options:     getTestOptions(),
        SelectedIdx: 1,
        Submitted:   true,
        CorrectIdx:  2, // Correct answer is C (index 2)
    }
    
    view := component.View()
    
    // Should show visual indicators
    assert.Contains(t, view, "✗") // Wrong indicator
    assert.Contains(t, view, "✓") // Correct indicator
    assert.Contains(t, view, "Explanation") // After submission
}
```

### 3.2 Session Management Tests (extend `internal/tui/session/`)

#### Mixed Session Handling
```go
func TestSessionMixedQuestionTypes(t *testing.T) {
    session := NewSession("testuser", "mixed")
    
    mcqQuestion := getMCQQuestion()
    subjectiveQuestion := getSubjectiveQuestion()
    
    session.Questions = append(session.Questions, mcqQuestion, subjectiveQuestion)
    
    // Test progress tracking
    assert.Equal(t, 0, session.Answered)
    
    // Answer MCQ question
    session.RecordResponse(&Response{
        QuestionID:   mcqQuestion.ID,
        Type:         "mcq",
        Content:      "B",
        IsCorrect:    boolPtr(true),
        SelectedOption: "B",
        TimeSpent:    30,
    })
    
    assert.Equal(t, 1, session.Answered)
    
    // Test session save/load
    savedData, err := json.Marshal(session)
    require.NoError(t, err)
    
    loadedSession := &Session{}
    err = json.Unmarshal(savedData, loadedSession)
    require.NoError(t, err)
    
    assert.Equal(t, session.Responses[0].QuestionType, "mcq")
    assert.Equal(t, session.Responses[0].IsCorrect, loadedSession.Responses[0].IsCorrect)
}
```

## Phase 4: Error Handling & Edge Cases

### 4.1 Error Scenarios
```go
func TestParserErrorHandling(t *testing.T) {
    tests := []struct {
        name        string
        content     string
        expectError string
    }{
        {
            name: "Missing answer in MCQ",
            content: `+++
type = "mcq"
+++

## Question
What is 2+2?

## Options
- A) 3
- B) 4
- C) 5
`,
            expectError: "answer not found",
        },
        {
            name: "Invalid option format", 
            content: `+++
type = "mcq"
+++

## Options
- Option 1 (no letter prefix)
- Option 2
`,
            expectError: "failed to parse MCQ options",
        },
        {
            name: "Empty options array",
            content: `+++
type = "mcq"
+++

## Options

## Answer
A
`,
            expectError: "no valid MCQ options found",
        },
    }
    
    parser := NewParser()
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tempFile := createTempFile(t, tt.content)
            _, err := parser.ParseQuestionFile(tempFile)
            assert.Error(t, err)
            assert.Contains(t, err.Error(), tt.expectError)
        })
    }
}
```

## Phase 5: Performance & Load Testing

### 5.1 Large Dataset Testing
```go
func TestLargeMCQDatasetPerformance(t *testing.T) {
    // Create 100 MCQ files
    generateMCQFiles(100, "test-data/performance/")
    
    start := time.Now()
    
    parser := NewParser()
    var questions []*models.Question
    
    for i := 1; i <= 100; i++ {
        q, err := parser.ParseQuestionFile(fmt.Sprintf("test-data/performance/mcq-%03d.md", i))
        require.NoError(t, err)
        questions = append(questions, q)
    }
    
    duration := time.Since(start)
    
    // Should parse 100 questions in under 1 second
    assert.Less(t, duration, time.Second)
    assert.Equal(t, 100, len(questions))
    
    // Verify all were parsed as MCQ
    for _, q := range questions {
        assert.Equal(t, "mcq", q.Type)
        assert.NotEmpty(t, q.Options)
        assert.NotEmpty(t, q.Answer)
    }
}
```

## Test Organization

### File Structure
```
internal/
├── markdown/
│   └── parser_test.go          # Extended with MCQ tests
├── srs/
│   └── mcq_test.go            # New MCQ SRS tests
├── tui/
│   ├── components/
│   │   └── mcq_test.go        # New MCQ component tests
│   └── session/
│       └── session_test.go    # Extended for mixed sessions
└── integration_test.go          # New end-to-end tests

test-data/
├── mcq/                       # MCQ test files
├── mixed/                     # Mixed session files
└── performance/               # Load test files
```

### Test Categories
1. **Unit Tests**: Individual component testing
2. **Integration Tests**: Cross-component interaction
3. **End-to-End Tests**: Complete workflow testing
4. **Performance Tests**: Load and timing validation
5. **Error Tests**: Edge case and failure handling

## Success Criteria

### Functional Requirements
- ✅ All MCQ formats parse correctly
- ✅ SRS tracks MCQ-specific metrics
- ✅ Mixed sessions work seamlessly
- ✅ UI component handles all states
- ✅ Session persistence includes all metadata

### Non-Functional Requirements
- ✅ Performance: <100ms per MCQ parse, <1s for 100 questions
- ✅ Coverage: >90% for MCQ-related code
- ✅ Reliability: All error conditions handled gracefully
- ✅ Compatibility: Both TOML and YAML frontmatter support

## Implementation Order

### Week 1: Foundation
1. Create test data directory and mock files
2. Extend parser tests for MCQ format
3. Implement SRS MCQ tracking tests

### Week 2: Integration  
4. Create end-to-end integration tests
5. Implement mixed session tests
6. Add UI component state tests

### Week 3: Polish
7. Add error handling and edge case tests
8. Implement performance testing
9. Achieve test coverage targets

## Mock Data Samples

### Basic MCQ Template
```markdown
+++
id = "basic-mcq-001"
title = "Test MCQ"
level = "L3"
category = "baseline"
type = "mcq"
+++

## Question
What is the capital of France?

## Options
- A) London
- B) Berlin
- C) Paris
- D) Madrid

## Answer
C

## Explanation
Paris is the capital and largest city of France.

## Hook
Geography knowledge is essential for understanding global contexts.
```

This plan ensures comprehensive testing of MCQ functionality while leveraging existing architecture and maintaining compatibility with the current codebase.
