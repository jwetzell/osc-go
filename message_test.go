package osc

import (
	"math"
	"reflect"
	"testing"
)

func TestGoodOSCMessageEncoding(t *testing.T) {

	testCases := []struct {
		name     string
		message  *OSCMessage
		expected []byte
	}{
		{
			name: "simple hello",
			message: &OSCMessage{
				Address: "/hello",
				Args:    []OSCArg{},
			},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 0, 0, 0},
		},
		{
			name: "simple address string arg",
			message: &OSCMessage{
				Address: "/hello",
				Args: []OSCArg{
					{
						Type:  "s",
						Value: "arg1",
					},
				},
			},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 115, 0, 0, 97, 114, 103, 49, 0, 0, 0, 0},
		},
		{
			name:     "simple address integer arg",
			message:  &OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "i", Value: 35}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 105, 0, 0, 0, 0, 0, 35},
		},
		{
			name:     "simple address float arg",
			message:  &OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "f", Value: 34.5}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 102, 0, 0, 66, 10, 0, 0},
		},
		{
			name:     "simple address blob arg",
			message:  &OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "b", Value: []byte{98, 108, 111, 98}}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 98, 0, 0, 0, 0, 0, 4, 98, 108, 111, 98},
		},
		{
			name:     "simple address True arg",
			message:  &OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "T", Value: true}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 84, 0, 0},
		},
		{
			name:     "simple address False arg",
			message:  &OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "F", Value: false}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 70, 0, 0},
		},
		{
			name:     "simple address color arg",
			message:  &OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "r", Value: OSCColor{r: 20, g: 21, b: 22, a: 10}}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 114, 0, 0, 20, 21, 22, 10},
		},
		{
			name:     "simple address nil arg",
			message:  &OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "N", Value: nil}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 78, 0, 0},
		},
		{
			name:     "simple address int64 arg",
			message:  &OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "h", Value: 281474976710655}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 104, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255},
		},
		{
			name:    "simple address float64 arg",
			message: &OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "d", Value: 12.7654763}}},
			expected: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 100, 0, 0, 0x40, 0x29, 0x87, 0xec, 0x82, 0x74, 0xb9, 0xe6,
			},
		},
		// TODO(jwetzell): get array args working working
		// {
		// 	name: "simple address array arg",
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
			name:    "osc 1.0 spec example 1",
			message: &OSCMessage{Address: "/oscillator/4/frequency", Args: []OSCArg{{Type: "f", Value: 440}}},
			expected: []byte{
				47, 111, 115, 99, 105, 108, 108, 97, 116, 111, 114, 47, 52, 47, 102, 114, 101, 113, 117, 101, 110, 99, 121, 0, 44,
				102, 0, 0, 67, 220, 0, 0,
			},
		},
		{
			name: "osc 1.0 spec example 2",
			message: &OSCMessage{
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
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.message.ToBytes()

			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Fatalf("failed to encode properly got '%v', expected '%v'", actual, testCase.expected)
			}
		})
	}

}

func TestGoodOSCMessageDecoding(t *testing.T) {
	testCases := []struct {
		name     string
		bytes    []byte
		expected OSCMessage
	}{
		{
			name:     "simple address no args",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 0, 0, 0},
			expected: OSCMessage{Address: "/hello", Args: []OSCArg{}},
		},
		{
			name:     "simple address string arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 115, 0, 0, 97, 114, 103, 49, 0, 0, 0, 0},
			expected: OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "s", Value: "arg1"}}},
		},
		{
			name:     "simple address integer arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 105, 0, 0, 0, 0, 0, 35},
			expected: OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "i", Value: int32(35)}}},
		},
		{
			name:     "simple address float arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 102, 0, 0, 66, 10, 0, 0},
			expected: OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "f", Value: float32(34.5)}}},
		},
		{
			name:     "simple address blob arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 98, 0, 0, 0, 0, 0, 4, 98, 108, 111, 98},
			expected: OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "b", Value: []byte{98, 108, 111, 98}}}},
		},
		{
			name:     "simple address True arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 84, 0, 0},
			expected: OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "T", Value: true}}},
		},
		{
			name:     "simple address False arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 70, 0, 0},
			expected: OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "F", Value: false}}},
		},
		{
			name:     "simple address color arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 114, 0, 0, 20, 21, 22, 10},
			expected: OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "r", Value: OSCColor{r: 20, g: 21, b: 22, a: 10}}}},
		},
		{
			name:     "simple address nil arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 78, 0, 0},
			expected: OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "N", Value: nil}}},
		},
		{
			name:     "simple address Inifinitum arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 73, 0, 0},
			expected: OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "I", Value: math.MaxInt32}}},
		},
		{
			name:     "simple address int64 arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 104, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255},
			expected: OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "h", Value: int64(281474976710655)}}},
		},
		{
			name: "simple address float64 arg",
			bytes: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 100, 0, 0, 0x40, 0x29, 0x87, 0xec, 0x82, 0x74, 0xb9, 0xe6,
			},
			expected: OSCMessage{Address: "/hello", Args: []OSCArg{{Type: "d", Value: float64(12.7654763)}}},
		},
		// TODO(jwetzell): support OSC array
		// {
		// 	name: "simple address array arg",
		// 	bytes: []byte{
		// 		47, 104, 101, 108, 108, 111, 0, 0, 44, 91, 100, 105, 93, 0, 0, 0, 0x40, 0x29, 0x87, 0xec, 0x82, 0x74, 0xb9, 0xe6,
		// 		0, 0, 3, 232,
		// 	},
		// 	expected: OSCMessage{
		// 		Address: "/hello",
		// 		Args: []OSCArg{
		// 			[]OSCArg{
		// 				{Type: "d", Value: 12.7654763},
		// 				{Type: "i", Value: 1000},
		// 			},
		// 		},
		// 	},
		// },
		{
			name:  "simple address no type string",
			bytes: []byte{47, 104, 101, 108, 108, 111, 0, 0},
			expected: OSCMessage{
				Address: "/hello",
				Args:    []OSCArg{},
			},
		},
		{
			name: "osc 1.0 spec example 1",
			bytes: []byte{
				47, 111, 115, 99, 105, 108, 108, 97, 116, 111, 114, 47, 52, 47, 102, 114, 101, 113, 117, 101, 110, 99, 121, 0, 44,
				102, 0, 0, 67, 220, 0, 0,
			},
			expected: OSCMessage{Address: "/oscillator/4/frequency", Args: []OSCArg{{Type: "f", Value: float32(440)}}},
		},
		{
			name: "osc 1.0 spec example 2",
			bytes: []byte{
				47, 102, 111, 111, 0, 0, 0, 0, 44, 105, 105, 115, 102, 102, 0, 0, 0, 0, 3, 232, 255, 255, 255, 255, 104, 101, 108,
				108, 111, 0, 0, 0, 63, 157, 243, 182, 64, 181, 178, 45,
			},
			expected: OSCMessage{
				Address: "/foo",
				Args: []OSCArg{
					{Type: "i", Value: int32(1000)},
					{Type: "i", Value: int32(-1)},
					{Type: "s", Value: "hello"},
					// thanks IEEE 754
					{Type: "f", Value: float32(1.2339999675750732421875)},
					{Type: "f", Value: float32(5.677999973297119140625)},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			actual, err := MessageFromBytes(testCase.bytes)

			if err != nil {
				t.Fatalf("failed to encode properly: %s", err.Error())
			}

			if !reflect.DeepEqual(actual.Address, testCase.expected.Address) {
				t.Fatalf("failed to encode address propertly got '%s', expected '%s'", actual.Address, testCase.expected.Address)
			}

			if !reflect.DeepEqual(actual.Args, testCase.expected.Args) {
				t.Fatalf("failed to encode args properly got '%+v', expected '%+v'", actual.Args, testCase.expected.Args)
			}
		})
	}
}

func TestBadOSCMessageDecoding(t *testing.T) {
	testCases := []struct {
		name        string
		bytes       []byte
		errorString string
	}{
		{
			name:        "empty byte array",
			bytes:       []byte{},
			errorString: "cannot create OSC Message from empty byte array",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := MessageFromBytes(testCase.bytes)

			if err == nil {
				t.Fatalf("MessageFromBytes expected to fail but got: %+v", got)
			}

			if err.Error() != testCase.errorString {
				t.Fatalf("MessageFromBytes got error '%s', expected '%s'", err.Error(), testCase.errorString)
			}
		})
	}
}
