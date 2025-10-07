package models

// Question represents a single question file from a chapters/ directory.
type Question struct {
	// --- File Metadata (not in YAML) ---
	Path      string `yaml:"-"`
	ChapterID string `yaml:"-"` // e.g., "ch3-storage"

	// --- Frontmatter ---
	ID             string   `yaml:"id"`
	Day            int      `yaml:"day"`
	Tags           []string `yaml:"tags"`
	RelatedStories []string `yaml:"related_stories"`

	// --- Content (not in YAML) ---
	ContentMarkdown string `yaml:"-"`
	Scenario        string `yaml:"-"`
	QuestionText    string `yaml:"-"`
	Explanation     string `yaml:"-"`
	Hook            string `yaml:"-"`
}
