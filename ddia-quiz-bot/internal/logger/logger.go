package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// SetupLogger configures logrus with dual writers (stdout + file) and appropriate formatting
func SetupLogger(logLevel string, logDir string) (*logrus.Logger, error) {
	logger := logrus.New()

	// Parse log level
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel // default to info
	}
	logger.SetLevel(level)

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFile := filepath.Join(logDir, "ddia-quiz-bot_"+timestamp+".log")

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// Create dual writer
	multiWriter := io.MultiWriter(os.Stdout, file)

	// Set up different formatters for stdout and file
	logger.SetOutput(multiWriter)

	// For stdout, use text formatter with colors
	consoleFormatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	}

	// For file, use JSON formatter
	jsonFormatter := &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	}

	// Use a custom hook to write different formats to different outputs
	logger.AddHook(&DualFormatHook{
		Writer:           file,
		JSONFormatter:    jsonFormatter,
		ConsoleFormatter: consoleFormatter,
	})

	// Set console formatter for stdout
	logger.SetFormatter(consoleFormatter)

	return logger, nil
}

// DualFormatHook writes JSON to file and text to console
type DualFormatHook struct {
	Writer           io.Writer
	JSONFormatter    logrus.Formatter
	ConsoleFormatter logrus.Formatter
}

func (hook *DualFormatHook) Fire(entry *logrus.Entry) error {
	// Format for file (JSON)
	jsonBytes, err := hook.JSONFormatter.Format(entry)
	if err != nil {
		return err
	}

	// Write JSON to file
	_, err = hook.Writer.Write(jsonBytes)
	return err
}

func (hook *DualFormatHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
