package go_wordle_solver

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestLen(t *testing.T) {
	w := WordFromString("hello")
	assert.Equal(t, w.Len(), 5)
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
