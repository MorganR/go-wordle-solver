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
