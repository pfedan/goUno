package gouno

import (
	"fmt"
)

type Color int

const (
	NoColor Color = 0
	Red     Color = 1
	Green   Color = 2
	Blue    Color = 3
	Yellow  Color = 4
	Black   Color = -1
)

var colorNames = map[Color]string{
	Black:   `Black`,
	Red:     `Red`,
	Green:   `Green`,
	Blue:    `Blue`,
	Yellow:  `Yellow`,
	NoColor: `NoColor`,
}

type CardType int

const (
	SkipCard         CardType = 100
	ReverseCard      CardType = 101
	DrawTwoCard      CardType = 102
	WildCard         CardType = 103
	WildDrawFourCard CardType = 104
)

var valueNames = map[CardType]string{
	0: `0`, 1: `1`, 2: `2`, 3: `3`, 4: `4`, 5: `5`, 6: `6`, 7: `7`, 8: `8`, 9: `9`,
	SkipCard:         `SkipCard`,
	ReverseCard:      `ReverseCard`,
	DrawTwoCard:      `DrawTwoCard`,
	WildCard:         `WildCard`,
	WildDrawFourCard: `WildDrawFourCard`,
}

type Card struct {
	color Color
	face  CardType
}

func (c Card) points() int {
	switch val := c.face; {
	case val == SkipCard:
		fallthrough
	case val == ReverseCard:
		fallthrough
	case val == DrawTwoCard:
		return 20
	case val == WildCard:
		return 50
	case val == WildDrawFourCard:
		return 70
	default:
		return int(val)
	}
}

func (c Card) String() string {
	return fmt.Sprintf("%s (%s)", valueNames[c.face], colorNames[c.color])
}

type Deck struct {
	cards []Card
}
