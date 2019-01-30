import { appEvents } from 'app/core/core';

export class LegalEntitiesCtrl {
  navModel: any;
  activeTabIndex = 0;
  logoSrc = '/public/img/logo.png';
  $scope;
  /** @ngInject */
  constructor($scope) {
    this.activeTabIndex = 0;
    this.$scope = $scope;
    $scope.getFile = this.getFile.bind(this);
  }

  activateTab(tabIndex) {
    this.activeTabIndex = tabIndex;
  }

  getFile(file) {
    if (!file) {
      return;
    }
    const fileReader = new FileReader();
    const that = this;
    fileReader.onloadend = e => {
      that.logoSrc = e.target['result'];
      this.$scope.$apply();
    };
    fileReader.readAsDataURL(file);
  }

  showSignatoryModal() {
    const text = 'Do you want to delete the ';

    appEvents.emit('signatory-modal', {
      text: text,
      icon: 'fa-trash',
    });
  }
  showBankModal() {
    const text = 'Do you want to delete the ';

    appEvents.emit('bank-modal', {
      text: text,
      icon: 'fa-trash',
    });
  }
}
