package go_wordle_solver

type Word struct {
	runes []rune
}

func WordFromString(s string) Word {
	return Word{[]rune(s)}
}

func (self Word) Len() uint8 {
	return uint8(len(self.runes))
}

func (self Word) String() string {
	return string(self.runes)
}

func (self Word) At(i uint8) rune {
	return self.runes[i]
}

func (self Word) ForEach(fn func(i int, letter rune) bool) {
	for i, l := range self.runes {
		shouldContinue := fn(i, l)
		if !shouldContinue {
			break
		}
	}
}
