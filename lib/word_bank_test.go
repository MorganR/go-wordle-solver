package go_wordle_solver

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestInitWordBankWithNilList(t *testing.T) {
	_, err := InitWordBank(nil)
	assert.Error(t, err, "At least one word must be provided.")
}

func TestInitWordBankWithEmptyList(t *testing.T) {
	_, err := InitWordBank([]string{})
	assert.Error(t, err, "At least one word must be provided.")
}

func TestInitWordBankWithDifferentLengthWords(t *testing.T) {
	_, err := InitWordBank([]string{"bad", "rad", "good"})
	assert.Error(t, err, "Words must all be the same length.")
}

func TestWordBankWords(t *testing.T) {
	bank, _ := InitWordBank([]string{"foo", "bar"})
	pw := bank.Words()

	assert.Equal(t, pw.Len(), 2)
}
