package biscuit

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUtil(t *testing.T) {
	testMap := make(map[string]float64)

	testMap["foo"] = 1
	testMap["bar"] = 2
	testMap["baz"] = 3
	testMap["lel"] = 4
	testMap["wat"] = 5

	Convey("Subject: Sorting maps", t, func() {
		Convey("Given an instance of a map[string]float64", func() {
			sorted := SortedKeys(testMap)
			Convey("The sort should return a []string", func() {
				So(sorted, ShouldHaveSameTypeAs, make([]string, 5))
			})
			Convey("The sort should return a []string of equal length to the specified map", func() {
				So(len(sorted), ShouldHaveSameTypeAs, len(testMap))
			})
			Convey("The sort should return a []string in descending order of the specified map's values", func() {
				So(sorted[0], ShouldEqual, "wat")
				So(sorted[len(sorted)-1], ShouldEqual, "foo")
			})
		})
	})
}
