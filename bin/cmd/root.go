package cmd

import (
	"fmt"
	"os"
	"time"

	gws "github.com/MorganR/go-wordle-solver/lib"
	"github.com/spf13/cobra"
)

var WordBankPath string
var Guesser string

var validGuessers [2]string = [2]string{"random", "max_eliminations"}

var wordBank gws.WordBank
var guesser gws.Guesser

var rootCmd = &cobra.Command{
	Use:   "gws",
	Short: "Go Wordle Solver is a tool for solving Wordle puzzles algorithmically.",
	Long: `Go Wordle Solver (GWS) is a tool for easily building and running different Wordle solving
algorithms, and using them to solve individual puzzles, or benchmarking them against many puzzles at
once.`,
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&WordBankPath, "word_bank", "w", "../data/improved-words.txt", "Path to a list of words to use as the word bank.")
	rootCmd.PersistentFlags().StringVarP(&Guesser, "guesser", "g", "max_eliminations", fmt.Sprintf("The guessing algorithm to use. Options: %s.", validGuessers))

	start := time.Now()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Printf("Command took %s.\n", elapsed)
}

func initRoot() {
	err := initWordBank()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	err = initGuesser()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initWordBank() error {
	f, err := os.Open(WordBankPath)
	if err != nil {
		return err
	}
	wordBank, err = gws.WordBankFromReader(f)
	return err
}

func initGuesser() error {
	switch Guesser {
	case "random":
		g := gws.InitRandomGuesser(&wordBank)
		guesser = &g
	case "max_eliminations":
		scorer, err := gws.InitMaxEliminationsScorer(&wordBank)
		if err != nil {
			return err
		}
		g := gws.InitMaxScoreGuesser(&wordBank, &scorer, gws.GuessModeAll)
		guesser = &g
	default:
		return fmt.Errorf("Did not recognize guesser type %s. Accepted options: %s", Guesser, validGuessers)
	}
	return nil
}
