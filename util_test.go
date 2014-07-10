package biscuit

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func BenchmarkRound(b *testing.B) {
	r := rand.New(rand.NewSource(99))

	for i := 0; i < b.N; i++ {
		Round(r.Float64(), 3)
	}
}

func ExampleRound() {
	fmt.Println(Round(math.Pi, 3))
	fmt.Println(Round(math.Pi, 4))

	// Output:
	// 3.14
	// 3.142
}

func BenchmarkSortedKeys(b *testing.B) {
	r := rand.New(rand.NewSource(99))

	things := make(map[string]float64)

	things["foo"] = r.Float64()
	things["bar"] = r.Float64()
	things["baz"] = r.Float64()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		SortedKeys(things)
	}
}

func ExampleSortedKeys() {
	things := make(map[string]float64)

	things["foo"] = 1
	things["bar"] = 8
	things["lel"] = 3
	things["baz"] = 10
	things["wut"] = 0.5

	sorted := SortedKeys(things)

	for _, thing := range sorted {
		fmt.Println(thing)
	}

	// Output:
	// baz
	// bar
	// lel
	// foo
	// wut
}

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
