export class BadgeRenderer {
  constructor() {}

  createHtml(badgesData, panel) {
    const dataToRender = this.manipulateData(badgesData, panel);
    if (dataToRender) {
      let retHtml = '';
      const keys = Object.keys(dataToRender);
      for (let i = 0; i < keys.length; i++) {
        const key = keys[i];
        const headerHTML = '<div class="head"><h2>' + key + '</h2></div>';
        let badgeInnerHTML = '<div class="content"><ul>';
        const data = dataToRender[key];
        for (let i = 0; i < data.length; i++) {
          badgeInnerHTML += '<li><p>' + data[i].label + '</p> <span>' + data[i].value + '</span></li>';
        }
        badgeInnerHTML += '</ul></div>';
        retHtml += '<div class="badge-conainer">' + headerHTML + badgeInnerHTML + '</div>';
      }
      return retHtml;
    } else {
      return '<div class="badge-conainer badge-data-loading query-transaction--loading">Your data is loading....</div>';
    }
  }

  manipulateData(data, panel) {
    if (data && data.length > 0) {
      const retData = [];
      for (let i = 0; i < data.length; i++) {
        const entry = data[i];
        let target = entry.target;
        let key = target;
        if (panel.keyRegex) {
          const keyRegex = new RegExp(panel.keyRegex, 'g');
          key = target.replace(keyRegex, '');
        }
        if (panel.labelRegex) {
          const regex = new RegExp(panel.labelRegex, 'g');
          target = target.replace(regex, '');
        }
        retData[key] = retData[key] || [];
        retData[key].push({
          label: target,
          value: entry.datapoints[0][0],
        });
      }
      return retData;
    }
    return null;
  }
}
