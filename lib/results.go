package go_wordle_solver

import "fmt"

// The result of a given letter at a specific location. There is some complexity here when a
// letter appears in a word more than once. See [`GuessResult`] for more details.
type LetterResult uint8

const (
	// The default value.
	LetterResultUnknown LetterResult = iota
	// This letter goes exactly here in the objective word.
	LetterResultCorrect
	// This letter is in the objective word, but not here.
	LetterResultPresentNotHere
	// This letter is not in the objective word, or is only in the word as many times as it was
	// marked either `PresentNotHere` or `Correct`.
	LetterResultNotPresent
)

// The result of a single word guess.
//
// There is some complexity here when the guess has duplicate letters. Duplicate letters are
// matched to [`LetterResult`]s as follows:
//
// 1. All letters in the correct location are marked `Correct`.
// 2. For any remaining letters, if the objective word has more letters than were marked correct,
//    then these letters are marked as `PresentNotHere` starting from the beginning of the word,
//    until all letters have been accounted for.
// 3. Any remaining letters are marked as `NotPresent`.
//
// For example, if the guess was "sassy" for the objective word "mesas", then the results would
// be: `[PresentNotHere, PresentNotHere, Correct, NotPresent, NotPresent]`.
type GuessResult struct {
	// The guess that was made.
	Guess Word
	// The result of each letter, provided in the same letter order as in the guess.
	Results []LetterResult
}

// Data about a single turn of a Wordle game.
type TurnData struct {
	// The guess that was made this turn.
	Guess Word
	// The number of possible words that remained at the start of this turn.
	NumPossibleWordsBeforeGuess uint
}

// The data from a game that was played.
type GameData struct {
	// Data for each turn that was played.
	Turns []TurnData
}

// Whether the game was won or lost.
type GameStatus int

const (
	// Indicates that the guesser won the game.
	GameSuccess GameStatus = iota
	// Indicates that the guesser failed to guess the word under the guess limit.
	GameFailure
	// Indicates that the given word was not in the guesser's word bank.
	UnknownWord
)

// The result of a Wordle game.
type GameResult struct {
	// Whether the game was won or lost.
	Status GameStatus
	// Additional data about the game. Only set if `Status` is `GameSuccess` or `GameFailure`.
	Data *GameData
}

/// Determines the result of the given `guess` when applied to the given `objective`.
///
/// ```
/// result := GetResultForGuess("mesas", "sassy")
///
/// TODO: Update example.
/// assert!(
///     matches!(
///         result,
///         Ok(GuessResult {
///             guess: []rune("sassy"),
///             results: _
///         })
///     )
/// );
/// assert_eq!(
///     result.unwrap().results,
///     vec![
///         LetterResult::PresentNotHere,
///         LetterResult::PresentNotHere,
///         LetterResult::Correct,
///         LetterResult::NotPresent,
///         LetterResult::NotPresent
///     ]
/// );
/// ```
func GetResultForGuess(objective, guess Word) (GuessResult, error) {
	// Convert to runes to properly handle unicode.
	guessLen := guess.Len()
	if objective.Len() != guessLen {
		return GuessResult{}, fmt.Errorf("The guess (%s) must be the same length as the objective (length: %v).", guess, objective.Len())
	}
	// This algorithm does the following:
	// * Assume none of the letters in the guess are present in the objective.
	// * For each letter in the objective:
	//   * Check if it's correct in the guess.
	//     * If this index was previously marked as `NotPresent`, then continue.
	//     * If it's currently marked as `PresentNotHere` then see if this can be forwarded to any
	//       matches later in the guess.
	//   * To forward the match to a later letter, this starts at the index after the objective
	//     letter. It then:
	//     * Checks if the guess letter equals the objective letter. If not, go to the next letter.
	//     * If the guess letter is equal, check if it's unset (i.e. "LetterResultNotPresent"). If
	//       so, set it to LetterResultPresentNotHere. Here we're done, so we move on to the next
	//       objective letter.
	results := make([]LetterResult, guessLen)
	fillSlice(results, LetterResultNotPresent)
	for oi := 0; oi < guessLen; oi++ {
		objectiveLetter := objective.At(oi)
		startI := 0
		if objective.At(oi) == guess.At(oi) {
			existingResult := results[oi]
			results[oi] = LetterResultCorrect
			if existingResult != LetterResultPresentNotHere || oi == guessLen-1 {
				continue
			}
			startI = oi + 1
		}
		for gi := startI; gi < guessLen; gi++ {
			guessLetter := guess.At(gi)
			// Continue if this letter doesn't match.
			if guessLetter != objectiveLetter {
				continue
			}
			existingResult := results[gi]
			if existingResult != LetterResultNotPresent {
				continue
			}
			results[gi] = LetterResultPresentNotHere
			break
		}
	}
	return GuessResult{Guess: guess, Results: results}, nil
}
