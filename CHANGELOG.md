## 0.4.1 - Unreleased

- Adds pgweb version on start [GH-65]
- Adds user detection from OS environment
- Adds simple memory profiles with --debug option
- Adds the session user and search path in connection info [GH-67]
- Removes query recording for internal queries [GH-67]
- Fixes default sslmode. Its not longer set to "disable"
- Fixes cells cropping on table indexes view
- Fixes connection URL generation using web interface
- Fixes SQL statements for table row count [GH-67]
- Fixes /tables JSON response if database does not have any tables

## 0.4.0

- Adds query escaping when exporting results to CSV [GH-38]
- Adds keyboard shortcut (ctrl+e, command+e on mac) for query explain action
- Adds HTTP basic authentication with --auth-user and --auth-pass flags
- Adds -skip-open/-s flag to disable automatic browser launch
- Adds --bind option to specify server listen hostname/ip
- Adds ssl mode parameters to url if ssl flag is set and not defined in the url
- Adds dependency management with Godep
- Adds Docker support
- Adds Heroku support
- Adds ability to connect to databases with no tables
- Adds precompiled assets into repository to simplify development
- Adds a connection details view
- Adds a new interface to specify connection settings or make a new connection
- Adds page favicon
- Adds ability to present cell data as text area by double clicking on it
- Fixes styles for query explain results
- Fixes sidebar navigation scrolling styles [GH-12]
- Fixes sidebar table name styles to support long names

## 0.3.1

- Adds proper exit code when printing version via -v/--version flag
- Adds --version and --debug long flag names
- Adds double quotes for table name when fetching table contents
- Adds support for DATABASE_URL environment variable if no --url is set
- Adds proper usage of jQuery .prop method
- Adds --pass flag to specify connection password
- Fixes --ssl flag usage, previous value was hardcoded

## 0.3.0

- Renamed `make deps` to `make setup` and fix issues with bootstrapping
- Removed hardcoded url for CSV export, it now detects application host:port
- Improved query history view table styles
- Moved table information view to the sidebar
- Added --listen flag to specify web server port, still defaults to 8080

## 0.2.0

- Design tweaks
- Automatically opens browser on OSX systems
- Adds query explain functionality
- Adds export to CSV

## 0.1.0

- Initial release