# 001-2025-10-07-123826-feature-implement-slack-notifier

## Overview
Implement Slack notifier for ddia-quiz-bot to send notifications to Slack channel via webhook.

## Tasks
- [x] Create slack.go in internal/notifier/
- [x] Update main.go to integrate SlackNotifier
- [x] Create/update .env.example with SLACK_WEBHOOK_URL
- [x] Build and test the project

## Notes
- SlackNotifier implements Notifier interface
- Uses Block Kit for formatting
- Conditionally enabled via SLACK_WEBHOOK_URL env var
- Timeout of 10s for HTTP requests