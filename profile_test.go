package biscuit

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"log"
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

			unknowns := make(map[string]string)

			unknowns["de"] = "Der Kanalinspektor natrlich!"
			unknowns["en"] = "His listeners knew that he was going to say this."
			// unknowns["es"] = "con sus aperos formados con prendas de procedencia."
			unknowns["es"] = "why, hello there, good sir. Fine day we're having!"
			unknowns["fr"] = "Monsieur le baron était un des plus puissants."
			unknowns["jp"] = "っともすずしく、さらになんの奇跡か、季節はずれというのにまだイチゴが食"
			unknowns["th"] = "ผู้ฟัง ของเขา รู้ว่า เขากำลังจะ พูดแบบนี้ พวกเขารู้มาก"

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

			profilesArray := make([]*Profile, 0, len(profiles))

			for _, p := range profiles {
				profilesArray = append(profilesArray, p)
			}

			for _, profile := range profiles {
				Convey(profile.Label+" sample text should score as "+profile.Label, func() {
					log.Print(profile.Label)
					unknown := NewProfileFromText("unknown", unknowns[profile.Label], n)
					So(unknown.Match(profilesArray), ShouldEqual, profile.Label)
				})
			}
		})
	})
}
