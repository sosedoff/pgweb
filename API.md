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
