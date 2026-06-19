package datautils

import (
	"reflect"
	"sync"
)

func MakeKeyTo[T any, K comparable](key K) KeyTo[T, K] {
	return KeyTo[T, K]{KeyToThe[T](MakeKeyForThe[T]()), key}
}

func MakeKeyForThe[T any]() KeyToThe[T] {
	return KeyToThe[T]{reflect.TypeFor[T]()}
}

type KeyTo[T any, K comparable] struct {
	KeyToThe[T]
	Key K
}

type KeyToSomething[K comparable] KeyTo[any, K]

type KeyToThe[T any] struct {
	valueType reflect.Type
}

func (receiver KeyTo[T, K]) Unspecified() KeyToSomething[K] {
	return KeyToSomething[K]{KeyToThe[any](receiver.KeyToThe), receiver.Key}
}

type ReadableMap[K comparable, T any] interface {
	map[KeyTo[T, K]]T | map[any]T | map[any]any | *sync.Map
}
