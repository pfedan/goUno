package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	gouno "github.com/pfedan/goUno/goUno"
	flag "github.com/spf13/pflag"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type unoGameParameters struct {
	playerNames []string
	humanFlags  []bool
	roundCount  int
	muteLog     bool
}

var gameParams unoGameParameters

func init() {

	playersDefault := []string{"A", "B"}
	humanFlagDefault := []bool{false, false}
	roundCountDefault := 1
	muteLogDefault := false

	flag.StringSliceVarP(&(gameParams.playerNames), "players", "p", playersDefault, "")
	flag.BoolSliceVarP(&(gameParams.humanFlags), "human", "h", humanFlagDefault, "")
	flag.IntVarP(&(gameParams.roundCount), "roundcount", "r", roundCountDefault, "")
	flag.BoolVarP(&(gameParams.muteLog), "mute", "m", muteLogDefault, "")
}

func main() {
	flag.Parse()
	if len(gameParams.humanFlags) > len(gameParams.playerNames) {
		os.Exit(-100)
	}

	rand.Seed(time.Now().UnixNano())

	log.SetFlags(0)

	var g gouno.UnoGame

	players := gameParams.playerNames

	cnt := map[string]int{}
	turns := map[int]int{}

	for round := 0; round < gameParams.roundCount; round++ {
		g.Muted = gameParams.muteLog

		g.Initialize(players)
		for i, humanFlag := range gameParams.humanFlags {
			g.Players[i].Human = humanFlag
		}

		stopGame := false
		for i := 1; !stopGame; i++ {
			g.Printf("Turn %d:\n", i)
			stopGame = g.PlayOneTurn()
			if stopGame {
				turns[i]++
				break
			}
		}

		cnt[g.GetActivePlayerName()]++

		g.Printf("Game over, Player %s has won.\n\n", g.GetActivePlayerName())
	}

	log.Printf("Wins per Player: %+v\n", cnt)
	log.Printf("Count of turns per game: %+v\n", turns)
}
