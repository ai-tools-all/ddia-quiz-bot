# Plan: Golang CLI for AI-Powered Quiz Response Evaluation

## Date: 2025-10-10
## Category: Feature Plan
## Status: In Planning

## Overview
Design and implement a Golang CLI tool that evaluates subjective quiz responses using AI by matching question IDs from CSV input with markdown question files and generating scores with detailed feedback.

## Project Structure
```
quiz-evaluator/
├── main.go                 # CLI entry point
├── go.mod                  # Go module definition
├── go.sum                  # Dependency locks
├── cmd/
│   └── evaluate.go         # Main evaluate command
├── internal/
│   ├── csv/
│   │   ├── reader.go       # CSV input reader
│   │   └── writer.go       # CSV output writer
│   ├── markdown/
│   │   ├── parser.go       # Markdown file parser
│   │   └── scanner.go      # Recursive directory scanner
│   ├── evaluator/
│   │   ├── ai_client.go    # AI service client
│   │   ├── evaluator.go    # Core evaluation logic
│   │   └── prompt.go       # Prompt builder
│   └── models/
│       ├── question.go     # Question model
│       ├── response.go     # Response model
│       └── evaluation.go   # Evaluation result model
├── config/
│   └── config.go           # Configuration management
└── README.md

```

## Implementation Plan

### Phase 1: Core Data Structures and Models

#### 1.1 Define Data Models
```go
// models/question.go
type Question struct {
    ID                string
    Title             string
    MainQuestion      string
    CoreConcepts      []string
    PeripheralConcepts []string
    SampleExcellent   string
    SampleAcceptable  string
    Level             string
    Category          string
    EvaluationRubric  map[string]string
}

// models/response.go
type UserResponse struct {
    QuestionID string
    Response   string
}

// models/evaluation.go
type Evaluation struct {
    QuestionID    string
    UserResponse  string
    Score         float64
    Feedback      string
    QuestionTitle string
    Level         string
    Strengths     []string
    Improvements  []string
}
```

#### 1.2 CSV Structures
```go
type InputRecord struct {
    QuestionID   string `csv:"question_id"`
    UserResponse string `csv:"user_response"`
}

type OutputRecord struct {
    QuestionID    string  `csv:"question_id"`
    UserResponse  string  `csv:"user_response"`
    Score         float64 `csv:"score"`
    Feedback      string  `csv:"feedback"`
    QuestionTitle string  `csv:"question_title"`
}
```

### Phase 2: Markdown Processing

#### 2.1 Markdown Scanner
- Recursively scan directory for `.md` files
- Filter for subjective question files
- Build index of question_id → file_path mapping

#### 2.2 Markdown Parser
- Parse YAML frontmatter for metadata
- Extract question components:
  - Question ID
  - Title
  - Main question text
  - Core concepts
  - Peripheral concepts
  - Sample answers
  - Evaluation rubric

```go
// markdown/parser.go
func ParseQuestionFile(filepath string) (*models.Question, error) {
    // 1. Read file content
    // 2. Split frontmatter and body
    // 3. Parse YAML frontmatter
    // 4. Parse markdown sections
    // 5. Build Question model
}
```

### Phase 3: AI Integration

#### 3.1 AI Client Interface
```go
type AIClient interface {
    EvaluateResponse(ctx context.Context, req EvaluationRequest) (*EvaluationResponse, error)
}

type EvaluationRequest struct {
    Question      *models.Question
    UserResponse  string
    MaxTokens     int
    Temperature   float64
}
```

#### 3.2 Supported AI Providers
- **OpenAI GPT-4**: Primary choice for quality
- **Anthropic Claude**: Alternative with good reasoning
- **Google Gemini**: Backup option
- **Local LLM**: Ollama integration for offline/privacy

#### 3.3 Prompt Engineering
```go
// evaluator/prompt.go
func BuildEvaluationPrompt(q *Question, response string) string {
    template := `
You are evaluating a subjective technical interview response.

Question: {{.MainQuestion}}
Level: {{.Level}}

Core Concepts (60% weight):
{{range .CoreConcepts}}- {{.}}
{{end}}

Peripheral Concepts (40% weight):
{{range .PeripheralConcepts}}- {{.}}
{{end}}

Sample Excellent Answer:
{{.SampleExcellent}}

Sample Acceptable Answer:
{{.SampleAcceptable}}

User Response:
{{.UserResponse}}

Evaluate based on:
1. Coverage of core concepts (60%)
2. Coverage of peripheral concepts (40%)
3. Technical accuracy
4. Clarity of explanation
5. Use of examples

Provide:
1. Score (0-100)
2. Key strengths (bullet points)
3. Areas for improvement (bullet points)
4. Overall feedback (2-3 sentences)

Output format:
SCORE: [number]
STRENGTHS:
- [strength 1]
- [strength 2]
IMPROVEMENTS:
- [improvement 1]
- [improvement 2]
FEEDBACK: [overall feedback]
`
}
```

### Phase 4: CLI Implementation

#### 4.1 Command Structure
```go
// Using cobra for CLI framework
var rootCmd = &cobra.Command{
    Use:   "quiz-eval",
    Short: "Evaluate subjective quiz responses using AI",
}

var evaluateCmd = &cobra.Command{
    Use:   "evaluate [input.csv] [questions-dir]",
    Short: "Evaluate responses from CSV against question bank",
    Args:  cobra.ExactArgs(2),
    Run:   runEvaluate,
}
```

#### 4.2 Configuration Options
```yaml
# config.yaml
ai:
  provider: "openai"  # openai, anthropic, gemini, ollama
  api_key: "${AI_API_KEY}"
  model: "gpt-4"
  max_tokens: 1500
  temperature: 0.3

evaluation:
  parallel_workers: 5
  retry_attempts: 3
  retry_delay: 2s
  
output:
  format: "csv"  # csv, json, markdown
  include_metadata: true
  verbose: false
```

### Phase 5: Core Evaluation Logic

#### 5.1 Main Evaluation Flow
```go
func runEvaluate(inputCSV, questionsDir string) error {
    // 1. Load configuration
    config := LoadConfig()
    
    // 2. Initialize AI client
    aiClient := NewAIClient(config.AI)
    
    // 3. Scan and index markdown files
    questionIndex := ScanQuestions(questionsDir)
    
    // 4. Read input CSV
    responses := ReadInputCSV(inputCSV)
    
    // 5. Process evaluations
    results := make([]Evaluation, 0)
    for _, response := range responses {
        // Find question
        question := questionIndex[response.QuestionID]
        if question == nil {
            log.Warn("Question not found:", response.QuestionID)
            continue
        }
        
        // Evaluate with AI
        eval := evaluateResponse(aiClient, question, response)
        results = append(results, eval)
    }
    
    // 6. Write output CSV
    WriteOutputCSV(results, config.Output)
    
    return nil
}
```

#### 5.2 Parallel Processing
```go
func processInParallel(responses []UserResponse, workers int) []Evaluation {
    jobs := make(chan UserResponse, len(responses))
    results := make(chan Evaluation, len(responses))
    
    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go evaluationWorker(jobs, results, &wg)
    }
    
    // Queue jobs
    for _, resp := range responses {
        jobs <- resp
    }
    close(jobs)
    
    // Collect results
    wg.Wait()
    close(results)
    
    // Aggregate
    evaluations := make([]Evaluation, 0)
    for eval := range results {
        evaluations = append(evaluations, eval)
    }
    
    return evaluations
}
```

### Phase 6: Error Handling and Validation

#### 6.1 Input Validation
- Validate CSV format and required columns
- Check question_id format
- Validate non-empty responses
- Verify questions directory exists

#### 6.2 Error Recovery
- Retry failed AI calls with exponential backoff
- Handle partial failures gracefully
- Log errors to separate error file
- Continue processing remaining questions

#### 6.3 Result Validation
- Validate score ranges (0-100)
- Ensure feedback is not empty
- Check for AI response parsing errors
- Flag suspicious evaluations for review

### Phase 7: Testing Strategy

#### 7.1 Unit Tests
- CSV parsing/writing
- Markdown parsing
- Prompt generation
- Score parsing

#### 7.2 Integration Tests
- End-to-end evaluation flow
- AI client mocking
- Error handling scenarios

#### 7.3 Test Data
- Sample questions in markdown format
- Sample responses CSV
- Expected evaluation outputs

### Phase 8: Advanced Features

#### 8.1 Batch Processing
- Support for large CSV files
- Progress indicators
- Resume capability for interrupted runs

#### 8.2 Reporting
- Generate HTML report with visualizations
- Statistics dashboard (score distribution, common weaknesses)
- Export to different formats (JSON, Excel)

#### 8.3 Caching
- Cache AI evaluations to avoid re-processing
- Question file caching for performance

#### 8.4 Multi-model Evaluation
- Compare evaluations from different AI models
- Ensemble scoring for higher accuracy
- A/B testing capabilities

## CLI Usage Examples

### Basic Usage
```bash
# Evaluate responses from CSV
quiz-eval evaluate responses.csv ./questions

# With custom output file
quiz-eval evaluate responses.csv ./questions -o results.csv

# With specific AI provider
quiz-eval evaluate responses.csv ./questions --provider anthropic --model claude-3
```

### Advanced Usage
```bash
# Parallel processing with 10 workers
quiz-eval evaluate responses.csv ./questions --workers 10

# Generate detailed report
quiz-eval evaluate responses.csv ./questions --format html --report detailed

# Dry run to validate inputs
quiz-eval evaluate responses.csv ./questions --dry-run

# Use local LLM
quiz-eval evaluate responses.csv ./questions --provider ollama --model llama2
```

## Dependencies

### Core Dependencies
```go
// go.mod
module github.com/user/quiz-evaluator

require (
    github.com/spf13/cobra v1.7.0        // CLI framework
    github.com/spf13/viper v1.16.0       // Configuration
    github.com/gocarina/gocsv v0.0.0     // CSV handling
    gopkg.in/yaml.v3 v3.0.1              // YAML parsing
    github.com/gomarkdown/markdown v0.0.0 // Markdown parsing
    github.com/sashabaranov/go-openai v1.15.0 // OpenAI client
)
```

### Optional Dependencies
```go
require (
    github.com/schollz/progressbar/v3 v3.13.0  // Progress bars
    github.com/fatih/color v1.15.0             // Colored output
    github.com/olekukonko/tablewriter v0.0.5   // Table formatting
    github.com/go-echarts/go-echarts/v2 v2.2.0 // Charts for reports
)
```

## Performance Considerations

### Optimization Strategies
1. **Concurrent Processing**: Process multiple evaluations in parallel
2. **Batching**: Group API calls when provider supports batch operations
3. **Caching**: Cache parsed markdown files and AI responses
4. **Streaming**: Stream large CSV files instead of loading into memory
5. **Connection Pooling**: Reuse HTTP connections for AI API calls

### Benchmarks
- Target: Process 100 responses in < 5 minutes
- Memory usage: < 500MB for 1000 questions
- API rate limiting: Implement adaptive rate limiting

## Security Considerations

1. **API Key Management**
   - Support environment variables
   - Use secure key storage (keyring)
   - Never log API keys

2. **Input Sanitization**
   - Validate and sanitize user responses
   - Prevent prompt injection attacks
   - Limit response lengths

3. **Data Privacy**
   - Option to use local LLMs for sensitive data
   - Audit logging for compliance
   - Data retention policies

## Deployment Options

### 1. Binary Distribution
```bash
# Build for multiple platforms
make build-all

# Outputs
dist/
├── quiz-eval-darwin-amd64
├── quiz-eval-darwin-arm64
├── quiz-eval-linux-amd64
└── quiz-eval-windows-amd64.exe
```

### 2. Docker Container
```dockerfile
FROM golang:1.21-alpine AS builder
# Build stage

FROM alpine:latest
# Runtime stage
COPY --from=builder /app/quiz-eval /usr/local/bin/
ENTRYPOINT ["quiz-eval"]
```

### 3. GitHub Action
```yaml
name: Quiz Evaluation
on:
  workflow_dispatch:
    inputs:
      csv_file:
        description: 'CSV file with responses'
        required: true

jobs:
  evaluate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: user/quiz-eval-action@v1
        with:
          csv-file: ${{ github.event.inputs.csv_file }}
          questions-dir: ./questions
```

## Success Metrics

1. **Accuracy**: 90% correlation with human evaluators
2. **Performance**: < 3 seconds per evaluation
3. **Reliability**: < 1% failure rate
4. **Usability**: < 5 minute setup time

## Timeline Estimate

- **Week 1**: Core models and CSV/Markdown processing
- **Week 2**: AI integration and prompt engineering
- **Week 3**: CLI implementation and error handling
- **Week 4**: Testing, optimization, and documentation
- **Week 5**: Advanced features and deployment

## Next Steps

1. [ ] Validate approach with stakeholders
2. [ ] Set up project repository and CI/CD
3. [ ] Implement Phase 1 (Core Data Structures)
4. [ ] Create test data and examples
5. [ ] Begin AI integration spike
