package gouno

import "math/rand"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type ShuffleSettings struct {
	DeckSize                  int
	RepetitionCount           int
	PreShuffleDivPos          float32
	PreShuffleDivUncertainty  int
	DoPostShuffleDiv          bool
	PostShuffleDivPos         float32
	PostShuffleDivUncertainty int
	RiffleRange               int
}

func (d Deck) divide(p float32, u int) ([]Card, []Card) {
	cut := int(float32(len(d.cards))*p) - u/2 + rand.Intn(u)

	return d.cards[:cut], d.cards[cut:]
}
func (d *Deck) shuffle(s ShuffleSettings) {
	a, b := d.divide(s.PreShuffleDivPos, s.PreShuffleDivUncertainty)

	// log.Printf("Before: %v\n", d.cards)
	d.cards = nil

	for {
		n := min(len(a), rand.Intn(s.RiffleRange))
		d.cards = append(d.cards, a[:n]...)
		a = a[n:]
		m := min(len(b), rand.Intn(s.RiffleRange))
		d.cards = append(d.cards, b[:m]...)
		b = b[m:]
		if len(a) == 0 && len(b) == 0 {
			break
		}
	}

	if s.DoPostShuffleDiv {
		a, b = d.divide(s.PostShuffleDivPos, s.PostShuffleDivUncertainty)
		d.cards = append(b, a...)
	}
}
