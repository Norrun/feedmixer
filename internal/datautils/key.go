package datautils

import (
	"reflect"
)

func MakeTypedKey[T any, K comparable](key K) TypeKey[T, K] {
	return TypeKey[T, K]{reflect.TypeFor[T](), key}
}

type TypeKey[T any, K comparable] struct {
	Structure reflect.Type
	Key       K
}

type ReadableMap[K comparable, T any] interface {
	map[TypeKey[T, K]]T | map[any]T | map[any]any
}
