package envconf

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
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

type tag struct {
	name     string
	value    string
	defvalue string
	required bool
}

func newTag(c *Config, s string, v reflect.Type) (*tag, error) {
	var t = &tag{}

	a := strings.Split(s, ",")
	switch len(a) {
	case 3:
		t.required = a[2] == "true"
		fallthrough
	case 2:
		t.setDefvalue(c.separator, a[1], v)
		fallthrough
	case 1:
		t.setName(c.prefix, a[0])
	case 0:
		return nil, errors.New("empty tag")
	default:
		return nil, errors.New("invalid tag value")
	}
	if err := t.isValid(); err != nil {
		return nil, err
	}
	if t.value == "" {
		t.value = t.defvalue
	}
	return t, nil
}

func (t *tag) setName(prefix, name string) {
	if prefix == "" {
		t.name = strings.ToUpper(name)
	} else {
		t.name = strings.ToUpper(prefix + "_" + name)
	}
}

func (t *tag) setDefvalue(separator, value string, v reflect.Type) {
	if value == "" {
		return
	}

	t.defvalue = value
	if v.Kind() == reflect.Array || v.Kind() == reflect.Slice {
		// Replace all ; to config separator so decoder could parse it
		t.defvalue = strings.Replace(t.defvalue, ";", separator, -1)
	}
}

func (t *tag) isValid() error {
	if t.required && t.name == "" {
		return fmt.Errorf("missing name")
	}

	if t.required && t.defvalue != "" {
		return fmt.Errorf("required not allowd with default value")
	}

	s, b := os.LookupEnv(t.name)
	if t.required && !b {
		return fmt.Errorf("env %s not found, but required", t.name)
	}

	t.value = s
	return nil
}
