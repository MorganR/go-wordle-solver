package go_wordle_solver

import "golang.org/x/exp/slices"

// PossibleWords provides a list of possible words for a Wordle puzzle.
//
// It provides easy operations to access words and to filter the list based on a [GuessResult].
// PossibleWords can be retrieved from a [WordBank].
type PossibleWords struct {
	words        []Word
	restrictions WordRestrictions
}

func initPossibleWords(words []Word) PossibleWords {
	return PossibleWords{
		slices.Clone(words),
		InitWordRestrictions(uint8(words[0].Len())),
	}
}

// Copies this set of possible words.
func (pw *PossibleWords) Copy() PossibleWords {
	return PossibleWords{
		slices.Clone(pw.words),
		pw.restrictions,
	}
}

// Len returns the number of possible words.
func (pw *PossibleWords) Len() int {
	if pw == nil {
		return 0
	}
	return len(pw.words)
}

// At retrieves the word at the given index.
func (pw *PossibleWords) At(i int) Word {
	return pw.words[i]
}

// Filter filters the possible words based on the given [GuessResult].
//
// Results from multiple calls to this method are accumulated to filter as many words as possible.
// If results conflict, an error is returned.
func (pw *PossibleWords) Filter(gr *GuessResult) error {
	err := pw.restrictions.Update(gr)
	if err != nil {
		return err
	}
	pw.words = filter(pw.words, pw.restrictions.IsSatisfiedBy)
	return nil
}

// Remove deletes the given word, if present.
//
// Returns true if the word was previously present and has now been removed.
func (pw *PossibleWords) Remove(w Word) bool {
	i := slices.IndexFunc(pw.words, w.Equal)
	if i >= 0 {
		pw.words = slices.Delete(pw.words, i, i+1)
		return true
	}
	return false
}

// Maximizing returns the word that maximizes the given function.
//
// This panics if there are no words in this [PossibleWords] object.
func (pw *PossibleWords) Maximizing(fn func(w Word) int64) Word {
	bestWord := pw.words[0]
	bestScore := fn(bestWord)
	length := len(pw.words)
	for i := 1; i < length; i++ {
		word := pw.words[i]
		score := fn(word)
		if bestScore != score {
			if bestScore < score {
				bestScore = score
				bestWord = word
			}
		}
	}
	return bestWord
}
