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
			message:     OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "b", Value: []byte{98, 108, 111, 98}}}},
			expected:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 98, 0, 0, 0, 0, 0, 4, 98, 108, 111, 98},
		},
		{
			description: "simple address True arg",
			message:     OSCMessage{Address: "/hello", Args: []OSCArg{OSCArg{Type: "T", Value: true}}},
			expected:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 84, 0, 0},
		},
		{
			description: "simple address False arg",
			message:     OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "F", Value: false}}},
			expected:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 70, 0, 0},
		},
		{
			description: "simple address color arg",
			message:     OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "r", Value: OSCColor{r: 20, g: 21, b: 22, a: 10}}}},
			expected:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 114, 0, 0, 20, 21, 22, 10},
		},
		{
			description: "simple address nil arg",
			message:     OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "N", Value: nil}}},
			expected:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 78, 0, 0},
		},
		{
			description: "simple address int64 arg",
			message:     OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "h", Value: 281474976710655}}},
			expected:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 104, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255},
		},
		{
			description: "simple address float64 arg",
			message:     OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "d", Value: 12.7654763}}},
			expected: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 100, 0, 0, 0x40, 0x29, 0x87, 0xec, 0x82, 0x74, 0xb9, 0xe6,
			},
		},
		// TODO(jwetzell): get array args working working
		// {
		// 	description: "simple address array arg",
		// 	message: OSCMessage{
		// 		Address: "/hello",
		// 		Args: []OSCArg{
		// 			[]OSCArg{
		//				{Type: "d", Value: 12.7654763},
		// 				{Type: "i", Value: 1000},
		//			},
		// 		},
		// 	},
		// 	expected: []byte{
		// 		47, 104, 101, 108, 108, 111, 0, 0, 44, 91, 100, 105, 93, 0, 0, 0, 0x40, 0x29, 0x87, 0xec, 0x82, 0x74, 0xb9, 0xe6,
		// 		0, 0, 3, 232,
		// 	},
		// },
		{
			description: "osc 1.0 spec example 1",
			message:     OSCMessage{Address: "/oscillator/4/frequency", Args: []OSCArg{{Type: "f", Value: 440}}},
			expected: []byte{
				47, 111, 115, 99, 105, 108, 108, 97, 116, 111, 114, 47, 52, 47, 102, 114, 101, 113, 117, 101, 110, 99, 121, 0, 44,
				102, 0, 0, 67, 220, 0, 0,
			},
		},
		{
			description: "osc 1.0 spec example 2",
			message: OSCMessage{
				Address: "/foo",
				Args: []OSCArg{
					{Type: "i", Value: 1000},
					{Type: "i", Value: -1},
					{Type: "s", Value: "hello"},
					// thanks IEEE 754
					{Type: "f", Value: 1.2339999675750732421875},
					{Type: "f", Value: 5.677999973297119140625},
				},
			},
			expected: []byte{
				47, 102, 111, 111, 0, 0, 0, 0, 44, 105, 105, 115, 102, 102, 0, 0, 0, 0, 3, 232, 255, 255, 255, 255, 104, 101, 108,
				108, 111, 0, 0, 0, 63, 157, 243, 182, 64, 181, 178, 45,
			},
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
