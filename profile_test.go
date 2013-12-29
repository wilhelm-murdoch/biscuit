package biscuit

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestProfile(t *testing.T) {
	testTable := make(map[string]int)
	testTable["  b"] = 1
	testTable[" bo"] = 1
	testTable["boo"] = 1
	testTable["ooy"] = 1
	testTable["oya"] = 1
	testTable["yah"] = 1
	testTable["ah "] = 1

	label := "english"
	text := "booyah"
	n := 3

	english := NewProfileFromText(label, text, n)

	Convey("Subject: Creating profiles", t, func() {
		Convey("Given a label, some text and an ngram length", func() {
			Convey("The profile label should equal the specified value", func() {
				So(english.Label, ShouldEqual, label)
			})

			Convey("The profile text should be parsed into an ngram table", func() {
				for sequence, frequency := range testTable {
					if f, ok := english.Ngrams[sequence]; ok {
						So(frequency, ShouldEqual, f)
					}
				}
			})

			Convey("The profile ngram length should equal the specified value", func() {
				So(english.N, ShouldEqual, n)
				for sequence := range english.Ngrams {
					So(len(sequence), ShouldEqual, n)
				}
			})
		})
	})

	Convey("Subject: Scoring and comparing profiles", t, func() {
		Convey("Given a corpora of sample text in various languages", func() {
			corpora := make(map[string]string)
			profiles := make(map[string]*Profile)
			samples, _ := filepath.Glob("./corpora/*.txt")

			for _, file := range samples {
				text, _ := ioutil.ReadFile(file)
				label = filepath.Base(file)[:2]
				corpora[label] = string(text)
				profiles[label] = NewProfileFromText(label, corpora[label], n)
			}

			profileInstances := make([]*Profile, 0, len(profiles))

			for _, profile := range profiles {
				profileInstances = append(profileInstances, profile)
			}

			Convey("Comparing a corpus against itself should yield an exact match", func() {
				for _, profile := range profiles {
					So(Round(profile.Subtract(profile), 1), ShouldBeGreaterThanOrEqualTo, 1)
				}
			})

			Convey("Comparing a corpus of one language against another should yield partial match", func() {
				difference := profiles["en"].Subtract(profiles["de"])
				So(difference, ShouldBeGreaterThanOrEqualTo, 0)
				So(difference, ShouldBeLessThan, 1)
			})

			Convey("Comparing a corpus of one language against another should yield the same score regardless of order", func() {
				So(profiles["de"].Subtract(profiles["en"]), ShouldEqual, profiles["en"].Subtract(profiles["de"]))
			})

			Convey("The vectors of each profile should be properly calculated", func() {
				for _, profile := range profiles {
					So(profile.length, ShouldEqual, profile.Length())
				}
			})

			Convey("DE sample text should score as DE", func() {
				unknown := NewProfileFromText("unknown", "Der Kanalinspektor natrlich!", n)
				So(unknown.Match(profileInstances), ShouldEqual, "de")
			})

			Convey("ES sample text should score as ES", func() {
				unknown := NewProfileFromText("unknown", "con sus aperos formados con prendas de procedencia.", n)
				So(unknown.Match(profileInstances), ShouldEqual, "es")
			})

			Convey("FR sample text should score as FR", func() {
				unknown := NewProfileFromText("unknown", "Monsieur le baron était un des plus puissants.", n)
				So(unknown.Match(profileInstances), ShouldEqual, "fr")
			})

			Convey("JP sample text should score as JP", func() {
				unknown := NewProfileFromText("unknown", "っともすずしく、さらになんの奇跡か、季節はずれというのにまだイチゴが食", n)
				So(unknown.Match(profileInstances), ShouldEqual, "jp")
			})

			Convey("TH sample text should score as TH", func() {
				unknown := NewProfileFromText("unknown", "ผู้ฟัง ของเขา รู้ว่า เขากำลังจะ พูดแบบนี้ พวกเขารู้มาก", n)
				So(unknown.Match(profileInstances), ShouldEqual, "th")
			})
		})
	})
}
