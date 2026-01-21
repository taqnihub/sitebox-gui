# SiteBox GUI

A cross-platform desktop application for downloading websites. Built with [Wails](https://wails.io/) and [Svelte](https://svelte.dev/).

SiteBox GUI provides a user-friendly interface for the [SiteBox](https://github.com/user/sitebox) command-line tool, making it easy to download entire websites for offline viewing.

## Features

- **Easy-to-use interface** - Simply enter a URL and click download
- **Real-time progress** - Watch downloads happen with live statistics
- **Download history** - Keep track of past downloads and re-run them
- **Presets** - Save your favorite configurations for quick access
- **Customizable settings** - Configure concurrent downloads, depth limits, retries, and more
- **Theme support** - Light, dark, and system theme options
- **Cross-platform** - Works on macOS, Windows, and Linux

## Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [Releases](https://github.com/user/sitebox-gui/releases) page:

- **macOS**: `SiteBox.app` (Universal binary for Intel and Apple Silicon)
- **Windows**: `SiteBox.exe`
- **Linux**: `SiteBox` (AppImage)

### macOS

1. Download `SiteBox.app.zip`
2. Extract and drag `SiteBox.app` to your Applications folder
3. On first launch, right-click and select "Open" to bypass Gatekeeper

### Windows

1. Download `SiteBox.exe`
2. Run the executable (you may need to allow it through Windows Defender)

### Linux

1. Download the AppImage
2. Make it executable: `chmod +x SiteBox.AppImage`
3. Run: `./SiteBox.AppImage`

## Usage

### Basic Download

1. **Enter URL** - Paste the website URL you want to download
2. **Domain Filter** - Automatically set from URL, limits downloads to this domain
3. **Output Directory** - Click "Browse" to select where to save files
4. **Start Download** - Click the button and watch the progress

### Advanced Options

Click "Show Advanced Options" to configure:

| Option | Description | Default |
|--------|-------------|---------|
| Concurrent Downloads | Number of parallel downloads | 5 |
| Max Depth | How deep to follow links | 50 |
| Retries | Retry attempts for failed requests | 3 |
| Include Images | Download images (PNG, JPG, SVG, etc.) | Off |
| Blacklist | URL patterns to exclude | github.com |
| Extensions | File types to download | .js, .css, .woff2, .woff, .ttf |

### Presets

Save your current configuration as a preset:

1. Configure your download options
2. Go to the **Presets** tab
3. Click "Save Current Settings"
4. Enter a name and save

Load a preset anytime by clicking "Load" next to it.

### History

The **History** tab shows all completed downloads. You can:

- **Re-run** - Start a new download with the same configuration
- **Open Folder** - Open the output directory in your file manager

### Settings

Configure default behavior in the **Settings** tab:

- **Theme** - Choose Light, Dark, or System
- **Default Output Directory** - Pre-fill the output path
- **Default Options** - Set default values for concurrent downloads, depth, retries
- **Include Images by Default** - Always download images
- **Desktop Notifications** - Get notified when downloads complete

## Building from Source

### Prerequisites

- [Go](https://golang.org/) 1.21 or later
- [Node.js](https://nodejs.org/) 18 or later
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Clone and Build

```bash
# Clone the repository
git clone https://github.com/user/sitebox-gui.git
cd sitebox-gui

# Install dependencies
make deps

# Build for current platform
make build

# Or build for all platforms
make build-all
```

### Development

```bash
# Run in development mode with hot reload
make dev
```

## Project Structure

```
sitebox-gui/
├── main.go                 # Wails entry point
├── wails.json              # Wails configuration
├── internal/gui/           # Go backend
│   ├── app.go              # App lifecycle
│   ├── download.go         # Download management
│   ├── history.go          # History persistence
│   ├── presets.go          # Presets management
│   └── settings.go         # Settings management
└── frontend/               # Svelte frontend
    └── src/
        ├── routes/         # Pages
        └── lib/
            ├── components/ # UI components
            └── stores/     # State management
```

## Data Storage

Settings, presets, and history are stored in:

- **macOS**: `~/Library/Application Support/sitebox/`
- **Windows**: `%APPDATA%\sitebox\`
- **Linux**: `~/.config/sitebox/`

## Requirements

- macOS 10.15 or later
- Windows 10 or later
- Linux with WebKit2GTK

## License

MIT License - see [LICENSE](LICENSE) for details.

## Related Projects

- [SiteBox CLI](https://github.com/user/sitebox) - Command-line version

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.
