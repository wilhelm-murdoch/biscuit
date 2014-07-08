package biscuit

import (
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func BenchmarkProfileFromFile(b *testing.B) {
	files := make(map[string]string)
	samples, _ := filepath.Glob("./corpora/*.txt")

	var labels = make([]string, len(samples))

	for i, file := range samples {
		label := filepath.Base(file)[:2]
		labels[i] = label
		files[label] = file
	}

	for i := 0; i < b.N; i++ {
		label := labels[rand.Intn(len(labels))]
		NewProfileFromFile(label, files[label], 3)
	}
}

func BenchmarkMatch(b *testing.B) {
	profiles := make(map[string]*Profile)
	files, _ := filepath.Glob("./corpora/*.txt")

	var labels = make([]string, len(files))
	var profileInstances = make([]*Profile, 0, len(profiles))

	for i, file := range files {
		label := filepath.Base(file)[:2]
		labels[i] = label
		profile, _ := NewProfileFromFile(label, file, 3)
		profiles[label] = profile
		profileInstances = append(profileInstances, profile)
	}

	for i := 0; i < b.N; i++ {
		label := labels[rand.Intn(len(labels))]
		profiles[label].Match(profileInstances)
	}
}

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

	Convey("Subject: Create ngram tables from files...", t, func() {
		Convey("Given a corpora of precalculated ngram tables", func() {
			Convey("Opening a text file should yield a new biscuit.Profile instance", func() {
				samples, _ := filepath.Glob("./corpora/*.txt")

				for _, file := range samples {
					p, err := NewProfileFromFile(filepath.Base(file)[:2], file, 3)
					So(err, ShouldEqual, nil)
					So(filepath.Base(file)[:2], ShouldEqual, p.Label)
					So(p.length, ShouldEqual, p.Length())
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

			Convey("Comparing profiles of different ngram lengths should raise an error", func() {
				unknown := NewProfileFromText("unknown", "sup", n+1)
				match, err := unknown.Match(profileInstances)
				So(match, ShouldEqual, "")
				So(err, ShouldNotEqual, nil)
			})

			Convey("Comparing a corpus against itself should yield an exact match", func() {
				for _, profile := range profiles {
					difference := profile.Subtract(profile)
					So(Round(difference, 1), ShouldBeGreaterThanOrEqualTo, 1)
				}
			})

			Convey("Comparing a corpus of one language against another should yield partial match", func() {
				difference := profiles["en"].Subtract(profiles["de"])
				So(difference, ShouldBeGreaterThanOrEqualTo, 0)
				So(difference, ShouldBeLessThan, 1)
			})

			Convey("Comparing a corpus of one language against another should yield the same score regardless of order", func() {
				difference1 := profiles["de"].Subtract(profiles["en"])
				difference2 := profiles["en"].Subtract(profiles["de"])

				So(difference1, ShouldEqual, difference2)
			})

			Convey("The vectors of each profile should be properly calculated", func() {
				for _, profile := range profiles {
					So(profile.length, ShouldEqual, profile.Length())
				}
			})

			Convey("DE sample text should score as DE", func() {
				unknown := NewProfileFromText("unknown", "Der Kanalinspektor natrlich!", n)
				match, err := unknown.Match(profileInstances)
				So(err, ShouldEqual, nil)
				So(match, ShouldEqual, "de")
			})

			Convey("ES sample text should score as ES", func() {
				unknown := NewProfileFromText("unknown", "con sus aperos formados con prendas de procedencia.", n)
				match, err := unknown.Match(profileInstances)
				So(err, ShouldEqual, nil)
				So(match, ShouldEqual, "es")
			})

			Convey("FR sample text should score as FR", func() {
				unknown := NewProfileFromText("unknown", "Monsieur le baron était un des plus puissants.", n)
				match, err := unknown.Match(profileInstances)
				So(err, ShouldEqual, nil)
				So(match, ShouldEqual, "fr")
			})

			Convey("JP sample text should score as JP", func() {
				unknown := NewProfileFromText("unknown", "っともすずしく、さらになんの奇跡か、季節はずれというのにまだイチゴが食", n)
				match, err := unknown.Match(profileInstances)
				So(err, ShouldEqual, nil)
				So(match, ShouldEqual, "jp")
			})

			Convey("TH sample text should score as TH", func() {
				unknown := NewProfileFromText("unknown", "ผู้ฟัง ของเขา รู้ว่า เขากำลังจะ พูดแบบนี้ พวกเขารู้มาก", n)
				match, err := unknown.Match(profileInstances)
				So(err, ShouldEqual, nil)
				So(match, ShouldEqual, "th")
			})
		})
	})
}
