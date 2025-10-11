package main

import (
	"fmt"
	"os"

	"github.com/abhishek/ddia-clicker/internal/toc"
	"github.com/spf13/cobra"
)

var (
	maxDepth         int
	format           string
	recursive        bool
	outputFile       string
	confirmThreshold int
	skipConfirmation bool
)

var rootCmd = &cobra.Command{
	Use:   "md-toc [file/directory]",
	Short: "Generate table of contents for markdown files",
	Long: `Generate table of contents (TOC) for markdown files.

This tool extracts headers from markdown files and generates a hierarchical
table of contents with proper indentation and anchor links.

Supports:
- Single file or directory processing
- Multiple output formats (markdown, json, text)
- Depth control (H1-H6)
- Recursive directory traversal

Examples:
  # Generate TOC for a single file
  md-toc README.md
  
  # Generate TOC for all markdown files in a directory
  md-toc ./docs
  
  # Limit to H1-H3 headers only
  md-toc --depth 3 README.md
  
  # Output as JSON
  md-toc --format json ./docs
  
  # Save to file
  md-toc --output toc.md ./docs
  
  # Non-recursive directory scan
  md-toc --recursive=false ./docs
  
  # Skip confirmation for large directories
  md-toc --yes ./large-docs
  
  # Set custom confirmation threshold
  md-toc --confirm-threshold 500 ./docs`,
	Args: cobra.ExactArgs(1),
	RunE: runTOC,
}

func init() {
	rootCmd.Flags().IntVarP(&maxDepth, "depth", "d", 6, "Maximum header depth (1-6)")
	rootCmd.Flags().StringVarP(&format, "format", "f", "markdown", "Output format (markdown, json, text)")
	rootCmd.Flags().BoolVarP(&recursive, "recursive", "r", true, "Process directories recursively")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
	rootCmd.Flags().IntVarP(&confirmThreshold, "confirm-threshold", "t", 1000, "Prompt for confirmation if files exceed this number (0 to disable)")
	rootCmd.Flags().BoolVarP(&skipConfirmation, "yes", "y", false, "Skip confirmation prompt for large number of files")
}

func runTOC(cmd *cobra.Command, args []string) error {
	path := args[0]

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}

	if maxDepth < 1 || maxDepth > 6 {
		return fmt.Errorf("depth must be between 1 and 6, got: %d", maxDepth)
	}

	validFormats := map[string]bool{
		"markdown": true,
		"json":     true,
		"text":     true,
	}
	if !validFormats[format] {
		return fmt.Errorf("invalid format: %s (must be markdown, json, or text)", format)
	}

	opts := toc.Options{
		MaxDepth:          maxDepth,
		Format:            format,
		Recursive:         recursive,
		OutputFile:        outputFile,
		ConfirmThreshold:  confirmThreshold,
		SkipConfirmation:  skipConfirmation,
	}

	results, err := toc.ProcessPath(path, opts)
	if err != nil {
		return fmt.Errorf("error processing path: %w", err)
	}

	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "No markdown files found")
		return nil
	}

	output := toc.GenerateCombinedTOC(results, opts)

	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
			return fmt.Errorf("error writing output file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "TOC written to: %s\n", outputFile)
	} else {
		fmt.Print(output)
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
