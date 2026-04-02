# ORDER BY + LIMIT for Advanced and Aggregate Panels

**Date:** 2026-04-02
**Status:** Approved

---

## Summary

Add multi-column ORDER BY and an explicit LIMIT input to both the Advanced search panel and the Aggregate panel in the Rows tab. Advanced search switches from the backend `getTableRows()` API to building its query fully client-side (consistent with how Aggregate already works), so ORDER BY appears in Show Query and no backend changes are needed.

---

## Approach

Option A (approved): client-side SQL generation for both panels. No Go backend changes.

---

## UI

### ORDER BY section (both panels)

- Placed below HAVING in Aggregate; below condition rows in Advanced.
- Section header: `ORDER BY` label + `+ Add` button.
- Each row:
  - Column dropdown
  - Direction toggle: `ASC` / `DESC` (default ASC)
  - Remove (`−`) button
- Rows are comma-separated in SQL (no AND/OR conjunction).
- Column options:
  - **Advanced:** all table columns (same list as condition column dropdown)
  - **Aggregate:** GROUP BY columns + aggregate aliases (PostgreSQL allows alias references in ORDER BY of grouped queries)

### LIMIT input (both panels)

- Inline in the panel footer, alongside Apply / Clear / Show Query.
- Accepts a positive integer; non-numeric or empty falls back to `getRowsLimit()`.
- Resets to empty on Clear.
- Included in the generated SQL (visible in Show Query).
- Does not affect the global paginator rows-per-page setting.

---

## SQL Generation

### Shared helper

```javascript
// Returns "ORDER BY \"col1\" ASC, \"col2\" DESC" or null.
function buildOrderByClause(rowsSelector) { ... }
```

Iterates ORDER BY rows under `rowsSelector`, skips rows with empty column. Returns `null` if no valid rows.

### Aggregate (`buildAggregateQuery`)

Clause order: `SELECT … FROM … WHERE … GROUP BY … HAVING … ORDER BY … LIMIT … OFFSET …`

ORDER BY aliases are legal in PostgreSQL after GROUP BY.

### Advanced — new `buildAdvancedQuery()`

Replaces `getTableRows()` call in `showTableContent()` when `advancedSearchActive` is true.

```sql
SELECT * FROM "schema"."table"
[WHERE <where clause>]
[ORDER BY "col1" ASC, "col2" DESC]
LIMIT n [OFFSET n]
```

- `applyAdvancedSearch()` stores ORDER BY rows alongside the WHERE data so pagination reruns pick them up.
- Column-header click sort is suppressed when `advancedSearchActive` is true and ORDER BY rows are present (mirrors existing aggregate guard).
- Uses `executeQuery()` instead of `getTableRows()`.

### LIMIT resolution (both panels)

```javascript
function getPanelLimit(panelSelector) {
  var v = parseInt($(panelSelector + " .panel-limit-input").val(), 10);
  return (v > 0) ? v : getRowsLimit();
}
```

---

## State & Reset

- ORDER BY rows and LIMIT input are cleared by their respective Clear buttons.
- Advanced: ORDER BY state stored in `#advanced_search_panel` data (alongside `where`), so `showPaginatedTableContent()` reruns correctly.
- Aggregate: ORDER BY rows live in `#agg_order_rows` DOM; `buildAggregateQuery()` reads them fresh on every run.

---

## Files Changed

| File | Change |
|---|---|
| `static/index.html` | Add ORDER BY sections + LIMIT inputs to both panels |
| `static/js/app.js` | `buildOrderByClause()`, `buildAdvancedQuery()`, update `buildAggregateQuery()`, update `applyAdvancedSearch()`, update reset handlers |
| `static/css/app.css` | Style ORDER BY rows (reuse `.adv-search-row` pattern), LIMIT input |

No backend changes required.

---

## Out of Scope

- OFFSET input (global paginator already handles paging)
- ORDER BY for the simple filter (not a panel)
- Saving ORDER BY state across table navigation
