package envconf

import (
	"os"
	"strings"
	"testing"
)

type testSimple struct {
	String     string    `envconf:"string,,true"`
	Bool       bool      `envconf:"bool,,true"`
	Int        int       `envconf:"int,,true"`
	StringPtr  *string   `envconf:"string_ptr,,true"`
	IntArray   []int     `envconf:"int_array,,true"`
	FloatArray []float64 `envconf:"float_array,,true"`
}

func TestSimple(t *testing.T) {
	var env testSimple
	os.Clearenv()
	os.Setenv("STRING", "test")
	os.Setenv("BOOL", "true")
	os.Setenv("INT", "1")
	os.Setenv("STRING_PTR", "test_ptr")
	os.Setenv("INT_ARRAY", "0,1")
	os.Setenv("FLOAT_ARRAY", "-0.1,0.1")

	if err := SetValues(&env); err != nil {
		t.Fatal(err)
	}

	if env.String != "test" {
		t.Fatalf("set string invalid: want %s, got %s", "test", env.String)
	}
	if env.Bool != true {
		t.Fatalf("set bool invalid: want %t, got %t", true, env.Bool)
	}
	if env.Int != 1 {
		t.Fatalf("set int invalid: want %d, got %d", 1, env.Int)
	}
	if *env.StringPtr != "test_ptr" {
		t.Fatalf("set string_ptr invalid: want %s, got %s", "test_ptr", *env.StringPtr)
	}
	if len(env.IntArray) != 2 {
		t.Fatalf("invalid int_array length: want %d, got %d", 2, len(env.IntArray))
	}
	if env.IntArray[0] != 0 || env.IntArray[1] != 1 {
		t.Fatalf("set int_array invalid: want [%d, %d], got [%d, %d]", 0, 1, env.IntArray[0], env.IntArray[1])
	}
	if len(env.FloatArray) != 2 {
		t.Fatalf("invalid float__array length: want %d, got %d", 2, len(env.FloatArray))
	}
	if env.FloatArray[0] != -0.1 || env.FloatArray[1] != 0.1 {
		t.Fatalf("set int_array invalid: want [%f, %f], got [%f, %f]", -0.1, 0.1, env.FloatArray[0], env.FloatArray[1])
	}
}

type testDefault struct {
	String string `envconf:"string,foobar"`
}

func TestDefault(t *testing.T) {
	var env testDefault
	os.Clearenv()
	if err := SetValues(&env); err != nil {
		t.Fatal(err)
	}

	if env.String != "foobar" {
		t.Fatalf("set string invalid: want %s, got %s", "foobar", env.String)
	}
}

type testPrefix struct {
	String string `envconf:"string"`
}

func TestPrefix(t *testing.T) {
	var env testPrefix
	os.Clearenv()
	os.Setenv("APP_STRING", "test")

	conf := NewConfig("app")
	if err := conf.SetValues(&env); err != nil {
		t.Fatal(err)
	}

	if env.String != "test" {
		t.Fatalf("set string invalid: want %s, got %s", "test", env.String)
	}
}

type testListSeparator struct {
	StringArray []string `envconf:"string_array"`
}

func TestListSeparator(t *testing.T) {
	var env testListSeparator
	os.Clearenv()
	os.Setenv("STRING_ARRAY", "a;b")

	conf := NewConfig("")
	conf.SetListSeparator(";")
	if err := conf.SetValues(&env); err != nil {
		t.Fatal(err)
	}
	if len(env.StringArray) != 2 {
		t.Fatalf("invalid string_array length: want %d, got %d", 2, len(env.StringArray))
	}
	if env.StringArray[0] != "a" || env.StringArray[1] != "b" {
		t.Fatalf("set string invalid: want [%s, %s], got [%s, %s]", "a", "b", env.StringArray[0], env.StringArray[1])
	}
}

type testDefaultArray struct {
	BoolArray []bool `envconf:"bool_array,false;true"`
}

func TestDefaultArray(t *testing.T) {
	var env testDefaultArray
	os.Clearenv()
	if err := SetValues(&env); err != nil {
		t.Fatal(err)
	}

	if len(env.BoolArray) != 2 {
		t.Fatalf("invalid bool_array length: want %d, got %d", 2, len(env.BoolArray))
	}
	if env.BoolArray[0] != false || env.BoolArray[1] != true {
		t.Fatalf("set bool_array invalid: want [%t, %t], got [%t, %t]", false, true, env.BoolArray[0], env.BoolArray[1])
	}
}

type testDefaultAndRequiredError struct {
	String string `envconf:"string,foobar,true"`
}

func TestDefaultAndRequiredError(t *testing.T) {
	var env testDefaultAndRequiredError
	if err := SetValues(&env); err == nil {
		t.Fatalf("expect field String required not allowd with default value")
	} else if !strings.Contains(err.Error(), "field String required not allowd with default value") {
		t.Fatalf("expect error: want \"field String required not allowd with default value\", got %v", err)
	}
}

type testRequiredError struct {
	String string `envconf:"string,,true"`
}

func TestRequiredError(t *testing.T) {
	var env testRequiredError
	if err := SetValues(&env); err == nil {
		t.Fatalf("expect env STRING not found, but required")
	} else if !strings.Contains(err.Error(), "env STRING not found, but require") {
		t.Fatalf("expect error: want \"env STRING not found, but required\", got %v", err)
	}
}

type testMissingNameTagError struct {
	String string `envconf:",,true"`
}

func TestMissingNameTagError(t *testing.T) {
	var env testMissingNameTagError
	if err := SetValues(&env); err == nil {
		t.Fatalf("expect field String missing name")
	} else if !strings.Contains(err.Error(), "field String missing name") {
		t.Fatalf("expect error: want \"field String missing name\", got %v", err)
	}
}

type testNoSettableError struct {
	s string `envconf:"string"`
}

func TestNoSettableError(t *testing.T) {
	var env testNoSettableError
	os.Clearenv()
	os.Setenv("STRING", "test")

	if err := SetValues(&env); err == nil {
		t.Fatalf("expect %v error", errNoSettable)
	} else if !strings.Contains(err.Error(), errNoSettable.Error()) {
		t.Fatalf("expect error: want %v, got %v", errNoSettable, err)
	}
}

type testInvalidValue struct {
	Struct struct{} `envconf:"struct"`
}

func TestInvalidValue(t *testing.T) {
	var env testInvalidValue
	os.Clearenv()
	os.Setenv("STRUCT", "t")

	if err := SetValues(&env); err == nil {
		t.Fatalf("expect %v error", errInvalidValue)
	} else if !strings.Contains(err.Error(), errInvalidValue.Error()) {
		t.Fatalf("expect error: want %v, got %v", errInvalidValue, err)
	}
}
