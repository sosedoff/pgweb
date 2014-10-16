function apiCall(path, cb) {
  $.getJSON(path, function(resp) { cb(resp); });
}

function getTables(cb)                { apiCall("/tables", cb); }
function getTableStructure(table, cb) { apiCall("/tables/" + table, cb); }
function getTableIndexes(table, cb)   { apiCall("/tables/" + table + "/indexes", cb); }
function getHistory(cb)               { apiCall("/history", cb); }

function executeQuery(query, cb) {
  $.ajax({
    url: "/query",
    method: "post",
    cache: false,
    data: { query: query, format: "json" },
    success: function(data) {
      cb(data);
    },
    error: function(xhr, status, data) {
      cb(jQuery.parseJSON(xhr.responseText));
    }
  });
}

function loadTables() {
  getTables(function(data) {
    data.forEach(function(item) {
      $("<li>" + item + "</li>").appendTo("#tables");
    });
  });
}

function escapeHtml(str) {
  if (str != null || str != undefined) {
    return jQuery("<div/>").text(str).html();
  }

  return "<span class='null'>null</span>";
}

function buildTable(results) {
  $("#results").text("").removeClass("empty");

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
  var rows = ""

  results.columns.forEach(function(col) {
    cols += "<th>" + col + "</th>";
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

function showTableContent() {
  var name = $("#tables li.selected").text();

  if (name.length == 0) {
    alert("Please select a table!");
    return;
  }

  var query = "SELECT * FROM " + name + " LIMIT 100;";

  executeQuery(query, function(data) {
    buildTable(data);
    setCurrentTab("table_content");

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

function runQuery() {
  setCurrentTab("table_query");

  $("#run").attr("disabled", "disabled");
  $("#explain").attr("disabled", "disabled");
  $("#query_progress").show();

  var query = $.trim(editor.getValue());

  if (query.length == 0) {
    return;
  }

  executeQuery(query, function(data) {
    buildTable(data);

    $("#run").removeAttr("disabled");
    $("#explain").removeAttr("disabled");
    $("#query_progress").hide();
    $("#input").show();
    $("#output").removeClass("full");
  });
}

function runExplain() {
  setCurrentTab("table_query");

  $("#run").attr("disabled", "disabled");
  $("#explain").attr("disabled", "disabled");
  $("#query_progress").show();

  var query = $.trim(editor.getValue());

  if (query.length == 0) {
    return;
  }

  query = "EXPLAIN ANALYZE " + query;

  executeQuery(query, function(data) {
    buildTable(data);

    $("#run").removeAttr("disabled");
    $("#explain").removeAttr("disabled");
    $("#query_progress").hide();
    $("#input").show();
    $("#output").removeClass("full");
  });
}

var editor;

$(document).ready(function() {
  editor = ace.edit("custom_query");
  editor.getSession().setMode("ace/mode/pgsql");
  editor.getSession().setTabSize(2);
  editor.getSession().setUseSoftTabs(true);

  $("#table_content").on("click",   function() { showTableContent();   });
  $("#table_structure").on("click", function() { showTableStructure(); });
  $("#table_indexes").on("click",   function() { showTableIndexes();   });
  $("#table_history").on("click",   function() { showQueryHistory();   });

  $("#table_query").on("click", function() {
    setCurrentTab("table_query");
    $("#input").show();
    $("#output").removeClass("full");
  });

  $("#run").on("click", function() {
    runQuery();
  });

  $("#explain").on("click", function() {
    runExplain();
  });

  $("#results").on("click", "tr", function() {
    $("#results tr.selected").removeClass();
    $(this).addClass("selected");
  });

  $("#tables").on("click", "li", function() {
    $("#tables li.selected").removeClass("selected");
    $(this).addClass("selected");
    showTableContent();
  });

  loadTables();
});