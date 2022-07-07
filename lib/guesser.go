package go_wordle_solver

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// A Guesser guesses words in order to solve a single Wordle.
type Guesser interface {
	// Copies this guesser.
	Copy() Guesser

	// Resets this guesser for solving a new puzzle.
	Reset()

	// Updates this guesser with information about a guess.
	Update(result *GuessResult) error

	// Selects a new guess for the Wordle.
	//
	// Returns an empty optional if no known words are possible given the known restrictions imposed
	// by previous calls to [Guesser.Update].
	SelectNextGuess() Optional[Word]

	// Provides read access to the remaining set of possible words in this guesser.
	PossibleWords() *PossibleWords
}

// Attempts to guess the given word within the maximum number of guesses, using the given [Guesser].
func PlayGameWithGuesser[G Guesser](
	objective Word,
	maxNumGuesses int,
	guesser G,
) (GameResult, error) {
	guesser.Reset()
	turns := make([]TurnData, 0, maxNumGuesses)
	for i := 0; i < maxNumGuesses; i++ {
		maybeGuess := guesser.SelectNextGuess()
		if !maybeGuess.HasValue() {
			return GameResult{}, errors.New("No more valid guesses.")
		}
		guess := maybeGuess.Value()
		numPossibleWordsBeforeGuess := guesser.PossibleWords().Len()
		result, err := GetResultForGuess(objective, guess)
		if err != nil {
			return GameResult{}, fmt.Errorf("Couldn't get result for guess %s, error: %s", guess, err)
		}
		turns = append(turns, TurnData{
			guess,
			uint(numPossibleWordsBeforeGuess),
		})
		if allValues(result.Results, func(lr LetterResult) bool {
			return lr == LetterResultCorrect
		}) {
			return GameResult{GameSuccess, turns}, nil
		}
		err = guesser.Update(&result)
		if err != nil {
			panic(fmt.Sprintf("Failed to update the guesser. Error: %s", err))
		}
	}
	return GameResult{GameFailure, turns}, nil
}

// Guesses at random from the possible words that meet the restrictions imposed by each guess.
//
// A sample benchmark against the `data/improved-words.txt` list performed as follows:
//
// |Num guesses to win|Num games|
// |------------------|---------|
// |1|1|
// |2|106|
// |3|816|
// |4|1628|
// |5|1248|
// |6|518|
// |7|180|
// |8|67|
// |9|28|
// |10|7|
// |11|2|
// |12|1|
//
// **Average number of guesses:** 4.49 +/- 1.26
type RandomGuesser struct {
	bank          *WordBank
	possibleWords PossibleWords
	rng           *rand.Rand
}

// Constructs a new [RandomGuesser] using the given word bank.
func InitRandomGuesser(bank *WordBank) RandomGuesser {
	return RandomGuesser{
		bank:          bank,
		possibleWords: bank.Words(),
		rng:           rand.New(rand.NewSource(time.Now().UnixMicro())),
	}
}

// Copies the [RandomGuesser].
//
// The current state of possible words is maintained, but the random source may or may not change.
func (self *RandomGuesser) Copy() Guesser {
	return &RandomGuesser{
		self.bank,
		self.possibleWords.Copy(),
		rand.New(rand.NewSource(time.Now().UnixMicro())),
	}
}

// Resets the [RandomGuesser]'s possible words.
func (self *RandomGuesser) Reset() {
	self.possibleWords = self.bank.Words()
}

// Updates this guesser's possible words based on the result.
func (self *RandomGuesser) Update(result *GuessResult) error {
	return self.possibleWords.Filter(result)
}

// Selects a new guess at random from the remaining possible words.
func (self *RandomGuesser) SelectNextGuess() Optional[Word] {
	if self.possibleWords.Len() == 0 {
		return Optional[Word]{}
	}
	random := self.rng.Int()
	return OptionalOf(self.possibleWords.At(random % self.possibleWords.Len()))
}

// Returns a pointer to the possible words for this guesser.
//
// This remains valid until [RandomGuesser.Reset] is called.
func (self *RandomGuesser) PossibleWords() *PossibleWords {
	return &self.possibleWords
}

// GuessMode determines how the best guess should be chosen.
type GuessMode int

const (
	// The best guess can be chosen from all words in the word bank.
	GuessModeAll GuessMode = iota
	// The best guess can only be chosen from remaining possible words.
	GuessModePossible
)

// String converts [GuessMode] to a readable string.
func (gm GuessMode) String() string {
	switch gm {
	case GuessModeAll:
		return "all"
	case GuessModePossible:
		return "possible"
	default:
		return "invalid GuessMode"
	}
}

// MaxScoreGuesser guesses the wordle answer by selecting the word that maximizes a score, as
// scored by the [WordScorer] implementation.
//
// This can support a large variety of algorithms.
type MaxScoreGuesser[S WordScorer] struct {
	bank           *WordBank
	possibleWords  PossibleWords
	scorer         S
	guessMode      GuessMode
	unguessedWords PossibleWords
}

// InitMaxScoreGuesser constructs a [MaxScoreGuesser] for the given bank, scorer and mode.
func InitMaxScoreGuesser[S WordScorer](bank *WordBank, scorer S, mode GuessMode) MaxScoreGuesser[S] {
	return MaxScoreGuesser[S]{
		bank:           bank,
		possibleWords:  bank.Words(),
		scorer:         scorer,
		guessMode:      mode,
		unguessedWords: bank.Words(),
	}
}

// Copy copies the [MaxScoreGuesser].
func (self *MaxScoreGuesser[S]) Copy() Guesser {
	return &MaxScoreGuesser[S]{
		self.bank,
		self.possibleWords.Copy(),
		self.scorer.Copy().(S),
		self.guessMode,
		self.unguessedWords.Copy(),
	}
}

// Reset resets the [MaxScoreGuesser]'s possible words so it can be used to solve a new Wordle.
func (self *MaxScoreGuesser[S]) Reset() {
	self.possibleWords = self.bank.Words()
	self.unguessedWords = self.bank.Words()
	self.scorer.Reset(&self.possibleWords)
}

// Update updates the current possible words based on the given result.
func (self *MaxScoreGuesser[S]) Update(result *GuessResult) error {
	self.unguessedWords.Remove(result.Guess)
	err := self.possibleWords.Filter(result)
	if err != nil {
		return err
	}
	return self.scorer.Update(result.Guess, &self.possibleWords)
}

// SelectNextGuess returns the guess that maximizes the owned [WordScorer]'s score.
//
// If there are no more possible words, this returns an empty optional. This should only happen if
// the objective word is not in this guesser's [WordBank].
func (self *MaxScoreGuesser[S]) SelectNextGuess() Optional[Word] {
	if self.possibleWords.Len() == 0 {
		return Optional[Word]{}
	}

	if self.guessMode == GuessModeAll && self.possibleWords.Len() > 2 {
		bestWord := self.unguessedWords.At(0)
		bestScore := self.scorer.ScoreWord(bestWord)
		scoresAllSame := true
		length := self.unguessedWords.Len()
		for i := 1; i < length; i++ {
			word := self.unguessedWords.At(i)
			score := self.scorer.ScoreWord(word)
			if bestScore != score {
				scoresAllSame = false
				if bestScore < score {
					bestScore = score
					bestWord = word
				}
			}
		}
		// If the scores are all the same, be sure to use a possible word so there is a chance of
		// getting it right.
		if scoresAllSame {
			return OptionalOf(self.possibleWords.At(0))
		}
		return OptionalOf(bestWord)
	}

	return OptionalOf(self.possibleWords.Maximizing(self.scorer.ScoreWord))
}

// PossibleWords provides a pointer to the possible words for this guesser.
//
// This remains valid until [MaxScoreGuesser.Reset] is called.
func (self *MaxScoreGuesser[S]) PossibleWords() *PossibleWords {
	return &self.possibleWords
}
