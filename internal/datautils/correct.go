package datautils

import (
	"reflect"
)

func CheckNil(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)

	for {
		switch rv.Kind() {
		case reflect.Interface, reflect.Pointer:
			if rv.IsNil() {
				return true
			}
			rv = rv.Elem()

		case reflect.Map,
			reflect.Slice,
			reflect.Chan,
			reflect.Func:
			return rv.IsNil()

		default:
			return false
		}
	}
}

func CheckValidZero(v any) int {
	t := reflect.ValueOf(v)
	if t.IsValid() {
		if t.IsZero() {
			return 0
		}
		return 1
	}
	return -1
}
