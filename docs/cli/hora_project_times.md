## hora project times

List time entries for a specific project

### Synopsis

List all time entries for a specific project, showing start time, end time, and duration. You can specify either the project ID (numeric) or name.

```
hora project times [PROJECT_ID_OR_NAME] [flags]
```

### Options

```
  -h, --help           help for times
  -l, --limit int      Maximum number of entries to show (default 50)
      --since string   Only show entries since this date (YYYY-MM-DD format)
      --sort string    Sort order: 'asc' (oldest first) or 'desc' (newest first) (default "desc")
```

### Options inherited from parent commands

```
  -c, --config string   Path to configuration file
```

### SEE ALSO

* [hora project](hora_project.md)	 - Manage projects

