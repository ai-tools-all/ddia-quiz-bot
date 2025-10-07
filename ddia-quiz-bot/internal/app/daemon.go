package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/your-username/ddia-quiz-bot/internal/config"
	"github.com/your-username/ddia-quiz-bot/internal/discovery"
	"github.com/your-username/ddia-quiz-bot/internal/models"
	"github.com/your-username/ddia-quiz-bot/internal/notifier"
	"github.com/your-username/ddia-quiz-bot/internal/presenter"
	"github.com/your-username/ddia-quiz-bot/internal/state"
	"github.com/your-username/ddia-quiz-bot/internal/store"
)

const reloadDebounceDuration = 2 * time.Second

type Daemon struct {
	contentPath    string
	logger         *log.Logger
	notifiers      []notifier.Notifier
	state          *state.Manager
	presenter      *presenter.SocialPresenter
	postHour       int
	postMinute     int
	reloadTimer    *time.Timer
	reloadTimerMux sync.Mutex

	// These are protected by the mutex
	store    *store.ContentStore
	schedule *models.Schedule
	mutex    sync.RWMutex
}

// NewDaemon initializes the application.
func NewDaemon(contentPath string, logger *log.Logger, notifiers []notifier.Notifier, stateMgr *state.Manager, postHour, postMinute int) *Daemon {
	return &Daemon{
		contentPath: contentPath,
		logger:      logger,
		notifiers:   notifiers,
		state:       stateMgr,
		presenter:   &presenter.SocialPresenter{},
		postHour:    postHour,
		postMinute:  postMinute,
	}
}

// Run starts all long-running processes.
func (d *Daemon) Run() error {
	d.logger.Println("Starting DDIA Quiz Daemon...")

	// 1. Initial Load
	if err := d.Reload(); err != nil {
		return fmt.Errorf("initial load failed: %w", err)
	}
	d.logger.Println("Initial content loaded successfully.")

	// 2. Start Filesystem Watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}
	defer watcher.Close()

	if err := filepath.Walk(d.contentPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to add paths to watcher: %w", err)
	}

	// 3. Start Scheduler Ticker and Signal Handler
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	d.logger.Println("Daemon is running. Watching for file changes and scheduled posts.")

	// 4. Main Event Loop
	for {
		select {
		case <-ticker.C:
			d.checkForScheduledPosts()
		case event := <-watcher.Events:
			d.handleFileChange(event)
		case err := <-watcher.Errors:
			d.logger.Printf("ERROR: Watcher error: %v", err)
		case <-signalChan:
			d.logger.Println("Shutdown signal received. Saving state and exiting.")
			if err := d.state.Save(); err != nil {
				d.logger.Printf("ERROR: Failed to save state on shutdown: %v", err)
			}
			return nil
		}
	}
}

// handleFileChange debounces reload events.
func (d *Daemon) handleFileChange(event fsnotify.Event) {
	if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
		d.reloadTimerMux.Lock()
		// If a timer is already running, reset it
		if d.reloadTimer != nil {
			d.reloadTimer.Stop()
		}
		// Set a new timer
		d.reloadTimer = time.AfterFunc(reloadDebounceDuration, func() {
			if err := d.Reload(); err != nil {
				d.logger.Printf("ERROR: Automatic reload failed: %v", err)
			}
		})
		d.reloadTimerMux.Unlock()
	}
}

// Reload attempts to load all config and content from disk.
func (d *Daemon) Reload() error {
	d.logger.Println("Change detected. Attempting to reload schedule and content...")

	newSchedule, err := config.LoadSchedule(filepath.Join(d.contentPath, "schedule.yml"))
	if err != nil {
		return fmt.Errorf("failed to parse schedule.yml: %w", err)
	}

	newStore, err := store.NewContentStore(d.contentPath)
	if err != nil {
		return fmt.Errorf("failed to load content store: %w", err)
	}

	// Validation (optional but recommended)
	// You could add a `Validate(schedule, store)` function here

	d.mutex.Lock()
	d.schedule = newSchedule
	d.store = newStore
	d.mutex.Unlock()

	d.logger.Println("Reload successful. Now running with updated configuration.")
	return nil
}

func (d *Daemon) checkForScheduledPosts() {
	d.mutex.RLock()
	schedule := d.schedule
	store := d.store
	d.mutex.RUnlock()

	if schedule == nil || store == nil {
		return // Not loaded yet
	}

	matcher := &discovery.Matcher{Store: store}
	now := time.Now().UTC()
	needsSave := false

	for _, chap := range schedule.Chapters {
		startDate, err := chap.GetStartDate()
		if err != nil {
			d.logger.Printf("WARN: Invalid start_date for chapter %d: %v", chap.Chapter, err)
			continue
		}

		for _, qSched := range chap.Questions {
			question, exists := store.QuestionsByID[qSched.File]
			if !exists {
				// The question file might have a different ID in its frontmatter, let's find it
				var foundQ *models.Question
				for _, q := range store.QuestionsByID {
					if filepath.Base(q.Path) == qSched.File {
						foundQ = q
						break
					}
				}
				if foundQ == nil {
					d.logger.Printf("WARN: Question file '%s' from schedule not found in store.", qSched.File)
					continue
				}
				question = foundQ
			}

			if d.state.HasPosted(question.ID) {
				continue
			}

			// Post time is start date + (day-1) days, at the configured time UTC
			postDate := startDate.AddDate(0, 0, qSched.Day-1)
			postTime := time.Date(postDate.Year(), postDate.Month(), postDate.Day(), d.postHour, d.postMinute, 0, 0, time.UTC)

			if now.After(postTime) {
				d.logger.Printf("Found pending post: %s", question.ID)

				stories, err := matcher.FindStoriesForQuestion(question.ID, &qSched)
				if err != nil {
					d.logger.Printf("ERROR: Could not find stories for %s: %v", question.ID, err)
					continue
				}

				dailyPost := &presenter.DailyPost{
					Question: question,
					Stories:  stories,
					Date:     now,
				}

				content, err := d.presenter.Format(dailyPost)
				if err != nil {
					d.logger.Printf("ERROR: Could not format post for %s: %v", question.ID, err)
					continue
				}

				// Send to all notifiers
				for _, n := range d.notifiers {
					if err := n.Notify(content); err != nil {
						d.logger.Printf("ERROR: Notifier failed for %s: %v", question.ID, err)
						// Decide on retry logic here if needed
					}
				}

				d.state.MarkAsPosted(question.ID)
				needsSave = true
			}
		}
	}

	if needsSave {
		if err := d.state.Save(); err != nil {
			d.logger.Printf("ERROR: Failed to save state after posting: %v", err)
		}
	}
}
