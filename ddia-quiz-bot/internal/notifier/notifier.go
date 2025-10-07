package notifier

import "github.com/sirupsen/logrus"

// Notifier is the interface for sending a formatted post to a platform.
type Notifier interface {
	Notify(content string) error
}

// LogNotifier is a simple implementation that prints to a logger.
type LogNotifier struct {
	Logger *logrus.Logger
}

func (n *LogNotifier) Notify(content string) error {
	n.Logger.Info("--- NEW POST ---")
	n.Logger.Info(content)
	n.Logger.Info("--- END POST ---")
	return nil
}
