package osc

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringArg(t *testing.T) {
	expected := "test"
	arg := StringArg(expected)
	require.Equal(t, "s", arg.Type)
	require.Equal(t, expected, arg.Value)
}

func TestInt32Arg(t *testing.T) {
	expected := int32(10)
	arg := Int32Arg(expected)
	require.Equal(t, "i", arg.Type)
	require.Equal(t, expected, arg.Value)
}

func TestFloatArg(t *testing.T) {
	expected := float32(3.14)
	arg := FloatArg(expected)
	require.Equal(t, "f", arg.Type)
	require.Equal(t, expected, arg.Value)
}

func TestBlobArg(t *testing.T) {
	expected := []byte{0x01, 0x02, 0x03}
	arg := BlobArg(expected)
	require.Equal(t, "b", arg.Type)
	require.True(t, reflect.DeepEqual(arg.Value, expected))
}

func TestTrueArg(t *testing.T) {
	arg := TrueArg()
	require.Equal(t, "T", arg.Type)
	require.Equal(t, true, arg.Value)
}

func TestFalseArg(t *testing.T) {
	arg := FalseArg()
	require.Equal(t, "F", arg.Type)
	require.Equal(t, false, arg.Value)
}

func TestNilArg(t *testing.T) {
	arg := NilArg()
	require.Equal(t, "N", arg.Type)
	require.Equal(t, nil, arg.Value)
}

func TestColorArg(t *testing.T) {
	expected := Color{R: 255, G: 123, B: 222, A: 111}
	arg := ColorArg(expected.R, expected.G, expected.B, expected.A)

	require.Equal(t, "r", arg.Type)
	require.Equal(t, expected, arg.Value)
}

func TestInt64Arg(t *testing.T) {
	expected := int64(100)
	arg := Int64Arg(expected)
	require.Equal(t, "h", arg.Type)
	require.Equal(t, expected, arg.Value)
}

func TestDoubleArg(t *testing.T) {
	expected := float64(3.14159)
	arg := DoubleArg(expected)
	require.Equal(t, "d", arg.Type)
	require.Equal(t, expected, arg.Value)
}

func TestTimeTagArg(t *testing.T) {
	expected := TimeTag{Seconds: 12345, FractionalSeconds: 67890}
	arg := TimeTagArg(expected.Seconds, expected.FractionalSeconds)

	require.Equal(t, "t", arg.Type)
	require.Equal(t, expected, arg.Value)
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

			require.NoError(t, err)
			require.Equal(t, testCase.expected, got)
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
			_, err := ArgFromStringAndType(testCase.rawArg, testCase.oscType)

			require.Error(t, err)
			require.EqualError(t, err, testCase.errorMsg)
		})
	}
}
