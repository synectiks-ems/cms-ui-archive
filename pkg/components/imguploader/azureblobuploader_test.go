package imguploader

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/xformation/cms-ui/pkg/setting"
)

func TestUploadToAzureBlob(t *testing.T) {
	SkipConvey("[Integration test] for external_image_store.azure_blob", t, func() {
		cfg := setting.NewCfg()
		err := cfg.Load(&setting.CommandLineArgs{
			HomePath: "../../../",
		})
		So(err, ShouldBeNil)

		uploader, _ := NewImageUploader()

		path, err := uploader.Upload(context.Background(), "../../../public/img/logo_transparent_400x.png")

		So(err, ShouldBeNil)
		So(path, ShouldNotEqual, "")
	})
}
