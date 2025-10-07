package notifier

import "log"

// Notifier is the interface for sending a formatted post to a platform.
type Notifier interface {
	Notify(content string) error
}

// LogNotifier is a simple implementation that prints to a logger.
type LogNotifier struct {
	Logger *log.Logger
}

func (n *LogNotifier) Notify(content string) error {
	n.Logger.Println("--- NEW POST ---")
	n.Logger.Println(content)
	n.Logger.Println("--- END POST ---")
	return nil
}
