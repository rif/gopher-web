'use strict';

/* Services */


// Demonstrate how to register services
// In this case it is a simple value service.
angular.module('gopher.services', ['ngResource']).
    value('version', '0.1').
    factory('Package', function ($resource) {
        return $resource('/api/query', {}, {
            query: {method: 'GET', params: {packageId: ''}, isArray: true},
            get: {method: 'GET', isArray: true}
        });
    });
