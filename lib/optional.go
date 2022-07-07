package go_wordle_solver

import (
	"fmt"
	"reflect"
)

// Optional represents a value that may or may not be present, similar to Rust's [Option] or Java's
// [Optional].
//
// This enables use of optional data without pushing it onto the heap and using a pointer. It
// should only be used with trivially copyable types.
//
// [Option]: https://doc.rust-lang.org/std/option/
// [Optional]: https://docs.oracle.com/javase/8/docs/api/java/util/Optional.html
type Optional[T any] struct {
	isPresent bool
	value     T
}

// OptionalOf creates an optional with a value.
func OptionalOf[T any](v T) Optional[T] {
	return Optional[T]{
		true,
		v,
	}
}

// HasValue returns true iff this optional has a value.
func (o *Optional[T]) HasValue() bool {
	return o.isPresent
}

// Value returns the stored value. If this optional has no value, this will return the type's
// default value.
func (o *Optional[T]) Value() T {
	return o.value
}

// Equal returns true iff this optional equals the other. Values are compared with
// `reflect.DeepEqual`.
func (o *Optional[T]) Equal(other *Optional[T]) bool {
	return o.isPresent == other.isPresent && reflect.DeepEqual(o.value, other.value)
}

// String converts optional (and its value, if present) to a readable string.
func (o *Optional[T]) String() string {
	if o.isPresent {
		return fmt.Sprintf("{ value: %v }", o.value)
	}
	return "{ }"
}
