[![Tests](https://github.com/nitschmann/hora/actions/workflows/test.yml/badge.svg)](https://github.com/nitschmann/hora/actions/workflows/test.yml)

# hora - Simple Time Tracking CLI

A simple and intuitive command-line time tracking tool built with Go. Track your project time with ease using a clean CLI interface.

## Features

- **Simple Time Tracking** - Start, stop, and pause time tracking for any project
- **Project Management** - Automatic project creation and management
- **Background Tracking** - Automatic pause/resume on screen lock (macOS)
- **Data Export** - Export time entries to CSV for further analysis
- **Category Support** - Organize time entries with custom categories
- **Rich Reporting** - View detailed time reports with pause information
- **Web Dashboard** - Interactive web UI with charts, analytics, and filtering
- **Cross-Platform** - Works on macOS and Linux

## Quick Start

### Installation

#### Option 1: Pre-compiled Binaries (Recommended)

Download the latest release binary for your platform from the [Releases page](https://github.com/nitschmann/hora/releases):

- **macOS**: `hora-darwin-amd64` or `hora-darwin-arm64`
- **Linux**: `hora-linux-amd64` or `hora-linux-arm64`

```bash
# Download and make executable (example for macOS ARM64)
curl -L -o hora https://github.com/nitschmann/hora/releases/latest/download/hora-darwin-arm64
chmod +x hora
sudo mv hora /usr/local/bin/

# Verify installation
hora version
```

#### Option 2: Build from Source

```bash
# Build from source
make build

# Or build directly
go build -o build/hora ./cmd/hora
```

### Basic Usage

```bash
# Start tracking time for a project
hora start "My Project"

# Check current status
hora status

# Stop tracking
hora stop

# List all time entries
hora times

# Export to CSV
hora export --output times.csv

# Launch web dashboard
hora ui

# Launch web dashboard on custom port
hora ui --port 3000
```

## Configuration

hora uses a YAML configuration file to customize behavior. The configuration file is automatically created in the default location, but you can initialize it manually or customize the settings.

### Default Configuration

By default, hora uses these settings:

```yaml
database_dir: "/path/to/database/directory"
debug: false
list_limit: 50
list_order: "desc"
use_background_tracker: true
```

### Configuration File Locations

hora looks for configuration files in the following order:

1. `~/.hora/config.yaml` (user's home directory)
2. `./.hora/config.yaml` (current directory)

### Initializing Configuration

Create a configuration file with default values:

```bash
# Initialize config in default location (~/.hora/)
hora config init

# Initialize config in a specific directory
hora config init --directory /path/to/config/dir

# Force overwrite existing config file
hora config init --force
```

### Configuration Options

| Option | Description | Default | Valid Values |
|--------|-------------|---------|--------------|
| `database_dir` | Directory where SQLite database is stored | Platform-specific | Any valid directory path |
| `debug` | Enable debug logging | `false` | `true`, `false` |
| `list_limit` | Maximum number of entries to show in lists | `50` | `1` or greater |
| `list_order` | Sort order for time entry lists | `desc` | `asc`, `desc` |
| `use_background_tracker` | Enable automatic pause/resume on screen lock | `true` | `true`, `false` |

### Using Custom Configuration

You can specify a custom configuration file:

```bash
# Use a specific config file
hora --config /path/to/custom/config.yaml start "My Project"

# Or use the short form
hora -c /path/to/custom/config.yaml start "My Project"
```

## Web Dashboard

hora includes a modern web dashboard that provides interactive analytics and visualization of your time tracking data. The dashboard offers a comprehensive view of your productivity patterns with beautiful charts and filtering capabilities.

### Starting the Web Dashboard

```bash
# Start the web dashboard on default port (8080)
hora ui

# Start on a custom port
hora ui --port 3000

# The dashboard will be available at http://localhost:8080 (or your custom port)
```

### Dashboard Features

#### Interactive Analytics
- **Daily Activity Chart** - Visualize your daily time tracking patterns
- **Project Distribution** - See how your time is distributed across projects
- **Category Breakdown** - Analyze time spent in different categories
- **Time Range Filtering** - View data for last 24hrs, 3 days, 7 days, 30 days, or all time

#### Key Metrics
- **Total Hours** - Cumulative time tracked
- **Project Count** - Number of active projects
- **Category Count** - Number of categories used
- **Time Entries** - Total number of tracking sessions
- **Average Session** - Average duration per tracking session

#### Advanced Filtering

- **Project Filtering** - Click on any project in charts or entries to filter by project
- **Category Filtering** - Click on any category to filter by category
- **URL Parameters** - Share filtered views with direct links
- **Clear Filters** - Reset all filters with one click

#### Recent Entries
- **Live Time Entries** - View all your recent time tracking sessions
- **Interactive Elements** - Click on projects or categories to filter
- **Duration Display** - See exact time spent on each session
- **Date & Time** - Full timestamp information for each entry

## Documentation

For complete usage information, command reference, and advanced features, see the [CLI Documentation](docs/cli/README.md).

## Data Storage

Time tracking data is stored in a SQLite database. The default database location is:

- **macOS**: `~/Library/Application Support/hora/hora.db`
- **Linux**: `~/.local/share/hora/hora.db`

You can customize the database location by setting the `database_dir` option in your configuration file (see [Configuration](#configuration) section above).

## Development

```bash
# Build
make build

# Run tests
make test

# Generate CLI documentation
make docs

# Clean build artifacts
make clean
```

## Requirements

(only for development)

- Go 1.25 or later
- SQLite3
- macOS or Linux

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the MIT License.
