# Hora - Simple Time Tracking CLI

Hora is a simple command-line time tracking tool built with Go and Cobra. It allows you to track time for different projects using a SQLite database.

## Features

- Start time tracking for any project
- Stop time tracking and see duration
- Show currently active session status
- Force start new session (stops existing session)
- Clear all time tracking data
- Project management with automatic creation/reuse
- List all projects with creation timestamps and last tracked time in a formatted table
- Remove projects and all their time entries
- SQLite database for persistent storage
- Service layer architecture for clean separation of concerns
- Model-based data structure
- Centralized database connection management
- Interface-based service layer for testability
- Query builder integration with goqu
- Simple and intuitive CLI interface

## Installation

Build the tool from source:

```bash
# Build using Makefile (recommended)
make build

# Or build directly with go
go build -o build/hora ./cmd/hora
```

The binary will be created in the `./build/` directory.

## Usage

### Start tracking time

```bash
# Start tracking with project name as argument
hora start "My Project"

# Start tracking with project name as flag
hora start --project "My Project"

# Start tracking and be prompted for project name
hora start

# Force start new session (stops any existing session)
hora start "New Project" --force
```

### Stop tracking time

```bash
hora stop
```

### Check current status

```bash
hora status
```

### Clear all data

```bash
# Clear with confirmation prompt
hora clear

# Clear without confirmation
hora clear --force
```

### Manage projects

```bash
# List all projects
hora project list

# Remove a project and all its time entries
hora project remove "Project Name"

# Remove without confirmation
hora project remove "Project Name" --force
```

### Get help

```bash
hora --help
hora start --help
hora stop --help
hora status --help
hora clear --help
hora project --help
hora project list --help
hora project remove --help
```

## Database

The tool stores all time tracking data in a SQLite database. The database location varies by operating system:

- **macOS**: `~/Library/Application Support/Hora/hora.db`
- **Windows**: `%LOCALAPPDATA%\Hora\hora.db`
- **Linux/Unix**: `~/.local/share/hora/hora.db`

The database and directory are automatically created on first use.

## Examples

```bash
# Build the tool first
$ make build

# Start tracking time for a project
$ ./build/hora start "Web Development"
Started tracking time for project: Web Development

# Stop tracking after some time
$ ./build/hora stop
Stopped tracking time for project: Web Development
Duration: 01:23:45

# Check current status
$ ./build/hora status
Active session:
  Project: Web Development
  Started: 2025-09-01 20:48:56
  Duration: 00:15:30

# Try to start tracking when already tracking
$ ./build/hora start "Another Project"
Error: already tracking time for project: Web Development

# Force start new session
$ ./build/hora start "Another Project" --force
Started tracking time for project: Another Project (stopped previous session)

# Try to stop when not tracking
$ ./build/hora stop
Error: no active time tracking session found

# Clear all data
$ ./build/hora clear
This will delete ALL time tracking data. Are you sure? (y/N): y
All time tracking data has been cleared.

# List projects
$ ./build/hora project list
┌─────────────────┬──────────────────┬──────────────────┐
│      NAME       │     CREATED      │   LAST TRACKED   │
├─────────────────┼──────────────────┼──────────────────┤
│ Mobile App      │ 2025-09-01 20:49 │ 2025-09-01 21:15 │
│ Web Development │ 2025-09-01 20:48 │ 2025-09-01 21:10 │
└─────────────────┴──────────────────┴──────────────────┘

# Remove a project
$ ./build/hora project remove "Web Development"
This will delete project 'Web Development' and ALL its time entries. Are you sure? (y/N): y
Project 'Web Development' and all its time entries have been removed.
```

## Project Structure

```
hora/
├── build/                    # Build output directory
│   └── .keep                # Keep directory in git
├── cmd/hora/main.go          # Main entry point
├── internal/
│   ├── cmd/
│   │   ├── main.go          # Root command
│   │   ├── start.go         # Start command
│   │   ├── stop.go          # Stop command
│   │   ├── status.go        # Status command
│   │   ├── clear.go         # Clear command
│   │   └── project.go       # Project management commands
│   ├── database/
│   │   ├── interface.go     # Database interface
│   │   ├── database.go      # SQLite implementation
│   │   ├── data_dir.go      # Cross-platform data directory handling
│   │   └── factory.go       # Database factory
│   ├── model/
│   │   ├── project.go       # Project model
│   │   └── time_entry.go    # TimeEntry model
│   └── service/
│       └── time_tracking.go # Business logic service layer
├── .gitignore               # Git ignore rules
├── go.mod
├── Makefile
└── README.md
```

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [sqlx](https://github.com/jmoiron/sqlx) - SQL extensions for Go
- [go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite driver
- [goqu](https://github.com/doug-martin/goqu) - SQL query builder
- [tablewriter](https://github.com/olekukonko/tablewriter) - ASCII table formatting

## Architecture

The application follows a clean architecture pattern with clear separation of concerns:

- **Commands** (`internal/cmd/`) - CLI command handlers with centralized database connection
- **Service Layer** (`internal/service/`) - Business logic interface with concrete implementation
- **Database Layer** (`internal/database/`) - Data persistence with interface-based design
- **Models** (`internal/model/`) - Data structures and domain objects

### Architecture Benefits

- **Centralized Connection Management** - Database connections are managed at the root command level using `PersistentPreRunE`
- **Interface-Based Services** - Service layer uses interfaces for better testability and dependency injection
- **Clean Separation** - Each layer has a single responsibility and clear boundaries
- **Resource Management** - Automatic database connection cleanup via `PersistentPostRun`

### Database Design

- **Projects Table** - Stores project information with unique names and creation timestamps
- **Time Entries Table** - References projects via foreign key with datetime fields
- **Automatic Project Management** - Projects are created on first use and reused thereafter
- **Cascading Deletes** - Removing a project removes all associated time entries

The database layer is designed to be extensible. Currently, only SQLite is supported, but the interface-based design allows for easy addition of other database backends (PostgreSQL, MySQL, etc.) in the future.
