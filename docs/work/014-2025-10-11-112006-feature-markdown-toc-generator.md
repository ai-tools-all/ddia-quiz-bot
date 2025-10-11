# Markdown TOC Generator - Plan

**Category**: feature  
**Created**: 2025-10-11 11:20:06  
**Status**: Planning

## Overview

A standalone utility binary that generates table of contents (TOC) for markdown files. Can process single files or entire directories containing markdown files.

## Goals

- Read markdown file(s) and extract headers (H1-H6)
- Generate a hierarchical table of contents
- Support both single file and directory modes
- Output TOC in markdown format with proper indentation and links
- Provide options for TOC depth control and formatting

## Technical Design

### Binary Location
- **Path**: `cmd/md-toc/main.go`
- **Build output**: `build/md-toc`
- **Build script**: Add entry to `scripts/build.sh`

### Command Structure

```bash
# Single file
md-toc README.md

# Directory (all markdown files)
md-toc ./docs

# With options
md-toc --depth 3 --output toc.md ./docs
md-toc --recursive --format json README.md
```

### CLI Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--depth` | `-d` | `6` | Maximum heading depth (1-6) |
| `--output` | `-o` | `stdout` | Output file path (default: print to stdout) |
| `--format` | `-f` | `markdown` | Output format: `markdown`, `json`, `text` |
| `--recursive` | `-r` | `true` | Process subdirectories |
| `--prepend` | `-p` | `false` | Prepend TOC to original file(s) |
| `--marker` | `-m` | `<!-- TOC -->` | TOC marker for insertion point |
| `--anchor-style` | `-a` | `github` | Anchor style: `github`, `gitlab`, `custom` |
| `--exclude` | `-e` | `nil` | Exclude pattern (glob) |

### Core Features

#### 1. Single File Mode
- Parse single markdown file
- Extract headers (H1-H6)
- Generate TOC with:
  - Proper indentation (2 spaces per level)
  - GitHub-style anchor links
  - Numbering options (optional)

#### 2. Directory Mode
- Find all `.md` files recursively
- Generate individual TOC per file, OR
- Generate combined TOC for all files with file sections

#### 3. TOC Formats

**Markdown Format** (default):
```markdown
## Table of Contents

- [Chapter 1](#chapter-1)
  - [Section 1.1](#section-11)
  - [Section 1.2](#section-12)
- [Chapter 2](#chapter-2)
```

**JSON Format**:
```json
{
  "files": [
    {
      "path": "README.md",
      "toc": [
        {"level": 1, "title": "Chapter 1", "anchor": "#chapter-1"},
        {"level": 2, "title": "Section 1.1", "anchor": "#section-11"}
      ]
    }
  ]
}
```

**Text Format**:
```
Chapter 1
  Section 1.1
  Section 1.2
Chapter 2
```

### Implementation Components

#### Package Structure
```
cmd/md-toc/
  main.go              # CLI entry point, cobra setup

internal/toc/
  parser.go            # Markdown parsing, header extraction
  generator.go         # TOC generation logic
  formatter.go         # Output formatting (markdown/json/text)
  anchor.go            # Anchor link generation
  file_walker.go       # Directory traversal
```

#### Core Functions

```go
// Parser
type Header struct {
    Level   int
    Title   string
    Anchor  string
    LineNum int
}

func ParseMarkdown(content string) ([]Header, error)
func ExtractHeaders(filepath string) ([]Header, error)

// Generator
type TOC struct {
    Headers []Header
    Depth   int
}

func GenerateTOC(headers []Header, depth int) string
func GenerateTOCForFile(filepath string, opts Options) (string, error)
func GenerateTOCForDir(dirpath string, opts Options) (map[string]string, error)

// Formatter
func FormatAsMarkdown(toc TOC) string
func FormatAsJSON(toc TOC) string
func FormatAsText(toc TOC) string

// Anchor generation
func GenerateAnchor(title string, style AnchorStyle) string
```

### Dependencies

- **github.com/spf13/cobra**: CLI framework (already in project)
- **github.com/gomarkdown/markdown**: Markdown parsing (or use regex for simple parsing)
- Standard library: `os`, `path/filepath`, `strings`, `encoding/json`

### Build Integration

Add to `scripts/build.sh`:
```bash
echo "Building md-toc..."
go build -o build/md-toc cmd/md-toc/main.go
```

### Testing Strategy

1. **Unit Tests**
   - Header extraction from markdown
   - Anchor generation (various styles)
   - TOC formatting
   - Edge cases: empty files, no headers, special chars

2. **Integration Tests**
   - Single file processing
   - Directory traversal
   - Output formats
   - File writing with --prepend option

3. **Test Data**
   - Create `internal/toc/testdata/` with sample markdown files
   - Various header levels and formats
   - Special characters in headers
   - Edge cases

### Future Enhancements (Not in MVP)

- Update existing TOC in-place (detect and replace between markers)
- Custom numbering schemes (1.1, 1.1.1, etc.)
- Ignore certain sections (comments, code blocks)
- Multi-language support (RTL languages)
- Watch mode (regenerate on file changes)

## Implementation Tasks

- [ ] Create package structure: `internal/toc/`
- [ ] Implement markdown parser with header extraction
- [ ] Implement anchor generation (GitHub style)
- [ ] Implement TOC generator (markdown format)
- [ ] Create CLI with cobra in `cmd/md-toc/main.go`
- [ ] Add single file mode
- [ ] Add directory mode with recursive traversal
- [ ] Implement JSON format output
- [ ] Implement text format output
- [ ] Add --prepend option (insert TOC into file)
- [ ] Add --depth flag
- [ ] Add --exclude pattern support
- [ ] Update build script
- [ ] Write unit tests
- [ ] Write integration tests
- [ ] Test on actual project docs
- [ ] Documentation in README

## Example Usage

```bash
# Generate TOC for project documentation
./build/md-toc --depth 3 --output docs/TOC.md ./docs/work

# Generate and prepend TOC to README
./build/md-toc --prepend --marker "<!-- TOC -->" README.md

# Generate JSON TOC for all markdown files
./build/md-toc --format json --recursive ./ddia-quiz-bot/content

# Quick TOC to stdout
./build/md-toc docs/work/014-*.md
```

## Success Criteria

1. Successfully parses markdown headers from files
2. Generates properly formatted TOC with correct indentation
3. Creates valid GitHub-style anchor links
4. Works with both single files and directories
5. Outputs in multiple formats (markdown, json, text)
6. Handles edge cases gracefully (no headers, special chars, etc.)
7. Integration with build script
8. Comprehensive test coverage (>80%)

## Notes

- Keep it simple - focus on common use cases first
- GitHub anchor style is most compatible
- Consider using simple regex parsing instead of full markdown parser for MVP
- Make it fast - should handle large directories quickly
- Follow project conventions (cobra CLI, internal packages)
