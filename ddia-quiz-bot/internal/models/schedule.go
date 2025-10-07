package models

import "time"

// Schedule represents the structure of schedule.yml.
type Schedule struct {
	Chapters []ChapterSchedule `yaml:"chapters"`
}

type ChapterSchedule struct {
	Chapter   int                `yaml:"chapter"`
	Path      string             `yaml:"path"`
	StartDate string             `yaml:"start_date"`
	Questions []QuestionSchedule `yaml:"questions"`
}

func (cs *ChapterSchedule) GetStartDate() (time.Time, error) {
	return time.Parse("2006-01-02", cs.StartDate)
}

type QuestionSchedule struct {
	Day            int      `yaml:"day"`
	File           string   `yaml:"file"`
	IncludeStories []string `yaml:"include_stories,omitempty"`
	ExcludeStories []string `yaml:"exclude_stories,omitempty"`
}
