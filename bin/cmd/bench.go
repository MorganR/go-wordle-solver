package cmd

import (
	"fmt"
	"os"
	"time"

	gws "github.com/MorganR/go-wordle-solver/lib"
	"github.com/spf13/cobra"
)

const (
	maxGuesses int = 128
)

var BenchListPath string

func init() {
	rootCmd.AddCommand(benchCmd)
	benchCmd.LocalFlags().StringVarP(&BenchListPath, "bench_list", "b", "../data/1000-improved-words-shuffled.txt", "Path to a list of objective words to benchmark this algorithm against.")
}

var benchCmd = &cobra.Command{
	Use:   "bench",
	Short: "Benchmarks an algorithm against a given word list.",
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := os.Open(BenchListPath)
		if err != nil {
			return err
		}
		benchBank, err := gws.WordBankFromReader(f)
		if err != nil {
			return err
		}
		benchWords := benchBank.Words()

		start := time.Now()
		guesser := gws.InitRandomGuesser(&wordBank)
		benchLen := benchWords.Len()
		countNumGuesses := make([]int, maxGuesses)
		for i := 0; i < benchLen; i++ {
			objective := benchWords.At(i)
			result, err := gws.PlayGameWithGuesser(objective, maxGuesses, &guesser)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while guessing %s during benchmark.\n", objective)
				return err
			}
			if result.Status != gws.GameSuccess {
				fmt.Fprintf(os.Stderr, "Failed to guess %s during benchmark.\n", objective)
				printGuesses(result.Turns)
				os.Exit(1)
				return nil
			}
			numGuesses := len(result.Turns)
			countNumGuesses[numGuesses-1] = countNumGuesses[numGuesses-1] + 1
		}
		end := time.Now()
		elapsed := end.Sub(start)
		fmt.Printf("Benchmark completed in %s.\n", elapsed)

		maxGuessIndex := findLastNonZeroIndex(countNumGuesses)
		countNumGuesses = countNumGuesses[0:maxGuessIndex]
		printNumGuessResults(countNumGuesses)

		return nil
	},
}

func findLastNonZeroIndex(counts []int) int {
	for i := len(counts) - 1; i >= 0; i-- {
		if counts[i] != 0 {
			return i
		}
	}
	return 0
}

func printNumGuessResults(counts []int) {
	fmt.Println("Num guesses | Count")
	totalGuesses := 0
	runningTotal := 0
	for i, n := range counts {
		fmt.Println("--|---")
		fmt.Printf("%v | %v\n", i+1, n)
		totalGuesses += n
		runningTotal += (i + 1) * n
	}
	fmt.Printf("Average: %v +/- %v\n", float32(runningTotal)/float32(totalGuesses), 0)
}
