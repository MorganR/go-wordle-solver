package go_wordle_solver

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestMaxEliminationsScoreWord(t *testing.T) {
	bank, err := WordBankFromSlice([]string{"cod", "wod", "mod"})
	assert.NilError(t, err)

	scorer, err := InitMaxEliminationsScorer(&bank)
	assert.NilError(t, err)
	assert.Equal(t, scorer.ScoreWord(WordFromString("cod")), int64(1333))
	assert.Equal(t, scorer.ScoreWord(WordFromString("mwc")), int64(2000))
	assert.Equal(t, scorer.ScoreWord(WordFromString("zzz")), int64(0))
}

func TestMaxEliminationsScoreWordWithUpdateAndReset(t *testing.T) {
	bank, err := WordBankFromSlice([]string{
		"abb", "abc", "bad", "zza", "zzz",
	})
	assert.NilError(t, err)
	scorer, err := InitMaxEliminationsScorer(&bank)
	assert.NilError(t, err)

	preUpdatePossibleWords := bank.Words()
	preUpdateScores := make([]int64, preUpdatePossibleWords.Len())
	for i := 0; i < preUpdatePossibleWords.Len(); i++ {
		preUpdateScores[i] = scorer.ScoreWord(preUpdatePossibleWords.At(i))
	}

	// Update
	result := GuessResult{
		Guess: WordFromString("zza"),
		Results: []LetterResult{
			LetterResultNotPresent,
			LetterResultNotPresent,
			LetterResultPresentNotHere,
		},
	}
	pw := bank.Words()
	err = pw.Filter(&result)
	assert.NilError(t, err)

	err = scorer.Update(result.Guess, &pw)
	assert.NilError(t, err)
	// Still possible: abb, abc, bad
	// Eliminates 2 in all cases.
	assert.Equal(t, scorer.ScoreWord(WordFromString("abb")), int64(2000))
	// Eliminates 2 in all cases.
	assert.Equal(t, scorer.ScoreWord(WordFromString("abc")), int64(2000))
	// Could be true in one case (elimnate 2), or false in 2 cases (eliminate 1)
	assert.Equal(t, scorer.ScoreWord(WordFromString("bad")), int64(1333))
	assert.Equal(t, scorer.ScoreWord(WordFromString("zzz")), int64(0))

	// Reset
	scorer.Reset(&preUpdatePossibleWords)
	for i := 0; i < preUpdatePossibleWords.Len(); i++ {
		assert.Equal(t, scorer.ScoreWord(preUpdatePossibleWords.At(i)), preUpdateScores[i])
	}
}

func BenchmarkInitMaxEliminationsScorer(b *testing.B) {
	f, err := os.Open("../data/1000-improved-words-shuffled.txt")
	if err != nil {
		b.Fatal(err)
	}
	bank, err := WordBankFromReader(f)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		_, err := InitMaxEliminationsScorer(&bank)
		if err != nil {
			b.Fatal(err)
		}
	}
}
