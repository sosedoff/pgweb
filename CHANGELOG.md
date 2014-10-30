## 0.4.0 (unreleased)

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