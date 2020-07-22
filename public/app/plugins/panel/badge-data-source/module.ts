import _ from 'lodash';
import { MetricsPanelCtrl } from 'app/plugins/sdk';
import { badgeEditor } from './editor';
import { BadgeRenderer } from './renderer';

class BadgeDataSourceCtrl extends MetricsPanelCtrl {
  static templateUrl = './partials/module.html';
  panelDefaults = {
    totalBadges: 1,
    badgesInfo: [],
    apiEndPoint: '',
  };
  data: null;
  isLoading: any;
  dummyData: any = [
    {
      title: 'Current month fee collection',
      data: [
        {
          label: 'Total',
          value: '52,08,222',
        },
        {
          label: 'Collected',
          value: '45,33,222',
        },
        {
          label: 'Pending',
          value: '4,33,222',
        },
        {
          label: 'Outstanding',
          value: '6,12,325',
        },
      ],
    },
    {
      title: 'Admission Details',
      data: [
        {
          label: 'Enquiry',
          value: '527',
        },
        {
          label: 'Application',
          value: '478',
        },
        {
          label: 'Admission',
          value: '300',
        },
      ],
    },
  ];
  constructor($scope, $injector) {
    super($scope, $injector);
    this.isLoading = false;
    _.defaults(this.panel, this.panelDefaults);
    this.events.on('data-received', this.onDataReceived.bind(this));
    this.events.on('init-edit-mode', this.onInitEditMode.bind(this));
  }

  onInitEditMode() {
    this.addEditorTab('Options', badgeEditor, 1);
  }

  onDataReceived(dataList: any) {
    this.data = dataList;
    this.data = this.dummyData;
    this.render();
  }

  render() {
    const badgeRenderer = new BadgeRenderer(this.panel);
    this.isLoading = false;
    const badgeHtml = badgeRenderer.createHtml(this.isLoading, this.data);
    return super.render(badgeHtml);
  }

  link(scope, elem, attrs, ctrl: BadgeDataSourceCtrl) {
    function renderPanel(renderData) {
      const parentElement = elem.find('.badges-main-container');
      parentElement.html(renderData);
    }

    ctrl.events.on('render', renderData => {
      if (renderData) {
        renderPanel(renderData);
      }
      ctrl.renderingCompleted();
    });
  }
}

export { BadgeDataSourceCtrl, BadgeDataSourceCtrl as PanelCtrl };