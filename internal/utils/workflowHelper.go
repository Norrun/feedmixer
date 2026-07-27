package utils

import (
	"runtime"
)

func FanOut[TIn, TOut any](done <-chan struct{}, source <-chan TIn, load, perTheread int, process func(in TIn) TOut) []<-chan TOut {

	num := min(load, runtime.NumCPU()*perTheread)
	res := make([]<-chan TOut, num)
	for i := 0; i < num; i++ {
		channel := make(chan TOut)
		go func() {
			select {
			case <-done:
				return
			case v := <-source:
				channel <- process(v)
			}
		}()
		res[i] = channel
	}
	return res
}
