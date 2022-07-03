package go_wordle_solver

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestRandomGuesserSelectNextGuess(t *testing.T) {
	bank, _ := WordBankFromSlice([]string{"abc", "bcd", "def"})
	guesser := InitRandomGuesser(&bank)

	w := guesser.SelectNextGuess()
	if !w.HasValue() {
		t.Errorf("Expected random guesser to select a valid guess.")
	}
	if !(w.Value().Equal(WordFromString("abc")) ||
		w.Value().Equal(WordFromString("bcd")) ||
		w.Value().Equal(WordFromString("def"))) {
		t.Errorf("Expected random guesser to select a guess from the word bank, but received %s.", w.Value())
	}
}

func TestRandomGuesserUpdateModifiesNextguess(t *testing.T) {
	bank, _ := WordBankFromSlice([]string{"abc", "bcd", "cde"})
	guesser := InitRandomGuesser(&bank)

	err := guesser.Update(&GuessResult{
		WordFromString("bcd"),
		[]LetterResult{
			LetterResultPresentNotHere,
			LetterResultPresentNotHere,
			LetterResultNotPresent,
		},
	})

	assert.NilError(t, err)
	got := guesser.SelectNextGuess()
	want := OptionalOf(WordFromString("abc"))
	assert.DeepEqual(t, &got, &want)
}

func TestRandomGuesserInvalidUpdateFails(t *testing.T) {
	bank, _ := WordBankFromSlice([]string{"abc", "bcd", "cde"})
	guesser := InitRandomGuesser(&bank)

	err := guesser.Update(&GuessResult{
		WordFromString("abc"),
		[]LetterResult{
			LetterResultNotPresent,
			LetterResultNotPresent,
			LetterResultPresentNotHere,
		},
	})
	assert.NilError(t, err)
	err = guesser.Update(&GuessResult{
		WordFromString("bcd"),
		[]LetterResult{
			LetterResultNotPresent,
			LetterResultPresentNotHere,
			LetterResultNotPresent,
		},
	})
	assert.NilError(t, err)
	// At this point, no words are possible, but the guess results are still valid.
	got := guesser.SelectNextGuess()
	want := Optional[Word]{}
	assert.DeepEqual(t, &got, &want)

	err = guesser.Update(&GuessResult{
		WordFromString("cde"),
		[]LetterResult{
			LetterResultPresentNotHere,
			LetterResultNotPresent,
			LetterResultNotPresent,
		},
	})
	assert.Error(t, err, "Can't set letter to not here at index 0 since it's already marked as here.")
}

func TestPlayGameWithUnknownWordRandom(t *testing.T) {
	bank, _ := WordBankFromSlice([]string{"abcz", "weyz", "defy", "ghix"})
	guesser := InitRandomGuesser(&bank)

	got := PlayGameWithGuesser(WordFromString("nope"), 10, &guesser)
	assert.Equal(t, got.Status, UnknownWord)
}

func TestPlayGameWithKnownWordRandom(t *testing.T) {
	bank, _ := WordBankFromSlice([]string{"abcz", "weyz", "defy", "ghix"})
	guesser := InitRandomGuesser(&bank)

	got := PlayGameWithGuesser(WordFromString("abcz"), 10, &guesser)

	assert.Equal(t, got.Status, GameSuccess)
	assert.Assert(t, len(got.Data.Turns) <= 4, "Random guesser took more than 4 guesses.")
	assert.DeepEqual(t, got.Data.Turns[len(got.Data.Turns)-1].Guess, WordFromString("abcz"))
}
