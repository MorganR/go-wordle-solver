package go_wordle_solver

type Optional[T any] struct {
	isPresent bool
	value     T
}

func OptionalOf[T any](v T) Optional[T] {
	return Optional[T]{
		true,
		v,
	}
}

func (o *Optional[T]) HasValue() bool {
	return o.isPresent
}

func (o *Optional[T]) Value() T {
	return o.value
}
