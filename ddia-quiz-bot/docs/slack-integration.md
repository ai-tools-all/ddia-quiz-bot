# Slack Integration Guide

## Overview
The ddia-quiz-bot can send notifications to a Slack channel via incoming webhooks. This is optional and enabled only when `SLACK_WEBHOOK_URL` is configured.

## Steps

### 1. Create Slack Webhook
- Go to [Slack App Directory](https://api.slack.com/apps)
- Create a new app or use existing Incoming Webhooks app
- Add Incoming Webhook to your workspace
- Choose the channel for notifications
- Copy the Webhook URL (starts with `https://hooks.slack.com/services/...`)

### 2. Configure Environment
- Copy `.env.example` to `.env`
- Add your webhook URL:
  ```
  SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/UNIQUE/URL
  ```

### 3. Run the Bot
- Start the daemon: `go run ./cmd/quiz-daemon/main.go`
- Check logs for: `"Slack notifier has been enabled."`
- Posts will now be sent to both console and Slack

## Notes
- If `SLACK_WEBHOOK_URL` is empty or unset, Slack notifications are disabled
- Uses Slack Block Kit for formatting
- HTTP requests timeout after 10 seconds
- Errors are logged but don't stop the bot