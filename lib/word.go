package go_wordle_solver

import "golang.org/x/exp/slices"

type Word struct {
	runes []rune
}

func WordFromString(s string) Word {
	return Word{[]rune(s)}
}

func (self Word) Equal(w Word) bool {
	return slices.Compare(self.runes, w.runes) == 0
}

func (self Word) Len() int {
	return len(self.runes)
}

func (self Word) String() string {
	return string(self.runes)
}

func (self Word) At(i int) rune {
	return self.runes[i]
}

func (self Word) AllLetters(fn func(rune) bool) bool {
	return allValues(self.runes, fn)
}
