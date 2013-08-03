'use strict';

/* Controllers */

angular.module('gopher.controllers', ['gopher.services']).
    controller('ListCtrl', ['$scope', 'Package', function ($scope, Package) {
        $scope.packages = Package.query();
    }])
    .controller('PackageCtrl', ['$scope', 'Package', '$routeParams', function ($scope, Package, $routeParams) {
        Package.get({packageId: $routeParams.packageId}, function (pack) {
                $scope.pack = pack;
        });
    }]);