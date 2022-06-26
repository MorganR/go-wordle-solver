package go_wordle_solver

import (
	"errors"
)

// WordBank provides a read-only set of equal length words.
type WordBank struct {
	allWords   [][]rune
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
	allWords := make([][]rune, len(words))
	for i, word := range words {
		wordRunes := []rune(word)
		if len(wordRunes) != wordLength {
			return WordBank{}, errors.New("Words must all be the same length.")
		}
		allWords[i] = wordRunes
	}
	return WordBank{allWords, uint(wordLength)}, nil
}

// Returns all possible words from this word bank.
func (wb *WordBank) Words() PossibleWords {
	return PossibleWords{wb.allWords}
}
