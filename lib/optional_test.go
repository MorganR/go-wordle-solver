package go_wordle_solver

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

func ExampleOptional() {
	maybeValue := Optional[int]{}
	fmt.Println(maybeValue.String())

	maybeValue = OptionalOf(5)
	fmt.Println(maybeValue.String())

	// Output:
	// { }
	// { value: 5 }
}

func TestDefaultOptional(t *testing.T) {
	opt := Optional[int]{}

	assert.Equal(t, opt.HasValue(), false)
	assert.Equal(t, opt.Value(), 0)
}

func TestOptionalWithValue(t *testing.T) {
	opt := OptionalOf(52)

	assert.Equal(t, opt.HasValue(), true)
	assert.Equal(t, opt.Value(), 52)
}
