package go_wordle_solver

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

func ExampleGetResultForGuess() {
	result, err := GetResultForGuess(WordFromString("mesas"), WordFromString("sassy"))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)

	// Output:
	// {sassy [present not here present not here correct not present not present]}
}

func TestGetResultForGuessCorrect(t *testing.T) {
	result, err := GetResultForGuess(WordFromString("abcb"), WordFromString("abcb"))

	assert.NilError(t, err)
	assert.DeepEqual(t,
		result,
		GuessResult{
			Guess: WordFromString("abcb"),
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
	result, err := GetResultForGuess(WordFromString("abc"), WordFromString("ab£"))

	assert.NilError(t, err)
	assert.DeepEqual(t,
		result,
		GuessResult{
			Guess: WordFromString("ab£"),
			Results: []LetterResult{
				LetterResultCorrect,
				LetterResultCorrect,
				LetterResultNotPresent,
			},
		},
	)
}

func TestGetResultForGuessPartial(t *testing.T) {
	result, err := GetResultForGuess(WordFromString("mesas"), WordFromString("sassy"))

	assert.NilError(t, err)
	assert.DeepEqual(t,
		result,
		GuessResult{
			Guess: WordFromString("sassy"),
			Results: []LetterResult{
				LetterResultPresentNotHere,
				LetterResultPresentNotHere,
				LetterResultCorrect,
				LetterResultNotPresent,
				LetterResultNotPresent,
			},
		})

	result, err = GetResultForGuess(WordFromString("abba"), WordFromString("babb"))
	assert.NilError(t, err)
	assert.DeepEqual(t,
		result,
		GuessResult{
			Guess: WordFromString("babb"),
			Results: []LetterResult{
				LetterResultPresentNotHere,
				LetterResultPresentNotHere,
				LetterResultCorrect,
				LetterResultNotPresent,
			},
		})

	result, err = GetResultForGuess(WordFromString("abcb"), WordFromString("bcce"))
	assert.NilError(t, err)
	assert.DeepEqual(t,
		result,
		GuessResult{
			Guess: WordFromString("bcce"),
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
	result, err := GetResultForGuess(WordFromString("abcb"), WordFromString("defg"))
	assert.NilError(t, err)
	assert.DeepEqual(t,
		result,
		GuessResult{
			Guess: WordFromString("defg"),
			Results: []LetterResult{
				LetterResultNotPresent,
				LetterResultNotPresent,
				LetterResultNotPresent,
				LetterResultNotPresent,
			},
		})
}

func TestGetResultForGuessInvalidGuess(t *testing.T) {
	_, err := GetResultForGuess(WordFromString("goal"), WordFromString("guess"))

	assert.Error(t, err, "The guess (guess) must be the same length as the objective (length: 4).")
}

func BenchmarkGetResultForGuessCorrect(b *testing.B) {
	objective := WordFromString("abcbd")
	for n := 0; n < b.N; n++ {
		GetResultForGuess(objective, objective)
	}
}

func BenchmarkGetResultForGuessPartial(b *testing.B) {
	objective := WordFromString("mesas")
	guess := WordFromString("sassy")
	for n := 0; n < b.N; n++ {
		GetResultForGuess(objective, guess)
	}
}

func BenchmarkGetResultForGuessNoMatch(b *testing.B) {
	objective := WordFromString("abcdefg")
	guess := WordFromString("hijklmn")
	for n := 0; n < b.N; n++ {
		GetResultForGuess(objective, guess)
	}
}

func TestCompressResultsEquality(t *testing.T) {
	correct, err := CompressResults([]LetterResult{
		LetterResultCorrect,
		LetterResultCorrect,
		LetterResultCorrect,
	})
	assert.NilError(t, err)
	notHere, err := CompressResults([]LetterResult{
		LetterResultPresentNotHere,
		LetterResultPresentNotHere,
		LetterResultPresentNotHere,
	})
	assert.NilError(t, err)
	notPresent, err := CompressResults([]LetterResult{
		LetterResultNotPresent,
		LetterResultNotPresent,
		LetterResultNotPresent,
	})
	assert.NilError(t, err)

	assert.Equal(t, correct, correct)
	assert.Equal(t, notHere, notHere)
	assert.Equal(t, notPresent, notPresent)
	assert.Assert(t, correct != notHere)
	assert.Assert(t, correct != notPresent)
	assert.Assert(t, notHere != notPresent)
}

func TestCompressResultsLimits(t *testing.T) {
	_, err := CompressResults(make([]LetterResult, MaxLettersInCompressedGuessResult))
	assert.NilError(t, err)

	_, err = CompressResults(make([]LetterResult, MaxLettersInCompressedGuessResult+1))
	assert.Error(t, err, "Results can only be compressed with up to 10 letters. This result has 11.")
}
