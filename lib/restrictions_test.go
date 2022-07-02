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

	assert.Assert(t, restrictions.IsSatisfiedBy("abcd"))
	assert.Assert(t, restrictions.IsSatisfiedBy("zzzz"))

	// Wrong length
	assert.Equal(t, restrictions.IsSatisfiedBy(""), false)
	assert.Equal(t, restrictions.IsSatisfiedBy("abcde"), false)
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

	assert.Assert(t, restrictions.IsSatisfiedBy("bdba"))
	assert.Assert(t, restrictions.IsSatisfiedBy("dabb"))

	assert.Equal(t, restrictions.IsSatisfiedBy("bbba"), false)
	assert.Equal(t, restrictions.IsSatisfiedBy("bcba"), false)
	assert.Equal(t, restrictions.IsSatisfiedBy("adbd"), false)
	assert.Equal(t, restrictions.IsSatisfiedBy("bdbd"), false)
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

	assert.Equal(t, restrictions.State(&LocatedLetter{Letter: 'a', Location: 0}), LetterRestrictionPresentNotHere)
	assert.Equal(t, restrictions.State(&LocatedLetter{Letter: 'a', Location: 1}), LetterRestrictionPresentMaybeHere)
	assert.Equal(t, restrictions.State(&LocatedLetter{Letter: 'a', Location: 2}), LetterRestrictionPresentNotHere)
	assert.Equal(t, restrictions.State(&LocatedLetter{Letter: 'b', Location: 0}), LetterRestrictionPresentMaybeHere)
	assert.Equal(t, restrictions.State(&LocatedLetter{Letter: 'b', Location: 1}), LetterRestrictionPresentNotHere)
	assert.Equal(t, restrictions.State(&LocatedLetter{Letter: 'b', Location: 2}), LetterRestrictionHere)
	assert.Equal(t, restrictions.State(&LocatedLetter{Letter: 'c', Location: 3}), LetterRestrictionNotPresent)
	assert.Equal(t, restrictions.State(&LocatedLetter{Letter: 'c', Location: 0}), LetterRestrictionNotPresent)
	assert.Equal(t, restrictions.State(&LocatedLetter{Letter: 'z', Location: 0}), LetterRestrictionUnknown)
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

	assert.Assert(t, restrictions.IsStateKnown(&LocatedLetter{Letter: 'a', Location: 0}))
	assert.Equal(t, restrictions.IsStateKnown(&LocatedLetter{Letter: 'a', Location: 1}), false)
	assert.Assert(t, restrictions.IsStateKnown(&LocatedLetter{Letter: 'b', Location: 2}))
	assert.Assert(t, restrictions.IsStateKnown(&LocatedLetter{Letter: 'c', Location: 3}))
	assert.Assert(t, restrictions.IsStateKnown(&LocatedLetter{Letter: 'c', Location: 0}))
	assert.Equal(t, restrictions.IsStateKnown(&LocatedLetter{Letter: 'z', Location: 0}), false)
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

	assert.Assert(t, restrictions.IsSatisfiedBy("edba"))
	assert.Assert(t, restrictions.IsSatisfiedBy("dabe"))
	assert.Assert(t, restrictions.IsSatisfiedBy("daba"))

	assert.Equal(t, restrictions.IsSatisfiedBy("bdba"), false)
	assert.Equal(t, restrictions.IsSatisfiedBy("dcba"), false)
	assert.Equal(t, restrictions.IsSatisfiedBy("adbd"), false)
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

	assert.Assert(t, restrictions.IsSatisfiedBy("beba"))
	assert.Assert(t, restrictions.IsSatisfiedBy("dabb"))

	assert.Equal(t, restrictions.IsSatisfiedBy("edba"), false)
	assert.Equal(t, restrictions.IsSatisfiedBy("ebbd"), false)
}

func TestWordRestrictionsEmptyThenMerge(t *testing.T) {
	restrictions := InitWordRestrictions(4)
	otherRestrictions := InitWordRestrictions(4)
	assert.NilError(t, otherRestrictions.Update(&GuessResult{
		Guess: WordFromString("abbc"),
		Results: []LetterResult{
			LetterResultPresentNotHere,
			LetterResultPresentNotHere,
			LetterResultCorrect,
			LetterResultNotPresent,
		},
	}))

	assert.NilError(t, restrictions.Merge(&otherRestrictions))

	assert.Assert(t, restrictions.IsSatisfiedBy("babd"))
	assert.Assert(t, restrictions.IsSatisfiedBy("baba"))
	assert.Equal(t, restrictions.IsSatisfiedBy("babc"), false)
	assert.Equal(t, restrictions.IsSatisfiedBy("badb"), false)
	assert.Equal(t, restrictions.IsSatisfiedBy("adbb"), false)
	assert.Equal(t, restrictions.IsSatisfiedBy("dbba"), false)
}

func TestWordRestrictionsMerge(t *testing.T) {
	restrictions := InitWordRestrictions(4)
	otherRestrictions := InitWordRestrictions(4)
	assert.NilError(t, restrictions.Update(&GuessResult{
		Guess: WordFromString("bade"),
		Results: []LetterResult{
			LetterResultCorrect,
			LetterResultCorrect,
			LetterResultNotPresent,
			LetterResultCorrect,
		},
	}))
	assert.NilError(t, otherRestrictions.Update(&GuessResult{
		Guess: WordFromString("abbc"),
		Results: []LetterResult{
			LetterResultPresentNotHere,
			LetterResultPresentNotHere,
			LetterResultCorrect,
			LetterResultNotPresent,
		},
	}))

	assert.NilError(t, restrictions.Merge(&otherRestrictions))

	assert.Assert(t, restrictions.IsSatisfiedBy("babe"))
	assert.Equal(t, restrictions.IsSatisfiedBy("baee"), false)
}

func TestWordRestrictionsMergeWrongLength(t *testing.T) {
	restrictions := InitWordRestrictions(4)
	otherRestrictions := InitWordRestrictions(5)

	assert.Error(t, restrictions.Merge(&otherRestrictions), "Can't merge restrictions with different word lengths (has: 4, received: 5).")
}

func TestWordRestrictionsConflictingMergePresentThenNotPresent(t *testing.T) {
	restrictions := InitWordRestrictions(4)
	otherRestrictions := InitWordRestrictions(4)
	assert.NilError(t, restrictions.Update(&GuessResult{
		Guess: WordFromString("abcd"),
		Results: []LetterResult{
			LetterResultPresentNotHere,
			LetterResultPresentNotHere,
			LetterResultPresentNotHere,
			LetterResultNotPresent,
		},
	}))
	assert.NilError(t, otherRestrictions.Update(&GuessResult{
		Guess: WordFromString("abbc"),
		Results: []LetterResult{
			LetterResultPresentNotHere,
			LetterResultPresentNotHere,
			LetterResultCorrect,
			LetterResultNotPresent,
		},
	}))

	assert.Error(t, restrictions.Merge(&otherRestrictions), "Can't merge incompatible restrictions.")
}

func TestWordRestrictionsConflictingMergeNotPresentThenPresent(t *testing.T) {
	restrictions := InitWordRestrictions(4)
	otherRestrictions := InitWordRestrictions(4)
	assert.NilError(t, restrictions.Update(&GuessResult{
		Guess: WordFromString("abcd"),
		Results: []LetterResult{
			LetterResultNotPresent,
			LetterResultPresentNotHere,
			LetterResultNotPresent,
			LetterResultNotPresent,
		},
	}))
	assert.NilError(t, otherRestrictions.Update(&GuessResult{
		Guess: WordFromString("abbc"),
		Results: []LetterResult{
			LetterResultPresentNotHere,
			LetterResultPresentNotHere,
			LetterResultCorrect,
			LetterResultNotPresent,
		},
	}))

	assert.Error(t, restrictions.Merge(&otherRestrictions), "Can't merge incompatible restrictions.")
}

func TestWordRestrictionsConflictingMergePresentDifferentPlace(t *testing.T) {
	restrictions := InitWordRestrictions(4)
	otherRestrictions := InitWordRestrictions(4)
	assert.NilError(t, restrictions.Update(&GuessResult{
		Guess: WordFromString("abbc"),
		Results: []LetterResult{
			LetterResultPresentNotHere,
			LetterResultCorrect,
			LetterResultCorrect,
			LetterResultNotPresent,
		},
	}))
	assert.NilError(t, otherRestrictions.Update(&GuessResult{
		Guess: WordFromString("abbc"),
		Results: []LetterResult{
			LetterResultPresentNotHere,
			LetterResultPresentNotHere,
			LetterResultCorrect,
			LetterResultNotPresent,
		},
	}))

	assert.Error(t, restrictions.Merge(&otherRestrictions), "Can't set letter to not here at index 1 since it's already marked as here.")
}
