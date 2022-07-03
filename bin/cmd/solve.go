package cmd

import (
	"fmt"
	"os"
	"time"

	gws "github.com/MorganR/go-wordle-solver/lib"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(solveCmd)
}

var solveCmd = &cobra.Command{
	Use:   "solve",
	Short: "Solves a single Wordle puzzle.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		objective := gws.WordFromString(args[0])

		if objective.Len() != int(wordBank.WordLength()) {
			fmt.Fprintf(os.Stderr, "The objective word's length (%v) must match the word bank (%v).\n", objective.Len(), wordBank.WordLength())
			os.Exit(1)
			return
		}
		pw := wordBank.Words()
		hasWord := false
		for i := 0; i < pw.Len(); i++ {
			if pw.At(i).Equal(objective) {
				hasWord = true
				break
			}
		}
		if !hasWord {
			fmt.Fprintf(os.Stderr, "The objective word (%s) is not present in the word bank. Try another.\n", objective)
			os.Exit(1)
			return
		}

		start := time.Now()
		guesser := gws.InitRandomGuesser(&wordBank)
		result := gws.PlayGameWithGuesser(objective, 128, &guesser)
		end := time.Now()
		elapsed := end.Sub(start)
		switch result.Status {
		case gws.GameSuccess:
			fmt.Printf("Solved! It took me %v guesses.\n", len(result.Data.Turns))
			printGuesses(result.Data)
		case gws.GameFailure:
			fmt.Println("Failed :( I couldn't guess the word within the guess limit.")
			printGuesses(result.Data)
		case gws.UnknownWord:
			// This should be impossible since the word is already verified to be in the word bank.
			fmt.Println("Internal error.")
			os.Exit(1)
		}
		fmt.Printf("Guessing took %s.\n", elapsed)
	},
}
