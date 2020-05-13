package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	gouno "github.com/pfedan/goUno/goUno"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	log.SetFlags(0)
	// log.SetOutput(ioutil.Discard)

	var g gouno.UnoGame

	players := []string{"Christin", "Daniel", "Julia"}

	cnt := map[string]int{}
	turns := map[int]int{}

	for round := 0; round < 1; round++ {
		g.Initialize(players)

		g.Players[1].Human = true

		// g.Players[2].strategy = StrategyAggressive + StrategyKeepColor
		// g.Players[2].strategy = StrategyKeepColor

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
