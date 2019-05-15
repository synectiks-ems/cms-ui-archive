import * as tslib_1 from "tslib";
import React from 'react';
import { shallow } from 'enzyme';
import { SideMenu } from './SideMenu';
import appEvents from '../../app_events';
import { contextSrv } from 'app/core/services/context_srv';
jest.mock('../../app_events', function () { return ({
    emit: jest.fn(),
}); });
jest.mock('app/core/services/context_srv', function () { return ({
    contextSrv: {
        sidemenu: true,
        user: {},
        isSignedIn: false,
        isGrafanaAdmin: false,
        isEditor: false,
        hasEditPermissionFolders: false,
        toggleSideMenu: jest.fn(),
    },
}); });
var setup = function (propOverrides) {
    var props = Object.assign({
        loginUrl: '',
        user: {},
        mainLinks: [],
        bottomeLinks: [],
        isSignedIn: false,
    }, propOverrides);
    return shallow(React.createElement(SideMenu, tslib_1.__assign({}, props)));
};
describe('Render', function () {
    it('should render component', function () {
        var wrapper = setup();
        expect(wrapper).toMatchSnapshot();
    });
});
describe('Functions', function () {
    describe('toggle side menu', function () {
        var wrapper = setup();
        var instance = wrapper.instance();
        instance.toggleSideMenu();
        it('should call contextSrv.toggleSideMenu', function () {
            expect(contextSrv.toggleSideMenu).toHaveBeenCalled();
        });
        it('should emit toggle sidemenu event', function () {
            expect(appEvents.emit).toHaveBeenCalledWith('toggle-sidemenu');
        });
    });
    describe('toggle side menu on mobile', function () {
        var wrapper = setup();
        var instance = wrapper.instance();
        instance.toggleSideMenuSmallBreakpoint();
        it('should emit toggle sidemenu event', function () {
            expect(appEvents.emit).toHaveBeenCalledWith('toggle-sidemenu-mobile');
        });
    });
});
//# sourceMappingURL=SideMenu.test.js.map