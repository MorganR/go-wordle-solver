package cmd

import (
	"fmt"
	"os"
	"runtime"
	"time"

	gws "github.com/MorganR/go-wordle-solver/lib"
	"github.com/spf13/cobra"
)

const (
	maxGuesses int = 128
)

var BenchListPath string

var maxThreads = runtime.NumCPU()

func init() {
	rootCmd.AddCommand(benchCmd)
	benchCmd.LocalFlags().StringVarP(&BenchListPath, "bench_list", "b", "../data/1000-improved-words-shuffled.txt", "Path to a list of objective words to benchmark this algorithm against.")
}

var benchCmd = &cobra.Command{
	Use:   "bench",
	Short: "Benchmarks an algorithm against a given word list.",
	RunE: func(cmd *cobra.Command, args []string) error {
		initRoot()
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
		benchLen := benchWords.Len()
		countNumGuesses := make([]int, maxGuesses)
		results := make(chan gws.GameResult, maxThreads)
		objectives := make(chan gws.Word, maxThreads)
		errs := make(chan error, maxThreads)
		done := make(chan bool)
		benchThreads := maxThreads - 1
		for i := 0; i < benchThreads; i++ {
			go benchGuesser(objectives, results, errs, done, guesser.Copy())
		}
		go collectResults(results, done, countNumGuesses)
		for i := 0; i < benchLen; i++ {
			objective := benchWords.At(i)
			select {
			case err = <-errs:
				break
			case objectives <- objective:
				continue
			}
		}
		close(objectives)
		// Wait for bench threads.
		for dones := 0; dones < benchThreads; {
			select {
			case err = <-errs:
				// Continue to finish the other routines.
			case <-done:
				dones++
			}
		}
		close(results)
		// Wait for the collector.
		<-done
		if err != nil {
			return err
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

func benchGuesser(objectives <-chan gws.Word, results chan<- gws.GameResult, errs chan<- error, done chan<- bool, guesser gws.Guesser) {
	for objective, more := <-objectives; more; objective, more = <-objectives {
		result, err := gws.PlayGameWithGuesser(objective, maxGuesses, guesser)
		if err != nil {
			errs <- err
			done <- true
			return
		}
		if result.Status != gws.GameSuccess {
			errs <- fmt.Errorf("Failed to guess %s during benchmark.\n", objective)
			printGuesses(result.Turns)
			done <- true
			return
		}
		results <- result
	}
	done <- true
}

func collectResults(results <-chan gws.GameResult, done chan<- bool, countNumGuesses []int) {
	for result, more := <-results; more; result, more = <-results {
		numGuesses := len(result.Turns)
		countNumGuesses[numGuesses-1] = countNumGuesses[numGuesses-1] + 1
	}
	done <- true
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
