package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	// Assign previous to 1 for the first iteration so the sum works out.
	previous, now := 1, 0
	return func() int {
		previous, now = now, previous+now
		return previous
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}