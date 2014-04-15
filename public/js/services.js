var phoneServices = angular.module('deviceServices', ['ngResource']);
phoneServices.factory('Devices', ['$resource',
    function($resource){
        return $resource('api/devices', {}, {
            query: {method:'GET', isArray:true}
        });
    }]);
