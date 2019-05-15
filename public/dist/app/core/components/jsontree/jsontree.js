import coreModule from 'app/core/core_module';
import { JsonExplorer } from '../json_explorer/json_explorer';
coreModule.directive('jsonTree', [
    function jsonTreeDirective() {
        return {
            restrict: 'E',
            scope: {
                object: '=',
                startExpanded: '@',
                rootName: '@',
            },
            link: function (scope, elem) {
                var jsonExp = new JsonExplorer(scope.object, 3, {
                    animateOpen: true,
                });
                var html = jsonExp.render(true);
                elem.html(html);
            },
        };
    },
]);
//# sourceMappingURL=jsontree.js.map