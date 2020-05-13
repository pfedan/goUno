package gouno

import (
	"fmt"
	"strings"
)

type Strategy int

const (
	StrategyRandom      Strategy = 1
	StrategyKeepColor   Strategy = 2
	StrategyChangeColor Strategy = 4
	StrategyAggressive  Strategy = 8
)

type CardCandidate struct {
	index int // index on player's hand
	score int // highest score will be played
}

type Player struct {
	name      string
	hand      []Card
	strategy  Strategy
	colorKept Color
	Human     bool
}

func (p *Player) getStrongestColor() Color {
	cnt := make([]int, 4)

	for _, el := range p.hand {
		if el.color >= 1 {
			cnt[el.color-1]++
		}
	}

	max := 0
	maxIndex := -1
	for i, value := range cnt {
		if value > max {
			max = value
			maxIndex = i
		}
	}

	return Color(maxIndex + 1)
}

func (g *UnoGame) scoreCandidates(candidates []CardCandidate, s Strategy) []CardCandidate {
	switch s {
	case StrategyAggressive:
		for i, c := range candidates {
			if g.Players[g.activePlayer].hand[c.index].value >= 100 {
				candidates[i].score++
			}
		}
	case StrategyChangeColor:
	case StrategyKeepColor:
		if g.Players[g.activePlayer].colorKept == NoColor {
			g.Players[g.activePlayer].colorKept = g.Players[g.activePlayer].getStrongestColor()
		}
		cnt := 0
		for i, c := range candidates {
			if g.Players[g.activePlayer].hand[c.index].color == g.Players[g.activePlayer].colorKept {
				candidates[i].score++
				cnt++
			}
		}
		if cnt == 0 {
			g.Players[g.activePlayer].colorKept = NoColor
		}
	case StrategyRandom:
	default:
	}
	return candidates
}

func (p *Player) handToString() string {
	var out strings.Builder
	for i, c := range p.hand {
		out.WriteString(fmt.Sprintf("(%d): %v", i+1, c)) // print from 1
		if i < len(p.hand)-1 {
			out.WriteString(", ")
		}
	}

	return out.String()
}
