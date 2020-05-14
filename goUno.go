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
	pointSum := map[string]int{}
	pointsPerGame := map[int]int{}
	turns := map[int]int{}

	for round := 0; round < gameParams.roundCount; round++ {
		g.Muted = gameParams.muteLog

		g.Initialize(players)
		for i, humanFlag := range gameParams.humanFlags {
			g.Players[i].Human = humanFlag
			// g.Players[i].Strategy = gouno.StrategyAggressive
		}

		stopGame := false
		for i := 1; !stopGame; i++ {
			g.Printf("Turn %d:\n", i)
			stopGame = g.PlayOneTurn()
			if stopGame {
				turns[g.Turns]++
				break
			}
		}

		cnt[g.GetActivePlayerName()]++

		points := 0
		for _, p := range g.Players {
			points += p.GetPoints()
		}
		pointsPerGame[points]++
		pointSum[g.GetActivePlayerName()] += points

		g.Printf("Game over, Player %s has won with %d points.\n\n", g.GetActivePlayerName(), points)
	}

	log.Printf("Wins per Player: %+v\n", cnt)
	log.Printf("Total points per Player: %+v\n", pointSum)
	log.Printf("Count of turns per game: %+v\n", turns)
	log.Printf("Count of points per game: %+v\n", pointsPerGame)
}
