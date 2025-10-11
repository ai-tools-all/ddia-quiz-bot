package toc

import (
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
