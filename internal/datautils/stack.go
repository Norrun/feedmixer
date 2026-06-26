package datautils

func Peek[T any](s []T) T {
	return s[len(s)-1]
}

func Pop[T any](s []T) T {
	v := Peek(s)
	s = s[:len(s)-1]
	return v
}
