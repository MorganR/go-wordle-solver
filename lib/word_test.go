package go_wordle_solver

import (
	"testing"

	"golang.org/x/exp/slices"
	"gotest.tools/v3/assert"
)

func TestLen(t *testing.T) {
	w := WordFromString("hello")
	assert.Equal(t, w.Len(), uint8(5))
}

func TestAt(t *testing.T) {
	w := WordFromString("hello")
	assert.Equal(t, w.At(0), 'h')
	assert.Equal(t, w.At(1), 'e')
	assert.Equal(t, w.At(2), 'l')
	assert.Equal(t, w.At(3), 'l')
	assert.Equal(t, w.At(4), 'o')

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic when indexing out of bounds.")
		}
	}()
	w.At(5)
}

func TestString(t *testing.T) {
	w := WordFromString("hello")
	assert.Equal(t, w.String(), "hello")
}

func TestFullForEach(t *testing.T) {
	w := WordFromString("hello")

	wordLength := w.Len()
	letters := make([]rune, wordLength)

	w.ForEach(func(i int, l rune) bool {
		letters[i] = l
		return true
	})

	want := []rune{'h', 'e', 'l', 'l', 'o'}
	if slices.Compare(letters, want) != 0 {
		t.Errorf("Expected letters to have value %v but found %v", want, letters)
	}
}

func TestForEachWithBreak(t *testing.T) {
	w := WordFromString("hello")

	wordLength := w.Len()
	letters := make([]rune, wordLength)

	w.ForEach(func(i int, l rune) bool {
		letters[i] = l
		return i != 2
	})

	want := []rune{'h', 'e', 'l', 0, 0}
	if slices.Compare(letters, want) != 0 {
		t.Errorf("Expected letters to have value %v but found %v", want, letters)
	}
}
