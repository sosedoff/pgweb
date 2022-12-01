# connect-backend-go

Example Golang backend for Pgweb Connect feature

## Usage

Run the backend:

```bash
go run main.go
```

Configure pgweb:

```bash
pgweb --sessions --connect-backend=http://localhost:4567 --connect-token=test
```
