package wire

type Result interface {
	Read() ([]struct {
		value any
		code  int
		err   error
	}, bool)
}

type Return interface {
	Write(value any, code int, err error)
}

type Output interface {
	Result
	Return
}

type Params interface {
	Read() (map[string]any, []any)
}

type Args interface {
	Write(map[string]any, []any)
}

type Input interface {
	Params
	Args
}
