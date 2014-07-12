package biscuit

import (
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func BenchmarkModelFromFile(b *testing.B) {
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
		NewModelFromFile(label, files[label], 3)
	}
}

func BenchmarkMatch(b *testing.B) {
	models := make(map[string]*Model)
	files, _ := filepath.Glob("./corpora/*.txt")

	var labels = make([]string, len(files))
	var modelInstances = make([]*Model, 0, len(models))

	for i, file := range files {
		label := filepath.Base(file)[:2]
		labels[i] = label
		model, _ := NewModelFromFile(label, file, 3)
		models[label] = model
		modelInstances = append(modelInstances, model)
	}

	for i := 0; i < b.N; i++ {
		label := labels[rand.Intn(len(labels))]
		models[label].MatchReturnBest(modelInstances)
	}
}

func TestModel(t *testing.T) {
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

	english := NewModelFromText(label, text, n)

	Convey("Subject: Creating models", t, func() {
		Convey("Given a label, some text and an ngram length", func() {
			Convey("The model label should equal the specified value", func() {
				So(english.Label, ShouldEqual, label)
			})

			Convey("The model text should be parsed into an ngram table", func() {
				for sequence, frequency := range testTable {
					if f, ok := english.Ngrams[sequence]; ok {
						So(frequency, ShouldEqual, f)
					}
				}
			})

			Convey("The model ngram length should equal the specified value", func() {
				So(english.N, ShouldEqual, n)
				for sequence := range english.Ngrams {
					So(len(sequence), ShouldEqual, n)
				}
			})
		})
	})

	Convey("Subject: Create ngram tables from files", t, func() {
		Convey("Given a corpora of text files", func() {
			Convey("Opening an invalid text file should yield an error", func() {
				_, err := NewModelFromFile("unknown", "/nothing/here.wut", 3)
				So(err, ShouldNotEqual, nil)
			})

			Convey("Opening a text file should yield a new biscuit.Model instance", func() {
				samples, _ := filepath.Glob("./corpora/*.txt")

				for _, file := range samples {
					p, err := NewModelFromFile(filepath.Base(file)[:2], file, 3)
					So(err, ShouldEqual, nil)
					So(filepath.Base(file)[:2], ShouldEqual, p.Label)
					So(p.length, ShouldEqual, p.Length())
				}
			})
		})
	})

	Convey("Subject: Scoring and comparing models", t, func() {
		Convey("Given a corpora of sample text in various languages", func() {
			corpora := make(map[string]string)
			models := make(map[string]*Model)
			samples, _ := filepath.Glob("./corpora/*.txt")

			for _, file := range samples {
				text, _ := ioutil.ReadFile(file)
				fileName := filepath.Base(file)
				label = fileName[0 : len(fileName)-len(filepath.Ext(fileName))]
				corpora[label] = string(text)
				models[label] = NewModelFromText(label, corpora[label], n)
			}

			modelInstances := make([]*Model, 0, len(models))

			for _, model := range models {
				modelInstances = append(modelInstances, model)
			}

			Convey("Subtracting one model from another with different ngram lengths should return zero(0)", func() {
				model1 := NewModelFromText("model1", "sup", 1)
				model2 := NewModelFromText("model2", "hey", 2)

				difference := model1.Subtract(model2)

				So(difference, ShouldEqual, 0)
			})

			Convey("Matching models of different ngram lengths should raise an error", func() {
				unknown := NewModelFromText("unknown", "sup", n+1)
				match, err := unknown.MatchReturnBest(modelInstances)
				So(match, ShouldEqual, "")
				So(err, ShouldNotEqual, nil)
			})

			Convey("Comparing a corpus against itself should yield an exact match", func() {
				for _, model := range models {
					difference := model.Subtract(model)
					So(Round(difference, 1), ShouldBeGreaterThanOrEqualTo, 1)
				}
			})

			Convey("Comparing a corpus of one language against another should yield partial match", func() {
				difference := models["english"].Subtract(models["german"])
				So(difference, ShouldBeGreaterThanOrEqualTo, 0)
				So(difference, ShouldBeLessThan, 1)
			})

			Convey("Comparing a corpus of one language against another should yield the same score regardless of order", func() {
				difference1 := models["german"].Subtract(models["english"])
				difference2 := models["english"].Subtract(models["german"])

				So(difference1, ShouldEqual, difference2)
			})

			Convey("The vectors of each model should be properly calculated", func() {
				for _, model := range models {
					So(model.length, ShouldEqual, model.Length())
				}
			})

			Convey("DE sample text should score as GERMAN", func() {
				unknown := NewModelFromText("unknown", "Der Kanalinspektor natrlich!", n)
				match, err := unknown.MatchReturnBest(modelInstances)
				So(err, ShouldEqual, nil)
				So(match, ShouldEqual, "german")
			})

			Convey("ES sample text should score as SPANISH", func() {
				unknown := NewModelFromText("unknown", "con sus aperos formados con prendas de procedencia.", n)
				match, err := unknown.MatchReturnBest(modelInstances)
				So(err, ShouldEqual, nil)
				So(match, ShouldEqual, "spanish")
			})

			Convey("FR sample text should score as FRENCH", func() {
				unknown := NewModelFromText("unknown", "Monsieur le baron était un des plus puissants.", n)
				match, err := unknown.MatchReturnBest(modelInstances)
				So(err, ShouldEqual, nil)
				So(match, ShouldEqual, "french")
			})

			Convey("JP sample text should score as JAPANESE", func() {
				unknown := NewModelFromText("unknown", "っともすずしく、さらになんの奇跡か、季節はずれというのにまだイチゴが食", n)
				match, err := unknown.MatchReturnBest(modelInstances)
				So(err, ShouldEqual, nil)
				So(match, ShouldEqual, "japanese")
			})

			Convey("TH sample text should score as THAI", func() {
				unknown := NewModelFromText("unknown", "ผู้ฟัง ของเขา รู้ว่า เขากำลังจะ พูดแบบนี้ พวกเขารู้มาก", n)
				match, err := unknown.MatchReturnBest(modelInstances)
				So(err, ShouldEqual, nil)
				So(match, ShouldEqual, "thai")
			})
		})
	})
}
