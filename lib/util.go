package go_wordle_solver

func fillSlice[T any](s []T, v T) {
	for i := range s {
		s[i] = v
	}
}

func allValues[T any](s []T, fn func(T) bool) bool {
	for _, v := range s {
		if !fn(v) {
			return false
		}
	}
	return true
}

func allPairs[K comparable, V any](m map[K]V, fn func(K, V) bool) bool {
	for k, v := range m {
		if !fn(k, v) {
			return false
		}
	}
	return true
}

func allLetters(s string, fn func(rune) bool) bool {
	for _, l := range s {
		if !fn(l) {
			return false
		}
	}
	return true
}
