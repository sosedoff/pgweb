# pgweb

Web-based PostgreSQL database browser written in Go.

[![Release](https://img.shields.io/github/release/sosedoff/pgweb.svg)](https://github.com/sosedoff/pgweb/releases)
[![Build Status](https://img.shields.io/travis/sosedoff/pgweb.svg)](https://travis-ci.org/sosedoff/pgweb)

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
- Save server bookmarks

Visit [wiki Pages](https://github.com/sosedoff/pgweb/wiki) for more details

## Installation

[Precompiled binaries]((https://github.com/sosedoff/pgweb/releases) for supported 
operating systems are available.

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

## Deploy on Heroku

[![Heroku Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy?template=https://github.com/sosedoff/pgweb)


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