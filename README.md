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

Get current database tables:

```
GET /tables
```

Get table details:

```
GET /tables/:name
```

Execute select query:

```
POST /select?query=SQL
 GET /select?query=SQL
```

### Response formats

Successful response:

```json
{
  "columns": [
    "column_name1",
    "column_name2",
    "column_name3"
  ],
  "rows": [
    [
      "column 1 value",
      "column 2 value",
      "column 3 value" 
    ]
  ]
}
```

Error response:

```json
{
  "error": "Error message"
}
```
