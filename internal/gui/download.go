package gui

import (
	"context"
	"sync"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"sitebox/pkg/downloader"
)

// DownloadConfig represents the download configuration from the frontend
type DownloadConfig struct {
	URL             string   `json:"url"`
	Domain          string   `json:"domain"`
	Output          string   `json:"output"`
	MaxDepth        int      `json:"maxDepth"`
	Concurrent      int      `json:"concurrent"`
	Retries         int      `json:"retries"`
	IncludeImages   bool     `json:"includeImages"`
	Blacklist       []string `json:"blacklist"`
	Extensions      []string `json:"extensions"`
	ImageExtensions []string `json:"imageExtensions"`
}

// ProgressEvent represents a download progress event sent to the frontend
type ProgressEvent struct {
	Type       string `json:"type"`
	URL        string `json:"url"`
	FilePath   string `json:"filePath"`
	Downloaded int64  `json:"downloaded"`
	Errors     int64  `json:"errors"`
	Message    string `json:"message"`
}

// DownloadManager manages download operations
type DownloadManager struct {
	ctx            context.Context
	historyManager *HistoryManager
	activeDownload *downloader.Client
	currentConfig  *DownloadConfig
	startTime      time.Time
	mu             sync.Mutex
	stopChan       chan struct{}
}

// NewDownloadManager creates a new download manager
func NewDownloadManager(ctx context.Context, historyManager *HistoryManager) *DownloadManager {
	return &DownloadManager{
		ctx:            ctx,
		historyManager: historyManager,
	}
}

// StartDownload starts a new download
func (a *App) StartDownload(config DownloadConfig) error {
	return a.downloadManager.Start(config)
}

// StopDownload stops the active download
func (a *App) StopDownload() error {
	a.downloadManager.Stop()
	return nil
}

// ValidateConfig validates the download configuration
func (a *App) ValidateConfig(config DownloadConfig) error {
	client := createClient(config)
	return client.Validate()
}

// Start begins a new download
func (dm *DownloadManager) Start(config DownloadConfig) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// Stop any existing download
	if dm.stopChan != nil {
		close(dm.stopChan)
	}

	dm.stopChan = make(chan struct{})
	dm.currentConfig = &config
	dm.startTime = time.Now()

	// Create client using public API
	client := createClient(config)

	// Validate
	if err := client.Validate(); err != nil {
		return err
	}

	dm.activeDownload = client

	// Set progress callback
	client.OnProgress(func(event downloader.Event) {
		dm.emitProgress(event)
	})

	// Run download in background
	go func() {
		var finalStatus string
		var finalDownloaded, finalErrors int64

		err := client.Download()
		if err != nil {
			finalStatus = "error"
			wailsRuntime.LogError(dm.ctx, "Download failed: "+err.Error())
		} else {
			finalStatus = "completed"
		}

		// Check if stopped
		select {
		case <-dm.stopChan:
			finalStatus = "stopped"
		default:
		}

		// Save to history
		if dm.historyManager != nil && dm.currentConfig != nil {
			entry := HistoryEntry{
				ID:          generateID(),
				URL:         config.URL,
				Domain:      config.Domain,
				Output:      config.Output,
				StartedAt:   dm.startTime.Format(time.RFC3339),
				CompletedAt: time.Now().Format(time.RFC3339),
				Downloaded:  finalDownloaded,
				Errors:      finalErrors,
				Status:      finalStatus,
				Config:      *dm.currentConfig,
			}
			dm.historyManager.Add(entry)
		}

		dm.mu.Lock()
		dm.activeDownload = nil
		dm.mu.Unlock()
	}()

	return nil
}

// Stop stops the active download
func (dm *DownloadManager) Stop() {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.stopChan != nil {
		close(dm.stopChan)
		dm.stopChan = nil
	}
}

// emitProgress sends a progress event to the frontend
func (dm *DownloadManager) emitProgress(event downloader.Event) {
	eventType := "download"
	switch event.Type {
	case downloader.EventStart:
		eventType = "start"
	case downloader.EventDownload:
		eventType = "download"
	case downloader.EventError:
		eventType = "error"
	case downloader.EventComplete:
		eventType = "complete"
	}

	wailsRuntime.EventsEmit(dm.ctx, "download:progress", ProgressEvent{
		Type:       eventType,
		URL:        event.URL,
		FilePath:   event.FilePath,
		Downloaded: event.Downloaded,
		Errors:     event.Errors,
		Message:    event.Message,
	})
}

// createClient creates a downloader client from frontend config
func createClient(config DownloadConfig) *downloader.Client {
	client := downloader.NewClient().
		SetDomain(config.Domain).
		SetURL(config.URL).
		SetOutput(config.Output).
		SetSuffix(""). // No suffix - save directly to output directory
		SetMaxDepth(config.MaxDepth).
		SetConcurrent(config.Concurrent).
		SetRetries(config.Retries).
		SetImages(config.IncludeImages).
		SetProgress(false). // Disable CLI progress bar in GUI mode
		SetVerbose(false)

	if len(config.Blacklist) > 0 {
		client.SetBlacklist(config.Blacklist)
	}
	if len(config.Extensions) > 0 {
		client.SetExtensions(config.Extensions)
	}
	if len(config.ImageExtensions) > 0 {
		client.SetImageExts(config.ImageExtensions)
	}

	return client
}
