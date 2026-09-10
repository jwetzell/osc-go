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
	expected := Color{r: 255, g: 0, b: 0, a: 255}
	arg := ColorArg(expected.r, expected.b, expected.g, expected.r)

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
	expected := TimeTag{seconds: 12345, fractionalSeconds: 67890}
	arg := TimeTagArg(expected.seconds, expected.fractionalSeconds)

	if arg.Type != "t" {
		t.Errorf("Expected type 't', got '%s'", arg.Type)
	}
	if !reflect.DeepEqual(arg.Value, expected) {
		t.Errorf("Expected value '%v', got '%v'", expected, arg.Value)
	}
}
