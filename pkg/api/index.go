package api

import (
	"bytes"
	"fmt"
	"github.com/synectiks-ems/ems-ui/pkg/components/simplejson"
	"github.com/synectiks-ems/ems-ui/pkg/log"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/synectiks-ems/ems-ui/pkg/api/dtos"
	"github.com/synectiks-ems/ems-ui/pkg/bus"
	m "github.com/synectiks-ems/ems-ui/pkg/models"
	"github.com/synectiks-ems/ems-ui/pkg/plugins"
	"github.com/synectiks-ems/ems-ui/pkg/setting"
)

const (
	// Themes
	lightName = "light"
	darkName  = "dark"
)

var isUiModulesExported = false

//var moduleMap = make(map[string]string)
//var delim = "~~~"

var externalSecurityServiceClient = &http.Client{
	Transport: &http.Transport{Proxy: http.ProxyFromEnvironment},
}

func (hs *HTTPServer) setIndexViewDataForRbacUser(externalUserId string, externalUserPw string, c *m.ReqContext) (*dtos.IndexViewData, error) {
	//response, err := externalSecurityServiceClient.Get(setting.ExternalSecurityUrl + "/security/public/login?username=" + externalUserId + "&password=" + externalUserPw)

	//response, err := externalSecurityServiceClient.Get(setting.CmsUrl + "/api/cmslogin?username=" + externalUserId + "&password=" + externalUserPw)
	//if err != nil {
	//	return nil, err
	//}
	//defer response.Body.Close()
	//bodyBytes, err := ioutil.ReadAll(response.Body)
	//
	//if err != nil {
	//	return nil, err
	//}
	//if response.StatusCode == 417 {
	//	return nil, err
	//}
	//bodyString := string(bodyBytes)
	//fmt.Println(bodyString)

	//var userInfo map[string]interface{}
	//errw := json.Unmarshal([]byte(bodyString), &userInfo)
	//if errw != nil {
	//	return nil, errw
	//}

	log.Info("Signed in user : " + c.SignedInUser.Name)
	var userInfo = c.Session.Get(c.SignedInUser.Name).(map[string]interface{})
	if userInfo == nil {
		return nil, nil
	}
	fmt.Println(userInfo)

	settings, err := hs.getFrontendSettingsMap(c)

	if err != nil {
		return nil, err
	}
	prefsQuery := m.GetPreferencesWithDefaultsQuery{User: c.SignedInUser}
	if err := bus.Dispatch(&prefsQuery); err != nil {
		return nil, err
	}
	prefs := prefsQuery.Result

	// Read locale from acccept-language
	acceptLang := c.Req.Header.Get("Accept-Language")
	locale := "en-US"

	if len(acceptLang) > 0 {
		parts := strings.Split(acceptLang, ",")
		locale = parts[0]
	}

	appURL := setting.AppUrl
	appSubURL := setting.AppSubUrl

	// special case when doing localhost call from phantomjs
	if c.IsRenderCall {
		appURL = fmt.Sprintf("%s://localhost:%s", setting.Protocol, setting.HttpPort)
		appSubURL = ""
		settings["appSubUrl"] = ""
	}

	//hasEditPermissionInFoldersQuery := m.HasEditPermissionInFoldersQuery{SignedInUser: c.SignedInUser}
	//if err := bus.Dispatch(&hasEditPermissionInFoldersQuery); err != nil {
	//	return nil, err
	//}

	var data = dtos.IndexViewData{
		User: &dtos.CurrentUser{
			Id:                         c.UserId,
			IsSignedIn:                 c.IsSignedIn,
			Login:                      c.Login,
			Email:                      c.Email,
			Name:                       c.Name,
			OrgCount:                   c.OrgCount,
			OrgId:                      c.OrgId,
			OrgName:                    c.OrgName,
			OrgRole:                    c.OrgRole,
			GravatarUrl:                dtos.GetGravatarUrl(c.Email),
			IsGrafanaAdmin:             c.IsGrafanaAdmin,
			LightTheme:                 true,
			Timezone:                   prefs.Timezone,
			Locale:                     locale,
			HelpFlags1:                 c.HelpFlags1,
			HasEditPermissionInFolders: false, //hasEditPermissionInFoldersQuery.Result,
		},
		Settings:  settings,
		Theme:     prefs.Theme,
		AppUrl:    appURL,
		AppSubUrl: appSubURL,
		//GoogleAnalyticsId:       setting.GoogleAnalyticsId,
		//GoogleTagManagerId:      setting.GoogleTagManagerId,
		//BuildVersion:            setting.BuildVersion,
		//BuildCommit:             setting.BuildCommit,
		//NewGrafanaVersion:       plugins.GrafanaLatestVersion,
		//NewGrafanaVersionExists: plugins.GrafanaHasUpdate,
		AppName:          setting.ApplicationName,
		AppNameBodyClass: getAppNameBodyClass(setting.ApplicationName),
	}

	if setting.DisableGravatar {
		data.User.GravatarUrl = setting.AppSubUrl + "/public/img/user_profile.png"
	}

	if len(data.User.Name) == 0 {
		data.User.Name = data.User.Login
	}

	themeURLParam := c.Query("theme")
	if themeURLParam == lightName {
		data.User.LightTheme = true
		data.Theme = lightName
	} else if themeURLParam == darkName {
		data.User.LightTheme = false
		data.Theme = darkName
	}

	//if hasEditPermissionInFoldersQuery.Result {
	//	children := []*dtos.NavLink{
	//		{Text: "Dashboard", Icon: "gicon gicon-dashboard-new", Url: setting.AppSubUrl + "/dashboard/new"},
	//	}
	//
	//	if c.OrgRole == m.ROLE_ADMIN || c.OrgRole == m.ROLE_EDITOR {
	//		children = append(children, &dtos.NavLink{Text: "Folder", SubTitle: "Create a new folder to organize your dashboards", Id: "folder", Icon: "gicon gicon-folder-new", Url: setting.AppSubUrl + "/dashboards/folder/new"})
	//	}
	//
	//	children = append(children, &dtos.NavLink{Text: "Import", SubTitle: "Import dashboard from file or Grafana.com", Id: "import", Icon: "gicon gicon-dashboard-import", Url: setting.AppSubUrl + "/dashboard/import"})
	//
	//	data.NavTree = append(data.NavTree, &dtos.NavLink{
	//		Text:     "Create",
	//		Id:       "create",
	//		Icon:     "fa fa-fw fa-plus",
	//		Url:      setting.AppSubUrl + "/dashboard/new",
	//		Children: children,
	//	})
	//}

	dashboardChildNavs := []*dtos.NavLink{
		{Text: "Home", Id: "home", Url: setting.AppSubUrl + "/", Icon: "gicon gicon-home", HideFromTabs: true},
		//{Text: "Divider", Divider: true, Id: "divider", HideFromTabs: true},
		//{Text: "Manage", Id: "manage-dashboards", Url: setting.AppSubUrl + "/dashboards", Icon: "gicon gicon-manage"},
		//{Text: "Playlists", Id: "playlists", Url: setting.AppSubUrl + "/playlists", Icon: "gicon gicon-playlists"},
		//{Text: "Snapshots", Id: "snapshots", Url: setting.AppSubUrl + "/dashboard/snapshots", Icon: "gicon gicon-snapshots"},
	}

	data.NavTree = append(data.NavTree, &dtos.NavLink{
		Text:     "Dashboards",
		Id:       "dashboards",
		SubTitle: "Manage dashboards & folders",
		Icon:     "gicon gicon-dashboard",
		Url:      setting.AppSubUrl + "/",
		Children: dashboardChildNavs,
	})

	//if setting.ExploreEnabled && (c.OrgRole == m.ROLE_ADMIN || c.OrgRole == m.ROLE_EDITOR) {
	//	data.NavTree = append(data.NavTree, &dtos.NavLink{
	//		Text:     "Explore",
	//		Id:       "explore",
	//		SubTitle: "Explore your data",
	//		Icon:     "fa fa-rocket",
	//		Url:      setting.AppSubUrl + "/explore",
	//	})
	//}

	if c.IsSignedIn {
		// Only set login if it's different from the name
		var login string
		if c.SignedInUser.Login != c.SignedInUser.NameOrFallback() {
			login = c.SignedInUser.Login
		}
		profileNode := &dtos.NavLink{
			Text:         c.SignedInUser.NameOrFallback(),
			SubTitle:     login,
			Id:           "profile",
			Img:          data.User.GravatarUrl,
			Url:          setting.AppSubUrl + "/profile",
			HideFromMenu: true,
			Children: []*dtos.NavLink{
				{Text: "Preferences", Id: "profile-settings", Url: setting.AppSubUrl + "/profile", Icon: "gicon gicon-preferences"},
				{Text: "Change Password", Id: "change-password", Url: setting.AppSubUrl + "/profile/password", Icon: "fa fa-fw fa-lock", HideFromMenu: true},
			},
		}

		if !setting.DisableSignoutMenu {
			// add sign out first
			profileNode.Children = append(profileNode.Children, &dtos.NavLink{
				Text: "Sign out", Id: "sign-out", Url: setting.AppSubUrl + "/logout", Icon: "fa fa-fw fa-sign-out", Target: "_self",
			})
		}

		data.NavTree = append(data.NavTree, profileNode)
	}

	enabledPlugins, err := plugins.GetEnabledPlugins(c.OrgId)
	if err != nil {
		return nil, err
	}
	var auth = userInfo["authz"].(map[string]interface{})
	//var permissions = fmt.Sprint(auth["permissions"])
	var roles = auth["roles"]
	c.Session.Set("myuserrole", roles)
	//fmt.Println(permissions)
	//fmt.Println(roles)
	//var replacer = strings.NewReplacer("[", "", "]", "")
	//permissions = replacer.Replace(permissions)
	//var permissionsAry = strings.Split(permissions, " ")
	var permissions = auth["permissions"].(interface{})
	keyMap := make(map[string]bool)
	list := []string{}
	for _, val := range permissions.([]interface{}) {
		var entry = fmt.Sprint(val)
		if _, value := keyMap[entry]; !value {
			keyMap[entry] = value
			list = append(list, entry)
		}
	}

	pluginResponse, err := externalSecurityServiceClient.Get(setting.CmsUrl + "/api/cmsparentmodules")
	if err != nil {
		return nil, err
	}
	defer pluginResponse.Body.Close()
	plgBytes, err := ioutil.ReadAll(pluginResponse.Body)
	if err != nil {
		return nil, err
	}
	pluginString := string(plgBytes)
	var replacer = strings.NewReplacer("\"", "", "[", "", "[ ", "", " [", "", " [ ", "", "]", "", "] ", "", " ]", "", " ] ", "")
	pluginString = replacer.Replace(pluginString)
	var pluginAry = strings.Split(pluginString, ",")
	pluginMap := make(map[string]interface{})
	for _, val := range pluginAry {
		pluginMap[strings.TrimSpace(val)] = strings.TrimSpace(val)
	}
	fmt.Println("plugins from cms : ", pluginMap)

	permissions = list
	fmt.Println(permissions)
	for _, plugin := range enabledPlugins.Apps {
		var plgNm = fmt.Sprint(pluginMap[plugin.Name])
		if plgNm == plugin.Name { //plugin.Pinned {
			appLink := &dtos.NavLink{
				Text: plugin.Name,
				Id:   "plugin-page-" + plugin.Id,
				Url:  plugin.DefaultNavUrl,
				Img:  plugin.Info.Logos.Small,
			}
			for _, val := range permissions.([]string) {
				var pr = fmt.Sprint(val)
				log.Debug("permission from security service : " + pr)
				for _, include := range plugin.Includes {
					if include.Type == "page" && include.AddToNav {
						var plgUrl = setting.AppSubUrl + "/plugins/" + plugin.Id + "/page/" + include.Slug
						log.Debug("Plugin url from grafana plugin : " + plgUrl)
						if include.Name == pr {
							link := &dtos.NavLink{
								Url:  plgUrl,
								Text: include.Name,
							}
							appLink.Children = append(appLink.Children, link)
						}
					}
				}
			}
			if len(appLink.Children) > 0 {
				data.NavTree = append(data.NavTree, appLink)
			}
		}
	}

	//for _, plugin := range enabledPlugins.Apps {
	//	if plugin.Pinned {
	//		appLink := &dtos.NavLink{
	//			Text: plugin.Name,
	//			Id:   "plugin-page-" + plugin.Id,
	//			Url:  plugin.DefaultNavUrl,
	//			Img:  plugin.Info.Logos.Small,
	//		}
	//
	//		for _, include := range plugin.Includes {
	//			if !c.HasUserRole(include.Role) {
	//				continue
	//			}
	//
	//			if include.Type == "page" && include.AddToNav {
	//				link := &dtos.NavLink{
	//					Url:  setting.AppSubUrl + "/plugins/" + plugin.Id + "/page/" + include.Slug,
	//					Text: include.Name,
	//				}
	//				appLink.Children = append(appLink.Children, link)
	//			}
	//
	//			//if include.Type == "dashboard" && include.AddToNav {
	//			//	link := &dtos.NavLink{
	//			//		Url:  setting.AppSubUrl + "/dashboard/db/" + include.Slug,
	//			//		Text: include.Name,
	//			//	}
	//			//	appLink.Children = append(appLink.Children, link)
	//			//}
	//		}
	//
	//		//if len(appLink.Children) > 0 && c.OrgRole == m.ROLE_ADMIN {
	//		//	appLink.Children = append(appLink.Children, &dtos.NavLink{Divider: true})
	//		//	appLink.Children = append(appLink.Children, &dtos.NavLink{Text: "Plugin Config", Icon: "gicon gicon-cog", Url: setting.AppSubUrl + "/plugins/" + plugin.Id + "/edit"})
	//		//}
	//
	//		if len(appLink.Children) > 0 {
	//			data.NavTree = append(data.NavTree, appLink)
	//		}
	//	}
	//}

	//if setting.AlertingEnabled && (c.OrgRole == m.ROLE_ADMIN || c.OrgRole == m.ROLE_EDITOR) {
	//	alertChildNavs := []*dtos.NavLink{
	//		{Text: "Alert Rules", Id: "alert-list", Url: setting.AppSubUrl + "/alerting/list", Icon: "gicon gicon-alert-rules"},
	//		{Text: "Notification channels", Id: "channels", Url: setting.AppSubUrl + "/alerting/notifications", Icon: "gicon gicon-alert-notification-channel"},
	//	}
	//
	//	data.NavTree = append(data.NavTree, &dtos.NavLink{
	//		Text:     "Alerting",
	//		SubTitle: "Alert rules & notifications",
	//		Id:       "alerting",
	//		Icon:     "gicon gicon-alert",
	//		Url:      setting.AppSubUrl + "/alerting/list",
	//		Children: alertChildNavs,
	//	})
	//}
	//if c.IsGrafanaAdmin || c.OrgRole == m.ROLE_ADMIN {
	log.Debug("Checking localconfig to be added or not")
	for _, val := range permissions.([]string) {
		var pr = fmt.Sprint(val)
		if pr == "Local Config" && fmt.Sprint(pluginMap["Configuration"]) == "Configuration" {
			cfgNode := &dtos.NavLink{
				Id:       "cfg",
				Text:     "Configuration",
				SubTitle: "Applications",
				Icon:     "gicon gicon-cog",
				Url:      setting.AppSubUrl + "/localapp",
				Children: []*dtos.NavLink{
					{
						Text:         "Local Config",
						Icon:         "gicon gicon-datasources",
						Description:  "Local Application",
						Id:           "localapp",
						Url:          setting.AppSubUrl + "/localapp",
						HideFromTabs: true,
					},
					//{
					//	Text:         "Global Config",
					//	Icon:         "gicon gicon-datasources",
					//	Description:  "Add and configure data sources",
					//	Id:           "datasources",
					//	Url:          setting.AppSubUrl + "/datasources",
					//	HideFromTabs: true,
					//},
					{
						Text:         "Data Sources",
						Icon:         "gicon gicon-datasources",
						Description:  "Add and configure data sources",
						Id:           "datasources",
						Url:          setting.AppSubUrl + "/datasources",
						HideFromMenu: true,
					},
					{
						Text:         "Users",
						Id:           "users",
						Description:  "Manage org members",
						Icon:         "gicon gicon-user",
						Url:          setting.AppSubUrl + "/org/users",
						HideFromMenu: true,
						HideFromTabs: true,
					},
					{
						Text:         "Teams",
						Id:           "teams",
						Description:  "Manage org groups",
						Icon:         "gicon gicon-team",
						Url:          setting.AppSubUrl + "/org/teams",
						HideFromMenu: true,
						HideFromTabs: true,
					},
					{
						Text:         "Plugins",
						Id:           "plugins",
						Description:  "View and configure plugins",
						Icon:         "gicon gicon-plugins",
						Url:          setting.AppSubUrl + "/plugins",
						HideFromMenu: true,
					},
					{
						Text:         "Preferences",
						Id:           "org-settings",
						Description:  "Organization preferences",
						Icon:         "gicon gicon-preferences",
						Url:          setting.AppSubUrl + "/org",
						HideFromMenu: true,
					},

					{
						Text:         "API Keys",
						Id:           "apikeys",
						Description:  "Create & manage API keys",
						Icon:         "gicon gicon-apikeys",
						Url:          setting.AppSubUrl + "/org/apikeys",
						HideFromMenu: true,
					},
				},
			}

			if c.OrgRole != m.ROLE_ADMIN {
				cfgNode = &dtos.NavLink{
					Id:       "cfg",
					Text:     "Configuration",
					SubTitle: "Organization: " + c.OrgName,
					Icon:     "gicon gicon-cog",
					Url:      setting.AppSubUrl + "/admin/users",
					Children: make([]*dtos.NavLink, 0),
				}
			}

			if c.OrgRole == m.ROLE_ADMIN && c.IsGrafanaAdmin {
				cfgNode.Children = append(cfgNode.Children, &dtos.NavLink{
					Divider: true, HideFromTabs: true, Id: "admin-divider", Text: "Text",
				})
			}

			//if c.IsGrafanaAdmin {
			//	cfgNode.Children = append(cfgNode.Children, &dtos.NavLink{
			//		Text:         "Server Admin",
			//		HideFromTabs: true,
			//		SubTitle:     "Manage all users & orgs",
			//		Id:           "admin",
			//		Icon:         "gicon gicon-shield",
			//		Url:          setting.AppSubUrl + "/admin/users",
			//		Children: []*dtos.NavLink{
			//			{Text: "Users", Id: "global-users", Url: setting.AppSubUrl + "/admin/users", Icon: "gicon gicon-user"},
			//			{Text: "Orgs", Id: "global-orgs", Url: setting.AppSubUrl + "/admin/orgs", Icon: "gicon gicon-org"},
			//			{Text: "Settings", Id: "server-settings", Url: setting.AppSubUrl + "/admin/settings", Icon: "gicon gicon-preferences"},
			//			{Text: "Stats", Id: "server-stats", Url: setting.AppSubUrl + "/admin/stats", Icon: "fa fa-fw fa-bar-chart"},
			//			{Text: "Style Guide", Id: "styleguide", Url: setting.AppSubUrl + "/styleguide", Icon: "fa fa-fw fa-eyedropper"},
			//		},
			//	})
			//}

			data.NavTree = append(data.NavTree, cfgNode)
		}
	}

	//}

	//data.NavTree = append(data.NavTree, &dtos.NavLink{
	//	Text:         "Help",
	//	SubTitle:     fmt.Sprintf(`%s v%s (%s)`, setting.ApplicationName, setting.BuildVersion, setting.BuildCommit),
	//	Id:           "help",
	//	Url:          "#",
	//	Icon:         "gicon gicon-question",
	//	HideFromMenu: true,
	//	Children: []*dtos.NavLink{
	//		{Text: "Keyboard shortcuts", Url: "/shortcuts", Icon: "fa fa-fw fa-keyboard-o", Target: "_self"},
	//		{Text: "Community site", Url: "http://community.grafana.com", Icon: "fa fa-fw fa-comment", Target: "_blank"},
	//		{Text: "Documentation", Url: "http://docs.grafana.org", Icon: "fa fa-fw fa-file", Target: "_blank"},
	//	},
	//})

	hs.HooksService.RunIndexDataHooks(&data)
	return &data, nil
}

func (hs *HTTPServer) pushModulesToCms(c *m.ReqContext) {
	type moduleArray []interface{}
	enabledPlugins, err := plugins.GetEnabledPlugins(c.OrgId)
	if err != nil {
		log.Error(500, "Failed to get all the enabled plugins", err)
	}
	var mary moduleArray
	for _, plugin := range enabledPlugins.Apps {
		if plugin.Pinned {
			for _, include := range plugin.Includes {
				if include.Type == "page" && include.AddToNav {
					message := map[string]interface{}{
						"moduleName":    plugin.Name,
						"subModuleName": include.Name,
						"url":           setting.AppSubUrl + "/plugins/" + plugin.Id + "/page/" + include.Slug,
						"status":        "ACTIVE",
					}
					mary = append(mary, message)
					//moduleMap[plugin.Name+delim+include.Name] = setting.AppSubUrl + "/plugins/" + plugin.Id + "/page/" + include.Slug
				}
			}
		}
	}

	message := map[string]interface{}{
		"moduleName":    "Configuration",
		"subModuleName": "Local Config",
		"url":           setting.AppSubUrl + "/localapp",
		"status":        "ACTIVE",
	}
	mary = append(mary, message)
	//moduleMap["Configuration"+delim+"Local Config"] = "/localapp"
	messageBytes, err := simplejson.NewFromAny(mary).Encode()
	if err != nil {
		log.Error(500, "Failed to encode message array in json", err)
	}
	httpClient := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, setting.CmsUrl+"/api/cmsmodules", bytes.NewBuffer(messageBytes))
	req.Header.Add("Content-Type", "application/json")
	httpClient.Do(req)
}

// setIndexViewDataForRbacCmsAdminUser() method to create index view for application admin user
func (hs *HTTPServer) setIndexViewDataForRbacCmsAdminUser(c *m.ReqContext) (*dtos.IndexViewData, error) {
	if !isUiModulesExported {
		hs.pushModulesToCms(c)
		isUiModulesExported = true
	}

	settings, err := hs.getFrontendSettingsMap(c)
	if err != nil {
		return nil, err
	}

	prefsQuery := m.GetPreferencesWithDefaultsQuery{User: c.SignedInUser}
	if err := bus.Dispatch(&prefsQuery); err != nil {
		return nil, err
	}
	prefs := prefsQuery.Result

	// Read locale from acccept-language
	acceptLang := c.Req.Header.Get("Accept-Language")
	locale := "en-US"

	if len(acceptLang) > 0 {
		parts := strings.Split(acceptLang, ",")
		locale = parts[0]
	}

	appURL := setting.AppUrl
	appSubURL := setting.AppSubUrl

	// special case when doing localhost call from phantomjs
	if c.IsRenderCall {
		appURL = fmt.Sprintf("%s://localhost:%s", setting.Protocol, setting.HttpPort)
		appSubURL = ""
		settings["appSubUrl"] = ""
	}

	hasEditPermissionInFoldersQuery := m.HasEditPermissionInFoldersQuery{SignedInUser: c.SignedInUser}
	if err := bus.Dispatch(&hasEditPermissionInFoldersQuery); err != nil {
		return nil, err
	}

	var data = dtos.IndexViewData{
		User: &dtos.CurrentUser{
			Id:                         c.UserId,
			IsSignedIn:                 c.IsSignedIn,
			Login:                      c.Login,
			Email:                      c.Email,
			Name:                       c.Name,
			OrgCount:                   c.OrgCount,
			OrgId:                      c.OrgId,
			OrgName:                    c.OrgName,
			OrgRole:                    c.OrgRole,
			GravatarUrl:                dtos.GetGravatarUrl(c.Email),
			IsGrafanaAdmin:             c.IsGrafanaAdmin,
			LightTheme:                 prefs.Theme == lightName,
			Timezone:                   prefs.Timezone,
			Locale:                     locale,
			HelpFlags1:                 c.HelpFlags1,
			HasEditPermissionInFolders: hasEditPermissionInFoldersQuery.Result,
		},
		Settings:                settings,
		Theme:                   prefs.Theme,
		AppUrl:                  appURL,
		AppSubUrl:               appSubURL,
		GoogleAnalyticsId:       setting.GoogleAnalyticsId,
		GoogleTagManagerId:      setting.GoogleTagManagerId,
		BuildVersion:            setting.BuildVersion,
		BuildCommit:             setting.BuildCommit,
		NewGrafanaVersion:       plugins.GrafanaLatestVersion,
		NewGrafanaVersionExists: plugins.GrafanaHasUpdate,
		AppName:                 setting.ApplicationName,
		AppNameBodyClass:        getAppNameBodyClass(setting.ApplicationName),
	}

	if setting.DisableGravatar {
		data.User.GravatarUrl = setting.AppSubUrl + "/public/img/user_profile.png"
	}

	if len(data.User.Name) == 0 {
		data.User.Name = data.User.Login
	}

	themeURLParam := c.Query("theme")
	if themeURLParam == lightName {
		data.User.LightTheme = true
		data.Theme = lightName
	} else if themeURLParam == darkName {
		data.User.LightTheme = false
		data.Theme = darkName
	}

	//if hasEditPermissionInFoldersQuery.Result {
	//	children := []*dtos.NavLink{
	//		{Text: "Dashboard", Icon: "gicon gicon-dashboard-new", Url: setting.AppSubUrl + "/dashboard/new"},
	//	}
	//
	//	if c.OrgRole == m.ROLE_ADMIN || c.OrgRole == m.ROLE_EDITOR {
	//		children = append(children, &dtos.NavLink{Text: "Folder", SubTitle: "Create a new folder to organize your dashboards", Id: "folder", Icon: "gicon gicon-folder-new", Url: setting.AppSubUrl + "/dashboards/folder/new"})
	//	}
	//
	//	children = append(children, &dtos.NavLink{Text: "Import", SubTitle: "Import dashboard from file or Grafana.com", Id: "import", Icon: "gicon gicon-dashboard-import", Url: setting.AppSubUrl + "/dashboard/import"})
	//
	//	data.NavTree = append(data.NavTree, &dtos.NavLink{
	//		Text:     "Create",
	//		Id:       "create",
	//		Icon:     "fa fa-fw fa-plus",
	//		Url:      setting.AppSubUrl + "/dashboard/new",
	//		Children: children,
	//	})
	//}

	dashboardChildNavs := []*dtos.NavLink{
		{Text: "Home", Id: "home", Url: setting.AppSubUrl + "/", Icon: "gicon gicon-home", HideFromTabs: true},
		//{Text: "Divider", Divider: true, Id: "divider", HideFromTabs: true},
		//{Text: "Manage", Id: "manage-dashboards", Url: setting.AppSubUrl + "/dashboards", Icon: "gicon gicon-manage"},
		//{Text: "Playlists", Id: "playlists", Url: setting.AppSubUrl + "/playlists", Icon: "gicon gicon-playlists"},
		//{Text: "Snapshots", Id: "snapshots", Url: setting.AppSubUrl + "/dashboard/snapshots", Icon: "gicon gicon-snapshots"},
	}

	data.NavTree = append(data.NavTree, &dtos.NavLink{
		Text:     "Dashboards",
		Id:       "dashboards",
		SubTitle: "Manage dashboards & folders",
		Icon:     "gicon gicon-dashboard",
		Url:      setting.AppSubUrl + "/",
		Children: dashboardChildNavs,
	})

	//if setting.ExploreEnabled && (c.OrgRole == m.ROLE_ADMIN || c.OrgRole == m.ROLE_EDITOR) {
	//	data.NavTree = append(data.NavTree, &dtos.NavLink{
	//		Text:     "Explore",
	//		Id:       "explore",
	//		SubTitle: "Explore your data",
	//		Icon:     "fa fa-rocket",
	//		Url:      setting.AppSubUrl + "/explore",
	//	})
	//}

	if c.IsSignedIn {
		// Only set login if it's different from the name
		var login string
		if c.SignedInUser.Login != c.SignedInUser.NameOrFallback() {
			login = c.SignedInUser.Login
		}
		profileNode := &dtos.NavLink{
			Text:         c.SignedInUser.NameOrFallback(),
			SubTitle:     login,
			Id:           "profile",
			Img:          data.User.GravatarUrl,
			Url:          setting.AppSubUrl + "/profile",
			HideFromMenu: true,
			//Children: []*dtos.NavLink{
			//	{Text: "Preferences", Id: "profile-settings", Url: setting.AppSubUrl + "/profile", Icon: "gicon gicon-preferences"},
			//	{Text: "Change Password", Id: "change-password", Url: setting.AppSubUrl + "/profile/password", Icon: "fa fa-fw fa-lock", HideFromMenu: true},
			//},
		}

		if !setting.DisableSignoutMenu {
			// add sign out first
			profileNode.Children = append(profileNode.Children, &dtos.NavLink{
				Text: "Sign out", Id: "sign-out", Url: setting.AppSubUrl + "/logout", Icon: "fa fa-fw fa-sign-out", Target: "_self",
			})
		}

		data.NavTree = append(data.NavTree, profileNode)
	}

	enabledPlugins, err := plugins.GetEnabledPlugins(c.OrgId)
	if err != nil {
		return nil, err
	}
	pluginResponse, err := externalSecurityServiceClient.Get(setting.CmsUrl + "/api/cmsparentmodules")
	if err != nil {
		return nil, err
	}
	defer pluginResponse.Body.Close()
	plgBytes, err := ioutil.ReadAll(pluginResponse.Body)
	if err != nil {
		return nil, err
	}
	pluginString := string(plgBytes)
	var replacer = strings.NewReplacer("\"", "", "[", "", "[ ", "", " [", "", " [ ", "", "]", "", "] ", "", " ]", "", " ] ", "")
	pluginString = replacer.Replace(pluginString)
	var pluginAry = strings.Split(pluginString, ",")
	pluginMap := make(map[string]interface{})
	for _, val := range pluginAry {
		pluginMap[strings.TrimSpace(val)] = strings.TrimSpace(val)
	}
	fmt.Println("plugins from cms : ", pluginMap)

	//permissions = list
	//fmt.Println(permissions)
	for _, plugin := range enabledPlugins.Apps {
		var plgNm = fmt.Sprint(pluginMap[plugin.Name])
		if plgNm == plugin.Name { //plugin.Pinned {
			appLink := &dtos.NavLink{
				Text: plugin.Name,
				Id:   "plugin-page-" + plugin.Id,
				Url:  plugin.DefaultNavUrl,
				Img:  plugin.Info.Logos.Small,
			}
			//for _, val := range permissions.([]string) {
			//	var pr = fmt.Sprint(val)
			//	log.Debug("permission from security service : " + pr)
			for _, include := range plugin.Includes {
				if include.Type == "page" && include.AddToNav {
					var plgUrl = setting.AppSubUrl + "/plugins/" + plugin.Id + "/page/" + include.Slug
					log.Debug("Plugin url from grafana plugin : " + plgUrl)
					//if include.Name == pr {
					link := &dtos.NavLink{
						Url:  plgUrl,
						Text: include.Name,
					}
					appLink.Children = append(appLink.Children, link)
					//}
				}
			}
			//}
			if len(appLink.Children) > 0 {
				data.NavTree = append(data.NavTree, appLink)
			}
		}
	}

	//if setting.AlertingEnabled && (c.OrgRole == m.ROLE_ADMIN || c.OrgRole == m.ROLE_EDITOR) {
	//	alertChildNavs := []*dtos.NavLink{
	//		{Text: "Alert Rules", Id: "alert-list", Url: setting.AppSubUrl + "/alerting/list", Icon: "gicon gicon-alert-rules"},
	//		{Text: "Notification channels", Id: "channels", Url: setting.AppSubUrl + "/alerting/notifications", Icon: "gicon gicon-alert-notification-channel"},
	//	}
	//
	//	data.NavTree = append(data.NavTree, &dtos.NavLink{
	//		Text:     "Alerting",
	//		SubTitle: "Alert rules & notifications",
	//		Id:       "alerting",
	//		Icon:     "gicon gicon-alert",
	//		Url:      setting.AppSubUrl + "/alerting/list",
	//		Children: alertChildNavs,
	//	})
	//}
	//if c.IsGrafanaAdmin || c.OrgRole == m.ROLE_ADMIN {

	cfgNode := &dtos.NavLink{
		Id:       "cfg",
		Text:     "Configuration",
		SubTitle: "Applications",
		Icon:     "gicon gicon-cog",
		Url:      setting.AppSubUrl + "/localapp",
		Children: []*dtos.NavLink{
			{
				Text:         "Local Config",
				Icon:         "gicon gicon-datasources",
				Description:  "Local Application",
				Id:           "localapp",
				Url:          setting.AppSubUrl + "/localapp",
				HideFromTabs: true,
			},
			//{
			//	Text:         "Global Config",
			//	Icon:         "gicon gicon-datasources",
			//	Description:  "Add and configure data sources",
			//	Id:           "datasources",
			//	Url:          setting.AppSubUrl + "/datasources",
			//	HideFromTabs: true,
			//},
			//{
			//	Text:         "Data Sources",
			//	Icon:         "gicon gicon-datasources",
			//	Description:  "Add and configure data sources",
			//	Id:           "datasources",
			//	Url:          setting.AppSubUrl + "/datasources",
			//	HideFromMenu: true,
			//},
			//{
			//	Text:         "Users",
			//	Id:           "users",
			//	Description:  "Manage org members",
			//	Icon:         "gicon gicon-user",
			//	Url:          setting.AppSubUrl + "/org/users",
			//	HideFromMenu: true,
			//	HideFromTabs: true,
			//},
			//{
			//	Text:         "Teams",
			//	Id:           "teams",
			//	Description:  "Manage org groups",
			//	Icon:         "gicon gicon-team",
			//	Url:          setting.AppSubUrl + "/org/teams",
			//	HideFromMenu: true,
			//	HideFromTabs: true,
			//},
			//{
			//	Text:         "Plugins",
			//	Id:           "plugins",
			//	Description:  "View and configure plugins",
			//	Icon:         "gicon gicon-plugins",
			//	Url:          setting.AppSubUrl + "/plugins",
			//	HideFromMenu: true,
			//},
			//{
			//	Text:         "Preferences",
			//	Id:           "org-settings",
			//	Description:  "Organization preferences",
			//	Icon:         "gicon gicon-preferences",
			//	Url:          setting.AppSubUrl + "/org",
			//	HideFromMenu: true,
			//},

			//{
			//	Text:         "API Keys",
			//	Id:           "apikeys",
			//	Description:  "Create & manage API keys",
			//	Icon:         "gicon gicon-apikeys",
			//	Url:          setting.AppSubUrl + "/org/apikeys",
			//	HideFromMenu: true,
			//},
		},
	}

	//if c.OrgRole != m.ROLE_ADMIN {
	//	cfgNode = &dtos.NavLink{
	//		Id:       "cfg",
	//		Text:     "Configuration",
	//		SubTitle: "Organization: " + c.OrgName,
	//		Icon:     "gicon gicon-cog",
	//		Url:      setting.AppSubUrl + "/admin/users",
	//		Children: make([]*dtos.NavLink, 0),
	//	}
	//}

	//if c.OrgRole == m.ROLE_ADMIN && c.IsGrafanaAdmin {
	//	cfgNode.Children = append(cfgNode.Children, &dtos.NavLink{
	//		Divider: true, HideFromTabs: true, Id: "admin-divider", Text: "Text",
	//	})
	//}

	//if c.IsGrafanaAdmin {
	//	cfgNode.Children = append(cfgNode.Children, &dtos.NavLink{
	//		Text:         "Server Admin",
	//		HideFromTabs: true,
	//		SubTitle:     "Manage all users & orgs",
	//		Id:           "admin",
	//		Icon:         "gicon gicon-shield",
	//		Url:          setting.AppSubUrl + "/admin/users",
	//		Children: []*dtos.NavLink{
	//			{Text: "Users", Id: "global-users", Url: setting.AppSubUrl + "/admin/users", Icon: "gicon gicon-user"},
	//			{Text: "Orgs", Id: "global-orgs", Url: setting.AppSubUrl + "/admin/orgs", Icon: "gicon gicon-org"},
	//			{Text: "Settings", Id: "server-settings", Url: setting.AppSubUrl + "/admin/settings", Icon: "gicon gicon-preferences"},
	//			{Text: "Stats", Id: "server-stats", Url: setting.AppSubUrl + "/admin/stats", Icon: "fa fa-fw fa-bar-chart"},
	//			{Text: "Style Guide", Id: "styleguide", Url: setting.AppSubUrl + "/styleguide", Icon: "fa fa-fw fa-eyedropper"},
	//		},
	//	})
	//}

	data.NavTree = append(data.NavTree, cfgNode)
	//}

	//data.NavTree = append(data.NavTree, &dtos.NavLink{
	//	Text:         "Help",
	//	SubTitle:     fmt.Sprintf(`%s v%s (%s)`, setting.ApplicationName, setting.BuildVersion, setting.BuildCommit),
	//	Id:           "help",
	//	Url:          "#",
	//	Icon:         "gicon gicon-question",
	//	HideFromMenu: true,
	//	Children: []*dtos.NavLink{
	//		{Text: "Keyboard shortcuts", Url: "/shortcuts", Icon: "fa fa-fw fa-keyboard-o", Target: "_self"},
	//		{Text: "Community site", Url: "http://community.grafana.com", Icon: "fa fa-fw fa-comment", Target: "_blank"},
	//		{Text: "Documentation", Url: "http://docs.grafana.org", Icon: "fa fa-fw fa-file", Target: "_blank"},
	//	},
	//})
	hs.HooksService.RunIndexDataHooks(&data)
	return &data, nil
}

//////////////////////////
func (hs *HTTPServer) setIndexViewData(c *m.ReqContext) (*dtos.IndexViewData, error) {
	if !isUiModulesExported {
		hs.pushModulesToCms(c)
		isUiModulesExported = true
	}

	log.Info("Signed in user name : " + c.SignedInUser.Login)
	if c.SignedInUser.Login != "admin" {
		externalUserId, ok := c.Session.Get("myuserid").(string)
		if ok && externalUserId != "<nil>" {
			log.Debug("rbac user id " + externalUserId)
			externalUserPw := c.Session.Get("myuserpw").(string)
			fmt.Println("external user id " + externalUserId)
			if c.SignedInUser.Login == setting.ApplicationAdminUser {
				return hs.setIndexViewDataForRbacCmsAdminUser(c)
			} else {
				return hs.setIndexViewDataForRbacUser(externalUserId, externalUserPw, c)
			}

		}
	}
	//externalUserId, ok := c.Session.Get("myuserid").(string)
	//if ok {
	//	log.Debug("rbac user id " + externalUserId)
	//	externalUserPw := c.Session.Get("myuserpw").(string)
	//	fmt.Println("external user id " + externalUserId)
	//	if externalUserId != "admin" {
	//		return hs.setIndexViewDataForRbacUser(externalUserId, externalUserPw, c)
	//	}
	//}
	settings, err := hs.getFrontendSettingsMap(c)
	if err != nil {
		return nil, err
	}

	prefsQuery := m.GetPreferencesWithDefaultsQuery{User: c.SignedInUser}
	if err := bus.Dispatch(&prefsQuery); err != nil {
		return nil, err
	}
	prefs := prefsQuery.Result

	// Read locale from acccept-language
	acceptLang := c.Req.Header.Get("Accept-Language")
	locale := "en-US"

	if len(acceptLang) > 0 {
		parts := strings.Split(acceptLang, ",")
		locale = parts[0]
	}

	appURL := setting.AppUrl
	appSubURL := setting.AppSubUrl

	// special case when doing localhost call from phantomjs
	if c.IsRenderCall {
		appURL = fmt.Sprintf("%s://localhost:%s", setting.Protocol, setting.HttpPort)
		appSubURL = ""
		settings["appSubUrl"] = ""
	}

	hasEditPermissionInFoldersQuery := m.HasEditPermissionInFoldersQuery{SignedInUser: c.SignedInUser}
	if err := bus.Dispatch(&hasEditPermissionInFoldersQuery); err != nil {
		return nil, err
	}

	var data = dtos.IndexViewData{
		User: &dtos.CurrentUser{
			Id:                         c.UserId,
			IsSignedIn:                 c.IsSignedIn,
			Login:                      c.Login,
			Email:                      c.Email,
			Name:                       c.Name,
			OrgCount:                   c.OrgCount,
			OrgId:                      c.OrgId,
			OrgName:                    c.OrgName,
			OrgRole:                    c.OrgRole,
			GravatarUrl:                dtos.GetGravatarUrl(c.Email),
			IsGrafanaAdmin:             c.IsGrafanaAdmin,
			LightTheme:                 prefs.Theme == lightName,
			Timezone:                   prefs.Timezone,
			Locale:                     locale,
			HelpFlags1:                 c.HelpFlags1,
			HasEditPermissionInFolders: hasEditPermissionInFoldersQuery.Result,
		},
		Settings:                settings,
		Theme:                   prefs.Theme,
		AppUrl:                  appURL,
		AppSubUrl:               appSubURL,
		GoogleAnalyticsId:       setting.GoogleAnalyticsId,
		GoogleTagManagerId:      setting.GoogleTagManagerId,
		BuildVersion:            setting.BuildVersion,
		BuildCommit:             setting.BuildCommit,
		NewGrafanaVersion:       plugins.GrafanaLatestVersion,
		NewGrafanaVersionExists: plugins.GrafanaHasUpdate,
		AppName:                 setting.ApplicationName,
		AppNameBodyClass:        getAppNameBodyClass(setting.ApplicationName),
	}

	if setting.DisableGravatar {
		data.User.GravatarUrl = setting.AppSubUrl + "/public/img/user_profile.png"
	}

	if len(data.User.Name) == 0 {
		data.User.Name = data.User.Login
	}

	themeURLParam := c.Query("theme")
	if themeURLParam == lightName {
		data.User.LightTheme = true
		data.Theme = lightName
	} else if themeURLParam == darkName {
		data.User.LightTheme = false
		data.Theme = darkName
	}

	if hasEditPermissionInFoldersQuery.Result {
		children := []*dtos.NavLink{
			{Text: "Dashboard", Icon: "gicon gicon-dashboard-new", Url: setting.AppSubUrl + "/dashboard/new"},
		}

		if c.OrgRole == m.ROLE_ADMIN || c.OrgRole == m.ROLE_EDITOR {
			children = append(children, &dtos.NavLink{Text: "Folder", SubTitle: "Create a new folder to organize your dashboards", Id: "folder", Icon: "gicon gicon-folder-new", Url: setting.AppSubUrl + "/dashboards/folder/new"})
		}

		children = append(children, &dtos.NavLink{Text: "Import", SubTitle: "Import dashboard from file or Grafana.com", Id: "import", Icon: "gicon gicon-dashboard-import", Url: setting.AppSubUrl + "/dashboard/import"})

		data.NavTree = append(data.NavTree, &dtos.NavLink{
			Text:     "Create",
			Id:       "create",
			Icon:     "fa fa-fw fa-plus",
			Url:      setting.AppSubUrl + "/dashboard/new",
			Children: children,
		})
	}

	dashboardChildNavs := []*dtos.NavLink{
		{Text: "Home", Id: "home", Url: setting.AppSubUrl + "/", Icon: "gicon gicon-home", HideFromTabs: true},
		{Text: "Divider", Divider: true, Id: "divider", HideFromTabs: true},
		{Text: "Manage", Id: "manage-dashboards", Url: setting.AppSubUrl + "/dashboards", Icon: "gicon gicon-manage"},
		{Text: "Playlists", Id: "playlists", Url: setting.AppSubUrl + "/playlists", Icon: "gicon gicon-playlists"},
		{Text: "Snapshots", Id: "snapshots", Url: setting.AppSubUrl + "/dashboard/snapshots", Icon: "gicon gicon-snapshots"},
	}

	data.NavTree = append(data.NavTree, &dtos.NavLink{
		Text:     "Dashboards",
		Id:       "dashboards",
		SubTitle: "Manage dashboards & folders",
		Icon:     "gicon gicon-dashboard",
		Url:      setting.AppSubUrl + "/",
		Children: dashboardChildNavs,
	})

	if setting.ExploreEnabled && (c.OrgRole == m.ROLE_ADMIN || c.OrgRole == m.ROLE_EDITOR) {
		data.NavTree = append(data.NavTree, &dtos.NavLink{
			Text:     "Explore",
			Id:       "explore",
			SubTitle: "Explore your data",
			Icon:     "fa fa-rocket",
			Url:      setting.AppSubUrl + "/explore",
		})
	}

	if c.IsSignedIn {
		// Only set login if it's different from the name
		var login string
		if c.SignedInUser.Login != c.SignedInUser.NameOrFallback() {
			login = c.SignedInUser.Login
		}
		profileNode := &dtos.NavLink{
			Text:         c.SignedInUser.NameOrFallback(),
			SubTitle:     login,
			Id:           "profile",
			Img:          data.User.GravatarUrl,
			Url:          setting.AppSubUrl + "/profile",
			HideFromMenu: true,
			Children: []*dtos.NavLink{
				{Text: "Preferences", Id: "profile-settings", Url: setting.AppSubUrl + "/profile", Icon: "gicon gicon-preferences"},
				{Text: "Change Password", Id: "change-password", Url: setting.AppSubUrl + "/profile/password", Icon: "fa fa-fw fa-lock", HideFromMenu: true},
			},
		}

		if !setting.DisableSignoutMenu {
			// add sign out first
			profileNode.Children = append(profileNode.Children, &dtos.NavLink{
				Text: "Sign out", Id: "sign-out", Url: setting.AppSubUrl + "/logout", Icon: "fa fa-fw fa-sign-out", Target: "_self",
			})
		}

		data.NavTree = append(data.NavTree, profileNode)
	}

	enabledPlugins, err := plugins.GetEnabledPlugins(c.OrgId)
	if err != nil {
		return nil, err
	}

	for _, plugin := range enabledPlugins.Apps {
		if plugin.Pinned {
			appLink := &dtos.NavLink{
				Text: plugin.Name,
				Id:   "plugin-page-" + plugin.Id,
				Url:  plugin.DefaultNavUrl,
				Img:  plugin.Info.Logos.Small,
			}

			for _, include := range plugin.Includes {
				if !c.HasUserRole(include.Role) {
					continue
				}

				if include.Type == "page" && include.AddToNav {
					link := &dtos.NavLink{
						Url:  setting.AppSubUrl + "/plugins/" + plugin.Id + "/page/" + include.Slug,
						Text: include.Name,
					}
					appLink.Children = append(appLink.Children, link)
				}

				if include.Type == "dashboard" && include.AddToNav {
					link := &dtos.NavLink{
						Url:  setting.AppSubUrl + "/dashboard/db/" + include.Slug,
						Text: include.Name,
					}
					appLink.Children = append(appLink.Children, link)
				}
			}

			if len(appLink.Children) > 0 && c.OrgRole == m.ROLE_ADMIN {
				appLink.Children = append(appLink.Children, &dtos.NavLink{Divider: true})
				appLink.Children = append(appLink.Children, &dtos.NavLink{Text: "Plugin Config", Icon: "gicon gicon-cog", Url: setting.AppSubUrl + "/plugins/" + plugin.Id + "/edit"})
			}

			if len(appLink.Children) > 0 {
				data.NavTree = append(data.NavTree, appLink)
			}
		}
	}
	if setting.AlertingEnabled && (c.OrgRole == m.ROLE_ADMIN || c.OrgRole == m.ROLE_EDITOR) {
		alertChildNavs := []*dtos.NavLink{
			{Text: "Alert Rules", Id: "alert-list", Url: setting.AppSubUrl + "/alerting/list", Icon: "gicon gicon-alert-rules"},
			{Text: "Notification channels", Id: "channels", Url: setting.AppSubUrl + "/alerting/notifications", Icon: "gicon gicon-alert-notification-channel"},
		}

		data.NavTree = append(data.NavTree, &dtos.NavLink{
			Text:     "Alerting",
			SubTitle: "Alert rules & notifications",
			Id:       "alerting",
			Icon:     "gicon gicon-alert",
			Url:      setting.AppSubUrl + "/alerting/list",
			Children: alertChildNavs,
		})
	}
	if c.IsGrafanaAdmin || c.OrgRole == m.ROLE_ADMIN {
		// cfgNode := &dtos.NavLink{
		// 	Id:       "cfg",
		// 	Text:     "Configuration",
		// 	SubTitle: "Organization: " + c.OrgName,
		// 	Icon:     "gicon gicon-cog",
		// 	Url:      setting.AppSubUrl + "/datasources",
		// 	Children: []*dtos.NavLink{
		// 		{
		// 			Text:        "Data Sources",
		// 			Icon:        "gicon gicon-datasources",
		// 			Description: "Add and configure data sources",
		// 			Id:          "datasources",
		// 			Url:         setting.AppSubUrl + "/datasources",
		// 		},
		// 		{
		// 			Text:        "Users",
		// 			Id:          "users",
		// 			Description: "Manage org members",
		// 			Icon:        "gicon gicon-user",
		// 			Url:         setting.AppSubUrl + "/org/users",
		// 		},
		// 		{
		// 			Text:        "Teams",
		// 			Id:          "teams",
		// 			Description: "Manage org groups",
		// 			Icon:        "gicon gicon-team",
		// 			Url:         setting.AppSubUrl + "/org/teams",
		// 		},
		// 		{
		// 			Text:        "Plugins",
		// 			Id:          "plugins",
		// 			Description: "View and configure plugins",
		// 			Icon:        "gicon gicon-plugins",
		// 			Url:         setting.AppSubUrl + "/plugins",
		// 		},
		// 		{
		// 			Text:        "Preferences",
		// 			Id:          "org-settings",
		// 			Description: "Organization preferences",
		// 			Icon:        "gicon gicon-preferences",
		// 			Url:         setting.AppSubUrl + "/org",
		// 		},

		// 		{
		// 			Text:        "API Keys",
		// 			Id:          "apikeys",
		// 			Description: "Create & manage API keys",
		// 			Icon:        "gicon gicon-apikeys",
		// 			Url:         setting.AppSubUrl + "/org/apikeys",
		// 		},
		// 		{
		// 			Text:        "Legal Entities",
		// 			Id:          "legalentities",
		// 			Description: "Add legal entities",
		// 			Icon:        "fa fa-balance-scale",
		// 			Url:         setting.AppSubUrl + "/legalentities",
		// 		},
		// 	},
		// }

		cfgNode := &dtos.NavLink{
			Id:       "cfg",
			Text:     "Configuration",
			SubTitle: "Applications",
			Icon:     "gicon gicon-cog",
			Url:      setting.AppSubUrl + "/localapp",
			Children: []*dtos.NavLink{
				{
					Text:         "Local Config",
					Icon:         "gicon gicon-datasources",
					Description:  "Local Application",
					Id:           "localapp",
					Url:          setting.AppSubUrl + "/localapp",
					HideFromTabs: true,
				},
				{
					Text:         "Global Config",
					Icon:         "gicon gicon-datasources",
					Description:  "Add and configure data sources",
					Id:           "datasources",
					Url:          setting.AppSubUrl + "/datasources",
					HideFromTabs: true,
				},
				{
					Text:         "Data Sources",
					Icon:         "gicon gicon-datasources",
					Description:  "Add and configure data sources",
					Id:           "datasources",
					Url:          setting.AppSubUrl + "/datasources",
					HideFromMenu: true,
				},
				{
					Text:         "Users",
					Id:           "users",
					Description:  "Manage org members",
					Icon:         "gicon gicon-user",
					Url:          setting.AppSubUrl + "/org/users",
					HideFromMenu: true,
					HideFromTabs: true,
				},
				{
					Text:         "Teams",
					Id:           "teams",
					Description:  "Manage org groups",
					Icon:         "gicon gicon-team",
					Url:          setting.AppSubUrl + "/org/teams",
					HideFromMenu: true,
					HideFromTabs: true,
				},
				{
					Text:         "Plugins",
					Id:           "plugins",
					Description:  "View and configure plugins",
					Icon:         "gicon gicon-plugins",
					Url:          setting.AppSubUrl + "/plugins",
					HideFromMenu: true,
				},
				{
					Text:         "Preferences",
					Id:           "org-settings",
					Description:  "Organization preferences",
					Icon:         "gicon gicon-preferences",
					Url:          setting.AppSubUrl + "/org",
					HideFromMenu: true,
				},

				{
					Text:         "API Keys",
					Id:           "apikeys",
					Description:  "Create & manage API keys",
					Icon:         "gicon gicon-apikeys",
					Url:          setting.AppSubUrl + "/org/apikeys",
					HideFromMenu: true,
				},
			},
		}

		if c.OrgRole != m.ROLE_ADMIN {
			cfgNode = &dtos.NavLink{
				Id:       "cfg",
				Text:     "Configuration",
				SubTitle: "Organization: " + c.OrgName,
				Icon:     "gicon gicon-cog",
				Url:      setting.AppSubUrl + "/admin/users",
				Children: make([]*dtos.NavLink, 0),
			}
		}

		if c.OrgRole == m.ROLE_ADMIN && c.IsGrafanaAdmin {
			cfgNode.Children = append(cfgNode.Children, &dtos.NavLink{
				Divider: true, HideFromTabs: true, Id: "admin-divider", Text: "Text",
			})
		}

		if c.IsGrafanaAdmin {
			cfgNode.Children = append(cfgNode.Children, &dtos.NavLink{
				Text:         "Server Admin",
				HideFromTabs: true,
				SubTitle:     "Manage all users & orgs",
				Id:           "admin",
				Icon:         "gicon gicon-shield",
				Url:          setting.AppSubUrl + "/admin/users",
				Children: []*dtos.NavLink{
					{Text: "Users", Id: "global-users", Url: setting.AppSubUrl + "/admin/users", Icon: "gicon gicon-user"},
					{Text: "Orgs", Id: "global-orgs", Url: setting.AppSubUrl + "/admin/orgs", Icon: "gicon gicon-org"},
					{Text: "Settings", Id: "server-settings", Url: setting.AppSubUrl + "/admin/settings", Icon: "gicon gicon-preferences"},
					{Text: "Stats", Id: "server-stats", Url: setting.AppSubUrl + "/admin/stats", Icon: "fa fa-fw fa-bar-chart"},
					{Text: "Style Guide", Id: "styleguide", Url: setting.AppSubUrl + "/styleguide", Icon: "fa fa-fw fa-eyedropper"},
				},
			})

			cfgNode.Children = append(cfgNode.Children, &dtos.NavLink{
				Text:         "Dashboard List",
				HideFromTabs: true,
				SubTitle:     "Manage all dashboards",
				Id:           "custom-dashboards",
				Icon:         "gicon gicon-dashboard",
				Url:          setting.AppSubUrl + "/customdashboards",
			})
		}

		data.NavTree = append(data.NavTree, cfgNode)
	}

	data.NavTree = append(data.NavTree, &dtos.NavLink{
		Text:         "Help",
		SubTitle:     fmt.Sprintf(`%s v%s (%s)`, setting.ApplicationName, setting.BuildVersion, setting.BuildCommit),
		Id:           "help",
		Url:          "#",
		Icon:         "gicon gicon-question",
		HideFromMenu: true,
		Children: []*dtos.NavLink{
			{Text: "Keyboard shortcuts", Url: "/shortcuts", Icon: "fa fa-fw fa-keyboard-o", Target: "_self"},
			{Text: "Community site", Url: "http://community.grafana.com", Icon: "fa fa-fw fa-comment", Target: "_blank"},
			{Text: "Documentation", Url: "http://docs.grafana.org", Icon: "fa fa-fw fa-file", Target: "_blank"},
		},
	})
	hs.HooksService.RunIndexDataHooks(&data)
	return &data, nil
}

func (hs *HTTPServer) Index(c *m.ReqContext) {
	data, err := hs.setIndexViewData(c)
	if err != nil {
		c.Handle(500, "Failed to get settings", err)
		return
	}
	c.HTML(200, "index", data)
}

func (hs *HTTPServer) NotFoundHandler(c *m.ReqContext) {
	if c.IsApiRequest() {
		c.JsonApiErr(404, "Not found", nil)
		return
	}

	data, err := hs.setIndexViewData(c)
	if err != nil {
		c.Handle(500, "Failed to get settings", err)
		return
	}

	c.HTML(404, "index", data)
}

func getAppNameBodyClass(name string) string {
	switch name {
	case setting.APP_NAME:
		return "app-grafana"
	case setting.APP_NAME_ENTERPRISE:
		return "app-enterprise"
	default:
		return ""
	}
}
