// Package biscuit is used for simple linguistic computations.
package biscuit

import (
	"errors"
	"io/ioutil"
	"math"
	"strings"
	"sync"
)

// Model is a structure we use to create an NGram data model. This stores all
// metadata associated with a processed corpus of text.
type Model struct {
	Label  string
	length float64
	N      int
	Ngrams map[string]int
}

// NewModelFromText is a factory function which returns a processed instance of a
// Model data model.
func NewModelFromText(label string, text string, n int) *Model {
	m := new(Model)

	m.N = n
	m.Label = label
	m.Ngrams = make(map[string]int)
	m.ParseTextToNgramTable(text)
	m.Length()

	return m
}

// NewModelFromFile attempts to open a file at the specified path and parse its
// contents as text using NewModelFromText. This is nothing more than a
// convienence method.
func NewModelFromFile(label string, filepath string, n int) (*Model, error) {
	bytes, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	return NewModelFromText(label, string(bytes), n), nil
}

// ParseTextToNgramTable creates an ngram table from the specified text. This
// table is a map whose keys are a distinct set of n-length character sequences
// associated with their frequency.
func (m *Model) ParseTextToNgramTable(text string) {
	chars := make([]rune, 2*m.N)

	k := 0
	for _, chars[k] = range strings.Join(strings.Fields(text), " ") + " " {
		chars[m.N+k] = chars[k]
		k = (k + 1) % m.N
		m.Ngrams[string(chars[k:k+m.N])]++
	}
}

// Length converts an ngram table into a vector and returns its magnitude. This
// will be used later when executing Match, or Subtract, as a vector based
// search
func (m *Model) Length() float64 {
	length := 0.0
	for _, frequency := range m.Ngrams {
		length += math.Pow(float64(frequency), 2)
	}

	m.length = math.Pow(length, 0.5)

	return m.length
}

// Subtract attempts to determine the difference between the specified model
// and the current model instance. This is done by using the angle between
// the two vector lengths and determining their cosine. This results in a float
// between 1 and 0. The closer the return value is to 1, the better the match.
func (m *Model) Subtract(model *Model) float64 {
	if m.N != model.N {
		return 0
	}

	total := 0.0
	for sequence, frequency := range m.Ngrams {
		if f, ok := model.Ngrams[sequence]; ok {
			total += float64(frequency * f)
		}
	}

	return total / (m.length * model.length)
}

// MatchReturnAll returns a sorted map of all resulting scores against the
// specified model instances. This will also return an error if any of
// the associated models' ngram lengths differ.
func (m *Model) MatchReturnAll(models []*Model) ([]string, map[string]float64, error) {
	scores := make(map[string]float64)

	var wg sync.WaitGroup

	for _, model := range models {
		if m.N != model.N {
			return nil, nil, errors.New("All models must be of the same ngram length.")
		}

		wg.Add(1)

		go func(model *Model) {
			defer wg.Done()

			scores[model.Label] = m.Subtract(model)
		}(model)
	}

	wg.Wait()

	return SortedKeys(scores), scores, nil
}

// MatchReturnBest returns the best possible match among the current model
// instance and the specified argument array of model instances. This will
// return any errors bubbled up from MatchReturnAll.
func (m *Model) MatchReturnBest(models []*Model) (string, error) {
	matches, _, err := m.MatchReturnAll(models)

	if err != nil {
		return "", err
	}

	return matches[0], nil
}
