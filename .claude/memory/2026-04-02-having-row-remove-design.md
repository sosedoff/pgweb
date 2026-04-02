---
name: HAVING Row Remove Button Design
description: Design for adding a remove button to all HAVING rows in the aggregate panel, including the first row
type: project
---

# HAVING Row Remove Button Design

**Date:** 2026-04-02
**Feature:** Allow removal of any HAVING row, including the first one

## Problem

`buildHavingRow(isFirst)` only appends the remove button when `!isFirst` (line 1278). The first HAVING row added via `+Add` never receives a remove button, making it impossible to remove without using the "Clear" button (which resets the entire aggregate panel).

## Design

### Change 1: `buildHavingRow()` — `static/js/app.js:1278-1280`

Remove the `if (!isFirst)` guard so every row gets a remove button:

```js
// Remove this guard:
// if (!isFirst) {
row.append('<button type="button" class="btn btn-default btn-xs adv-remove-row"><i class="fa fa-minus"></i></button>');
// }
```

### Change 2: `.adv-remove-row` handler — `static/js/app.js:2307-2309`

After removing the row, if it was the first child, promote the new first row to show the "WHERE" label instead of AND/OR conjunction buttons:

```js
$("#agg_having_rows").on("click", ".adv-remove-row", function() {
  var row = $(this).closest(".adv-search-row");
  var wasFirst = row.is(":first-child");
  row.remove();
  if (wasFirst) {
    var newFirst = $("#agg_having_rows .adv-search-row").first();
    if (newFirst.length) {
      newFirst.find(".adv-row-conj").replaceWith(
        '<div class="adv-row-conj adv-row-conj-first"><span>WHERE</span></div>'
      );
    }
  }
});
```

## Scope

- 2 changes, 1 file: `static/js/app.js`
- No HTML, CSS, or other changes needed
- Consistent with `buildAggregateRow()` which always includes a remove button
