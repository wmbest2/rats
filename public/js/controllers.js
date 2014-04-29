var ratsApp = angular.module('RatsApp',['ngRoute', 'deviceServices']);

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
 
ratsApp.controller('DeviceController', ['$scope', 'Devices', function ($scope, Devices) {
    $scope.devices = Devices.query();
}]);

ratsApp.controller('RunsController', ['$scope', '$routeParams', 'Runs', function ($scope, $routeParams, Runs) {
    if ($routeParams.id === undefined && $routeParams.device === undefined) {
        $scope.runs = Runs.query();
    } else {
        $scope.run = Runs.get({id: $routeParams.id, device: $routeParams.device});
        console.log($scope.run)
    }

    $scope.getSuiteName = function(suite) {
        if (suite.device === undefined) {
            return suite.name
        }
        return suite.device.manufacturer + " " + suite.device.model
    }
}]);
