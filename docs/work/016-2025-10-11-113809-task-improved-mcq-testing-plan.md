# Improved MCQ Quiz Testing Plan - Behavioral Focus

## Date: 2025-10-11 11:38:09
## Category: Task
## Priority: High

## Overview
This plan focuses on behavioral end-to-end testing to ensure the MCQ quiz system does exactly what it's supposed to do. Testing philosophy: verify the complete user journey from question loading to result tracking.

## Core Testing Strengths - Expanded

### 1. MCQ Format Variation Testing - Comprehensive Coverage

#### 1.1 Format Compatibility Matrix
Test every possible valid MCQ format the system should accept:

```go
func TestMCQFormatCompatibilityMatrix(t *testing.T) {
    formats := []struct {
        name           string
        content        string
        expectedParse  bool
        description    string
    }{
        {
            name: "Standard TOML frontmatter with dashes",
            content: `+++
id = "test-001"
type = "mcq"
+++

## Question
Text here?

## Options
- A) Option 1
- B) Option 2
- C) Option 3
- D) Option 4

## Answer
B`,
            expectedParse: true,
            description: "Most common format in existing content",
        },
        {
            name: "YAML frontmatter variant",
            content: `---
id: test-002
type: mcq
---

## Question
Text here?

## Options
- A) Option 1
- B) Option 2

## Answer
A`,
            expectedParse: true,
            description: "Alternative frontmatter syntax support",
        },
        {
            name: "Options without dash prefix",
            content: `+++
id = "test-003"
type = "mcq"
+++

## Question
Text?

## Options
A) Option 1
B) Option 2
C) Option 3

## Answer
C`,
            expectedParse: true,
            description: "Cleaner option format without bullet points",
        },
        {
            name: "Options with asterisk bullets",
            content: `+++
id = "test-004"
type = "mcq"
+++

## Question
Text?

## Options
* A) Option 1
* B) Option 2
* C) Option 3

## Answer
A`,
            expectedParse: true,
            description: "Alternative bullet style",
        },
        {
            name: "Mixed case section headers",
            content: `+++
id = "test-005"
type = "mcq"
+++

## question
Text?

## OPTIONS
- A) Option 1
- B) Option 2

## answer
B`,
            expectedParse: true,
            description: "Case-insensitive header parsing",
        },
        {
            name: "Multi-line options",
            content: `+++
id = "test-006"
type = "mcq"
+++

## Question
Which statement is correct?

## Options
- A) This is a very long option that spans
     multiple lines and contains detailed
     technical information
- B) Short option
- C) Another multi-line option with
     code examples and explanations

## Answer
A`,
            expectedParse: true,
            description: "Options can contain line breaks",
        },
        {
            name: "Options with markdown formatting",
            content: `+++
id = "test-007"
type = "mcq"
+++

## Question
Select the correct code:

## Options
- A) Use \`fmt.Println()\` for output
- B) Use **bold** for emphasis
- C) Use _italic_ for notes
- D) All contain `code blocks`

## Answer
D`,
            expectedParse: true,
            description: "Markdown within options preserved",
        },
    }

    parser := NewParser()
    for _, tt := range formats {
        t.Run(tt.name, func(t *testing.T) {
            tempFile := createTempMCQFile(t, tt.content)
            q, err := parser.ParseQuestionFile(tempFile)
            
            if tt.expectedParse {
                require.NoError(t, err, "Format should parse: %s", tt.description)
                assert.Equal(t, "mcq", q.Type)
                assert.NotEmpty(t, q.Options, "Options should be extracted")
                assert.NotEmpty(t, q.Answer, "Answer should be found")
            } else {
                assert.Error(t, err, "Format should fail: %s", tt.description)
            }
        })
    }
}
```

### 2. Edge Cases and Boundary Testing - Detailed Scenarios

#### 2.1 Option Count Variations
```go
func TestMCQOptionCountBehavior(t *testing.T) {
    testCases := []struct {
        name         string
        optionCount  int
        shouldWork   bool
        description  string
    }{
        {
            name:        "Two options (binary choice)",
            optionCount: 2,
            shouldWork:  true,
            description: "True/False style questions",
        },
        {
            name:        "Standard four options",
            optionCount: 4,
            shouldWork:  true,
            description: "Most common MCQ format",
        },
        {
            name:        "Six options",
            optionCount: 6,
            shouldWork:  true,
            description: "Extended choice questions",
        },
        {
            name:        "Ten options",
            optionCount: 10,
            shouldWork:  true,
            description: "Maximum reasonable options",
        },
        {
            name:        "Single option",
            optionCount: 1,
            shouldWork:  false,
            description: "Not a valid MCQ",
        },
        {
            name:        "Zero options",
            optionCount: 0,
            shouldWork:  false,
            description: "Empty options section",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            mcqContent := generateMCQWithNOptions(tc.optionCount)
            parser := NewParser()
            
            q, err := parser.ParseQuestionFile(createTempFile(t, mcqContent))
            
            if tc.shouldWork {
                require.NoError(t, err, tc.description)
                assert.Len(t, q.Options, tc.optionCount)
            } else {
                assert.Error(t, err, tc.description)
            }
        })
    }
}
```

#### 2.2 Answer Validation Scenarios
```go
func TestMCQAnswerValidation(t *testing.T) {
    tests := []struct {
        name          string
        options       []string
        answer        string
        shouldAccept  bool
        description   string
    }{
        {
            name:    "Single letter uppercase",
            options: []string{"A) Yes", "B) No"},
            answer:  "A",
            shouldAccept: true,
            description: "Standard single letter answer",
        },
        {
            name:    "Single letter lowercase",
            options: []string{"A) Yes", "B) No"},
            answer:  "a",
            shouldAccept: true,
            description: "Case-insensitive answer matching",
        },
        {
            name:    "Letter with parenthesis",
            options: []string{"A) Yes", "B) No"},
            answer:  "A)",
            shouldAccept: true,
            description: "Answer includes formatting",
        },
        {
            name:    "Full option text",
            options: []string{"A) Yes", "B) No"},
            answer:  "A) Yes",
            shouldAccept: true,
            description: "Complete option as answer",
        },
        {
            name:    "Multiple answers (comma-separated)",
            options: []string{"A) One", "B) Two", "C) Three"},
            answer:  "A, C",
            shouldAccept: true,
            description: "Multiple correct answers",
        },
        {
            name:    "Out of range answer",
            options: []string{"A) One", "B) Two"},
            answer:  "E",
            shouldAccept: false,
            description: "Answer not in options",
        },
        {
            name:    "Empty answer",
            options: []string{"A) One", "B) Two"},
            answer:  "",
            shouldAccept: false,
            description: "Missing answer field",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mcq := &MCQuestion{
                Options: tt.options,
                Answer:  tt.answer,
            }
            
            err := validateMCQAnswer(mcq)
            
            if tt.shouldAccept {
                assert.NoError(t, err, tt.description)
            } else {
                assert.Error(t, err, tt.description)
            }
        })
    }
}
```

### 3. End-to-End User Journey Testing - Complete Workflows

#### 3.1 Complete Quiz Session Flow
```go
func TestCompleteQuizSessionUserJourney(t *testing.T) {
    // Step 1: User starts quiz
    user := "testuser"
    quizType := "mcq"
    level := "L3"
    
    // Initialize session
    session, err := StartQuizSession(user, quizType, level)
    require.NoError(t, err)
    assert.NotNil(t, session.ID)
    assert.Equal(t, user, session.User)
    
    // Step 2: Load questions for the session
    questions, err := LoadQuestionsForSession(session)
    require.NoError(t, err)
    assert.Greater(t, len(questions), 0, "Should load at least one question")
    
    // Step 3: User navigates through each question
    for idx, question := range questions {
        t.Run(fmt.Sprintf("Question_%d", idx+1), func(t *testing.T) {
            // Verify question loaded correctly
            assert.NotEmpty(t, question.ID)
            assert.NotEmpty(t, question.Content)
            assert.Equal(t, "mcq", question.Type)
            assert.GreaterOrEqual(t, len(question.Options), 2)
            
            // Simulate user thinking time
            startTime := time.Now()
            time.Sleep(2 * time.Second) // Simulate reading
            
            // User selects an option
            selectedOption := simulateUserChoice(question.Options)
            
            // Record the response
            response := &Response{
                QuestionID:     question.ID,
                SelectedOption: selectedOption,
                TimeSpent:      time.Since(startTime).Seconds(),
                Timestamp:      time.Now(),
            }
            
            err := session.RecordResponse(response)
            assert.NoError(t, err)
            
            // Check if answer is correct
            isCorrect := checkAnswer(selectedOption, question.Answer)
            response.IsCorrect = &isCorrect
            
            // Update SRS based on result
            if session.TrackProgress {
                srsUpdate := UpdateSRS(question.ID, isCorrect, response.TimeSpent)
                assert.NotNil(t, srsUpdate)
                
                // Verify SRS calculations
                if isCorrect {
                    assert.Greater(t, srsUpdate.NextInterval, 0.0)
                } else {
                    assert.Equal(t, 0.0, srsUpdate.NextInterval)
                }
            }
        })
    }
    
    // Step 4: Complete session
    summary, err := session.Complete()
    require.NoError(t, err)
    
    // Verify session summary
    assert.Equal(t, len(questions), summary.TotalQuestions)
    assert.Equal(t, len(questions), summary.Answered)
    assert.GreaterOrEqual(t, summary.CorrectAnswers, 0)
    assert.LessOrEqual(t, summary.CorrectAnswers, summary.TotalQuestions)
    assert.Greater(t, summary.TotalTime, 0.0)
    
    // Step 5: Verify session persistence
    savedSession, err := LoadSession(session.ID)
    require.NoError(t, err)
    assert.Equal(t, session.ID, savedSession.ID)
    assert.Equal(t, len(questions), len(savedSession.Responses))
    
    // Step 6: Verify SRS updates persisted
    for _, response := range savedSession.Responses {
        if response.IsCorrect != nil && *response.IsCorrect {
            card, err := GetSRSCard(response.QuestionID)
            require.NoError(t, err)
            assert.Greater(t, card.Interval, 0.0)
            assert.Contains(t, []string{"Learning", "Review", "Mature"}, card.State)
        }
    }
}
```

#### 3.2 Mixed Question Type Session
```go
func TestMixedQuestionTypeSession(t *testing.T) {
    // Create a session with both MCQ and subjective questions
    session := &Session{
        User: "testuser",
        Type: "mixed",
    }
    
    // Load mixed question set
    mcqQuestions := loadMCQQuestions("test-data/mcq/", 5)
    subjectiveQuestions := loadSubjectiveQuestions("test-data/subjective/", 3)
    
    allQuestions := append(mcqQuestions, subjectiveQuestions...)
    session.Questions = shuffleQuestions(allQuestions)
    
    // Process each question based on type
    for _, q := range session.Questions {
        switch q.Type {
        case "mcq":
            // MCQ flow
            t.Run(fmt.Sprintf("MCQ_%s", q.ID), func(t *testing.T) {
                assert.NotEmpty(t, q.Options, "MCQ must have options")
                assert.NotEmpty(t, q.Answer, "MCQ must have answer")
                
                // Simulate MCQ interaction
                selected := q.Options[0] // Select first option
                isCorrect := extractLetter(selected) == q.Answer
                
                response := &Response{
                    QuestionID:     q.ID,
                    Type:           "mcq",
                    SelectedOption: selected,
                    IsCorrect:      &isCorrect,
                    TimeSpent:      15.5,
                }
                
                err := session.RecordResponse(response)
                assert.NoError(t, err)
            })
            
        case "subjective":
            // Subjective flow
            t.Run(fmt.Sprintf("Subjective_%s", q.ID), func(t *testing.T) {
                assert.Empty(t, q.Options, "Subjective should not have options")
                assert.NotEmpty(t, q.Content, "Subjective must have content")
                
                // Simulate subjective answer
                response := &Response{
                    QuestionID: q.ID,
                    Type:       "subjective",
                    Content:    "This is a subjective answer with explanation...",
                    TimeSpent:  45.0,
                }
                
                err := session.RecordResponse(response)
                assert.NoError(t, err)
            })
        }
    }
    
    // Verify mixed session handling
    assert.Equal(t, len(allQuestions), len(session.Responses))
    
    mcqCount := 0
    subjectiveCount := 0
    
    for _, resp := range session.Responses {
        switch resp.Type {
        case "mcq":
            mcqCount++
            assert.NotNil(t, resp.SelectedOption)
            assert.NotNil(t, resp.IsCorrect)
        case "subjective":
            subjectiveCount++
            assert.NotEmpty(t, resp.Content)
            assert.Nil(t, resp.SelectedOption)
        }
    }
    
    assert.Equal(t, 5, mcqCount, "Should have 5 MCQ responses")
    assert.Equal(t, 3, subjectiveCount, "Should have 3 subjective responses")
}
```

### 4. Performance Benchmarks - Realistic Scenarios

#### 4.1 Question Loading Performance
```go
func BenchmarkMCQLoadingPerformance(b *testing.B) {
    // Test different batch sizes
    batchSizes := []int{1, 10, 50, 100, 500}
    
    for _, size := range batchSizes {
        b.Run(fmt.Sprintf("Load_%d_MCQs", size), func(b *testing.B) {
            // Prepare test data
            questions := generateMCQFiles(size, "bench-data/")
            
            b.ResetTimer()
            
            for i := 0; i < b.N; i++ {
                parser := NewParser()
                loaded := make([]*Question, 0, size)
                
                for _, file := range questions {
                    q, err := parser.ParseQuestionFile(file)
                    if err != nil {
                        b.Fatal(err)
                    }
                    loaded = append(loaded, q)
                }
                
                // Verify all loaded
                if len(loaded) != size {
                    b.Fatalf("Expected %d questions, got %d", size, len(loaded))
                }
            }
        })
    }
}

func TestMCQLoadingTimeConstraints(t *testing.T) {
    constraints := []struct {
        count       int
        maxDuration time.Duration
        description string
    }{
        {1, 10 * time.Millisecond, "Single MCQ should load in <10ms"},
        {10, 50 * time.Millisecond, "10 MCQs should load in <50ms"},
        {100, 500 * time.Millisecond, "100 MCQs should load in <500ms"},
        {1000, 5 * time.Second, "1000 MCQs should load in <5s"},
    }
    
    for _, c := range constraints {
        t.Run(fmt.Sprintf("%d_questions", c.count), func(t *testing.T) {
            files := generateMCQFiles(c.count, "perf-test/")
            
            start := time.Now()
            
            parser := NewParser()
            for _, file := range files {
                _, err := parser.ParseQuestionFile(file)
                require.NoError(t, err)
            }
            
            duration := time.Since(start)
            assert.LessOrEqual(t, duration, c.maxDuration, c.description)
        })
    }
}
```

#### 4.2 SRS Update Performance
```go
func BenchmarkSRSUpdatePerformance(b *testing.B) {
    // Setup scheduler with existing cards
    scheduler := NewScheduler()
    
    // Pre-load cards
    for i := 0; i < 1000; i++ {
        card := &Card{
            QuestionID: fmt.Sprintf("mcq-%04d", i),
            Type:       "mcq",
            State:      "Learning",
            Interval:   1.0,
        }
        scheduler.AddCard(card)
    }
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        // Simulate random MCQ review
        questionID := fmt.Sprintf("mcq-%04d", rand.Intn(1000))
        isCorrect := rand.Float32() > 0.5
        timeSpent := rand.Float64() * 60 // 0-60 seconds
        
        _, err := scheduler.RecordMCQReview(questionID, isCorrect, timeSpent)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### 5. Clear Test Organization - Behavioral Focus

#### Test File Structure (Simplified)
```
internal/
├── markdown/
│   ├── parser_mcq_test.go      # MCQ parsing behavior tests
│   └── parser_formats_test.go   # Format compatibility tests
├── srs/
│   └── mcq_behavior_test.go    # SRS MCQ tracking behavior
├── tui/
│   └── mcq_flow_test.go        # MCQ UI flow testing
└── e2e/
    ├── quiz_session_test.go     # Complete session flows
    ├── mixed_session_test.go    # Mixed question handling
    └── user_journey_test.go     # Full user journeys

test-data/
├── valid-mcqs/                  # Valid format samples
├── invalid-mcqs/                # Invalid format samples
├── edge-cases/                  # Boundary condition samples
└── realistic-content/           # Real quiz content samples
```

## Test Execution Strategy

### 1. Continuous Testing
```bash
# Run all behavioral tests
go test ./... -tags=behavior

# Run only MCQ tests
go test ./... -run="MCQ"

# Run with coverage
go test ./... -cover -coverprofile=coverage.out
```

### 2. Test Data Management
```go
func TestMain(m *testing.M) {
    // Setup: Create test data
    setupTestData()
    
    // Run tests
    code := m.Run()
    
    // Cleanup: Remove test data
    cleanupTestData()
    
    os.Exit(code)
}
```

### 3. Behavioral Test Categories

#### Core Behaviors to Verify
1. **Question Loading**: Can the system load all MCQ formats from files?
2. **Option Parsing**: Are all option formats correctly extracted?
3. **Answer Validation**: Does answer checking work for all formats?
4. **Session Flow**: Can users complete a full quiz session?
5. **Progress Tracking**: Does SRS correctly track MCQ performance?
6. **Persistence**: Do sessions and progress save/load correctly?
7. **Mixed Content**: Can the system handle mixed question types?

## Success Metrics

### Functional Behavior Verification
- ✅ All documented MCQ formats parse correctly
- ✅ User can complete full quiz session without errors
- ✅ Answers are validated accurately
- ✅ Progress is tracked and persisted
- ✅ Mixed sessions work seamlessly
- ✅ Performance meets user expectations (<100ms response)

### Test Quality Metrics
- ✅ All user journeys covered
- ✅ All format variations tested
- ✅ Edge cases handled gracefully
- ✅ Tests run in <30 seconds
- ✅ Tests are maintainable and clear

## Implementation Priority

### Phase 1: Core Behavior (Week 1)
1. MCQ format parsing tests
2. Answer validation tests
3. Basic session flow tests

### Phase 2: Complete Journeys (Week 2)
4. End-to-end user journey tests
5. Mixed question session tests
6. SRS integration behavior tests

### Phase 3: Edge Cases & Performance (Week 3)
7. Edge case handling tests
8. Performance benchmarks
9. Test cleanup and documentation

This behavioral testing approach ensures the MCQ quiz system does exactly what users expect, with comprehensive coverage of all real-world usage patterns.
