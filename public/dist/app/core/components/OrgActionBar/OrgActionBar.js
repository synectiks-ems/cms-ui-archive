import * as tslib_1 from "tslib";
import React, { PureComponent } from 'react';
import LayoutSelector from '../LayoutSelector/LayoutSelector';
var OrgActionBar = /** @class */ (function (_super) {
    tslib_1.__extends(OrgActionBar, _super);
    function OrgActionBar() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    OrgActionBar.prototype.render = function () {
        var _a = this.props, searchQuery = _a.searchQuery, layoutMode = _a.layoutMode, onSetLayoutMode = _a.onSetLayoutMode, linkButton = _a.linkButton, setSearchQuery = _a.setSearchQuery, target = _a.target;
        var linkProps = { href: linkButton.href, target: undefined };
        if (target) {
            linkProps.target = target;
        }
        return (React.createElement("div", { className: "page-action-bar" },
            React.createElement("div", { className: "gf-form gf-form--grow" },
                React.createElement("label", { className: "gf-form--has-input-icon" },
                    React.createElement("input", { type: "text", className: "gf-form-input width-20", value: searchQuery, onChange: function (event) { return setSearchQuery(event.target.value); }, placeholder: "Filter by name or type" }),
                    React.createElement("i", { className: "gf-form-input-icon fa fa-search" })),
                React.createElement(LayoutSelector, { mode: layoutMode, onLayoutModeChanged: function (mode) { return onSetLayoutMode(mode); } })),
            React.createElement("div", { className: "page-action-bar__spacer" }),
            React.createElement("a", tslib_1.__assign({ className: "btn btn-success" }, linkProps), linkButton.title)));
    };
    return OrgActionBar;
}(PureComponent));
export default OrgActionBar;
//# sourceMappingURL=OrgActionBar.js.map