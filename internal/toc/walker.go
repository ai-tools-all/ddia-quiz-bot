package toc

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ProcessPath(path string, opts Options) (map[string][]Header, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]Header)

	if !info.IsDir() {
		if isMarkdownFile(path) {
			headers, err := ParseFile(path)
			if err != nil {
				return nil, err
			}
			result[path] = headers
		}
		return result, nil
	}

	var markdownFiles []string
	err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && !opts.Recursive && filePath != path {
			return filepath.SkipDir
		}

		if !info.IsDir() && isMarkdownFile(filePath) {
			markdownFiles = append(markdownFiles, filePath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if opts.ConfirmThreshold > 0 && len(markdownFiles) > opts.ConfirmThreshold && !opts.SkipConfirmation {
		fmt.Fprintf(os.Stderr, "\nFound %d markdown files. This may take a while.\n", len(markdownFiles))
		fmt.Fprintf(os.Stderr, "Continue? [y/N]: ")

		var response string
		fmt.Scanln(&response)
		response = strings.ToLower(strings.TrimSpace(response))

		if response != "y" && response != "yes" {
			return nil, fmt.Errorf("operation cancelled by user")
		}
	}

	for _, filePath := range markdownFiles {
		headers, err := ParseFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("error parsing %s: %w", filePath, err)
		}
		result[filePath] = headers
	}

	return result, nil
}

func isMarkdownFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md" || ext == ".markdown"
}

func GenerateCombinedTOC(results map[string][]Header, opts Options) string {
	if len(results) == 0 {
		return ""
	}

	var sb strings.Builder

	switch opts.Format {
	case "json":
		return generateCombinedJSON(results)
	case "text":
		for path, headers := range results {
			sb.WriteString(fmt.Sprintf("%s\n", filepath.Base(path)))
			sb.WriteString(GenerateText(FilterByDepth(headers, opts.MaxDepth)))
			sb.WriteString("\n")
		}
	default:
		sb.WriteString("## Table of Contents\n\n")
		for path, headers := range results {
			sb.WriteString(fmt.Sprintf("### %s\n\n", filepath.Base(path)))
			headers = FilterByDepth(headers, opts.MaxDepth)
			for _, h := range headers {
				indent := strings.Repeat("  ", h.Level-1)
				sb.WriteString(fmt.Sprintf("%s- [%s](#%s)\n", indent, h.Title, h.Anchor))
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func generateCombinedJSON(results map[string][]Header) string {
	type FileEntry struct {
		Path    string     `json:"path"`
		Headers []TOCEntry `json:"headers"`
	}

	type CombinedOutput struct {
		Files []FileEntry `json:"files"`
		Total int         `json:"total_headers"`
	}

	var files []FileEntry
	totalCount := 0

	for path, headers := range results {
		entries := make([]TOCEntry, len(headers))
		for i, h := range headers {
			entries[i] = TOCEntry{
				Level:   h.Level,
				Title:   h.Title,
				Anchor:  h.Anchor,
				LineNum: h.LineNum,
			}
		}

		files = append(files, FileEntry{
			Path:    path,
			Headers: entries,
		})
		totalCount += len(headers)
	}

	output := CombinedOutput{
		Files: files,
		Total: totalCount,
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	return string(data)
}
