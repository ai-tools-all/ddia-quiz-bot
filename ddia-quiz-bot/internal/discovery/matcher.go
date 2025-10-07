package discovery

import (
	"fmt"

	"github.com/your-username/ddia-quiz-bot/internal/models"
	"github.com/your-username/ddia-quiz-bot/internal/store"
)

// Matcher finds relationships between stories and questions.
type Matcher struct {
	Store *store.ContentStore
}

// FindStoriesForQuestion implements the two-way linking logic.
func (m *Matcher) FindStoriesForQuestion(
	questionID string,
	scheduleRules *models.QuestionSchedule,
) ([]*models.Story, error) {
	question, exists := m.Store.QuestionsByID[questionID]
	if !exists {
		return nil, fmt.Errorf("question with ID '%s' not found in store", questionID)
	}

	storiesFound := make(map[string]*models.Story)

	// 1. Story -> Question (story's `relates_to` points to this question)
	for _, story := range m.Store.AllStories {
		for _, relatesToID := range story.RelatesTo {
			if relatesToID == questionID {
				storiesFound[story.Filename] = story
			}
		}
	}

	// 2. Question -> Story (question's `related_stories` points to a story)
	for _, storyFilename := range question.RelatedStories {
		if story, ok := m.Store.StoriesByFilename[storyFilename]; ok {
			storiesFound[story.Filename] = story
		}
	}

	// 3. Schedule overrides
	if scheduleRules != nil {
		// Include stories
		for _, storyFilename := range scheduleRules.IncludeStories {
			if story, ok := m.Store.StoriesByFilename[storyFilename]; ok {
				storiesFound[story.Filename] = story
			}
		}
		// Exclude stories
		for _, storyFilename := range scheduleRules.ExcludeStories {
			delete(storiesFound, storyFilename)
		}
	}

	// Convert map to slice
	result := make([]*models.Story, 0, len(storiesFound))
	for _, story := range storiesFound {
		result = append(result, story)
	}
	return result, nil
}
