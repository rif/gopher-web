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
    		Package.save({name: pkg.name, repo: pkg.repo, description: pkg.description}, function(response){
                window.location.assign("#/");
    		});
    	}
    }])
    .controller('RemoveCtrl', ['$scope','Package', function ($scope, Package)  {
        $scope.removePackage = function(pkg){
            Package.remove({repo: pkg.repo, reason: pkg.reason}, function(response){
                $scope.removeResponse = response;
            });
        }
    }]);