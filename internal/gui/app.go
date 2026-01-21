package gui

import (
	"context"
	"os"
	"path/filepath"
	goruntime "runtime"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	configDir      string
	downloadManager *DownloadManager
	historyManager  *HistoryManager
	presetManager   *PresetManager
	settingsManager *SettingsManager
}

// NewApp creates a new App instance
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize config directory
	configDir, err := getConfigDir()
	if err != nil {
		wailsRuntime.LogError(ctx, "Failed to get config directory: "+err.Error())
		configDir = "."
	}
	a.configDir = configDir

	// Ensure config directory exists
	if err := os.MkdirAll(a.configDir, 0755); err != nil {
		wailsRuntime.LogError(ctx, "Failed to create config directory: "+err.Error())
	}

	// Initialize managers
	a.historyManager = NewHistoryManager(filepath.Join(a.configDir, "history.json"))
	a.presetManager = NewPresetManager(filepath.Join(a.configDir, "presets.json"))
	a.settingsManager = NewSettingsManager(filepath.Join(a.configDir, "settings.json"))
	a.downloadManager = NewDownloadManager(ctx, a.historyManager)

	wailsRuntime.LogInfo(ctx, "SiteBox GUI started, config dir: "+a.configDir)
}

// Shutdown is called when the app is closing
func (a *App) Shutdown(ctx context.Context) {
	// Stop any active downloads
	if a.downloadManager != nil {
		a.downloadManager.Stop()
	}
	wailsRuntime.LogInfo(ctx, "SiteBox GUI shutting down")
}

// SelectOutputDirectory opens a directory picker dialog
func (a *App) SelectOutputDirectory() (string, error) {
	return wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Select Output Directory",
	})
}

// OpenFolder opens a folder in the file manager
func (a *App) OpenFolder(path string) {
	wailsRuntime.BrowserOpenURL(a.ctx, "file://"+path)
}

// getConfigDir returns the platform-specific config directory
func getConfigDir() (string, error) {
	var configDir string

	switch goruntime.GOOS {
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(home, "Library", "Application Support", "sitebox")
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		configDir = filepath.Join(appData, "sitebox")
	default: // linux and others
		configHome := os.Getenv("XDG_CONFIG_HOME")
		if configHome == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			configHome = filepath.Join(home, ".config")
		}
		configDir = filepath.Join(configHome, "sitebox")
	}

	return configDir, nil
}
