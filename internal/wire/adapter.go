package wire

import (
	"fmt"
	"reflect"
)

type IsAdaptable interface {
	encode(any) error
	decode(any) error
}

type Adaptable struct {
	data []any
}
func (receiver Adaptable) encode(v any) error {
	panic("unimplemented")
}

func Encode(v any) Adaptable {
	panic("unimplemented")
}

func (receiver Adaptable) decode(v any) error {

	ptr := reflect.ValueOf(v)
	if ptr.Kind() != reflect.Pointer || ptr.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("Error decoding %s: %T is not a pointer to a struct", "Adaptable", v)
	}
	s := ptr.Elem()

	if len(receiver.data) > s.NumField() {
		return fmt.Errorf("Struct is too small")
	}

	idx := 0
	for _, v := range s.Fields() {
		if len(receiver.data) == idx {
			return fmt.Errorf("Struct is too big")
		}
		if false == v.CanSet() {
			continue
		}
		n := reflect.ValueOf(receiver.data[idx])

		if false == n.Type().AssignableTo(v.Type()) {
			return fmt.Errorf(
				"field %d: got %v want %v",
				idx,
				n.Type(),
				v.Type(),
			)
		}
		v.Set(n)
		idx++

	}
	return nil

}

