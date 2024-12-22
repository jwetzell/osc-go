package osc

import (
	"fmt"
	"reflect"
	"testing"
)

func TestOSCBundleEncoding(t *testing.T) {

	testCases := []struct {
		description string
		bundle      OSCBundle
		expected    []byte
	}{
		{
			"simple contents single message",
			OSCBundle{
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

		actual := testCase.bundle.ToBytes()

		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("Test '%s' failed to encode properly", testCase.description)
			fmt.Printf("expected: %v\n", testCase.expected)
			fmt.Printf("actual: %v\n", actual)
		}
	}

}

func TestOSCBundleDecoding(t *testing.T) {
	testCases := []struct {
		description string
		expected    OSCBundle
		bytes       []byte
	}{
		{
			"simple contents single message",
			OSCBundle{
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

		actual, remainingBytes, error := BundleFromBytes(testCase.bytes)

		if error != nil {
			fmt.Println(error)
			t.Errorf("Test '%s' failed to encode properly", testCase.description)
		}

		if len(remainingBytes) > 0 {
			t.Errorf("Test '%s' should not have any remaining bytes", testCase.description)
		}

		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("Test '%s' failed to encode bundle properly", testCase.description)
			fmt.Printf("expected: %v\n", testCase.expected)
			fmt.Printf("actual: %v\n", actual)
		}

	}
}
