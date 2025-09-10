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

## Documentation

For complete usage information, command reference, and advanced features, see the [CLI Documentation](docs/cli/README.md).

## Data Storage

Time tracking data is stored in a SQLite database:

- **macOS**: `~/Library/Application Support/hora/hora.db`
- **Linux**: `~/.local/share/hora/hora.db`

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
