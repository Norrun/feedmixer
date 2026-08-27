package wire

import (
	"net/http"
)

type Leaf[T any] interface {
	func() (T, error) | func() (T, bool) | func() (T, Possible[error, bool])
}

func EmptyG[T any]() T {
	var empty T
	return empty
}

type Possible[TA, TB any] struct {
	value any
}

func NewA[TA, TB any](a TA) Possible[TA, TB] {
	return Possible[TA, TB]{value: a}
}
func NewB[TA, TB any](a TB) Possible[TA, TB] {
	return Possible[TA, TB]{value: a}
}

func (receiver Possible[TA, TB]) Do(a func(TA), b func(TB)) {
	switch v := receiver.value.(type) {
	case TA:
		a(v)
	case TB:
		b(v)
	}
}
func (receiver Possible[TA, TB]) Raw() any {
	return receiver.value
}

func (receiver Possible[TA, TB]) TryA() (TA, bool) {
	if val, ok := receiver.value.(TA); ok {
		return val, ok
	}
	return EmptyG[TA](), false
}
func (receiver Possible[TA, TB]) TryB() (TB, bool) {
	if val, ok := receiver.value.(TB); ok {
		return val, ok
	}
	return EmptyG[TB](), false
}

func NewSnitchResponceWriter(w http.ResponseWriter) *SnitchResponseWriter {
	return &SnitchResponseWriter{inner: w}
}

type SnitchResponseWriter struct {
	inner  http.ResponseWriter
	status int
	hasBod bool
}

// Write implements [io.Writer].

func (receiver SnitchResponseWriter) IsWritten() bool {
	return receiver.hasBod
}

func (receiver SnitchResponseWriter) Status() int {
	return receiver.status
}

func (receiver *SnitchResponseWriter) Header() http.Header {
	return receiver.inner.Header()
}

func (receiver *SnitchResponseWriter) Write(b []byte) (int, error) {
	receiver.hasBod = true
	return receiver.inner.Write(b)
}

func (receiver *SnitchResponseWriter) WriteHeader(status int) {
	receiver.status = status
	receiver.WriteHeader(status)
}

func HandlerFallback(handlers ...http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := NewSnitchResponceWriter(w)
		for _, v := range handlers {
			v(sw, r)
			if sw.hasBod || sw.Status() >= 199 {
				break
			}
		}
	})
}
