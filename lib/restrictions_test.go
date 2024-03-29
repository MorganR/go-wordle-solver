package go_wordle_solver

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestNewPresentLetter(t *testing.T) {
	presence := newPresentLetter(3)

	assert.Equal(t, presence.state(0), llsUnknown)
	assert.Equal(t, presence.state(1), llsUnknown)
	assert.Equal(t, presence.state(2), llsUnknown)
}

func TestPresentLetterSetHere(t *testing.T) {
	letter := newPresentLetter(3)

	assert.NilError(t, letter.setMustBeAt(1))

	assert.Equal(t, letter.state(0), llsUnknown)
	assert.Equal(t, letter.state(1), llsHere)
	assert.Equal(t, letter.state(2), llsUnknown)
}

func TestPresentLetterSetHereCanBeRepeated(t *testing.T) {
	letter := newPresentLetter(3)

	assert.NilError(t, letter.setMustBeAt(1))
	assert.NilError(t, letter.setMustBeAt(1))
	assert.NilError(t, letter.setMustBeAt(1))
	assert.NilError(t, letter.setMustBeAt(1))

	assert.Equal(t, letter.state(0), llsUnknown)
	assert.Equal(t, letter.state(1), llsHere)
	assert.Equal(t, letter.state(2), llsUnknown)
}

func TestPresentLetterSetNotHere(t *testing.T) {
	letter := newPresentLetter(3)

	assert.NilError(t, letter.setMustNotBeAt(1))

	assert.Equal(t, letter.state(0), llsUnknown)
	assert.Equal(t, letter.state(1), llsNotHere)
	assert.Equal(t, letter.state(2), llsUnknown)
}

func TestPresentLetterSetNotHereCanBeRepeated(t *testing.T) {
	letter := newPresentLetter(3)

	assert.NilError(t, letter.setMustNotBeAt(1))
	assert.NilError(t, letter.setMustNotBeAt(1))
	assert.NilError(t, letter.setMustNotBeAt(1))
	assert.NilError(t, letter.setMustNotBeAt(1))

	assert.Equal(t, letter.state(0), llsUnknown)
	assert.Equal(t, letter.state(1), llsNotHere)
	assert.Equal(t, letter.state(2), llsUnknown)
}

func TestPresentLetterInferMustBeHere(t *testing.T) {
	letter := newPresentLetter(3)

	assert.NilError(t, letter.setMustNotBeAt(1))
	assert.NilError(t, letter.setMustNotBeAt(2))

	assert.Equal(t, letter.state(0), llsHere)
	assert.Equal(t, letter.state(1), llsNotHere)
	assert.Equal(t, letter.state(2), llsNotHere)
}

func TestPresentLetterMustBeHereWholeWord(t *testing.T) {
	letter := newPresentLetter(3)

	assert.NilError(t, letter.setMustBeAt(0))
	assert.NilError(t, letter.setMustBeAt(1))
	assert.NilError(t, letter.setMustBeAt(2))

	assert.Equal(t, letter.state(0), llsHere)
	assert.Equal(t, letter.state(1), llsHere)
	assert.Equal(t, letter.state(2), llsHere)
}

func TestPresentLetterMaxCountThenHereFillsRemainderNotHere(t *testing.T) {
	letter := newPresentLetter(3)

	assert.NilError(t, letter.setRequiredCount(2))
	assert.NilError(t, letter.setMustBeAt(1))

	assert.Equal(t, letter.state(0), llsUnknown)
	assert.Equal(t, letter.state(1), llsHere)
	assert.Equal(t, letter.state(2), llsUnknown)

	// Same location, no change.
	assert.NilError(t, letter.setMustBeAt(1))
	assert.Equal(t, letter.state(0), llsUnknown)
	assert.Equal(t, letter.state(1), llsHere)
	assert.Equal(t, letter.state(2), llsUnknown)

	assert.NilError(t, letter.setMustBeAt(0))
	assert.Equal(t, letter.state(0), llsHere)
	assert.Equal(t, letter.state(1), llsHere)
	assert.Equal(t, letter.state(2), llsNotHere)
}

func TestPresentLetterHereThenMaxCountFillsRemainderNotHere(t *testing.T) {
	letter := newPresentLetter(3)

	assert.NilError(t, letter.setMustBeAt(1))
	assert.NilError(t, letter.setRequiredCount(1))

	assert.Equal(t, letter.state(0), llsNotHere)
	assert.Equal(t, letter.state(1), llsHere)
	assert.Equal(t, letter.state(2), llsNotHere)
}

func TestPresentLetterMaxCountThenNotHereFillsRemainderHere(t *testing.T) {
	letter := newPresentLetter(3)

	assert.NilError(t, letter.setMustBeAt(1))
	assert.NilError(t, letter.setRequiredCount(2))
	assert.NilError(t, letter.setMustNotBeAt(0))

	assert.Equal(t, letter.state(0), llsNotHere)
	assert.Equal(t, letter.state(1), llsHere)
	assert.Equal(t, letter.state(2), llsHere)
}

func TestPresentLetterMaxCountLessThanHereErrors(t *testing.T) {
	letter := newPresentLetter(3)

	assert.NilError(t, letter.setMustBeAt(0))
	assert.NilError(t, letter.setMustBeAt(1))
	assert.Error(t, letter.setRequiredCount(1), "Can't set required count to 1 since that would be less than the minimum count (2).")
}

func TestPresentLetterMaxCountMoreThanPossibleErrors(t *testing.T) {
	letter := newPresentLetter(3)

	assert.NilError(t, letter.setMustNotBeAt(0))
	assert.NilError(t, letter.setMustNotBeAt(1))
	assert.Error(t, letter.setRequiredCount(2), "Can't set required count to 2 since it's already 1.")
}

func TestPresentLetterHereAfterNotHereErrors(t *testing.T) {
	letter := newPresentLetter(3)

	assert.NilError(t, letter.setMustNotBeAt(0))
	assert.Error(t, letter.setMustBeAt(0), "Can't set letter to here at index 0 since it's already marked as not here.")
}

func TestPresentLetterNotHereAfterHereErrors(t *testing.T) {
	letter := newPresentLetter(3)

	assert.NilError(t, letter.setMustBeAt(0))
	assert.Error(t, letter.setMustNotBeAt(0), "Can't set letter to not here at index 0 since it's already marked as here.")
}

func TestWordRestrictionsIsSatisfiedByNoRestrictions(t *testing.T) {
	restrictions := InitWordRestrictions(4)

	assert.Assert(t, restrictions.IsSatisfiedBy(WordFromString("abcd")))
	assert.Assert(t, restrictions.IsSatisfiedBy(WordFromString("zzzz")))

	// Wrong length
	assert.Equal(t, restrictions.IsSatisfiedBy(WordFromString("")), false)
	assert.Equal(t, restrictions.IsSatisfiedBy(WordFromString("abcde")), false)
}

func TestWordRestrictionsIsSatisfiedByWithRestrictions(t *testing.T) {
	restrictions := InitWordRestrictions(4)

	assert.NilError(t, restrictions.Update(&GuessResult{
		Guess: WordFromString("abbc"),
		Results: []LetterResult{
			LetterResultPresentNotHere,
			LetterResultPresentNotHere,
			LetterResultCorrect,
			LetterResultNotPresent,
		},
	}))

	assert.Assert(t, restrictions.IsSatisfiedBy(WordFromString("bdba")))
	assert.Assert(t, restrictions.IsSatisfiedBy(WordFromString("dabb")))

	assert.Equal(t, restrictions.IsSatisfiedBy(WordFromString("bbba")), false)
	assert.Equal(t, restrictions.IsSatisfiedBy(WordFromString("bcba")), false)
	assert.Equal(t, restrictions.IsSatisfiedBy(WordFromString("adbd")), false)
	assert.Equal(t, restrictions.IsSatisfiedBy(WordFromString("bdbd")), false)
}

func TestWordRestrictionsWithDuplicateNotHereLetter(t *testing.T) {
	restrictions := InitWordRestrictions(5)

	assert.NilError(t, restrictions.Update(&GuessResult{
		Guess: WordFromString("emcee"),
		Results: []LetterResult{
			LetterResultNotPresent,
			LetterResultNotPresent,
			LetterResultNotPresent,
			LetterResultNotPresent,
			LetterResultCorrect,
		},
	}))

	assert.Equal(t, restrictions.State('e', 0), LetterRestrictionPresentNotHere)
	assert.Equal(t, restrictions.State('e', 1), LetterRestrictionPresentNotHere)
	assert.Equal(t, restrictions.State('e', 2), LetterRestrictionPresentNotHere)
	assert.Equal(t, restrictions.State('e', 3), LetterRestrictionPresentNotHere)
	assert.Equal(t, restrictions.State('e', 4), LetterRestrictionHere)
	assert.Equal(t, restrictions.State('s', 0), LetterRestrictionUnknown)
	assert.Equal(t, restrictions.State('t', 1), LetterRestrictionUnknown)
	assert.Equal(t, restrictions.State('a', 2), LetterRestrictionUnknown)
	assert.Equal(t, restrictions.State('v', 3), LetterRestrictionUnknown)
	assert.Assert(t, restrictions.IsSatisfiedBy(WordFromString("stave")))
}

func TestWordRestrictionsState(t *testing.T) {
	restrictions := InitWordRestrictions(4)

	assert.NilError(t, restrictions.Update(&GuessResult{
		Guess: WordFromString("abbc"),
		Results: []LetterResult{
			LetterResultPresentNotHere,
			LetterResultPresentNotHere,
			LetterResultCorrect,
			LetterResultNotPresent,
		},
	}))

	assert.Equal(t, restrictions.State('a', 0), LetterRestrictionPresentNotHere)
	assert.Equal(t, restrictions.State('a', 1), LetterRestrictionPresentMaybeHere)
	assert.Equal(t, restrictions.State('a', 2), LetterRestrictionPresentNotHere)
	assert.Equal(t, restrictions.State('b', 0), LetterRestrictionPresentMaybeHere)
	assert.Equal(t, restrictions.State('b', 1), LetterRestrictionPresentNotHere)
	assert.Equal(t, restrictions.State('b', 2), LetterRestrictionHere)
	assert.Equal(t, restrictions.State('c', 3), LetterRestrictionNotPresent)
	assert.Equal(t, restrictions.State('c', 0), LetterRestrictionNotPresent)
	assert.Equal(t, restrictions.State('z', 0), LetterRestrictionUnknown)
}

func TestWordRestrictionsIsStateKnown(t *testing.T) {
	restrictions := InitWordRestrictions(4)

	assert.NilError(t, restrictions.Update(&GuessResult{
		Guess: WordFromString("abbc"),
		Results: []LetterResult{
			LetterResultPresentNotHere,
			LetterResultPresentNotHere,
			LetterResultCorrect,
			LetterResultNotPresent,
		},
	}))

	assert.Assert(t, restrictions.IsStateKnown('a', 0))
	assert.Equal(t, restrictions.IsStateKnown('a', 1), false)
	assert.Assert(t, restrictions.IsStateKnown('b', 2))
	assert.Assert(t, restrictions.IsStateKnown('c', 3))
	assert.Assert(t, restrictions.IsStateKnown('c', 0))
	assert.Equal(t, restrictions.IsStateKnown('z', 0), false)
}

func TestWordRestrictionsIsSatisfiedByWithKnownRequiredCount(t *testing.T) {
	restrictions := InitWordRestrictions(4)

	assert.NilError(t, restrictions.Update(&GuessResult{
		Guess: WordFromString("abbc"),
		Results: []LetterResult{
			LetterResultPresentNotHere,
			LetterResultNotPresent,
			LetterResultCorrect,
			LetterResultNotPresent,
		},
	}))

	assert.Assert(t, restrictions.IsSatisfiedBy(WordFromString("edba")))
	assert.Assert(t, restrictions.IsSatisfiedBy(WordFromString("dabe")))
	assert.Assert(t, restrictions.IsSatisfiedBy(WordFromString("daba")))

	assert.Equal(t, restrictions.IsSatisfiedBy(WordFromString("bdba")), false)
	assert.Equal(t, restrictions.IsSatisfiedBy(WordFromString("dcba")), false)
	assert.Equal(t, restrictions.IsSatisfiedBy(WordFromString("adbd")), false)
}

func TestWordRestrictionsIsSatisfiedByWithMinCount(t *testing.T) {
	restrictions := InitWordRestrictions(4)

	assert.NilError(t, restrictions.Update(&GuessResult{
		Guess: WordFromString("abbc"),
		Results: []LetterResult{
			LetterResultPresentNotHere,
			LetterResultPresentNotHere,
			LetterResultCorrect,
			LetterResultNotPresent,
		},
	}))

	assert.Assert(t, restrictions.IsSatisfiedBy(WordFromString("beba")))
	assert.Assert(t, restrictions.IsSatisfiedBy(WordFromString("dabb")))

	assert.Equal(t, restrictions.IsSatisfiedBy(WordFromString("edba")), false)
	assert.Equal(t, restrictions.IsSatisfiedBy(WordFromString("ebbd")), false)
}
