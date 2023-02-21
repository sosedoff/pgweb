# 0.14.0 - 2023-02-21

- FIX: History page query loading fixup, GH-632
- NEW: Display cell content via context menu, GH-634
- NEW: Handle support/permissions errors in info call, GH-635
- NEW: Show error message when API calls fail, GH-636
- NEW: Add bookmark options to load username/password from env vars, GH-638
- NEW: Add context menu to display database tables stats, GH-639
- NEW: Added Local Queries feature, GH-641
- FIX: Ensure that objects are sorted by schema and name, GH-648
- FIX: Fetch local queries on db connect, GH-650

# 0.13.1 - 2022-12-27

- Fix connect flow when `~/.pgweb/bookmarks` directory is not available, GH-631

# 0.13.0 - 2022-12-25

- Add support for .pgpass file, GH-617
- Request logging additions (request id, forwarded user), GH-618
- Establish connections using bookmark ID only, GH-619
- Display empty schemas on the sidebar, GH-621
- Configure timeout and retries when testing connection status, GH-623
- Setup basic prom metrics endpoint, GH-624
- Add default connect_timeout option to connection string, GH-626
- Add duration_ms to log entries, GH-628
- Add query execution stats to api endpoint, GH-629

# 0.12.0 - 2022-12-13

- Deprecate usage of Gox for binary builds, GH-571
- Add netcat install in dockerfile to provide a way to healthcheck, GH-572
- Install latest postgres client in docker image, GH-577
- Add support for `PGWEB_` prefix environment variables, GH-585
- Fix export URL generation, refactor export code, GH-588
- Add logrus-based request logger, GH-589
- Configure logger for connect backend, GH-591
- Set LDFLAGS for make build/release commands, GH-592
- Add internal sessions manager, GH-593
- Include index size on the index list view, GH-595
- Fix flaky backend connection test, GH-596
- Add ability to view and copy views/materialized views definitions, GH-594
- Enable dev assets mode with PGWEB_ASSETS_DEVMODE env var, GH-597
- Make query input box resizable, GH-599
- Deprecate Heroku demo deployments and switch to Fly, GH-600
- Handle returning values in update/delete queries, GH-601
- Fix panic with invalid time marshaling, GH-602
- Configure logging level and format, GH-605
- Use go embed to load queries from static files, GH-607
- Switch go build target to 1.19, GH-603
- Add support for user functions, GH-608
- Implement global query timeout option, GH-609
- Switch windows tests from Appveyor to Github Actions, GH-611
- Fix activity endpoint panic when server version is not detected, GH-612

# 0.11.12 - 2022-07-05

- Update base docker image (alpine), update deps, GH-558
- Refactor docker images building, include ARM, GH-568

# 0.11.11 - 2022-03-29

- Auto-detect the query from the query source based on user selection, GH-547
- Added binary codec base58 as well as improving the help for --binary-codec flag, GH-548
- Change binary codec back to none, GH-555

## 0.11.10 - 2022-01-20

- Removes alert on column copy value, GH-536
- Migrate test suite to Github Action, GH-540
- Serialize binary bytea cols into hex/base64, GH-537
- Include build time into version string, GH-541
- Explain analyze dropdown button, GH-532
- Switch to go 1.17, GH-543
- Use HTTP 302 status code for successful backend redirect, GH-544
- Add connect backend tests, GH-546

## 0.11.9 - 2021-11-08

- Releases are built on Go 1.17
- Build time correction, GH-521
- Fix broken assets URL path prefix, GH-525
- Update docker build image to alpine:3.14, GH-522
- Upgrade gin dependency to v1.7.4, GH-527
- Add FreeBSD startup script, GH-520

## 0.11.8 - 2021-07-07

- Releases are built with Go 1.16
- Add ARM64 v7 build target, GH-497
- Switch to Go modules for dependency management, GH-509
- Switch to Go embed for static assets management, GH-510
- Add Darwin/ARM64 build target (Apple Silicon), GH-513

## 0.11.7 - 2020-10-18

- Releases are built with Go 1.15
- Show results row context menu on custom query results, GH-457
- Do not terminate if local authentication failed on start, GH-463
- Do not show other databases if session is locked, GH-470
- Strip debug information from binary to reduce size, GH-489
- Disable autocomplete on database search field, GH-492
- Improve windows connection error matching during start, GH-493

## 0.11.6 - 2020-02-19

- Add CLI options for SSL key, cert and root certs, GH-452
- Remove double click action on cell, GH-455

## 0.11.5 - 2019-12-16

- Add basic SQL keyword autocompletion, GH-443
- SSH Private Key handling update (encrypted keys are supported now), GH-445
- Include Go version into `pgweb --version` output, GH-447
- Fix long table name bug in the sidebar, GH-448
- Add SQL objects (table,views,etc) autocompletion, GH-449
- Include Go version into info API endpoint, GH-450

## 0.11.4 - 2019-10-05

- Fix SQL export filename, GH-438
- Update Docker image to alpine:3.10, GH-439
- Drop unsupported pg_dump options from connection string, GH-441
- Misc code cleanup and formatting, GH-442

## 0.11.3 - 2019-07-24

- Misc: add script to update homebrew formula version, GH-423
- Destructive keyword restriction in read-only mode, GH-421
- Make database object searchable in sidebar, GH-434
- Update lib/pg to 1.1.1, GH-435

## 0.11.2 - 2019-02-15

- Fix table row estimation query for camelcase schemas, GH-414

## 0.11.1 - 2019-01-28

- Typo fixes
- Add Base64 javascript encoder/decoder to replace deprecated window.atob call, GH-405
- Fix startup error when DATABASE_URL is set, GH-406
- Fix user auto detection when USER env var is not set, GH-408
- Switch bindata dependency to use maintained fork: github.com/go-bindata/go-bindata, GH-409

## 0.11.0 - 2018-12-24

- Tweak sidebar database object counters styles, GH-400
- Do not exit with error if local server is not running, GH-399
- Fix SSH host verification check, GH-398
- Scope activity list to current database only, GH-397
- Show current release version and check for updates, GH-396
- Force switch back to default connection settings view, GH-395
- Fix row count estimation bug, GH-394
- Print out failed query SQL and args with --debug flag, GH-393

## 0.10.0 - 2018-11-28

- Fixes relation not found errors when dealing with table names that have uppercase characters, GH-356
- Dockerfile updates, GH-357
- Check if pg_dump is available before running database export, GH-358
- Improvements to CockroachDB integration, GH-365
- Add EstimatedTableRowsCount to avoid count in large tables, GH-366
- Automatically set table filter option to 'equals' if its not set, GH-370
- Dependencies update and switch to dep, GH-375
- Add column context menu item to get numeric stats, GH-377
- Fix issues with connection string builder, GH-378
- Include rows count to numeric stats view on table column, GH-379
- Make localhost to be a default db host, GH-380
- Clear out connection settings/bookmark on login screen when running in session/connect mode
- Add table row context menu with actions, GH-381
- Allow settings url prefix with URL_PREFIX env var, GH-387
- Fix JSON marshal panic when dealing with NaN values, GH-388
- Fix startup behavior when user did not provide a database name, GH-389

## 0.9.12 - 2018-04-23

- Add link to view database connection string format on login page
- Include constraint name under "constraints" tab, GH-343
- Misc CI and config changes

## 0.9.11 - 2017-12-07

- Fix ssl mode for the connection url in the bookmarks, GH-320
- Add support for CORS, GH-321
- Fix custom query results counter for empty queries, GH-322
- Reorganize the table context menu, GH-323
- Disable database connection string text field autocomplete, GH-327
- Add db prefix to the table export files, GH-329
- Add database view context menu with export actions, GH-330

## 0.9.10 - 2017-11-03

- Make idle connection timeout configurable, [GH-282]
- Fix panics when sshinfo is not set on bookmarks, [GH-296]
- Dot now allow using startup bookmark in multi-session mode, [GH-300]
- Add ability to copy table name from the sidebar, [GH-301]

## 0.9.9 - 2017-09-28

- Automatically format JSON data exports, GH-255
- Update Docker image to alpine:3.6, GH-256
- Print out PostgreSQL server version on start in a single-session mode, GH-264
- Record last query timestamp for the client connection, GH-265
- Add context menu for table headers in browse mode (copy name, see unique values), GH-268
- Add ability to export current database dump, GH-270
- Automatically open pgweb in browser on start if its already running, GH-272
- Connect to the database with credentials provided by a third-party backend, GH-266
- Automatically close idle sessions (no activity in 1 hour), GH-275
- Allow connecting via SSH with a custom private key and other fixes, GH-277
- Add options to disable SSH connections, GH-279

## 0.9.8 - 2017-08-04

- Fixed error checking in the API, GH-234
- Fixed activity tab to support PG 9.x versions, GH-237
- Remember sort column and order for pagination, GH-240
- Use `sslmode=disable` for bookmarks without sslmode option, GH-244
- Javascript fixes for IE9-11, GH-245
- Require confirmation for the disconnect, GH-246
- Clean the results table on manual disconnect

## 0.9.7 - 2017-04-04

- Fixed issue with locked session and empty db url, GH-206
- Fixed path rewrite on DB change, GH-212
- Upgraded dependencies, GH-217
- Added ability to specify bookmarks path, GH-218
- Added counter for the number of rows from a custom SQL query, GH-224
- Added new behavior for removing table rows view on custom SQL query page, GH-225

## 0.9.6 - 2016-11-18

- Fixed bug in query base64-encoding, GH-186
- Fixed rows pagination visibility bug, GH-190
- Fixed issue with query order escaping, GH-191
- Fixed invalid query selection for explain command, GH-198
- Fixed issue with empty sidebar, now it shows empty state, GH-202
- Added new flag --readonly to enable read only transaction mode, GH-193
- Added ability to kill any running query, GH-194
- Added session database connection locking, GH-195
- Added ability to switch between databases, GH-196
- Added feature to keep last selected tab when switching between tables, GH-197
- Added new flag --bookmark (-b) to specify server connection from bookmark, GH-201

## 0.9.5 - 2016-10-01

- Only view schema with USAGE privileges, GH-167
- Fixed broken export to CSV/JSON/XML if hashmark in URL, GH-175
- Added example service configuration for systemd, GH-177
- Allow setting auth user and pass using variables

## 0.9.4 - 2016-07-29

- Fixes CSV/JSON/XML export buttons when pgweb is running with url prefix, GH-170

## 0.9.3 - 2016-06-30

- Uses Go 1.6 for development, GH-155
- Fixes timestamp formatting in CSV export, GH-163
- Included PostgreSQL 9.6 for integration testing
- Switches docker image to Alpine to reduce image size
- Adds support for ARMv5

## 0.9.2 - 2016-03-01

- Fixes bug with unsafe base64 encoded sql queries
- Fixes issue with session id not being included in multi-session mode
- Fixes visual issue with long table names in sidebar
- Fixes visual issue with a scrollbar in table information widget
- Fixes issue with database connection form being reset by clicking on 'cancel' button
- Adds ability to close connection
- Adds display message for number of affected rows for update/delete queries, GH-133
- Adds web server url prefix as a CLI option, GH-135

## 0.9.1 - 2016-01-25

- Fixes bug with tables context menu
- Fixes JS bug when query returns no rows
- Fixes bug with switching between different connection modes
- Adds AJAX timeout to 5s
- Adds sidebar reload action on any CREATE/DROP action

## 0.9.0 - 2016-01-19

- Add support for multiple schemas. GH-112
- Add support for native ssh tunnes. GH-114
- Add materialized views to list of schema objects
- Adds a few design tweaks and cleanups
- Fixes bug with nil result set when fetching rows

## 0.8.0 - 2016-01-11

- Fixes bug with bigint conversions in javascript. Now bigints are encoded as strings. GH-109
- Adds pagination and simple column filtering to table rows browser. GH-110
- Adds ability to use pgweb with multiple database sessions. GH-111
- Adds a few design tweaks and cleanups

## 0.7.0 - 2016-01-05

- Adds sequences to the sidebar panel - GH-100
- Adds table constrains view - GH-104
- Adds ability to export table and query rows as JSON/XML - GH-107
- Updates to UI theme and SQL editor

## 0.6.3 - 2015-08-16

- Adds PostgreSQL password escaping in web ui, GH-96
- Adds base64 query encoding for CSV export, GH-95
- Adds automatic saving of last executed query to localStorage
- Adds request middleware to log incoming form params in debug mode

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
