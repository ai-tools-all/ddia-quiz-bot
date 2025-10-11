package models

// Question represents a quiz question with evaluation criteria
type Question struct {
	ID                 string            `toml:"id" json:"id"`
	Title              string            `toml:"title" json:"title"`
	MainQuestion       string            `toml:"question" json:"question"`
	CoreConcepts       []string          `toml:"core_concepts" json:"core_concepts"`
	PeripheralConcepts []string          `toml:"peripheral_concepts" json:"peripheral_concepts"`
	SampleExcellent    string            `toml:"sample_excellent" json:"sample_excellent"`
	SampleAcceptable   string            `toml:"sample_acceptable" json:"sample_acceptable"`
	Level              string            `toml:"level" json:"level"`
	Category           string            `toml:"category" json:"category"`
	EvaluationRubric   map[string]string `toml:"evaluation_rubric" json:"evaluation_rubric"`
	FilePath           string            `toml:"-" json:"-"` // Source file path
	
	// MCQ specific fields
	Type        string   `toml:"type" json:"type"`                   // "subjective" or "mcq"
	Options     []string `toml:"options" json:"options"`             // MCQ options
	Answer      string   `toml:"answer" json:"answer"`               // Correct answer (A, B, C, D)
	Explanation string   `toml:"explanation" json:"explanation"`     // MCQ explanation
	Hook        string   `toml:"hook" json:"hook"`                   // Engagement hook
}

// QuestionIndex maps question IDs to Question objects
type QuestionIndex map[string]*Question
