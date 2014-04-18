(function($) {
    function initPage() {
        $('ul.nav > li > a[href="' + document.location.pathname + '"]').parent().addClass('active');
        $('a').tooltip();
    }
    $(initPage);

})(jQuery);

var myModule = angular.module('MediatorApp', ['ngAnimate', 'ui.bootstrap'], function ($interpolateProvider) {
  $interpolateProvider.startSymbol('[[');
  $interpolateProvider.endSymbol(']]');
});

function MediatorCtrl($scope, $timeout) {
    $scope.stories = [];
    $scope.errors = [];
    $scope.connection = null;
    $scope.messages = [];
    $scope.orderProp = '-LastMentioned';

    $scope.NewConnection = function() {
        var wsproto = "";
        if (document.location.protocol == "https:") {
            wsproto = "wss";
        } else {
            wsproto = "ws";
        }
        connection = new WebSocket(wsproto+"://"+document.location.host+'/messages');

        connection.onopen = function () {
            $scope.connection = connection;
        };

        connection.onclose = function (e) {
            $scope.connection = null;
            $scope.NewConnection();
        };

        connection.onerror = function (error) {
            console.log('WebSocket Error ' + error);
            $scope.$apply(function () {
                $scope.errors.push(error);
            });
        };

        connection.onmessage = function(e) {
            $scope.$apply(function () {
                var msg = JSON.parse(e.data);
                if ("Tweet" in msg) {
                    var updated = false;
                    for (var i = 0; i < $scope.stories.length; i++) {
                        if ($scope.stories[i].Url == msg.Story.Url) {
                            $scope.stories[i].Count = msg.Count;
                            $scope.stories[i].LastMentioned = msg.Tweet.Published;
                            updated = true;
                            break;
                        }
                    }
                    if (updated === false) {
                        var s = msg.Story;
                        s.Count = msg.Count;
                        s.LastMentioned = msg.Tweet.Published;
                        $scope.stories.unshift(s);
                    }
                }
            });
        };
    };

    $(window).on("pageshow", function() {
        $scope.NewConnection();
    });

    $(window).on("pagehide", function() {
        if ($scope.connection !== null) {
            $scope.connection.close();
        }
    });

    $scope.displayMessage = function(message) {
        $scope.messages.push(message);
        $timeout(function() {
            $scope.messages.shift();
        }, 10000);
    };

    $scope.formatWhen = function(when) {
        return when.substring(11,19);
    };

    $scope.storyID = function(url) {
        return url.replace("https://medium.com/", "");
    };

}
