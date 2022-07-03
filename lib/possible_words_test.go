package go_wordle_solver

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestPossibleWordsLen(t *testing.T) {
	pw := initPossibleWords([]Word{WordFromString("foo"), WordFromString("bar")})
	assert.Equal(t, pw.Len(), 2)

	var pwPointer *PossibleWords = nil
	assert.Equal(t, pwPointer.Len(), 0)
}

func TestPossibleWordsAt(t *testing.T) {
	pw := initPossibleWords([]Word{WordFromString("foo"), WordFromString("bar")})

	assert.DeepEqual(t, pw.At(0), WordFromString("foo"))
	assert.DeepEqual(t, pw.At(1), WordFromString("bar"))
}

func TestPossibleWordsFilter(t *testing.T) {
	pw := initPossibleWords([]Word{
		WordFromString("mad"),
		WordFromString("bad"),
		WordFromString("and"),
		WordFromString("cat"),
	})

	gr, _ := GetResultForGuess(WordFromString("mad"), WordFromString("add"))
	pw.Filter(&gr)

	assert.Equal(t, pw.Len(), 2)
	assert.DeepEqual(t, pw.At(0), WordFromString("mad"))
	assert.DeepEqual(t, pw.At(1), WordFromString("bad"))
}
