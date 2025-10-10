package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/abhishek/quiz-evaluator/internal/models"
)

// Reader handles CSV input operations
type Reader struct {
	filepath string
}

// NewReader creates a new CSV reader
func NewReader(filepath string) *Reader {
	return &Reader{filepath: filepath}
}

// ReadResponses reads user responses from CSV file
func (r *Reader) ReadResponses() ([]models.UserResponse, error) {
	file, err := os.Open(r.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	
	// Read header
	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Validate header
	if err := validateHeader(header); err != nil {
		return nil, err
	}

	// Find column indices
	questionIDIdx := findColumnIndex(header, "question_id")
	responseIdx := findColumnIndex(header, "user_response")

	if questionIDIdx == -1 || responseIdx == -1 {
		return nil, fmt.Errorf("required columns 'question_id' and 'user_response' not found")
	}

	var responses []models.UserResponse
	lineNum := 1

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading line %d: %w", lineNum+1, err)
		}
		lineNum++

		if len(record) <= questionIDIdx || len(record) <= responseIdx {
			return nil, fmt.Errorf("invalid record at line %d: insufficient columns", lineNum)
		}

		// Skip empty responses
		if record[responseIdx] == "" {
			continue
		}

		response := models.UserResponse{
			QuestionID:   record[questionIDIdx],
			UserResponse: record[responseIdx],
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func validateHeader(header []string) error {
	if len(header) < 2 {
		return fmt.Errorf("CSV must have at least 2 columns")
	}
	return nil
}

func findColumnIndex(header []string, columnName string) int {
	for i, name := range header {
		if name == columnName {
			return i
		}
	}
	return -1
}
