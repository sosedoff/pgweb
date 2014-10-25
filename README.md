# pgweb

Web-based PostgreSQL database browser written in Go.

## Overview

This is a web-based browser for PostgreSQL database server. Its written in Go
and works on Mac OSX, Linux and Windows machines. Main idea behind using Go for the backend
is to utilize language's ability for cross-compile source code for multiple platforms. 
This project is an attempt to create a very simple and portable application to work with 
PostgreSQL databases.

<img src="screenshots/browse.png" width="345px" />
<img src="screenshots/query.png" width="345px" />

Features:

- Connect to local or remote server
- Browse tables and table rows
- Get table details: structure, size, indices, row count
- Execute SQL query and run analyze on it
- Export query results to CSV
- View query history

## Installation

Please visit [Github Releases](https://github.com/sosedoff/pgweb/releases) to download a 
precompiled binary for your operating system.

Currently supported:

- Mac OSX 64bit
- Linux 32/64bit
- Windows 32/64bit

Supported PostgreSQL versions:

- 9.1
- 9.2
- 9.3

Older versions of PostgreSQL might also work but this project is not tested on 
8.x branches.

## Usage

Start server:

```
pgweb --host localhost --user myuser --db mydb
```

You can also specify a connection URI instead of individual flags:

```
pgweb --url postgres://user:password@host:port/database
```

It works great with [Heroku Postgres](https://postgres.heroku.com) if you need 
to troubleshoot production database or simply run a few queries.

Full CLI options:

```
Usage:
  pgweb [OPTIONS]

Application Options:
  -v          Print version
  -d          Enable debugging mode (false)
      --url=  Database connection string
      --host= Server hostname or IP (localhost)
      --port= Server port (5432)
      --user= Database user (postgres)
      --db=   Database name (postgres)
      --ssl=  SSL option (disable)
```

## Compile from source

Go 1.3+ is required. You can install Go with `homebrew`:

```
brew install go
```

To compile source code run the following command:

```
make setup
make dev
```

This will produce `pgweb` binary in the current directory.

There's also a task to compile binaries for other operating systems:

```
make build
```

Under the hood it uses [gox](https://github.com/mitchellh/gox). Compiled binaries
will be stored into `./bin` directory.

## Contributors

- Dan Sosedoff - https://twitter.com/sosedoff
- Masha Safina - https://twitter.com/mashasafina

## License

The MIT License (MIT)

Copyright (c) 2014 Dan Sosedoff, <dan.sosedoff@gmail.com>
