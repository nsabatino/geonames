package geonames

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAlternateNames(t *testing.T) {
	Convey("Test AlternateNames handler", t, func() {
		result, err := AlternateNames()
		Convey("We should not get an error", func() {
			So(err, ShouldBeNil)
		})
		Convey("Result should be correct", func() {
			So(result, ShouldNotBeNil)
			So(result, ShouldNotBeEmpty)
			So(result[0].Name, ShouldEqual, "Zamīn Sūkhteh")
		})
	})
}
