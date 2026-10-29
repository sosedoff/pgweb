var appInfo             = {};
var appFeatures         = {};
var editor              = null;
var connected           = false;
var bookmarks           = {};
var default_rows_limit  = 100;
var currentObject       = null;
var autocompleteObjects = [];
var inputResizing       = false;
var inputResizeOffset   = null;
var currentMetadata     = null;
var metadataCache       = {};
var tableFilters        = [];
var rowModalContext     = null;

var filterOptions = {
  "equal":      "= 'DATA'",
  "not_equal":  "!= 'DATA'",
  "greater":    "> 'DATA'" ,
  "greater_eq": ">= 'DATA'",
  "less":       "< 'DATA'",
  "less_eq":    "<= 'DATA'",
  "like":       "LIKE 'DATA'",
  "ilike":      "ILIKE 'DATA'",
  "in":         "IN (DATA)",
  "null":       "IS NULL",
  "not_null":   "IS NOT NULL"
};

// Human readable labels for the active filter chips
var filterLabels = {
  "equal":      "=",
  "not_equal":  "\u2260",
  "greater":    ">",
  "greater_eq": "\u2265",
  "less":       "<",
  "less_eq":    "\u2264",
  "like":       "LIKE",
  "ilike":      "ILIKE",
  "in":         "IN",
  "null":       "IS NULL",
  "not_null":   "IS NOT NULL"
};

// Operators that do not require a value
function filterOpHasValue(op) {
  return ["null", "not_null"].indexOf(op) < 0;
}

function getSessionId() {
  var id = sessionStorage.getItem("session_id");

  if (!id) {
    id = guid();
    sessionStorage.setItem("session_id", id);
  }

  return id;
}

function setRowsLimit(num) {
  localStorage.setItem("rows_limit", num);
}

function getRowsLimit() {
  return parseInt(localStorage.getItem("rows_limit") || default_rows_limit);
}

function getPaginationOffset() {
  var page  = $(".current-page").data("page");
  var limit = getRowsLimit();
  return (page - 1) * limit;
}

function getPagesCount(rowsCount) {
  var limit = getRowsLimit();
  var num = parseInt(rowsCount / limit);

  if ((num * limit) < rowsCount) {
    num++;
  }

  return num;
}

function apiCall(method, path, params, cb) {
  apiRequest(method, path, { data: params }, cb);
}

// apiJsonCall sends a JSON encoded body, which is required by the row editing
// endpoints since they need typed values and explicit nulls.
function apiJsonCall(method, path, body, cb) {
  apiRequest(method, path, {
    data: JSON.stringify(body),
    contentType: "application/json"
  }, cb);
}

function apiRequest(method, path, settings, cb) {
  var timeout = appFeatures.query_timeout;
  if (timeout == null) {
    timeout = 300; // in seconds
  }

  $.ajax($.extend({
    timeout: timeout * 1000, // in milliseconds
    url: "api" + path,
    method: method,
    cache: false,
    headers: {
      "x-session-id": getSessionId()
    },
    success: cb,
    error: function(xhr, status, data) {
      switch(status) {
        case "error":
          if (xhr.readyState == 0) { // 0 = UNSENT
            showErrorBanner("Sorry, something went wrong with your request. Refresh the page and try again!");
          }
          break;
        case "timeout":
          return cb({ error: "Query timeout after " + timeout + "s" });
      }

      var responseText;
      try {
        responseText = jQuery.parseJSON(xhr.responseText);
      }
      catch {
        responseText = { error: "Failed to parse the JSON response." };
      }
      cb(responseText);
    }
  }, settings));
}

function getInfo(cb)                        { apiCall("get", "/info", {}, cb); }
function getConnection(cb)                  { apiCall("get", "/connection", {}, cb); }
function getServerSettings(cb)              { apiCall("get", "/server_settings", {}, cb); }
function getSchemas(cb)                     { apiCall("get", "/schemas", {}, cb); }
function getObjects(cb)                     { apiCall("get", "/objects", {}, cb); }
function getTables(cb)                      { apiCall("get", "/tables", {}, cb); }
function getTableRows(table, opts, cb)      { apiCall("get", "/tables/" + table + "/rows", opts, cb); }
function getTableStructure(table, opts, cb) { apiCall("get", "/tables/" + table, opts, cb); }
function getTableMetadata(table, cb)        { apiCall("get", "/tables/" + table + "/metadata", {}, cb); }
function createTableRow(table, values, cb)  { apiJsonCall("post", "/tables/" + table + "/rows", { values: values }, cb); }
function updateTableRow(table, key, values, cb) { apiJsonCall("put", "/tables/" + table + "/rows", { key: key, values: values }, cb); }
function deleteTableRow(table, key, cb)     { apiJsonCall("delete", "/tables/" + table + "/rows", { key: key }, cb); }
function getTableIndexes(table, cb)         { apiCall("get", "/tables/" + table + "/indexes", {}, cb); }
function getTableConstraints(table, cb)     { apiCall("get", "/tables/" + table + "/constraints", {}, cb); }
function getTablesStats(cb)                 { apiCall("get", "/tables_stats", {}, cb); }
function getFunction(id, cb)                { apiCall("get", "/functions/" + id, {}, cb); }
function getHistory(cb)                     { apiCall("get", "/history", {}, cb); }
function getBookmarks(cb)                   { apiCall("get", "/bookmarks", {}, cb); }
function executeQuery(query, cb)            { apiCall("post", "/query", { query: query }, cb); }
function explainQuery(query, cb)            { apiCall("post", "/explain", { query: query }, cb); }
function analyzeQuery(query, cb)            { apiCall("post", "/analyze", { query: query }, cb); }
function disconnect(cb)                     { apiCall("post", "/disconnect", {}, cb); }

function encodeQuery(query) {
  return Base64.encode(query).replace(/\+/g, "-").replace(/\//g, "_").replace(/=/g, ".");
}

function showErrorBanner(text) {
  if (window.errBannerTimeout != null) {
    clearTimeout(window.errBannerTimeout);
  }

  window.errBannerTimeout = setTimeout(function() {
    $("#error_banner").fadeOut("fast").text("");
  }, 3000);

  $("#error_banner").removeClass("success").text(text).show();
}

function showSuccessBanner(text) {
  if (window.errBannerTimeout != null) {
    clearTimeout(window.errBannerTimeout);
  }

  window.errBannerTimeout = setTimeout(function() {
    $("#error_banner").fadeOut("fast").text("").removeClass("success");
  }, 3000);

  $("#error_banner").addClass("success").text(text).show();
}

function buildSchemaSection(name, objects) {
  var section = "";

  var titles = {
    "table":             "Tables",
    "view":              "Views",
    "materialized_view": "Materialized Views",
    "function":          "Functions",
    "sequence":          "Sequences"
  };

  var icons = {
    "table":             '<i class="fa fa-table"></i>',
    "view":              '<i class="fa fa-table"></i>',
    "materialized_view": '<i class="fa fa-table"></i>',
    "function":          '<i class="fa fa-bolt"></i>',
    "sequence":          '<i class="fa fa-circle-o"></i>'
  };

  var klass = "";
  if (name == "public") klass = "expanded";

  section += "<div class='schema " + klass + "'>";
  section += "<div class='schema-name'><i class='fa fa-folder-o'></i><i class='fa fa-folder-open-o'></i> " + name + "</div>";
  section += "<div class='schema-container'>";

  ["table", "view", "materialized_view", "function", "sequence"].forEach(function(group) {
    group_klass = "";
    if (name == "public" && group == "table") group_klass = "expanded";

    section += "<div class='schema-group " + group_klass + "'>";
    section += "<div class='schema-group-title'><i class='fa fa-chevron-right'></i><i class='fa fa-chevron-down'></i> " + titles[group] + " <span class='schema-group-count'>" + objects[group].length + "</span></div>";
    section += "<ul data-group='" + group + "'>";

    if (objects[group]) {
      objects[group].forEach(function(item) {
        var id = name + "." + item.name;

        // Use function OID since multiple functions with the same name might exist
        if (group == "function") {
          id = item.oid;
        }

        section += "<li class='schema-item schema-" + group + "' data-type='" + group + "' data-id='" + id + "' data-name='" + item.name + "'>" + icons[group] + "&nbsp;" + item.name + "</li>";
      });
      section += "</ul></div>";
    }
  });

  section += "</div></div>";

  return section;
}

function loadLocalQueries() {
  if (!appFeatures.local_queries) return;

  $("body").on("click", "a.load-local-query", function(e) {
    var id = $(this).data("id");

    apiCall("get", "/local_queries/" + id, {}, function(resp) {
      editor.setValue(resp.query);
      editor.clearSelection();
    });
  });

  apiCall("get", "/local_queries", {}, function(resp) {
    if (resp.error) return;

    var container = $("#load-query-dropdown").find(".dropdown-menu");

    resp.forEach(function(item) {
      var title = item.title || item.id;
      $("<li><a href='#' class='load-local-query' data-id='" + item.id + "'>" + title + "</a></li>").appendTo(container);
    });

    if (resp.length > 0) $("#load-local-query").prop("disabled", "");
    $("#load-query-dropdown").show();
  });
}

function loadSchemas() {
  $("#objects").html("");
  resetMetadataCache();

  var emptyObjectList = function() {
    return {
      table: [],
      view: [],
      materialized_view: [],
      function: [],
      sequence: []
    }
  }

  getSchemas(function(schemasData) {
    if (schemasData.error) {
      alert("Error while fetching schemas: " + schemasData.error);
      return;
    }

    getObjects(function(data) {
      if (data.error) {
        alert("Error while fetching database objects: " + data.error);
        return;
      }

      if (Object.keys(data).length == 0) {
        data["public"] = emptyObjectList();
      }

      for (schemaName of schemasData) {
        // Allow users to see empty schemas if we dont have any objects in them
        if (!data[schemaName]) {
          data[schemaName] = emptyObjectList();
        }

        $(buildSchemaSection(schemaName, data[schemaName])).appendTo("#objects");
      }

      if (Object.keys(data).length == 1) {
        $(".schema").addClass("expanded");
      }

      // Clear out all autocomplete objects
      autocompleteObjects = [];
      for (schema in data) {
        for (kind in data[schema]) {
          if (!(kind == "table" || kind == "view" || kind == "materialized_view" || kind == "function")) {
            continue
          }

          for (item in data[schema][kind]) {
            autocompleteObjects.push({
              caption: data[schema][kind][item].name,
              value: data[schema][kind][item].name,
              meta: kind
            });
          }
        }
      }

      bindContextMenus();
    });
  });
}

function escapeHtml(str) {
  if (str != null || str != undefined) {
    return jQuery("<div/>").text(str).html();
  }

  return "<span class='null'>null</span>";
}

function unescapeHtml(str){
  var e = document.createElement("div");
  e.innerHTML = str;
  return e.childNodes.length === 0 ? "" : e.childNodes[0].nodeValue;
}

function getCurrentObject() {
  return currentObject || { name: "", type: "" };
}

function resetTable() {
  $("#results_header").html("");
  $("#results_body").html("");
  $("#results_view").html("").hide();

  $("#results").
    data("mode", "").
    removeClass("empty").
    removeClass("no-crop").
    show();
}

function performTableAction(table, action, el) {
  if (action == "truncate" || action == "delete") {
    var message = "Are you sure you want to " + action + " table " + table + " ?";
    if (!confirm(message)) return;
  }

  switch(action) {
    case "truncate":
      executeQuery("TRUNCATE TABLE " + table, function(data) {
        if (data.error) alert(data.error);
        resetTable();
      });
      break;
    case "delete":
      executeQuery("DROP TABLE " + table, function(data) {
        if (data.error) alert(data.error);
        loadSchemas();
        resetTable();
      });
      break;
    case "export":
      var format = el.data("format");
      var db = $("#current_database").text();
      var filename = db + "." + table + "." + format;
      var query = "SELECT * FROM " + table;
      openInNewWindow("api/query", { "format": format, "filename": filename, "query": query });
      break;
    case "dump":
      openInNewWindow("api/export", { "table": table });
      break;
    case "copy":
      copyToClipboard(table.split('.')[1]);
      break;
    case "analyze":
      executeQuery("ANALYZE " + table, function(data) {
        if (data.error) alert(data.error);
        resetTable();
      });
      break;
  }
}

function performViewAction(view, action, el) {
  if (action == "delete") {
    var message = "Are you sure you want to " + action + " view " + view + " ?";
    if (!confirm(message)) return;
  }

  switch(action) {
    case "delete":
      executeQuery("DROP VIEW " + view, function(data) {
        if (data.error) alert(data.error);
        loadSchemas();
        resetTable();
      });
      break;
    case "export":
      var format = el.data("format");
      var db = $("#current_database").text();
      var filename = db + "." + view + "." + format;
      var query = "SELECT * FROM " + view;
      openInNewWindow("api/query", { "format": format, "filename": filename, "query": query });
      break;
    case "copy":
      copyToClipboard(view.split('.')[1]);
      break;
    case "copy_def":
      executeQuery("SELECT pg_get_viewdef('" + view + "', true);", function(data) {
        if (data.error) {
          alert(data.error);
          return;
        }
        copyToClipboard(data.rows[0]);
      });
      break;
    case "view_def":
      executeQuery("SELECT pg_get_viewdef('" + view + "', true);", function(data) {
        if (data.error) {
          alert(data.error);
          return;
        }
        showViewDefinition(view, data.rows[0]);
      });
      break;
  }
}

function performRowAction(action, value) {
  if (action == "stop_query") {
    if (!confirm("Are you sure you want to stop the query?")) return;
    executeQuery("SELECT pg_cancel_backend(" + value + ");", function(data) {
      if (data.error) alert(data.error);
      setTimeout(showActivityPanel, 1000);
    });
  }
}

function sortArrow(direction) {
  switch (direction) {
    case "ASC":
      return "&#x25B2;";
    case "DESC":
      return "&#x25BC;";
    default:
      return "";
  }
}

function buildTable(results, sortColumn, sortOrder, options) {
  if (!options) options = {};
  var action = options.action;

  resetTable();

  if (results.error) {
    $("#results_header").html("");
    $("#results_body").html("<tr><td>ERROR: " + results.error + "</tr></tr>");
    return;
  }

  if (results.rows.length == 0) {
    $("#results_header").html("");
    $("#results_body").html("<tr><td>No records found</td></tr>");
    if (results.stats) {
      $("#result-rows-count").html(results.stats.query_duration_ms + " ms");
    } else {
      $("#result-rows-count").html("");
    }
    $("#results").addClass("empty");
    return;
  }

  var cols = "";
  var rows = "";
  var rowKeys = [];

  // Primary key columns are tracked so that the grid rows can be addressed
  // individually by the inline editor and the row action buttons.
  var pkColumns = options.pkColumns || [];
  var pkIndexes = pkColumns.map(function(col) { return results.columns.indexOf(col); });
  var trackKeys = pkColumns.length > 0 && pkIndexes.indexOf(-1) < 0;

  results.columns.forEach(function(col) {
    if (col === sortColumn) {
      cols += "<th class='table-header-col active' data-name='" + col + "' data-order=" + sortOrder + ">" + col + "&nbsp;" + sortArrow(sortOrder) + "</th>";
    } else {
      cols += "<th class='table-header-col' data-name='" + col + "'>" + col + "</th>";
    }
  });

  // No header to make the column non-sortable
  if (action) {
    cols += "<th></th>";

    // Determine which column contains the data attribute
    action.dataColumn = results.columns.indexOf(action.data);
  }

  results.rows.forEach(function(row) {
    var r = "";

    if (trackKeys) {
      var key = {};
      pkColumns.forEach(function(col, idx) {
        key[col] = row[pkIndexes[idx]];
      });
      rowKeys.push(key);
    }

    // Add all actual row data here
    for (i in row) {
      r += "<td data-col='" + i + "'><div>" + escapeHtml(row[i]) + "</div></td>";
    }

    // Add row action button
    if (action) {
      r += "<td><a class='btn btn-xs btn-" + action.style + " row-action' data-action='" + action.name + "' data-value='" + row[action.dataColumn] + "' href='#'>" + action.title + "</a></td>";
    }

    rows += "<tr>" + r + "</tr>";
  });

  $("#results_header").html(cols);
  $("#results_body").html(rows);
  $("#results").data("columns", results.columns);

  // Attach the keys via jQuery data to avoid any HTML attribute escaping issues
  if (trackKeys) {
    $("#results_body > tr").each(function(idx) {
      $(this).data("pk", rowKeys[idx]);
    });
  }

  // Show number of rows rendered on the page
  if (results.stats) {
    $("#result-rows-count").html(results.stats.rows_count + " rows in " + results.stats.query_duration_ms + " ms");
  } else {
    $("#result-rows-count").html(results.rows.length + " rows");
  }
}

function setCurrentTab(id) {
  // Pagination should only be visible on rows tab
  if (id != "table_content") {
    $("#body").removeClass("with-pagination");
  }

  $("#nav ul li.selected").removeClass("selected");
  $("#" + id).addClass("selected");

  // Persist tab selection into the session storage
  sessionStorage.setItem("tab", id);
}

function showQueryHistory() {
  getHistory(function(data) {
    var rows = [];

    for(i in data) {
      rows.unshift([parseInt(i) + 1, data[i].query, data[i].timestamp]);
    }

    buildTable({ columns: ["id", "query", "timestamp"], rows: rows });

    setCurrentTab("table_history");
    $("#input").hide();
    $("#body").prop("class", "full");
    $("#results").addClass("no-crop");
  });
}

function showTableIndexes() {
  var name = getCurrentObject().name;

  if (name.length == 0) {
    alert("Please select a table!");
    return;
  }

  getTableIndexes(name, function(data) {
    setCurrentTab("table_indexes");
    buildTable(data);

    $("#input").hide();
    $("#body").prop("class", "full");
    $("#results").addClass("no-crop");
  });
}

function showTableConstraints() {
  var name = getCurrentObject().name;

  if (name.length == 0) {
    alert("Please select a table!");
    return;
  }

  getTableConstraints(name, function(data) {
    setCurrentTab("table_constraints");
    buildTable(data);

    $("#input").hide();
    $("#body").prop("class", "full");
    $("#results").addClass("no-crop");
  });
}

function showTableInfo() {
  var name = getCurrentObject().name;

  if (name.length == 0) {
    alert("Please select a table!");
    return;
  }

  apiCall("get", "/tables/" + name + "/info", {}, function(data) {
    $(".table-information .lines").show();
    $("#table_total_size").text(data.total_size);
    $("#table_data_size").text(data.data_size);
    $("#table_index_size").text(data.index_size);
    $("#table_rows_count").text(data.rows_count);
    $("#table_encoding").text("Unknown");
  });

  buildTableFilters(name, getCurrentObject().type);
}

function updatePaginator(pagination) {
  if (!pagination) {
    $(".current-page").data("page", 1).data("pages", 1);
    $("button.page").text("1 of 1");
    $(".prev-page, .next-page").prop("disabled", "disabled");
    return;
  }

  $(".current-page").
    data("page", pagination.page).
    data("pages", pagination.pages_count);

  if (pagination.page > 1) {
    $(".prev-page").prop("disabled", "");
  }
  else {
    $(".prev-page").prop("disabled", "disabled");
  }

  if (pagination.pages_count > 1 && pagination.page < pagination.pages_count) {
    $(".next-page").prop("disabled", "");
  }
  else {
    $(".next-page").prop("disabled", "disabled");
  }

  $("#total_records").text(pagination.rows_count);
  if (pagination.pages_count == 0) pagination.pages_count = 1;
  $("button.page").text(pagination.page + " of " + pagination.pages_count);
}

function showTableContent(sortColumn, sortOrder) {
  var name = getCurrentObject().name;

  if (name.length == 0) {
    alert("Please select a table!");
    return;
  }

  if (getCurrentObject().type == "function") {
    alert("Cant view rows for a function");
    return;
  }

  var opts = {
    limit:       getRowsLimit(),
    offset:      getPaginationOffset(),
    sort_column: sortColumn,
    sort_order:  sortOrder
  };

  // Structured filters are pushed down to PostgreSQL as parameterized conditions
  if (tableFilters.length > 0) {
    opts["filters"] = JSON.stringify(tableFilters);
  }

  // Metadata is cached, so this does not add a request per page
  loadTableMetadata(name, function() {
    getTableRows(name, opts, function(data) {
      $("#input").hide();
      $("#body").prop("class", "with-pagination");

      buildTable(data, sortColumn, sortOrder, { pkColumns: editablePrimaryKey() });
      setCurrentTab("table_content");
      updatePaginator(data.pagination);

      $("#results").data("mode", "browse").data("table", name);
      updateRowActions();
    });
  });
}

// editablePrimaryKey returns the primary key columns of the current table when
// the table can actually be modified, and an empty list otherwise.
function editablePrimaryKey() {
  if (!currentMetadata) return [];
  if (!currentMetadata.editable) return [];
  if (!currentMetadata.primary_key || currentMetadata.primary_key.length == 0) return [];

  return currentMetadata.primary_key;
}

// rowEditingAvailable reports whether the current grid supports row level writes
function rowEditingAvailable() {
  return $("#results").data("mode") == "browse" &&
    $("#table_content").hasClass("selected") &&
    editablePrimaryKey().length > 0;
}

function updateRowActions() {
  var insertable = currentMetadata && currentMetadata.insertable;

  if (!insertable && !rowEditingAvailable()) {
    $("#pagination .row-actions").hide();
    return;
  }

  $("#pagination .row-actions").show();
  $("#pagination .new-row").prop("disabled", !insertable);
  $("#pagination .edit-row, #pagination .delete-row").prop("disabled", !rowEditingAvailable());
}

function showPaginatedTableContent() {
  var activeColumn = $("#results th.active");
  var sortColumn = null;
  var sortOrder = null;

  if (activeColumn.length) {
    sortColumn = activeColumn.data("name");
    sortOrder = activeColumn.data("order");
  }

  showTableContent(sortColumn, sortOrder);
}

function showDatabaseStats() {
  getTablesStats(function(data) {
    buildTable(data);

    setCurrentTab("table_structure");
    $("#input").hide();
    $("#body").prop("class", "full");
    $("#results").addClass("no-crop");
  });
}

function downloadDatabaseStats() {
  openInNewWindow("api/tables_stats", { format: "csv", export: "true" });
}

function showServerSettings() {
  getServerSettings(function(data) {
    buildTable(data);

    setCurrentTab("table_content");
    $("#input").hide();
    $("#body").prop("class", "full");
    $("#results").addClass("no-crop");
  });
}

function showTableStructure() {
  var name = getCurrentObject().name;

  if (name.length == 0) {
    alert("Please select a table!");
    return;
  }

  setCurrentTab("table_structure");

  $("#input").hide();
  $("#body").prop("class", "full");

  getTableStructure(name, { type: getCurrentObject().type }, function(data) {
    if (getCurrentObject().type == "function") {
      var name = data.rows[0][data.columns.indexOf("proname")];
      var definition = data.rows[0][data.columns.indexOf("functiondef")];
      showFunctionDefinition(name, definition);
      return
    }

    buildTable(data);
    $("#results").addClass("no-crop");
  });
}

function showViewDefinition(viewName, viewDefintion) {
  setCurrentTab("table_structure");
  renderResultsView("View definition for: <strong>" + viewName + "</strong>", viewDefintion);
}

function showFunctionDefinition(functionName, definition) {
  setCurrentTab("table_structure");
  renderResultsView("Function definition for: <strong>" + functionName + "</strong>", definition)
}

function renderResultsView(title, content) {
  $("#results").addClass("no-crop");
  $("#input").hide();
  $("#body").prop("class", "full");
  $("#results").hide();

  var title = $("<div/>").prop("class", "title").html(title);
  var content = $("<pre/>").text(content);

  $("<div/>").
    html("<i class='fa fa-copy'></i>").
    addClass("copy").
    appendTo(content);

  $("#results_view").html("");
  title.appendTo("#results_view");
  content.appendTo("#results_view");
  $("#results_view").show();
}

function showQueryPanel() {
  if (!$("#table_query").hasClass("selected")) {
    resetTable();
  }

  setCurrentTab("table_query");
  editor.focus();

  $("#input").show();
  $("#body").prop("class", "")
}

function showConnectionPanel() {
  setCurrentTab("table_connection");
  $("#input").hide();
  $("#body").addClass("full");

  getConnection(function(data) {
    var rows = [];

    for(key in data) {
      rows.push([key, data[key]]);
    }

    buildTable({
      columns: ["attribute", "value"],
      rows: rows
    });
  });
}

function showActivityPanel() {
  var options = {
    action: {
      name: "stop_query",
      title: "stop",
      data: "pid",
      style: "danger"
    }
  }

  setCurrentTab("table_activity");
  $("#input").hide();
  $("#body").addClass("full");

  apiCall("get", "/activity", {}, function(data) {
    buildTable(data, null, null, options);
  });
}

function showQueryProgressMessage() {
  $("#run, #explain-dropdown-toggle, #csv, #json, #xml, #load-local-query").prop("disabled", true);
  $("#explain-dropdown").removeClass("open");
  $("#query_progress").show();
}

function hideQueryProgressMessage() {
  $("#run, #explain-dropdown-toggle, #csv, #json, #xml, #load-local-query").prop("disabled", false);
  $("#query_progress").hide();
}

function getEditorSelection() {
  // Return the exact selection if user has one
  var query = $.trim(editor.getSelectedText());
  if (query.length > 0) {
    return query;
  }

  query = editor.getValue();

  // Determine which query we should run when there are multiple queries without a delimiter
  if (query.indexOf(";") == -1) {
    var subquery = getSubquery(query, editor.getCursorPosition());

    if (subquery) {
      // Highlight query selection so user knows what is being executed
      if (subquery.numChunks > 1) {
        editor.selection.setSelectionRange({
          start: { row: subquery.startRow, column: 0 },
          end: { row: subquery.endRow, column: 0 },
        })
      }

      return subquery.text;
    }
  }

  return query;
}

function getSubquery(text, cursor) {
  var lines = text.split("\n");
  var startRow = undefined;
  var numChunks = 0;
  var ranges = [];

  for (i = 0; i < lines.length; i++) {
    if (lines[i].trim().length == 0) {
      if (startRow >= 0 && cursor.row >= startRow && cursor.row <= i) {
        ranges.push([startRow, i]);
      }

      numChunks++;
      startRow = undefined;
      continue;
    }

    if (startRow === undefined) {
      startRow = i;
    }

    if (i == lines.length - 1) {
      ranges.push([startRow, i + 1]);
      numChunks++;
    }
  }

  if (ranges.length > 0) {
    return {
      text: lines.slice(ranges[0][0], ranges[0][1]).join("\n"),
      startRow: ranges[0][0],
      endRow: ranges[0][1],
      numChunks: numChunks
    };
  }
}

function runQuery() {
  setCurrentTab("table_query");
  showQueryProgressMessage();

  var query = getEditorSelection();
  if (query.length == 0) {
    hideQueryProgressMessage();
    return;
  }

  executeQuery(query, function(data) {
    buildTable(data);

    hideQueryProgressMessage();
    $("#input").show();
    $("#body").removeClass("full");
    $("#results").data("mode", "query");

    if (query.toLowerCase().indexOf("explain") != -1) {
      $("#results").addClass("no-crop");
    }

    // Reload objects list if anything was created/deleted
    if (query.match(/(create|drop)\s/i)) {
      loadSchemas();
    }
  });
}

function runExplain() {
  setCurrentTab("table_query");
  showQueryProgressMessage();

  var query = getEditorSelection();
  if (query.length == 0) {
    hideQueryProgressMessage();
    return;
  }

  explainQuery(query, function(data) {
    buildTable(data);

    hideQueryProgressMessage();
    $("#input").show();
    $("#body").removeClass("full");
    $("#results").addClass("no-crop");
  });
}

function runAnalyze() {
  setCurrentTab("table_query");
  showQueryProgressMessage();

  var query = getEditorSelection();
  if (query.length == 0) {
    hideQueryProgressMessage();
    return;
  }

  analyzeQuery(query, function(data) {
    buildTable(data);

    hideQueryProgressMessage();
    $("#input").show();
    $("#body").removeClass("full");
    $("#results").addClass("no-crop");
  });
}

function generateURL(path, params) {
  var url = new URL(window.location.href.split("#")[0]);

  url.pathname += path;
  for (key in params) {
    url.searchParams.append(key, params[key]);
  }

  // Automatically append session id so we dont have to do that everywhere
  url.searchParams.append("_session_id", getSessionId());

  return url.toString();
}

function openInNewWindow(path, params) {
  var url = generateURL(path, params);
  var win = window.open(url, '_blank');
  win.focus();
}

function exportTo(format) {
  var query = getEditorSelection();
  if (query.length == 0) {
    return;
  }

  setCurrentTab("table_query");

  openInNewWindow("api/query", {
    "format": format,
    "query": encodeQuery(query)
  })
}

// Fetch all unique values for the selected column in the table
function showUniqueColumnsValues(table, column, showCounts) {
  var query = 'SELECT DISTINCT "' + column + '" FROM ' + table;

  // Display results ordered by counts.
  // This could be slow on large sets without an index.
  if (showCounts) {
    query = 'SELECT DISTINCT "' + column + '", COUNT(1) AS total_count FROM ' + table + ' GROUP BY "' + column + '" ORDER BY total_count DESC';
  }

  executeQuery(query, function(data) {
    $("#input").hide();
    $("#body").prop("class", "full");
    $("#results").data("mode", "query");
    buildTable(data);
  });
}

// Show numeric stats on the field
function showFieldNumStats(table, column) {
  var query = 'SELECT count(1), min(' + column + '), max(' + column + '), avg(' + column + ') FROM ' + table;

  executeQuery(query, function(data) {
    $("#input").hide();
    $("#body").prop("class", "full");
    $("#results").data("mode", "query");
    buildTable(data);
  });
}

function buildTableFilters(name, type) {
  tableFilters = [];
  renderActiveFilters();

  loadTableMetadata(name, function(meta) {
    if (meta) {
      renderFilterColumns(meta.columns.map(function(col) { return col.name; }));
      updateRowActions();
      return;
    }

    // Objects without column metadata (sequences, custom types) fall back to
    // the structure endpoint so that filtering keeps working as before.
    getTableStructure(name, { type: type }, function(data) {
      var columns = [];

      if (data && data.rows) {
        columns = data.rows.map(function(row) { return row[0]; });
      }

      renderFilterColumns(columns);
      updateRowActions();
    });
  });
}

// loadTableMetadata fetches and caches the column metadata for a table. The
// cache keeps insert/edit forms and filters from re-querying the schema.
function loadTableMetadata(name, cb) {
  if (metadataCache[name] !== undefined) {
    currentMetadata = metadataCache[name];
    return cb(currentMetadata);
  }

  getTableMetadata(name, function(data) {
    var meta = null;

    if (data && !data.error && data.columns && data.columns.length > 0) {
      meta = data;
    }

    metadataCache[name] = meta;
    currentMetadata = meta;
    cb(meta);
  });
}

function resetMetadataCache() {
  metadataCache = {};
  currentMetadata = null;
}

function resetTableFilters() {
  tableFilters = [];
  $(".filters select, .filters input").val("");
  $(".filters input").show();
  $(".current-page").data("page", 1);
  renderActiveFilters();
}

function renderFilterColumns(columns) {
  if (columns.length == 0) {
    $("#pagination .filters").hide();
  }
  else {
    $("#pagination .filters").show();
  }

  var select = $("#pagination select.column");
  select.html("<option value='' selected>Select column</option>");

  columns.forEach(function(col) {
    select.append($("<option/>").attr("value", col).text(col));
  });
}

// renderActiveFilters draws the list of currently applied conditions above the
// results grid. The filters array is the single source of truth.
function renderActiveFilters() {
  var container = $("#active_filters").empty();

  if (tableFilters.length == 0) {
    container.hide();
    return;
  }

  tableFilters.forEach(function(filter, idx) {
    var chip = $("<span/>").addClass("filter-chip");

    chip.append($("<span/>").addClass("column").text(filter.column));
    chip.append($("<span/>").addClass("op").text(filterLabels[filter.op] || filter.op));

    if (filterOpHasValue(filter.op)) {
      chip.append($("<span/>").addClass("value").text(filter.value));
    }

    chip.append(
      $("<i/>").addClass("fa fa-times remove-filter").attr("data-index", idx).attr("title", "Remove this filter")
    );

    container.append(chip);
  });

  container.append($("<a/>").addClass("clear-filters").attr("href", "#").text("Clear all"));
  container.show();
}

// addFilterFromForm moves the condition built in the toolbar into the active
// filter list. Returns false when the form is incomplete.
function addFilterFromForm() {
  var column = $(".filters select.column").val();
  var op     = $(".filters select.filter").val();
  var value  = $.trim($(".filters input").val());

  if (!column || !op) return false;

  if (filterOpHasValue(op) && value == "") {
    alert("Please specify filter query");
    return false;
  }

  tableFilters.push({ column: column, op: op, value: filterOpHasValue(op) ? value : "" });

  $(".filters select, .filters input").val("");
  $(".filters input").show();
  renderActiveFilters();

  return true;
}

var objectAutocompleter = {
  getCompletions: function (editor, session, pos, prefix, callback) {
    callback(null, autocompleteObjects);
  }
}

function initEditor() {
  var writeQueryTimeout = null;

  editor = ace.edit("custom_query");
  editor.setOptions({
    enableBasicAutocompletion: true,
    enableLiveAutocompletion: true,
  });
  editor.completers.push(objectAutocompleter);

  editor.setFontSize(13);
  editor.setTheme("ace/theme/tomorrow");
  editor.setShowPrintMargin(false);
  editor.getSession().setMode("ace/mode/pgsql");
  editor.getSession().setTabSize(2);
  editor.getSession().setUseSoftTabs(true);

  editor.commands.addCommands([{
    name: "run_query",
    bindKey: {
      win: "Ctrl-Enter",
      mac: "Command-Enter"
    },
    exec: function(editor) {
      runQuery();
    }
  }, {
    name: "explain_query",
    bindKey: {
      win: "Ctrl-E",
      mac: "Command-E"
    },
    exec: function(editor) {
      runExplain();
    }
  }]);

  editor.on("change", function() {
    if (writeQueryTimeout) {
      clearTimeout(writeQueryTimeout);
    }

    writeQueryTimeout = setTimeout(function() {
      localStorage.setItem("pgweb_query", editor.getValue());
    }, 1000);
  });

  var query = localStorage.getItem("pgweb_query");
  if (query && query.length > 0) {
    editor.setValue(query);
    editor.clearSelection();
  }
}

function addShortcutTooltips() {
  if (navigator.userAgent.indexOf("OS X") > 0) {
    $("#run").attr("title", "Shortcut: ⌘+Enter");
    $("#explain").attr("title", "Shortcut: ⌘+E");
  }
  else {
    $("#run").attr("title", "Shortcut: Ctrl+Enter");
    $("#explain").attr("title", "Shortcut: Ctrl+E");
  }
}

// Get the latest release from Github API
function getLatestReleaseInfo(current) {
  try {
    $.get("https://api.github.com/repos/sosedoff/pgweb/releases/latest", function(release) {
      if (release.name != current.version) {
        var message = "Update available. Check out " + release.tag_name + " on <a target='_blank' href='" + release.html_url + "'>Github</a>";
        $(".connection-settings .update").html(message).fadeIn();
      }
    });
  }
  catch(error) {
    console.log("Cant get last release from github:", error);
  }
}

function showConnectionSettings() {
  // Show the current postgres version
  $(".connection-settings .version").text("v" + appInfo.version).show();
  $("#connection_window").show();
  initConnectionWindow();

  // Check github release page for updates
  getLatestReleaseInfo(appInfo);

  getBookmarks(function(data) {
    if (data.error) {
      console.log("Error while fetching bookmarks:", data.error);
      return;
    }

    if (data.length > 0) {
      // Set bookmarks in global var
      bookmarks = data;

      // Remove all existing bookmark options
      $("#connection_bookmarks").html("");

      // Add blank option
      $("<option value=''>Select a bookmarked database to connect to</option>").appendTo("#connection_bookmarks");

      // Add all available bookmarks
      for (key of data) {
        $("<option value='" + key + "''>" + key + "</option>").appendTo("#connection_bookmarks");
      }

      $(".bookmarks").show();
    }
    else {
      if (appFeatures.bookmarks_only) {
        $("#connection_error").html("Running in <b>bookmarks-only</b> mode but <b>NO</b> bookmarks configured.").show();
        $(".open-connection").hide();
      } else {
        $(".bookmarks").hide();
      }
    }
  });
}

function initConnectionWindow() {
  if (appFeatures.bookmarks_only) {
    $(".connection-group-switch").hide();
    $(".connection-scheme-group").hide();
    $(".connection-bookmarks-group").show();
    $(".connection-standard-group").hide();
    $(".connection-ssh-group").hide();
  } else {
    $(".connection-group-switch").show();
    $(".connection-scheme-group").hide();
    $(".connection-bookmarks-group").show();
    $(".connection-standard-group").show();
    $(".connection-ssh-group").hide();
  }
}

function getConnectionString() {
  var url  = $.trim($("#connection_url").val());
  var mode = $(".connection-group-switch button.active").attr("data");
  var ssl  = $("#connection_ssl").val();

  if (mode == "standard" || mode == "ssh") {
    var host = $("#pg_host").val();
    var port = $("#pg_port").val();
    var user = $("#pg_user").val();
    var pass = encodeURIComponent($("#pg_password").val());
    var db   = $("#pg_db").val();

    if (port.length == 0) {
      port = "5432";
    }

    url = "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + db + "?sslmode=" + ssl;
  }
  else {
    var local = url.indexOf("localhost") != -1 || url.indexOf("127.0.0.1") != -1;

    if (local && url.indexOf("sslmode") == -1) {
      url += "?sslmode=" + ssl;
    }
  }

  return url;
}

// Add a context menu to the results table header columns
function bindTableHeaderMenu() {
  $("#results_header").contextmenu({
    scopes: "th",
    target: "#results_header_menu",
    before: function(e, element, target) {
      // Enable menu for browsing table rows view only.
      if ($("#results").data("mode") != "browse") {
        e.preventDefault();
        this.closemenu();
        return false;
      }
    },
    onItem: function(context, e) {
      var menuItem = $(e.target);

      switch(menuItem.data("action")) {
        case "copy_name":
          copyToClipboard($(context).data("name"));
          break;

        case "unique_values":
          showUniqueColumnsValues(
            $("#results").data("table"), // table name
            $(context).data("name"),     // column name
            menuItem.data("counts")      // display counts
          );
          break;

        case "num_stats":
          showFieldNumStats(
            $("#results").data("table"), // table name
            $(context).data("name")      // column name
          );
          break;
      }
    }
  });

  $("#results_body").contextmenu({
    scopes: "td",
    target: "#results_row_menu",
    before: function(e, element, target) {
      var browseMode = $("#results").data("mode");
      var isEmpty    = $("#results").hasClass("empty");
      var isAllowed  = browseMode == "browse" || browseMode == "query";

      if (isEmpty || !isAllowed) {
        e.preventDefault();
        this.closemenu();
        return false;
      }
    },
    onItem: function(context, e) {
      var menuItem = $(e.target);

      switch(menuItem.data("action")) {
        case "display_value":
          var value = $(context).text();
          $("#content_modal .content").text(value);
          $("#content_modal").show();
          break;
        case "copy_value":
          copyToClipboard($(context).text());
          break;
        case "filter_by_value":
          var colIdx   = $(context).data("col");
          var colValue = $(context).text();
          var colName  = $("#results_header th").eq(colIdx).data("name");

          $("select.column").val(colName);
          $("select.filter").val("equal");
          $("#table_filter_value").val(colValue);
          $("#rows_filter").submit();
      }
    }
  });
}

function bindCurrentDatabaseMenu() {
  $("#current_database").contextmenu({
    target: "#current_database_context_menu",
    onItem: function(context, e) {
      var menuItem = $(e.target);

      switch(menuItem.data("action")) {
        case "show_db_stats":
          showDatabaseStats();
          break;
        case "download_db_stats":
          downloadDatabaseStats();
          break;
        case "server_settings":
          showServerSettings();
          break;
        case "export":
          openInNewWindow("api/export");
          break;
      }
    }
  });
}

function bindDatabaseObjectsFilter() {
  var filterTimeout = null;

  $("#filter_database_objects").on("keyup", function (e) {
    clearTimeout(filterTimeout);

    var val = $(this).val().trim();

    // Reset search on ESC
    if (e.keyCode == 27 || val == "") {
      resetObjectsFilter();
      return;
    }

    $(".clear-objects-filter").show();
    $(".schema-group").addClass("expanded");

    filterTimeout = setTimeout(function() {
      filterObjectsByName(val)
    }, 200);
  });

  $(".clear-objects-filter").on("click", function(e) {
    resetObjectsFilter();
  });
}

function resetObjectsFilter() {
  $("#filter_database_objects").val("");
  $("#objects li.schema-item").show();
  $(".clear-objects-filter").hide();
}

function filterObjectsByName(query) {
  $("#objects li.schema-item").each(function (idx, el) {
    var item = $(el);
    var name = $(el).data("name");

    if (name.indexOf(query) < 0) {
      item.hide();
    } else {
      item.show();
    }
  });
}

function getQuotedSchemaTableName(table) {
  if (typeof table === "string" && table.indexOf(".") > -1) {
    var schemaTableComponents = table.split(".");
    return ['"', schemaTableComponents[0], '"."', schemaTableComponents[1], '"'].join('');
  }
  return table;
}

function bindContextMenus() {
  bindTableHeaderMenu();
  bindCurrentDatabaseMenu();

  $(".schema-group ul").each(function(id, el) {
    var group = $(el).data("group");

    if (group == "table") {
      $(el).contextmenu({
        target: "#tables_context_menu",
        scopes: "li.schema-table",
        onItem: function(context, e) {
          var el      = $(e.target);
          var table   = getQuotedSchemaTableName($(context[0]).data("id"));
          var action  = el.data("action");
          performTableAction(table, action, el);
        }
      });
    }

    if (group == "view") {
      $(el).contextmenu({
        target: "#view_context_menu",
        scopes: "li.schema-view",
        onItem: function(context, e) {
          var el      = $(e.target);
          var table   = getQuotedSchemaTableName($(context[0]).data("id"));
          var action  = el.data("action");
          performViewAction(table, action, el);
        }
      });
    }

    if (group == "materialized_view") {
      $(el).contextmenu({
        target: "#view_context_menu",
        scopes: "li.schema-materialized_view",
        onItem: function(context, e) {
          var el      = $(e.target);
          var table   = getQuotedSchemaTableName($(context[0]).data("id"));
          var action  = el.data("action");
          performViewAction(table, action, el);
        }
      });
    }
  });
}

function toggleDatabaseSearch() {
  $("#current_database").toggle();
  $("#database_search").toggle();
}

function enableDatabaseSearch(data) {
  var input = $("#database_search");

  input.typeahead("destroy");

  input.typeahead({
    source: data,
    minLength: 0,
    items: "all",
    autoSelect: false,
    fitToElement: true
  });

  input.typeahead("lookup").focus();

  input.on("focusout", function(e){
    toggleDatabaseSearch();
    input.off("focusout");
  });
}

function bindInputResizeEvents() {
  var height = sessionStorage.getItem("input_height");
  if (height) {
    resizeInput(height);
    checkInputSize();
  }

  $("body").on("mousemove", onInputResize);
  $("body").on("mouseup", endInputResize);
  $("#input_resize_handler").on("mousedown", beginInputResize);
  $(window).on("resize", checkInputSize);
}

function checkInputSize() {
  var inputHeight = $("#input").height();
  var bodyHeight = $("#body").height();

  if (bodyHeight == 0 || inputHeight == 0) return;

  if (inputHeight > bodyHeight || bodyHeight - inputHeight < 200) {
    resizeInput(bodyHeight - 200);
  }
}

function resizeInput(height) {
  if (height < 100) height = 100;

  var diff = 50 + 12; // actions box + padding

  $("#input").height(height);
  $("#input .input-wrapper").height(height - diff);
  $("#custom_query").height(height - diff);
  $("#output").css("top", height + "px");

  if (editor) {
    editor.resize();
  }
}

function beginInputResize() {
  inputResizing = true;
  inputResizeOffset = $("#input").offset().top;

  $("html").css("cursor", "row-resize");
  $("#input_resize_handler").addClass("dragging");
}

function endInputResize() {
  if (!inputResizing) return;

  inputResizing = false;
  inputResizeOffset = null;

  $("html").css("cursor", "auto");
  $("#input_resize_handler").removeClass("dragging");

  // Save current settings for page reloads
  sessionStorage.setItem("input_height", $("#input").height());
}

function onInputResize(event) {
  if (!inputResizing) return;

  var computedHeight = event.clientY - inputResizeOffset;
  if (computedHeight < 150) computedHeight = 150;

  resizeInput(computedHeight);
}

// ----------------------------------------------------------------------------
// Row editing: inline cell editor and the insert/edit/delete row modal.
// All write operations go through the parameterized /tables/:table/rows API.
// ----------------------------------------------------------------------------

function metaColumn(name) {
  if (!currentMetadata) return null;

  var found = null;

  currentMetadata.columns.forEach(function(col) {
    if (col.name == name) found = col;
  });

  return found;
}

function columnWritable(col) {
  return !!col && col.updatable && !col.generated;
}

// Binary columns are not editable as text in v1.0
function columnEditableType(col) {
  return !!col && (col.udt_name || "").toLowerCase() != "bytea";
}

function cellIsNull(td) {
  return td.find("> div > span.null").length > 0;
}

function cellColumnName(td) {
  var columns = $("#results").data("columns") || [];
  var idx = parseInt(td.data("col"));

  if (isNaN(idx)) return null;
  return columns[idx];
}

function rowValuesFromGrid(tr) {
  var values = {};

  tr.find("> td").each(function() {
    var name = cellColumnName($(this));
    if (name == null) return;

    values[name] = cellIsNull($(this)) ? null : $(this).find("> div").text();
  });

  return values;
}

// patchRow updates a single grid row in place, avoiding a full table reload
function patchRow(tr, row) {
  if (!tr || !row) return;

  tr.find("> td").each(function() {
    var name = cellColumnName($(this));
    if (name == null || !(name in row)) return;

    $(this).find("> div").html(escapeHtml(row[name]));
  });

  // The primary key itself may have been updated
  var pkColumns = editablePrimaryKey();

  if (pkColumns.length > 0) {
    var key = {};
    pkColumns.forEach(function(col) { key[col] = row[col]; });
    tr.data("pk", key);
  }
}

function startCellEdit(td) {
  if (!rowEditingAvailable()) return false;
  if (td.hasClass("cell-editing")) return true;

  var name = cellColumnName(td);
  var col  = metaColumn(name);

  if (!columnWritable(col) || !columnEditableType(col)) return false;

  var tr  = td.closest("tr");
  var key = tr.data("pk");
  if (!key) return false;

  var div   = td.find("> div");
  var value = cellIsNull(td) ? "" : div.text();
  var input = $("<input/>").attr("type", "text").addClass("cell-editor").val(value);

  input.data("original", value).data("column", name).data("key", key).data("cell", td);

  td.addClass("cell-editing");
  div.hide();
  td.append(input);
  input.focus().select();

  return true;
}

function finishCellEdit(input, save) {
  var td  = input.data("cell");
  var div = td.find("> div");

  function restore() {
    input.remove();
    td.removeClass("cell-editing");
    div.show();
  }

  var column   = input.data("column");
  var original = input.data("original");
  var value    = input.val();

  if (!save || value === original) {
    restore();
    return;
  }

  var col     = metaColumn(column);
  var payload = {};

  // An empty input clears a nullable column. Use the edit row modal to set an
  // empty string on a nullable text column.
  payload[column] = (value === "" && col && col.nullable) ? null : value;

  input.prop("disabled", true);

  updateTableRow($("#results").data("table"), input.data("key"), payload, function(data) {
    if (data.error) {
      input.prop("disabled", false).focus();
      showErrorBanner(data.error);
      return;
    }

    restore();
    patchRow(td.closest("tr"), data.row);
    td.addClass("cell-modified");
    showSuccessBanner("Row updated");
  });
}

function buildRowForm(mode, values) {
  var body = $("#row_modal .row-modal-body").empty();

  currentMetadata.columns.forEach(function(col) {
    var writable = columnWritable(col) && columnEditableType(col);

    // Read-only columns are not worth showing on a new row form
    if (mode == "insert" && !writable) return;

    var value  = values ? values[col.name] : null;
    var isNull = (value === null || value === undefined);

    var field = $("<div/>").addClass("field");
    var label = $("<label/>").text(col.name);

    label.append($("<span/>").addClass("type").text(col.data_type));
    if (col.primary_key) label.append($("<span/>").addClass("pk").text("PK"));
    field.append(label);

    var input;

    if ((col.udt_name || "").toLowerCase() == "bool") {
      input = $("<select/>");
      // The blank option lets the column default (or the current value) stand
      input.append($("<option/>").attr("value", "").text(col.default ? "(default)" : "(unset)"));
      input.append($("<option/>").attr("value", "true").text("true"));
      input.append($("<option/>").attr("value", "false").text("false"));
      input.val(isNull ? "" : (String(value) === "true" ? "true" : "false"));
    }
    else {
      input = $("<input/>").attr("type", "text").val(isNull ? "" : String(value));
      if (col.default) input.attr("placeholder", "default: " + col.default);
    }

    input.addClass("row-field").attr("data-column", col.name);
    if (!writable) input.prop("disabled", true);
    field.append(input);

    if (col.nullable && writable) {
      var checkbox = $("<input/>").attr("type", "checkbox").addClass("null-checkbox");
      if (isNull && mode == "edit") checkbox.prop("checked", true);

      field.append(
        $("<label/>").addClass("null-toggle").append(checkbox).append(document.createTextNode(" NULL"))
      );
    }

    body.append(field);
  });
}

function buildDeleteSummary(values, key) {
  var body = $("#row_modal .row-modal-body").empty();

  body.append($("<p/>").text("This will permanently delete the row below. The action cannot be undone."));

  var summary = $("<div/>").addClass("summary");
  var shown   = 0;

  Object.keys(key).forEach(function(name) {
    summary.append(
      $("<div/>").append($("<b/>").text(name + ": ")).append(document.createTextNode(String(key[name])))
    );
  });

  Object.keys(values).forEach(function(name) {
    if (name in key || shown >= 8) return;
    shown++;

    var value = values[name];

    summary.append(
      $("<div/>").append($("<b/>").text(name + ": ")).append(document.createTextNode(value === null ? "null" : String(value)))
    );
  });

  body.append(summary);
}

function collectRowValues(mode, original) {
  var values = {};

  $("#row_modal .field").each(function() {
    var input = $(this).find(".row-field");
    if (input.prop("disabled")) return;

    var column = input.attr("data-column");
    var isNull = $(this).find(".null-checkbox").prop("checked") === true;
    var value  = isNull ? null : input.val();

    // An unset dropdown means "leave this column alone"
    if (!isNull && value === "" && input.is("select")) return;

    if (mode == "insert") {
      // Empty fields are omitted so that column defaults are applied
      if (!isNull && value === "") return;

      values[column] = value;
      return;
    }

    var before = original ? original[column] : undefined;
    if (before === undefined) before = null;
    if (before !== null) before = String(before);

    // Only the modified fields are sent
    if (value === before) return;

    values[column] = value;
  });

  return values;
}

function showRowModal(mode) {
  if (!currentMetadata) return;

  var tr = null, original = null, key = null;

  if (mode != "insert") {
    tr = $("#results_body > tr.selected");

    if (tr.length != 1) {
      alert("Please select a row first");
      return;
    }

    key = tr.data("pk");

    if (!key) {
      alert("No primary key detected for this table, rows cannot be modified.");
      return;
    }

    original = rowValuesFromGrid(tr);
  }

  rowModalContext = { mode: mode, tr: tr, original: original, key: key };

  var submit = $("#row_modal .row-modal-action[data-action=submit]");
  var title  = currentMetadata.schema + "." + currentMetadata.table;

  $("#row_modal .row-modal-error").text("");

  if (mode == "delete") {
    $("#row_modal .title").text("Delete row from " + title);
    buildDeleteSummary(original, key);
    submit.text("Delete").removeClass("btn-primary").addClass("btn-danger");
  }
  else {
    $("#row_modal .title").text((mode == "insert" ? "New row in " : "Edit row in ") + title);
    buildRowForm(mode, original);
    submit.text("Save").removeClass("btn-danger").addClass("btn-primary");
  }

  submit.prop("disabled", false);
  $("#row_modal").show();
}

function closeRowModal() {
  $("#row_modal").hide();
  $("#row_modal .row-modal-body").empty();
  rowModalContext = null;
}

function submitRowModal() {
  if (!rowModalContext) return;

  var ctx    = rowModalContext;
  var table  = $("#results").data("table") || getCurrentObject().name;
  var submit = $("#row_modal .row-modal-action[data-action=submit]").prop("disabled", true);

  function done(data) {
    submit.prop("disabled", false);

    if (data.error) {
      $("#row_modal .row-modal-error").text(data.error);
      return;
    }

    closeRowModal();

    switch (ctx.mode) {
      case "insert":
        showSuccessBanner("Row created");
        showPaginatedTableContent();
        break;
      case "edit":
        patchRow(ctx.tr, data.row);
        ctx.tr.find("> td").addClass("cell-modified");
        showSuccessBanner("Row updated");
        break;
      case "delete":
        ctx.tr.remove();
        showSuccessBanner("Row deleted");
        break;
    }
  }

  if (ctx.mode == "delete") {
    deleteTableRow(table, ctx.key, done);
    return;
  }

  var values = collectRowValues(ctx.mode, ctx.original);

  if (Object.keys(values).length == 0) {
    submit.prop("disabled", false);
    $("#row_modal .row-modal-error").text(ctx.mode == "insert" ? "Please fill in at least one column" : "No changes to save");
    return;
  }

  if (ctx.mode == "insert") {
    createTableRow(table, values, done);
  }
  else {
    updateTableRow(table, ctx.key, values, done);
  }
}

function bindRowEditingEvents() {
  $("#pagination").on("click", ".new-row",    function() { showRowModal("insert"); });
  $("#pagination").on("click", ".edit-row",   function() { showRowModal("edit");   });
  $("#pagination").on("click", ".delete-row", function() { showRowModal("delete"); });

  $("#row_modal").on("click", ".row-modal-action", function() {
    switch ($(this).data("action")) {
      case "submit":
        submitRowModal();
        break;
      case "close":
        closeRowModal();
        break;
    }
  });

  $("#row_modal").on("click", function(e) {
    // Clicking the backdrop closes the modal
    if (e.target === this) closeRowModal();
  });

  $("#row_modal").on("change", ".null-checkbox", function() {
    $(this).closest(".field").find(".row-field").prop("readonly", $(this).prop("checked"));
  });

  $("#results").on("keydown", "input.cell-editor", function(e) {
    switch (e.keyCode) {
      case 13: // enter
        e.preventDefault();
        finishCellEdit($(this), true);
        break;
      case 27: // escape
        e.preventDefault();
        finishCellEdit($(this), false);
        break;
    }
  });

  $("#results").on("blur", "input.cell-editor", function() {
    // Cancel on focus loss to avoid accidental writes
    if (!$(this).prop("disabled")) finishCellEdit($(this), false);
  });

  $("#results").on("click", "input.cell-editor", function(e) {
    e.stopPropagation();
  });
}

function bindContentModalEvents() {
  var contentModal = document.getElementById("content_modal");

  $(window).on("click", function(e) {
    // Automatically hide the modal on any click outside of the modal window
    if (e.target && !contentModal.contains(e.target)) {
      $("#content_modal").hide();
    }
  });

  $("#content_modal .content-modal-action").on("click", function() {
    switch ($(this).data("action")) {
      case "copy":
        copyToClipboard($("#content_modal pre").text());
        break;
      case "close":
        $("#content_modal").hide();
        break;
    }
  });

  $("#results").on("dblclick", "td > div", function() {
    // Editable cells enter inline edit mode, everything else opens the viewer
    if (startCellEdit($(this).parent())) return;

    var value = unescapeHtml($(this).html());
    if (!value) return;

    $("#content_modal pre").html(value);
    $("#content_modal").show();
  })
}

$(document).ready(function() {
  bindInputResizeEvents();
  bindContentModalEvents();
  bindRowEditingEvents();

  $("#table_content").on("click",     function() { showTableContent();     });
  $("#table_structure").on("click",   function() { showTableStructure();   });
  $("#table_indexes").on("click",     function() { showTableIndexes();     });
  $("#table_constraints").on("click", function() { showTableConstraints(); });
  $("#table_history").on("click",     function() { showQueryHistory();     });
  $("#table_query").on("click",       function() { showQueryPanel();       });
  $("#table_connection").on("click",  function() { showConnectionPanel();  });
  $("#table_activity").on("click",    function() { showActivityPanel();    });

  $("#run").on("click", function() {
    runQuery();
  });

  $("#explain").on("click", function() {
    runExplain();
  });

  $("#analyze").on("click", function() {
    runAnalyze();
  });

  $("#csv").on("click", function() {
    exportTo("csv");
  });

  $("#json").on("click", function() {
    exportTo("json");
  });

  $("#xml").on("click", function() {
    exportTo("xml");
  });

  $("#results_view").on("click", ".copy", function() {
    copyToClipboard($(this).parent().text());
  });

  $("#results").on("click", "tr", function(e) {
    $("#results tr.selected").removeClass();
    $(this).addClass("selected");
  });

  $("#objects").on("click", ".schema-group-title", function(e) {
    $(this).parent().toggleClass("expanded");
  });

  $("#objects").on("click", ".schema-name", function(e) {
    $(this).parent().toggleClass("expanded");
  });

  $("#objects").on("click", "li", function(e) {
    currentObject = {
      name: $(this).data("id"),
      type: $(this).data("type")
    };

    $("#objects li").removeClass("active");
    $(this).addClass("active");
    $(".current-page").data("page", 1);
    resetTableFilters();

    if (currentObject.type == "function") {
      sessionStorage.setItem("tab", "table_structure");
    } else {
      showTableInfo();
    }

    switch(sessionStorage.getItem("tab")) {
      case "table_content":
        showTableContent();
        break;
      case "table_structure":
        showTableStructure();
        break;
      case "table_constraints":
        showTableConstraints();
        break;
      case "table_indexes":
        showTableIndexes();
        break;
      default:
        showTableContent();
    }
  });

  $("#results").on("click", "a.row-action", function(e) {
    e.preventDefault();

    var action = $(this).data("action");
    var value  = $(this).data("value");

    performRowAction(action, value);
  })

  $("#results").on("click", "th", function(e) {
    if (!$("#table_content").hasClass("selected")) return;

    var sortColumn = $(this).data("name");
    var sortOrder  = $(this).data("order") === "ASC" ? "DESC" : "ASC";

    $(this).data("order", sortOrder);
    showTableContent(sortColumn, sortOrder);
  });

  $("#refresh_tables").on("click", function() {
    loadSchemas();
  });

  $("#rows_filter").on("submit", function(e) {
    e.preventDefault();

    // Move any condition built in the toolbar into the active filter list
    if ($(this).find("select.column").val()) {
      if (!addFilterFromForm()) return;
    }

    $(".current-page").data("page", 1);
    showTableContent();
  });

  $("#pagination").on("click", ".add-filter", function() {
    addFilterFromForm();
  });

  $("#active_filters").on("click", ".remove-filter", function() {
    tableFilters.splice($(this).data("index"), 1);
    renderActiveFilters();

    $(".current-page").data("page", 1);
    showTableContent();
  });

  $("#active_filters").on("click", ".clear-filters", function(e) {
    e.preventDefault();
    resetTableFilters();
    showTableContent();
  });

  $(".change-limit").on("click", function() {
    var limit = prompt("Please specify a new rows limit", getRowsLimit());

    if (limit && limit >= 1) {
      $(".current-page").data("page", 1);
      setRowsLimit(limit);
      showTableContent();
    }
  });

  $("select.filter").on("change", function(e) {
    var val = $(this).val();

    if (["null", "not_null"].indexOf(val) >= 0) {
      $(".filters input").hide().val("");
    }
    else {
      $(".filters input").show();
    }
  });

  $("button.reset-filters").on("click", function() {
    resetTableFilters();
    showTableContent();
  });

  // Automatically prefill the filter if it's not set yet
  $("select.column").on("change", function() {
    if ($("select.filter").val() == "") {
      $("select.filter").val("equal");
      $("#table_filter_value").focus();
    }
  });

  $("#pagination .next-page").on("click", function() {
    var current = $(".current-page").data("page");
    var total   = $(".current-page").data("pages");

    if (total > current) {
      $(".current-page").data("page", current + 1);
      showPaginatedTableContent();

      if (current + 1 == total) {
        $(this).prop("disabled", "disabled");
      }
    }

    if (current > 1) {
      $(".prev-page").prop("disabled", "");
    }
  });

  $("#pagination .prev-page").on("click", function() {
    var current = $(".current-page").data("page");

    if (current > 1) {
      $(".current-page").data("page", current - 1);
      $(".next-page").prop("disabled", "");
      showPaginatedTableContent();
    }

    if (current == 1) {
      $(this).prop("disabled", "disabled");
    }
  });

  $("#current_database").on("click", function(e) {
    apiCall("get", "/databases", {}, function(resp) {
      toggleDatabaseSearch();
      enableDatabaseSearch(resp);
    });
  });

  $("#database_search").change(function(e) {
    var current = $("#database_search").typeahead("getActive");
    if (current && current == $("#database_search").val()) {
      apiCall("post", "/switchdb", { db: current }, function(resp) {
        if (resp.error) {
          alert(resp.error);
          return;
        };
        window.location.reload();
      });
    };
  });

  $("#edit_connection").on("click", function() {
    if (connected) {
      $("#close_connection_window").show();
    }

    showConnectionSettings();
  });

  $("#close_connection").on("click", function() {
    if (!confirm("Are you sure you want to disconnect?")) return;

    disconnect(function() {
      showConnectionSettings();
      resetTable();
      $("#close_connection_window").hide();
    });
  });

  $("#close_connection_window").on("click", function() {
    $("#connection_window").hide();
  });

  $("#connection_url").on("change", function() {
    if ($(this).val().indexOf("localhost") != -1) {
      $("#connection_ssl").val("disable");
    }
  });

  $("#pg_host").on("change", function() {
    var value = $(this).val();

    if (value.indexOf("localhost") != -1 || value.indexOf("127.0.0.1") != -1) {
      $("#connection_ssl").val("disable");
    }
  });

  $(".connection-group-switch button").on("click", function() {
    $(".connection-group-switch button").removeClass("active");
    $(this).addClass("active");

    switch($(this).attr("data")) {
      case "scheme":
        $(".connection-scheme-group").show();
        $(".connection-standard-group").hide();
        $(".connection-ssh-group").hide();
        return;
      case "standard":
        $(".connection-scheme-group").hide();
        $(".connection-standard-group").show();
        $(".connection-ssh-group").hide();
        return;
      case "ssh":
        $(".connection-scheme-group").hide();
        $(".connection-standard-group").show();
        $(".connection-ssh-group").show();
        return;
    }
  });

  $("#connection_bookmarks").on("change", function(e) {
    var selection = $(this).val();

    var inputs = [
      $("#connection_form input[type='text']"),
      $("#connection_form input[type='password']"),
      $("#connection_ssl")
    ];

    inputs.forEach(function(selector) {
      selector.val("").prop("disabled", selection == "" ? "" : "disabled");
    });
  });

  $("#connection_form").on("submit", function(e) {
    e.preventDefault();

    var button = $(this).find("button.open-connection");
    var params = {};
    var bookmarkID = $.trim($("#connection_bookmarks").val());

    if (bookmarkID != "") {
      params["bookmark_id"] = $("#connection_bookmarks").val();
    }
    else {
      params.url = getConnectionString();
      if (params.url.length == 0) {
        return;
      }

      if ($(".connection-group-switch button.active").attr("data") == "ssh") {
        params["ssh"]              = 1
        params["ssh_host"]         = $("#ssh_host").val();
        params["ssh_port"]         = $("#ssh_port").val();
        params["ssh_user"]         = $("#ssh_user").val();
        params["ssh_password"]     = $("#ssh_password").val();
        params["ssh_key"]          = $("#ssh_key").val();
        params["ssh_key_password"] = $("#ssh_key_password").val()
      }
    }

    $("#connection_error").hide();
    button.prop("disabled", true).text("Please wait...");

    apiCall("post", "/connect", params, function(resp) {
      button.prop("disabled", false).text("Connect");

      if (resp.error) {
        connected = false;
        $("#connection_error").text(resp.error).show();
      }
      else {
        connected = true;
        loadSchemas();
        loadLocalQueries();

        $("#connection_window").hide();
        $("#current_database").text(resp.current_database);
        $("#main").show();
      }
    });
  });

  initEditor();
  addShortcutTooltips();
  bindDatabaseObjectsFilter();

  // Set session from the url
  var reqUrl = new URL(window.location);
  var sessionId = reqUrl.searchParams.get("session");

  if (sessionId && sessionId != "") {
    sessionStorage.setItem("session_id", sessionId);
    window.history.pushState({}, document.title, window.location.pathname);
  }

  getInfo(function(resp) {
    if (resp.error) {
      alert("Unable to fetch app info: " + resp.error + ". Please reload the browser page.");
      return;
    }

    appInfo = resp.app;
    appFeatures = resp.features;

    getConnection(function(resp) {
      if (resp.error) {
        connected = false;
        showConnectionSettings();
        $(".connection-actions").show();
        return;
      }

      connected = true;
      loadSchemas();
      loadLocalQueries();

      $("#current_database").text(resp.current_database);
      $("#main").show();

      if (!appFeatures.session_lock) {
        $(".connection-actions").show();
      }
    });
  });
});

