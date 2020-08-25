package api

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/synectiks-ems/ems-ui/pkg/setting"
)

func TestHTTPServer(t *testing.T) {
	Convey("Given a HTTPServer", t, func() {
		ts := &HTTPServer{
			Cfg: setting.NewCfg(),
		}

		Convey("Given that basic auth on the metrics endpoint is enabled", func() {
			ts.Cfg.MetricsEndpointBasicAuthUsername = "foo"
			ts.Cfg.MetricsEndpointBasicAuthPassword = "bar"

			So(ts.metricsEndpointBasicAuthEnabled(), ShouldBeTrue)
		})

		Convey("Given that basic auth on the metrics endpoint is disabled", func() {
			ts.Cfg.MetricsEndpointBasicAuthUsername = ""
			ts.Cfg.MetricsEndpointBasicAuthPassword = ""

			So(ts.metricsEndpointBasicAuthEnabled(), ShouldBeFalse)
		})
	})
}
