package store

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/your-username/ddia-quiz-bot/internal/models"
	"github.com/your-username/ddia-quiz-bot/internal/parser"
)

// ContentStore acts as an in-memory database for all parsed quiz content.
type ContentStore struct {
	QuestionsByID     map[string]*models.Question
	StoriesByFilename map[string]*models.Story
	AllStories        []*models.Story
}

// NewContentStore creates and populates a store by scanning a content directory.
func NewContentStore(contentDir string) (*ContentStore, error) {
	store := &ContentStore{
		QuestionsByID:     make(map[string]*models.Question),
		StoriesByFilename: make(map[string]*models.Story),
		AllStories:        []*models.Story{},
	}

	if err := store.loadStories(filepath.Join(contentDir, "stories")); err != nil {
		return nil, err
	}

	if err := store.loadChapters(filepath.Join(contentDir, "chapters")); err != nil {
		return nil, err
	}

	return store, nil
}

func (cs *ContentStore) loadStories(storiesDir string) error {
	return filepath.Walk(storiesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			yamlData, mdData, err := parser.ReadAndParseFile(path)
			if err != nil {
				return err // Or log and continue
			}

			var story models.Story
			if err := parser.UnmarshalFrontmatter(yamlData, &story); err != nil {
				return err // Or log and continue
			}

			story.Path = path
			story.Filename = strings.TrimSuffix(info.Name(), ".md")
			story.ContentMarkdown = string(mdData)

			cs.StoriesByFilename[story.Filename] = &story
			cs.AllStories = append(cs.AllStories, &story)
		}
		return nil
	})
}

func (cs *ContentStore) loadChapters(chaptersDir string) error {
	return filepath.Walk(chaptersDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") && info.Name() != "_meta.md" {
			yamlData, mdData, err := parser.ReadAndParseFile(path)
			if err != nil {
				return err
			}

			var question models.Question
			if err := parser.UnmarshalFrontmatter(yamlData, &question); err != nil {
				return err
			}

			question.Path = path
			question.ChapterID = filepath.Base(filepath.Dir(path))
			parser.ParseQuestionSections(mdData, &question) // Populate structured fields

			cs.QuestionsByID[question.ID] = &question
		}
		return nil
	})
}
