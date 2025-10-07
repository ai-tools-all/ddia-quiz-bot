package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/your-username/ddia-quiz-bot/internal/models"
	"gopkg.in/yaml.v3"
)

// LoadSchedule reads and parses the schedule.yml file with environment variable substitution.
func LoadSchedule(path string) (*models.Schedule, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Replace environment variables in the YAML content
	expandedData := expandEnvVars(string(data))

	var schedule models.Schedule
	if err := yaml.Unmarshal([]byte(expandedData), &schedule); err != nil {
		return nil, err
	}

	return &schedule, nil
}

// expandEnvVars replaces ${VAR} or $VAR with environment variable values
func expandEnvVars(content string) string {
	return os.Expand(content, func(key string) string {
		// Handle ${VAR} and $VAR syntax
		if strings.HasPrefix(key, "{") && strings.HasSuffix(key, "}") {
			key = key[1 : len(key)-1]
		}
		return os.Getenv(key)
	})
}
