package osc_test

import (
	"testing"

	"github.com/jwetzell/osc-go"
	"github.com/stretchr/testify/require"
)

func TestGoodPacketFromBytes(t *testing.T) {

	testCases := []struct {
		name     string
		expected osc.Packet
		bytes    []byte
	}{
		{
			name: "message with no args",
			expected: &osc.Message{
				Address: "/hello",
				Args:    []osc.Arg{},
			},
			bytes: []byte{47, 104, 101, 108, 108, 111, 0, 0},
		},
		{
			name: "bundle with one message with no args",
			expected: &osc.Bundle{
				TimeTag: osc.TimeTag{
					Seconds:           32,
					FractionalSeconds: 0,
				},
				Contents: []osc.Packet{&osc.Message{Address: "/oscillator/4/frequency", Args: []osc.Arg{{Type: "f", Value: float32(440)}}}},
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
			got, remainingBytes, err := osc.PacketFromBytes(testCase.bytes)

			require.NoError(t, err)
			require.Empty(t, remainingBytes)
			require.Equal(t, testCase.expected, got)
		})
	}
}

func TestBadPacketFromBytes(t *testing.T) {

	testCases := []struct {
		name        string
		bytes       []byte
		errorString string
	}{
		{
			name:        "empty bytes",
			bytes:       []byte{},
			errorString: "cannot create OSC Packet from empty byte array",
		},
		{
			name:        "packet that does not start with / or #",
			bytes:       []byte{0, 1, 2, 3},
			errorString: "OSC Packet must start with # for bundle or / for message",
		},
		{
			name:        "bundle with not enough bytes",
			bytes:       []byte{35, 98, 117, 110, 100, 108, 101, 0},
			errorString: "OSC Bundle has to be at least 20 bytes",
		},
		{
			name:        "message without null terminated address",
			bytes:       []byte{47, 104, 101, 108, 108, 111},
			errorString: "OSC string must be null-terminated",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, _, err := osc.PacketFromBytes(testCase.bytes)

			require.Error(t, err)
			require.EqualError(t, err, testCase.errorString)
		})
	}
}
