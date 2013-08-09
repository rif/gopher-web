'use strict';

/* Services */


// Demonstrate how to register services
// In this case it is a simple value service.
angular.module('gopher.services', ['ngResource'])
.value('version', '0.1')
.factory('Package', function ($resource) {
	return $resource('/api/pkg', {}, {
		query: {method: 'GET', params: {repo: 'all'}, isArray: true},
	});
})
.factory('Category', function ($resource) {
	return $resource('/api/cat', {}, {
		query: {method: 'GET', isArray: true},
	});
});
