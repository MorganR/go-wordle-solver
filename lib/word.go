package go_wordle_solver

import "golang.org/x/exp/slices"

// A Word represents a read-only string, optimized to work with runes instead of bytes.
type Word struct {
	runes []rune
}

// WordFromString constructs a word from a given string.
//
// This requires allocating a slice, so it's best to convert strings to [Word]s once, and then only
// use Words from there.
func WordFromString(s string) Word {
	return Word{[]rune(s)}
}

// Equal determines if this word is equal to the given word.
func (self Word) Equal(w Word) bool {
	return slices.Compare(self.runes, w.runes) == 0
}

// Len returns the number of runes (i.e. letters) in this word.
func (self Word) Len() int {
	return len(self.runes)
}

// String converts this word back to a string.
func (self Word) String() string {
	return string(self.runes)
}

// At returns the rune (i.e. letter) at the given index in the word.
func (self Word) At(i int) rune {
	return self.runes[i]
}

// AllLetters determines if all the letters in this word satisfy the given function.
func (self Word) AllLetters(fn func(rune) bool) bool {
	return allValues(self.runes, fn)
}
