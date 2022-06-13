package go_wordle_solver

import (
	"testing"
)

func TestNewWordBankWithNilList(t *testing.T) {
	_, err := NewWordBank(nil)
	if err == nil {
		t.Fatalf("Expected an error.")
	}
}

func TestNewWordBankWithEmptyList(t *testing.T) {
	_, err := NewWordBank([]string{})
	if err == nil {
		t.Fatalf("Expected an error.")
	}
}

func TestNewWordBankWithDifferentLengthWords(t *testing.T) {
	_, err := NewWordBank([]string{"bad", "rad", "good"})
	if err == nil {
		t.Fatalf("Expected an error.")
	}
}

func TestWordBankWords(t *testing.T) {
	bank, _ := NewWordBank([]string{"foo", "bar"})
	pw := bank.Words()

	if pw.Len() != 2 {
		t.Fatal("Expected 2 possible words.")
	}
}
