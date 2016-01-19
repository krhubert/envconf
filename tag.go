package envconf

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

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
	if err := t.getValue(); err != nil {
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
	return nil
}

func (t *tag) getValue() error {
	s, b := os.LookupEnv(t.name)
	if t.required && !b {
		return fmt.Errorf("env %s not found, but required", t.name)
	}

	t.value = s
	return nil
}
