# Group By Dropdown Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the GROUP BY checkbox grid in the aggregate panel with a native `<select multiple>` dropdown.

**Architecture:** Three files change in concert — HTML swaps the container element, JS updates four functions to read from the select instead of checkboxes, CSS removes the old grid rules and adds two new rules for the select and hint. No new dependencies; no other logic is touched.

**Tech Stack:** jQuery, Bootstrap 3, plain HTML/CSS. No JS test framework exists for the frontend — verification is via Go build success and manual browser inspection.

---

### Task 1: Replace GROUP BY HTML

**Files:**
- Modify: `static/index.html:156-158`

- [ ] **Step 1: Replace the checkbox grid div with a `<select multiple>` and hint**

In `static/index.html`, replace lines 156-158:

```html
            <div id="agg_group_by_grid" class="agg-group-by-grid">
              <!-- Group-by checkboxes injected by buildGroupBySection() -->
            </div>
```

with:

```html
            <select id="agg_group_by_select" multiple class="form-control agg-group-by-select"></select>
            <div class="agg-group-by-hint">Hold Ctrl / Cmd to select multiple columns</div>
```

- [ ] **Step 2: Commit**

```bash
git add static/index.html
git commit -S -m "feat: replace GROUP BY checkbox grid with multi-select element"
```

---

### Task 2: Update CSS

**Files:**
- Modify: `static/css/app.css:1140-1154`

- [ ] **Step 1: Remove old grid rules and add new select rules**

In `static/css/app.css`, replace lines 1140-1154:

```css
.agg-group-by-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 6px 16px;
}

.agg-group-by-grid label {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: normal;
  margin: 0;
  cursor: pointer;
}
```

with:

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

- [ ] **Step 2: Commit**

```bash
git add static/css/app.css
git commit -S -m "feat: swap GROUP BY grid CSS for multi-select styles"
```

---

### Task 3: Update `buildGroupBySection()`

**Files:**
- Modify: `static/js/app.js:1181-1192`

- [ ] **Step 1: Rewrite to populate `<option>` elements**

In `static/js/app.js`, replace lines 1181-1192:

```js
function buildGroupBySection() {
  var grid = $("#agg_group_by_grid");
  grid.empty();
  $("#pagination select.column option").each(function() {
    var val = $(this).val();
    if (!val) return;
    var label = $("<label>");
    var cb = $("<input>", { type: "checkbox", class: "agg-group-col", value: val });
    label.append(cb).append(" " + val);
    grid.append(label);
  });
}
```

with:

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

- [ ] **Step 2: Commit**

```bash
git add static/js/app.js
git commit -S -m "feat: rewrite buildGroupBySection to populate select options"
```

---

### Task 4: Update `buildGroupByClause()`

**Files:**
- Modify: `static/js/app.js:1321-1328`

- [ ] **Step 1: Read from multi-select `.val()` instead of checked checkboxes**

In `static/js/app.js`, replace lines 1321-1328:

```js
function buildGroupByClause() {
  var cols = [];
  $(".agg-group-col:checked").each(function() {
    cols.push('"' + $(this).val() + '"');
  });
  if (cols.length === 0) return null;
  return "GROUP BY " + cols.join(", ");
}
```

with:

```js
function buildGroupByClause() {
  var cols = $("#agg_group_by_select").val() || [];
  if (cols.length === 0) return null;
  return "GROUP BY " + cols.map(function(c) { return '"' + c + '"'; }).join(", ");
}
```

- [ ] **Step 2: Commit**

```bash
git add static/js/app.js
git commit -S -m "feat: update buildGroupByClause to read from multi-select"
```

---

### Task 5: Update `buildAggregateSelectClause()`

**Files:**
- Modify: `static/js/app.js:1298-1302`

- [ ] **Step 1: Replace `.agg-group-col:checked` iteration with `.val()` array**

In `static/js/app.js`, replace lines 1298-1302:

```js
function buildAggregateSelectClause() {
  var parts = [];
  $(".agg-group-col:checked").each(function() {
    parts.push('"' + $(this).val() + '"');
  });
```

with:

```js
function buildAggregateSelectClause() {
  var parts = [];
  var cols = $("#agg_group_by_select").val() || [];
  cols.forEach(function(c) { parts.push('"' + c + '"'); });
```

- [ ] **Step 2: Commit**

```bash
git add static/js/app.js
git commit -S -m "feat: update buildAggregateSelectClause to read from multi-select"
```

---

### Task 6: Update `resetAggregate()`

**Files:**
- Modify: `static/js/app.js:1493`

- [ ] **Step 1: Replace `buildGroupBySection()` call with deselect-all**

In `static/js/app.js`, replace line 1493:

```js
  buildGroupBySection();
```

with:

```js
  $("#agg_group_by_select").val([]);
```

The full function after the change:

```js
function resetAggregate() {
  aggregateActive = false;
  $("#agg_active_badge").hide();
  $("#agg_group_by_select").val([]);
  $("#agg_expr_rows").empty();
  $("#agg_having_rows").empty();
  $("#aggregate_panel").hide();
  $("#aggregate-toggle").removeClass("agg-open");
  $("#pagination").removeClass("agg-panel-open");
  adjustOutputTop();
}
```

- [ ] **Step 2: Commit**

```bash
git add static/js/app.js
git commit -S -m "feat: update resetAggregate to deselect multi-select instead of rebuilding"
```

---

### Task 7: Build and verify

- [ ] **Step 1: Build the binary**

```bash
GOROOT=/opt/homebrew/Cellar/go/1.26.1/libexec GOPROXY=https://proxy.golang.org,direct GONOSUMDB='*' GOOS=linux GOARCH=amd64 go build -o ./bin/pgweb_linux_amd64
```

Expected: exits 0, no output.

- [ ] **Step 2: Run Go tests**

```bash
GOROOT=/opt/homebrew/Cellar/go/1.26.1/libexec go test ./... 2>&1 | grep -v "^?"
```

Expected: all packages pass except `pkg/client` (requires a live PostgreSQL connection — pre-existing infrastructure failure, not related to this change).

- [ ] **Step 3: Manual smoke check**

Start the server locally and open a table:
1. Click **Aggregate** button — GROUP BY section shows a `<select>` populated with column names
2. Select one column (single-click) → click **Apply** → results group correctly
3. Select multiple columns (Ctrl+click) → click **Apply** → GROUP BY clause includes all selected columns
4. Click **Show Query** → SQL reflects selected columns in GROUP BY
5. Click **Clear** → select has no selection, aggregate mode deactivates

- [ ] **Step 4: Invoke `superpowers:finishing-a-development-branch`**
