package toc

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var headerRegex = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)

func ParseFile(filepath string) ([]Header, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var headers []Header
	scanner := bufio.NewScanner(file)
	lineNum := 0

	inCodeBlock := false
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inCodeBlock = !inCodeBlock
			continue
		}

		if inCodeBlock {
			continue
		}

		matches := headerRegex.FindStringSubmatch(line)
		if matches != nil {
			level := len(matches[1])
			title := strings.TrimSpace(matches[2])

			headers = append(headers, Header{
				Level:   level,
				Title:   title,
				Anchor:  GenerateAnchor(title),
				LineNum: lineNum,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return headers, nil
}

func GenerateAnchor(title string) string {
	anchor := strings.ToLower(title)
	anchor = strings.TrimSpace(anchor)

	anchor = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(anchor, "")
	anchor = regexp.MustCompile(`[\s]+`).ReplaceAllString(anchor, "-")
	anchor = regexp.MustCompile(`-+`).ReplaceAllString(anchor, "-")
	anchor = strings.Trim(anchor, "-")

	return anchor
}

func ParseMarkdown(content string) ([]Header, error) {
	var headers []Header
	lines := strings.Split(content, "\n")
	inCodeBlock := false

	for lineNum, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inCodeBlock = !inCodeBlock
			continue
		}

		if inCodeBlock {
			continue
		}

		matches := headerRegex.FindStringSubmatch(line)
		if matches != nil {
			level := len(matches[1])
			title := strings.TrimSpace(matches[2])

			headers = append(headers, Header{
				Level:   level,
				Title:   title,
				Anchor:  GenerateAnchor(title),
				LineNum: lineNum + 1,
			})
		}
	}

	return headers, nil
}

func FilterByDepth(headers []Header, maxDepth int) []Header {
	if maxDepth <= 0 || maxDepth > 6 {
		return headers
	}

	filtered := make([]Header, 0, len(headers))
	for _, h := range headers {
		if h.Level <= maxDepth {
			filtered = append(filtered, h)
		}
	}
	return filtered
}

func GenerateAnchorWithStyle(title, style string) string {
	switch style {
	case "github":
		return GenerateAnchor(title)
	default:
		return GenerateAnchor(title)
	}
}

func countLeadingHashes(line string) int {
	count := 0
	for _, ch := range line {
		if ch == '#' {
			count++
		} else if ch == ' ' {
			break
		} else {
			return 0
		}
	}
	return count
}

func extractTitle(line string, hashCount int) string {
	if hashCount == 0 {
		return ""
	}
	title := strings.TrimPrefix(line, strings.Repeat("#", hashCount))
	return strings.TrimSpace(title)
}

func isValidHeaderLevel(level int) bool {
	return level >= 1 && level <= 6
}

func normalizeWhitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func removeMarkdownFormatting(title string) string {
	title = regexp.MustCompile(`\[([^\]]+)\]\([^\)]+\)`).ReplaceAllString(title, "$1")
	title = regexp.MustCompile(`[*_~`+"`"+`]+`).ReplaceAllString(title, "")
	title = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(title, "")
	return title
}

func parseHeaderLevel(line string) (int, string, bool) {
	matches := headerRegex.FindStringSubmatch(line)
	if matches == nil {
		return 0, "", false
	}

	level, _ := strconv.Atoi(strconv.Itoa(len(matches[1])))
	title := strings.TrimSpace(matches[2])
	title = removeMarkdownFormatting(title)

	return level, title, true
}
