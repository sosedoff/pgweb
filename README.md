# pgweb

Web-based PostgreSQL database browser written in Go.

[![Release](https://img.shields.io/github/release/sosedoff/pgweb.svg?label=Release)](https://github.com/sosedoff/pgweb/releases)
[![Linux Build](https://img.shields.io/travis/sosedoff/pgweb.svg?label=Linux)](https://travis-ci.org/sosedoff/pgweb)
[![Windows Build](https://img.shields.io/appveyor/ci/sosedoff/pgweb/master.svg?label=Windows)](https://ci.appveyor.com/project/sosedoff/pgweb)

## Overview

Pgweb is a web-based database browser for PostgreSQL, written in Go and works
on OSX, Linux and Windows machines. Main idea behind using Go for backend development
is to utilize ability of the compiler to produce zero-dependency binaries for 
multiple platforms. Pgweb was created as an attempt to build very simple and portable
application to work with local or remote PostgreSQL databases.

<img src="screenshots/browse.png" width="345px" />
<img src="screenshots/query.png" width="345px" />

## Features

- Works on OSX, Linux and Windows
- Zero dependencies
- Simple installation (distributes as a single binary)
- Connect to local or remote servers
- Browse tables and table data
- Get table details: structure, size, indeces, row count
- Run / analyze custom queries
- Export query results to CSV
- Query history
- Server bookmarks

Visit [WIKI](https://github.com/sosedoff/pgweb/wiki) for more details

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
```

## Deploy on Heroku

[![Heroku Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy?template=https://github.com/sosedoff/pgweb)

## Testing

Run tests:

```
make test
```

## Contribute

- Fork repository
- Create a feature or bugfix branch
- Open a new pull request
- Use [github issues](https://github.com/sosedoff/pgweb/issues) for any questions

## Contact

- Dan Sosedoff
- [dan.sosedoff@gmail.com](mailto:dan.sosedoff@gmail.com)
- [http://twitter.com/sosedoff](http://twitter.com/sosedoff)

## License

The MIT License (MIT)

Copyright (c) 2014-2015 Dan Sosedoff, <dan.sosedoff@gmail.com>
