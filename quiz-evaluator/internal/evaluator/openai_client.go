package evaluator

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient implements the AIClient interface for OpenAI
type OpenAIClient struct {
	client      *openai.Client
	config      AIClientConfig
	promptBuilder *PromptBuilder
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(config AIClientConfig) (*OpenAIClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	client := openai.NewClient(config.APIKey)
	
	if config.Model == "" {
		config.Model = "gpt-4-turbo-preview"
	}

	if config.MaxTokens == 0 {
		config.MaxTokens = 1500
	}

	if config.Temperature == 0 {
		config.Temperature = 0.3
	}

	return &OpenAIClient{
		client: client,
		config: config,
		promptBuilder: NewPromptBuilder(),
	}, nil
}

// GetProviderName returns the provider name
func (c *OpenAIClient) GetProviderName() string {
	return "OpenAI"
}

// EvaluateResponse evaluates a user response using OpenAI
func (c *OpenAIClient) EvaluateResponse(ctx context.Context, req EvaluationRequest) (*EvaluationResponse, error) {
	// Build the evaluation prompt
	prompt := c.promptBuilder.BuildEvaluationPrompt(req.Question, req.UserResponse)

	// Prepare the API request
	apiReq := openai.ChatCompletionRequest{
		Model: c.config.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an expert technical interviewer evaluating responses to technical questions. Provide detailed, constructive feedback.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: c.config.Temperature,
		MaxTokens:   c.config.MaxTokens,
	}

	// Add timeout to context if not already present
	if c.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.config.Timeout)
		defer cancel()
	}

	// Call OpenAI API with retries
	var resp openai.ChatCompletionResponse
	var err error
	
	for retry := 0; retry <= c.config.MaxRetries; retry++ {
		resp, err = c.client.CreateChatCompletion(ctx, apiReq)
		if err == nil {
			break
		}
		
		// Check if error is retryable
		if !isRetryableError(err) {
			return nil, fmt.Errorf("OpenAI API error: %w", err)
		}
		
		if retry < c.config.MaxRetries {
			// Exponential backoff
			backoff := time.Duration(1<<uint(retry)) * time.Second
			time.Sleep(backoff)
		}
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed after %d retries: %w", c.config.MaxRetries, err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Parse the AI response
	rawResponse := resp.Choices[0].Message.Content
	evaluation, err := c.parseEvaluationResponse(rawResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	evaluation.RawResponse = rawResponse
	evaluation.Timestamp = time.Now()

	return evaluation, nil
}

// parseEvaluationResponse parses the structured response from AI
func (c *OpenAIClient) parseEvaluationResponse(response string) (*EvaluationResponse, error) {
	eval := &EvaluationResponse{
		Strengths:    []string{},
		Improvements: []string{},
	}

	// Try to parse JSON response first
	if json.Valid([]byte(response)) {
		var jsonResp struct {
			Score        float64  `json:"score"`
			Strengths    []string `json:"strengths"`
			Improvements []string `json:"improvements"`
			Feedback     string   `json:"feedback"`
		}
		if err := json.Unmarshal([]byte(response), &jsonResp); err == nil {
			eval.Score = jsonResp.Score
			eval.Strengths = jsonResp.Strengths
			eval.Improvements = jsonResp.Improvements
			eval.Feedback = jsonResp.Feedback
			return eval, nil
		}
	}

	// Fall back to text parsing
	lines := strings.Split(response, "\n")
	currentSection := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse score
		if strings.HasPrefix(strings.ToUpper(line), "SCORE:") {
			scoreStr := strings.TrimSpace(strings.TrimPrefix(strings.ToUpper(line), "SCORE:"))
			scoreStr = strings.TrimSuffix(scoreStr, "/100")
			score, err := strconv.ParseFloat(scoreStr, 64)
			if err == nil {
				eval.Score = score
			}
			continue
		}

		// Detect sections
		upperLine := strings.ToUpper(line)
		if strings.HasPrefix(upperLine, "STRENGTHS:") {
			currentSection = "strengths"
			continue
		} else if strings.HasPrefix(upperLine, "IMPROVEMENTS:") || strings.HasPrefix(upperLine, "AREAS FOR IMPROVEMENT:") {
			currentSection = "improvements"
			continue
		} else if strings.HasPrefix(upperLine, "FEEDBACK:") || strings.HasPrefix(upperLine, "OVERALL FEEDBACK:") {
			currentSection = "feedback"
			feedbackContent := strings.TrimSpace(strings.TrimPrefix(line, "FEEDBACK:"))
			if feedbackContent != "" {
				eval.Feedback = feedbackContent
			}
			continue
		}

		// Parse content based on current section
		switch currentSection {
		case "strengths":
			if item := parseBulletItem(line); item != "" {
				eval.Strengths = append(eval.Strengths, item)
			}
		case "improvements":
			if item := parseBulletItem(line); item != "" {
				eval.Improvements = append(eval.Improvements, item)
			}
		case "feedback":
			if eval.Feedback != "" {
				eval.Feedback += " "
			}
			eval.Feedback += line
		}
	}

	// Validate parsed response
	if eval.Score < 0 || eval.Score > 100 {
		return nil, fmt.Errorf("invalid score: %.2f", eval.Score)
	}

	if eval.Feedback == "" {
		eval.Feedback = "Evaluation completed."
	}

	return eval, nil
}

// parseBulletItem extracts content from a bullet point line
func parseBulletItem(line string) string {
	// Match various bullet formats
	patterns := []string{
		`^[-*•]\s+(.+)$`,
		`^\d+\.\s+(.+)$`,
		`^\d+\)\s+(.+)$`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(line); len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	// If no bullet format matched but line is not empty, return as is
	// This handles cases where AI doesn't use bullet points
	if line != "" && !strings.HasSuffix(line, ":") {
		return line
	}

	return ""
}

// isRetryableError determines if an error is worth retrying
func isRetryableError(err error) bool {
	errStr := err.Error()
	retryableErrors := []string{
		"rate limit",
		"timeout",
		"connection",
		"temporary",
		"503",
		"504",
		"429",
	}

	for _, keyword := range retryableErrors {
		if strings.Contains(strings.ToLower(errStr), keyword) {
			return true
		}
	}

	return false
}
