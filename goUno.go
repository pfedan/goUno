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
}

var gameParams unoGameParameters

func init() {

	playersDefault := []string{"A", "B"}
	humanFlagDefault := []bool{false, false}
	roundCountDefault := 1

	flag.StringSliceVarP(&(gameParams.playerNames), "players", "p", playersDefault, "")
	flag.BoolSliceVarP(&(gameParams.humanFlags), "human", "h", humanFlagDefault, "")
	flag.IntVarP(&(gameParams.roundCount), "roundcount", "r", roundCountDefault, "")
}

func main() {
	flag.Parse()
	if len(gameParams.humanFlags) > len(gameParams.playerNames) {
		os.Exit(-100)
	}

	rand.Seed(time.Now().UnixNano())

	log.SetFlags(0)
	// log.SetOutput(ioutil.Discard)

	var g gouno.UnoGame

	players := gameParams.playerNames

	cnt := map[string]int{}
	turns := map[int]int{}

	for round := 0; round < gameParams.roundCount; round++ {
		g.Initialize(players)
		for i, humanFlag := range gameParams.humanFlags {
			g.Players[i].Human = humanFlag
		}

		stopGame := false
		for i := 1; !stopGame; i++ {
			log.Printf("Turn %d:\n", i)
			stopGame = g.PlayOneTurn()
			if stopGame {
				turns[i]++
				break
			}
		}

		cnt[g.GetActivePlayerName()]++

		log.Printf("Game over, Player %s has won.\n\n", g.GetActivePlayerName())
	}

	log.SetOutput(os.Stderr)

	log.Printf("%+v\n", cnt)
	log.Printf("%+v\n", turns)
}
