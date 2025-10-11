package markdown

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseQuestionFile(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		wantID       string
		wantTitle    string
		wantLevel    string
		wantQuestion string
		wantErr      bool
	}{
		{
			name: "complete question with frontmatter",
			content: `---
id: test-001
type: subjective
level: L3
category: baseline
topic: gfs
---

# question_title - Test Question

## main_question - Core Question
This is the main question text.

## core_concepts - Must Mention
- Concept 1
- Concept 2

## sample_excellent - Example Excellence
This is an excellent answer.
`,
			wantID:       "test-001",
			wantLevel:    "L3",
			wantQuestion: "This is the main question text.",
			wantErr:      false,
		},
		{
			name: "question with simple section headers",
			content: `---
id: test-002
level: L4
---

## main question
Simple question text here.

## core concepts
- First concept
- Second concept
`,
			wantID:       "test-002",
			wantLevel:    "L4",
			wantQuestion: "Simple question text here.",
			wantErr:      false,
		},
		{
			name: "question without frontmatter",
			content: `
question_id: test-003

## question
What is the answer?

## core concepts
- Key concept
`,
			wantID:       "test-003",
			wantQuestion: "What is the answer?",
			wantErr:      false,
		},
		{
			name: "actual GFS question format",
			content: `---
id: gfs-subjective-L3-001
type: subjective
level: L3
category: baseline
topic: gfs
subtopic: replication
estimated_time: 5-7 minutes
---

# question_title - GFS Replication Strategy

## main_question - Core Question
"Explain how GFS ensures data durability when a chunk server fails."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **3x Replication**: GFS maintains 3 copies
- **Master Detection**: Master detects failure

## peripheral_concepts - Nice to Have (40%)
- **Rack Awareness**: Replicas spread across racks

## sample_excellent - Example Excellence
This is an excellent answer with details.

## sample_acceptable - Minimum Acceptable
This is a minimal acceptable answer.
`,
			wantID:       "gfs-subjective-L3-001",
			wantLevel:    "L3",
			wantQuestion: "\"Explain how GFS ensures data durability when a chunk server fails.\"",
			wantErr:      false,
		},
		{
			name: "missing id should error",
			content: `---
level: L3
---

## main_question
Question without ID.
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.md")
			err := os.WriteFile(tmpFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Parse the file
			parser := NewParser()
			question, err := parser.ParseQuestionFile(tmpFile)

			// Check error expectation
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify parsed data
			if question.ID != tt.wantID {
				t.Errorf("ID: got %q, want %q", question.ID, tt.wantID)
			}

			if tt.wantLevel != "" && question.Level != tt.wantLevel {
				t.Errorf("Level: got %q, want %q", question.Level, tt.wantLevel)
			}

			if tt.wantQuestion != "" && question.MainQuestion != tt.wantQuestion {
				t.Errorf("MainQuestion: got %q, want %q", question.MainQuestion, tt.wantQuestion)
			}

			// Check that MainQuestion is not empty when we expect it
			if tt.wantQuestion != "" && question.MainQuestion == "" {
				t.Error("MainQuestion is empty but should have content")
			}
		})
	}
}

func TestExtractFrontmatter(t *testing.T) {
	tests := []struct {
		name            string
		content         string
		wantFrontmatter string
		wantBody        string
	}{
		{
			name: "valid frontmatter",
			content: `---
id: test-001
level: L3
---

Body content here.
`,
			wantFrontmatter: "id: test-001\nlevel: L3",
			wantBody:        "\n\nBody content here.\n",
		},
		{
			name: "no frontmatter",
			content: `Just body content.
No frontmatter here.
`,
			wantFrontmatter: "",
			wantBody:        "Just body content.\nNo frontmatter here.\n",
		},
		{
			name: "frontmatter without closing",
			content: `---
id: test
Body without closing marker
`,
			wantFrontmatter: "",
			wantBody:        "---\nid: test\nBody without closing marker\n",
		},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFrontmatter, gotBody, _ := parser.extractFrontmatter(tt.content)

			if gotFrontmatter != tt.wantFrontmatter {
				t.Errorf("Frontmatter:\ngot:  %q\nwant: %q", gotFrontmatter, tt.wantFrontmatter)
			}

			if gotBody != tt.wantBody {
				t.Errorf("Body:\ngot:  %q\nwant: %q", gotBody, tt.wantBody)
			}
		})
	}
}

func TestParseListSection(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name: "bullet list with dashes",
			content: `- First item
- Second item
- Third item`,
			want: []string{"First item", "Second item", "Third item"},
		},
		{
			name: "bullet list with asterisks",
			content: `* First item
* Second item`,
			want: []string{"First item", "Second item"},
		},
		{
			name: "numbered list",
			content: `1. First item
2. Second item
3. Third item`,
			want: []string{"First item", "Second item", "Third item"},
		},
		{
			name: "nested markdown in list",
			content: `- **Bold concept**: Description
- Another concept with text`,
			want: []string{"**Bold concept**: Description", "Another concept with text"},
		},
		{
			name:    "empty content",
			content: "",
			want:    nil,
		},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.parseListSection(tt.content)

			if len(got) != len(tt.want) {
				t.Errorf("Length mismatch: got %d items, want %d items", len(got), len(tt.want))
				t.Logf("Got: %v", got)
				t.Logf("Want: %v", tt.want)
				return
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Item %d: got %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestParseSectionWithSuffix(t *testing.T) {
	// Test that sections with suffixes like "main_question - Core Question" are handled
	content := `---
id: test-suffix-001
level: L3
---

## main_question - Core Question
This is the actual question text.

## core_concepts - Must Mention (60%)
- Concept one
- Concept two

## sample_excellent - Example Excellence  
This is an excellent answer.
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.md")
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	parser := NewParser()
	question, err := parser.ParseQuestionFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to parse question: %v", err)
	}

	if question.MainQuestion == "" {
		t.Error("MainQuestion is empty - parser didn't handle section suffix")
	}

	if question.MainQuestion != "This is the actual question text." {
		t.Errorf("MainQuestion: got %q, want %q",
			question.MainQuestion,
			"This is the actual question text.")
	}

	if len(question.CoreConcepts) != 2 {
		t.Errorf("CoreConcepts: got %d items, want 2", len(question.CoreConcepts))
	}

	if question.SampleExcellent == "" {
		t.Error("SampleExcellent is empty - parser didn't handle section suffix")
	}
}

func TestParseRealQuestionFile(t *testing.T) {
	// Test with an actual question file if it exists
	questionPath := filepath.Join("..", "..", "ddia-quiz-bot", "content", "chapters",
		"09-distributed-systems-gfs", "subjective", "L3-baseline", "01-replication-basics.md")

	if _, err := os.Stat(questionPath); os.IsNotExist(err) {
		t.Skip("Skipping real file test - file not found")
	}

	parser := NewParser()
	question, err := parser.ParseQuestionFile(questionPath)
	if err != nil {
		t.Fatalf("Failed to parse real question file: %v", err)
	}

	// Verify essential fields
	if question.ID == "" {
		t.Error("ID is empty")
	}

	if question.Level == "" {
		t.Error("Level is empty")
	}

	if question.MainQuestion == "" {
		t.Error("MainQuestion is empty - this is the bug!")
		t.Logf("Question struct: %+v", question)
	} else {
		t.Logf("Successfully parsed MainQuestion: %q", question.MainQuestion)
	}

	if len(question.CoreConcepts) == 0 {
		t.Log("Warning: CoreConcepts is empty")
	}
}

// TestParseMCQWithTOML tests parsing MCQ questions with TOML frontmatter
func TestParseMCQWithTOML(t *testing.T) {
	tests := []struct {
		name            string
		content         string
		wantType        string
		wantOptions     int
		wantAnswer      string
		wantExplanation bool
		wantHook        bool
		wantError       bool
	}{
		{
			name: "complete MCQ with TOML frontmatter",
			content: `+++
id = "test-mcq-001"
title = "Test MCQ"
level = "L3"
category = "baseline"
type = "mcq"
+++

## Question

What is the capital of France?

## Options

- A) London
- B) Paris
- C) Berlin
- D) Madrid

## Answer

B

## Explanation

Paris is the capital and largest city of France.

## Hook

Geography basics matter!

## Core Concepts

- European capitals
- Geography

## Peripheral Concepts

- French culture
- European cities
`,
			wantType:        "mcq",
			wantOptions:     4,
			wantAnswer:      "B",
			wantExplanation: true,
			wantHook:        true,
		},
		{
			name: "MCQ with YAML frontmatter (backward compatibility)",
			content: `---
id: test-mcq-002
title: Test MCQ YAML
level: L4
category: baseline
type: mcq
---

## Question

Which is a programming language?

## Options

- A) HTML
- B) CSS
- C) Python
- D) JSON

## Answer

C

## Explanation

Python is a general-purpose programming language.
`,
			wantType:        "mcq",
			wantOptions:     4,
			wantAnswer:      "C",
			wantExplanation: true,
			wantHook:        false,
		},
		{
			name: "MCQ with various option formats",
			content: `+++
id = "test-mcq-003"
type = "mcq"
+++

## Question

Test question

## Options

A) First option
* B) Second option
- C) Third option
  D) Fourth option

## Answer

A
`,
			wantType:    "mcq",
			wantOptions: 4,
			wantAnswer:  "A",
		},
		{
			name: "MCQ missing answer field",
			content: `+++
id = "test-mcq-004"
type = "mcq"
+++

## Question

Test question

## Options

- A) Option 1
- B) Option 2
`,
			wantType:    "mcq",
			wantOptions: 2,
			wantAnswer:  "", // Should handle gracefully
		},
		{
			name: "Mixed MCQ and subjective fields",
			content: `+++
id = "test-mcq-005"
type = "mcq"
+++

## Question

Test question

## Options

- A) Option A
- B) Option B

## Answer

A

## Explanation

Test explanation

## Sample Excellent Answer

This is ignored for MCQ

## Core Concepts

- Concept 1
- Concept 2
`,
			wantType:    "mcq",
			wantOptions: 2,
			wantAnswer:  "A",
		},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmpFile, err := os.CreateTemp("", "mcq-test-*.md")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(tt.content); err != nil {
				t.Fatalf("Failed to write temp file: %v", err)
			}
			tmpFile.Close()

			// Parse the file
			question, err := parser.ParseQuestionFile(tmpFile.Name())

			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify type
			if question.Type != tt.wantType {
				t.Errorf("Type: got %q, want %q", question.Type, tt.wantType)
			}

			// Verify options
			if len(question.Options) != tt.wantOptions {
				t.Errorf("Options count: got %d, want %d", len(question.Options), tt.wantOptions)
				t.Logf("Options: %v", question.Options)
			}

			// Verify answer
			if question.Answer != tt.wantAnswer {
				t.Errorf("Answer: got %q, want %q", question.Answer, tt.wantAnswer)
			}

			// Verify explanation
			if tt.wantExplanation && question.Explanation == "" {
				t.Error("Expected explanation but got empty string")
			}

			// Verify hook
			if tt.wantHook && question.Hook == "" {
				t.Error("Expected hook but got empty string")
			}

			// Verify question text is not empty
			if question.MainQuestion == "" {
				t.Error("MainQuestion should not be empty")
			}
		})
	}
}

// TestParseMCQOptions tests the parseMCQOptions function with various formats
func TestParseMCQOptions(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name: "standard dash format",
			content: `- A) First option
- B) Second option
- C) Third option
- D) Fourth option`,
			want: []string{"A) First option", "B) Second option", "C) Third option", "D) Fourth option"},
		},
		{
			name: "asterisk format",
			content: `* A) First option
* B) Second option`,
			want: []string{"A) First option", "B) Second option"},
		},
		{
			name: "no bullet format",
			content: `A) First option
B) Second option
C) Third option`,
			want: []string{"A) First option", "B) Second option", "C) Third option"},
		},
		{
			name: "mixed spacing",
			content: `  - A) Indented option
-B) No space option
  C)   Extra spaces`,
			want: []string{"A) Indented option", "B) No space option", "C) Extra spaces"},
		},
		{
			name: "lowercase letters (should convert to uppercase)",
			content: `- a) First option
- b) Second option`,
			want: []string{"A) First option", "B) Second option"},
		},
		{
			name: "empty lines and extra content",
			content: `- A) First option

- B) Second option

Some extra text here
- C) Third option`,
			want: []string{"A) First option", "B) Second option", "C) Third option"},
		},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.parseMCQOptions(tt.content)

			if len(got) != len(tt.want) {
				t.Errorf("Length mismatch: got %d, want %d", len(got), len(tt.want))
				t.Logf("Got: %v", got)
				t.Logf("Want: %v", tt.want)
				return
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Option %d: got %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestTOMLFrontmatterExtraction tests TOML vs YAML frontmatter detection
func TestTOMLFrontmatterExtraction(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantTOML bool
		wantFM   string
		wantBody string
	}{
		{
			name: "TOML frontmatter",
			content: `+++
id = "test"
title = "Test"
+++

Body content`,
			wantTOML: true,
			wantFM:   "id = \"test\"\ntitle = \"Test\"",
			wantBody: "\n\nBody content",
		},
		{
			name: "YAML frontmatter",
			content: `---
id: test
title: Test
---

Body content`,
			wantTOML: false,
			wantFM:   "id: test\ntitle: Test",
			wantBody: "\n\nBody content",
		},
		{
			name:     "No frontmatter",
			content:  "Just body content",
			wantTOML: false,
			wantFM:   "",
			wantBody: "Just body content",
		},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFM, gotBody, gotTOML := parser.extractFrontmatter(tt.content)

			if gotTOML != tt.wantTOML {
				t.Errorf("TOML flag: got %v, want %v", gotTOML, tt.wantTOML)
			}

			if gotFM != tt.wantFM {
				t.Errorf("Frontmatter:\ngot:  %q\nwant: %q", gotFM, tt.wantFM)
			}

			if gotBody != tt.wantBody {
				t.Errorf("Body:\ngot:  %q\nwant: %q", gotBody, tt.wantBody)
			}
		})
	}
}

// Helper function to create temporary MCQ files for testing
func createTempMCQFile(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test-mcq.md")
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	return tmpFile
}

// TestMCQFormatCompatibilityMatrix tests every possible valid MCQ format
func TestMCQFormatCompatibilityMatrix(t *testing.T) {
	formats := []struct {
		name           string
		content        string
		expectedParse  bool
		description    string
	}{
		{
			name: "Standard TOML frontmatter with dashes",
			content: `+++
id = "test-001"
type = "mcq"
level = "L3"
+++

## Question
What is the capital of France?

## Options
- A) Paris
- B) London
- C) Berlin
- D) Madrid

## Answer
A`,
			expectedParse: true,
			description: "Most common format in existing content",
		},
		{
			name: "YAML frontmatter variant",
			content: `---
id: test-002
type: mcq
level: "L4"
---

## Question
What is 2 + 2?

## Options
- A) 3
- B) 4

## Answer
B`,
			expectedParse: true,
			description: "Alternative frontmatter syntax support",
		},
		{
			name: "Options without dash prefix",
			content: `+++
id = "test-003"
type = "mcq"
+++

## Question
Select the correct statement:

## Options
A) Option 1
B) Option 2
C) Option 3

## Answer
C`,
			expectedParse: true,
			description: "Cleaner option format without bullet points",
		},
		{
			name: "Options with asterisk bullets",
			content: `+++
id = "test-004"
type = "mcq"
+++

## Question
Which is correct?

## Options
* A) Answer A
* B) Answer B
* C) Answer C

## Answer
B`,
			expectedParse: true,
			description: "Alternative bullet style",
		},
		{
			name: "Mixed case section headers",
			content: `+++
id = "test-005"
type = "mcq"
+++

## question
Text here?

## OPTIONS
- A) First option
- B) Second option

## answer
B`,
			expectedParse: true,
			description: "Case-insensitive header parsing",
		},
		{
			name: "MCQ with auto-detection (no type field)",
			content: `+++
id = "test-006"
level = "L3"
+++

## Question
This should be auto-detected as MCQ:

## Options
- A) Yes
- B) No
- C) Maybe

## Answer
A`,
			expectedParse: true,
			description: "Auto-detect MCQ when options and answer present",
		},
		{
			name: "MCQ with explanation",
			content: `+++
id = "test-007"
type = "mcq"
+++

## Question
What causes seasons?

## Options
- A) Distance from sun
- B) Earth's tilt
- C) Moon's gravity

## Answer
B

## Explanation
Earth's tilt causes different amounts of sunlight to reach different parts of the planet throughout the year.`,
			expectedParse: true,
			description: "MCQ with detailed explanation",
		},
		{
			name: "Invalid MCQ - Missing answer",
			content: `+++
id = "test-008"
type = "mcq"
+++

## Question
Missing answer field:

## Options
- A) Option 1
- B) Option 2

## [Should have Answer section]`,
			expectedParse: true, // Should parse but answer field will be empty
			description: "MCQ with missing answer should still parse",
		},
	}

	parser := NewParser()
	for _, tt := range formats {
		t.Run(tt.name, func(t *testing.T) {
			tempFile := createTempMCQFile(t, tt.content)
			q, err := parser.ParseQuestionFile(tempFile)
			
			if tt.expectedParse {
				require.NoError(t, err, "Format should parse: %s", tt.description)
				assert.Equal(t, "mcq", q.Type, "Should be detected as MCQ type")
				assert.NotEmpty(t, q.Options, "Options should be extracted")
				assert.NotNil(t, q.Options, "Options array should not be nil")
				
				// Auto-detection test
				if tt.name == "MCQ with auto-detection (no type field)" {
					assert.Equal(t, "mcq", q.Type, "Should auto-detect as MCQ")
				}
				
				// Check answer field (may be empty for some valid formats)
				if tt.name != "Invalid MCQ - Missing answer" {
					assert.NotEmpty(t, q.Answer, "Answer should be found")
				}
			} else {
				assert.Error(t, err, "Format should fail: %s", tt.description)
			}
		})
	}
}

// TestMCQOptionCountBehavior tests different option counts
func TestMCQOptionCountBehavior(t *testing.T) {
	testCases := []struct {
		name         string
		optionCount  int
		shouldWork   bool
		description  string
	}{
		{
			name:        "Two options (binary choice)",
			optionCount: 2,
			shouldWork:  true,
			description: "True/False style questions",
		},
		{
			name:        "Standard four options",
			optionCount: 4,
			shouldWork:  true,
			description: "Most common MCQ format",
		},
		{
			name:        "Six options",
			optionCount: 6,
			shouldWork:  true,
			description: "Extended choice questions",
		},
		{
			name:        "Empty options",
			optionCount: 0,
			shouldWork:  true,
			description: "Empty options section should still parse",
		},
	}

	// Helper function to generate MCQ content with N options
	generateMCQWithNOptions := func(n int) string {
		var optionsContent string
		if n > 0 {
			optionsContent = "## Options\n"
			for i := 0; i < n; i++ {
				letter := fmt.Sprintf("%c", 'A'+i)
				optionsContent += fmt.Sprintf("- %s) Option %d\n", letter, i+1)
			}
			optionsContent += "\n## Answer\n" + fmt.Sprintf("%c", 'A') + "\n"
		} else {
			optionsContent = "## Options\n\n## Answer\nA\n"
		}
		
		return fmt.Sprintf(`+++
id = "count-test-%d"
type = "mcq"
level = "L3"
+++

## Question
Test question with %d options:

%s`, n, n, optionsContent)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mcqContent := generateMCQWithNOptions(tc.optionCount)
			tempFile := createTempMCQFile(t, mcqContent)
			parser := NewParser()
			
			q, err := parser.ParseQuestionFile(tempFile)
			require.NoError(t, err, "Should parse: %s", tc.description)
			
			assert.Equal(t, "mcq", q.Type)
			if tc.optionCount > 0 {
				assert.Len(t, q.Options, tc.optionCount, "Should have %d options", tc.optionCount)
			} else {
				assert.Len(t, q.Options, 0, "Should have no options")
			}
		})
	}
}

// TestMCQAnswerValidation tests various answer formats
func TestMCQAnswerValidation(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		expectedAnswer string
		description   string
	}{
		{
			name: "Single letter uppercase",
			content: `+++
id = "answer-test-1"
type = "mcq"
+++

## Question
Test?

## Options
- A) Yes
- B) No

## Answer
A`,
			expectedAnswer: "A",
			description: "Standard single letter answer",
		},
		{
			name: "Single letter lowercase",
			content: `+++
id = "answer-test-2"
type = "mcq"
+++

## Question
Test?

## Options
- A) Yes
- B) No

## Answer
b`,
			expectedAnswer: "b",
			description: "Lowercase answer",
		},
		{
			name: "Letter with parenthesis",
			content: `+++
id = "answer-test-3"
type = "mcq"
+++

## Question
Test?

## Options
- A) Yes
- B) No

## Answer
A)`,
			expectedAnswer: "A)",
			description: "Answer includes formatting",
		},
		{
			name: "Full option text as answer",
			content: `+++
id = "answer-test-4"
type = "mcq"
+++

## Question
Test?

## Options
- A) Yes, this is correct
- B) No, this is wrong

## Answer
A) Yes, this is correct`,
			expectedAnswer: "A) Yes, this is correct",
			description: "Complete option as answer",
		},
		{
			name: "Empty answer",
			content: `+++
id = "answer-test-5"
type = "mcq"
+++

## Question
Test?

## Options
- A) Yes
- B) No

## Answer

## Explanation
Explanation here`,
			expectedAnswer: "",
			description: "Empty answer field",
		},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempFile := createTempMCQFile(t, tt.content)
			q, err := parser.ParseQuestionFile(tempFile)
			
			require.NoError(t, err, "Should parse: %s", tt.description)
			assert.Equal(t, tt.expectedAnswer, q.Answer, "Answer mismatch")
		})
	}
}
