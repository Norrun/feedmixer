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

type SafeIn[T any] interface {
	WaitWrite(v T) (err error)
	Write(v T) (_ bool, err error)
	RawIn() chan<- T
}
type SafeOut[T any] interface {
	WaitRead() T
	Read() (T, bool)
	RawOut() <-chan T
}

type Safe[T any] interface {
	SafeIn[T]
	SafeOut[T]
	Raw() chan T
}
type SafeChan[T any] chan T

func (receiver SafeChan[T]) Out() SafeChanOut[T] {
	return receiver.Raw()
}

func (receiver SafeChan[T]) In() SafeChanIn[T] {
	return receiver.Raw()
}

func (receiver SafeChan[T]) Raw() chan T {
	return receiver
}

type SafeChanIn[T any] chan<- T

func (s SafeChanIn[T]) RawIn() chan<- T {
	return s
}

// WaitWrite implements [SafeIn].
func (s SafeChanIn[T]) WaitWrite(v T) (err error) {
	defer panicToError(&err)

	s <- v
	return nil
}

// Write implements [SafeIn].
func (s SafeChanIn[T]) Write(v T) (_ bool, err error) {
	defer panicToError(&err)
	select {
	case s <- v:
		return true, nil
	default:
		return false, nil

	}
}

type SafeChanOut[T any] <-chan T

func (s SafeChanOut[T]) RawOut() <-chan T {
	return s
}

// Read implements [SafeOut].
func (s SafeChanOut[T]) Read() (T, bool) {
	var emp T
	select {
	case o := <-s:
		return o, true
	default:
		return emp, false
	}
}

// WaitRead implements [SafeOut].
func (s SafeChanOut[T]) WaitRead() T {
	return <-s
}

func (receiver SafeChan[T]) WaitRead() T {
	return <-receiver
}

func (receiver SafeChan[T]) Read() (T, bool) {
	var emp T
	select {
	case o := <-receiver:
		return o, true
	default:
		return emp, false
	}
}

func (receiver SafeChan[T]) WaitWrite(v T) (err error) {
	defer panicToError(&err)

	receiver <- v
	return nil
}

func (receiver SafeChan[T]) Write(v T) (_ bool, err error) {
	defer panicToError(&err)
	select {
	case receiver <- v:
		return true, nil
	default:
		return false, nil

	}
}

func (receiver SafeChan[T]) Close() {
	close(receiver)
}

type TwoWay[I, O any] struct {
	SafeChanIn[I]
	SafeChanOut[O]
}

func panicToError(err *error) {
	r := recover()
	if r != nil {
		*err = fmt.Errorf("%v", r)
	}
}
