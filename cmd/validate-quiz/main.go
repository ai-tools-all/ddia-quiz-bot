package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/abhishek/ddia-clicker/internal/quiz"
)

var (
	outputFormat string
	verbose      bool
	recursive    bool
	quiet        bool
	strictMode   bool
)

var rootCmd = &cobra.Command{
	Use:   "validate-quiz [file/directory]",
	Short: "Validate markdown quiz files for TUI compatibility",
	Long: `Validate quiz markdown files to ensure they can be properly parsed by the TUI application.
	
This tool checks for:
- Valid YAML frontmatter (if present)
- Required fields (question ID, main question)
- Proper markdown structure
- Section formatting
- List formatting in concepts sections
- Evaluation rubric structure

Examples:
  # Validate a single file
  validate-quiz question.md
  
  # Validate all files in a directory
  validate-quiz ./questions
  
  # Validate recursively with JSON output
  validate-quiz -r -f json ./content
  
  # Quiet mode - only show errors
  validate-quiz -q ./content`,
	Args: cobra.ExactArgs(1),
	RunE: runValidation,
}

func init() {
	rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "text", "Output format (text, json)")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed validation information")
	rootCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively validate directories")
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Only show files with errors")
	rootCmd.Flags().BoolVarP(&strictMode, "strict", "s", false, "Treat warnings as errors")
}

func runValidation(cmd *cobra.Command, args []string) error {
	path := args[0]

	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot access path: %w", err)
	}

	// Create validator and configure it
	validator := quiz.NewValidator()
	validator.SetStrictMode(strictMode)

	// Validate files
	var validations []quiz.FileValidation
	if fileInfo.IsDir() {
		validations, err = validator.ValidateDirectory(path, recursive)
		if err != nil {
			return err
		}
	} else {
		validation := validator.ValidateFile(path)
		validations = []quiz.FileValidation{validation}
	}

	// Generate and output report
	reporter := quiz.NewReporter(verbose, quiet)
	report := reporter.Generate(validations)

	switch outputFormat {
	case "json":
		err = reporter.OutputJSON(os.Stdout, report)
	default:
		err = reporter.OutputText(os.Stdout, report)
	}

	if err != nil {
		return err
	}

	// Exit with non-zero code if there are invalid files
	if report.InvalidFiles > 0 {
		os.Exit(1)
	}

	return nil
}



func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
