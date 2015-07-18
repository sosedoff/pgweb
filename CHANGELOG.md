## 0.6.2 - 2015-07-15

- Adds ability to specify connection strings prefixed by `postgresql://`, [GH-92]
- Updates configuration for Heroku, [GH-89], [GH-90]
- Updates postgresql library dependency to latest, [GH-91]
- Fixes password field to not display plaintext passwords, [GH-87]

## 0.6.1 - 2015-06-18

- This release is repackage-release targeted to fix binary downloads

## 0.6.0 - 2015-05-31

- Adds ability to execute only selected SQL query in run command view, [GH-85]
- Adds ability to delete/truncate table via context meny on sidebar view 
- Adds ability to export table contents to CSV via context menu on sidebar view
- Changes sidebar color scheme to a lighter and better looking one

## 0.5.3 - 2015-05-06

- Changes default server port from 8080 to 8081 to avoil conflict with RethinkDB
- Changes styles for table rows and connection settings window
- Adds highlighting styles for columns with sort order
- Adds git sha into program version output
- Add new endpoint /api/info to get build details

## 0.5.2 - 2015-04-13

- Adds a new endpoint /activity that retuns active queries
- Adds tab to view active queries
- Adds column sorting when browsing table contents
- Fixes SQL query view when switching to table structure view

## 0.5.1 - 2015-02-23

- Upgrades Gin framework dependency to 0.5.0
- Fixes server crash if another pgweb server is running

## 0.5.0 - 2015-01-13

- Adds Go 1.4 support
- Adds connection string printing in debug mode
- Adds initial bookmarks support
- Adds /api prefix for all API calls
- Adds makefile usage task
- Adds windows CI to verify build process
- Adds example sql database to codebase
- Adds timestamped filenames when exporting results to CSV [GH-75]
- Adds connection checking on each request to prevent api panics
- Adds timestamps to query history records
- Adds current database name to the sidebar
- Adds button to refresh tables list to the sidebar
- Updates all application dependencies
- Changes /api/info endpoint to /api/connection
- Fixes issues with connection string/options parsing
- Fixes capitalized column names in table view
- Fixes connection string validation in /api/connect endpoint

## 0.4.1 - 2014-12-01

- Adds pgweb version on start [GH-65]
- Adds user detection from OS environment
- Adds simple memory profiles with --debug option
- Adds the session user and search path in connection info [GH-67]
- Adds table list reloading after CREATE/DROP TABLE queries [GH-69]
- Adds font awesome icons for the sidebar
- Removes query recording for internal queries [GH-67]
- Fixes default sslmode. Its not longer set to "disable"
- Fixes cells cropping on table indexes view
- Fixes connection URL generation using web interface
- Fixes SQL statements for table row count [GH-67]
- Fixes /tables JSON response if database does not have any tables

## 0.4.0 - 2014-11-11

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

## 0.3.1 - 2014-10-28

- Adds proper exit code when printing version via -v/--version flag
- Adds --version and --debug long flag names
- Adds double quotes for table name when fetching table contents
- Adds support for DATABASE_URL environment variable if no --url is set
- Adds proper usage of jQuery .prop method
- Adds --pass flag to specify connection password
- Fixes --ssl flag usage, previous value was hardcoded

## 0.3.0 - 2014-10-26

- Renamed `make deps` to `make setup` and fix issues with bootstrapping
- Removed hardcoded url for CSV export, it now detects application host:port
- Improved query history view table styles
- Moved table information view to the sidebar
- Added --listen flag to specify web server port, still defaults to 8080

## 0.2.0 - 2014-10-23

- Design tweaks
- Automatically opens browser on OSX systems
- Adds query explain functionality
- Adds export to CSV

## 0.1.0 - 2014-10-14

- Initial release