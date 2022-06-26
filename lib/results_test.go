package go_wordle_solver

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestGetResultForGuessCorrect(t *testing.T) {
	result, err := GetResultForGuess([]rune("abcb"), []rune("abcb"))

	assert.NilError(t, err)
	assert.DeepEqual(t,
		result,
		GuessResult{
			Guess: []rune("abcb"),
			Results: []LetterResult{
				LetterResultCorrect,
				LetterResultCorrect,
				LetterResultCorrect,
				LetterResultCorrect,
			},
		},
	)
}

func TestGetResultForGuessSupportsUnicode(t *testing.T) {
	result, err := GetResultForGuess([]rune("abc"), []rune("ab£"))

	assert.NilError(t, err)
	assert.DeepEqual(t,
		result,
		GuessResult{
			Guess: []rune("ab£"),
			Results: []LetterResult{
				LetterResultCorrect,
				LetterResultCorrect,
				LetterResultNotPresent,
			},
		},
	)
}

func TestGetResultForGuessPartial(t *testing.T) {
	result, err := GetResultForGuess([]rune("mesas"), []rune("sassy"))

	assert.NilError(t, err)
	assert.DeepEqual(t,
		result,
		GuessResult{
			Guess: []rune("sassy"),
			Results: []LetterResult{
				LetterResultPresentNotHere,
				LetterResultPresentNotHere,
				LetterResultCorrect,
				LetterResultNotPresent,
				LetterResultNotPresent,
			},
		})

	result, err = GetResultForGuess([]rune("abba"), []rune("babb"))
	assert.NilError(t, err)
	assert.DeepEqual(t,
		result,
		GuessResult{
			Guess: []rune("babb"),
			Results: []LetterResult{
				LetterResultPresentNotHere,
				LetterResultPresentNotHere,
				LetterResultCorrect,
				LetterResultNotPresent,
			},
		})

	result, err = GetResultForGuess([]rune("abcb"), []rune("bcce"))
	assert.NilError(t, err)
	assert.DeepEqual(t,
		result,
		GuessResult{
			Guess: []rune("bcce"),
			Results: []LetterResult{
				LetterResultPresentNotHere,
				LetterResultNotPresent,
				LetterResultCorrect,
				LetterResultNotPresent,
			},
		},
	)
}

func TestGetResultForGuessNoneMatch(t *testing.T) {
	result, err := GetResultForGuess([]rune("abcb"), []rune("defg"))
	assert.NilError(t, err)
	assert.DeepEqual(t,
		result,
		GuessResult{
			Guess: []rune("defg"),
			Results: []LetterResult{
				LetterResultNotPresent,
				LetterResultNotPresent,
				LetterResultNotPresent,
				LetterResultNotPresent,
			},
		})
}

func TestGetResultForGuessInvalidGuess(t *testing.T) {
	_, err := GetResultForGuess([]rune("goal"), []rune("guess"))

	assert.Error(t, err, "The guess (guess) must be the same length as the objective (length: 4).")
}

func BenchmarkGetResultForGuessCorrect(b *testing.B) {
	objective := []rune("abcbd")
	for n := 0; n < b.N; n++ {
		GetResultForGuess(objective, objective)
	}
}

func BenchmarkGetResultForGuessPartial(b *testing.B) {
	objective := []rune("mesas")
	guess := []rune("sassy")
	for n := 0; n < b.N; n++ {
		GetResultForGuess(objective, guess)
	}
}

func BenchmarkGetResultForGuessNoMatch(b *testing.B) {
	objective := []rune("abcdefg")
	guess := []rune("hijklmn")
	for n := 0; n < b.N; n++ {
		GetResultForGuess(objective, guess)
	}
}
