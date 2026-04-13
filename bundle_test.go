package osc

import (
	"reflect"
	"testing"
)

func TestOSCBundleEncoding(t *testing.T) {

	testCases := []struct {
		name     string
		bundle   *OSCBundle
		expected []byte
	}{
		{
			"simple contents single message",
			&OSCBundle{
				TimeTag: OSCTimeTag{
					seconds:           32,
					fractionalSeconds: 0,
				},
				Contents: []OSCPacket{&OSCMessage{Address: "/oscillator/4/frequency", Args: []OSCArg{{Type: "f", Value: float32(440)}}}},
			},
			[]byte{35, 98, 117, 110, 100, 108, 101, 0, 0, 0, 0,
				32, 0, 0, 0, 0, 0, 0, 0, 32, 47, 111,
				115, 99, 105, 108, 108, 97, 116, 111, 114, 47, 52,
				47, 102, 114, 101, 113, 117, 101, 110, 99, 121, 0,
				44, 102, 0, 0, 67, 220, 0, 0},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			got, err := testCase.bundle.ToBytes()

			if err != nil {
				t.Fatalf("failed to encode properly: %s", err.Error())
			}

			if !reflect.DeepEqual(got, testCase.expected) {
				t.Fatalf("failed to encode properly got '%v', expected '%v'", got, testCase.expected)
			}
		})
	}
}

func TestOSCBundleDecoding(t *testing.T) {
	testCases := []struct {
		name     string
		expected *OSCBundle
		bytes    []byte
	}{
		{
			name: "simple contents single message",
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
			actual, remainingBytes, error := BundleFromBytes(testCase.bytes)

			if error != nil {
				t.Fatalf("failed to decode properly: %s", error.Error())
			}

			if len(remainingBytes) > 0 {
				t.Fatalf("should not have any remaining bytes")
			}

			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Fatalf("failed to decode properly got '%v', expected '%v'", actual, testCase.expected)
			}
		})
	}
}
