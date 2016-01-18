package envconf

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var (
	errNoSettable   = errors.New("is not settable")
	errInvalidValue = errors.New("has invalid type")
)

type decoder struct {
	separator string
}

func (d *decoder) decode(v reflect.Value, s string) error {
	if !v.CanSet() {
		return errNoSettable
	}

	v = d.indirect(v)
	switch v.Kind() {
	case reflect.Bool:
		return d.decodeBool(v, s)
	case reflect.String:
		return d.decodeString(v, s)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return d.decodeUint(v, s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return d.decodeInt(v, s)
	case reflect.Float32, reflect.Float64:
		return d.decodeFloat(v, s)
	case reflect.Array, reflect.Slice:
		return d.decodeArray(v, s)
	default:
		return errInvalidValue
	}
}

// indirect walks down v allocating pointers as needed, until it gets to a non-pointer.
func (d *decoder) indirect(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	return v
}

func (d *decoder) decodeArray(v reflect.Value, s string) error {
	t := strings.Split(s, d.separator)
	v.Set(reflect.MakeSlice(v.Type(), len(t), len(t)))

	for i := 0; i < len(t); i++ {
		if err := d.decode(v.Index(i), t[i]); err != nil {
			return err
		}
	}
	return nil
}

func (d *decoder) decodeString(v reflect.Value, s string) error {
	v.SetString(s)
	return nil
}

func (d *decoder) decodeBool(v reflect.Value, s string) error {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	v.SetBool(b)
	return nil
}

func (d *decoder) decodeInt(v reflect.Value, s string) error {
	i, err := strconv.ParseInt(s, 0, v.Type().Bits())
	if err != nil {
		return err
	}
	v.SetInt(i)
	return nil
}

func (d *decoder) decodeUint(v reflect.Value, s string) error {
	i, err := strconv.ParseUint(s, 0, v.Type().Bits())
	if err != nil {
		return err
	}
	v.SetUint(i)
	return nil
}

func (d *decoder) decodeFloat(v reflect.Value, s string) error {
	f, err := strconv.ParseFloat(s, v.Type().Bits())
	if err != nil {
		return err
	}
	v.SetFloat(f)
	return nil
}
