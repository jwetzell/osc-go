package osc

import (
	"reflect"
	"testing"
)

func TestStringArg(t *testing.T) {
	expected := "test"
	arg := StringArg(expected)
	if arg.Type != "s" {
		t.Errorf("Expected type 's', got '%s'", arg.Type)
	}
	if arg.Value != expected {
		t.Errorf("Expected value '%s', got '%v'", expected, arg.Value)
	}
}

func TestInt32Arg(t *testing.T) {
	expected := int32(10)
	arg := Int32Arg(expected)
	if arg.Type != "i" {
		t.Errorf("Expected type 'i', got '%s'", arg.Type)
	}
	if arg.Value != expected {
		t.Errorf("Expected value '%d', got '%v'", expected, arg.Value)
	}
}

func TestFloatArg(t *testing.T) {
	expected := float32(3.14)
	arg := FloatArg(expected)
	if arg.Type != "f" {
		t.Errorf("Expected type 'f', got '%s'", arg.Type)
	}
	if arg.Value != expected {
		t.Errorf("Expected value '%f', got '%v'", expected, arg.Value)
	}
}

func TestBlobArg(t *testing.T) {
	expected := []byte{0x01, 0x02, 0x03}
	arg := BlobArg(expected)
	if arg.Type != "b" {
		t.Errorf("Expected type 'b', got '%s'", arg.Type)
	}

	if !reflect.DeepEqual(arg.Value, expected) {
		t.Errorf("Expected value '%v', got '%v'", expected, arg.Value)
	}
}

func TestTrueArg(t *testing.T) {
	arg := TrueArg()
	if arg.Type != "T" {
		t.Errorf("Expected type 'T', got '%s'", arg.Type)
	}
	if arg.Value != true {
		t.Errorf("Expected value '%v', got '%v'", true, arg.Value)
	}
}

func TestFalseArg(t *testing.T) {
	arg := FalseArg()
	if arg.Type != "F" {
		t.Errorf("Expected type 'F', got '%s'", arg.Type)
	}
	if arg.Value != false {
		t.Errorf("Expected value '%v', got '%v'", false, arg.Value)
	}
}

func TestNilArg(t *testing.T) {
	arg := NilArg()
	if arg.Type != "N" {
		t.Errorf("Expected type 'N', got '%s'", arg.Type)
	}
	if arg.Value != nil {
		t.Errorf("Expected value '%v', got '%v'", nil, arg.Value)
	}
}

func TestColorArg(t *testing.T) {
	expected := Color{R: 255, G: 123, B: 222, A: 111}
	arg := ColorArg(expected.R, expected.G, expected.B, expected.A)

	if arg.Type != "r" {
		t.Errorf("Expected type 'r', got '%s'", arg.Type)
	}
	if !reflect.DeepEqual(arg.Value, expected) {
		t.Errorf("Expected value '%v', got '%v'", expected, arg.Value)
	}
}

func TestInt64Arg(t *testing.T) {
	expected := int64(100)
	arg := Int64Arg(expected)
	if arg.Type != "h" {
		t.Errorf("Expected type 'h', got '%s'", arg.Type)
	}
	if arg.Value != expected {
		t.Errorf("Expected value '%d', got '%v'", expected, arg.Value)
	}
}

func TestDoubleArg(t *testing.T) {
	expected := float64(3.14159)
	arg := DoubleArg(expected)
	if arg.Type != "d" {
		t.Errorf("Expected type 'd', got '%s'", arg.Type)
	}
	if arg.Value != expected {
		t.Errorf("Expected value '%f', got '%v'", expected, arg.Value)
	}
}

func TestTimeTagArg(t *testing.T) {
	expected := TimeTag{Seconds: 12345, FractionalSeconds: 67890}
	arg := TimeTagArg(expected.Seconds, expected.FractionalSeconds)

	if arg.Type != "t" {
		t.Errorf("Expected type 't', got '%s'", arg.Type)
	}
	if !reflect.DeepEqual(arg.Value, expected) {
		t.Errorf("Expected value '%v', got '%v'", expected, arg.Value)
	}
}

func TestGoodArgFromStringAndType(t *testing.T) {
	testCases := []struct {
		name     string
		rawArg   string
		oscType  string
		expected Arg
	}{
		{
			name:    "string arg",
			rawArg:  "hello",
			oscType: "s",
			expected: Arg{
				Type:  "s",
				Value: "hello",
			},
		},
		{
			name:    "int32 arg",
			rawArg:  "42",
			oscType: "i",
			expected: Arg{
				Type:  "i",
				Value: int32(42),
			},
		},
		{
			name:    "float32 arg",
			rawArg:  "42.1",
			oscType: "f",
			expected: Arg{
				Type:  "f",
				Value: float32(42.1),
			},
		},
		{
			name:    "blob arg",
			rawArg:  "68656c6c6f", // "hello" in hex
			oscType: "b",
			expected: Arg{
				Type:  "b",
				Value: []byte("hello"),
			},
		},
		{
			name:    "int64 arg",
			rawArg:  "42",
			oscType: "h",
			expected: Arg{
				Type:  "h",
				Value: int64(42),
			},
		},
		{
			name:    "float64 arg",
			rawArg:  "42.1",
			oscType: "d",
			expected: Arg{
				Type:  "d",
				Value: float64(42.1),
			},
		},
		{
			name:    "true arg",
			rawArg:  "",
			oscType: "T",
			expected: Arg{
				Type:  "T",
				Value: true,
			},
		},
		{
			name:    "false arg",
			rawArg:  "",
			oscType: "F",
			expected: Arg{
				Type:  "F",
				Value: false,
			},
		},
		{
			name:    "nil arg",
			rawArg:  "",
			oscType: "N",
			expected: Arg{
				Type:  "N",
				Value: nil,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := ArgFromStringAndType(testCase.rawArg, testCase.oscType)

			if err != nil {
				t.Fatalf("failed to encode properly: %s", err.Error())
			}

			if !reflect.DeepEqual(got, testCase.expected) {
				t.Fatalf("failed to encode properly got '%v', expected '%v'", got, testCase.expected)
			}
		})
	}
}

func TestBadArgFromStringAndType(t *testing.T) {
	testCases := []struct {
		name     string
		rawArg   string
		oscType  string
		errorMsg string
	}{
		{
			name:     "invalid int32",
			rawArg:   "abc",
			oscType:  "i",
			errorMsg: "strconv.ParseInt: parsing \"abc\": invalid syntax",
		},
		{
			name:     "invalid float32",
			rawArg:   "abc",
			oscType:  "f",
			errorMsg: "strconv.ParseFloat: parsing \"abc\": invalid syntax",
		},
		{
			name:     "invalid blob",
			rawArg:   "zz",
			oscType:  "b",
			errorMsg: "encoding/hex: invalid byte: U+007A 'z'",
		},
		{
			name:     "invalid int64",
			rawArg:   "abc",
			oscType:  "h",
			errorMsg: "strconv.ParseInt: parsing \"abc\": invalid syntax",
		},
		{
			name:     "invalid float64",
			rawArg:   "abc",
			oscType:  "d",
			errorMsg: "strconv.ParseFloat: parsing \"abc\": invalid syntax",
		},
		{
			name:     "unsupported OSC type",
			rawArg:   "42",
			oscType:  "x",
			errorMsg: "unsupported OSC arg type: x",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := ArgFromStringAndType(testCase.rawArg, testCase.oscType)

			if err == nil {
				t.Fatalf("expected error but got: %+v", got)
			}

			if err.Error() != testCase.errorMsg {
				t.Fatalf("expected error '%v', got '%v'", testCase.errorMsg, err.Error())
			}
		})
	}
}
