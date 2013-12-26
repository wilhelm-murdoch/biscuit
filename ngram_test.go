package biscuit

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNgram(t *testing.T) {
	ngram := NewNgram("foo")
	Convey("Subject: Creating ngrams", t, func() {
		Convey("Given a sequence of characters", func() {
			Convey("The sequence should equal the specified value", func() {
				So(ngram.sequence, ShouldEqual, "foo")
			})
			Convey("The default frequency should equal 1", func() {
				So(ngram.frequency, ShouldEqual, 1)
			})
			Convey("The frequency should increase by one when incremented", func() {
				So(ngram.frequency, ShouldEqual, 1)
				ngram.Increment()
				So(ngram.frequency, ShouldEqual, 2)
			})
		})
	})
}
