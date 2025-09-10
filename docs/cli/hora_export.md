## hora export

Export time entries to CSV

### Synopsis

Export all time entries across projects to a CSV file

```
hora export [flags]
```

### Options

```
      --category string   Filter by category
  -h, --help              help for export
  -l, --limit int         Maximum number of entries to show (default 50)
  -o, --output string     Output file path (default: TIMESTAMP_times.csv)
      --since string      Only show entries since this date (YYYY-MM-DD format)
  -s, --sort string       Sort order: 'asc' (oldest first) or 'desc' (newest first) (default "desc")
```

### Options inherited from parent commands

```
  -c, --config string   Path to configuration file
```

### SEE ALSO

* [hora](README.md)	 - hora is a simple time tracking CLI tool

