'use strict';


// Declare app level module which depends on filters, and services
angular.module('gopher', ['gopher.filters', 'gopher.services', 'gopher.directives', 'gopher.controllers']).
  config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/', {templateUrl: 'partials/list.html', controller: 'ListCtrl'});
    $routeProvider.when('/package/*repo', {templateUrl: 'partials/package.html', controller: 'PackageCtrl'});
    $routeProvider.when('/add', {templateUrl: 'partials/add.html', controller: 'AddCtrl'});
    $routeProvider.when('/remove', {templateUrl: 'partials/remove.html', controller: 'RemoveCtrl'});
    $routeProvider.otherwise({redirectTo: '/'});
  }]);
