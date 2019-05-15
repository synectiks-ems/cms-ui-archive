// Utils
import config from 'app/core/config';
import appEvents from 'app/core/app_events';
import coreModule from 'app/core/core_module';
import { removePanel } from 'app/features/dashboard/utils/panel';
var DashboardCtrl = /** @class */ (function () {
    /** @ngInject */
    function DashboardCtrl($scope, keybindingSrv, timeSrv, variableSrv, alertingSrv, dashboardSrv, unsavedChangesSrv, dashboardViewStateSrv, annotationsSrv, playlistSrv) {
        this.$scope = $scope;
        this.keybindingSrv = keybindingSrv;
        this.timeSrv = timeSrv;
        this.variableSrv = variableSrv;
        this.alertingSrv = alertingSrv;
        this.dashboardSrv = dashboardSrv;
        this.unsavedChangesSrv = unsavedChangesSrv;
        this.dashboardViewStateSrv = dashboardViewStateSrv;
        this.annotationsSrv = annotationsSrv;
        this.playlistSrv = playlistSrv;
        // temp hack due to way dashboards are loaded
        // can't use controllerAs on route yet
        $scope.ctrl = this;
        // TODO: break out settings view to separate view & controller
        this.editTab = 0;
        // funcs called from React component bindings and needs this binding
        this.getPanelContainer = this.getPanelContainer.bind(this);
    }
    DashboardCtrl.prototype.setupDashboard = function (data) {
        try {
            this.setupDashboardInternal(data);
        }
        catch (err) {
            this.onInitFailed(err, 'Dashboard init failed', true);
        }
    };
    DashboardCtrl.prototype.setupDashboardInternal = function (data) {
        var _this = this;
        var dashboard = this.dashboardSrv.create(data.dashboard, data.meta);
        this.dashboardSrv.setCurrent(dashboard);
        // init services
        this.timeSrv.init(dashboard);
        this.alertingSrv.init(dashboard, data.alerts);
        this.annotationsSrv.init(dashboard);
        // template values service needs to initialize completely before
        // the rest of the dashboard can load
        this.variableSrv
            .init(dashboard)
            // template values failes are non fatal
            .catch(this.onInitFailed.bind(this, 'Templating init failed', false))
            // continue
            .finally(function () {
            _this.dashboard = dashboard;
            _this.dashboard.processRepeats();
            _this.dashboard.updateSubmenuVisibility();
            _this.dashboard.autoFitPanels(window.innerHeight);
            _this.unsavedChangesSrv.init(dashboard, _this.$scope);
            // TODO refactor ViewStateSrv
            _this.$scope.dashboard = dashboard;
            _this.dashboardViewState = _this.dashboardViewStateSrv.create(_this.$scope);
            _this.keybindingSrv.setupDashboardBindings(_this.$scope, dashboard);
            _this.setWindowTitleAndTheme();
            appEvents.emit('dashboard-initialized', dashboard);
        })
            .catch(this.onInitFailed.bind(this, 'Dashboard init failed', true));
    };
    DashboardCtrl.prototype.onInitFailed = function (msg, fatal, err) {
        console.log(msg, err);
        if (err.data && err.data.message) {
            err.message = err.data.message;
        }
        else if (!err.message) {
            err = { message: err.toString() };
        }
        this.$scope.appEvent('alert-error', [msg, err.message]);
        // protect against  recursive fallbacks
        if (fatal && !this.loadedFallbackDashboard) {
            this.loadedFallbackDashboard = true;
            this.setupDashboard({ dashboard: { title: 'Dashboard Init failed' } });
        }
    };
    DashboardCtrl.prototype.templateVariableUpdated = function () {
        this.dashboard.processRepeats();
    };
    DashboardCtrl.prototype.setWindowTitleAndTheme = function () {
        window.document.title = config.windowTitlePrefix + this.dashboard.title;
    };
    DashboardCtrl.prototype.showJsonEditor = function (evt, options) {
        var model = {
            object: options.object,
            updateHandler: options.updateHandler,
        };
        this.$scope.appEvent('show-dash-editor', {
            src: 'public/app/partials/edit_json.html',
            model: model,
        });
    };
    DashboardCtrl.prototype.getDashboard = function () {
        return this.dashboard;
    };
    DashboardCtrl.prototype.getPanelContainer = function () {
        return this;
    };
    DashboardCtrl.prototype.onRemovingPanel = function (evt, options) {
        options = options || {};
        if (!options.panelId) {
            return;
        }
        var panelInfo = this.dashboard.getPanelInfoById(options.panelId);
        removePanel(this.dashboard, panelInfo.panel, true);
    };
    DashboardCtrl.prototype.onDestroy = function () {
        if (this.dashboard) {
            this.dashboard.destroy();
        }
    };
    DashboardCtrl.prototype.init = function (dashboard) {
        this.$scope.onAppEvent('show-json-editor', this.showJsonEditor.bind(this));
        this.$scope.onAppEvent('template-variable-value-updated', this.templateVariableUpdated.bind(this));
        this.$scope.onAppEvent('panel-remove', this.onRemovingPanel.bind(this));
        this.$scope.$on('$destroy', this.onDestroy.bind(this));
        this.setupDashboard(dashboard);
    };
    return DashboardCtrl;
}());
export { DashboardCtrl };
coreModule.controller('DashboardCtrl', DashboardCtrl);
//# sourceMappingURL=dashboard_ctrl.js.map