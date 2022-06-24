package go_wordle_solver

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestNewWordBankWithNilList(t *testing.T) {
	_, err := NewWordBank(nil)
	assert.Error(t, err, "At least one word must be provided.")
}

func TestNewWordBankWithEmptyList(t *testing.T) {
	_, err := NewWordBank([]string{})
	assert.Error(t, err, "At least one word must be provided.")
}

func TestNewWordBankWithDifferentLengthWords(t *testing.T) {
	_, err := NewWordBank([]string{"bad", "rad", "good"})
	assert.Error(t, err, "Words must all be the same length.")
}

func TestWordBankWords(t *testing.T) {
	bank, _ := NewWordBank([]string{"foo", "bar"})
	pw := bank.Words()

	assert.Equal(t, pw.Len(), 2)
}
