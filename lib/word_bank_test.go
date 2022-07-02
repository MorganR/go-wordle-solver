package go_wordle_solver

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestWordBankFromReaderWhenEmpty(t *testing.T) {
	_, err := WordBankFromReader(strings.NewReader(""))
	assert.Error(t, err, "At least one word must be provided.")
}

func TestWordBankFromReaderWithOnlyWhiteSpace(t *testing.T) {
	_, err := WordBankFromReader(strings.NewReader("  \n  "))
	assert.Error(t, err, "At least one word must be provided.")
}

func TestWordBankFromReaderDifferentLengthWords(t *testing.T) {
	_, err := WordBankFromReader(strings.NewReader("abc\nbcd\nefgh"))
	assert.Error(t, err, "Words must all be the same length. Encountered word with length 4 when expecting length 3.")
}

func TestWordBankFromReaderValidWords(t *testing.T) {
	bank, err := WordBankFromReader(strings.NewReader("abc\n bcd \nef£"))

	assert.NilError(t, err)
	pw := bank.Words()
	assert.Equal(t, pw.Len(), 3)
}

func TestWordBankFromSliceWithNilList(t *testing.T) {
	_, err := WordBankFromSlice(nil)
	assert.Error(t, err, "At least one word must be provided.")
}

func TestWordBankFromSliceWithEmptyList(t *testing.T) {
	_, err := WordBankFromSlice([]string{})
	assert.Error(t, err, "At least one word must be provided.")
}

func TestWordBankFromSliceWithDifferentLengthWords(t *testing.T) {
	_, err := WordBankFromSlice([]string{"bad", "rad", "good"})
	assert.Error(t, err, "Words must all be the same length. Encountered word with length 4 when expecting length 3.")
}

func TestWordBankWords(t *testing.T) {
	bank, _ := WordBankFromSlice([]string{"foo", "bar"})
	pw := bank.Words()

	assert.Equal(t, pw.Len(), 2)
}
