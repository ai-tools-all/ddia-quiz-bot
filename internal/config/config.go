package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	AI         AIConfig         `mapstructure:"ai"`
	Evaluation EvaluationConfig `mapstructure:"evaluation"`
	Output     OutputConfig     `mapstructure:"output"`
}

// AIConfig holds AI provider configuration
type AIConfig struct {
	Provider    string        `mapstructure:"provider"`
	APIKey      string        `mapstructure:"api_key"`
	Model       string        `mapstructure:"model"`
	MaxTokens   int           `mapstructure:"max_tokens"`
	Temperature float32       `mapstructure:"temperature"`
	Timeout     time.Duration `mapstructure:"timeout"`
	MaxRetries  int           `mapstructure:"max_retries"`
}

// EvaluationConfig holds evaluation settings
type EvaluationConfig struct {
	ParallelWorkers int           `mapstructure:"parallel_workers"`
	RetryAttempts   int           `mapstructure:"retry_attempts"`
	RetryDelay      time.Duration `mapstructure:"retry_delay"`
}

// OutputConfig holds output settings
type OutputConfig struct {
	Format          string `mapstructure:"format"`
	IncludeMetadata bool   `mapstructure:"include_metadata"`
	Verbose         bool   `mapstructure:"verbose"`
	OutputFile      string `mapstructure:"output_file"`
}

// Load loads configuration from file and environment
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Set config file if provided
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// Look for config in standard locations
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("$HOME/.quiz-evaluator")
	}

	// Read config file if it exists
	if err := v.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist, we'll use defaults and env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Read environment variables
	v.SetEnvPrefix("QUIZ_EVAL")
	v.AutomaticEnv()

	// Override with environment variables
	bindEnvironmentVariables(v)

	// Unmarshal configuration
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// AI defaults
	v.SetDefault("ai.provider", "openai")
	v.SetDefault("ai.model", "gpt-4-turbo-preview")
	v.SetDefault("ai.max_tokens", 1500)
	v.SetDefault("ai.temperature", 0.3)
	v.SetDefault("ai.timeout", "30s")
	v.SetDefault("ai.max_retries", 3)

	// Evaluation defaults
	v.SetDefault("evaluation.parallel_workers", 5)
	v.SetDefault("evaluation.retry_attempts", 3)
	v.SetDefault("evaluation.retry_delay", "2s")

	// Output defaults
	v.SetDefault("output.format", "csv")
	v.SetDefault("output.include_metadata", true)
	v.SetDefault("output.verbose", false)
	v.SetDefault("output.output_file", "evaluations.csv")
}

// bindEnvironmentVariables binds specific environment variables
func bindEnvironmentVariables(v *viper.Viper) {
	// Allow API key from environment
	v.BindEnv("ai.api_key", "OPENAI_API_KEY", "AI_API_KEY", "QUIZ_EVAL_API_KEY")

	// Allow provider override
	v.BindEnv("ai.provider", "AI_PROVIDER", "QUIZ_EVAL_PROVIDER")

	// Allow model override
	v.BindEnv("ai.model", "AI_MODEL", "QUIZ_EVAL_MODEL")
}

// validateConfig validates the configuration
func validateConfig(cfg *Config) error {
	// Check API key for cloud providers
	if cfg.AI.Provider != "ollama" && cfg.AI.APIKey == "" {
		// Try to get from environment
		for _, envVar := range []string{"OPENAI_API_KEY", "AI_API_KEY", "ANTHROPIC_API_KEY"} {
			if key := os.Getenv(envVar); key != "" {
				cfg.AI.APIKey = key
				break
			}
		}

		if cfg.AI.APIKey == "" {
			return fmt.Errorf("API key is required for provider %s", cfg.AI.Provider)
		}
	}

	// Validate provider
	validProviders := map[string]bool{
		"openai":    true,
		"anthropic": true,
		"gemini":    true,
		"ollama":    true,
	}

	if !validProviders[cfg.AI.Provider] {
		return fmt.Errorf("unsupported AI provider: %s", cfg.AI.Provider)
	}

	// Validate output format
	validFormats := map[string]bool{
		"csv":      true,
		"json":     true,
		"markdown": true,
	}

	if !validFormats[cfg.Output.Format] {
		return fmt.Errorf("unsupported output format: %s", cfg.Output.Format)
	}

	// Ensure positive values
	if cfg.Evaluation.ParallelWorkers < 1 {
		cfg.Evaluation.ParallelWorkers = 1
	}

	if cfg.AI.MaxTokens < 100 {
		cfg.AI.MaxTokens = 1500
	}

	return nil
}

// SaveExample saves an example configuration file
func SaveExample(filepath string) error {
	exampleConfig := `# Quiz Evaluator Configuration

# AI Provider Settings
ai:
  provider: "openai"  # Options: openai, anthropic, gemini, ollama
  api_key: "${OPENAI_API_KEY}"  # Can also use environment variable
  model: "gpt-4-turbo-preview"
  max_tokens: 1500
  temperature: 0.3
  timeout: "30s"
  max_retries: 3

# Evaluation Settings
evaluation:
  parallel_workers: 5
  retry_attempts: 3
  retry_delay: "2s"

# Output Settings
output:
  format: "csv"  # Options: csv, json, markdown
  include_metadata: true
  verbose: false
  output_file: "evaluations.csv"
`

	return os.WriteFile(filepath, []byte(exampleConfig), 0644)
}
