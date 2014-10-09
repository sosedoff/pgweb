# pgweb

Experiments with PostgreSQL and GO

## Usage

CLI options:

```
-h, --host= Server hostname or IP (localhost)
-p, --port= Server port (5432)
-u, --user= Database user (postgres)
-d, --db=   Database name (postgres)
```

## Compile

Go 1.3+ is required. To complire source execute:

```
go get
go build
```

This will produce `pgweb` binary in the current workdir.

## API

Get databases:

```
GET /databases
```

Get current database tables:

```
GET /tables
```

Execute select query:

```
POST /query?query=SQL
GET  /query?query=SQL
```