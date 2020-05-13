package gouno

import (
	"fmt"
	"log"
	"math/rand"
)

type UnoGame struct {
	Players      []Player
	drawPile     Deck
	discardPile  Deck
	direction    bool // true: clockwise, false: anti-clockwise
	activePlayer int
	forcedColor  Color
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
	g.Players[playerNumber].hand = append(g.Players[playerNumber].hand, g.drawPile.cards[0])
	g.drawPile.cards = g.drawPile.cards[1:]

	log.Printf("Player %s draws card: %+v\n", g.Players[playerNumber].name, g.Players[playerNumber].hand[len(g.Players[playerNumber].hand)-1])

	if len(g.drawPile.cards) == 0 {
		g.drawPile.cards = append(g.drawPile.cards, g.discardPile.cards[1:]...)
		g.discardPile.cards = g.discardPile.cards[:1]

		log.Printf("Reshuffling discard pile: %d cards.\n", len(g.drawPile.cards))

		var s = ShuffleSettings{
			DeckSize:                  len(g.drawPile.cards),
			RepetitionCount:           1,
			PreShuffleDivPos:          0.5,
			PreShuffleDivUncertainty:  10,
			DoPostShuffleDiv:          true,
			PostShuffleDivPos:         0.5,
			PostShuffleDivUncertainty: len(g.drawPile.cards) - 1,
			RiffleRange:               2,
		}

		for i := 0; i < 5; i++ {
			g.drawPile.shuffle(s)
		}
	}
}

func (g *UnoGame) Initialize(playerList []string) {
	g.Players = make([]Player, len(playerList))
	for i := range g.Players {
		g.Players[i].name = playerList[i]
		g.Players[i].strategy = StrategyRandom
		g.Players[i].colorKept = NoColor
	}

	g.generateNewDeck()

	var s = ShuffleSettings{
		DeckSize:                  len(g.drawPile.cards),
		RepetitionCount:           1,
		PreShuffleDivPos:          0.5,
		PreShuffleDivUncertainty:  10,
		DoPostShuffleDiv:          true,
		PostShuffleDivPos:         0.5,
		PostShuffleDivUncertainty: len(g.drawPile.cards) - 1,
		RiffleRange:               2,
	}

	for i := 0; i < 5; i++ {
		g.drawPile.shuffle(s)
	}

	for i := 0; i < 7; i++ {
		for j := range g.Players {
			g.playerDrawsCard(j)
		}
	}

	g.discardPile.cards = append(g.discardPile.cards, g.drawPile.cards[0])
	g.drawPile.cards = g.drawPile.cards[1:]

	log.Printf("First card: %+v\n", g.discardPile.cards[0])

	g.activePlayer = 0
}

func (g *UnoGame) getNextPlayerIndex() int {
	nextPlayerIndex := g.activePlayer
	if g.direction {
		nextPlayerIndex++
		if nextPlayerIndex == len(g.Players) {
			nextPlayerIndex = 0
		}
	} else {
		nextPlayerIndex--
		if nextPlayerIndex < 0 {
			nextPlayerIndex = len(g.Players) - 1
		}
	}

	return nextPlayerIndex
}

func (g *UnoGame) setNextPlayer() {
	if len(g.Players[g.activePlayer].hand) > 0 {
		g.activePlayer = g.getNextPlayerIndex()
	}
}

func (g *UnoGame) getCardCandidates() []CardCandidate {
	topColor := g.discardPile.cards[0].color
	topValue := g.discardPile.cards[0].value
	var candidates []CardCandidate
	playerNumber := g.activePlayer

	for i, thisCard := range g.Players[playerNumber].hand {
		if g.forcedColor == 0 && (thisCard.color == topColor || thisCard.value == topValue) {
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

func (g *UnoGame) playOutCard(candidates []CardCandidate) bool {
	maxScore := -1
	for _, c := range candidates {
		if c.score > maxScore {
			maxScore = c.score
		}
	}
	var candidateIndices = make([]int, 0)
	for _, c := range candidates {
		if c.score == maxScore {
			candidateIndices = append(candidateIndices, c.index)
		}
	}

	if len(candidateIndices) == 0 {
		g.playerDrawsCard(g.activePlayer)
		candidates = g.getCardCandidates()
		if len(candidates) == 0 {
			return false
		} else {
			candidateIndices = append(candidateIndices, candidates[0].index)
		}
	}

	playedCardIndex := candidateIndices[rand.Intn(len(candidateIndices))]
	playedCardValue := g.Players[g.activePlayer].hand[playedCardIndex].value

	g.discardPile.cards = append([]Card{g.Players[g.activePlayer].hand[playedCardIndex]}, g.discardPile.cards...)

	lenHand := len(g.Players[g.activePlayer].hand)
	g.Players[g.activePlayer].hand[playedCardIndex] = g.Players[g.activePlayer].hand[lenHand-1] // Copy last element to index i.
	g.Players[g.activePlayer].hand = g.Players[g.activePlayer].hand[:lenHand-1]                 // Truncate slice.
	lenHand--

	g.forcedColor = NoColor

	if lenHand == 1 {
		log.Printf("Player %s: UNO!\n", g.GetActivePlayerName())
	}
	log.Printf("Player %s plays %+v\n", g.GetActivePlayerName(), g.discardPile.cards[0])

	switch playedCardValue {
	case SkipCard:
		g.setNextPlayer() // skip next player
	case ReverseCard:
		g.direction = !g.direction
		if len(g.Players) == 2 {
			g.setNextPlayer() // skip next player
		}
	case DrawTwoCard:
		g.playerDrawsCard(g.getNextPlayerIndex())
		g.playerDrawsCard(g.getNextPlayerIndex())
		g.setNextPlayer() // skip next player
	case WildCard:
		g.forcedColor = g.Players[g.activePlayer].getStrongestColor()
		log.Printf("Player %s: WildCard wish: %s\n", g.GetActivePlayerName(), colorNames[g.forcedColor])
	case WildDrawFourCard:
		g.playerDrawsCard(g.getNextPlayerIndex())
		g.playerDrawsCard(g.getNextPlayerIndex())
		g.playerDrawsCard(g.getNextPlayerIndex())
		g.playerDrawsCard(g.getNextPlayerIndex())
		g.forcedColor = g.Players[g.activePlayer].getStrongestColor()
		log.Printf("Player %s: WildCard wish: %s\n", g.GetActivePlayerName(), colorNames[g.forcedColor])
		g.setNextPlayer() // skip next player
	}

	switch lenHand {
	case 0:
		return true
	default:
		return false
	}
}

func (g *UnoGame) getHumanChoice() []CardCandidate {
	log.Printf("%s's hand:\n", g.GetActivePlayerName())
	log.Printf("%s\n", g.Players[g.activePlayer].handToString())

	var choice int
	fmt.Scanf("%d", &choice)
	return []CardCandidate{{choice, 1}}
}

func (g *UnoGame) PlayOneTurn() bool {
	// log.Printf("Hands: %s(%d), %s(%d), %s(%d)\n", g.players[0].name, len(g.players[0].hand), g.players[1].name, len(g.players[1].hand), g.players[2].name, len(g.players[2].hand))
	candidates := g.getCardCandidates()
	if g.Players[g.activePlayer].Human {
		candidates = g.getHumanChoice()
	} else {
		candidates = g.scoreCandidates(candidates, g.Players[g.activePlayer].strategy)
	}

	if g.playOutCard(candidates) {
		return true
	}
	g.setNextPlayer()
	return false
}

func (g *UnoGame) GetActivePlayerIndex() int {
	return g.activePlayer
}

func (g *UnoGame) GetActivePlayerName() string {
	return g.Players[g.activePlayer].name
}
