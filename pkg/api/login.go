package api

import (
	"encoding/json"
	"github.com/synectiks-ems/ems-ui/pkg/login"
	"io/ioutil"
	"net/url"

	"github.com/synectiks-ems/ems-ui/pkg/api/dtos"
	"github.com/synectiks-ems/ems-ui/pkg/bus"
	"github.com/synectiks-ems/ems-ui/pkg/log"
	"github.com/synectiks-ems/ems-ui/pkg/metrics"
	m "github.com/synectiks-ems/ems-ui/pkg/models"
	"github.com/synectiks-ems/ems-ui/pkg/services/session"
	"github.com/synectiks-ems/ems-ui/pkg/setting"
)

const (
	ViewIndex = "index"
)

func (hs *HTTPServer) LoginView(c *m.ReqContext) {
	viewData, err := hs.setIndexViewData(c)
	if err != nil {
		c.Handle(500, "Failed to get settings", err)
		return
	}

	enabledOAuths := make(map[string]interface{})
	for key, oauth := range setting.OAuthService.OAuthInfos {
		enabledOAuths[key] = map[string]string{"name": oauth.Name}
	}

	viewData.Settings["oauth"] = enabledOAuths
	viewData.Settings["disableUserSignUp"] = !setting.AllowUserSignUp
	viewData.Settings["loginHint"] = setting.LoginHint
	viewData.Settings["disableLoginForm"] = setting.DisableLoginForm

	if loginError, ok := c.Session.Get("loginError").(string); ok {
		c.Session.Delete("loginError")
		viewData.Settings["loginError"] = loginError
	}

	if tryOAuthAutoLogin(c) {
		return
	}

	if !tryLoginUsingRememberCookie(c) {
		c.HTML(200, ViewIndex, viewData)
		return
	}

	if redirectTo, _ := url.QueryUnescape(c.GetCookie("redirect_to")); len(redirectTo) > 0 {
		c.SetCookie("redirect_to", "", -1, setting.AppSubUrl+"/")
		c.Redirect(redirectTo)
		return
	}

	c.Redirect(setting.AppSubUrl + "/")
}

func tryOAuthAutoLogin(c *m.ReqContext) bool {
	if !setting.OAuthAutoLogin {
		return false
	}
	oauthInfos := setting.OAuthService.OAuthInfos
	if len(oauthInfos) != 1 {
		log.Warn("Skipping OAuth auto login because multiple OAuth providers are configured.")
		return false
	}
	for key := range setting.OAuthService.OAuthInfos {
		redirectUrl := setting.AppSubUrl + "/login/" + key
		log.Info("OAuth auto login enabled. Redirecting to " + redirectUrl)
		c.Redirect(redirectUrl, 307)
		return true
	}
	return false
}

func tryLoginUsingRememberCookie(c *m.ReqContext) bool {
	// Check auto-login.
	uname := c.GetCookie(setting.CookieUserName)
	if len(uname) == 0 {
		return false
	}

	isSucceed := false
	defer func() {
		if !isSucceed {
			log.Trace("auto-login cookie cleared: %s", uname)
			c.SetCookie(setting.CookieUserName, "", -1, setting.AppSubUrl+"/")
			c.SetCookie(setting.CookieRememberName, "", -1, setting.AppSubUrl+"/")
			return
		}
	}()

	userQuery := m.GetUserByLoginQuery{LoginOrEmail: uname}
	if err := bus.Dispatch(&userQuery); err != nil {
		return false
	}

	user := userQuery.Result

	// validate remember me cookie
	signingKey := user.Rands + user.Password
	if len(signingKey) < 10 {
		c.Logger.Error("Invalid user signingKey")
		return false
	}

	if val, _ := c.GetSuperSecureCookie(signingKey, setting.CookieRememberName); val != user.Login {
		return false
	}

	isSucceed = true
	loginUserWithUser(user, c)
	return true
}

func LoginAPIPing(c *m.ReqContext) {
	if !tryLoginUsingRememberCookie(c) {
		c.JsonApiErr(401, "Unauthorized", nil)
		return
	}

	c.JsonOK("Logged in")
}

func LoginPost(c *m.ReqContext, cmd dtos.LoginCommand) Response {

	result := map[string]interface{}{
		"message": "Logged in",
	}
	if cmd.User == setting.ApplicationAdminUser {
		authQuery := &m.LoginUserQuery{
			ReqContext: c,
			Username:   cmd.User,
			Password:   cmd.Password,
			IpAddress:  c.Req.RemoteAddr,
		}
		if cmd.Password != "cmsadmin" {
			return Error(401, "Invalid username or password", nil)
		}

		var usr = &m.User{
			Password: cmd.Password,
			Login:    cmd.User,
			Name:     cmd.User,
			OrgId:    1,
			IsAdmin:  true,
			Id:       100,
		}
		authQuery.User = usr
		c.SignedInUser = &m.SignedInUser{
			Name:   cmd.User,
			Login:  cmd.User,
			OrgId:  usr.OrgId,
			UserId: usr.Id,
			Email:  cmd.User,
		}
		var userInfo = make(map[string]interface{})
		loginUserWithRbacUser(authQuery.User, c, userInfo)

		if redirectTo, _ := url.QueryUnescape(c.GetCookie("redirect_to")); len(redirectTo) > 0 {
			result["redirectUrl"] = redirectTo
			c.SetCookie("redirect_to", "", -1, setting.AppSubUrl+"/")
		}

		metrics.M_Api_Login_Post.Inc()
	} else {
		if cmd.User != "admin" {
			log.Info("Login authentication endpoint hit")
			//response, err := externalSecurityServiceClient.Get(setting.ExternalSecurityUrl + "/security/public/login?username=" + query.Username + "&password=" + query.Password)
			response, err := externalSecurityServiceClient.Get(setting.CmsUrl + "/api/cmslogin?username=" + cmd.User + "&password=" + cmd.Password)
			if err != nil || response.StatusCode == 417 {
				log.Error(1, "User authentication failed. Please check the login credentials", err)
				return Error(401, "Invalid username or password", err)
			}
			defer response.Body.Close()
			bodyBytes, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Error(1, "Invalid response", err)
				return Error(401, "Invalid response", err)
			}
			bodyString := string(bodyBytes)
			var userInfo = make(map[string]interface{})
			errw := json.Unmarshal([]byte(bodyString), &userInfo)
			if errw != nil {
				log.Error(1, "Response unmarshalling to JSON failed", errw)
				return Error(401, "Response unmarshalling to JSON failed", errw)
			}
			//c.Session.Set(cmd.User, userInfo)
			//fmt.Println("User Info : ",userInfo)
			authQuery := &m.LoginUserQuery{
				ReqContext: c,
				Username:   cmd.User,
				Password:   cmd.Password,
				IpAddress:  c.Req.RemoteAddr,
			}
			var usr = &m.User{
				Password: cmd.Password,
				Login:    cmd.User,
				Name:     cmd.User,
				OrgId:    1,
				IsAdmin:  false,
				Id:       100,
			}
			//usr.Name = cmd.User
			//usr.Password = cmd.Password
			//usr.Login = cmd.User
			//usr.OrgId = 100
			//usr.IsAdmin = false
			//usr.Id = 100
			authQuery.User = usr

			c.SignedInUser = &m.SignedInUser{
				Name:   cmd.User,
				Login:  cmd.User,
				OrgId:  usr.OrgId,
				UserId: usr.Id,
				Email:  cmd.User,
			}

			loginUserWithRbacUser(authQuery.User, c, userInfo)

			if redirectTo, _ := url.QueryUnescape(c.GetCookie("redirect_to")); len(redirectTo) > 0 {
				result["redirectUrl"] = redirectTo
				c.SetCookie("redirect_to", "", -1, setting.AppSubUrl+"/")
			}

			metrics.M_Api_Login_Post.Inc()

			//authQuery.ReqContext.IsSignedIn = true
			//authQuery.ReqContext.SignedInUser = &m.SignedInUser{
			//	Name: cmd.User,
			//	Login: cmd.User,
			//	OrgId: usr.OrgId,
			//	UserId: usr.Id,
			//	Email: cmd.User,
			//}
		} else {
			if setting.DisableLoginForm {
				return Error(401, "Login is disabled", nil)
			}

			authQuery := &m.LoginUserQuery{
				ReqContext: c,
				Username:   cmd.User,
				Password:   cmd.Password,
				IpAddress:  c.Req.RemoteAddr,
			}

			if err := bus.Dispatch(authQuery); err != nil {
				if err == login.ErrInvalidCredentials || err == login.ErrTooManyLoginAttempts {
					return Error(401, "Invalid username or password", err)
				}

				return Error(500, "Error while trying to authenticate user", err)
			}

			user := authQuery.User

			loginUserWithUser(user, c)

			if redirectTo, _ := url.QueryUnescape(c.GetCookie("redirect_to")); len(redirectTo) > 0 {
				result["redirectUrl"] = redirectTo
				c.SetCookie("redirect_to", "", -1, setting.AppSubUrl+"/")
			}

			metrics.M_Api_Login_Post.Inc()

		}
	}
	return JSON(200, result)
}

func loginUserWithRbacUser(user *m.User, c *m.ReqContext, userInfo map[string]interface{}) {
	if user == nil {
		log.Error(3, "User login with nil user")
	}

	c.Resp.Header().Del("Set-Cookie")

	days := 86400 * setting.LogInRememberDays
	if days > 0 {
		c.SetCookie(setting.CookieUserName, user.Login, days, setting.AppSubUrl+"/")
		c.SetSuperSecureCookie(user.Rands+user.Password, setting.CookieRememberName, user.Login, days, setting.AppSubUrl+"/")
	}

	c.Session.RegenerateId(c.Context)
	c.Session.Set(session.SESS_KEY_USERID, user.Id)
	c.Session.Set("myuserid", user.Name)
	c.Session.Set("myuserpw", user.Password)
	c.Session.Set(user.Name, userInfo)

}

func loginUserWithUser(user *m.User, c *m.ReqContext) {
	if user == nil {
		log.Error(3, "User login with nil user")
	}

	c.Resp.Header().Del("Set-Cookie")

	days := 86400 * setting.LogInRememberDays
	if days > 0 {
		c.SetCookie(setting.CookieUserName, user.Login, days, setting.AppSubUrl+"/")
		c.SetSuperSecureCookie(user.Rands+user.Password, setting.CookieRememberName, user.Login, days, setting.AppSubUrl+"/")
	}

	c.Session.RegenerateId(c.Context)
	c.Session.Set(session.SESS_KEY_USERID, user.Id)
	c.Session.Set("myuserid", user.Name)
	c.Session.Set("myuserpw", user.Password)
}

func Logout(c *m.ReqContext) {
	externalSecurityServiceClient.Get(setting.CmsUrl + "/api/cmslogout?username=" + c.SignedInUser.Login)
	c.SetCookie(setting.CookieUserName, "", -1, setting.AppSubUrl+"/")
	c.SetCookie(setting.CookieRememberName, "", -1, setting.AppSubUrl+"/")
	c.Session.Destory(c.Context)
	if setting.SignoutRedirectUrl != "" {
		c.Redirect(setting.SignoutRedirectUrl)
	} else {
		c.Redirect(setting.AppSubUrl + "/login")
	}
}
