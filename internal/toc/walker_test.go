package toc

import (
	"strings"
	"testing"
)

func TestProcessPath(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		recursive bool
		wantErr   bool
		minFiles  int
	}{
		{
			name:      "single markdown file",
			path:      "testdata/sample1.md",
			recursive: false,
			wantErr:   false,
			minFiles:  1,
		},
		{
			name:      "directory with multiple files",
			path:      "testdata",
			recursive: false,
			wantErr:   false,
			minFiles:  3,
		},
		{
			name:      "non-existent path",
			path:      "testdata/does_not_exist",
			recursive: false,
			wantErr:   true,
			minFiles:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{
				MaxDepth:  6,
				Recursive: tt.recursive,
			}
			results, err := ProcessPath(tt.path, opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(results) < tt.minFiles {
				t.Errorf("ProcessPath() got %d files, want at least %d", len(results), tt.minFiles)
			}
		})
	}
}

func TestIsMarkdownFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"markdown .md", "file.md", true},
		{"markdown .markdown", "file.markdown", true},
		{"uppercase .MD", "file.MD", true},
		{"text file", "file.txt", false},
		{"go file", "file.go", false},
		{"no extension", "README", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isMarkdownFile(tt.path)
			if got != tt.expected {
				t.Errorf("isMarkdownFile(%s) = %v, want %v", tt.path, got, tt.expected)
			}
		})
	}
}

func TestGenerateCombinedTOC(t *testing.T) {
	results := map[string][]Header{
		"file1.md": {
			{Level: 1, Title: "File 1", Anchor: "file-1"},
			{Level: 2, Title: "Section 1", Anchor: "section-1"},
		},
		"file2.md": {
			{Level: 1, Title: "File 2", Anchor: "file-2"},
		},
	}

	tests := []struct {
		name   string
		format string
		check  func(string) bool
	}{
		{
			name:   "markdown format with multiple files",
			format: "markdown",
			check: func(s string) bool {
				return len(s) > 0
			},
		},
		{
			name:   "json format with multiple files",
			format: "json",
			check: func(s string) bool {
				return len(s) > 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{
				MaxDepth: 6,
				Format:   tt.format,
			}
			result := GenerateCombinedTOC(results, opts)
			if !tt.check(result) {
				t.Errorf("GenerateCombinedTOC() check failed")
			}
		})
	}
}

func TestGenerateCombinedTOCEmpty(t *testing.T) {
	results := map[string][]Header{}
	opts := Options{MaxDepth: 6, Format: "markdown"}
	result := GenerateCombinedTOC(results, opts)
	if result != "" {
		t.Error("Empty results should produce empty output")
	}
}

func TestGenerateCombinedTOCSingleFile(t *testing.T) {
	results := map[string][]Header{
		"single.md": {
			{Level: 1, Title: "Title", Anchor: "title"},
		},
	}
	opts := Options{MaxDepth: 6, Format: "markdown"}
	result := GenerateCombinedTOC(results, opts)
	if result == "" {
		t.Error("Single file results should produce output")
	}
}

func TestGenerateCombinedTOCAlwaysShowsFilenames(t *testing.T) {
	tests := []struct {
		name     string
		results  map[string][]Header
		format   string
		wantFile string
	}{
		{
			name: "single file shows filename in markdown",
			results: map[string][]Header{
				"testdata/sample1.md": {
					{Level: 1, Title: "Main Title", Anchor: "main-title"},
					{Level: 2, Title: "Section 1", Anchor: "section-1"},
				},
			},
			format:   "markdown",
			wantFile: "sample1.md",
		},
		{
			name: "multiple files show all filenames in markdown",
			results: map[string][]Header{
				"testdata/file1.md": {
					{Level: 1, Title: "File 1", Anchor: "file-1"},
				},
				"testdata/file2.md": {
					{Level: 1, Title: "File 2", Anchor: "file-2"},
				},
			},
			format:   "markdown",
			wantFile: "file1.md",
		},
		{
			name: "single file shows filename in text format",
			results: map[string][]Header{
				"testdata/sample2.md": {
					{Level: 1, Title: "Title", Anchor: "title"},
				},
			},
			format:   "text",
			wantFile: "sample2.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{
				MaxDepth: 6,
				Format:   tt.format,
			}
			result := GenerateCombinedTOC(tt.results, opts)

			if result == "" {
				t.Error("Result should not be empty")
			}

			if !strings.Contains(result, tt.wantFile) {
				t.Errorf("Result should contain filename '%s', got:\n%s", tt.wantFile, result)
			}
		})
	}
}

func TestGenerateCombinedTOCDirectoryFormat(t *testing.T) {
	results := map[string][]Header{
		"docs/chapter1.md": {
			{Level: 1, Title: "Chapter 1", Anchor: "chapter-1"},
			{Level: 2, Title: "Section 1.1", Anchor: "section-11"},
		},
		"docs/chapter2.md": {
			{Level: 1, Title: "Chapter 2", Anchor: "chapter-2"},
		},
	}

	opts := Options{MaxDepth: 6, Format: "markdown"}
	result := GenerateCombinedTOC(results, opts)

	// Check that both filenames appear
	if !strings.Contains(result, "chapter1.md") {
		t.Error("Result should contain 'chapter1.md'")
	}
	if !strings.Contains(result, "chapter2.md") {
		t.Error("Result should contain 'chapter2.md'")
	}

	// Check for TOC structure
	if !strings.Contains(result, "## Table of Contents") {
		t.Error("Result should contain TOC header")
	}

	// Check for individual file sections
	if !strings.Contains(result, "### chapter1.md") {
		t.Error("Result should have section header for chapter1.md")
	}
	if !strings.Contains(result, "### chapter2.md") {
		t.Error("Result should have section header for chapter2.md")
	}

	// Check for TOC entries
	if !strings.Contains(result, "[Chapter 1](#chapter-1)") {
		t.Error("Result should contain TOC entry for Chapter 1")
	}
	if !strings.Contains(result, "[Section 1.1](#section-11)") {
		t.Error("Result should contain TOC entry for Section 1.1")
	}
}

func TestGenerateCombinedTOCTextFormat(t *testing.T) {
	results := map[string][]Header{
		"path/to/file.md": {
			{Level: 1, Title: "Title", Anchor: "title"},
			{Level: 2, Title: "Subtitle", Anchor: "subtitle"},
		},
	}

	opts := Options{MaxDepth: 6, Format: "text"}
	result := GenerateCombinedTOC(results, opts)

	// Should show just the basename
	if !strings.Contains(result, "file.md") {
		t.Error("Result should contain filename 'file.md'")
	}

	// Should contain the titles
	if !strings.Contains(result, "Title") {
		t.Error("Result should contain 'Title'")
	}
	if !strings.Contains(result, "Subtitle") {
		t.Error("Result should contain 'Subtitle'")
	}

	// Should NOT contain markdown link syntax in text format
	if strings.Contains(result, "[") || strings.Contains(result, "]") {
		t.Error("Text format should not contain markdown link syntax")
	}
}
