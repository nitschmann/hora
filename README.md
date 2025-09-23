# hora - Simple Time Tracking CLI

A simple and intuitive command-line time tracking tool built with Go. Track your project time with ease using a clean CLI interface.

## Features

- **Simple Time Tracking** - Start, stop, and pause time tracking for any project
- **Project Management** - Automatic project creation and management
- **Background Tracking** - Automatic pause/resume on screen lock (macOS)
- **Data Export** - Export time entries to CSV for further analysis
- **Category Support** - Organize time entries with custom categories
- **Rich Reporting** - View detailed time reports with pause information
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
