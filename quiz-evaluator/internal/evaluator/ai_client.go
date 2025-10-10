package evaluator

import (
	"context"
	"fmt"
	"time"

	"github.com/abhishek/quiz-evaluator/internal/models"
)

// AIProvider represents different AI service providers
type AIProvider string

const (
	ProviderOpenAI   AIProvider = "openai"
	ProviderAnthropic AIProvider = "anthropic"
	ProviderGemini   AIProvider = "gemini"
	ProviderOllama   AIProvider = "ollama"
)

// AIClient interface for AI evaluation services
type AIClient interface {
	EvaluateResponse(ctx context.Context, req EvaluationRequest) (*EvaluationResponse, error)
	GetProviderName() string
}

// EvaluationRequest contains the data needed for evaluation
type EvaluationRequest struct {
	Question     *models.Question
	UserResponse string
	MaxTokens    int
	Temperature  float32
}

// EvaluationResponse contains the AI evaluation result
type EvaluationResponse struct {
	Score        float64
	Strengths    []string
	Improvements []string
	Feedback     string
	RawResponse  string // For debugging
	Timestamp    time.Time
}

// AIClientConfig contains configuration for AI clients
type AIClientConfig struct {
	Provider    AIProvider
	APIKey      string
	Model       string
	MaxTokens   int
	Temperature float32
	Timeout     time.Duration
	MaxRetries  int
}

// NewAIClient creates an appropriate AI client based on provider
func NewAIClient(config AIClientConfig) (AIClient, error) {
	switch config.Provider {
	case ProviderOpenAI:
		return NewOpenAIClient(config)
	case ProviderAnthropic:
		return nil, fmt.Errorf("anthropic provider not yet implemented")
	case ProviderGemini:
		return nil, fmt.Errorf("gemini provider not yet implemented")
	case ProviderOllama:
		return nil, fmt.Errorf("ollama provider not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", config.Provider)
	}
}
