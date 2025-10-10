package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/abhishek/quiz-evaluator/config"
	"github.com/abhishek/quiz-evaluator/internal/csv"
	"github.com/abhishek/quiz-evaluator/internal/evaluator"
	"github.com/abhishek/quiz-evaluator/internal/markdown"
	"github.com/abhishek/quiz-evaluator/internal/models"
	"github.com/spf13/cobra"
)

var (
	outputFile      string
	provider        string
	model           string
	workers         int
	dryRun          bool
	generateExample bool
)

// evaluateCmd represents the evaluate command
var evaluateCmd = &cobra.Command{
	Use:   "evaluate [input.csv] [questions-dir]",
	Short: "Evaluate quiz responses from CSV against question bank",
	Long: `Evaluate user responses from a CSV file by matching them with
question definitions in markdown files and using AI to generate scores
and detailed feedback.

The CSV file should have at minimum two columns:
- question_id: The ID of the question
- user_response: The user's response to evaluate

The questions directory should contain markdown files with question
definitions including evaluation criteria and sample answers.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if generateExample {
			return nil
		}
		if len(args) != 2 {
			return fmt.Errorf("requires exactly 2 arguments: input CSV file and questions directory")
		}
		return nil
	},
	RunE: runEvaluate,
}

func init() {
	rootCmd.AddCommand(evaluateCmd)

	// Command-specific flags
	evaluateCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output CSV file (default: input_evaluated.csv)")
	evaluateCmd.Flags().StringVar(&provider, "provider", "", "AI provider (openai, anthropic, gemini, ollama)")
	evaluateCmd.Flags().StringVar(&model, "model", "", "AI model to use")
	evaluateCmd.Flags().IntVarP(&workers, "workers", "w", 0, "number of parallel workers")
	evaluateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "validate inputs without running evaluation")
	evaluateCmd.Flags().BoolVar(&generateExample, "generate-example", false, "generate example configuration file")
}

func runEvaluate(cmd *cobra.Command, args []string) error {
	// Handle example generation
	if generateExample {
		if err := config.SaveExample("config.example.yaml"); err != nil {
			return fmt.Errorf("failed to generate example config: %w", err)
		}
		fmt.Println("Example configuration saved to config.example.yaml")
		return nil
	}

	inputCSV := args[0]
	questionsDir := args[1]

	// Set default output file if not specified
	if outputFile == "" {
		baseName := strings.TrimSuffix(filepath.Base(inputCSV), filepath.Ext(inputCSV))
		outputFile = baseName + "_evaluated.csv"
	}

	// Load configuration
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Override configuration with command flags
	if provider != "" {
		cfg.AI.Provider = provider
	}
	if model != "" {
		cfg.AI.Model = model
	}
	if workers > 0 {
		cfg.Evaluation.ParallelWorkers = workers
	}
	if verbose {
		cfg.Output.Verbose = true
	}
	cfg.Output.OutputFile = outputFile

	// Print configuration if verbose
	if cfg.Output.Verbose {
		printConfiguration(cfg)
	}

	// Step 1: Read input CSV
	fmt.Println("📖 Reading input CSV...")
	csvReader := csv.NewReader(inputCSV)
	responses, err := csvReader.ReadResponses()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}
	fmt.Printf("   Found %d responses to evaluate\n", len(responses))

	if len(responses) == 0 {
		fmt.Println("No responses to evaluate")
		return nil
	}

	// Step 2: Scan and index question files
	fmt.Println("\n🔍 Scanning questions directory...")
	scanner := markdown.NewScanner(questionsDir)
	questionIndex, err := scanner.ScanQuestions()
	if err != nil {
		return fmt.Errorf("failed to scan questions: %w", err)
	}
	fmt.Printf("   Found %d questions in the bank\n", len(questionIndex))

	// Validate that all response question IDs exist
	missingQuestions := validateQuestions(responses, questionIndex)
	if len(missingQuestions) > 0 {
		fmt.Printf("\n⚠️  Warning: %d question IDs not found in question bank:\n", len(missingQuestions))
		for _, id := range missingQuestions {
			fmt.Printf("   - %s\n", id)
		}
		if !askToContinue() {
			return fmt.Errorf("evaluation cancelled")
		}
	}

	// Dry run mode - stop here
	if dryRun {
		fmt.Println("\n✅ Dry run completed successfully")
		fmt.Println("   Input CSV and questions directory validated")
		return nil
	}

	// Step 3: Initialize AI client
	fmt.Println("\n🤖 Initializing AI client...")
	aiConfig := evaluator.AIClientConfig{
		Provider:    evaluator.AIProvider(cfg.AI.Provider),
		APIKey:      cfg.AI.APIKey,
		Model:       cfg.AI.Model,
		MaxTokens:   cfg.AI.MaxTokens,
		Temperature: cfg.AI.Temperature,
		Timeout:     cfg.AI.Timeout,
		MaxRetries:  cfg.AI.MaxRetries,
	}

	aiClient, err := evaluator.NewAIClient(aiConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize AI client: %w", err)
	}
	fmt.Printf("   Using %s provider with model %s\n", cfg.AI.Provider, cfg.AI.Model)

	// Step 4: Create evaluator
	evalConfig := evaluator.Config{
		ParallelWorkers: cfg.Evaluation.ParallelWorkers,
		Timeout:         30 * time.Second,
		Verbose:         cfg.Output.Verbose,
	}
	eval := evaluator.NewEvaluator(aiClient, questionIndex, evalConfig)

	// Step 5: Process evaluations
	fmt.Printf("\n⚙️  Evaluating responses (using %d parallel workers)...\n", cfg.Evaluation.ParallelWorkers)
	
	ctx := context.Background()
	startTime := time.Now()
	
	// Show progress
	done := make(chan struct{})
	go showProgress(len(responses), done)

	evaluations, err := eval.EvaluateResponses(ctx, responses)
	close(done)
	
	if err != nil {
		return fmt.Errorf("evaluation failed: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("\n   Completed %d evaluations in %s\n", len(evaluations), duration.Round(time.Second))

	// Step 6: Calculate statistics
	stats := eval.GetStatistics(evaluations)
	printStatistics(stats)

	// Step 7: Write output
	fmt.Printf("\n💾 Writing results to %s...\n", outputFile)
	csvWriter := csv.NewWriter(outputFile)
	if err := csvWriter.WriteEvaluations(evaluations); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	// Write summary if requested
	if cfg.Output.IncludeMetadata {
		if err := csvWriter.WriteSummary(evaluations); err != nil {
			fmt.Printf("   Warning: Failed to write summary: %v\n", err)
		} else {
			summaryFile := strings.TrimSuffix(outputFile, ".csv") + "_summary.txt"
			fmt.Printf("   Summary written to %s\n", summaryFile)
		}
	}

	fmt.Println("\n✅ Evaluation complete!")
	return nil
}

func printConfiguration(cfg *config.Config) {
	fmt.Println("\n📋 Configuration:")
	fmt.Printf("   AI Provider: %s\n", cfg.AI.Provider)
	fmt.Printf("   Model: %s\n", cfg.AI.Model)
	fmt.Printf("   Parallel Workers: %d\n", cfg.Evaluation.ParallelWorkers)
	fmt.Printf("   Output Format: %s\n", cfg.Output.Format)
	fmt.Println()
}

func validateQuestions(responses []models.UserResponse, index models.QuestionIndex) []string {
	var missing []string
	seen := make(map[string]bool)

	for _, response := range responses {
		if seen[response.QuestionID] {
			continue
		}
		seen[response.QuestionID] = true

		if _, exists := index[response.QuestionID]; !exists {
			missing = append(missing, response.QuestionID)
		}
	}

	return missing
}

func askToContinue() bool {
	fmt.Print("\nContinue anyway? (y/N): ")
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
}

func showProgress(total int, done chan struct{}) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	start := time.Now()
	count := 0

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			count++
			elapsed := time.Since(start).Seconds()
			rate := float64(count) / elapsed
			fmt.Printf("   Progress: ~%d/%d evaluations (%.1f/sec)\r", 
				int(rate*elapsed), total, rate)
		}
	}
}

func printStatistics(stats evaluator.Statistics) {
	fmt.Println("\n📊 Evaluation Statistics:")
	fmt.Printf("   Total Evaluated: %d\n", stats.TotalEvaluated)
	fmt.Printf("   Average Score: %.1f\n", stats.AverageScore)
	fmt.Printf("   Score Range: %.1f - %.1f\n", stats.MinScore, stats.MaxScore)
	
	if len(stats.ScoreDistribution) > 0 {
		fmt.Println("\n   Score Distribution:")
		for bucket, count := range stats.ScoreDistribution {
			percentage := float64(count) / float64(stats.TotalEvaluated) * 100
			fmt.Printf("     %s: %d (%.1f%%)\n", bucket, count, percentage)
		}
	}
}
