'use strict';

angular.module('pgweb', ['ui.router.state', 'ui.router', 'ui.ace'])
.config(function($stateProvider, $urlRouterProvider) {

  $stateProvider.state('root', {
    url: '',
    abstract: true,
    views: {
      "sidebar@": {
        templateUrl: '/static/tpl/sidebar-ctrl.html',
        controller: 'SidebarCtrl'
      }
    }
  }).state('root.login', {
    url: '/',
    views: {
      "content@": {
        templateUrl: '/static/tpl/home-ctrl.html',
        controller: 'HomeCtrl'
      }
    }
  }).state('root.query', {
    url: '/table/:table/query',
    views: {
      "content@": {
        templateUrl: '/static/tpl/query-ctrl.html',
        controller: 'QueryCtrl'
      }
    }
  }).state('root.data', {
    url: '/table/:table/data',
    views: {
      "content@": {
        templateUrl: '/static/tpl/data-ctrl.html',
        controller: 'DataCtrl'
      }
    }
  }).state('root.structure', {
    url: '/table/:table/structure',
    views: {
      "content@": {
        templateUrl: '/static/tpl/structure-ctrl.html',
        controller: 'StructureCtrl'
      }
    }
  }).state('root.history', {
    url: '/history',
    views: {
      "content@": {
        templateUrl: '/static/tpl/history-ctrl.html',
        controller: 'HistoryCtrl',
        resolve: {"results": function($http) {
          return $http.get('/history');
        }}
      }
    }
  });

  $urlRouterProvider.otherwise('/');

})


.run(function($rootScope) {
  $rootScope.hello = 'world';
})


.controller('SidebarCtrl', function($scope, $rootScope, $http, $state) {
  $scope.tables = [];
  $scope.tableinfo = {};
  $http.get('/tables').success(function(data, status) {
    $scope.tables = data;
  })
  $rootScope.currentTable = null;
  $scope.selectTable = function(table) {
    $rootScope.currentTable = table;
    // Load data, transition to Data view for that table.
    var goToName = 'root.data';
    if ($state.current.name != 'root.home') {
      goToName = $state.current.name;
    }
    $state.go(goToName, {table: table});

    $http.get('/tables/' + table + '/info').success(function(data, status) {
      console.log(data, status);
      $scope.tableinfo = data;
    });
  };

})


.controller('HomeCtrl', function($scope, $rootScope, $http, $state) {
})


.controller('DataCtrl', function($scope, $http, $stateParams) {
  $scope.results = {columns: [], rows: []};
  // WARNING: we need to escape "$stateParams.table" or use prepared statements!
  $http.post('/query', {query: "SELECT * FROM \"" + $stateParams.table + "\" LIMIT 100"}).success(function(data, status) {
    $scope.results = data;
  }).error(function(data, status) {
    alert("Error: " + data.error);
  });
})


.controller('StructureCtrl', function($scope, $http, $stateParams) {
  $scope.structure = {columns: [], rows: []};
  $http.get('/tables/' + $stateParams.table).success(function(data, status) {
    $scope.structure = data;
  })

  $scope.indexes = {columns: [], rows: []};
  $http.get('/tables/' + $stateParams.table + '/indexes').success(function(data, status) {
    $scope.indexes = data;
  })
})


.controller('QueryCtrl', function($scope, $http, $stateParams) {
  $scope.query = "SELECT * FROM \"" + $stateParams.table + "\" LIMIT 10;";
  $scope.results = {columns: [], rows: []};
  $scope.loading = false;

  $scope.aceLoaded = function(_editor) {
    console.log("ACE LOADDDEEED");
    _editor.getSession().setMode("ace/mode/pgsql");
    _editor.getSession().setTabSize(2);
    _editor.getSession().setUseSoftTabs(true);
  };

  $scope.doQuery = function(query, explain) {
    $scope.loading = true;
    explain = !!explain;
    $http.post('/query', {query: query, explain: explain}).success(function(data, status) {
      $scope.results = data;
      $scope.loading = false;
    }).error(function(data, status) {
      // Use some "angular-toastr" goodness instead
      alert("Error: " + status + "\n" + data.error);
      $scope.loading = false;
    });
  };

  $scope.downloadCsv = function(query) {
    query = query.replace(/\n/g, " ");

    var url = "http://" + window.location.host + "/query?format=csv&query=" + query;
    var win = window.open(url, '_blank');
    win.focus();
  }
})


.controller('HistoryCtrl', function($scope, $http, results) {
  // Fetch history from somewhere ?
  console.log("History endpoint: ", results);

  var rows = [];
  for(var i in results.data) {
    rows.unshift([parseInt(i) + 1, results.data[i]]);
  }

  $scope.results = {
    columns: ["id", "query"],
    rows: rows
  };
})


.directive('pgContentNavigation', function($state, $rootScope, $stateParams) {
  return {
    templateUrl: "/static/tpl/content-navigation-directive.html",
    link: function($scope, $element, $attrs) {
      $scope.active = $state.current.name;
      $scope.go = function(where) {
        $state.transitionTo(where, {table: $rootScope.currentTable});
      };
    }
  }
})


.directive('pgTableView', function() {
  return {
    scope: {
      results: "=pgTableView",
    },
    templateUrl: "/static/tpl/table-view-directive.html",
    link: function($scope, $element, $attrs) {
      $scope.showResults = function() {
        return !$scope.results.error && $scope.results.rows && $scope.results.rows.length;
      }
      $scope.is_null = function(input) {
        return (input == null || input == undefined);
      };
    }
  }
})
;
