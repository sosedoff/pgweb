var editor;

function apiCall(method, path, params, cb) {
  $.ajax({
    url: path, 
    method: method,
    cache: false,
    data: params,
    success: function(data) {
      cb(data);
    },
    error: function(xhr, status, data) {
      cb(jQuery.parseJSON(xhr.responseText));
    }
  });
}

function getTables(cb)                { apiCall("get", "/tables", {}, cb); }
function getTableStructure(table, cb) { apiCall("get", "/tables/" + table, {}, cb); }
function getTableIndexes(table, cb)   { apiCall("get", "/tables/" + table + "/indexes", {}, cb); }
function getHistory(cb)               { apiCall("get", "/history", {}, cb); }

function executeQuery(query, cb) {
  apiCall("post", "/query", { query: query }, cb);
}

function explainQuery(query, cb) {
  apiCall("post", "/explain", { query: query }, cb);
}

function loadTables() {
  getTables(function(data) {
    data.forEach(function(item) {
      $("<li><span>" + item + "</span></li>").appendTo("#tables");
    });
  });
}

function escapeHtml(str) {
  if (str != null || str != undefined) {
    return jQuery("<div/>").text(str).html();
  }

  return "<span class='null'>null</span>";
}

function resetTable() {
  $("#results").
    attr("data-mode", "").
    text("").
    removeClass("empty").
    removeClass("history");
}

function buildTable(results) {
  resetTable();

  if (results.error) {
    $("<tr><td>ERROR: " + results.error + "</tr></tr>").appendTo("#results");
    $("#results").addClass("empty");
    return;
  }

  if (!results.rows) {
    $("<tr><td>No records found</tr></tr>").appendTo("#results");
    $("#results").addClass("empty");
    return;
  }

  var cols = "";
  var rows = "";

  results.columns.forEach(function(col) {
    cols += "<th data='" + col + "'>" + col + "</th>";
  });

  results.rows.forEach(function(row) {
    var r = "";
    for (i in row) { r += "<td><div>" + escapeHtml(row[i]) + "</div></td>"; }
    rows += "<tr>" + r + "</tr>";
  });

  $("<thead>" + cols + "</thead><tbody>" + rows + "</tobdy>").appendTo("#results");
}

function setCurrentTab(id) {
  $("#nav ul li.selected").removeClass("selected");
  $("#" + id).addClass("selected");
}

function showQueryHistory() {
  getHistory(function(data) {
    var rows = [];

    for(i in data) {
      rows.unshift([parseInt(i) + 1, data[i]]);
    }

    buildTable({ columns: ["id", "query"], rows: rows });
  
    setCurrentTab("table_history");  
    $("#input").hide();
    $("#output").addClass("full");
    $("#results").addClass("history");
  });
}

function showTableIndexes() {
  var name = $("#tables li.selected").text();

  if (name.length == 0) {
    alert("Please select a table!");
    return;
  }

  getTableIndexes(name, function(data) {
    setCurrentTab("table_indexes");
    buildTable(data);

    $("#input").hide();
    $("#output").addClass("full");
  });
}

function showTableInfo() {
  var name = $("#tables li.selected").text();

  if (name.length == 0) {
    alert("Please select a table!");
    return;
  }

  apiCall("get", "/tables/" + name + "/info", {}, function(data) {
    $(".table-information ul").show();
    $("#table_total_size").text(data.total_size);
    $("#table_data_size").text(data.data_size);
    $("#table_index_size").text(data.index_size);
    $("#table_rows_count").text(data.rows_count);
    $("#table_encoding").text("Unknown");
  });
}

function showTableContent() {
  var name = $("#tables li.selected").text();

  if (name.length == 0) {
    alert("Please select a table!");
    return;
  }

  var query = "SELECT * FROM \"" + name + "\" LIMIT 100;";

  executeQuery(query, function(data) {
    buildTable(data);
    setCurrentTab("table_content");

    $("#results").attr("data-mode", "browse");
    $("#input").hide();
    $("#output").addClass("full");
  });
}

function showTableStructure() {
  var name = $("#tables li.selected").text();

  if (name.length == 0) {
    alert("Please select a table!");
    return;
  }

  getTableStructure(name, function(data) {
    setCurrentTab("table_structure");
    buildTable(data);
  });
}

function showQueryPanel() {
  setCurrentTab("table_query");
  editor.focus();

  $("#input").show();
  $("#output").removeClass("full");
}

function runQuery() {
  setCurrentTab("table_query");

  $("#run, #explain, #csv").prop("disabled", true);
  $("#query_progress").show();

  var query = $.trim(editor.getValue());

  if (query.length == 0) {
    $("#run, #explain, #csv").prop("disabled", false);
    $("#query_progress").hide();
    return;
  }

  executeQuery(query, function(data) {
    buildTable(data);

    $("#run, #explain, #csv").prop("disabled", false);
    $("#query_progress").hide();
    $("#input").show();
    $("#output").removeClass("full");
  });
}

function runExplain() {
  setCurrentTab("table_query");

  $("#run, #explain, #csv").prop("disabled", true);
  $("#query_progress").show();

  var query = $.trim(editor.getValue());

  if (query.length == 0) {
    $("#run, #explain, #csv").prop("disabled", false);
    $("#query_progress").hide();
    return;
  }

  explainQuery(query, function(data) {
    buildTable(data);

    $("#run, #explain, #csv").prop("disabled", false);
    $("#query_progress").hide();
    $("#input").show();
    $("#output").removeClass("full");
  });
}

function exportToCSV() {
  var query = $.trim(editor.getValue());

  if (query.length == 0) {
    return;
  }

  // Replace line breaks with spaces and properly encode query
  query = window.encodeURI(query.replace(/\n/g, " "));

  var url = "http://" + window.location.host + "/query?format=csv&query=" + query;
  var win = window.open(url, '_blank');

  setCurrentTab("table_query");
  win.focus();
}

function initEditor() {
  editor = ace.edit("custom_query");

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

$(document).ready(function() {
  initEditor();
  addShortcutTooltips();

  $("#table_content").on("click",   function() { showTableContent();   });
  $("#table_structure").on("click", function() { showTableStructure(); });
  $("#table_indexes").on("click",   function() { showTableIndexes();   });
  $("#table_history").on("click",   function() { showQueryHistory();   });
  $("#table_query").on("click",     function() { showQueryPanel();     });

  $("#run").on("click", function() {
    runQuery();
  });

  $("#explain").on("click", function() {
    runExplain();
  });

  $("#csv").on("click", function() {
    exportToCSV();
  });

  $("#results").on("click", "tr", function() {
    $("#results tr.selected").removeClass();
    $(this).addClass("selected");
  });

  $("#tables").on("click", "li", function() {
    $("#tables li.selected").removeClass("selected");
    $(this).addClass("selected");
    showTableContent();
    showTableInfo();
  });

  loadTables();
});