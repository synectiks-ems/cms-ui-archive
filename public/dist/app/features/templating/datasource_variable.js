import kbn from 'app/core/utils/kbn';
import { containsVariable, assignModelProperties, variableTypes } from './variable';
var DatasourceVariable = /** @class */ (function () {
    /** @ngInject */
    function DatasourceVariable(model, datasourceSrv, variableSrv, templateSrv) {
        this.model = model;
        this.datasourceSrv = datasourceSrv;
        this.variableSrv = variableSrv;
        this.templateSrv = templateSrv;
        this.defaults = {
            type: 'datasource',
            name: '',
            hide: 0,
            label: '',
            current: {},
            regex: '',
            options: [],
            query: '',
            refresh: 1,
            skipUrlSync: false,
        };
        assignModelProperties(this, model, this.defaults);
        this.refresh = 1;
    }
    DatasourceVariable.prototype.getSaveModel = function () {
        assignModelProperties(this.model, this, this.defaults);
        // don't persist options
        this.model.options = [];
        return this.model;
    };
    DatasourceVariable.prototype.setValue = function (option) {
        return this.variableSrv.setOptionAsCurrent(this, option);
    };
    DatasourceVariable.prototype.updateOptions = function () {
        var options = [];
        var sources = this.datasourceSrv.getMetricSources({ skipVariables: true });
        var regex;
        if (this.regex) {
            regex = this.templateSrv.replace(this.regex, null, 'regex');
            regex = kbn.stringToJsRegex(regex);
        }
        for (var i = 0; i < sources.length; i++) {
            var source = sources[i];
            // must match on type
            if (source.meta.id !== this.query) {
                continue;
            }
            if (regex && !regex.exec(source.name)) {
                continue;
            }
            options.push({ text: source.name, value: source.name });
        }
        if (options.length === 0) {
            options.push({ text: 'No data sources found', value: '' });
        }
        this.options = options;
        return this.variableSrv.validateVariableSelectionState(this);
    };
    DatasourceVariable.prototype.dependsOn = function (variable) {
        if (this.regex) {
            return containsVariable(this.regex, variable.name);
        }
        return false;
    };
    DatasourceVariable.prototype.setValueFromUrl = function (urlValue) {
        return this.variableSrv.setOptionFromUrl(this, urlValue);
    };
    DatasourceVariable.prototype.getValueForUrl = function () {
        return this.current.value;
    };
    return DatasourceVariable;
}());
export { DatasourceVariable };
variableTypes['datasource'] = {
    name: 'Datasource',
    ctor: DatasourceVariable,
    description: 'Enabled you to dynamically switch the datasource for multiple panels',
};
//# sourceMappingURL=datasource_variable.js.map