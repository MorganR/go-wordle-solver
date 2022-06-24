package go_wordle_solver

/// The result of a given letter at a specific location. There is some complexity here when a
/// letter appears in a word more than once. See [`GuessResult`] for more details.
type LetterResult uint8

const (
	LetterResultUnknown LetterResult = iota
	/// This letter goes exactly here in the objective word.
	LetterResultCorrect
	/// This letter is in the objective word, but not here.
	LetterResultPresentNotHere
	/// This letter is not in the objective word, or is only in the word as many times as it was
	/// marked either `PresentNotHere` or `Correct`.
	LetterResultNotPresent
)

/// The result of a single word guess.
///
/// There is some complexity here when the guess has duplicate letters. Duplicate letters are
/// matched to [`LetterResult`]s as follows:
///
/// 1. All letters in the correct location are marked `Correct`.
/// 2. For any remaining letters, if the objective word has more letters than were marked correct,
///    then these letters are marked as `PresentNotHere` starting from the beginning of the word,
///    until all letters have been accounted for.
/// 3. Any remaining letters are marked as `NotPresent`.
///
/// For example, if the guess was "sassy" for the objective word "mesas", then the results would
/// be: `[PresentNotHere, PresentNotHere, Correct, NotPresent, NotPresent]`.
type GuessResult struct {
	/// The guess that was made.
	Guess string
	/// The result of each letter, provided in the same letter order as in the guess.
	Results []LetterResult
}
