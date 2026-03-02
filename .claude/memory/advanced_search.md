# Advanced Search Feature

## Status: Implemented and working (v2)

## What was built
Multi-condition advanced search panel for the Rows tab. Accessible via an "Advanced" toggle button next to the existing Apply/× buttons.

### Features (v1)
- Multiple filter conditions (unlimited rows)
- "Advanced Filter Active" badge in pagination row when active
- Right-click "Filter Rows By Value" appends a pre-filled row when panel is open
- Resets on table switch and basic × reset button
- `#output` top offset recalculated when panel opens/closes

### Features added in v2
- **Per-row AND/OR connector pills** — each row after the first has AND/OR buttons; first row shows "WHERE" label; enables `(A) AND (B) OR (C)` style expressions
- **Show Query button** — toggles a `<pre>` box showing the full `SELECT * FROM "schema"."table" WHERE ...;` with a Copy button (reuses `copyToClipboard()`)
- **Expanded operator set** (18 operators, grouped with `<optgroup>`):
  - Comparison: `=`, `<>`, `<`, `>`, `<=`, `>=`
  - List: `IN` (comma-sep → `IN ('a','b')`), `NOT IN`
  - Null: `IS NULL`, `IS NOT NULL`
  - Range: `BETWEEN` (From/To inputs), `NOT BETWEEN`
  - Pattern: Contains, Not contains, Has prefix, Has suffix + case-insensitive variants (ILIKE)
- **`getOpInputType(op)`** helper: returns `"none" | "single" | "list" | "range"` — controls which input variant is shown
- **`buildFullQuery()`** — builds full SELECT string using `getCurrentObject().name`

## Files changed
- `static/index.html` — Advanced button, `#advanced_search_panel`, `#adv_search_active_badge`, Show Query button, `#adv_query_display`
- `static/css/app.css` — appended styles for panel, connector pills, range inputs, query display box
- `static/js/app.js` — see key functions below

## Key JS functions (app.js)
- `var advancedSearchActive = false` — global flag
- `escapeSqlLiteral(val)` — doubles single-quotes (also applied to simple filter)
- `getOpInputType(op)` — returns input variant type for an operator
- `buildAdvancedSearchRow(isFirst)` — builds a condition row; `isFirst=true` → WHERE label, no pill
- `buildAdvancedWhereClause()` — iterates rows, reads `data-row-conj` per row, builds SQL
- `buildFullQuery()` — wraps WHERE clause in full SELECT statement
- `applyAdvancedSearch()` — stores WHERE in panel `.data("where")`, sets flag, reloads
- `resetAdvancedSearch()` — clears flag/badge, empties rows, adds `buildAdvancedSearchRow(true)`
- `adjustOutputTop()` — sets `#output` CSS top to `#pagination` outerHeight
- `bindAdvancedOpHandlers()` — delegated handler showing correct input variant per operator

## Key JS edits (app.js)
- `showTableContent()` — advanced takes precedence over simple filter when `advancedSearchActive`
- `buildTableFilters()` — syncs columns into existing advanced rows; passes `isFirst=true`
- Objects click handler — calls `resetAdvancedSearch()` on table switch
- `reset-filters` button — calls `resetAdvancedSearch()`
- `filter_by_value` context menu — appends pre-filled row when panel visible

## No backend changes needed
The existing `where` param on `GET /api/tables/:table/rows` accepts raw SQL.
`TableRows()` in `pkg/client/client.go` appends `WHERE <opts.Where>` directly.

## Build command (GOROOT is broken in this environment)
```bash
GOROOT=/opt/homebrew/Cellar/go/1.25.7_1/libexec \
GOPROXY=https://proxy.golang.org,direct \
GONOSUMDB='*' \
GOOS=linux GOARCH=amd64 \
go build -o ./bin/pgweb_linux_amd64
```
Output: `bin/pgweb_linux_amd64` (~28MB, ELF 64-bit, statically linked)
