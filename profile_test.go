package biscuit

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"math"
	"testing"
)

func Round(x float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)
	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow
}

func TestProfile(t *testing.T) {
	testTable := []*Ngram{
		&Ngram{"  b", 1},
		&Ngram{" bo", 1},
		&Ngram{"boo", 1},
		&Ngram{"ooy", 1},
		&Ngram{"oya", 1},
		&Ngram{"yah", 1},
		&Ngram{"ah ", 1},
	}
	label := "english"
	text := "booyah"
	n := 3

	english := NewProfile(label, text, n)

	Convey("Subject: Creating profiles", t, func() {
		Convey("Given a label, some text and an ngram length", func() {
			Convey("The profile label should equal the specified value", func() {
				So(english.label, ShouldEqual, label)
			})

			Convey("The profile text should be parsed into an ngram table", func() {
				for _, ngram := range testTable {
					if f, ok := english.FindNgram(ngram.sequence); ok {
						So(ngram.sequence, ShouldEqual, f.sequence)
						So(ngram.frequency, ShouldEqual, f.frequency)
					}
				}
			})

			Convey("The profile ngram length should equal the specified value", func() {
				So(english.n, ShouldEqual, n)
				So(len(english.table[0].sequence), ShouldEqual, n)
			})

			Convey("The profile should be able to locate an associated ngram by sequence", func() {
				ngram, ok := english.FindNgram(testTable[4].sequence)

				So(ngram.sequence, ShouldEqual, testTable[4].sequence)
				So(ok, ShouldEqual, true)
			})
		})
	})

	Convey("Subject: Scoring and comparing profiles", t, func() {
		Convey("Given a corpora of sample text in French, English and German", func() {
			text, _ := ioutil.ReadFile("./corpora/en/angel-island.txt")
			en := string(text)

			text, _ = ioutil.ReadFile("./corpora/fr/candide.txt")
			fr := string(text)

			text, _ = ioutil.ReadFile("./corpora/de/coriolanus.txt")
			de := string(text)

			english := NewProfile("en", en, 3)
			german := NewProfile("de", de, 3)
			french := NewProfile("fr", fr, 3)

			Convey("Comparing a corpus against itself should yield an exact match", func() {
				So(Round(german.Subtract(german), 1), ShouldBeGreaterThanOrEqualTo, 1)
				So(Round(english.Subtract(english), 1), ShouldBeGreaterThanOrEqualTo, 1)
			})

			Convey("Comparing a corpus of one language against another should yield partial match", func() {
				difference := english.Subtract(german)
				So(difference, ShouldBeGreaterThanOrEqualTo, 0)
				So(difference, ShouldBeLessThan, 1)
			})

			Convey("Comparing a corpus of one language against another should yield the same score regardless of order", func() {
				So(german.Subtract(english), ShouldEqual, english.Subtract(german))
			})

			Convey("The vectors of each profile should be properly calculated", func() {
				So(french.length, ShouldEqual, french.Length())
				So(english.length, ShouldEqual, english.Length())
				So(german.length, ShouldEqual, german.Length())
			})

			Convey("French sample text should score as FR", func() {
				unknown := NewProfile("unknown", "Voulez-vous coucher avec moi ce soir?", n)
				So(unknown.Match(french, english, german), ShouldEqual, "fr")
			})

			Convey("German sample text should score as DE", func() {
				unknown := NewProfile("unknown", "Iche bin ein Berliner.", n)
				So(unknown.Match(french, english, german), ShouldEqual, "de")
			})

			Convey("English sample text should score as EN", func() {
				unknown := NewProfile("unknown", "the rain in spain falls mainly on the plane.", n)
				So(unknown.Match(french, english, german), ShouldEqual, "en")
			})
		})
	})
}
