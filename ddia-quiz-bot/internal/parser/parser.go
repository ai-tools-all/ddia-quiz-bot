package parser

import (
	"bytes"
	"errors"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/your-username/ddia-quiz-bot/internal/models"
	"gopkg.in/yaml.v3"
)

var (
	frontmatterDelimiter = []byte("---")
	// Regex to extract specific sections from markdown
	sectionRegex = regexp.MustCompile(`(?s)##\s*(scenario|question|explanation|hook)\n(.*?)(?:\n##|$)`)
)

// ParseFile reads a byte slice and separates the YAML frontmatter from the markdown body.
func ParseFile(data []byte) (yamlData, markdownData []byte, err error) {
	if !bytes.HasPrefix(data, frontmatterDelimiter) {
		return nil, nil, errors.New("file does not have a YAML frontmatter")
	}

	parts := bytes.SplitN(data, frontmatterDelimiter, 3)
	if len(parts) < 3 {
		return nil, nil, errors.New("invalid frontmatter format")
	}

	return parts[1], bytes.TrimSpace(parts[2]), nil
}

// UnmarshalFrontmatter parses the YAML frontmatter into a given struct.
func UnmarshalFrontmatter(yamlData []byte, v interface{}) error {
	return yaml.Unmarshal(yamlData, v)
}

// ParseQuestionSections extracts specific h2 sections from a question's markdown body.
func ParseQuestionSections(markdownData []byte, q *models.Question) {
	matches := sectionRegex.FindAllSubmatch(markdownData, -1)
	rawContent := string(markdownData)

	for _, match := range matches {
		sectionTitle := strings.ToLower(string(match[1]))
		sectionContent := strings.TrimSpace(string(match[2]))

		switch sectionTitle {
		case "scenario":
			q.Scenario = sectionContent
		case "question":
			q.QuestionText = sectionContent
		case "explanation":
			q.Explanation = sectionContent
		case "hook":
			q.Hook = sectionContent
		}

		// Remove the parsed section from the raw content to isolate the main question title
		rawContent = strings.Replace(rawContent, string(match[0]), "", 1)
	}
	q.ContentMarkdown = strings.TrimSpace(rawContent)
}

// ReadAndParseFile is a helper to read a file and parse it.
func ReadAndParseFile(path string) ([]byte, []byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	return ParseFile(data)
}
