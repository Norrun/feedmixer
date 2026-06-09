package wire

import (
	"bytes"
	"net/http"
)

type ApproveResponseWriter struct {
	inner     http.ResponseWriter
	status    int
	approved  bool
	tmpBuff   bytes.Buffer
	predicate func(*ApproveResponseWriter) bool
}

func NewApprovalWriter(inner http.ResponseWriter, predicate func(*ApproveResponseWriter) bool) *ApproveResponseWriter {
	return &ApproveResponseWriter{inner: inner, predicate: predicate}
}

func (receiver *ApproveResponseWriter) Header() http.Header {
	return receiver.inner.Header()
}

func (receiver *ApproveResponseWriter) Write(b []byte) (int, error) {
	defer receiver.ApplyPredicate()
	if receiver.approved {
		return receiver.inner.Write(b)
	}
	return receiver.tmpBuff.Write(b)
}

func (receiver *ApproveResponseWriter) ApplyPredicate() {
	if receiver.predicate == nil {
		return
	}
	if receiver.predicate(receiver) {
		receiver.Approve()
	}
}

func (receiver *ApproveResponseWriter) WriteHeader(status int) {
	defer receiver.ApplyPredicate()
	if receiver.approved {
		receiver.WriteHeader(status)
		return
	}
	receiver.status = status

}

func (receiver *ApproveResponseWriter) Approve() {
	if receiver.approved {
		return
	}
	receiver.approved = true
	if receiver.tmpBuff.Len() > 0 {
		receiver.inner.Write(receiver.tmpBuff.Bytes())
	}
	if receiver.status > 99 {
		receiver.WriteHeader(receiver.status)
	}
}

func (receiver *ApproveResponseWriter) Approved() bool {
	return receiver.approved
}

func (receiver *ApproveResponseWriter) Status() int {
	return receiver.status
}

func (receiver *ApproveResponseWriter) ClearCache() {
	receiver.tmpBuff.Reset()
}
