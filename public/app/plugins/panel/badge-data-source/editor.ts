import _ from 'lodash';

export class BadgeEditorCtrl {
  panel: any;
  panelCtrl: any;
  keyRegex: any;
  labelRegex: any;
  /** @ngInject */
  constructor($scope) {
    $scope.editor = this;
    this.panelCtrl = $scope.ctrl;
    this.panel = this.panelCtrl.panel;
    this.keyRegex = this.panel.keyRegex;
    this.labelRegex = this.panel.labelRegex;
  }

  onChangeKeyRegex() {
    this.panel.keyRegex = this.keyRegex;
  }

  onChangeLabelRegex() {
    this.panel.labelRegex = this.labelRegex;
  }

  render() {
    this.panelCtrl.render();
  }
}

/** @ngInject */
export function badgeEditor($q, uiSegmentSrv) {
  'use strict';
  return {
    restrict: 'E',
    scope: true,
    templateUrl: 'public/app/plugins/panel/badge-data-source/partials/editor.html',
    controller: BadgeEditorCtrl,
  };
}
