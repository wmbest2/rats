var ratsApp = angular.module('RatsApp',['ngRoute', 'deviceServices', 'ui.bootstrap']);

ratsApp.config(function($routeProvider, $locationProvider) {
    //$locationProvider.html5Mode(true);
    $routeProvider

    // route for the home page
    .when('/', {
        templateUrl : 'pages/devices-list.html',
        controller  : 'DeviceController'
    })

    // route for the devices page
    .when('/runs/:id/:device', {
        templateUrl : 'pages/suite-details.html',
        controller  : 'RunController'
    })
    .when('/runs/:id', {
        templateUrl : 'pages/run-details.html',
        controller  : 'RunController'
    })
    .when('/runs', {
        templateUrl : 'pages/runs-list.html',
        controller  : 'RunsController',
        reloadOnSearch : false
    })
    .when('/devices', {
        templateUrl : 'pages/devices-list.html',
        controller  : 'DeviceController'
    })
    .otherwise({
        redirectTo: '/'
    });
});

ratsApp.controller('MainController', function ($scope) {
    $scope.menu = 'pages/menu.html';
});
 
ratsApp.controller('DeviceController', ['$scope', '$interval', '$window', 'Devices', function ($scope, $interval, $window, Devices) {
    $scope.devices = [];

    $scope.refreshing = true;

    $scope.toggleRefresh = function() {
        $scope.refreshing=!$scope.refreshing;
        if ($scope.refreshing) {
            $scope.promise = $interval(tick, 1000);
        } else {
            $scope.cancelRefresh();
        }
    };

    var tick = function() {
        Devices.query(function(data){
            $scope.devices = data;
        });
    };

    $scope.cancelRefresh = function(){
        $interval.cancel($scope.promise);
    };

    $scope.refreshClass = function() {
        if ($scope.refreshing) {
            return "glyphicon-refresh spin";
        } 
        return "glyphicon-pause blink";
    };

    $scope.copySerial = function(serial) {
        $window.prompt("Copy to clipboard: Ctrl+C, Enter", "-s " + serial);
    }

    tick();
    $scope.promise = $interval(tick, 1000);

    $scope.predicate = ['manufacturer','model','version'];

    $scope.$on('$destroy', $scope.cancelRefresh);
}]);

ratsApp.controller('RunsController', ['$scope', '$routeParams', '$location', 'Runs', function ($scope, $routeParams, $location, Runs) {
    $scope.runs = [];
    $scope.maxPages = 5;
    $scope.currentPage = $routeParams.page;
    if ($scope.currentPage === undefined) {
        $scope.currentPage = 1;
    }

    // set minimum total so that paging works correctly
    $scope.meta = {"total": $scope.currentPage * 10};

    Runs.get({page: $scope.currentPage}, function(data) {
        $scope.runs = data.runs;
        $scope.meta = data.meta;
    });

    $scope.onPageChange = function() {
        $location.search({page: $scope.currentPage})
        Runs.get({page: $scope.currentPage}, function(data) {
            $scope.runs = data.runs;
            $scope.meta = data.meta;
        });
    };

    $scope.testSuccess = function(test) {
        return {
            'progress-bar-danger': test.failure !== undefined || test.error != undefined, 
            'progress-bar-success': test.failure === undefined && test.error === undefined
        };
    }

    $scope.prettyPrintTime = function(timestamp) {
        return moment(timestamp).calendar();
    }
}]);

ratsApp.controller('RunController', ['$scope', '$routeParams', 'Runs', function ($scope, $routeParams, Runs) {
    $scope.run = {};
    $scope.run = Runs.get({id: $routeParams.id, device: $routeParams.device});

    $scope.getSuiteName = function(suite) {
        if (suite.device === undefined) {
            return suite.name
        }
        return suite.device.manufacturer + " " + suite.device.model
    }

    $scope.testSuccess = function(test) {
        return {
            'progress-bar-danger': test.failure !== undefined || test.error != undefined, 
            'progress-bar-success': test.failure === undefined && test.error === undefined
        };
    }
}]);

ratsApp.controller('ErrorCtrl', ['$scope', function($scope) {
    $scope.hide_errors = true;
}]);
