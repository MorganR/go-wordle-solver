package cmd

import (
	"fmt"

	gws "github.com/MorganR/go-wordle-solver/lib"
)

func printGuesses(turns []gws.TurnData) {
	for i, td := range turns {
		fmt.Printf("\t%v: %s (%v remaining)\n", i+1, td.Guess, td.NumPossibleWordsBeforeGuess)
	}
}
