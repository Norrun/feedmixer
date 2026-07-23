package datautils

import (
	"errors"
	"reflect"
)

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
	if CheckNil(value) {
		return empty, errors.New("<nil> not allowed")
	}
	return NotNil[T]{value: value}, nil
}

func (receiver NotNil[T]) Value() T {
	if CheckNil(receiver.value) {
		panic("nil refrence")
	}
	return receiver.value
}

func CheckNil(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)

	for {
		switch rv.Kind() {
		case reflect.Interface, reflect.Pointer:
			if rv.IsNil() {
				return true
			}
			rv = rv.Elem()

		case reflect.Map,
			reflect.Slice,
			reflect.Chan,
			reflect.Func:
			return rv.IsNil()

		default:
			return false
		}
	}
}

func CheckValidZero(v any) int {
	t := reflect.ValueOf(v)
	if t.IsValid() {
		if t.IsZero() {
			return 0
		}
		return 1
	}
	return -1
}
