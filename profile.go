// Package biscuit is used for simple linguistic computations.
package biscuit

import (
	"io/ioutil"
	"math"
	"strings"
)

// Profile is a structure we use to create an NGram data model. This stores all
// metadata associated with a processed corpus of text.
type Profile struct {
	Label  string
	length float64
	N      int
	Ngrams map[string]int
}

// NewProfileFromText is a factory function which returns a processed instance of a
// Profile data model.
func NewProfileFromText(label string, text string, n int) *Profile {
	p := new(Profile)

	p.N = n
	p.Label = label
	p.Ngrams = make(map[string]int)
	p.ParseTextToNgramTable(text)
	p.Length()

	return p
}

// NewProfileFromFile attempts to open a file at the specified path and parse its
// contents as text using NewProfileFromText. This is nothing more than a
// convienence method.
func (p *Profile) NewProfileFromFile(label string, filepath string, n int) *Profile {
	bytes, err := ioutil.ReadFile(filepath)

	if err != nil {
		panic(err)
	}

	return NewProfileFromText(label, string(bytes), n)
}

// NewProfileFromNgramCSV can be used to speed up the calculation process by
// specifying a precalculated ngram table stored in a CSV file. This method
// simply buffers the file and creates the table in memory on-the-fly.
func (p *Profile) NewProfileFromNgramCSV(label string, filepath string, n int) *Profile {
	// bytes, err := ioutil.ReadFile(filepath)

	// if err != nil {
	// 	panic(err)
	// }

	return nil
}

// Parse creates an ngram table from the specified text. This table
// is a map whose keys are a distinct set of n-length character sequences
// associated with their frequency.
func (p *Profile) ParseTextToNgramTable(text string) {
	chars := make([]rune, 2*p.N)

	k := 0
	for _, chars[k] = range strings.Join(strings.Fields(text), " ") + " " {
		chars[p.N+k] = chars[k]
		k = (k + 1) % p.N
		p.Ngrams[string(chars[k:k+p.N])]++
	}
}

// Length converts an ngram table into a vector and returns its magnitude. This
// will be used later when executing Match, or Subtract, as a vector based
// search
func (p *Profile) Length() float64 {
	length := 0.0
	for _, frequency := range p.Ngrams {
		length += math.Pow(float64(frequency), 2)
	}

	p.length = math.Pow(length, 0.5)

	return p.length
}

// Subtract attempts to determine the difference between the specified profile
// and the current profile instance. This is done by using the angle between
// the two vector lengths and determining their cosine. This results in a float
// between 1 and 0. The closer the return value is to 1, the better the match.
func (p *Profile) Subtract(profile *Profile) float64 {
	total := 0.0
	for sequence, frequency := range p.Ngrams {
		if f, ok := profile.Ngrams[sequence]; ok {
			total += float64(frequency * f)
		}
	}

	return total / (p.length * profile.length)
}

// Match returns the best possible match among the current profile instance
// and the specified argument array of profile instances.
func (p *Profile) Match(profiles []*Profile) string {
	scores := make(map[string]float64)

	for _, profile := range profiles {
		scores[profile.Label] = p.Subtract(profile)
	}

	return SortedKeys(scores)[0]
}
