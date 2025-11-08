package cheenv

import (
	"os"
	"testing"
	"time"

	"github.com/comfortablynumb/che/pkg/chetest"
)

func TestGet(t *testing.T) {
	key := "TEST_STRING"
	os.Setenv(key, "hello")
	defer os.Unsetenv(key)

	chetest.RequireEqual(t, Get(key, "default"), "hello")
	chetest.RequireEqual(t, Get("NON_EXISTENT", "default"), "default")
}

func TestMustGet(t *testing.T) {
	key := "TEST_MUST_STRING"
	os.Setenv(key, "hello")
	defer os.Unsetenv(key)

	chetest.RequireEqual(t, MustGet(key), "hello")
}

func TestMustGet_Panic(t *testing.T) {
	defer func() {
		r := recover()
		chetest.RequireEqual(t, r != nil, true)
	}()

	MustGet("NON_EXISTENT_KEY")
}

func TestGetInt(t *testing.T) {
	key := "TEST_INT"
	os.Setenv(key, "42")
	defer os.Unsetenv(key)

	chetest.RequireEqual(t, GetInt(key, 10), 42)
	chetest.RequireEqual(t, GetInt("NON_EXISTENT", 10), 10)

	os.Setenv(key, "invalid")
	chetest.RequireEqual(t, GetInt(key, 10), 10)
}

func TestMustGetInt(t *testing.T) {
	key := "TEST_MUST_INT"
	os.Setenv(key, "42")
	defer os.Unsetenv(key)

	chetest.RequireEqual(t, MustGetInt(key), 42)
}

func TestMustGetInt_Panic(t *testing.T) {
	defer func() {
		r := recover()
		chetest.RequireEqual(t, r != nil, true)
	}()

	MustGetInt("NON_EXISTENT_INT")
}

func TestMustGetInt_PanicOnInvalidValue(t *testing.T) {
	key := "TEST_INVALID_INT"
	os.Setenv(key, "not_a_number")
	defer os.Unsetenv(key)

	defer func() {
		r := recover()
		chetest.RequireEqual(t, r != nil, true)
	}()

	MustGetInt(key)
}

func TestGetInt64(t *testing.T) {
	key := "TEST_INT64"
	os.Setenv(key, "9223372036854775807")
	defer os.Unsetenv(key)

	chetest.RequireEqual(t, GetInt64(key, 10), int64(9223372036854775807))
	chetest.RequireEqual(t, GetInt64("NON_EXISTENT", 10), int64(10))
}

func TestGetFloat(t *testing.T) {
	key := "TEST_FLOAT"
	os.Setenv(key, "3.14159")
	defer os.Unsetenv(key)

	chetest.RequireEqual(t, GetFloat(key, 0.0), 3.14159)
	chetest.RequireEqual(t, GetFloat("NON_EXISTENT", 1.5), 1.5)
}

func TestGetBool(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
	}{
		{"true", true},
		{"TRUE", true},
		{"True", true},
		{"1", true},
		{"yes", true},
		{"YES", true},
		{"on", true},
		{"ON", true},
		{"y", true},
		{"t", true},
		{"false", false},
		{"FALSE", false},
		{"0", false},
		{"no", false},
		{"NO", false},
		{"off", false},
		{"n", false},
		{"f", false},
	}

	key := "TEST_BOOL"
	for _, tt := range tests {
		os.Setenv(key, tt.value)
		result := GetBool(key, false)
		chetest.RequireEqual(t, result, tt.expected)
	}
	os.Unsetenv(key)

	// Test default value
	chetest.RequireEqual(t, GetBool("NON_EXISTENT", true), true)
	chetest.RequireEqual(t, GetBool("NON_EXISTENT", false), false)

	// Test invalid value returns default
	os.Setenv(key, "invalid")
	chetest.RequireEqual(t, GetBool(key, true), true)
	os.Unsetenv(key)
}

func TestMustGetBool(t *testing.T) {
	key := "TEST_MUST_BOOL"
	os.Setenv(key, "true")
	defer os.Unsetenv(key)

	chetest.RequireEqual(t, MustGetBool(key), true)
}

func TestMustGetBool_Panic(t *testing.T) {
	defer func() {
		r := recover()
		chetest.RequireEqual(t, r != nil, true)
	}()

	MustGetBool("NON_EXISTENT_BOOL")
}

func TestGetDuration(t *testing.T) {
	key := "TEST_DURATION"
	os.Setenv(key, "5m")
	defer os.Unsetenv(key)

	chetest.RequireEqual(t, GetDuration(key, time.Second), 5*time.Minute)
	chetest.RequireEqual(t, GetDuration("NON_EXISTENT", time.Second), time.Second)

	os.Setenv(key, "invalid")
	chetest.RequireEqual(t, GetDuration(key, time.Second), time.Second)
}

func TestGetStringList(t *testing.T) {
	key := "TEST_STRING_LIST"
	os.Setenv(key, "a,b,c")
	defer os.Unsetenv(key)

	result := GetStringList(key, ",", []string{"default"})
	chetest.RequireEqual(t, len(result), 3)
	chetest.RequireEqual(t, result[0], "a")
	chetest.RequireEqual(t, result[1], "b")
	chetest.RequireEqual(t, result[2], "c")

	// Test with whitespace
	os.Setenv(key, " a , b , c ")
	result = GetStringList(key, ",", []string{"default"})
	chetest.RequireEqual(t, len(result), 3)
	chetest.RequireEqual(t, result[0], "a")
	chetest.RequireEqual(t, result[1], "b")
	chetest.RequireEqual(t, result[2], "c")

	// Test default
	result = GetStringList("NON_EXISTENT", ",", []string{"default"})
	chetest.RequireEqual(t, len(result), 1)
	chetest.RequireEqual(t, result[0], "default")

	// Test empty value
	os.Setenv(key, "")
	result = GetStringList(key, ",", []string{"default"})
	chetest.RequireEqual(t, len(result), 1)
	chetest.RequireEqual(t, result[0], "default")
}

func TestGetIntList(t *testing.T) {
	key := "TEST_INT_LIST"
	os.Setenv(key, "1,2,3")
	defer os.Unsetenv(key)

	result := GetIntList(key, ",", []int{0})
	chetest.RequireEqual(t, len(result), 3)
	chetest.RequireEqual(t, result[0], 1)
	chetest.RequireEqual(t, result[1], 2)
	chetest.RequireEqual(t, result[2], 3)

	// Test with whitespace
	os.Setenv(key, " 1 , 2 , 3 ")
	result = GetIntList(key, ",", []int{0})
	chetest.RequireEqual(t, len(result), 3)
	chetest.RequireEqual(t, result[0], 1)

	// Test default on non-existent
	result = GetIntList("NON_EXISTENT", ",", []int{99})
	chetest.RequireEqual(t, len(result), 1)
	chetest.RequireEqual(t, result[0], 99)

	// Test default on invalid value
	os.Setenv(key, "1,invalid,3")
	result = GetIntList(key, ",", []int{99})
	chetest.RequireEqual(t, len(result), 1)
	chetest.RequireEqual(t, result[0], 99)
}

func TestSet(t *testing.T) {
	key := "TEST_SET"
	err := Set(key, "value")
	defer os.Unsetenv(key)

	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, os.Getenv(key), "value")
}

func TestUnset(t *testing.T) {
	key := "TEST_UNSET"
	os.Setenv(key, "value")

	err := Unset(key)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, os.Getenv(key), "")
}

func TestHas(t *testing.T) {
	key := "TEST_HAS"
	os.Setenv(key, "value")
	defer os.Unsetenv(key)

	chetest.RequireEqual(t, Has(key), true)
	chetest.RequireEqual(t, Has("NON_EXISTENT"), false)

	// Test empty value
	os.Setenv(key, "")
	chetest.RequireEqual(t, Has(key), true)
}

func TestGetAll(t *testing.T) {
	key1 := "TEST_GET_ALL_1"
	key2 := "TEST_GET_ALL_2"
	os.Setenv(key1, "value1")
	os.Setenv(key2, "value2")
	defer os.Unsetenv(key1)
	defer os.Unsetenv(key2)

	all := GetAll()

	chetest.RequireEqual(t, all[key1], "value1")
	chetest.RequireEqual(t, all[key2], "value2")
}

func TestGetWithPrefix(t *testing.T) {
	key1 := "APP_NAME"
	key2 := "APP_VERSION"
	key3 := "OTHER_KEY"

	os.Setenv(key1, "myapp")
	os.Setenv(key2, "1.0.0")
	os.Setenv(key3, "value")
	defer os.Unsetenv(key1)
	defer os.Unsetenv(key2)
	defer os.Unsetenv(key3)

	result := GetWithPrefix("APP_")

	chetest.RequireEqual(t, len(result), 2)
	chetest.RequireEqual(t, result["NAME"], "myapp")
	chetest.RequireEqual(t, result["VERSION"], "1.0.0")
	chetest.RequireEqual(t, result["OTHER_KEY"], "")
}

func TestParseBool(t *testing.T) {
	trueValues := []string{"true", "TRUE", "True", "1", "yes", "YES", "on", "ON", "y", "t"}
	for _, v := range trueValues {
		result, err := parseBool(v)
		chetest.RequireEqual(t, err, nil)
		chetest.RequireEqual(t, result, true)
	}

	falseValues := []string{"false", "FALSE", "False", "0", "no", "NO", "off", "OFF", "n", "f"}
	for _, v := range falseValues {
		result, err := parseBool(v)
		chetest.RequireEqual(t, err, nil)
		chetest.RequireEqual(t, result, false)
	}

	invalidValues := []string{"invalid", "2", "maybe", ""}
	for _, v := range invalidValues {
		_, err := parseBool(v)
		chetest.RequireEqual(t, err != nil, true)
	}
}
