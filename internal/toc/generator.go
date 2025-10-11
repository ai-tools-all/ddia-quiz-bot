package toc

import (
	"encoding/json"
	"fmt"
	"strings"
)

func GenerateTOC(headers []Header, opts Options) string {
	headers = FilterByDepth(headers, opts.MaxDepth)

	switch opts.Format {
	case "json":
		return GenerateJSON(headers)
	case "text":
		return GenerateText(headers)
	default:
		return GenerateMarkdown(headers)
	}
}

func GenerateMarkdown(headers []Header) string {
	if len(headers) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("## Table of Contents\n\n")

	for _, h := range headers {
		indent := strings.Repeat("  ", h.Level-1)
		sb.WriteString(fmt.Sprintf("%s- [%s](#%s)\n", indent, h.Title, h.Anchor))
	}

	return sb.String()
}

func GenerateText(headers []Header) string {
	if len(headers) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, h := range headers {
		indent := strings.Repeat("  ", h.Level-1)
		sb.WriteString(fmt.Sprintf("%s%s\n", indent, h.Title))
	}

	return sb.String()
}

type TOCEntry struct {
	Level   int    `json:"level"`
	Title   string `json:"title"`
	Anchor  string `json:"anchor"`
	LineNum int    `json:"line_number"`
}

type TOCOutput struct {
	Headers []TOCEntry `json:"headers"`
	Count   int        `json:"count"`
}

func GenerateJSON(headers []Header) string {
	entries := make([]TOCEntry, len(headers))
	for i, h := range headers {
		entries[i] = TOCEntry{
			Level:   h.Level,
			Title:   h.Title,
			Anchor:  h.Anchor,
			LineNum: h.LineNum,
		}
	}

	output := TOCOutput{
		Headers: entries,
		Count:   len(entries),
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	return string(data)
}
