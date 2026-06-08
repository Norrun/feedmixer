package datautils

type Pipe[T any] interface {
	PipeIn[T]
	PipeOut[T]
	In() PipeIn[T]
	Out() PipeOut[T]
}
type PipeIn[T any] interface {
	Write(v T) bool
}
type PipeOut[T any] interface {
	Read() (T, bool)
	ReadAll() []T
}

type pipe[T any] struct {
	data []T
	mCap int
}

type PipeAny = Pipe[any]
type PipeInAny = PipeIn[any]
type PipeOutAny = PipeOut[any]

func NewPipe[T any](cap int) Pipe[T] {
	return &pipe[T]{mCap: cap, data: make([]T, 0, cap)}
}

func (receiver *pipe[T]) Write(v T) bool {
	if receiver.mCap <= len(receiver.data) && receiver.mCap > 0 {
		return false
	}
	receiver.data = append(receiver.data, v)
	return true
}

func (receiver *pipe[T]) Read() (T, bool) {
	if len(receiver.data) > 0 {
		res := receiver.data[0]
		receiver.data = receiver.data[1:]
		return res, true
	}
	return ZeroG[T](), false
}

func (receiver *pipe[T]) ReadAll() []T {
	defer func() {
		receiver.data = make([]T, 0, receiver.mCap)
	}()
	return receiver.data
}

type pipeIn[T any] struct {
	pipe *pipe[T]
}

func (receiver pipeIn[T]) Write(v T) bool {
	return receiver.pipe.Write(v)
}

type pipeOut[T any] struct {
	pipe *pipe[T]
}

func (receiver pipeOut[T]) Read() (T, bool) {
	return receiver.pipe.Read()
}
func (receiver pipeOut[T]) ReadAll() []T {
	return receiver.pipe.ReadAll()
}

func (receiver *pipe[T]) In() PipeIn[T] {
	return pipeIn[T]{pipe: receiver}
}
func (receiver *pipe[T]) Out() PipeOut[T] {
	return pipeOut[T]{pipe: receiver}
}

type OmnyPipe[I, O any] struct {
	PipeIn[I]
	PipeOut[O]
}

// might make a stack named Bucket later
