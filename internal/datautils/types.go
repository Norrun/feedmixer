package datautils

type Result[T any] struct {
	Value T
	Err   error
}

func ResultValue[T any](v T) Result[T] {
	return Result[T]{Value: v}
}

func ResultError[T any](err error) Result[T] {
	return Result[T]{Err: err}
}

func NewResult[T any](v T, err error) Result[T] {
	return Result[T]{v, err}
}

type Option[T any] struct {
	Value T
	Ok    bool
}

func NewOption[T any](v T, ok bool) Option[T] {
	return Option[T]{v, ok}
}

func OptionValue[T any](v T) Option[T] {
	return Option[T]{v, true}
}
func OptionZero[T any]() Option[T] {
	return Option[T]{}
}

type Duo[T0, T1 any] struct {
	V0 T0
	V1 T1
}
type Trio[T0, T1, T2 any] struct {
	Duo[T0, T1]
	V2 T2
}

func NewDuo[T0, T1 any](v0 T0, v1 T1) Duo[T0, T1] {
	return Duo[T0, T1]{
		V0: v0,
		V1: v1,
	}
}

func NewTrio[T0, T1, T2 any](v0 T0, v1 T1, v2 T2) Trio[T0, T1, T2] {
	return Trio[T0, T1, T2]{
		Duo: NewDuo(v0, v1),
		V2:  v2,
	}
}

func ZeroG[T any]() T {
	var empty T
	return empty
}

func _dev() {

}
