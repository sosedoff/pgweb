---
name: bug-query-tab-buttons
description: Run Query / Explain Query buttons unclickable after switching from Rows tab — root cause and fix
metadata:
  type: project
---

## Bug: Query tab buttons blocked after using Rows tab

**Symptom:** After using the Rows tab (especially with Advanced or Aggregate filters), switching to the Query tab leaves the Run Query and Explain Query buttons unclickable — cursor stays as arrow, clicks don't register. Refreshing the page restores them.

**Root cause:** `adjustOutputTop()` (called by `showTableContent`) sets `#output` top CSS to the pagination panel height (e.g. 150px) so results clear the filter UI on the Rows tab. When switching to the Query tab, `showQueryPanel()` shows `#input` but never resets `#output` top. With `#output` positioned too high (position: absolute within #body), it overlaps and invisibly intercepts all mouse events on the `.actions` buttons inside `#input`.

**Fix:** One line added to `showQueryPanel()` in `static/js/app.js`:
```javascript
$("#output").css("top", $("#input").height() + "px");
```
This resets `#output` top to align with the bottom of `#input` when entering the Query tab.

**Why:** `adjustOutputTop()` is only called for Rows tab layout; no corresponding reset existed for Query tab. The `#output` div (overflow:auto, position:absolute) silently covered the button bar.

**How to apply:** When debugging future UI event issues (cursor, click failures) on the Query tab after Rows tab use — check `#output` top CSS first.
