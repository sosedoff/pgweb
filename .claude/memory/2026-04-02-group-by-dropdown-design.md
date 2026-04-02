# Group By Dropdown Design

**Date:** 2026-04-02
**Feature:** Replace GROUP BY checkbox grid with native multi-select dropdown

## Context

The aggregate panel's GROUP BY section currently renders a flat grid of checkboxes (one per column), injected by `buildGroupBySection()`. This design scales poorly for tables with many columns and lacks a clear affordance for multi-selection.

## Goal

Replace the checkbox grid with a native `<select multiple>` element. Selection alone is required — column ordering is not.

## Design

### HTML (`static/index.html`)

Replace the `#agg_group_by_grid` div with a `<select multiple>` and a hint label inside the existing `.agg-section`:

```html
<div class="agg-section">
  <div class="agg-section-header">GROUP BY</div>
  <select id="agg_group_by_select" multiple class="form-control agg-group-by-select"></select>
  <div class="agg-group-by-hint">Hold Ctrl / Cmd to select multiple columns</div>
</div>
```

### JS (`static/js/app.js`)

**`buildGroupBySection()`** — populate `<option>` elements instead of checkbox labels:

```js
function buildGroupBySection() {
  var sel = $("#agg_group_by_select");
  sel.empty();
  $("#pagination select.column option").each(function() {
    var val = $(this).val();
    if (!val) return;
    sel.append($("<option>", { value: val, text: val }));
  });
}
```

**`buildGroupByClause()`** — read `.val()` on the multi-select (returns array or null):

```js
function buildGroupByClause() {
  var cols = $("#agg_group_by_select").val() || [];
  if (cols.length === 0) return null;
  return "GROUP BY " + cols.map(function(c) { return '"' + c + '"'; }).join(", ");
}
```

**`buildAggregateSelectClause()`** — iterate `.val()` instead of `.agg-group-col:checked`:

```js
var cols = $("#agg_group_by_select").val() || [];
cols.forEach(function(c) { parts.push('"' + c + '"'); });
```

**`resetAggregate()`** — deselect all: `$("#agg_group_by_select").val([]);`

### CSS (`static/css/app.css`)

Remove `.agg-group-by-grid` and `.agg-group-by-grid label` rules. Add:

```css
.agg-group-by-select {
  height: 90px;
  font-size: 12px;
}
.agg-group-by-hint {
  font-size: 11px;
  color: #aaa;
  margin-top: 3px;
}
```

## Scope

- 4 JS functions modified: `buildGroupBySection`, `buildGroupByClause`, `buildAggregateSelectClause`, `resetAggregate`
- 1 HTML block replaced: `#agg_group_by_grid` → `#agg_group_by_select`
- 2 CSS rules removed, 2 added
- No changes to `buildAggregateQuery`, HAVING logic, or aggregate expression logic
- No new dependencies
