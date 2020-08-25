package login

import (
	"time"

	"github.com/synectiks-ems/ems-ui/pkg/bus"
	m "github.com/synectiks-ems/ems-ui/pkg/models"
	"github.com/synectiks-ems/ems-ui/pkg/setting"
)

var (
	maxInvalidLoginAttempts int64 = 5
	loginAttemptsWindow           = time.Minute * 5
)

var validateLoginAttempts = func(username string) error {
	if setting.DisableBruteForceLoginProtection {
		return nil
	}

	loginAttemptCountQuery := m.GetUserLoginAttemptCountQuery{
		Username: username,
		Since:    time.Now().Add(-loginAttemptsWindow),
	}

	if err := bus.Dispatch(&loginAttemptCountQuery); err != nil {
		return err
	}

	if loginAttemptCountQuery.Result >= maxInvalidLoginAttempts {
		return ErrTooManyLoginAttempts
	}

	return nil
}

var saveInvalidLoginAttempt = func(query *m.LoginUserQuery) {
	if setting.DisableBruteForceLoginProtection {
		return
	}

	loginAttemptCommand := m.CreateLoginAttemptCommand{
		Username:  query.Username,
		IpAddress: query.IpAddress,
	}

	bus.Dispatch(&loginAttemptCommand)
}
