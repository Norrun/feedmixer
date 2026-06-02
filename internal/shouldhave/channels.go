package shouldhave

import (
	"time"
)

func SendAfter[T any](ch chan<- T, val T, delay time.Duration) {
	go func() {
		time.Sleep(delay)
		ch <- val
	}()
}

func Relay[T any](in <-chan T, out chan<- T, delay time.Duration, ops ...func(T) (T, bool)) {

	go func() {
		time.Sleep(delay)
		for v := range in {
			go func() {
				i := v
				for _, op := range ops {
					t, ok := op(i)
					if ok {
						i = t
					} else {
						return
					}
				}
				out <- i
			}()
			time.Sleep(delay)
		}
	}()
}

func Mod[T any](mod func(T) T) func(T) (T, bool) {
	return func(t T) (T, bool) {
		return mod(t), true
	}
}
func Effect[T any](fx func(T)) func(T) (T, bool) {
	return func(t T) (T, bool) {
		fx(t)
		return t, true
	}
}
func Filter[T any](f func(T) bool) func(T) (T, bool) {
	return func(t T) (T, bool) {
		return t, f(t)
	}
}
