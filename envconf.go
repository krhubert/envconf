package envconf

import (
	"fmt"
	"reflect"
	"runtime"
)

// StdListSeparator is default list separator.
// It can be changed only by NewConfig add SetListSeparator call.
var StdListSeparator = ","

var std = &Config{"", StdListSeparator, &decoder{StdListSeparator}}

// SetValues sets config value based on enviroment variable.
func SetValues(v interface{}) error {
	return std.SetValues(v)
}

// Config for process config reading from envrioment variables.
type Config struct {
	prefix    string
	separator string

	d *decoder
}

// NewConfig creates new config with given prefix.
// The PREFIX_ will be add to every enviroment variable name.
func NewConfig(prefix string) *Config {
	return &Config{
		prefix,
		StdListSeparator,
		&decoder{StdListSeparator},
	}
}

// SetListSeparator changes the list separator.
func (c *Config) SetListSeparator(separator string) {
	c.separator = separator
	c.d.separator = separator
}

// SetValues sets config value based on enviroment variable.
func (c *Config) SetValues(v interface{}) (err error) {

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	rv := reflect.ValueOf(v)
	if err := c.isValidValue(rv); err != nil {
		return err
	}
	rv = rv.Elem()
	tv := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		tagval := tv.Field(i).Tag.Get("envconf")
		if tagval == "" {
			continue
		}

		t, err := newTag(c, tagval, tv.Field(i).Type)
		if err != nil {
			return fmt.Errorf("envconf: field %s %s", tv.Field(i).Name, err)
		}

		if err := c.d.decode(rv.Field(i), t.value); err != nil {
			return fmt.Errorf("envconf: field %s %s", tv.Field(i).Name, err)
		}
	}
	return nil
}

func (c *Config) isValidValue(v reflect.Value) error {
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct || v.IsNil() {
		if v.Type() == nil {
			return fmt.Errorf("envconf: nil")
		}

		if v.Type().Kind() != reflect.Ptr {
			return fmt.Errorf("envconf: non-pointer %s", v.Type())
		}
		if v.Type().Kind() != reflect.Struct {
			return fmt.Errorf("envconf: non-struct %s", v.Type())
		}
		return fmt.Errorf("envconf: nil %s", v.Type())
	}
	return nil
}
