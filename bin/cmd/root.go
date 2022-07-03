package cmd

import (
	"fmt"
	"os"

	gws "github.com/MorganR/go-wordle-solver/lib"
	"github.com/spf13/cobra"
)

var WordBankPath string
var wordBank gws.WordBank

var rootCmd = &cobra.Command{
	Use:   "gws",
	Short: "Go Wordle Solver is a tool for solving Wordle puzzles algorithmically.",
	Long: `Go Wordle Solver (GWS) is a tool for easily building and running different Wordle solving
algorithms, and using them to solve individual puzzles, or benchmarking them against many puzzles at
once.`,
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&WordBankPath, "word_bank", "w", "../data/improved-words.txt", "Path to a list of words to use as the word bank.")
	f, err := os.Open(WordBankPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	wordBank, err = gws.WordBankFromReader(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
