package nginx

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGraphite(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Load the sample config file", t, func() {

		config := ReadConfig("../config.toml")
		So(config, ShouldNotBeNil)

		Convey("Connect to sample graphite server", func() {

		})
	})
}
