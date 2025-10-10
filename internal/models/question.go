package models

// Question represents a quiz question with evaluation criteria
type Question struct {
	ID                 string            `yaml:"id" json:"id"`
	Title              string            `yaml:"title" json:"title"`
	MainQuestion       string            `yaml:"question" json:"question"`
	CoreConcepts       []string          `yaml:"core_concepts" json:"core_concepts"`
	PeripheralConcepts []string          `yaml:"peripheral_concepts" json:"peripheral_concepts"`
	SampleExcellent    string            `yaml:"sample_excellent" json:"sample_excellent"`
	SampleAcceptable   string            `yaml:"sample_acceptable" json:"sample_acceptable"`
	Level              string            `yaml:"level" json:"level"`
	Category           string            `yaml:"category" json:"category"`
	EvaluationRubric   map[string]string `yaml:"evaluation_rubric" json:"evaluation_rubric"`
	FilePath           string            `yaml:"-" json:"-"` // Source file path
}

// QuestionIndex maps question IDs to Question objects
type QuestionIndex map[string]*Question
