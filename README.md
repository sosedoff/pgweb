# pgweb

Web-based PostgreSQL database browser written in Go.

[![Build Status](https://travis-ci.org/sosedoff/pgweb.svg?branch=master)](https://travis-ci.org/sosedoff/pgweb)

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

## Run on Heroku

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy?template=https://github.com/sosedoff/pgweb)

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
```

It works great with [Heroku Postgres](https://postgres.heroku.com) if you need 
to troubleshoot production database or simply run a few queries.

### CLI options

```
Usage:
  pgweb [OPTIONS]

Application Options:
  -v, --version    Print version
  -d, --debug      Enable debugging mode (false)
  -s, --skip-open  Skip browser open on start
      --url=       Database connection string
      --host=      Server hostname or IP (localhost)
      --port=      Server port (5432)
      --user=      Database user (postgres)
      --pass=      Password for user
      --db=        Database name (postgres)
      --ssl=       SSL option (disable)
      --bind=      HTTP server host (localhost)
      --listen=    HTTP server listen port (8080)
      --auth-user= HTTP basic auth user
      --auth-pass= HTTP basic auth password
```

## Build from source

Go 1.3+ is required. You can install Go with `homebrew`:

```
brew install go

brew install mercurial (if need)
yum install mercurial (on centos 6.5)
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


## Use in Docker

Build the image:

```
docker build -t pgweb .
```
Or using make to build images

```
make docker
```

Start container:

```
docker run [OPTIONS of docker] pgweb [OPTIONS of pgweb]
```

### Demo

#### Run postgresql container:

```
docker run -d \
           --name="postgresql" \
           -p 5432:5432 \
           -e USER="super" \
           -e DB="testdb" \
           -e PASS="test123456" \
           paintedfox/postgresql
```

#### Create tables in postgresql:

##### 1. Create postgresql-client container:
```
docker run -ti --link postgresql:db ubuntu bash

```
##### 2. In the ubuntu container:

```
$ apt-get update
$ apt-get install -y postgresql-client
$ psql 	-U super	\
		-h	your-ip	\
		-p	5432	\
		-d	testdb

```
Then you can create some tables in the postgresql


#### Run pgweb container:

```
docker run -d \
           -p 8080:8080 pgweb \
           --url postgres://super:test123456@your-ip:5432/testdb \
           --bind 0.0.0.0
```

Then open [http://your-ip:8080](#) in your browser.

## Contributing

- Fork repository
- Create a feature or bugfix branch
- Open a new pull request
- Use github issues for any questions

## Contact

- Dan Sosedoff
- [dan.sosedoff@gmail.com](mailto:dan.sosedoff@gmail.com)
- [http://twitter.com/sosedoff](http://twitter.com/sosedoff)

## License

The MIT License (MIT)

Copyright (c) 2014 Dan Sosedoff, <dan.sosedoff@gmail.com>