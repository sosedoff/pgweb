# ORDER BY + LIMIT Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add multi-column ORDER BY sections and explicit LIMIT inputs to both the Advanced search panel and Aggregate panel, switching Advanced to client-side SQL generation for consistency.

**Architecture:** Both panels build their full SQL queries client-side. A shared `buildOrderByClause(containerSelector)` helper reads ORDER BY rows from the DOM. `buildAdvancedQuery()` mirrors `buildAggregateQuery()` and replaces the `getTableRows()` path in `showTableContent()` when `advancedSearchActive` is true.

**Tech Stack:** jQuery, Bootstrap 3, vanilla JS, no backend changes.

---

## File Map

| File | Changes |
|---|---|
| `static/index.html` | Add `#adv_order_rows` section + `#adv-limit-input` to Advanced panel; add `#agg_order_rows` section + `#agg-limit-input` to Aggregate panel |
| `static/js/app.js` | Add `buildOrderByRow()`, `buildOrderByClause()`, `updateAdvOrderByColumns()`, `updateAggOrderByColumns()`, `buildAdvancedQuery()`; update `buildAggregateQuery()`, `showTableContent()`, `buildFullQuery()`, `resetAdvancedSearch()`, `resetAggregate()`, `applyAdvancedSearch()`; add event handlers |
| `static/css/app.css` | Add `.order-by-row`, `.order-by-col`, `.order-by-dir`, `.panel-limit-wrap` styles |

---

## Task 1: Add ORDER BY HTML sections to both panels

**Files:**
- Modify: `static/index.html`

- [ ] **Step 1: Add ORDER BY section + LIMIT to Advanced panel**

In `static/index.html`, replace the Advanced panel footer:

```html
<!-- BEFORE -->
        <div class="adv-search-footer">
            <button type="button" class="btn btn-default btn-xs" id="adv-add-condition"><i class="fa fa-plus"></i> Add Condition</button>
            <button type="button" class="btn btn-primary btn-sm" id="adv-apply">Apply</button>
            <button type="button" class="btn btn-default btn-sm" id="adv-reset"><i class="fa fa-times"></i> Clear</button>
            <button type="button" class="btn btn-default btn-sm" id="adv-show-query"><i class="fa fa-code"></i> Show Query</button>
          </div>
```

```html
<!-- AFTER -->
        <div class="agg-section" id="adv-order-section">
            <div class="agg-section-header">
              ORDER BY
              <button type="button" class="btn btn-default btn-xs" id="adv-add-order"><i class="fa fa-plus"></i> Add</button>
            </div>
            <div id="adv_order_rows"></div>
          </div>
          <div class="adv-search-footer">
            <button type="button" class="btn btn-default btn-xs" id="adv-add-condition"><i class="fa fa-plus"></i> Add Condition</button>
            <button type="button" class="btn btn-primary btn-sm" id="adv-apply">Apply</button>
            <button type="button" class="btn btn-default btn-sm" id="adv-reset"><i class="fa fa-times"></i> Clear</button>
            <button type="button" class="btn btn-default btn-sm" id="adv-show-query"><i class="fa fa-code"></i> Show Query</button>
            <span class="panel-limit-wrap">LIMIT <input type="text" id="adv-limit-input" class="form-control panel-limit-input" placeholder="default" /></span>
          </div>
```

- [ ] **Step 2: Add ORDER BY section + LIMIT to Aggregate panel**

In `static/index.html`, replace the Aggregate panel footer:

```html
<!-- BEFORE -->
          <div class="adv-search-footer">
            <button type="button" class="btn btn-primary btn-sm" id="agg-apply">Apply</button>
            <button type="button" class="btn btn-default btn-sm" id="agg-reset"><i class="fa fa-times"></i> Clear</button>
            <button type="button" class="btn btn-default btn-sm" id="agg-show-query"><i class="fa fa-code"></i> Show Query</button>
          </div>
```

```html
<!-- AFTER -->
          <div class="agg-section" id="agg-order-section">
            <div class="agg-section-header">
              ORDER BY
              <button type="button" class="btn btn-default btn-xs" id="agg-add-order"><i class="fa fa-plus"></i> Add</button>
            </div>
            <div id="agg_order_rows"></div>
          </div>
          <div class="adv-search-footer">
            <button type="button" class="btn btn-primary btn-sm" id="agg-apply">Apply</button>
            <button type="button" class="btn btn-default btn-sm" id="agg-reset"><i class="fa fa-times"></i> Clear</button>
            <button type="button" class="btn btn-default btn-sm" id="agg-show-query"><i class="fa fa-code"></i> Show Query</button>
            <span class="panel-limit-wrap">LIMIT <input type="text" id="agg-limit-input" class="form-control panel-limit-input" placeholder="default" /></span>
          </div>
```

- [ ] **Step 3: Commit**

```bash
git add static/index.html
git commit -S -m "feat: add ORDER BY sections and LIMIT inputs to Advanced and Aggregate panels (HTML)"
```

---

## Task 2: Add CSS for ORDER BY rows and LIMIT inputs

**Files:**
- Modify: `static/css/app.css`

- [ ] **Step 1: Add styles**

Append after the `.agg-col-hidden` rule in `static/css/app.css`:

```css
.order-by-row {
  display: flex;
  align-items: center;
  margin-bottom: 4px;
  gap: 6px;
}

.order-by-row .order-by-col {
  flex: 1 1 auto;
  min-width: 0;
  height: 28px;
  font-size: 12px;
  padding: 2px 6px;
}

.order-by-row .order-by-dir {
  width: 70px;
  flex-shrink: 0;
  height: 28px;
  font-size: 12px;
  padding: 2px 6px;
}

.panel-limit-wrap {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 12px;
  color: #555;
  margin-left: 6px;
}

.panel-limit-input {
  width: 70px;
  height: 28px;
  font-size: 12px;
  padding: 2px 6px;
  display: inline-block;
}
```

- [ ] **Step 2: Commit**

```bash
git add static/css/app.css
git commit -S -m "feat: add CSS for ORDER BY rows and LIMIT inputs"
```

---

## Task 3: Add ORDER BY JS helpers

**Files:**
- Modify: `static/js/app.js` (after `updateHavingAliasDropdowns`, around line 1330)

- [ ] **Step 1: Add `buildOrderByRow`, `buildOrderByClause`, `updateAdvOrderByColumns`, `updateAggOrderByColumns`**

Insert after the closing `}` of `updateHavingAliasDropdowns()`:

```javascript
// Build a single ORDER BY row element. colHtml is the innerHTML for the column <select>.
function buildOrderByRow(colHtml) {
  var colSelect = $('<select class="order-by-col form-control"></select>').html(
    colHtml || '<option value="">Select column</option>'
  );
  var dirSelect = $('<select class="order-by-dir form-control"></select>').append(
    $('<option value="ASC">ASC</option>'),
    $('<option value="DESC">DESC</option>')
  );
  var removeBtn = $('<button type="button" class="btn btn-default btn-xs order-by-remove"><i class="fa fa-minus"></i></button>');
  return $('<div class="order-by-row"></div>').append(colSelect, dirSelect, removeBtn);
}

// Returns "ORDER BY \"col1\" ASC, \"col2\" DESC" or null if no valid rows.
// containerSelector: CSS selector for the element containing .order-by-row elements.
function buildOrderByClause(containerSelector) {
  var parts = [];
  $(containerSelector + " .order-by-row").each(function() {
    var col = $(this).find(".order-by-col").val() || "";
    var dir = $(this).find(".order-by-dir").val() || "ASC";
    if (!col) return;
    parts.push('"' + col + '" ' + dir);
  });
  if (parts.length === 0) return null;
  return "ORDER BY " + parts.join(", ");
}

// Sync column options into existing Advanced ORDER BY rows when the table column list changes.
function updateAdvOrderByColumns() {
  var colHtml = $("#pagination select.column").html() || '<option value="">Select column</option>';
  $("#adv_order_rows .order-by-row .order-by-col").each(function() {
    var prev = $(this).val();
    $(this).html(colHtml);
    $(this).val(prev);
  });
}

// Rebuild Aggregate ORDER BY column dropdowns: GROUP BY selected cols + aggregate aliases.
function updateAggOrderByColumns() {
  var options = [];
  var gbCols = $("#agg_group_by_select").val() || [];
  gbCols.forEach(function(c) { options.push(c); });
  getAggregateAliases().forEach(function(a) { options.push(a); });

  $("#agg_order_rows .order-by-row").each(function() {
    var sel  = $(this).find(".order-by-col");
    var prev = sel.val();
    sel.empty();
    if (options.length === 0) {
      sel.append('<option value=""></option>');
    } else {
      options.forEach(function(o) {
        sel.append($('<option></option>').val(o).text(o));
      });
    }
    if (prev && options.indexOf(prev) !== -1) sel.val(prev);
  });
}
```

- [ ] **Step 2: Commit**

```bash
git add static/js/app.js
git commit -S -m "feat: add buildOrderByRow, buildOrderByClause, updateAdvOrderByColumns, updateAggOrderByColumns"
```

---

## Task 4: Update `buildAggregateQuery` to include ORDER BY and panel LIMIT

**Files:**
- Modify: `static/js/app.js` (`buildAggregateQuery`, around line 1410)

- [ ] **Step 1: Replace `buildAggregateQuery`**

Replace the current function body:

```javascript
// BEFORE (lines ~1410-1441):
function buildAggregateQuery() {
  var obj   = getCurrentObject();
  var parts = obj.name.split(".");
  var tableRef;
  if (parts.length === 2) {
    tableRef = '"' + parts[0] + '"."' + parts[1] + '"';
  } else {
    tableRef = '"' + parts[0] + '"';
  }

  var selectClause = buildAggregateSelectClause();
  var groupBy      = buildGroupByClause();
  if (!groupBy) return null;

  var sql = "SELECT " + selectClause + " FROM " + tableRef;

  if (advancedSearchActive) {
    var where = $("#advanced_search_panel").data("where");
    if (where) sql += " WHERE " + where;
  }

  sql += " " + groupBy;

  var having = buildHavingClause();
  if (having) sql += " " + having;

  sql += " LIMIT " + getRowsLimit();
  var offset = getPaginationOffset();
  if (offset > 0) sql += " OFFSET " + offset;

  return sql;
}
```

```javascript
// AFTER:
function buildAggregateQuery() {
  var obj   = getCurrentObject();
  var parts = obj.name.split(".");
  var tableRef;
  if (parts.length === 2) {
    tableRef = '"' + parts[0] + '"."' + parts[1] + '"';
  } else {
    tableRef = '"' + parts[0] + '"';
  }

  var selectClause = buildAggregateSelectClause();
  var groupBy      = buildGroupByClause();
  if (!groupBy) return null;

  var sql = "SELECT " + selectClause + " FROM " + tableRef;

  if (advancedSearchActive) {
    var where = $("#advanced_search_panel").data("where");
    if (where) sql += " WHERE " + where;
  }

  sql += " " + groupBy;

  var having = buildHavingClause();
  if (having) sql += " " + having;

  var orderBy = buildOrderByClause("#agg_order_rows");
  if (orderBy) sql += " " + orderBy;

  var limitVal = parseInt($("#agg-limit-input").val(), 10);
  sql += " LIMIT " + (limitVal > 0 ? limitVal : getRowsLimit());

  var offset = getPaginationOffset();
  if (offset > 0) sql += " OFFSET " + offset;

  return sql;
}
```

- [ ] **Step 2: Commit**

```bash
git add static/js/app.js
git commit -S -m "feat: inject ORDER BY and panel LIMIT into buildAggregateQuery"
```

---

## Task 5: Add `buildAdvancedQuery` and switch `showTableContent` to use it

**Files:**
- Modify: `static/js/app.js`

- [ ] **Step 1: Add `buildAdvancedQuery` after `applyAdvancedSearch`**

Insert after the closing `}` of `applyAdvancedSearch()` (around line 1496):

```javascript
// Build the full SELECT query for Advanced search mode using stored WHERE state
// and current ORDER BY rows + LIMIT input from the panel DOM.
function buildAdvancedQuery() {
  var obj   = getCurrentObject();
  var parts = obj.name.split(".");
  var tableRef = parts.length === 2
    ? '"' + parts[0] + '"."' + parts[1] + '"'
    : '"' + parts[0] + '"';

  var sql = "SELECT * FROM " + tableRef;

  var where = $("#advanced_search_panel").data("where");
  if (where) sql += " WHERE " + where;

  var orderBy = buildOrderByClause("#adv_order_rows");
  if (orderBy) sql += " " + orderBy;

  var limitVal = parseInt($("#adv-limit-input").val(), 10);
  sql += " LIMIT " + (limitVal > 0 ? limitVal : getRowsLimit());

  var offset = getPaginationOffset();
  if (offset > 0) sql += " OFFSET " + offset;

  return sql;
}
```

- [ ] **Step 2: Replace the `advancedSearchActive` branch in `showTableContent`**

In `showTableContent`, replace:

```javascript
  // Advanced search takes precedence over the simple filter
  if (advancedSearchActive) {
    var advWhere = $("#advanced_search_panel").data("where");
    if (advWhere) opts["where"] = advWhere;
  } else {
```

```javascript
  if (advancedSearchActive) {
    var advQuery = buildAdvancedQuery();
    executeQuery(advQuery, function(data) {
      $("#input").hide();
      $("#body").prop("class", "with-pagination");
      adjustOutputTop();
      buildTable(data);
      setCurrentTab("table_content");
      updatePaginator(data.pagination);
      $("#results").data("mode", "browse").data("table", name);
    });
    return;
  }

  // Simple filter (no advanced search active)
  {
```

Also change the closing `}` of the else-block (which was `}`) to `}` — the wrapping braces now close the new block. Specifically: the old code had:

```javascript
  } else {
    var filter = { ... };
    if (filter.column && filter.op) {
      opts["where"] = where;
    }
  }

  getTableRows(name, opts, ...);
```

After the change, the `else` becomes a bare block and `getTableRows` is still called for the non-advanced path:

```javascript
  } else {
    var filter = {
      column: $(".filters select.column").val(),
      op:     $(".filters select.filter").val(),
      input:  $(".filters input").val()
    };

    if (filter.column && filter.op) {
      var where = [
        '"' + filter.column + '"',
        filterOptions[filter.op].replace("DATA", escapeSqlLiteral(filter.input))
      ].join(" ");

      opts["where"] = where;
    }
  }

  getTableRows(name, opts, function(data) { ... });
```

This block is unchanged — only the `advancedSearchActive` branch is replaced with the early return.

- [ ] **Step 3: Update `buildFullQuery` to include ORDER BY + LIMIT for the Advanced preview**

Replace `buildFullQuery`:

```javascript
// BEFORE:
function buildFullQuery() {
  if (aggregateActive) {
    var aggQ = buildAggregateQuery();
    return aggQ ? aggQ + ";" : null;
  }
  var where = buildAdvancedWhereClause();
  if (!where) return null;
  var table = getCurrentObject().name;
  var nameParts = table.split(".");
  var sql;
  if (nameParts.length === 2) {
    sql = 'SELECT * FROM "' + nameParts[0] + '"."' + nameParts[1] + '"';
  } else {
    sql = 'SELECT * FROM "' + table + '"';
  }
  return sql + " WHERE " + where + ";";
}
```

```javascript
// AFTER:
function buildFullQuery() {
  if (aggregateActive) {
    var aggQ = buildAggregateQuery();
    return aggQ ? aggQ + ";" : null;
  }
  // Advanced: build from DOM (WHERE from current row state, not stored state,
  // so the preview reflects unsaved edits too).
  var where = buildAdvancedWhereClause();
  if (!where) return null;
  var table = getCurrentObject().name;
  var nameParts = table.split(".");
  var tableRef = nameParts.length === 2
    ? '"' + nameParts[0] + '"."' + nameParts[1] + '"'
    : '"' + table + '"';
  var sql = "SELECT * FROM " + tableRef + " WHERE " + where;
  var orderBy = buildOrderByClause("#adv_order_rows");
  if (orderBy) sql += " " + orderBy;
  var limitVal = parseInt($("#adv-limit-input").val(), 10);
  sql += " LIMIT " + (limitVal > 0 ? limitVal : getRowsLimit());
  return sql + ";";
}
```

- [ ] **Step 4: Commit**

```bash
git add static/js/app.js
git commit -S -m "feat: add buildAdvancedQuery and switch showTableContent to client-side SQL for advanced mode"
```

---

## Task 6: Update reset functions and add event handlers

**Files:**
- Modify: `static/js/app.js`

- [ ] **Step 1: Update `resetAdvancedSearch` to clear ORDER BY rows and LIMIT**

Replace:

```javascript
function resetAdvancedSearch() {
  advancedSearchActive = false;
  $("#advanced_search_panel").data("where", null);
  $("#adv_search_active_badge").hide();
  $("#advanced-search-toggle").removeClass("adv-open");

  $("#adv_search_rows").empty();
  $("#adv_search_rows").append(buildAdvancedSearchRow(true));
}
```

```javascript
function resetAdvancedSearch() {
  advancedSearchActive = false;
  $("#advanced_search_panel").data("where", null);
  $("#adv_search_active_badge").hide();
  $("#advanced-search-toggle").removeClass("adv-open");

  $("#adv_search_rows").empty();
  $("#adv_search_rows").append(buildAdvancedSearchRow(true));
  $("#adv_order_rows").empty();
  $("#adv-limit-input").val("");
}
```

- [ ] **Step 2: Update `resetAggregate` to clear ORDER BY rows and LIMIT**

Replace:

```javascript
function resetAggregate() {
  aggregateActive = false;
  $("#agg_active_badge").hide();
  $("#agg_group_by_select").val([]);
  $("#agg_expr_rows").empty();
  $("#agg_having_rows").empty();
  $("#aggregate_panel").hide();
  $("#aggregate-toggle").removeClass("agg-open");
  $("#pagination").removeClass("agg-panel-open");
```

```javascript
function resetAggregate() {
  aggregateActive = false;
  $("#agg_active_badge").hide();
  $("#agg_group_by_select").val([]);
  $("#agg_expr_rows").empty();
  $("#agg_having_rows").empty();
  $("#agg_order_rows").empty();
  $("#agg-limit-input").val("");
  $("#aggregate_panel").hide();
  $("#aggregate-toggle").removeClass("agg-open");
  $("#pagination").removeClass("agg-panel-open");
```

- [ ] **Step 3: Add Advanced ORDER BY event handlers**

After the `#adv-add-condition` click handler, add:

```javascript
  // Add ORDER BY row to Advanced panel
  $("#adv-add-order").on("click", function() {
    var colHtml = $("#pagination select.column").html() || '<option value="">Select column</option>';
    $("#adv_order_rows").append(buildOrderByRow(colHtml));
    adjustOutputTop();
  });

  // Remove ORDER BY row from Advanced panel
  $("#adv_order_rows").on("click", ".order-by-remove", function() {
    $(this).closest(".order-by-row").remove();
    adjustOutputTop();
  });
```

- [ ] **Step 4: Add Aggregate ORDER BY event handlers and wire column sync**

After the `#agg-add-having` click handler, add:

```javascript
  // Add ORDER BY row to Aggregate panel
  $("#agg-add-order").on("click", function() {
    var options = [];
    var gbCols = $("#agg_group_by_select").val() || [];
    gbCols.forEach(function(c) { options.push(c); });
    getAggregateAliases().forEach(function(a) { options.push(a); });
    var colHtml = options.map(function(o) {
      return '<option value="' + o + '">' + o + '</option>';
    }).join("") || '<option value="">Select column</option>';
    $("#agg_order_rows").append(buildOrderByRow(colHtml));
  });

  // Remove ORDER BY row from Aggregate panel
  $("#agg_order_rows").on("click", ".order-by-remove", function() {
    $(this).closest(".order-by-row").remove();
  });
```

- [ ] **Step 5: Call `updateAggOrderByColumns` wherever `updateHavingAliasDropdowns` is called**

Find all three call sites of `updateHavingAliasDropdowns()` and add `updateAggOrderByColumns()` immediately after each:

```javascript
// Site 1 — agg-fn change handler
    updateHavingAliasDropdowns();
    updateAggOrderByColumns();

// Site 2 — agg-alias input handler
    updateHavingAliasDropdowns();
    updateAggOrderByColumns();

// Site 3 — agg-remove-row click handler
    updateHavingAliasDropdowns();
    updateAggOrderByColumns();
```

- [ ] **Step 6: Call `updateAdvOrderByColumns` wherever `#adv_search_rows .adv-col` columns are synced**

Find the block around line 1101 that syncs `.adv-col` dropdowns and add the call after it:

```javascript
    // Sync columns into any existing advanced search rows
    var colHtml = $("#pagination select.column").html();
    $("#adv_search_rows .adv-col").each(function() {
      var prev = $(this).val();
      $(this).html(colHtml);
      $(this).val(prev);
    });
    updateAdvOrderByColumns();   // ← add this line
```

Also call `updateAggOrderByColumns()` wherever `buildGroupBySection()` is called (GROUP BY selection changes affect the ORDER BY column list):

```javascript
    buildGroupBySection();
    updateAggOrderByColumns();   // ← add after every buildGroupBySection() call
```

Find all `buildGroupBySection()` call sites and add the line after each.

- [ ] **Step 7: Commit**

```bash
git add static/js/app.js
git commit -S -m "feat: wire ORDER BY event handlers, reset functions, and column sync for Advanced and Aggregate"
```

---

## Task 7: Build and verify

- [ ] **Step 1: Build**

```bash
GOROOT=/opt/homebrew/Cellar/go/1.26.1/libexec \
GOPROXY=https://proxy.golang.org,direct \
GONOSUMDB='*' \
GOOS=linux GOARCH=amd64 \
go build -o ./bin/pgweb_linux_amd64
```

Expected: exits 0, `bin/pgweb_linux_amd64` updated.

- [ ] **Step 2: Manual smoke-test checklist**

Start the server and open a table's Rows tab. Verify:

**Advanced panel:**
1. Click `+ Add` in ORDER BY section → row appears with column dropdown + ASC/DESC
2. Select a column, set DESC, click Apply → query runs; click Show Query → `ORDER BY "col" DESC` appears in SQL
3. Add a second ORDER BY row → `ORDER BY "col1" DESC, "col2" ASC` in SQL
4. Set LIMIT to `5` → `LIMIT 5` appears in SQL and only 5 rows return
5. Click Clear → ORDER BY rows gone, LIMIT input empty
6. Changing table reloads columns into ORDER BY dropdowns correctly

**Aggregate panel:**
1. Select GROUP BY columns, add aggregate expression
2. Click `+ Add` in ORDER BY section → dropdown contains GROUP BY cols + aggregate aliases
3. Set `count DESC`, click Apply → `ORDER BY "count" DESC` in generated SQL (Show Query)
4. Set LIMIT to `3` → `LIMIT 3` in SQL
5. Click Clear → ORDER BY rows gone, LIMIT empty
6. Adding/removing/renaming aggregate expressions updates ORDER BY dropdowns

- [ ] **Step 3: Commit (if any fixups were needed)**

```bash
git add static/js/app.js static/css/app.css static/index.html
git commit -S -m "fix: ORDER BY + LIMIT post-smoke-test fixups"
```

---

## Self-Review Notes

- **Spec coverage:** ORDER BY for both panels ✓ | Multi-column ✓ | LIMIT override ✓ | Advanced switches to client-side SQL ✓ | Show Query reflects new clauses ✓ | Resets clear new state ✓ | Column sync on table change ✓
- **No placeholders:** All steps include exact code
- **Type consistency:** `buildOrderByClause` signature `(containerSelector)` used consistently in Tasks 3, 4, 5 | `buildOrderByRow(colHtml)` used consistently in Tasks 3, 6 | `updateAggOrderByColumns()` / `updateAdvOrderByColumns()` zero-arg calls throughout
- **Edge case noted:** When `advancedSearchActive` is true, column-header click sort is ignored in favour of the panel ORDER BY rows (consistent with aggregate mode behaviour)
