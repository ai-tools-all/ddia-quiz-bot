package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// SlackNotifier sends notifications to a Slack channel via an incoming webhook.
type SlackNotifier struct {
	WebhookURL string
	Logger     *logrus.Logger
}

// slackPayload defines the JSON structure for a Slack message using Block Kit.
// This allows for better formatting than a simple text message.
type slackPayload struct {
	Blocks []slackBlock `json:"blocks"`
}

type slackBlock struct {
	Type string    `json:"type"`
	Text slackText `json:"text"`
}

type slackText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// NewSlackNotifier creates a new, configured instance of the Slack notifier.
// It returns the notifier if the webhookURL is provided, otherwise it returns nil.
func NewSlackNotifier(webhookURL string, logger *logrus.Logger) *SlackNotifier {
	if webhookURL == "" {
		logger.Debug("SLACK_WEBHOOK_URL not set, Slack notifier is disabled.")
		return nil
	}
	return &SlackNotifier{
		WebhookURL: webhookURL,
		Logger:     logger,
	}
}

// Notify sends the formatted content to the configured Slack webhook URL.
// It implements the Notifier interface.
func (n *SlackNotifier) Notify(content string) error {
	// Slack's Block Kit uses a markdown-like syntax called "mrkdwn".
	// The output from our presenter is generally compatible.
	payload := slackPayload{
		Blocks: []slackBlock{
			{
				Type: "section",
				Text: slackText{
					Type: "mrkdwn",
					Text: content,
				},
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		n.Logger.WithError(err).Error("Failed to marshal Slack payload")
		return fmt.Errorf("failed to marshal slack payload: %w", err)
	}

	req, err := http.NewRequest("POST", n.WebhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		n.Logger.WithError(err).Error("Failed to create Slack HTTP request")
		return fmt.Errorf("failed to create slack http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Use a client with a timeout to prevent the application from hanging.
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		n.Logger.WithError(err).Error("Failed to send notification to Slack")
		return fmt.Errorf("failed to send notification to slack: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful response from Slack.
	if resp.StatusCode != http.StatusOK {
		n.Logger.WithField("status_code", resp.StatusCode).Error("Slack returned a non-200 status code")
		return fmt.Errorf("slack returned a non-ok status: %d", resp.StatusCode)
	}

	n.Logger.Info("Successfully sent notification to Slack.")
	return nil
}
