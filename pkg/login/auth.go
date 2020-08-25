package login

import (
	"errors"
	"net/http"

	"github.com/synectiks-ems/ems-ui/pkg/bus"
	m "github.com/synectiks-ems/ems-ui/pkg/models"
)

var (
	ErrEmailNotAllowed       = errors.New("Required email domain not fulfilled")
	ErrInvalidCredentials    = errors.New("Invalid Username or Password")
	ErrNoEmail               = errors.New("Login provider didn't return an email address")
	ErrProviderDeniedRequest = errors.New("Login provider denied login request")
	ErrSignUpNotAllowed      = errors.New("Signup is not allowed for this adapter")
	ErrTooManyLoginAttempts  = errors.New("Too many consecutive incorrect login attempts for user. Login for user temporarily blocked")
	ErrPasswordEmpty         = errors.New("No password provided.")
	ErrUsersQuotaReached     = errors.New("Users quota reached")
	ErrGettingUserQuota      = errors.New("Error getting user quota")
)

func Init() {
	bus.AddHandler("auth", AuthenticateUser)
	loadLdapConfig()
}

var externalSecurityServiceClient = &http.Client{
	//Timeout:   time.Second * 5,
	Transport: &http.Transport{Proxy: http.ProxyFromEnvironment},
}

func AuthenticateUser(query *m.LoginUserQuery) error {

	if query.Username != "admin" {
		//log.Info("Login authentication endpoint hit")
		//response, err := externalSecurityServiceClient.Get(setting.ExternalSecurityUrl + "/security/public/login?username=" + query.Username + "&password=" + query.Password)
		//response, err := externalSecurityServiceClient.Get(setting.CmsUrl + "/api/cmslogin?username=" + query.Username + "&password=" + query.Password)
		//defer response.Body.Close()
		//bodyBytes, err := ioutil.ReadAll(response.Body)
		//bodyString := string(bodyBytes)
		//if err != nil || response.StatusCode == 417 {
		//	log.Error(1, "User authentication failed. Please check the login credentials", bodyString)
		//	return ErrInvalidCredentials
		//}

		//if response.StatusCode == 200 ||  response.StatusCode == 201 {
		//	fmt.Println(bodyString)
		//}
		var usr = m.User{}
		usr.Name = query.Username
		usr.Password = query.Password
		usr.Login = query.Username
		usr.OrgId = 100
		usr.IsAdmin = false
		usr.Id = 100
		query.User = &usr

		return nil

	} else {
		if err := validateLoginAttempts(query.Username); err != nil {
			return err
		}

		if err := validatePasswordSet(query.Password); err != nil {
			return err
		}

		err := loginUsingGrafanaDB(query)
		if err == nil || (err != m.ErrUserNotFound && err != ErrInvalidCredentials) {
			return err
		}

		ldapEnabled, ldapErr := loginUsingLdap(query)
		if ldapEnabled {
			if ldapErr == nil || ldapErr != ErrInvalidCredentials {
				return ldapErr
			}

			err = ldapErr
		}

		if err == ErrInvalidCredentials {
			saveInvalidLoginAttempt(query)
		}

		if err == m.ErrUserNotFound {
			return ErrInvalidCredentials
		}
		return err
	}

}
func validatePasswordSet(password string) error {
	if len(password) == 0 {
		return ErrPasswordEmpty
	}

	return nil
}
