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

func TestGoodPacketFromBytes(t *testing.T) {

	testCases := []struct {
		name     string
		expected OSCPacket
		bytes    []byte
	}{
		{
			name: "message with no args",
			expected: &OSCMessage{
				Address: "/hello",
				Args:    []OSCArg{},
			},
			bytes: []byte{47, 104, 101, 108, 108, 111, 0, 0},
		},
		{
			name: "bundle with one message with no args",
			expected: &OSCBundle{
				TimeTag: OSCTimeTag{
					seconds:           32,
					fractionalSeconds: 0,
				},
				Contents: []OSCPacket{&OSCMessage{Address: "/oscillator/4/frequency", Args: []OSCArg{{Type: "f", Value: float32(440)}}}},
			},
			bytes: []byte{35, 98, 117, 110, 100, 108, 101, 0, 0, 0, 0,
				32, 0, 0, 0, 0, 0, 0, 0, 32, 47, 111,
				115, 99, 105, 108, 108, 97, 116, 111, 114, 47, 52,
				47, 102, 114, 101, 113, 117, 101, 110, 99, 121, 0,
				44, 102, 0, 0, 67, 220, 0, 0},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, remainingBytes, err := PacketFromBytes(testCase.bytes)

			if err != nil {
				t.Fatalf("failed to decode properly: %s", err.Error())
			}

			if len(remainingBytes) != 0 {
				t.Fatalf("failed to decode properly, expected no remaining bytes but got: %v", remainingBytes)
			}

			if !reflect.DeepEqual(got, testCase.expected) {
				t.Fatalf("failed to decode properly got '%v', expected '%v'", got, testCase.expected)
			}
		})
	}
}

func TestBadPacketFromBytes(t *testing.T) {

	testCases := []struct {
		name        string
		bytes       []byte
		errorString string
	}{
		{name: "empty bytes",
			bytes:       []byte{},
			errorString: "cannot create OSC Packet from empty byte array",
		},
		{name: "packet that does not start with / or #",
			bytes:       []byte{0, 1, 2, 3},
			errorString: "OSC Packet must start with # for bundle or / for message",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, _, err := PacketFromBytes(testCase.bytes)

			if err == nil {
				t.Fatalf("PacketFromBytes expected to fail but got: %+v", got)
			}

			if err.Error() != testCase.errorString {
				t.Fatalf("PacketFromBytes got error '%s', expected '%s'", err.Error(), testCase.errorString)
			}
		})
	}
}
