package main

import (
	"fmt"
	"math/rand"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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
	name     string
	hand     []Card
	strategy Strategy
}

type UnoGame struct {
	players      []Player
	drawPile     Deck
	discardPile  Deck
	direction    bool // true: clockwise, false: anti-clockwise
	activePlayer int
	forcedColor  Color
}

type ShuffleSettings struct {
	deckSize                  int
	repetitionCount           int
	preShuffleDivPos          float32
	preShuffleDivUncertainty  int
	doPostShuffleDiv          bool
	postShuffleDivPos         float32
	postShuffleDivUncertainty int
	riffleRange               int
}

func (d Deck) divide(p float32, u int) ([]Card, []Card) {
	cut := int(float32(len(d.cards))*p) - u/2 + rand.Intn(u)

	return d.cards[:cut], d.cards[cut:]
}
func (d *Deck) shuffle(s ShuffleSettings) {
	a, b := d.divide(s.preShuffleDivPos, s.preShuffleDivUncertainty)

	// fmt.Printf("Before: %v\n", d.cards)
	d.cards = nil

	for {
		n := min(len(a), rand.Intn(s.riffleRange))
		d.cards = append(d.cards, a[:n]...)
		a = a[n:]
		m := min(len(b), rand.Intn(s.riffleRange))
		d.cards = append(d.cards, b[:m]...)
		b = b[m:]
		if len(a) == 0 && len(b) == 0 {
			break
		}
	}

	if s.doPostShuffleDiv {
		a, b = d.divide(s.postShuffleDivPos, s.postShuffleDivUncertainty)
		d.cards = append(b, a...)
	}
}

func (g *UnoGame) generateNewDeck() {
	g.drawPile.cards = make([]Card, 108)

	i := 0
	// Generate number cards
	for c := Red; c <= Yellow; c++ {
		for v := 0; v <= 9; v++ {
			g.drawPile.cards[i] = Card{c, CardValue(v)}
			i++
		}
		for v := 1; v <= 9; v++ {
			g.drawPile.cards[i] = Card{c, CardValue(v)}
			i++
		}
	}
	// Generate color special cards
	for c := Red; c <= Yellow; c++ {
		for v := SkipCard; v <= DrawTwoCard; v++ {
			g.drawPile.cards[i] = Card{c, v}
			i++
			g.drawPile.cards[i] = Card{c, v}
			i++
		}
	}
	// Generate black cards
	for j := 0; j < 4; j++ {
		g.drawPile.cards[i] = Card{Black, WildCard}
		i++
		g.drawPile.cards[i] = Card{Black, WildDrawFourCard}
		i++
	}
}

func (g *UnoGame) playerDrawsCard(playerNumber int) {
	g.players[playerNumber].hand = append(g.players[playerNumber].hand, g.drawPile.cards[0])
	g.drawPile.cards = g.drawPile.cards[1:]

	fmt.Printf("Player %s: draws card: %+v\n", g.players[playerNumber].name, g.players[playerNumber].hand[len(g.players[playerNumber].hand)-1])

	if len(g.drawPile.cards) == 0 {
		g.drawPile.cards = append(g.drawPile.cards, g.discardPile.cards[1:]...)
		g.discardPile.cards = g.discardPile.cards[:1]

		fmt.Printf("Reshuffling discard pile: %d cards.\n", len(g.drawPile.cards))

		var s = ShuffleSettings{
			deckSize:                  len(g.drawPile.cards),
			repetitionCount:           1,
			preShuffleDivPos:          0.5,
			preShuffleDivUncertainty:  10,
			doPostShuffleDiv:          true,
			postShuffleDivPos:         0.5,
			postShuffleDivUncertainty: len(g.drawPile.cards),
			riffleRange:               2,
		}

		for i := 0; i < 5; i++ {
			g.drawPile.shuffle(s)
		}
	}
}

func (g *UnoGame) initialize(playerList []string) {
	g.players = make([]Player, len(playerList))
	for i := range g.players {
		g.players[i].name = playerList[i]
		g.players[i].strategy = StrategyRandom
	}

	g.generateNewDeck()

	var s = ShuffleSettings{
		deckSize:                  len(g.drawPile.cards),
		repetitionCount:           1,
		preShuffleDivPos:          0.5,
		preShuffleDivUncertainty:  10,
		doPostShuffleDiv:          true,
		postShuffleDivPos:         0.5,
		postShuffleDivUncertainty: len(g.drawPile.cards),
		riffleRange:               2,
	}

	for i := 0; i < 5; i++ {
		g.drawPile.shuffle(s)
	}

	for i := 0; i < 7; i++ {
		for j := range g.players {
			g.playerDrawsCard(j)
		}
	}

	g.discardPile.cards = append(g.discardPile.cards, g.drawPile.cards[0])
	g.drawPile.cards = g.drawPile.cards[1:]

	g.activePlayer = 0
}

func (g *UnoGame) getNextPlayerIndex() int {
	nextPlayerIndex := g.activePlayer
	if g.direction {
		nextPlayerIndex++
		if nextPlayerIndex == len(g.players) {
			nextPlayerIndex = 0
		}
	} else {
		nextPlayerIndex--
		if nextPlayerIndex < 0 {
			nextPlayerIndex = len(g.players) - 1
		}
	}

	return nextPlayerIndex
}

func (g *UnoGame) setNextPlayer() {
	g.activePlayer = g.getNextPlayerIndex()
}

func (g *UnoGame) getCardCandidates() []CardCandidate {
	topColor := g.discardPile.cards[0].color
	topValue := g.discardPile.cards[0].value
	var candidates []CardCandidate
	playerNumber := g.activePlayer

	for i, thisCard := range g.players[playerNumber].hand {
		if g.forcedColor < 0 && (thisCard.color == topColor || thisCard.value == topValue) {
			candidates = append(candidates, CardCandidate{i, 0})
		}
		if g.forcedColor > 0 && thisCard.color == g.forcedColor {
			candidates = append(candidates, CardCandidate{i, 0})
		}
		if thisCard.color == Black {
			candidates = append(candidates, CardCandidate{i, 0})
		}
	}

	return candidates
}

func (g *UnoGame) scoreCandidates(candidates []CardCandidate, s Strategy) []CardCandidate {
	return candidates
}

func (g *UnoGame) playOutCard(candidates []CardCandidate) bool {
	maxScore := -1
	for _, c := range candidates {
		if c.score > maxScore {
			maxScore = c.score
		}
	}
	for i, c := range candidates {
		if c.score < maxScore {
			candidates[i] = candidates[len(candidates)-1] // Copy last element to index i.
			candidates = candidates[:len(candidates)-1]   // Truncate slice.
		}
	}

	if len(candidates) == 0 {
		g.playerDrawsCard(g.activePlayer)
		return false
	}

	playedCardIndex := candidates[rand.Intn(len(candidates))].index
	playedCardValue := g.players[g.activePlayer].hand[playedCardIndex].value

	g.discardPile.cards = append([]Card{g.players[g.activePlayer].hand[playedCardIndex]}, g.discardPile.cards...)

	lenHand := len(g.players[g.activePlayer].hand)
	g.players[g.activePlayer].hand[playedCardIndex] = g.players[g.activePlayer].hand[lenHand-1] // Copy last element to index i.
	g.players[g.activePlayer].hand = g.players[g.activePlayer].hand[:lenHand-1]                 // Truncate slice.

	g.forcedColor = NoColor

	fmt.Printf("Player %s plays %+v\n", g.players[g.activePlayer].name, g.discardPile.cards[0])

	//TODO React to played Card
	switch playedCardValue {
	case SkipCard:
		g.setNextPlayer()
	case ReverseCard:
		g.direction = !g.direction
	case DrawTwoCard:
		g.playerDrawsCard(g.getNextPlayerIndex())
		g.playerDrawsCard(g.getNextPlayerIndex())
	case WildCard:
		g.forcedColor = Color(rand.Intn(len(candidates))) + 1
		fmt.Printf("Player %s: WildCard wish: %d\n", g.players[g.activePlayer].name, g.forcedColor)
	case WildDrawFourCard:
		g.playerDrawsCard(g.getNextPlayerIndex())
		g.playerDrawsCard(g.getNextPlayerIndex())
		g.playerDrawsCard(g.getNextPlayerIndex())
		g.playerDrawsCard(g.getNextPlayerIndex())
		g.forcedColor = Color(rand.Intn(len(candidates))) + 1
		fmt.Printf("Player %s: WildCard wish: %d\n", g.players[g.activePlayer].name, g.forcedColor)
	}

	switch lenHand {
	case 2:
		fmt.Printf("Player %s: UNO!\n", g.players[g.activePlayer].name)
		return false
	case 1:
		fmt.Printf("Player %s: WON!\n", g.players[g.activePlayer].name)
		return true
	default:
		return false
	}
}

func (g *UnoGame) playOneTurn() bool {
	candidates := g.getCardCandidates()
	candidates = g.scoreCandidates(candidates, StrategyRandom)
	if g.playOutCard(candidates) {
		return true
	}
	g.setNextPlayer()
	return false
}

func main() {
	rand.Seed(time.Now().UnixNano())

	g := UnoGame{
		players:      nil,
		drawPile:     Deck{nil},
		discardPile:  Deck{nil},
		direction:    true,
		activePlayer: 0,
		forcedColor:  -1,
	}
	g.initialize([]string{"A", "B", "C"})
	//fmt.Print(g)

	stopGame := false
	for i := 1; !stopGame; i++ {
		fmt.Printf("\nTurn %d:\n", i)
		stopGame = g.playOneTurn()
		if stopGame {
			break
		}
	}

	fmt.Printf("Game over, Player %s has won.", g.players[g.activePlayer].name)
}
