package wire

import (
	"errors"
	"fmt"
)

type Leaf[T any] interface {
	func() (T, error) | func() (T, bool)
}

func NewLeaf[T any, F Leaf[T]](f F) func(int) (T, error) {
	var empty T
	switch r := any(f).(type) {
	case func() (T, error):
		return func(i int) (T, error) {
			res, err := r()
			if err != nil {
				return empty, fmt.Errorf("Error at %d: %w", i, err)
			}
			return res, nil
		}
	case func() (T, bool):
		return func(i int) (T, error) {
			res, ok := r()
			if ok {
				return res, nil
			}
			return empty, fmt.Errorf("False at %d", i)
		}
	default:
		return func(i int) (T, error) {
			return empty, errors.New("Invalid function")
		}

	}

}

func Fallback[T any](leafs ...func(int) (T, error)) (T, error) {
	var empty T
	var errl error
	for i, v := range leafs {
		res, err := v(i)
		if err != nil {
			errl = errors.Join(errl, err)
			continue
		}
		return res, nil

	}
	return empty, errl
}
