'use strict';

/* Directives */

angular.module('gopher.directives', []).
  directive('autoComplete', [function() {
    return function(scope, iElement, iAttrs) {
            iElement.autocomplete({
                source: scope[iAttrs.uiItems],
                select: function() {
                    $timeout(function() {
                      iElement.trigger('input');
                    }, 0);
                }
            });
    };
  }]);
