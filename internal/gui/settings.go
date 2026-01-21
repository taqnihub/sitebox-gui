package gui

import (
	"encoding/json"
	"os"
	"sync"
)

// Settings represents user preferences
type Settings struct {
	Theme             string `json:"theme"`
	DefaultOutput     string `json:"defaultOutput"`
	DefaultConcurrent int    `json:"defaultConcurrent"`
	DefaultMaxDepth   int    `json:"defaultMaxDepth"`
	DefaultRetries    int    `json:"defaultRetries"`
	IncludeImages     bool   `json:"includeImages"`
	ShowNotifications bool   `json:"showNotifications"`
}

// DefaultSettings returns default settings
func DefaultSettings() Settings {
	return Settings{
		Theme:             "system",
		DefaultOutput:     "",
		DefaultConcurrent: 5,
		DefaultMaxDepth:   50,
		DefaultRetries:    3,
		IncludeImages:     false,
		ShowNotifications: true,
	}
}

// SettingsManager manages settings persistence
type SettingsManager struct {
	filePath string
	settings Settings
	mu       sync.RWMutex
}

// NewSettingsManager creates a new settings manager
func NewSettingsManager(filePath string) *SettingsManager {
	sm := &SettingsManager{
		filePath: filePath,
		settings: DefaultSettings(),
	}
	sm.load()
	return sm
}

// GetSettings returns the current settings
func (a *App) GetSettings() Settings {
	return a.settingsManager.Get()
}

// SaveSettings saves settings
func (a *App) SaveSettings(settings Settings) error {
	return a.settingsManager.Save(settings)
}

// GetDefaultConfig returns default download config based on settings
func (a *App) GetDefaultConfig() DownloadConfig {
	settings := a.settingsManager.Get()
	return DownloadConfig{
		Output:          settings.DefaultOutput,
		MaxDepth:        settings.DefaultMaxDepth,
		Concurrent:      settings.DefaultConcurrent,
		Retries:         settings.DefaultRetries,
		IncludeImages:   settings.IncludeImages,
		Blacklist:       []string{"https://github"},
		Extensions:      []string{".js", ".css", ".woff2", ".woff", ".ttf"},
		ImageExtensions: []string{".png", ".svg", ".jpg", ".jpeg", ".gif", ".webp"},
	}
}

// Get returns the current settings
func (sm *SettingsManager) Get() Settings {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.settings
}

// Save saves settings
func (sm *SettingsManager) Save(settings Settings) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.settings = settings
	return sm.save()
}

// load loads settings from file
func (sm *SettingsManager) load() error {
	data, err := os.ReadFile(sm.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &sm.settings)
}

// save saves settings to file
func (sm *SettingsManager) save() error {
	data, err := json.MarshalIndent(sm.settings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(sm.filePath, data, 0644)
}
