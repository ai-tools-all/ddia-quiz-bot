package state

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

// Manager handles the persistent state of posted items.
type Manager struct {
	filepath string
	posted   map[string]time.Time // Question ID -> Post Timestamp
	mutex    sync.RWMutex
}

// NewManager creates a manager and loads initial state from disk.
func NewManager(filepath string) (*Manager, error) {
	m := &Manager{
		filepath: filepath,
		posted:   make(map[string]time.Time),
	}
	if err := m.Load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return m, nil
}

// Load reads the state file from disk.
func (m *Manager) Load() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	data, err := ioutil.ReadFile(m.filepath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &m.posted)
}

// Save writes the current state to disk atomically.
func (m *Manager) Save() error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	data, err := json.MarshalIndent(m.posted, "", "  ")
	if err != nil {
		return err
	}
	// Write to a temporary file first, then rename for atomicity
	tempFile := m.filepath + ".tmp"
	if err := ioutil.WriteFile(tempFile, data, 0644); err != nil {
		return err
	}
	return os.Rename(tempFile, m.filepath)
}

// HasPosted checks if a question ID has already been posted.
func (m *Manager) HasPosted(questionID string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	_, exists := m.posted[questionID]
	return exists
}

// MarkAsPosted records that a question has been posted.
func (m *Manager) MarkAsPosted(questionID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.posted[questionID] = time.Now().UTC()
}
