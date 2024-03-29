package go_wordle_solver

import (
	"fmt"
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func ExamplePlayGameWithGuesser() {
	bank, err := WordBankFromSlice([]string{"abc", "bcd", "cde"})
	if err != nil {
		fmt.Print(err)
		return
	}
	guesser := InitRandomGuesser(&bank)

	maxGuesses := 3
	result, err := PlayGameWithGuesser(WordFromString("bcd"), maxGuesses, &guesser)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(result.Status)

	// Output:
	// success
}

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

	_, err := PlayGameWithGuesser(WordFromString("nope"), 10, &guesser)
	assert.Error(t, err, "No more valid guesses.")
}

func TestPlayGameWithKnownWordRandom(t *testing.T) {
	bank, _ := WordBankFromSlice([]string{"abcz", "weyz", "defy", "ghix"})
	guesser := InitRandomGuesser(&bank)

	got, err := PlayGameWithGuesser(WordFromString("abcz"), 10, &guesser)

	assert.NilError(t, err)
	assert.Equal(t, got.Status, GameSuccess)
	assert.Assert(t, len(got.Turns) <= 4, "Random guesser took more than 4 guesses.")
	assert.DeepEqual(t, got.Turns[len(got.Turns)-1].Guess, WordFromString("abcz"))
}

func TestPlayGameWithUnknownWordMaxScore(t *testing.T) {
	bank, _ := WordBankFromSlice([]string{"abcz", "weyz", "defy", "ghix"})
	scorer, err := InitMaxEliminationsScorer(&bank)
	assert.NilError(t, err)
	guesser := InitMaxScoreGuesser(&bank, &scorer, GuessModeAll)

	_, err = PlayGameWithGuesser(WordFromString("nope"), 10, &guesser)
	assert.Error(t, err, "No more valid guesses.")
}

func TestPlayGameWithKnownWordMaxScore(t *testing.T) {
	bank, _ := WordBankFromSlice([]string{"abcz", "weyz", "defy", "ghix"})
	scorer, err := InitMaxEliminationsScorer(&bank)
	assert.NilError(t, err)
	guesser := InitMaxScoreGuesser(&bank, &scorer, GuessModeAll)

	got, err := PlayGameWithGuesser(WordFromString("abcz"), 10, &guesser)

	assert.NilError(t, err)
	assert.Equal(t, got.Status, GameSuccess)
	assert.Assert(t, len(got.Turns) <= 4, "Max score guesser took more than 4 guesses.")
	assert.DeepEqual(t, got.Turns[len(got.Turns)-1].Guess, WordFromString("abcz"))
}

func BenchmarkPlayGameWithRandom(b *testing.B) {
	f, err := os.Open("../data/1000-improved-words-shuffled.txt")
	if err != nil {
		b.Fatal(err)
	}
	bank, err := WordBankFromReader(f)
	if err != nil {
		b.Fatal(err)
	}
	guesser := InitRandomGuesser(&bank)
	allWords := bank.Words()
	numWords := allWords.Len()

	for i := 0; i < b.N; i++ {
		guesser.Reset()
		word := allWords.At(i % numWords)
		result, err := PlayGameWithGuesser(word, 128, &guesser)
		if err != nil {
			b.Fatal(err)
		}
		if result.Status != GameSuccess {
			b.Fatalf("Game failed for word %s, result: %v", word, result)
		}
	}
}

func BenchmarkPlayGameWithMaxScore(b *testing.B) {
	f, err := os.Open("../data/1000-improved-words-shuffled.txt")
	if err != nil {
		b.Fatal(err)
	}
	bank, err := WordBankFromReader(f)
	if err != nil {
		b.Fatal(err)
	}
	scorer, err := InitMaxEliminationsScorer(&bank)
	if err != nil {
		b.Fatal(err)
	}
	guesser := InitMaxScoreGuesser(&bank, &scorer, GuessModeAll)
	allWords := bank.Words()
	numWords := allWords.Len()

	for i := 0; i < b.N; i++ {
		guesser.Reset()
		word := allWords.At(i % numWords)
		result, err := PlayGameWithGuesser(word, 128, &guesser)
		if err != nil {
			b.Fatal(err)
		}
		if result.Status != GameSuccess {
			b.Fatalf("Game failed for word %s, result: %v", word, result)
		}
	}
}
