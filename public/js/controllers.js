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
        controller  : 'RunsController'
    })
    .when('/runs/:id', {
        templateUrl : 'pages/run-details.html',
        controller  : 'RunsController'
    })
    .when('/runs', {
        templateUrl : 'pages/runs-list.html',
        controller  : 'RunsController'
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
 
ratsApp.controller('DeviceController', ['$scope', '$interval', 'Devices', function ($scope, $interval, Devices) {
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

    tick();
    $scope.promise = $interval(tick, 1000);

    $scope.predicate = ['manufacturer','model','version'];

    $scope.$on('$destroy', $scope.cancelRefresh);
}]);

ratsApp.controller('RunsController', ['$scope', '$routeParams', 'Runs', function ($scope, $routeParams, Runs) {
    $scope.runs = [];
    $scope.run = {};
    if ($routeParams.id === undefined && $routeParams.device === undefined) {
        Runs.query(function(data) {
            $scope.runs = data;
        });
    } else {
        $scope.run = Runs.get({id: $routeParams.id, device: $routeParams.device});
    }

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
