package datautils

func ConvertSlice[Tin, Tout any](in []Tin, to func(Tin) Tout) []Tout {
	var res []Tout
	for _, v := range in {
		res = append(res, to(v))
	}
	return res
}

func ConvertSliceErr[Tin, Tout any](in []Tin, to func(Tin) (Tout, error)) ([]Tout, error) {
	var res []Tout
	for _, v := range in {
		c, err := to(v)
		if err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func AnySlice[T any](a []T) []any {
	return ConvertSlice(a, func(v T) any { return any(v) })
}
