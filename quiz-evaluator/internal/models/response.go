package models

// UserResponse represents a user's answer to a question
type UserResponse struct {
	QuestionID   string `csv:"question_id" json:"question_id"`
	UserResponse string `csv:"user_response" json:"user_response"`
}

// InputRecord represents a row in the input CSV
type InputRecord struct {
	QuestionID   string `csv:"question_id"`
	UserResponse string `csv:"user_response"`
}
