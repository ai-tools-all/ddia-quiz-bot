package markdown

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseQuestionFile(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantID      string
		wantTitle   string
		wantLevel   string
		wantQuestion string
		wantErr     bool
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
		name             string
		content          string
		wantFrontmatter  string
		wantBody         string
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
			gotFrontmatter, gotBody := parser.extractFrontmatter(tt.content)

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
