Of course. Evolving the tool from a one-shot CLI into a long-running, stateful daemon that intelligently reloads on-the-fly is an excellent next step. This requires a more robust architecture focused on state management, concurrent operations, and filesystem monitoring.

Here is a detailed coat plan for building this daemon in Go.

### **Part 1: Architectural Vision (Daemon Model)**

1.  **Stateful & Persistent**: The daemon's primary job is to remember what it has posted. This state must survive restarts. We will use a simple, local JSON file for this (`post_history.json`), which is robust and easy to inspect.
2.  **Watch, Don't Poll**: The daemon will use a filesystem watcher to get instant notifications when files change. This is far more efficient than constantly scanning the directories.
3.  **Graceful Reloads**: When a file changes, the daemon will attempt to reload the entire content and schedule *in memory*. If the new configuration is invalid (e.g., bad YAML syntax), it will discard the changes, log a detailed error, and continue running with the last known good configuration. This ensures a single typo doesn't crash the service.
4.  **Concurrent & Safe**: The filesystem watcher and the main scheduler will run in separate goroutines. All access to shared data (like the current schedule and content store) must be protected by mutexes to prevent race conditions.
5.  **Scheduler-Driven**: A central ticker will "wake up" the daemon every minute (or a configurable interval) to check if any scheduled posts are due.
6.  **Pluggable Notifiers**: The core logic will prepare the content to be posted. The actual sending (to Telegram, WhatsApp, etc.) will be handled by a "Notifier" interface, allowing us to easily add new platforms without changing the core daemon.

### **Part 2: Updated Project Structure**

We'll add new packages for the daemon's specific responsibilities and a new entrypoint.

```
ddia-quiz-bot/
├── cmd/
│   ├── quiz/           # The original CLI tool (still useful for manual validation)
│   │   └── main.go
│   └── quiz-daemon/    # NEW: The long-running service entrypoint
│       └── main.go
│
├── internal/
│   ├── app/            # Main application orchestrator (the Daemon struct)
│   ├── cli/
│   ├── config/
│   ├── discovery/
│   ├── models/
│   ├── notifier/       # NEW: Interface and implementations for sending posts
│   ├── parser/
│   ├── presenter/
│   ├── state/          # NEW: Manages the post_history.json file
│   └── store/
│
├── content/            # The content directory the daemon will watch
│   ├── ...
│   ├── schedule.yml
│   └── post_history.json # NEW: Created and managed by the daemon
│
├── go.mod
└── go.sum
```

### **Part 3: New & Updated Core Packages**

#### **1. `internal/state` (State Management)**

This package is responsible for persistence.

*   **`post_history.json` format**:
    ```json
    {
      "ch3_kv_01": "2025-10-08T10:00:00Z",
      "ch3_compaction_02": "2025-10-09T10:00:00Z"
    }
    ```
*   **`manager.go`**:
    ```go
    package state

    import "sync"

    // Manager handles the state of posted items. It is safe for concurrent use.
    type Manager struct {
        filepath   string
        posted     map[string]time.Time // Question ID -> Post Timestamp
        mutex      sync.RWMutex
    }

    // NewManager creates a manager and loads initial state from disk.
    func NewManager(filepath string) (*Manager, error)

    // Load reads the state file from disk.
    func (m *Manager) Load() error

    // Save writes the current state to disk.
    func (m *Manager) Save() error

    // HasPosted checks if a question ID has already been posted.
    func (m *Manager) HasPosted(questionID string) bool

    // MarkAsPosted records that a question has been posted.
    func (m *Manager) MarkAsPosted(questionID string)
    ```

#### **2. `internal/notifier` (Pluggable Notifiers)**

This package decouples the posting logic from the destination platform.

*   **`notifier.go`**:
    ```go
    package notifier

    // Notifier is the interface for sending a formatted post to a platform.
    type Notifier interface {
        Notify(content string) error
    }
    ```
*   **Implementations**:
    *   `log_notifier.go`: A simple implementation that just prints the content to the console. Perfect for development and testing.
    *   `telegram_notifier.go`: (Future) An implementation that sends the content to a Telegram channel via its API.
    *   `whatsapp_notifier.go`: (Future) An implementation for WhatsApp.

#### **3. `internal/app` (The Daemon Orchestrator)**

This is the new heart of the application. It ties everything together.

*   **`daemon.go`**:
    ```go
    package app

    import (
        "log"
        "sync"
        "time"
    )

    type Daemon struct {
        contentPath string
        logger      *log.Logger
        notifiers   []notifier.Notifier
        state       *state.Manager

        // These are protected by the mutex
        store    *store.ContentStore
        schedule *models.Schedule
        mutex    sync.RWMutex
    }

    // NewDaemon initializes the application.
    func NewDaemon(...) (*Daemon, error)

    // Run starts all long-running processes: the file watcher and the scheduler.
    // It blocks until a shutdown signal is received.
    func (d *Daemon) Run() error

    // handleFileChange is triggered by the watcher. It debounces events
    // and then calls Reload.
    func (d *Daemon) handleFileChange(event fsnotify.Event)

    // Reload attempts to load all config and content from disk. If successful,
    // it swaps the in-memory data. If not, it logs the error and keeps the old data.
    func (d *Daemon) Reload() error

    // checkForScheduledPosts is called by the ticker. It checks the schedule
    // against the state and posts anything that is due.
    func (d *Daemon) checkForScheduledPosts()
    ```

### **Part 4: Execution Flow**

#### **Startup (`cmd/quiz-daemon/main.go`)**

1.  Initialize logger to output to `stdout`.
2.  Get the content path from a command-line flag (e.g., `./quiz-daemon --path ./content`).
3.  Create the state manager: `state.NewManager("./content/post_history.json")`.
4.  Create the notifiers (e.g., a `notifier.LogNotifier`).
5.  Create the main daemon instance: `app.NewDaemon(...)`.
6.  Call `daemon.Run()`.

#### **Inside `daemon.Run()`**

1.  **Initial Load**: Call `d.Reload()` immediately to perform the first load. If it fails, exit with an error. If successful, log "Daemon started. Initial content loaded successfully."
2.  **Start Filesystem Watcher**:
    *   Use the `fsnotify` library to create a new watcher.
    *   Recursively add the entire `contentPath` to the watch list.
    *   Start a new goroutine that listens for events on the watcher's channel.
    *   This goroutine will call `d.handleFileChange(event)`.
3.  **Start Scheduler Ticker**:
    *   `ticker := time.NewTicker(1 * time.Minute)`
4.  **Start Signal Handler**:
    *   Listen for `os.Interrupt` (Ctrl+C) and `syscall.SIGTERM` signals.
5.  **Main `select` Loop**:
    ```go
    for {
        select {
        case <-ticker.C:
            d.checkForScheduledPosts()
        case event := <-watcher.Events: // Simplified for clarity
            d.handleFileChange(event)
        case err := <-watcher.Errors:
            d.logger.Printf("ERROR: Watcher error: %v", err)
        case <-signalChan:
            d.logger.Println("Shutdown signal received. Saving state and exiting.")
            d.state.Save()
            return nil
        }
    }
    ```

#### **Scenario: User updates `schedule.yml`**

1.  The filesystem watcher detects a `WRITE` event for `content/schedule.yml`.
2.  The watcher's goroutine receives the event. To prevent reloading multiple times for a single save, it "debounces" by waiting ~2 seconds. If another event comes in, the timer resets.
3.  After the debounce period, it calls `d.Reload()`.
4.  **Inside `d.Reload()`**:
    *   `d.mutex.Lock()` // Prevent the scheduler from running during the reload.
    *   `defer d.mutex.Unlock()`
    *   `d.logger.Println("Change detected. Attempting to reload schedule and content...")`
    *   It tries to parse the new `schedule.yml`.
        *   **FAILURE**: The YAML is malformed. The function logs a detailed error: `d.logger.Println("ERROR: Failed to reload schedule.yml: yaml: line 15: could not find expected ':'")`. It then `return`s, leaving the old, valid schedule in memory.
        *   **SUCCESS**: The new schedule is parsed successfully.
    *   It then re-loads the entire content store (`store.NewContentStore(...)`) to catch any new or changed `.md` files.
    *   It runs the full validation logic (e.g., checking that all stories in the new schedule exist).
    *   If everything is valid, it replaces the old pointers: `d.schedule = newSchedule` and `d.store = newStore`.
    *   It logs a success message: `d.logger.Println("Reload successful. Now running with updated configuration.")`

#### **Scenario: The Scheduler Ticks**

1.  The `select` loop triggers `d.checkForScheduledPosts()`.
2.  **Inside `checkForScheduledPosts()`**:
    *   `d.mutex.RLock()` // Get a read-lock. Multiple schedulers could (in theory) run at once.
    *   `defer d.mutex.RUnlock()`
    *   It iterates through the current `d.schedule`.
    *   For each question, it calculates the `postDateTime`.
    *   It checks two conditions:
        1.  `time.Now().After(postDateTime)`
        2.  `!d.state.HasPosted(question.ID)`
    *   If both are true, it has found a pending post.
    *   It uses the `discovery.Matcher` to find related stories.
    *   It uses a `presenter` to format the final post content.
    *   It iterates through `d.notifiers` and calls `notifier.Notify(content)`.
    *   If `Notify` returns no error, it calls `d.state.MarkAsPosted(question.ID)`.
    *   After the loop, it calls `d.state.Save()` to ensure the history is written to disk immediately.