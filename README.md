# pgweb

Simple web-based and cross platform PostgreSQL database explorer.

[![Release](https://img.shields.io/github/release/sosedoff/pgweb.svg?label=Release)](https://github.com/sosedoff/pgweb/releases)
[![Linux Build](https://github.com/sosedoff/pgweb/actions/workflows/checks.yml/badge.svg)](https://github.com/sosedoff/pgweb/actions?query=branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/sosedoff/pgweb)](https://goreportcard.com/report/github.com/sosedoff/pgweb)
[![GoDoc](https://pkg.go.dev/badge/github.com/sosedoff/pgweb)](https://pkg.go.dev/github.com/sosedoff/pgweb)
[![Docker Pulls](https://img.shields.io/docker/pulls/sosedoff/pgweb.svg)](https://hub.docker.com/r/sosedoff/pgweb/)

## Overview

Pgweb is a web-based database explorer for PostgreSQL, written in Go, and works
on Mac, Linux and Windows machines. Distributed as a simple binary with zero dependencies.
Very easy to use and packs just the right amount of features.

[See application screenshots](SCREENS.md)

## Features

- Cross-platform: Mac/Linux/Windows (64bit).
- Simple installation (distributed as a single binary).
- Zero dependencies.
- Works with PostgreSQL 9.6+.
- Supports native SSH tunnels.
- Multiple database sessions.
- Execute and analyze custom SQL queries.
- Table and query data export to CSV/JSON/XML.
- Query history.
- Server bookmarks.
- Inline cell editing, row insert/edit/delete, and multi-column filtering (see [Row editing](#row-editing)).

Visit [WIKI](https://github.com/sosedoff/pgweb/wiki) for more details.

## Demo

Visit https://pgweb-demo.fly.dev/ to see Pgweb in action.

## Installation

- [Precompiled binaries](https://github.com/sosedoff/pgweb/releases) for supported operating systems are available.
- [More installation options](https://github.com/sosedoff/pgweb/wiki/Installation)

## Usage

Start server:

```
pgweb
```

You can also provide connection flags:

```
pgweb --host localhost --user myuser --db mydb
```

Connection URL scheme is also supported:

```
pgweb --url postgres://user:password@host:port/database?sslmode=[mode]
pgweb --url "postgres:///database?host=/absolute/path/to/unix/socket/dir"
```

### Multiple database sessions

To enable multiple database sessions in pgweb, start the server with:

```
pgweb --sessions
```

Or set environment variable:

```
PGWEB_SESSIONS=1 pgweb
```

## Row editing

Pgweb can modify table data directly from the browser, without writing SQL.
The feature is enabled by default and only applies to the table rows view.

**Inline cell editing**

Double-click a cell in the rows view to edit it. Press `Enter` to save, `Escape`
to cancel. Only the edited column is written, and only the affected row is
refreshed in the grid. An empty input clears a nullable column to `NULL`; use
the row editor to store an empty string instead.

**Insert, edit and delete rows**

The `New`, `Edit` and `Delete` buttons in the rows toolbar open a form generated
from the table metadata. `Edit` and `Delete` act on the currently selected row,
and `Delete` always asks for confirmation and shows a summary of the row.
Leaving a field empty on a new row lets the column default apply.

**Requirements and limitations**

- A primary key is required to edit or delete a row. Composite keys are supported.
- Insert is available on any insertable relation, including tables without a primary key.
- Generated and non-updatable columns are read-only, and `bytea` columns are not editable.
- Every statement is parameterized, and each operation is limited to a single row.

**Disabling the feature**

Row editing is turned off automatically in `--readonly` mode. It can also be
disabled explicitly:

```
pgweb --no-row-editing
```

Or set environment variable:

```
PGWEB_NO_ROW_EDITING=1 pgweb
```

When disabled, the write endpoints return `403` and the UI hides all row actions.

See [docs/row-editing-and-filtering.md](docs/row-editing-and-filtering.md) for
the API reference, migration notes and performance assessment.

## Multi-column filtering

The rows view supports combining multiple conditions with `AND`. Pick a column,
an operator and a value, then press `Add Filter` to stack another condition, or
`Apply` to run the query. Active filters are listed above the table and can be
removed individually or cleared at once.

Supported operators: `=`, `!=`, `>`, `>=`, `<`, `<=`, `LIKE`, `ILIKE`, `IN`,
`IS NULL`, `IS NOT NULL`. For `IN`, provide a comma separated list of values.

Filtering is executed by PostgreSQL using parameterized queries, so pagination
and large tables behave exactly as they do without filters.

## Testing

Before running tests, make sure you have PostgreSQL server running on `localhost:5432`
interface. Also, you must have `postgres` user that could create new databases
in your local environment. Pgweb server should not be running at the same time.

Execute test suite:

```
make test
```

If you're using Docker locally, you might also run pgweb test suite against
all supported PostgreSQL version with a single command:

```
make test-all
```

## Contribute

- Fork this repository
- Create a new feature branch for a new functionality or bugfix
- Commit your changes
- Execute test suite
- Push your code and open a new pull request
- Use [issues](https://github.com/sosedoff/pgweb/issues) for any questions
- Check [wiki](https://github.com/sosedoff/pgweb/wiki) for extra documentation

## License

The MIT License (MIT). See [LICENSE](LICENSE) file for more details.
