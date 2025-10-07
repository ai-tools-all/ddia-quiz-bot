package presenter

import (
	"bytes"
	"text/template"
	"time"

	"github.com/your-username/ddia-quiz-bot/internal/models"
)

// DailyPost is a generic struct containing everything needed for a day's output.
type DailyPost struct {
	Question *models.Question
	Stories  []*models.Story
	Date     time.Time
}

// SocialPresenter formats output for a social media channel.
type SocialPresenter struct{}

const socialTemplate = `🔥 DDIA {{.Question.ChapterID}} - Day {{.Question.Day}}

{{.Question.QuestionText}}

---
{{if .Stories}}
📖 Real-World Example(s):
{{range .Stories}}
🎬 {{.Title}}
{{.EngagementHook}}
Read more: [link to {{.Filename}} story]
{{end}}
#DDIA #RealWorld
{{end}}`

func (p *SocialPresenter) Format(post *DailyPost) (string, error) {
	tmpl, err := template.New("social").Parse(socialTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, post); err != nil {
		return "", err
	}

	return buf.String(), nil
}
