# pgweb

Web-based PostgreSQL database browser written in Go.

## Usage

CLI options:

```
-h, --host= Server hostname or IP (localhost)
-p, --port= Server port (5432)
-u, --user= Database user (postgres)
-d, --db=   Database name (postgres)
    --url=  Database connection string (postgresql://...)
    --ssl=  SSL option (disable)
```

## Compile from source

Go 1.3+ is required. You can install Go with `homebrew`:

```
brew install go
```

To compile source code run the following command:

```
go get
go build
```

This will produce `pgweb` binary in the current directory.

## TODO

- [ ] Change server port to something better
- [ ] Add ability to switch between databases
- [ ] More detailed usage / help section
- [ ] Build for linux and windows
- [ ] Add view to see table details (not only structure)
- [ ] Make results table rows sortable