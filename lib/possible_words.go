package go_wordle_solver

type PossibleWords struct {
	words [][]rune
}

func (pw *PossibleWords) Len() int {
	if pw == nil {
		return 0
	}
	return len(pw.words)
}
