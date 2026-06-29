package datautils

func ConvertSlice[Tin, Tout any](in []Tin, to func(Tin) Tout) []Tout {
	var res []Tout
	for _, v := range in {
		res = append(res, to(v))
	}
	return res
}

func AnySlice[T any](a []T) []any {
	return ConvertSlice(a, func(v T) any { return any(v) })
}
