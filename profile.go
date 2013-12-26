package biscuit

import (
	"math"
	"strings"
)

type Profile struct {
	label  string
	length float64
	n      int
	table  map[string]int
}

func NewProfile(label string, text string, n int) *Profile {
	p := new(Profile)

	p.n = n
	p.label = label
	p.table = make(map[string]int)
	p.Parse(text)
	p.Length()

	return p
}

func (p *Profile) Parse(text string) {
	chars := make([]rune, 2*p.n)

	k := 0
	for _, chars[k] = range strings.Join(strings.Fields(text), " ") + " " {
		chars[p.n+k] = chars[k]
		k = (k + 1) % p.n
		p.table[string(chars[k:k+p.n])]++
	}
}

func (p *Profile) Length() float64 {
	length := 0.0
	for _, frequency := range p.table {
		length += math.Pow(float64(frequency), 2)
	}

	p.length = math.Pow(length, 0.5)

	return p.length
}

func (p *Profile) Subtract(profile *Profile) float64 {
	total := 0.0
	for sequence, frequency := range p.table {
		if f, ok := profile.table[sequence]; ok {
			total += float64(frequency * f)
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
