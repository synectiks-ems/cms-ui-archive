import React, { SFC } from 'react';

export interface Props {
  link: any;
  user: any;
}

const PluginItem: SFC<Props> = props => {
  const { link, user } = props;
  return (
    <div className="plugin-item" title={link.text}>
      <a className="sidemenu-link" href={link.url + '?signedInUser=' + user.login} target={link.target}>
        <span className="icon-circle sidemenu-icon">
          <i className={link.icon} />
          {link.img && <img src={link.img} />}
        </span>
        <div className="plugin-name">{link.text}</div>
      </a>
    </div>
  );
};

export default PluginItem;
