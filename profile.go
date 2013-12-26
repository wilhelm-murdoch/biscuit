package biscuit

import (
	"math"
	"strings"
)

type Profile struct {
	label  string
	length float64
	n      int
	table  []*Ngram
}

func NewProfile(label string, text string, n int) *Profile {
	p := new(Profile)

	p.n = n
	p.label = label
	p.Parse(text)
	p.Length()

	return p
}

func (p *Profile) FindNgram(sequence string) (*Ngram, bool) {
	for _, ngram := range p.table {
		if ngram.sequence == sequence {
			return ngram, true
		}
	}
	return nil, false
}

func (p *Profile) Parse(text string) {
	chars := []rune(strings.Repeat(" ", p.n))

	for _, letter := range strings.Join(strings.Fields(text), " ") + " " {
		chars = append(chars[1:], letter)

		sequence := string(chars)
		if ngram, ok := p.FindNgram(sequence); ok {
			ngram.Increment()
		} else {
			p.table = append(p.table, NewNgram(sequence))
		}
	}
}

func (p *Profile) Length() {
	length := 0.0
	for _, ngram := range p.table {
		length += math.Pow(float64(ngram.frequency), 2)
	}

	p.length = math.Pow(length, 0.5)
}

func (p *Profile) Subtract(profile *Profile) float64 {
	total := 0.0
	for _, ngram1 := range p.table {
		if ngram2, ok := profile.FindNgram(ngram1.sequence); ok {
			total += float64(ngram1.frequency * ngram2.frequency)
		}
	}

	return total / (p.length * profile.length)
}

func (p *Profile) Match(profiles ...*Profile) string {
	scores := make(map[string]float64)

	for _, profile := range profiles {
		scores[profile.label] = p.Subtract(profile)
	}

	return sortedKeys(scores)[0]
}
