package gui

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"sync"
)

// HistoryEntry represents a completed download in history
type HistoryEntry struct {
	ID          string         `json:"id"`
	URL         string         `json:"url"`
	Domain      string         `json:"domain"`
	Output      string         `json:"output"`
	StartedAt   string         `json:"startedAt"`
	CompletedAt string         `json:"completedAt,omitempty"`
	Downloaded  int64          `json:"downloaded"`
	Errors      int64          `json:"errors"`
	Status      string         `json:"status"`
	Config      DownloadConfig `json:"config"`
}

// HistoryManager manages download history persistence
type HistoryManager struct {
	filePath string
	entries  []HistoryEntry
	mu       sync.RWMutex
}

// NewHistoryManager creates a new history manager
func NewHistoryManager(filePath string) *HistoryManager {
	hm := &HistoryManager{
		filePath: filePath,
		entries:  []HistoryEntry{},
	}
	hm.load()
	return hm
}

// GetHistory returns all history entries
func (a *App) GetHistory() []HistoryEntry {
	return a.historyManager.GetAll()
}

// ClearHistory clears all history entries
func (a *App) ClearHistory() error {
	return a.historyManager.Clear()
}

// RerunFromHistory loads config from a history entry
func (a *App) RerunFromHistory(id string) (DownloadConfig, error) {
	entry, err := a.historyManager.Get(id)
	if err != nil {
		return DownloadConfig{}, err
	}
	return entry.Config, nil
}

// GetAll returns all history entries
func (hm *HistoryManager) GetAll() []HistoryEntry {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	// Return in reverse order (newest first)
	result := make([]HistoryEntry, len(hm.entries))
	for i, entry := range hm.entries {
		result[len(hm.entries)-1-i] = entry
	}
	return result
}

// Get returns a history entry by ID
func (hm *HistoryManager) Get(id string) (HistoryEntry, error) {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	for _, entry := range hm.entries {
		if entry.ID == id {
			return entry, nil
		}
	}
	return HistoryEntry{}, os.ErrNotExist
}

// Add adds a new history entry
func (hm *HistoryManager) Add(entry HistoryEntry) error {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	hm.entries = append(hm.entries, entry)

	// Keep only last 100 entries
	if len(hm.entries) > 100 {
		hm.entries = hm.entries[len(hm.entries)-100:]
	}

	return hm.save()
}

// Clear removes all history entries
func (hm *HistoryManager) Clear() error {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	hm.entries = []HistoryEntry{}
	return hm.save()
}

// load loads history from file
func (hm *HistoryManager) load() error {
	data, err := os.ReadFile(hm.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &hm.entries)
}

// save saves history to file
func (hm *HistoryManager) save() error {
	data, err := json.MarshalIndent(hm.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(hm.filePath, data, 0644)
}

// generateID generates a random ID
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
