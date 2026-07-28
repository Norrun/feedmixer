package utils

import (
	"runtime"
	"sync"
)

func FanOut[TIn, TOut any](done <-chan struct{}, source <-chan TIn, load, perTheread int, process func(in TIn) TOut) []<-chan TOut {

	num := min(load, runtime.NumCPU()*perTheread)
	res := make([]<-chan TOut, num)
	for i := range num {
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

func Fan[TIn, TOut any](done <-chan struct{}, load, perTheread int, process func(in TIn) TOut, sources ...<-chan TIn) []<-chan TOut {
	return FanOut(done, FanIn(done, sources...), max(len(sources), load), perTheread, process)
}

func FanIn[T any](done <-chan struct{}, sources ...<-chan T) <-chan T {
	wg := sync.WaitGroup{}
	fannedIn := make(chan T)
	transfer := func(c <-chan T) {
		defer wg.Done()
		for v := range c {
			select {
			case <-done:
				return
			case fannedIn <- v:
			}
		}
	}
	for _, v := range sources {
		wg.Add(1)
		go transfer(v)
	}
	go func() {
		wg.Wait()
		close(fannedIn)
	}()
	return fannedIn
}
