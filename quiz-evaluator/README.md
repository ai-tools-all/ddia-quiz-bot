# Quiz Evaluator CLI

AI-powered CLI tool for evaluating subjective quiz responses by matching question IDs from CSV input with markdown question files and generating scores with detailed feedback.

## Features

- 🤖 **AI-Powered Evaluation**: Uses OpenAI GPT-4 (and other providers) to evaluate responses
- 📊 **Batch Processing**: Efficiently processes multiple responses in parallel
- 📝 **Detailed Feedback**: Provides scores, strengths, improvements, and constructive feedback
- 🔄 **Multiple AI Providers**: Supports OpenAI, Anthropic, Gemini, and local LLMs (Ollama)
- 📈 **Statistics & Reports**: Generates evaluation statistics and summary reports
- ⚡ **High Performance**: Concurrent processing with configurable worker pools

## Installation

### From Source

```bash
git clone https://github.com/abhishek/quiz-evaluator.git
cd quiz-evaluator
go build -o quiz-eval
```

### Using Go Install

```bash
go install github.com/abhishek/quiz-evaluator@latest
```

## Quick Start

1. **Set up your API key** (for OpenAI):
   ```bash
   export OPENAI_API_KEY="your-api-key-here"
   ```

2. **Prepare your input CSV** with columns:
   - `question_id`: The ID of the question
   - `user_response`: The user's response to evaluate

3. **Prepare question files** in markdown format with:
   - Question ID in frontmatter or content
   - Core concepts and peripheral concepts
   - Sample answers and evaluation rubric

4. **Run the evaluation**:
   ```bash
   quiz-eval evaluate responses.csv ./questions
   ```

## Usage

### Basic Command

```bash
quiz-eval evaluate [input.csv] [questions-dir]
```

### Options

```
Flags:
  -o, --output string      Output CSV file (default: input_evaluated.csv)
      --provider string    AI provider (openai, anthropic, gemini, ollama)
      --model string       AI model to use
  -w, --workers int       Number of parallel workers (default: 5)
      --dry-run           Validate inputs without running evaluation
      --generate-example  Generate example configuration file
  -v, --verbose          Verbose output
      --config string     Config file (default: ./config.yaml)
  -h, --help             Help for evaluate
```

### Examples

```bash
# Basic evaluation
quiz-eval evaluate responses.csv ./questions

# With custom output file
quiz-eval evaluate responses.csv ./questions -o results.csv

# Using different AI provider
quiz-eval evaluate responses.csv ./questions --provider anthropic --model claude-3

# Parallel processing with 10 workers
quiz-eval evaluate responses.csv ./questions --workers 10

# Dry run to validate inputs
quiz-eval evaluate responses.csv ./questions --dry-run

# Generate example configuration
quiz-eval evaluate --generate-example
```

## Configuration

Create a `config.yaml` file for persistent settings:

```yaml
ai:
  provider: "openai"
  api_key: "${OPENAI_API_KEY}"
  model: "gpt-4-turbo-preview"
  max_tokens: 1500
  temperature: 0.3

evaluation:
  parallel_workers: 5
  retry_attempts: 3
  retry_delay: "2s"

output:
  format: "csv"
  include_metadata: true
  verbose: false
```

## Question File Format

Questions should be in markdown format with YAML frontmatter:

```markdown
---
id: "q001"
title: "Distributed Consensus"
level: "L5"
category: "Distributed Systems"
---

## Question
Explain the Raft consensus algorithm...

## Core Concepts
- Leader election
- Log replication
- Safety properties

## Peripheral Concepts
- Performance optimizations
- Practical implementations

## Sample Excellent Answer
[Detailed sample answer...]

## Sample Acceptable Answer
[Basic sample answer...]

## Evaluation Rubric
Technical Accuracy: 40%
Completeness: 30%
Clarity: 20%
Examples: 10%
```

## Input CSV Format

```csv
question_id,user_response
q001,"User's answer to question 1..."
q002,"User's answer to question 2..."
```

## Output

The tool generates:
1. **Evaluated CSV** with scores, feedback, and improvements
2. **Summary Report** with statistics and score distribution

### Output CSV Format

```csv
question_id,question_title,level,user_response,score,strengths,improvements,feedback
q001,"Distributed Consensus","L5","User answer...",85.0,"Clear explanation; Good examples","Missing safety properties","Good understanding..."
```

## Environment Variables

- `OPENAI_API_KEY`: OpenAI API key
- `AI_API_KEY`: Generic AI API key
- `QUIZ_EVAL_API_KEY`: Tool-specific API key
- `AI_PROVIDER`: Override AI provider
- `AI_MODEL`: Override AI model

## Performance

- Processes 20-30 evaluations per minute (with default settings)
- Supports concurrent evaluation with configurable workers
- Automatic retry on transient failures
- Caching support for repeated evaluations (future feature)

## Troubleshooting

### Common Issues

1. **API Key not found**
   - Set environment variable: `export OPENAI_API_KEY="your-key"`
   - Or add to config.yaml

2. **Question not found**
   - Ensure question IDs in CSV match those in markdown files
   - Check markdown file format and frontmatter

3. **Rate limiting**
   - Reduce parallel workers: `--workers 2`
   - Add delay in configuration

4. **Timeout errors**
   - Increase timeout in config
   - Reduce response length

## Development

### Building from Source

```bash
# Clone repository
git clone https://github.com/abhishek/quiz-evaluator.git
cd quiz-evaluator

# Install dependencies
go mod download

# Build
go build -o quiz-eval

# Run tests
go test ./...
```

### Project Structure

```
quiz-evaluator/
├── cmd/              # CLI commands
├── internal/         # Internal packages
│   ├── csv/         # CSV I/O
│   ├── markdown/    # Markdown parsing
│   ├── evaluator/   # AI evaluation
│   └── models/      # Data models
├── config/          # Configuration
└── main.go          # Entry point
```

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details

## Support

For issues and questions:
- GitHub Issues: [github.com/abhishek/quiz-evaluator/issues](https://github.com/abhishek/quiz-evaluator/issues)
- Documentation: [github.com/abhishek/quiz-evaluator/wiki](https://github.com/abhishek/quiz-evaluator/wiki)
