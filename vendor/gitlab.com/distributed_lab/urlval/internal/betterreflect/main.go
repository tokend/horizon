package betterreflect

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

type Struct struct {
	src interface{}
}

func NewStruct(src interface{}) *Struct {
	return &Struct{
		src,
	}
}

func (s *Struct) srcvalue() reflect.Value {
	return reflect.Indirect(reflect.ValueOf(s.src))
}

func (s *Struct) srctype() reflect.Type {
	return s.srcvalue().Type()
}

func (s *Struct) NumField() int {
	return s.srcvalue().Type().NumField()
}

func (s *Struct) Tag(i int, key string) string {
	return s.srctype().Field(i).Tag.Get(key)
}

func (s *Struct) Type(i int) reflect.Type {
	return s.srctype().Field(i).Type
}

func (s *Struct) Value(i int) reflect.Value {
	return s.srcvalue().Field(i)
}

func (s *Struct) Set(i int, value interface{}) (err error) {

	// 	unmarshaler, ok := fieldvalue.Interface().(encoding.TextUnmarshaler)
	// 	if ok {
	// 		if err := unmarshaler.UnmarshalText([]byte(value)); err != nil {
	// 			panic(errors.Wrap(err, "failed to unmarshal"))
	// 		}
	// 		continue
	// 	}

	kind := s.Type(i).Kind()
	if s.Type(i).Kind() == reflect.Ptr {
		kind = s.Type(i).Elem().Kind()
	}

	switch kind {
	case reflect.String:
	case reflect.Bool:
		value, err = cast.ToBoolE(value)
	case reflect.Int:
		value, err = cast.ToIntE(value)
	case reflect.Int8:
		value, err = cast.ToInt8E(value)
	case reflect.Int16:
		value, err = cast.ToInt16E(value)
	case reflect.Int32:
		value, err = cast.ToInt32E(value)
	case reflect.Int64:
		value, err = cast.ToInt64E(value)
	case reflect.Uint:
		value, err = cast.ToUintE(value)
	case reflect.Uint8:
		value, err = cast.ToUint8E(value)
	case reflect.Uint16:
		value, err = cast.ToUint16E(value)
	case reflect.Uint32:
		value, err = cast.ToUint32E(value)
	case reflect.Uint64:
		value, err = cast.ToUint64E(value)
	case reflect.Float32:
		value, err = cast.ToFloat32E(value)
	case reflect.Float64:
		value, err = cast.ToFloat64E(value)
	case reflect.Slice: // TODO: major feature!
		panic(fmt.Sprintf("not (yet) implemented: %v", kind))
	case reflect.Complex64, reflect.Complex128, reflect.Chan,
		reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr,
		reflect.Array, reflect.Struct, reflect.UnsafePointer:
	default:
		panic(fmt.Sprintf("unknown field kind: %v", kind))
	}
	if err != nil {
		return err
	}

	if s.Type(i).Kind() == reflect.Ptr {
		zero := reflect.New(s.Type(i).Elem())
		s.Value(i).Set(zero)
		s.Value(i).Elem().Set(reflect.ValueOf(value))
	} else {
		s.Value(i).Set(reflect.ValueOf(value))
	}

	return nil
}
