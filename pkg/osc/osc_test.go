package osc

import (
	"fmt"
	"reflect"
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestOSCEncoding(t *testing.T) {

	testCases := []struct {
		description string
		message     OSCMessage
		expected    []byte
	}{
		{
			"simple hello",
			OSCMessage{
				Address: "/hello",
				Args:    []OSCArg{},
			},
			[]byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 0, 0, 0},
		},
		{
			"simple address string arg",
			OSCMessage{
				Address: "/hello",
				Args: []OSCArg{
					{
						Type:  "s",
						Value: "arg1",
					},
				},
			},
			[]byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 115, 0, 0, 97, 114, 103, 49, 0, 0, 0, 0},
		},
		{
			description: "simple address integer arg",
			message:     OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "i", Value: 35}}},
			expected:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 105, 0, 0, 0, 0, 0, 35},
		},
		{
			description: "simple address float arg",
			message:     OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "f", Value: 34.5}}},
			expected:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 102, 0, 0, 66, 10, 0, 0},
		},
		{
			description: "simple address blob arg",
			message:     OSCMessage{Address: "/hello", Args: []OSCArg{OSCArg{Type: "b", Value: []byte{98, 108, 111, 98}}}},
			expected:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 98, 0, 0, 0, 0, 0, 4, 98, 108, 111, 98},
		},
	}

	for _, testCase := range testCases {

		actual := ToBytes(testCase.message)

		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("Test '%s' failed to encode properly", testCase.description)
			fmt.Printf("expected: %v\n", testCase.expected)
			fmt.Printf("actual: %v\n", actual)
		}
	}

}
