# Tools

This directory contains utility tools for the Hora project.

## gendocs

Generates markdown documentation for all CLI commands using the `spf13/cobra/doc` package.

### Usage

```bash
# Generate docs using make
make docs

# Or run directly
go run ./tools/gendocs

# Clean generated docs
make clean-docs
```

### Output

The generated documentation will be placed in `docs/cli/` directory with individual markdown files for each command:

- `README.md` - Main command documentation (renamed from `hora.md` for GitHub display)
- `hora_export.md` - Export command documentation
- `hora_project.md` - Project command documentation
- `hora_project_export-times.md` - Project export command documentation
- And more...

Each file contains:
- Command description and synopsis
- Available flags and options
- Inherited parent command options
- Cross-references to related commands
