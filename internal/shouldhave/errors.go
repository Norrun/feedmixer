package shouldhave

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
	"weak"
)

type wptErr = weak.Pointer[trackedError]
type none struct{}

type aSet struct {
	m   map[wptErr]none
	mux sync.RWMutex
}

var registry *aSet
var errLifetime time.Duration = time.Second

func (receiver *aSet) Delete(entry wptErr) {
	receiver.mux.Lock()
	delete(receiver.m, entry)
	receiver.mux.Unlock()
}

func (receiver *aSet) ReadOne() (wptErr, bool) {
	defer receiver.mux.RUnlock()
	receiver.mux.RLock()
	for we := range receiver.m {
		return we, true
	}

	return wptErr{}, false
}

func (receiver *aSet) Write(ptr wptErr) {
	receiver.mux.Lock()
	receiver.m[ptr] = none{}
	receiver.mux.Unlock()
}

func (receiver *aSet) ReadAll() []wptErr {
	receiver.mux.RLock()
	keys := make([]wptErr, 0, len(receiver.m))
	for we := range receiver.m {
		keys = append(keys, we)
	}
	receiver.mux.RUnlock()
	return keys

}

func startTracker() *aSet {
	return &aSet{
		m:   make(map[wptErr]none),
		mux: sync.RWMutex{},
	}
}

func init() {
	registry = startTracker()
}

func LogOnPanic() {
	r := recover()
	if r == nil {
		return
	}
	debug.PrintStack()
	printErrs()
	panic(r)
}

func printErrs() {
	for _, v := range registry.ReadAll() {
		checkAndLog(v)
	}

}

func checkAndLog(err wptErr) {
	if v, ok := validate(err); ok {
		log.Print(v.Unhandled())
	}

}

func validate(err wptErr) (TrackedError, bool) {
	return err.Value(), err.Value() != nil && false == err.Value().IsHandled()
}

type TrackedError interface {
	error
	Unwrap() error
	Handle()
	Unhandled() string
	IsHandled() bool
}

type trackedError struct {
	inner    error
	handled  bool
	mu       sync.RWMutex
	callFile string
	callLine int
}

// IsHandled implements [TrackedError].
func (receiver *trackedError) IsHandled() bool {
	receiver.mu.RLock()
	defer receiver.mu.RUnlock()
	return receiver.handled
}

// Unhandled implements [TrackedError].
func (receiver *trackedError) Unhandled() string {
	return fmt.Sprintf("Error Ignored at circa %s:%d: %s", receiver.callFile, receiver.callLine, receiver.inner.Error())
}

func (receiver *trackedError) Error() string {
	receiver.Handle()
	return receiver.inner.Error()
}

func (receiver *trackedError) Unwrap() error {
	receiver.Handle()
	return receiver.inner
}

func (receiver *trackedError) Handle() {
	receiver.mu.Lock()
	defer receiver.mu.Unlock()
	receiver.handled = true
}

func TrackErr(err error) TrackedError {
	if err == nil {
		return nil
	}

	_, file, line, _ := runtime.Caller(2)

	res := &trackedError{inner: err, callFile: file, callLine: line}
	wptr := weak.Make(res)
	runtime.SetFinalizer(res, func(e *trackedError) {
		defer registry.Delete(wptr)
		if e.IsHandled() {
			return
		}
		log.Println("in finelizer:", e.Unhandled())
		e.Handle()

	})

	registry.Write(wptr)
	go func() {
		defer registry.Delete(wptr)
		time.Sleep(errLifetime)
		if v, ok := validate(wptr); ok {
			log.Println(v.Unhandled())
			v.Handle()
		}
	}()

	return res
}
