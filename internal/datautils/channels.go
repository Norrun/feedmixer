package datautils

import (
	"fmt"
	"time"
)

func SendAfter[T any](ch chan<- T, val T, delay time.Duration) {
	go func() {
		time.Sleep(delay)
		ch <- val
	}()
}

// WaitWrite implements [SafeIn].
func SafeWrite[T any](ch chan<- T, v T) (err error) {
	defer panicToError(&err)

	ch <- v
	return nil
}

// Write implements [SafeIn].
func SafeTryWrite[T any](ch chan<- T, v T) (_ bool, err error) {
	defer panicToError(&err)
	select {
	case ch <- v:
		return true, nil
	default:
		return false, nil

	}
}

// Read implements [SafeOut].
func TryRead[T any](ch <-chan T) (T, bool) {
	var emp T
	select {
	case o := <-ch:
		return o, true
	default:
		return emp, false
	}
}

func panicToError(err *error) {
	r := recover()
	if r != nil {
		*err = fmt.Errorf("%v", r)
	}
}
