package go_wordle_solver

// Gives words a score, where the maximum score indicates the best guess.
type WordScorer interface {
	// Updates the scorer with the latest guess, the updated set of restrictions, and the updated
	// list of possible words.
	Update(latestGuess Word, possibleWords *PossibleWords) error
	// Determines a score for the given word. The higher the score, the better the guess.
	ScoreWord(word Word) int64
}
