package models

// Story represents a single story file from the stories/ directory.
type Story struct {
	// --- File Metadata (not in YAML) ---
	Path     string `yaml:"-"`
	Filename string `yaml:"-"`

	// --- Frontmatter ---
	Title          string   `yaml:"title"`
	Source         string   `yaml:"source"`
	URL            string   `yaml:"url"`
	Date           string   `yaml:"date"` // Keep as string for parsing flexibility
	Author         string   `yaml:"author"`
	Tags           []string `yaml:"tags"`
	RelatesTo      []string `yaml:"relates_to"` // Links to specific Question IDs
	Topics         []string `yaml:"topics"`
	Type           string   `yaml:"type"`
	Difficulty     string   `yaml:"difficulty"`
	BestPairedWith []string `yaml:"best_paired_with"`
	EngagementHook string   `yaml:"engagement_hook"`

	// --- Content (not in YAML) ---
	ContentMarkdown string `yaml:"-"`
}
