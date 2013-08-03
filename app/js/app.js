'use strict';


// Declare app level module which depends on filters, and services
angular.module('gopher', ['gopher.filters', 'gopher.services', 'gopher.directives', 'gopher.controllers']).
  config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/', {templateUrl: 'partials/list.html', controller: 'ListCtrl'});
    $routeProvider.when('/package/*repo', {templateUrl: 'partials/package.html', controller: 'PackageCtrl'});
    $routeProvider.otherwise({redirectTo: '/'});
  }]);
