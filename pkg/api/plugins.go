package api

import (
	"bytes"
	"github.com/synectiks-ems/ems-ui/pkg/components/simplejson"
	"github.com/synectiks-ems/ems-ui/pkg/log"
	"net/http"
	"sort"

	"github.com/synectiks-ems/ems-ui/pkg/api/dtos"
	"github.com/synectiks-ems/ems-ui/pkg/bus"
	m "github.com/synectiks-ems/ems-ui/pkg/models"
	"github.com/synectiks-ems/ems-ui/pkg/plugins"
	"github.com/synectiks-ems/ems-ui/pkg/setting"
)

var moduleRestClient = &http.Client{
	Transport: &http.Transport{Proxy: http.ProxyFromEnvironment},
}

func (hs *HTTPServer) GetPluginList(c *m.ReqContext) Response {
	typeFilter := c.Query("type")
	enabledFilter := c.Query("enabled")
	embeddedFilter := c.Query("embedded")
	coreFilter := c.Query("core")

	pluginSettingsMap, err := plugins.GetPluginSettings(c.OrgId)

	if err != nil {
		return Error(500, "Failed to get list of plugins", err)
	}

	result := make(dtos.PluginList, 0)
	for _, pluginDef := range plugins.Plugins {
		// filter out app sub plugins
		if embeddedFilter == "0" && pluginDef.IncludedInAppId != "" {
			continue
		}

		// filter out core plugins
		if coreFilter == "0" && pluginDef.IsCorePlugin {
			continue
		}

		// filter on type
		if typeFilter != "" && typeFilter != pluginDef.Type {
			continue
		}

		if pluginDef.State == plugins.PluginStateAlpha && !hs.Cfg.EnableAlphaPanels {
			continue
		}

		listItem := dtos.PluginListItem{
			Id:            pluginDef.Id,
			Name:          pluginDef.Name,
			Type:          pluginDef.Type,
			Info:          &pluginDef.Info,
			LatestVersion: pluginDef.GrafanaNetVersion,
			HasUpdate:     pluginDef.GrafanaNetHasUpdate,
			DefaultNavUrl: pluginDef.DefaultNavUrl,
			State:         pluginDef.State,
		}

		if pluginSetting, exists := pluginSettingsMap[pluginDef.Id]; exists {
			listItem.Enabled = pluginSetting.Enabled
			listItem.Pinned = pluginSetting.Pinned
		}

		if listItem.DefaultNavUrl == "" || !listItem.Enabled {
			listItem.DefaultNavUrl = setting.AppSubUrl + "/plugins/" + listItem.Id + "/edit"
		}

		// filter out disabled
		if enabledFilter == "1" && !listItem.Enabled {
			continue
		}

		// filter out built in data sources
		if ds, exists := plugins.DataSources[pluginDef.Id]; exists {
			if ds.BuiltIn {
				continue
			}
		}

		result = append(result, listItem)
	}

	sort.Sort(result)
	return JSON(200, result)
}

func GetPluginSettingByID(c *m.ReqContext) Response {
	pluginID := c.Params(":pluginId")

	def, exists := plugins.Plugins[pluginID]
	if !exists {
		return Error(404, "Plugin not found, no installed plugin with that id", nil)
	}

	dto := &dtos.PluginSetting{
		Type:          def.Type,
		Id:            def.Id,
		Name:          def.Name,
		Info:          &def.Info,
		Dependencies:  &def.Dependencies,
		Includes:      def.Includes,
		BaseUrl:       def.BaseUrl,
		Module:        def.Module,
		DefaultNavUrl: def.DefaultNavUrl,
		LatestVersion: def.GrafanaNetVersion,
		HasUpdate:     def.GrafanaNetHasUpdate,
		State:         def.State,
	}

	query := m.GetPluginSettingByIdQuery{PluginId: pluginID, OrgId: c.OrgId}
	if err := bus.Dispatch(&query); err != nil {
		if err != m.ErrPluginSettingNotFound {
			return Error(500, "Failed to get login settings", nil)
		}
	} else {
		dto.Enabled = query.Result.Enabled
		dto.Pinned = query.Result.Pinned
		dto.JsonData = query.Result.JsonData
	}

	return JSON(200, dto)
}

func UpdatePluginSetting(c *m.ReqContext, cmd m.UpdatePluginSettingCmd) Response {
	pluginID := c.Params(":pluginId")

	cmd.OrgId = c.OrgId
	cmd.PluginId = pluginID

	if _, ok := plugins.Apps[cmd.PluginId]; !ok {
		return Error(404, "Plugin not installed.", nil)
	}
	var plg = plugins.Apps[cmd.PluginId]
	var status = "DEACTIVE"
	if cmd.Enabled {
		status = "ACTIVE"
	}
	// call rhe rest service to changes the plugin status in cms
	message := map[string]interface{}{
		"moduleName": plg.Name,
		"status":     status,
	}
	messageBytes, err := simplejson.NewFromAny(message).Encode()
	if err != nil {
		log.Error(500, "Failed to encode message in bytes", err)
		return Error(500, "Failed to update plugin setting", err)
	}
	httpClient := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, setting.CmsUrl+"/api/cmsmodules", bytes.NewBuffer(messageBytes))
	//req.Header.Set("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	response, err := httpClient.Do(req)
	if err != nil {
		log.Error(500, "Failed to return rest response", err)
		return Error(500, "Failed to update plugin setting", err)
	}
	if response.StatusCode == 200 || response.StatusCode == 201 {
		log.Info("Plugin settings updated in cms")
		if err := bus.Dispatch(&cmd); err != nil {
			log.Error(500, "Failed to update plugin setting", err)
			return Error(500, "Failed to update plugin setting", err)
		}
		return Success("Plugin settings updated")
	}
	return Error(500, "Failed to update plugin setting", err)

}

func GetPluginDashboards(c *m.ReqContext) Response {
	pluginID := c.Params(":pluginId")

	list, err := plugins.GetPluginDashboards(c.OrgId, pluginID)
	if err != nil {
		if notfound, ok := err.(plugins.PluginNotFoundError); ok {
			return Error(404, notfound.Error(), nil)
		}

		return Error(500, "Failed to get plugin dashboards", err)
	}

	return JSON(200, list)
}

func GetPluginMarkdown(c *m.ReqContext) Response {
	pluginID := c.Params(":pluginId")
	name := c.Params(":name")

	content, err := plugins.GetPluginMarkdown(pluginID, name)
	if err != nil {
		if notfound, ok := err.(plugins.PluginNotFoundError); ok {
			return Error(404, notfound.Error(), nil)
		}

		return Error(500, "Could not get markdown file", err)
	}

	resp := Respond(200, content)
	resp.Header("Content-Type", "text/plain; charset=utf-8")
	return resp
}

func ImportDashboard(c *m.ReqContext, apiCmd dtos.ImportDashboardCommand) Response {

	cmd := plugins.ImportDashboardCommand{
		OrgId:     c.OrgId,
		User:      c.SignedInUser,
		PluginId:  apiCmd.PluginId,
		Path:      apiCmd.Path,
		Inputs:    apiCmd.Inputs,
		Overwrite: apiCmd.Overwrite,
		FolderId:  apiCmd.FolderId,
		Dashboard: apiCmd.Dashboard,
	}

	if err := bus.Dispatch(&cmd); err != nil {
		return Error(500, "Failed to import dashboard", err)
	}

	return JSON(200, cmd.Result)
}
