package models

import "time"

// Evaluation represents the AI evaluation result for a user response
type Evaluation struct {
	QuestionID    string    `json:"question_id" csv:"question_id"`
	UserResponse  string    `json:"user_response" csv:"user_response"`
	Score         float64   `json:"score" csv:"score"`
	Feedback      string    `json:"feedback" csv:"feedback"`
	QuestionTitle string    `json:"question_title" csv:"question_title"`
	Level         string    `json:"level" csv:"level"`
	Strengths     []string  `json:"strengths"`
	Improvements  []string  `json:"improvements"`
	Timestamp     time.Time `json:"timestamp"`
}

// OutputRecord represents a row in the output CSV
type OutputRecord struct {
	QuestionID    string  `csv:"question_id"`
	UserResponse  string  `csv:"user_response"`
	Score         float64 `csv:"score"`
	Feedback      string  `csv:"feedback"`
	QuestionTitle string  `csv:"question_title"`
	Level         string  `csv:"level"`
	Strengths     string  `csv:"strengths"`    // Comma-separated
	Improvements  string  `csv:"improvements"` // Comma-separated
}
