package shouldhave

type Unimplemented struct{}

func (receiver Unimplemented) Error() string {
	return "an unimplemented type struct"
}
