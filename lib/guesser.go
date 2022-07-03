package go_wordle_solver

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Guesses words in order to solve a single Wordle.
type Guesser interface {
	// Resets this guesser for solving a new puzzle.
	Reset()

	// Updates this guesser with information about a word.
	Update(result *GuessResult) error

	// Selects a new guess for the Wordle.
	//
	// Returns an empty optional if no known words are possible given the known restrictions imposed
	// by previous calls to [`Self::update()`].
	SelectNextGuess() Optional[Word]

	// Provides read access to the remaining set of possible words in this guesser.
	PossibleWords() *PossibleWords
}

// Attempts to guess the given word within the maximum number of guesses, using the given word
// guesser.
//
// ```
// use rs_wordle_solver::GameResult;
// use rs_wordle_solver::RandomGuesser;
// use rs_wordle_solver::WordBank;
// use rs_wordle_solver::play_game_with_guesser;
//
// let bank = WordBank::from_iterator(&["abc", "def", "ghi"]).unwrap();
// let mut guesser = RandomGuesser::new(&bank);
// let result = play_game_with_guesser("def", 4, guesser.clone());
//
// assert!(matches!(result, GameResult::Success(_guesses)));
//
// let result = play_game_with_guesser("zzz", 4, guesser.clone());
//
// assert!(matches!(result, GameResult::UnknownWord));
//
// let result = play_game_with_guesser("other", 4, guesser);
//
// assert!(matches!(result, GameResult::UnknownWord));
// ```
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

// Guesses at random from the possible words that meet the restrictions.
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

// Constructs a new `RandomGuesser` using the given word bank.
//
// ```
// use rs_wordle_solver::RandomGuesser;
// use rs_wordle_solver::WordBank;
//
// let bank = WordBank::from_iterator(&["abc", "def", "ghi"]).unwrap();
// let guesser = RandomGuesser::new(&bank);
// ```
func InitRandomGuesser(bank *WordBank) RandomGuesser {
	return RandomGuesser{
		bank:          bank,
		possibleWords: bank.Words(),
		rng:           rand.New(rand.NewSource(time.Now().UnixMicro())),
	}
}

func (self *RandomGuesser) Reset() {
	self.possibleWords = self.bank.Words()
}

func (self *RandomGuesser) Update(result *GuessResult) error {
	return self.possibleWords.Filter(result)
}

func (self *RandomGuesser) SelectNextGuess() Optional[Word] {
	if self.possibleWords.Len() == 0 {
		return Optional[Word]{}
	}
	random := self.rng.Int()
	return OptionalOf(self.possibleWords.At(random % self.possibleWords.Len()))
}

func (self *RandomGuesser) PossibleWords() *PossibleWords {
	return &self.possibleWords
}

// Determines how the best guess should be chosen.
type GuessMode int

const (
	// The best guess can be chosen from all words in the word bank.
	GuessModeAll GuessMode = iota
	// The best guess can only be chosen from remaining possible words.
	GuessModePossible
)

type MaxScoreGuesser[S WordScorer] struct {
	bank           *WordBank
	possibleWords  PossibleWords
	scorer         S
	guessMode      GuessMode
	unguessedWords PossibleWords
}

func InitMaxScoreGuesser[S WordScorer](bank *WordBank, scorer S, mode GuessMode) MaxScoreGuesser[S] {
	return MaxScoreGuesser[S]{
		bank:           bank,
		possibleWords:  bank.Words(),
		scorer:         scorer,
		guessMode:      mode,
		unguessedWords: bank.Words(),
	}
}

func (self *MaxScoreGuesser[S]) Reset() {
	self.possibleWords = self.bank.Words()
	self.unguessedWords = self.bank.Words()
}

func (self *MaxScoreGuesser[S]) Update(result *GuessResult) error {
	self.unguessedWords.Remove(result.Guess)
	err := self.possibleWords.Filter(result)
	if err != nil {
		return err
	}
	return self.scorer.Update(result.Guess, &self.possibleWords)
}

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

func (self *MaxScoreGuesser[S]) PossibleWords() *PossibleWords {
	return &self.possibleWords
}
