'use strict';

/* Controllers */

angular.module('gopher.controllers', ['gopher.services'])
    .controller('ListCtrl', ['$scope', 'Package', function ($scope, Package) {
        $scope.packages = Package.query();
    }])
    .controller('PackageCtrl', ['$scope', 'Package', '$routeParams', function ($scope, Package, $routeParams) {
        Package.get({repo: $routeParams.repo}, function (pack) {
                $scope.pack = pack;
        });
    }])
    .controller('AddCtrl', ['$scope', 'Package', function ($scope, Package) {
    	$scope.addPackage = function(pkg){
    		Package.save(pkg, function(response){
                window.location.assign("#/"); // TODO: some flash thing
    		});
    	}
    }])
    .controller('RemoveCtrl', ['$scope','Package', function ($scope, Package)  {
        $scope.removePackage = function(pkg){
            Package.remove(pkg, function(response){
                $scope.removeResponse = response;
            });
        }
    }]);