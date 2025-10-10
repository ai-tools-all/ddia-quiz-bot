package markdown

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/abhishek/ddia-clicker/internal/models"
	"gopkg.in/yaml.v3"
)

// Parser handles markdown file parsing
type Parser struct {
	frontmatterRegex *regexp.Regexp
}

// NewParser creates a new markdown parser
func NewParser() *Parser {
	return &Parser{
		frontmatterRegex: regexp.MustCompile(`(?s)^---\n(.*?)\n---`),
	}
}

// ParseQuestionFile parses a markdown file containing a quiz question
func (p *Parser) ParseQuestionFile(filepath string) (*models.Question, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Extract frontmatter
	frontmatter, body := p.extractFrontmatter(string(content))
	if frontmatter == "" {
		// No frontmatter, try to parse as pure markdown
		return p.parseMarkdownBody(body)
	}

	// Parse YAML frontmatter
	var metadata map[string]interface{}
	if err := yaml.Unmarshal([]byte(frontmatter), &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	question := &models.Question{}

	// Extract metadata fields
	if id, ok := metadata["id"].(string); ok {
		question.ID = id
	} else if id, ok := metadata["question_id"].(string); ok {
		question.ID = id
	}

	if title, ok := metadata["title"].(string); ok {
		question.Title = title
	}

	if level, ok := metadata["level"].(string); ok {
		question.Level = level
	}

	if category, ok := metadata["category"].(string); ok {
		question.Category = category
	}

	// Parse the markdown body for question details
	p.parseBody(body, question)

	// Validate question
	if question.ID == "" {
		return nil, fmt.Errorf("question ID not found")
	}

	return question, nil
}

// extractFrontmatter separates YAML frontmatter from markdown body
func (p *Parser) extractFrontmatter(content string) (string, string) {
	if !strings.HasPrefix(content, "---\n") {
		return "", content
	}

	matches := p.frontmatterRegex.FindStringSubmatch(content)
	if len(matches) >= 2 {
		frontmatter := matches[1]
		body := strings.TrimPrefix(content, matches[0])
		return frontmatter, body
	}

	return "", content
}

// parseBody extracts question details from markdown body
func (p *Parser) parseBody(body string, question *models.Question) {
	scanner := bufio.NewScanner(strings.NewReader(body))
	currentSection := ""
	var sectionContent strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// Check for section headers
		if strings.HasPrefix(line, "## ") {
			// Save previous section if any
			p.saveSection(currentSection, sectionContent.String(), question)
			
			// Start new section
			currentSection = strings.TrimSpace(strings.TrimPrefix(line, "## "))
			sectionContent.Reset()
			continue
		}

		// Accumulate content for current section
		if currentSection != "" {
			sectionContent.WriteString(line)
			sectionContent.WriteString("\n")
		}
	}

	// Save last section
	p.saveSection(currentSection, sectionContent.String(), question)
}

// saveSection processes and saves section content to the question
func (p *Parser) saveSection(section, content string, question *models.Question) {
	content = strings.TrimSpace(content)
	if content == "" {
		return
	}

	// Normalize section name: lowercase and extract prefix before any dash or space
	sectionLower := strings.ToLower(section)
	
	// Split by dash to get the base section name (e.g., "main_question - Core Question" -> "main_question")
	if idx := strings.Index(sectionLower, " -"); idx != -1 {
		sectionLower = strings.TrimSpace(sectionLower[:idx])
	}
	
	// Also normalize underscores to spaces for comparison
	sectionNormalized := strings.ReplaceAll(sectionLower, "_", " ")

	// Check if section starts with expected prefix
	if sectionLower == "question" || sectionNormalized == "main question" || sectionLower == "main_question" {
		question.MainQuestion = content
	} else if sectionNormalized == "core concepts" || sectionLower == "core_concepts" {
		question.CoreConcepts = p.parseListSection(content)
	} else if sectionNormalized == "peripheral concepts" || sectionLower == "peripheral_concepts" {
		question.PeripheralConcepts = p.parseListSection(content)
	} else if strings.Contains(sectionNormalized, "sample excellent") || strings.Contains(sectionLower, "excellent answer") {
		question.SampleExcellent = content
	} else if strings.Contains(sectionNormalized, "sample acceptable") || strings.Contains(sectionLower, "acceptable answer") {
		question.SampleAcceptable = content
	} else if sectionNormalized == "evaluation rubric" || sectionLower == "rubric" {
		question.EvaluationRubric = p.parseRubricSection(content)
	}
}

// parseListSection parses a bullet list into string slice
func (p *Parser) parseListSection(content string) []string {
	var items []string
	lines := strings.Split(content, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "•") {
			item := strings.TrimSpace(line[1:])
			if item != "" {
				items = append(items, item)
			}
		} else if matches := regexp.MustCompile(`^\d+\.\s+(.+)`).FindStringSubmatch(line); len(matches) > 1 {
			items = append(items, matches[1])
		}
	}
	
	return items
}

// parseRubricSection parses evaluation rubric into a map
func (p *Parser) parseRubricSection(content string) map[string]string {
	rubric := make(map[string]string)
	lines := strings.Split(content, "\n")
	currentKey := ""
	currentValue := ""
	
	for _, line := range lines {
		// Check if it's a new rubric item
		if strings.Contains(line, ":") {
			// Save previous item if exists
			if currentKey != "" {
				rubric[currentKey] = strings.TrimSpace(currentValue)
			}
			
			parts := strings.SplitN(line, ":", 2)
			currentKey = strings.TrimSpace(parts[0])
			if len(parts) > 1 {
				currentValue = strings.TrimSpace(parts[1])
			} else {
				currentValue = ""
			}
		} else if currentKey != "" {
			// Continue accumulating value for current key
			currentValue += " " + strings.TrimSpace(line)
		}
	}
	
	// Save last item
	if currentKey != "" {
		rubric[currentKey] = strings.TrimSpace(currentValue)
	}
	
	return rubric
}

// parseMarkdownBody attempts to parse a pure markdown file without frontmatter
func (p *Parser) parseMarkdownBody(content string) (*models.Question, error) {
	question := &models.Question{}
	
	// Look for question ID in the content
	idRegex := regexp.MustCompile(`(?i)question[_\s]?id:\s*([^\n]+)`)
	if matches := idRegex.FindStringSubmatch(content); len(matches) > 1 {
		question.ID = strings.TrimSpace(matches[1])
	}
	
	// Parse the body
	p.parseBody(content, question)
	
	if question.ID == "" {
		return nil, fmt.Errorf("no question ID found in file")
	}
	
	return question, nil
}
