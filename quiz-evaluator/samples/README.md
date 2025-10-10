# Quiz Evaluator Sample Files

This directory contains sample files to help you get started with the Quiz Evaluator CLI.

## Files

- `responses.csv` - Sample user responses to evaluate
- `questions/` - Sample question definitions in markdown format
- `config.yaml` - Sample configuration file

## Quick Start

1. **Set your OpenAI API key:**
   ```bash
   export OPENAI_API_KEY="your-api-key-here"
   ```

2. **Run evaluation with sample data:**
   ```bash
   # From the quiz-evaluator directory
   ./quiz-eval evaluate samples/responses.csv samples/questions
   
   # Or with custom config
   ./quiz-eval evaluate samples/responses.csv samples/questions --config samples/config.yaml
   
   # With verbose output
   ./quiz-eval evaluate samples/responses.csv samples/questions -v
   ```

3. **Check the results:**
   - Output will be saved to `responses_evaluated.csv`
   - Summary statistics will be displayed in the terminal
   - If metadata is enabled, a summary file will also be created

## Expected Output

The evaluation will:
1. Read 2 sample responses from the CSV
2. Match them with question definitions
3. Use AI to evaluate each response
4. Generate scores and detailed feedback
5. Save results to CSV

Sample output format:
```csv
question_id,question_title,level,user_response,score,strengths,improvements,feedback
raft-l5-consensus,"Raft Consensus Algorithm","L5","...",85.0,"Good understanding","Could elaborate on safety","Strong grasp of core concepts..."
```

## Customization

You can customize the evaluation by:
- Adding more questions to the `questions/` directory
- Modifying the `config.yaml` file
- Using different AI providers or models
- Adjusting the number of parallel workers

## Testing Different Scenarios

### Dry Run
Test without making API calls:
```bash
./quiz-eval evaluate samples/responses.csv samples/questions --dry-run
```

### Different Output Format
Specify a custom output file:
```bash
./quiz-eval evaluate samples/responses.csv samples/questions -o my_results.csv
```

### Using Different AI Models
Override the model from command line:
```bash
./quiz-eval evaluate samples/responses.csv samples/questions --model gpt-3.5-turbo
```
