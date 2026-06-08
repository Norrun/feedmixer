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
	switch rv.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
		return rv.IsNil()
	}
	return false
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
