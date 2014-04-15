var ratsApp = angular.module('RatsApp',['ngRoute', 'deviceServices']);

ratsApp.config(function($routeProvider, $locationProvider) {
    $locationProvider.html5Mode(true);
    $routeProvider

    // route for the home page
    .when('/', {
        templateUrl : 'pages/devices-list.html',
        controller  : 'DeviceController'
    })

    // route for the devices page
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
