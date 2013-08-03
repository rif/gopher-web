'use strict';

/* Filters */

angular.module('gopher.filters', []).
  filter('urlencode', [function() {
    return function(text) {
      return encodeURIComponent(text);
    }
  }]);
