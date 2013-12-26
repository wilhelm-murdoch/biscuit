package biscuit

type Ngram struct {
	sequence  string
	frequency int
}

func (n *Ngram) Increment() {
	n.frequency++
}

func NewNgram(sequence string) *Ngram {
	return &Ngram{sequence, 1}
}
