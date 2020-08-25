package sqlstore

import (
	"github.com/synectiks-ems/ems-ui/pkg/bus"
	m "github.com/synectiks-ems/ems-ui/pkg/models"
)

func init() {
	bus.AddHandler("sql", GetDBHealthQuery)
}

func GetDBHealthQuery(query *m.GetDBHealthQuery) error {
	return x.Ping()
}
