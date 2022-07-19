package go_wordle_solver

import "fmt"

// LetterResult indicates the result of a given letter at a specific location.
//
// There is some complexity here when a letter appears in a word more than once. See [GuessResult]
// for more details.
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

// String converts [LetterResult] into a readable string.
func (lr LetterResult) String() string {
	switch lr {
	case LetterResultUnknown:
		return "unknown"
	case LetterResultCorrect:
		return "correct"
	case LetterResultPresentNotHere:
		return "present not here"
	case LetterResultNotPresent:
		return "not present"
	default:
		return "invalid LetterResult"
	}
}

// A compressed form of LetterResults. Can only store vectors of up to
// [MaxLettersInCompressedGuessResult] results.
type CompressedGuessResult uint32

// How many letters can be stored in [CompressedGuessResult].
const MaxLettersInCompressedGuessResult uint8 = 32 / 3

// Creates a compressed form of the given letter results.
//
// Returns an error if letterResults has more than [MaxLettersInCompressedGuessResult] values.
func CompressResults(letterResults []LetterResult) (CompressedGuessResult, error) {
	if len(letterResults) > int(MaxLettersInCompressedGuessResult) {
		return 0, fmt.Errorf("Results can only be compressed with up to %v letters. This result has %v.", MaxLettersInCompressedGuessResult, len(letterResults))
	}
	var data CompressedGuessResult = 0
	index := 0
	for _, result := range letterResults {
		data |= 1 << (index + int(result))
		index += 3
	}
	return data, nil
}

// GuessResult is the result of a single word guess.
//
// There is some complexity here when the guess has duplicate letters. Duplicate letters are
// matched to [LetterResult]s as follows:
//
// 1. All letters in the correct location are marked [LetterResultCorrect].
// 2. For any remaining letters, if the objective word has more letters than were marked correct,
//    then these letters are marked as [LetterResultPresentNotHere] starting from the beginning of
//    the word, until all letters have been accounted for.
// 3. Any remaining letters are marked as [LetterResultNotPresent].
//
// For example, if the guess was "sassy" for the objective word "mesas", then the results would
// be: [PresentNotHere, PresentNotHere, Correct, NotPresent, NotPresent].
type GuessResult struct {
	// The guess that was made.
	Guess Word
	// The result of each letter, provided in the same letter order as in the guess.
	Results []LetterResult
}

// TurnData provides data about a single turn of a Wordle game.
type TurnData struct {
	// The guess that was made this turn.
	Guess Word
	// The number of possible words that remained at the start of this turn.
	NumPossibleWordsBeforeGuess uint
}

// GameStatus indicates whether the game was won or lost.
type GameStatus int

const (
	// Indicates that the guesser won the game.
	GameSuccess GameStatus = iota
	// Indicates that the guesser failed to guess the word under the guess limit.
	GameFailure
)

// String prints [GameStatus] as a readable string.
func (gs GameStatus) String() string {
	switch gs {
	case GameSuccess:
		return "success"
	case GameFailure:
		return "failure"
	default:
		return "invalid status"
	}
}

// GameResult is the result of a Wordle game.
type GameResult struct {
	// Whether the game was won or lost.
	Status GameStatus
	// Data for each turn that was played.
	Turns []TurnData
}

// GetResultForGuess determines the result of the given guess when applied to the given objective.
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
