package go_wordle_solver

import (
	"errors"
)

// WordBank provides a read-only set of equal length words.
type WordBank struct {
	allWords   []string
	wordLength uint
}

// Creates a WordBank from the given word list.
//
// The list must be non-empty, and the words must all have the same length. Duplicate words are
// removed.
func InitWordBank(words []string) (WordBank, error) {
	if len(words) == 0 {
		return WordBank{}, errors.New("At least one word must be provided.")
	}
	wordLength := len(words[0])
	for _, word := range words {
		if len(word) != wordLength {
			return WordBank{}, errors.New("Words must all be the same length.")
		}
	}
	return WordBank{words, uint(wordLength)}, nil
}

// Returns all possible words from this word bank.
func (wb *WordBank) Words() PossibleWords {
	return PossibleWords{wb.allWords}
}
