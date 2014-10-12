# API

Current endpoint:

```
 GET /info
 GET /tables
 GET /tables/:table
 GET /tables/:table/indexes
 GET /query
POST /query
 GET /explain
POST /explain
 GET /history
```

# Query Response

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
