package go_wordle_solver

import (
	"fmt"
	"runtime"
)

var maxThreads = runtime.NumCPU()

// Gives words a score, where the maximum score indicates the best guess.
type WordScorer interface {
	// Provides a copy of this scorer.
	Copy() WordScorer
	// Resets this scorer for a new puzzle.
	Reset(pw *PossibleWords)
	// Updates the scorer with the latest guess, the updated set of restrictions, and the updated
	// list of possible words.
	Update(latestGuess Word, pw *PossibleWords) error
	// Determines a score for the given word. The higher the score, the better the guess.
	ScoreWord(word Word) int64
}

// This probabilistically calculates the expectation value for how many words will be eliminated by
// each guess, and chooses the word that eliminates the most other guesses.
//
// This is a highly effective scoring strategy, but also quite expensive to compute. On my
// machine, constructing the scorer for about 4600 words takes about 1.4 seconds, but each
// subsequent game can be played in about 27ms if the scorer is then cloned before each game.
//
// When benchmarked against the 4602 words in `data/improved-words.txt`, this has the following
// results:
//
// |Num guesses|Num games (Guess from: `PossibleWords`)|Num games (Guess from: `AllUnguessedWords`)|
// |-----------|---------|---------------|
// |1|1|1|
// |2|180|53|
// |3|1452|1426|
// |4|1942|2635|
// |5|666|468|
// |6|220|19|
// |7|93|0|
// |8|33|0|
// |9|10|0|
// |10|4|0|
// |11|1|0|
//
// **Average guesses:**
//
// GuessModePossible: 3.95 +/- 1.10
//
// GuessModeAll: 3.78 +/- 0.65
type MaxEliminationsScorer struct {
	possibleWords                    *PossibleWords
	firstExpectedEliminationsPerWord map[string]float64
	isFirstRound                     bool
}

const maxEliminationsChunkSize int = 128

// Constructs a `MaxEliminationsScorer`. **Be careful, this is expensive to compute!**
//
// Once constructed for a given set of words, this precomputation can be reused by simply
// cloning a new version of the scorer for each game.
//
// The cost of this function scales in approximately *O*(*n*<sup>2</sup>), where *n* is the
// number of words.
//
// ```
// use rs_wordle_solver::GuessFrom;
// use rs_wordle_solver::Guesser;
// use rs_wordle_solver::MaxScoreGuesser;
// use rs_wordle_solver::WordBank;
// use rs_wordle_solver::scorers::MaxEliminationsScorer;
//
// let bank = WordBank::from_iterator(&["abc", "def", "ghi"]).unwrap();
// let scorer = MaxEliminationsScorer::new(&bank).unwrap();
// let mut guesser = MaxScoreGuesser::new(GuessFrom::AllUnguessedWords, &bank, scorer);
//
// assert!(guesser.select_next_guess().is_some());
// ```
func InitMaxEliminationsScorer(bank *WordBank) (MaxEliminationsScorer, error) {
	words := bank.Words()
	numWords := words.Len()
	orderedExpectedEliminations := make([]float64, numWords)

	chunks := make(chan int)
	errs := make(chan error)
	done := make(chan bool)

	for i := 0; i < maxThreads; i++ {
		go computeExpectedEliminationsChunk(chunks, errs, done, &words, orderedExpectedEliminations)
	}
	var err error = nil
	for start := 0; start < numWords; start += maxEliminationsChunkSize {
		select {
		case err = <-errs:
			break
		case chunks <- start:
			continue
		}
	}
	close(chunks)
	for dones := 0; dones < maxThreads; {
		select {
		case err = <-errs:
			// Continue to clear out dones
		case <-done:
			dones++
		}
	}
	if err != nil {
		return MaxEliminationsScorer{}, err
	}

	expectedEliminationsPerWord := make(map[string]float64, numWords)
	for i := 0; i < numWords; i++ {
		word := words.At(i)
		expectedEliminationsPerWord[word.String()] = orderedExpectedEliminations[i]
	}
	return MaxEliminationsScorer{
		possibleWords:                    &words,
		firstExpectedEliminationsPerWord: expectedEliminationsPerWord,
		isFirstRound:                     true,
	}, nil
}

func (self *MaxEliminationsScorer) Copy() WordScorer {
	pwCopy := self.possibleWords.Copy()
	return &MaxEliminationsScorer{
		&pwCopy,
		self.firstExpectedEliminationsPerWord,
		self.isFirstRound,
	}
}

func (self *MaxEliminationsScorer) Reset(pw *PossibleWords) {
	self.possibleWords = pw
	self.isFirstRound = true
}

func (self *MaxEliminationsScorer) Update(latestGuess Word, pw *PossibleWords) error {
	self.possibleWords = pw
	self.isFirstRound = false
	return nil
}

func (self *MaxEliminationsScorer) ScoreWord(w Word) int64 {
	if self.isFirstRound {
		if expectedEliminations, isPresent := self.firstExpectedEliminationsPerWord[w.String()]; isPresent {
			return int64(expectedEliminations * 1000.0)
		}
	}

	expectedEliminations, err := computeExpectedEliminations(w, self.possibleWords)
	if err != nil {
		panic(fmt.Sprintf("Failed to compute expectations for word: %s, error: %s", w, err))
	}
	return int64(expectedEliminations * 1000.0)
}

func computeExpectedEliminationsChunk(startIndices <-chan int, errs chan<- error, done chan<- bool, pw *PossibleWords, results []float64) {
	pwLen := pw.Len()
	for startIndex, ok := <-startIndices; ok; startIndex, ok = <-startIndices {
		endIndex := startIndex + maxEliminationsChunkSize
		if endIndex > pwLen {
			endIndex = pwLen
		}
		for i := startIndex; i < endIndex; i++ {
			word := pw.At(i)
			expectation, err := computeExpectedEliminations(word, pw)
			if err != nil {
				errs <- err
				done <- true
				return
			}
			results[i] = expectation
		}
	}
	done <- true
}

func computeExpectedEliminations(guess Word, possibleWords *PossibleWords) (float64, error) {
	numPossible := possibleWords.Len()
	matchingResults := make(map[CompressedGuessResult]uint, numPossible)
	for i := 0; i < numPossible; i++ {
		objective := possibleWords.At(i)
		result, err := GetResultForGuess(objective, guess)
		if err != nil {
			return 0.0, err
		}
		compressed, err := CompressResults(result.Results)
		if err != nil {
			return 0.0, err
		}
		var count uint = 1
		if knownCount, isPresent := matchingResults[compressed]; isPresent {
			count += knownCount
		}
		matchingResults[compressed] = count
	}
	numerator := uint(0)
	for _, numMatched := range matchingResults {
		numEliminated := uint(numPossible) - numMatched
		numerator += numEliminated * numMatched
	}
	return float64(numerator) / float64(numPossible), nil
}
