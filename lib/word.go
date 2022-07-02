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

type WordIterator struct {
	runes   []rune
	current int
}

func (self Word) AsIterator() WordIterator {
	return WordIterator{self.runes, -1}
}

func (self Word) AsIteratorFrom(start int) WordIterator {
	return WordIterator{self.runes[start:], -1}
}

func (self *WordIterator) Next() bool {
	self.current += 1
	return self.current < len(self.runes)
}

func (self *WordIterator) Get() (index int, letter rune) {
	return self.current, self.runes[self.current]
}
