package osc

import (
	"reflect"
	"testing"
)

func TestGoodOSCArgsToBuffer(t *testing.T) {

	testCases := []struct {
		name     string
		args     []OSCArg
		expected []byte
	}{
		{
			name: "int arg",
			args: []OSCArg{
				{
					Type:  "i",
					Value: int(123),
				},
			},
			expected: []byte{0, 0, 0, 123},
		},
		{
			name: "int32 arg",
			args: []OSCArg{
				{
					Type:  "i",
					Value: int32(123),
				},
			},
			expected: []byte{0, 0, 0, 123},
		},
		{
			name: "float32 arg",
			args: []OSCArg{
				{
					Type:  "f",
					Value: float32(123),
				},
			},
			expected: []byte{66, 246, 0, 0},
		},
		{
			name: "float32 arg with int value",
			args: []OSCArg{
				{
					Type:  "f",
					Value: int(123),
				},
			},
			expected: []byte{66, 246, 0, 0},
		},
		{
			name: "float32 arg with int32 value",
			args: []OSCArg{
				{
					Type:  "f",
					Value: int32(123),
				},
			},
			expected: []byte{66, 246, 0, 0},
		},
		{
			name: "float32 arg with int64 value",
			args: []OSCArg{
				{
					Type:  "f",
					Value: int64(123),
				},
			},
			expected: []byte{66, 246, 0, 0},
		},
		{
			name: "float64 arg",
			args: []OSCArg{
				{
					Type:  "d",
					Value: float64(123),
				},
			},
			expected: []byte{64, 94, 192, 0, 0, 0, 0, 0},
		},
		{
			name: "float64 arg with float32 value",
			args: []OSCArg{
				{
					Type:  "d",
					Value: float32(123),
				},
			},
			expected: []byte{64, 94, 192, 0, 0, 0, 0, 0},
		},
		{
			name: "float64 arg with int value",
			args: []OSCArg{
				{
					Type:  "d",
					Value: int(123),
				},
			},
			expected: []byte{64, 94, 192, 0, 0, 0, 0, 0},
		},
		{
			name: "float64 arg with int32 value",
			args: []OSCArg{
				{
					Type:  "d",
					Value: int32(123),
				},
			},
			expected: []byte{64, 94, 192, 0, 0, 0, 0, 0},
		},
		{
			name: "float64 arg with int64 value",
			args: []OSCArg{
				{
					Type:  "d",
					Value: int64(123),
				},
			},
			expected: []byte{64, 94, 192, 0, 0, 0, 0, 0},
		},
		{
			name: "int64 arg",
			args: []OSCArg{
				{
					Type:  "h",
					Value: int64(123),
				},
			},
			expected: []byte{0, 0, 0, 0, 0, 0, 0, 123},
		},
		{
			name: "int64 arg with int32 value",
			args: []OSCArg{
				{
					Type:  "h",
					Value: int32(123),
				},
			},
			expected: []byte{0, 0, 0, 0, 0, 0, 0, 123},
		},
		{
			name: "blob arg",
			args: []OSCArg{
				{
					Type:  "b",
					Value: []byte{1, 2, 3},
				},
			},
			expected: []byte{0, 0, 0, 3, 1, 2, 3, 0},
		},
		{
			name: "true arg",
			args: []OSCArg{
				{
					Type:  "T",
					Value: true,
				},
			},
			expected: []byte{},
		},
		{
			name: "false arg",
			args: []OSCArg{
				{
					Type:  "F",
					Value: false,
				},
			},
			expected: []byte{},
		},
		{
			name: "nil arg",
			args: []OSCArg{
				{
					Type:  "N",
					Value: nil,
				},
			},
			expected: []byte{},
		},
		{
			name: "inifinitum arg",
			args: []OSCArg{
				{
					Type:  "I",
					Value: nil,
				},
			},
			expected: []byte{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			got, err := argsToBuffer(testCase.args)

			if err != nil {
				t.Fatalf("failed to encode properly: %s", err.Error())
			}

			if !reflect.DeepEqual(got, testCase.expected) {
				t.Fatalf("failed to encode properly got '%v', expected '%v'", got, testCase.expected)
			}
		})
	}
}

func TestBadOSCArgsToBuffer(t *testing.T) {

	testCases := []struct {
		name        string
		args        []OSCArg
		errorString string
	}{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			got, err := argsToBuffer(testCase.args)

			if err == nil {
				t.Fatalf("argsToBuffer expected to fail but got: %+v", got)
			}

			if err.Error() != testCase.errorString {
				t.Fatalf("argsToBuffer got error '%s', expected '%s'", err.Error(), testCase.errorString)
			}
		})
	}
}
