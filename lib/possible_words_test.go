package go_wordle_solver

import (
	"testing"
)

func TestPossibleWordsLen(t *testing.T) {
	pw := &PossibleWords{[]Word{WordFromString("foo"), WordFromString("bar")}}
	want := 2
	got := pw.Len()
	if got != want {
		t.Fatalf("Expected len %v, got %v", want, got)
	}

	pw = nil
	want = 0
	got = pw.Len()
	if got != want {
		t.Fatalf("Expected len %v, got %v", want, got)
	}
}
