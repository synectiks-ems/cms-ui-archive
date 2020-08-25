package middleware

import (
	"fmt"

	"gopkg.in/macaron.v1"

	m "github.com/synectiks-ems/ems-ui/pkg/models"
	"github.com/synectiks-ems/ems-ui/pkg/services/quota"
)

func Quota(target string) macaron.Handler {
	return func(c *m.ReqContext) {
		limitReached, err := quota.QuotaReached(c, target)
		if err != nil {
			c.JsonApiErr(500, "failed to get quota", err)
			return
		}
		if limitReached {
			c.JsonApiErr(403, fmt.Sprintf("%s Quota reached", target), nil)
			return
		}
	}
}
