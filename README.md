# pgweb

Web-based PostgreSQL database browser written in Go.

## Overview

This is a web-based browser for PostgreSQL database server. Its written in Go
and works on Mac OSX, Linux and Windows machines. Main idea behind using Go for the backend
is to utilize language's ability for cross-compile source code for multiple platforms. 
This project is an attempt to create a very simple and portable application to work with 
PostgreSQL databases.

## Installation

Please visit [Github Releases](https://github.com/sosedoff/pgweb/releases) to download a 
precompiled binary for your operating system.

Currently supported:

- OSX 64bit
- Linux 32/64bit
- Windows 32/64bit

## Usage

To start a server, type: `pgweb`. Application will try to connect to the local PostgreSQL
server with `postgresql` user and select `postgresql` database. You can specify connection
flags, like database, host or user. See `CLI` section of this readme.

You can also specify a connection URI instead of settings individual connection settings:

```
pgweb --url postgresql://user:password@host:port/database
```

It works great with [Heroku Postgres](https://postgres.heroku.com) if you need 
to troubleshoot production database or simply run a few queries.

## CLI

CLI options:

```
Usage:
  pgweb [OPTIONS]

Application Options:
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
make deps
make dev
```

This will produce `pgweb` binary in the current directory.

## Contributors

- Dan Sosedoff - https://twitter.com/sosedoff
- Masha Safina - https://twitter.com/mashasafina
- Jeff Canty - https://twitter.com/cantyjeffrey

## License

The MIT License (MIT)

Copyright (c) 2014 Dan Sosedoff, <dan.sosedoff@gmail.com>