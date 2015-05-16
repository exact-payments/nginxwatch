package nginx

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTail(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Load the sample config file", t, func() {

		// config := ReadConfig("../config.toml")
		// So(config, ShouldNotBeNil)

		// Convey("Connect to sample graphite server", func() {

		// 	conn, addr, err := connectToGraphite(config.Server)
		// 	So(err, ShouldBeNil)

		// 	Convey("Send data", func() {
		// 		err := writeData("Test", 100, conn, addr)
		// 		So(err, ShouldBeNil)
		// 	})
		// })
	})
}
