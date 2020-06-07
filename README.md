# pgweb

Web-based PostgreSQL database browser written in Go.

[![Release](https://img.shields.io/github/release/sosedoff/pgweb.svg?label=Release)](https://github.com/sosedoff/pgweb/releases)
[![Linux Build](https://img.shields.io/travis/sosedoff/pgweb/master.svg?label=Linux)](https://travis-ci.org/sosedoff/pgweb)
[![Windows Build](https://img.shields.io/appveyor/ci/sosedoff/pgweb/master.svg?label=Windows)](https://ci.appveyor.com/project/sosedoff/pgweb)
[![Go Report Card](https://goreportcard.com/badge/github.com/sosedoff/pgweb)](https://goreportcard.com/report/github.com/sosedoff/pgweb)
[![GoDoc](https://godoc.org/github.com/sosedoff/pgweb?status.svg)](https://godoc.org/github.com/sosedoff/pgweb)
[![Docker Pulls](https://img.shields.io/docker/pulls/sosedoff/pgweb.svg)](https://hub.docker.com/r/sosedoff/pgweb/)

## Overview

Pgweb is a web-based database browser for PostgreSQL, written in Go and works
on OSX, Linux and Windows machines. Main idea behind using Go for backend development
is to utilize ability of the compiler to produce zero-dependency binaries for 
multiple platforms. Pgweb was created as an attempt to build very simple and portable
application to work with local or remote PostgreSQL databases.

[See application screenshots](SCREENS.md)

## Features

- Cross-platform support OSX/Linux/Windows 32/64-bit
- Simple installation (distributed as a single binary)
- Zero dependencies
- Works with PostgreSQL 9.1+
- SSH Connections
- Multiple database sessions
- Simple database browser
- Execute and analyze custom SQL queries
- Table and query data export to CSV/JSON/XML
- Query history
- Server bookmarks

Visit [WIKI](https://github.com/sosedoff/pgweb/wiki) for more details

## Pgweb Pro

Pgweb Pro is the next major version of Pgweb and includes features like:

- Table structure editing
- Data editing (update row content)
- Charting
- History persistence
- Multiple tabs

Please get in touch via: https://goo.gl/forms/euQOGWg5uPdPH70b2

## Demo

Visit https://pgweb-demo.herokuapp.com to see pgweb in action.

## Installation

[Precompiled binaries](https://github.com/sosedoff/pgweb/releases) for supported 
operating systems are available.

[More installation options](https://github.com/sosedoff/pgweb/wiki/Installation)

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
SESSIONS=1 pgweb
```

## Deploy on Heroku

[![Heroku Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/sosedoff/pgweb)

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
