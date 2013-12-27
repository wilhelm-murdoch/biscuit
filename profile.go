// Package biscuit is used for simple linguistic computations.
package biscuit

import (
	"math"
	"strings"
)

// Profile is a structure we use to create an NGram data model. This stores all
// metadata associated with a processed corpus of text.
type Profile struct {
	label  string
	length float64
	n      int
	ngrams map[string]int
}

// NewProfile is a factory function which returns a processed instance of a
// Profile data model.
func NewProfile(label string, text string, n int) *Profile {
	p := new(Profile)

	p.n = n
	p.label = label
	p.ngrams = make(map[string]int)
	p.ParseFromText(text)
	p.Length()

	return p
}

// ParseFromText creates an ngram table from the specified text. This table
// is a map whose keys are a distinct set of n-length character sequences
// associated with their frequency.
func (p *Profile) ParseFromText(text string) {
	chars := make([]rune, 2*p.n)

	k := 0
	for _, chars[k] = range strings.Join(strings.Fields(text), " ") + " " {
		chars[p.n+k] = chars[k]
		k = (k + 1) % p.n
		p.ngrams[string(chars[k:k+p.n])]++
	}
}

// Length converts an ngram table into a vector and returns its magnitude. This
// will be used later when executing Match, or Subtract, as a vector based
// search
func (p *Profile) Length() float64 {
	length := 0.0
	for _, frequency := range p.ngrams {
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
	for sequence, frequency := range p.ngrams {
		if f, ok := profile.ngrams[sequence]; ok {
			total += float64(frequency * f)
		}
	}

	return total / (p.length * profile.length)
}

// Match returns the best possible match among the current profile instance
// and the specified argument array of profile instances.
func (p *Profile) Match(profiles ...*Profile) string {
	scores := make(map[string]float64)

	for _, profile := range profiles {
		scores[profile.label] = p.Subtract(profile)
	}

	return SortedKeys(scores)[0]
}
