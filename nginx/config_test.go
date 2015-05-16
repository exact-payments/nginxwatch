package nginx

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConfig(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Load the sample config file", t, func() {

		config := ReadConfig("../config.toml")
		So(config, ShouldNotBeNil)

		Convey("Verify Config Contents", func() {

		})
	})
}
