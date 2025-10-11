package toc

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGenerateMarkdown(t *testing.T) {
	headers := []Header{
		{Level: 1, Title: "Main Title", Anchor: "main-title"},
		{Level: 2, Title: "Section 1", Anchor: "section-1"},
		{Level: 2, Title: "Section 2", Anchor: "section-2"},
	}

	result := GenerateMarkdown(headers)

	if !strings.Contains(result, "## Table of Contents") {
		t.Error("Markdown output should contain 'Table of Contents' header")
	}

	if !strings.Contains(result, "[Main Title](#main-title)") {
		t.Error("Markdown output should contain main title link")
	}

	if !strings.Contains(result, "[Section 1](#section-1)") {
		t.Error("Markdown output should contain section 1 link")
	}
}

func TestGenerateText(t *testing.T) {
	headers := []Header{
		{Level: 1, Title: "Main Title", Anchor: "main-title"},
		{Level: 2, Title: "Section 1", Anchor: "section-1"},
		{Level: 3, Title: "Subsection", Anchor: "subsection"},
	}

	result := GenerateText(headers)

	if !strings.Contains(result, "Main Title") {
		t.Error("Text output should contain main title")
	}

	if !strings.Contains(result, "Section 1") {
		t.Error("Text output should contain section 1")
	}

	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}
}

func TestGenerateJSON(t *testing.T) {
	headers := []Header{
		{Level: 1, Title: "Main Title", Anchor: "main-title", LineNum: 1},
		{Level: 2, Title: "Section 1", Anchor: "section-1", LineNum: 3},
	}

	result := GenerateJSON(headers)

	var output TOCOutput
	err := json.Unmarshal([]byte(result), &output)
	if err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if output.Count != 2 {
		t.Errorf("Expected count 2, got %d", output.Count)
	}

	if len(output.Headers) != 2 {
		t.Errorf("Expected 2 headers, got %d", len(output.Headers))
	}

	if output.Headers[0].Title != "Main Title" {
		t.Errorf("Expected first title 'Main Title', got '%s'", output.Headers[0].Title)
	}
}

func TestGenerateTOC(t *testing.T) {
	headers := []Header{
		{Level: 1, Title: "Title", Anchor: "title"},
		{Level: 2, Title: "Section", Anchor: "section"},
	}

	tests := []struct {
		name   string
		format string
		check  func(string) bool
	}{
		{
			name:   "markdown format",
			format: "markdown",
			check: func(s string) bool {
				return strings.Contains(s, "## Table of Contents")
			},
		},
		{
			name:   "text format",
			format: "text",
			check: func(s string) bool {
				return strings.Contains(s, "Title") && !strings.Contains(s, "[")
			},
		},
		{
			name:   "json format",
			format: "json",
			check: func(s string) bool {
				return strings.Contains(s, `"title"`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{
				MaxDepth: 6,
				Format:   tt.format,
			}
			result := GenerateTOC(headers, opts)
			if !tt.check(result) {
				t.Errorf("GenerateTOC() format check failed for %s", tt.format)
			}
		})
	}
}

func TestGenerateMarkdownWithIndentation(t *testing.T) {
	headers := []Header{
		{Level: 1, Title: "H1", Anchor: "h1"},
		{Level: 2, Title: "H2", Anchor: "h2"},
		{Level: 3, Title: "H3", Anchor: "h3"},
		{Level: 4, Title: "H4", Anchor: "h4"},
	}

	result := GenerateMarkdown(headers)
	lines := strings.Split(result, "\n")

	var foundH1, foundH2, foundH3, foundH4 bool
	for _, line := range lines {
		if strings.Contains(line, "[H1]") && !strings.HasPrefix(line, " ") {
			foundH1 = true
		}
		if strings.Contains(line, "[H2]") && strings.HasPrefix(line, "  ") {
			foundH2 = true
		}
		if strings.Contains(line, "[H3]") && strings.HasPrefix(line, "    ") {
			foundH3 = true
		}
		if strings.Contains(line, "[H4]") && strings.HasPrefix(line, "      ") {
			foundH4 = true
		}
	}

	if !foundH1 {
		t.Error("H1 should have no indentation")
	}
	if !foundH2 {
		t.Error("H2 should have 2 spaces indentation")
	}
	if !foundH3 {
		t.Error("H3 should have 4 spaces indentation")
	}
	if !foundH4 {
		t.Error("H4 should have 6 spaces indentation")
	}
}

func TestGenerateEmptyHeaders(t *testing.T) {
	headers := []Header{}

	markdown := GenerateMarkdown(headers)
	if markdown != "" {
		t.Error("Empty headers should produce empty markdown")
	}

	text := GenerateText(headers)
	if text != "" {
		t.Error("Empty headers should produce empty text")
	}
}
