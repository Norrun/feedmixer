package serverutils

import "github.com/Norrun/feedmixer/internal/shouldhave"

type ServerError interface {
	error
	Unwrap() error
	ServiceMessage() string
	UserMessage() string
	StatusCode() int

	//Metadata() map[any]any //Possible
	//ServiceObject() any //Possible
	//InternalCode() int //Possible

}

type serverErr struct {
	inner error
	smsg  string
	umsg  string
	scode int
}

// Error implements [ServerError].
func (s serverErr) Error() string {
	return s.inner.Error()
}

// ServiceMessage implements [ServerError].
func (s serverErr) ServiceMessage() string {
	return s.smsg
}

// StatusCode implements [ServerError].
func (s serverErr) StatusCode() int {
	return s.scode
}

// Unwrap implements [ServerError].
func (s serverErr) Unwrap() error {
	return s.inner
}

// UserMessage implements [ServerError].
func (s serverErr) UserMessage() string {
	return s.umsg
}

func NewTrackedServerError(inner error, serMsg, usrMsg string, status int) error {
	return shouldhave.TrackErr(serverErr{inner: inner, smsg: serMsg, umsg: usrMsg, scode: status})
}

func _test() ServerError {

	return serverErr{}

}
