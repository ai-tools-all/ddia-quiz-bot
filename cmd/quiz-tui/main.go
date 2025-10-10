package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/abhishek/ddia-clicker/internal/config"
	"github.com/abhishek/ddia-clicker/internal/tui/screens"
)

var (
	userName   string
	configPath string
)

var rootCmd = &cobra.Command{
	Use:   "quiz-tui",
	Short: "Interactive TUI for subjective quiz questions",
	Long: `Quiz TUI is an interactive terminal user interface for taking subjective quizzes.
It allows you to progress through questions, save your answers, and resume sessions.`,
	RunE: runTUI,
}

func init() {
	rootCmd.Flags().StringVarP(&userName, "user", "u", "", "Username (required)")
	rootCmd.MarkFlagRequired("user")
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to config file (default: config/tui.toml)")
}

func runTUI(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.LoadTUIConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create initial model
	model := screens.NewImprovedAppModel(userName, cfg)

	// Initialize the program
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
