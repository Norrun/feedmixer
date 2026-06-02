package shouldhave

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"time"
	"weak"
)

type wptErr = weak.Pointer[trackedError]

// var registry map[weak.Pointer[trackedError]]errEntry
var fastCh chan weak.Pointer[trackedError]
var slowCh chan weak.Pointer[trackedError]

func init() {

	fastCh = make(chan weak.Pointer[trackedError], 100)
	slowCh = make(chan weak.Pointer[trackedError], 10)
	Relay(fastCh, slowCh,
		time.Second,
		Filter(func(e wptErr) bool { return e.Value() != nil && false == e.Value().handled }),
	)
	startScanner()

}

func LogOnPanic() {

	r := recover()
	if r == nil {
		return
	}
	debug.PrintStack()
	printErrs()

}

func printErrs() {
	for {
		select {
		case t := <-slowCh:
			checkAndLog(t)
			close(slowCh)
		default:
			select {
			case t := <-fastCh:
				checkAndLog(t)
			default:
				return
			}

		}
	}

}

func checkAndLog(err wptErr) {
	if err.Value() != nil && false == err.Value().handled {
		return
	}
	log.Print(err.Value().Unhandled())
}

// Should be called once at the start of the program.
func startScanner() {

	go func() {

		for v := range slowCh {
			err := v.Value()
			if err != nil {
				continue
			}
			log.Println(err.Unhandled())
			err.Handle()

		}
	}()

}

type TrackedError interface {
	error
	Unwrap() error
	Handle()
	Unhandled() string
}

type trackedError struct {
	inner    error
	handled  bool
	callFile string
	callLine int
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
	receiver.handled = true
}

func TrackErr(err error) TrackedError {
	if err == nil {
		return nil
	}

	_, file, line, _ := runtime.Caller(2)

	res := &trackedError{inner: err, callFile: file, callLine: line}
	runtime.SetFinalizer(res, func(e *trackedError) {
		if e.handled {
			return
		}
		log.Println("in finelizer:", e.Unhandled())
		e.Handle()
	})
	select {
	case fastCh <- weak.Make(res):

	default:
		log.Print("Too many ingnored errors")
	}
	return res
}
