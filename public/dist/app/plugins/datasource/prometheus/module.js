import { PrometheusDatasource } from './datasource';
import { PrometheusQueryCtrl } from './query_ctrl';
import { PrometheusConfigCtrl } from './config_ctrl';
import PrometheusStartPage from './components/PromStart';
import PromQueryField from './components/PromQueryField';
var PrometheusAnnotationsQueryCtrl = /** @class */ (function () {
    function PrometheusAnnotationsQueryCtrl() {
    }
    PrometheusAnnotationsQueryCtrl.templateUrl = 'partials/annotations.editor.html';
    return PrometheusAnnotationsQueryCtrl;
}());
export { PrometheusDatasource as Datasource, PrometheusQueryCtrl as QueryCtrl, PrometheusConfigCtrl as ConfigCtrl, PrometheusAnnotationsQueryCtrl as AnnotationsQueryCtrl, PromQueryField as ExploreQueryField, PrometheusStartPage as ExploreStartPage, };
//# sourceMappingURL=module.js.map