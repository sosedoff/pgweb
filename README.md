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
make deps
make dev
```

This will produce `pgweb` binary in the current directory.