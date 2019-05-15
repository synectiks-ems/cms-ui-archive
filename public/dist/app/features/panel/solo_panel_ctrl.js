import angular from 'angular';
import locationUtil from 'app/core/utils/location_util';
import appEvents from 'app/core/app_events';
var SoloPanelCtrl = /** @class */ (function () {
    /** @ngInject */
    function SoloPanelCtrl($scope, $routeParams, $location, dashboardLoaderSrv, contextSrv, backendSrv) {
        var panelId;
        $scope.init = function () {
            contextSrv.sidemenu = false;
            appEvents.emit('toggle-sidemenu-hidden');
            var params = $location.search();
            panelId = parseInt(params.panelId, 10);
            appEvents.on('dashboard-initialized', $scope.initPanelScope);
            // if no uid, redirect to new route based on slug
            if (!($routeParams.type === 'script' || $routeParams.type === 'snapshot') && !$routeParams.uid) {
                backendSrv.getDashboardBySlug($routeParams.slug).then(function (res) {
                    if (res) {
                        var url = locationUtil.stripBaseFromUrl(res.meta.url.replace('/d/', '/d-solo/'));
                        $location.path(url).replace();
                    }
                });
                return;
            }
            dashboardLoaderSrv.loadDashboard($routeParams.type, $routeParams.slug, $routeParams.uid).then(function (result) {
                result.meta.soloMode = true;
                $scope.initDashboard(result, $scope);
            });
        };
        $scope.initPanelScope = function () {
            var panelInfo = $scope.dashboard.getPanelInfoById(panelId);
            // fake row ctrl scope
            $scope.ctrl = {
                dashboard: $scope.dashboard,
            };
            $scope.panel = panelInfo.panel;
            $scope.panel.soloMode = true;
            $scope.$index = 0;
            if (!$scope.panel) {
                $scope.appEvent('alert-error', ['Panel not found', '']);
                return;
            }
        };
        $scope.init();
    }
    return SoloPanelCtrl;
}());
export { SoloPanelCtrl };
angular.module('grafana.routes').controller('SoloPanelCtrl', SoloPanelCtrl);
//# sourceMappingURL=solo_panel_ctrl.js.map