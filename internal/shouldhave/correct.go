package shouldhave

import "errors"

type NotNil[T any] struct {
	value T
}

func MustNotNil[T any](value T) NotNil[T] {
	res, err := AsNotNil(value)
	if err != nil {
		panic(err)
	}
	return res
}

func AsNotNil[T any](value T) (NotNil[T], error) {
	var empty NotNil[T]
	if any(value) == nil {
		return empty, errors.New("<nil> not allowed")
	}
	return NotNil[T]{value: value}, nil
}

func (receiver NotNil[T]) Value() T {
	if any(receiver.value) == nil {
		panic("nil refrence")
	}
	return receiver.value
}
