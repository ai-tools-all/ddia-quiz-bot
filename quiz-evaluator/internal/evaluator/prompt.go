package evaluator

import (
	"fmt"
	"strings"

	"github.com/abhishek/quiz-evaluator/internal/models"
)

// PromptBuilder constructs evaluation prompts for AI
type PromptBuilder struct{}

// NewPromptBuilder creates a new prompt builder
func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{}
}

// BuildEvaluationPrompt creates a detailed evaluation prompt
func (pb *PromptBuilder) BuildEvaluationPrompt(question *models.Question, userResponse string) string {
	var prompt strings.Builder

	prompt.WriteString("You are evaluating a technical interview response. ")
	prompt.WriteString("Please provide a thorough and constructive evaluation.\n\n")

	// Question details
	prompt.WriteString("=== QUESTION DETAILS ===\n")
	if question.Title != "" {
		fmt.Fprintf(&prompt, "Title: %s\n", question.Title)
	}
	if question.Level != "" {
		fmt.Fprintf(&prompt, "Level: %s\n", question.Level)
	}
	if question.Category != "" {
		fmt.Fprintf(&prompt, "Category: %s\n", question.Category)
	}
	fmt.Fprintf(&prompt, "\nQuestion: %s\n", question.MainQuestion)

	// Evaluation criteria
	prompt.WriteString("\n=== EVALUATION CRITERIA ===\n")
	
	if len(question.CoreConcepts) > 0 {
		prompt.WriteString("\nCore Concepts (60% weight) - These must be covered well:\n")
		for _, concept := range question.CoreConcepts {
			fmt.Fprintf(&prompt, "• %s\n", concept)
		}
	}

	if len(question.PeripheralConcepts) > 0 {
		prompt.WriteString("\nPeripheral Concepts (40% weight) - Nice to have:\n")
		for _, concept := range question.PeripheralConcepts {
			fmt.Fprintf(&prompt, "• %s\n", concept)
		}
	}

	// Sample answers for reference
	if question.SampleExcellent != "" {
		prompt.WriteString("\n=== SAMPLE EXCELLENT ANSWER ===\n")
		prompt.WriteString(question.SampleExcellent)
		prompt.WriteString("\n")
	}

	if question.SampleAcceptable != "" {
		prompt.WriteString("\n=== SAMPLE ACCEPTABLE ANSWER ===\n")
		prompt.WriteString(question.SampleAcceptable)
		prompt.WriteString("\n")
	}

	// Rubric if available
	if len(question.EvaluationRubric) > 0 {
		prompt.WriteString("\n=== SPECIFIC RUBRIC ===\n")
		for criterion, description := range question.EvaluationRubric {
			fmt.Fprintf(&prompt, "%s: %s\n", criterion, description)
		}
	}

	// User response
	prompt.WriteString("\n=== USER RESPONSE TO EVALUATE ===\n")
	prompt.WriteString(userResponse)
	prompt.WriteString("\n")

	// Evaluation instructions
	prompt.WriteString("\n=== EVALUATION INSTRUCTIONS ===\n")
	prompt.WriteString("Please evaluate the user's response based on the following criteria:\n")
	prompt.WriteString("1. Coverage of core concepts (60% of score)\n")
	prompt.WriteString("2. Coverage of peripheral concepts (40% of score)\n")
	prompt.WriteString("3. Technical accuracy and correctness\n")
	prompt.WriteString("4. Clarity and structure of explanation\n")
	prompt.WriteString("5. Use of relevant examples or real-world applications\n")
	prompt.WriteString("6. Depth of understanding demonstrated\n\n")

	prompt.WriteString("Provide your evaluation in the following format:\n\n")
	prompt.WriteString("SCORE: [0-100]\n\n")
	prompt.WriteString("STRENGTHS:\n")
	prompt.WriteString("- [Key strength 1]\n")
	prompt.WriteString("- [Key strength 2]\n")
	prompt.WriteString("- [Additional strengths...]\n\n")
	prompt.WriteString("IMPROVEMENTS:\n")
	prompt.WriteString("- [Area for improvement 1]\n")
	prompt.WriteString("- [Area for improvement 2]\n")
	prompt.WriteString("- [Additional improvements...]\n\n")
	prompt.WriteString("FEEDBACK: [2-3 sentences of overall constructive feedback]\n\n")

	prompt.WriteString("Important guidelines:\n")
	prompt.WriteString("• Be specific and actionable in your feedback\n")
	prompt.WriteString("• Focus on the technical content, not writing style\n")
	prompt.WriteString("• Acknowledge what the candidate did well\n")
	prompt.WriteString("• Provide concrete suggestions for improvement\n")
	prompt.WriteString("• Consider the difficulty level when scoring\n")

	return prompt.String()
}

// BuildBatchPrompt creates a prompt for evaluating multiple responses
func (pb *PromptBuilder) BuildBatchPrompt(questions []*models.Question, responses []string) string {
	if len(questions) != len(responses) {
		return ""
	}

	var prompt strings.Builder
	prompt.WriteString("Evaluate the following technical interview responses. ")
	prompt.WriteString("Provide scores and brief feedback for each.\n\n")

	for i, question := range questions {
		fmt.Fprintf(&prompt, "=== QUESTION %d ===\n", i+1)
		fmt.Fprintf(&prompt, "ID: %s\n", question.ID)
		fmt.Fprintf(&prompt, "Question: %s\n", question.MainQuestion)
		fmt.Fprintf(&prompt, "Response: %s\n\n", responses[i])
	}

	prompt.WriteString("For each question, provide:\n")
	prompt.WriteString("1. Score (0-100)\n")
	prompt.WriteString("2. Brief feedback (1-2 sentences)\n")

	return prompt.String()
}
