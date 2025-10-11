package toc

import (
	"testing"
)

func TestParseMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name: "simple headers",
			content: `# Title
## Section 1
### Subsection 1.1
## Section 2`,
			expected: 4,
		},
		{
			name: "with code blocks",
			content: `# Title
## Section 1
` + "```" + `
# Not a header
## Also not a header
` + "```" + `
## Section 2`,
			expected: 3,
		},
		{
			name: "no headers",
			content: `This is just text
with no headers at all`,
			expected: 0,
		},
		{
			name: "with YAML frontmatter",
			content: `---
id: test
tags: [test]
---

# Title
## Section`,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers, err := ParseMarkdown(tt.content)
			if err != nil {
				t.Fatalf("ParseMarkdown() error = %v", err)
			}
			if len(headers) != tt.expected {
				t.Errorf("ParseMarkdown() got %d headers, want %d", len(headers), tt.expected)
			}
		})
	}
}

func TestGenerateAnchor(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected string
	}{
		{
			name:     "simple title",
			title:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "with special chars",
			title:    "Hello, World!",
			expected: "hello-world",
		},
		{
			name:     "with multiple spaces",
			title:    "Hello    World",
			expected: "hello-world",
		},
		{
			name:     "with hyphens",
			title:    "Hello-World-Test",
			expected: "hello-world-test",
		},
		{
			name:     "with numbers",
			title:    "Section 1.2.3",
			expected: "section-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateAnchor(tt.title)
			if got != tt.expected {
				t.Errorf("GenerateAnchor() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFilterByDepth(t *testing.T) {
	headers := []Header{
		{Level: 1, Title: "H1"},
		{Level: 2, Title: "H2"},
		{Level: 3, Title: "H3"},
		{Level: 4, Title: "H4"},
	}

	tests := []struct {
		name     string
		maxDepth int
		expected int
	}{
		{"depth 1", 1, 1},
		{"depth 2", 2, 2},
		{"depth 3", 3, 3},
		{"depth 0 (all)", 0, 4},
		{"depth 10 (all)", 10, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := FilterByDepth(headers, tt.maxDepth)
			if len(filtered) != tt.expected {
				t.Errorf("FilterByDepth(%d) = %d headers, want %d", tt.maxDepth, len(filtered), tt.expected)
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		wantErr  bool
		minCount int
	}{
		{
			name:     "sample file with headers",
			filepath: "testdata/sample1.md",
			wantErr:  false,
			minCount: 3,
		},
		{
			name:     "file with frontmatter",
			filepath: "testdata/sample2.md",
			wantErr:  false,
			minCount: 3,
		},
		{
			name:     "file with no headers",
			filepath: "testdata/no_headers.md",
			wantErr:  false,
			minCount: 0,
		},
		{
			name:     "file with code blocks",
			filepath: "testdata/with_code.md",
			wantErr:  false,
			minCount: 2,
		},
		{
			name:     "non-existent file",
			filepath: "testdata/does_not_exist.md",
			wantErr:  true,
			minCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers, err := ParseFile(tt.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(headers) < tt.minCount {
				t.Errorf("ParseFile() got %d headers, want at least %d", len(headers), tt.minCount)
			}
		})
	}
}

func TestParseMarkdownIgnoresCodeBlocks(t *testing.T) {
	content := `# Real Header
` + "```" + `markdown
# Fake Header in Code
## Another Fake
` + "```" + `
## Real Section`

	headers, err := ParseMarkdown(content)
	if err != nil {
		t.Fatalf("ParseMarkdown() error = %v", err)
	}

	if len(headers) != 2 {
		t.Errorf("Expected 2 headers, got %d", len(headers))
	}

	for _, h := range headers {
		if h.Title == "Fake Header in Code" || h.Title == "Another Fake" {
			t.Errorf("Code block header was not ignored: %s", h.Title)
		}
	}
}
