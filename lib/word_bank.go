package go_wordle_solver

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

// WordBank provides a read-only set of equal length words.
type WordBank struct {
	allWords   [][]rune
	wordLength uint
}

const defaultWordBuffer int = 100

// Constructs a new `WordBank` struct by reading words from the given reader.
//
// The reader should provide one word per line. Each word will be trimmed and converted to
// lower case. Empty lines are skipped. At least one word must be provided.
//
// After trimming, all words must be the same length, else this returns an error.
func WordBankFromReader(r io.Reader) (WordBank, error) {
	s := bufio.NewScanner(r)
	words := make([][]rune, 0, defaultWordBuffer)
	n := 0
	wordLength := 0
	for ok := s.Scan(); ok; ok = s.Scan() {
		thisWord := []rune(strings.ToLower(strings.TrimSpace(s.Text())))
		thisWordLength := len(thisWord)
		if thisWordLength == 0 {
			continue
		}
		words = append(words, thisWord)
		if n == 0 {
			wordLength = thisWordLength
		}
		if thisWordLength != wordLength {
			return WordBank{}, fmt.Errorf("Words must all be the same length. Encountered word with length %v when expecting length %v.", thisWordLength, wordLength)
		}
		n++
	}
	if err := s.Err(); err != nil {
		return WordBank{}, err
	}
	if len(words) == 0 {
		return WordBank{}, errors.New("At least one word must be provided.")
	}
	return WordBank{words[:], uint(wordLength)}, nil
}

// Constructs a new `WordBank` struct using the words from the given vector.
//
// Each word will be trimmed and converted to lower case. At least one word must be provided.
//
// After trimming, all words must be the same length, else this returns an error.
func WordBankFromSlice(words []string) (WordBank, error) {
	if len(words) == 0 {
		return WordBank{}, errors.New("At least one word must be provided.")
	}
	wordLength := len(words[0])
	allWords := make([][]rune, len(words))
	for i, word := range words {
		wordRunes := []rune(strings.ToLower(strings.TrimSpace(word)))
		if len(wordRunes) != wordLength {
			return WordBank{}, fmt.Errorf("Words must all be the same length. Encountered word with length %v when expecting length %v.", len(wordRunes), wordLength)
		}
		allWords[i] = wordRunes
	}
	return WordBank{allWords, uint(wordLength)}, nil
}

// Returns all possible words from this word bank.
func (wb *WordBank) Words() PossibleWords {
	return PossibleWords{wb.allWords}
}
