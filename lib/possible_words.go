package go_wordle_solver

import "golang.org/x/exp/slices"

type PossibleWords struct {
	words        []Word
	restrictions WordRestrictions
}

func initPossibleWords(words []Word) PossibleWords {
	return PossibleWords{
		slices.Clone(words),
		InitWordRestrictions(uint8(words[0].Len())),
	}
}

func (pw *PossibleWords) Len() int {
	if pw == nil {
		return 0
	}
	return len(pw.words)
}

func (pw *PossibleWords) At(i int) Word {
	return pw.words[i]
}

func (pw *PossibleWords) Filter(gr *GuessResult) error {
	err := pw.restrictions.Update(gr)
	if err != nil {
		return err
	}
	pw.words = filter(pw.words, pw.restrictions.IsSatisfiedBy)
	return nil
}
