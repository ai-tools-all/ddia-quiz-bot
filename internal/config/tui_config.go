package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// TUIConfig represents the TUI application configuration
type TUIConfig struct {
	AutoSaveInterval time.Duration `mapstructure:"auto_save_interval"`
	SessionsDir      string        `mapstructure:"sessions_dir"`
	ContentPath      string        `mapstructure:"content_path"`       // Deprecated: use ChaptersRootPath
	ChaptersRootPath string        `mapstructure:"chapters_root_path"` // Path to chapters directory
	DefaultMode      string        `mapstructure:"default_mode"`       // Default question mode: "mcq", "subjective", or "mixed"
}

// LoadTUIConfig loads TUI configuration from file
func LoadTUIConfig(configPath string) (*TUIConfig, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("auto_save_interval", 30)
	v.SetDefault("sessions_dir", "sessions")
	v.SetDefault("content_path", "ddia-quiz-bot/content/chapters/09-distributed-systems-gfs/subjective")
	v.SetDefault("chapters_root_path", "ddia-quiz-bot/content/chapters")
	v.SetDefault("default_mode", "mixed")

	// Set config file if provided
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// Look for config in standard locations
		v.SetConfigName("tui")
		v.SetConfigType("toml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("$HOME/.quiz-tui")
	}

	// Read config file if it exists
	if err := v.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist, we'll use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Read environment variables
	v.SetEnvPrefix("QUIZ_TUI")
	v.AutomaticEnv()

	var config TUIConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Convert auto_save_interval from seconds to duration
	if seconds := v.GetInt("auto_save_interval"); seconds > 0 {
		config.AutoSaveInterval = time.Duration(seconds) * time.Second
	}

	// Backward compatibility: if ChaptersRootPath is not set but ContentPath is,
	// keep ContentPath for single-topic mode
	if config.ChaptersRootPath == "" && config.ContentPath != "" {
		// Use old behavior - ContentPath points to specific topic
		config.ChaptersRootPath = ""
	}

	return &config, nil
}
