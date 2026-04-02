# HAVING Row Remove Button Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a remove (−) button to every HAVING row in the aggregate panel, including the first row, and correctly promote the new first row's display when the original first row is deleted.

**Architecture:** Two surgical edits to `static/js/app.js` — remove the `isFirst` guard in `buildHavingRow()` so all rows get a remove button, and extend the existing `.adv-remove-row` click handler for `#agg_having_rows` to repair the first-row conjunction display after deletion. No HTML or CSS changes needed.

**Tech Stack:** jQuery, Bootstrap 3, plain JS. No frontend test framework — verification is via Go build success and manual browser inspection.

---

### Task 1: Add remove button to all HAVING rows

**Files:**
- Modify: `static/js/app.js:1278-1280`

- [ ] **Step 1: Apply the change**

In `static/js/app.js`, find `buildHavingRow`. The end of the function currently reads:

```js
  row.append(exprSelect).append(opSelect).append(valInput);

  if (!isFirst) {
    row.append('<button type="button" class="btn btn-default btn-xs adv-remove-row"><i class="fa fa-minus"></i></button>');
  }

  return row;
```

Replace with:

```js
  row.append(exprSelect).append(opSelect).append(valInput);

  row.append('<button type="button" class="btn btn-default btn-xs adv-remove-row"><i class="fa fa-minus"></i></button>');

  return row;
```

- [ ] **Step 2: Commit**

```bash
git add static/js/app.js
git commit -S -m "feat: add remove button to first HAVING row"
```

---

### Task 2: Fix first-row display after first-row deletion

**Files:**
- Modify: `static/js/app.js:2307-2309`

- [ ] **Step 1: Apply the change**

In `static/js/app.js`, find the `#agg_having_rows` remove handler. It currently reads:

```js
  $("#agg_having_rows").on("click", ".adv-remove-row", function() {
    $(this).closest(".adv-search-row").remove();
  });
```

Replace with:

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

- [ ] **Step 2: Commit**

```bash
git add static/js/app.js
git commit -S -m "fix: promote new first HAVING row after first-row deletion"
```

---

### Task 3: Build and verify

- [ ] **Step 1: Build**

```bash
GOROOT=/opt/homebrew/Cellar/go/1.26.1/libexec GOPROXY=https://proxy.golang.org,direct GONOSUMDB='*' GOOS=linux GOARCH=amd64 go build -o ./bin/pgweb_linux_amd64
```

Expected: exits 0, no output.

- [ ] **Step 2: Run Go tests**

```bash
GOROOT=/opt/homebrew/Cellar/go/1.26.1/libexec go test ./... 2>&1 | grep -v "^?"
```

Expected: all packages pass except `pkg/client` (requires live PostgreSQL — pre-existing).

- [ ] **Step 3: Manual smoke check**

1. Open Aggregate panel on any table
2. Click **HAVING +Add** → first row appears with a `−` button
3. Click `−` on the first row → row is removed, panel is empty again
4. Click **HAVING +Add** twice → two rows appear, both with `−` buttons
5. Click `−` on the first row → second row becomes first; its label changes to "WHERE" (no AND/OR buttons)
6. Click `−` on that row → panel is empty
7. Add two rows, remove the second (non-first) row → first row unchanged, second removed

- [ ] **Step 4: Invoke `superpowers:finishing-a-development-branch`**
