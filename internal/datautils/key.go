package datautils

import (
	"reflect"
	"sync"
)

func MakeKeyTo[T any, K comparable](key K) KeyTo[T, K] {
	return KeyTo[T, K]{KeyToThe[T](MakeKeyToThe[T]()), key}
}

func MakeKeyToThe[T any]() KeyToThe[T] {
	return KeyToThe[T]{}
}

type TypeKey interface {
	Validate(any) bool
}
type TypeKeyValue[K comparable] interface {
	TypeKey
	Key() K
}

type KeyTo[T any, K comparable] struct {
	KeyToThe[T]
	key K
}

func (k KeyToThe[T]) Validate(v any) bool {
	return reflect.TypeOf(v).AssignableTo(reflect.TypeFor[T]())
}

type KeyToThe[T any] struct{}

type ReadableMap[K comparable, T any] interface {
	map[KeyTo[T, K]]T | map[any]T | map[any]any | *sync.Map
}
