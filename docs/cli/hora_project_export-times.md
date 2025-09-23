## hora project export-times

Export project time entries to CSV

### Synopsis

Export time entries for a specific project to a CSV file

```
hora project export-times [project] [flags]
```

### Options

```
  -h, --help            help for export-times
  -l, --limit int       Maximum number of entries to show (default 50)
  -o, --output string   Output file path (default: TIMESTAMP_PROJECT_times.csv)
      --since string    Only show entries since this date (YYYY-MM-DD format)
      --sort string     Sort order: 'asc' (oldest first) or 'desc' (newest first) (default "desc")
```

### Options inherited from parent commands

```
  -c, --config string   Path to configuration file
```

### SEE ALSO

* [hora project](hora_project.md)	 - Manage projects

