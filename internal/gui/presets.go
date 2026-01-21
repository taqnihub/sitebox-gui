package gui

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Preset represents a saved download configuration
type Preset struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	CreatedAt string         `json:"createdAt"`
	Config    DownloadConfig `json:"config"`
}

// PresetManager manages preset persistence
type PresetManager struct {
	filePath string
	presets  []Preset
	mu       sync.RWMutex
}

// NewPresetManager creates a new preset manager
func NewPresetManager(filePath string) *PresetManager {
	pm := &PresetManager{
		filePath: filePath,
		presets:  []Preset{},
	}
	pm.load()
	return pm
}

// GetPresets returns all presets
func (a *App) GetPresets() []Preset {
	return a.presetManager.GetAll()
}

// SavePreset saves a new preset
func (a *App) SavePreset(name string, config DownloadConfig) error {
	return a.presetManager.Save(name, config)
}

// LoadPreset returns a preset by ID
func (a *App) LoadPreset(id string) (DownloadConfig, error) {
	preset, err := a.presetManager.Get(id)
	if err != nil {
		return DownloadConfig{}, err
	}
	return preset.Config, nil
}

// DeletePreset deletes a preset by ID
func (a *App) DeletePreset(id string) error {
	return a.presetManager.Delete(id)
}

// GetAll returns all presets
func (pm *PresetManager) GetAll() []Preset {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	result := make([]Preset, len(pm.presets))
	copy(result, pm.presets)
	return result
}

// Get returns a preset by ID
func (pm *PresetManager) Get(id string) (Preset, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	for _, preset := range pm.presets {
		if preset.ID == id {
			return preset, nil
		}
	}
	return Preset{}, os.ErrNotExist
}

// Save saves a new preset
func (pm *PresetManager) Save(name string, config DownloadConfig) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	preset := Preset{
		ID:        generateID(),
		Name:      name,
		CreatedAt: time.Now().Format(time.RFC3339),
		Config:    config,
	}

	pm.presets = append(pm.presets, preset)
	return pm.save()
}

// Delete removes a preset by ID
func (pm *PresetManager) Delete(id string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for i, preset := range pm.presets {
		if preset.ID == id {
			pm.presets = append(pm.presets[:i], pm.presets[i+1:]...)
			return pm.save()
		}
	}
	return os.ErrNotExist
}

// load loads presets from file
func (pm *PresetManager) load() error {
	data, err := os.ReadFile(pm.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &pm.presets)
}

// save saves presets to file
func (pm *PresetManager) save() error {
	data, err := json.MarshalIndent(pm.presets, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(pm.filePath, data, 0644)
}
