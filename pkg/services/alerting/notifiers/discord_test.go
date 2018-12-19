package notifiers

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/xformation/cms-ui/pkg/components/simplejson"
	m "github.com/xformation/cms-ui/pkg/models"
)

func TestDiscordNotifier(t *testing.T) {
	Convey("Telegram notifier tests", t, func() {

		Convey("Parsing alert notification from settings", func() {
			Convey("empty settings should return error", func() {
				json := `{ }`

				settingsJSON, _ := simplejson.NewJson([]byte(json))
				model := &m.AlertNotification{
					Name:     "discord_testing",
					Type:     "discord",
					Settings: settingsJSON,
				}

				_, err := NewDiscordNotifier(model)
				So(err, ShouldNotBeNil)
			})

			Convey("settings should trigger incident", func() {
				json := `
				{
          "url": "https://web.hook/"
				}`

				settingsJSON, _ := simplejson.NewJson([]byte(json))
				model := &m.AlertNotification{
					Name:     "discord_testing",
					Type:     "discord",
					Settings: settingsJSON,
				}

				not, err := NewDiscordNotifier(model)
				discordNotifier := not.(*DiscordNotifier)

				So(err, ShouldBeNil)
				So(discordNotifier.Name, ShouldEqual, "discord_testing")
				So(discordNotifier.Type, ShouldEqual, "discord")
				So(discordNotifier.WebhookURL, ShouldEqual, "https://web.hook/")
			})
		})
	})
}
