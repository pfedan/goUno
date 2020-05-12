package gouno

import (
	"fmt"
)

type Color int

const (
	Black   Color = 0
	Red     Color = 1
	Green   Color = 2
	Blue    Color = 3
	Yellow  Color = 4
	NoColor Color = -1
)

var colorNames = map[Color]string{
	Black:   `Black`,
	Red:     `Red`,
	Green:   `Green`,
	Blue:    `Blue`,
	Yellow:  `Yellow`,
	NoColor: `NoColor`,
}

type CardValue int

const (
	SkipCard         CardValue = 100
	ReverseCard      CardValue = 101
	DrawTwoCard      CardValue = 102
	WildCard         CardValue = 103
	WildDrawFourCard CardValue = 104
)

var valueNames = map[CardValue]string{
	0: `0`, 1: `1`, 2: `2`, 3: `3`, 4: `4`, 5: `5`, 6: `6`, 7: `7`, 8: `8`, 9: `9`,
	SkipCard:         `SkipCard`,
	ReverseCard:      `ReverseCard`,
	DrawTwoCard:      `DrawTwoCard`,
	WildCard:         `WildCard`,
	WildDrawFourCard: `WildDrawFourCard`,
}

type Card struct {
	color Color
	value CardValue
}

func (c Card) String() string {
	return fmt.Sprintf("%s (%s)", valueNames[c.value], colorNames[c.color])
}

type Deck struct {
	cards []Card
}
