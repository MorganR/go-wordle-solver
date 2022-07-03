package cmd

import (
	"fmt"

	gws "github.com/MorganR/go-wordle-solver/lib"
)

func printGuesses(data *gws.GameData) {
	for i, td := range data.Turns {
		fmt.Printf("\t%v: %s (%v remaining)\n", i+1, td.Guess, td.NumPossibleWordsBeforeGuess)
	}
}
