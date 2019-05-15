import _ from 'lodash';
import coreModule from 'app/core/core_module';
import impressionSrv from 'app/core/services/impression_srv';
import store from 'app/core/store';
import { contextSrv } from 'app/core/services/context_srv';
var SearchSrv = /** @class */ (function () {
    /** @ngInject */
    function SearchSrv(backendSrv, $q) {
        this.backendSrv = backendSrv;
        this.$q = $q;
        this.recentIsOpen = store.getBool('search.sections.recent', true);
        this.starredIsOpen = store.getBool('search.sections.starred', true);
    }
    SearchSrv.prototype.getRecentDashboards = function (sections) {
        var _this = this;
        return this.queryForRecentDashboards().then(function (result) {
            if (result.length > 0) {
                sections['recent'] = {
                    title: 'Recent',
                    icon: 'fa fa-clock-o',
                    score: -1,
                    removable: true,
                    expanded: _this.recentIsOpen,
                    toggle: _this.toggleRecent.bind(_this),
                    items: result,
                };
            }
        });
    };
    SearchSrv.prototype.queryForRecentDashboards = function () {
        var dashIds = _.take(impressionSrv.getDashboardOpened(), 5);
        if (dashIds.length === 0) {
            return Promise.resolve([]);
        }
        return this.backendSrv.search({ dashboardIds: dashIds }).then(function (result) {
            return dashIds
                .map(function (orderId) {
                return _.find(result, { id: orderId });
            })
                .filter(function (hit) { return hit && !hit.isStarred; });
        });
    };
    SearchSrv.prototype.toggleRecent = function (section) {
        this.recentIsOpen = section.expanded = !section.expanded;
        store.set('search.sections.recent', this.recentIsOpen);
        if (!section.expanded || section.items.length) {
            return Promise.resolve(section);
        }
        return this.queryForRecentDashboards().then(function (result) {
            section.items = result;
            return Promise.resolve(section);
        });
    };
    SearchSrv.prototype.toggleStarred = function (section) {
        this.starredIsOpen = section.expanded = !section.expanded;
        store.set('search.sections.starred', this.starredIsOpen);
        return Promise.resolve(section);
    };
    SearchSrv.prototype.getStarred = function (sections) {
        var _this = this;
        if (!contextSrv.isSignedIn) {
            return Promise.resolve();
        }
        return this.backendSrv.search({ starred: true, limit: 5 }).then(function (result) {
            if (result.length > 0) {
                sections['starred'] = {
                    title: 'Starred',
                    icon: 'fa fa-star-o',
                    score: -2,
                    expanded: _this.starredIsOpen,
                    toggle: _this.toggleStarred.bind(_this),
                    items: result,
                };
            }
        });
    };
    SearchSrv.prototype.search = function (options) {
        var _this = this;
        var sections = {};
        var promises = [];
        var query = _.clone(options);
        var hasFilters = options.query ||
            (options.tag && options.tag.length > 0) ||
            options.starred ||
            (options.folderIds && options.folderIds.length > 0);
        if (!options.skipRecent && !hasFilters) {
            promises.push(this.getRecentDashboards(sections));
        }
        if (!options.skipStarred && !hasFilters) {
            promises.push(this.getStarred(sections));
        }
        query.folderIds = query.folderIds || [];
        if (!hasFilters) {
            query.folderIds = [0];
        }
        promises.push(this.backendSrv.search(query).then(function (results) {
            return _this.handleSearchResult(sections, results);
        }));
        return this.$q.all(promises).then(function () {
            return _.sortBy(_.values(sections), 'score');
        });
    };
    SearchSrv.prototype.handleSearchResult = function (sections, results) {
        if (results.length === 0) {
            return sections;
        }
        // create folder index
        for (var _i = 0, results_1 = results; _i < results_1.length; _i++) {
            var hit = results_1[_i];
            if (hit.type === 'dash-folder') {
                sections[hit.id] = {
                    id: hit.id,
                    uid: hit.uid,
                    title: hit.title,
                    expanded: false,
                    items: [],
                    toggle: this.toggleFolder.bind(this),
                    url: hit.url,
                    icon: 'fa fa-folder',
                    score: _.keys(sections).length,
                };
            }
        }
        for (var _a = 0, results_2 = results; _a < results_2.length; _a++) {
            var hit = results_2[_a];
            if (hit.type === 'dash-folder') {
                continue;
            }
            var section = sections[hit.folderId || 0];
            if (!section) {
                if (hit.folderId) {
                    section = {
                        id: hit.folderId,
                        uid: hit.folderUid,
                        title: hit.folderTitle,
                        url: hit.folderUrl,
                        items: [],
                        icon: 'fa fa-folder-open',
                        toggle: this.toggleFolder.bind(this),
                        score: _.keys(sections).length,
                    };
                }
                else {
                    section = {
                        id: 0,
                        title: 'General',
                        items: [],
                        icon: 'fa fa-folder-open',
                        toggle: this.toggleFolder.bind(this),
                        score: _.keys(sections).length,
                    };
                }
                // add section
                sections[hit.folderId || 0] = section;
            }
            section.expanded = true;
            section.items.push(hit);
        }
    };
    SearchSrv.prototype.toggleFolder = function (section) {
        section.expanded = !section.expanded;
        section.icon = section.expanded ? 'fa fa-folder-open' : 'fa fa-folder';
        if (section.items.length) {
            return Promise.resolve(section);
        }
        var query = {
            folderIds: [section.id],
        };
        return this.backendSrv.search(query).then(function (results) {
            section.items = results;
            return Promise.resolve(section);
        });
    };
    SearchSrv.prototype.getDashboardTags = function () {
        return this.backendSrv.get('/api/dashboards/tags');
    };
    return SearchSrv;
}());
export { SearchSrv };
coreModule.service('searchSrv', SearchSrv);
//# sourceMappingURL=search_srv.js.map