package osc_test

import (
	"math"
	"testing"

	"github.com/jwetzell/osc-go"
	"github.com/stretchr/testify/require"
)

func TestGoodOSCMessageEncoding(t *testing.T) {

	testCases := []struct {
		name     string
		message  *osc.Message
		expected []byte
	}{
		{
			name: "simple hello",
			message: &osc.Message{
				Address: "/hello",
				Args:    []osc.Arg{},
			},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 0, 0, 0},
		},
		{
			name: "simple address string arg",
			message: &osc.Message{
				Address: "/hello",
				Args: []osc.Arg{
					{
						Type:  's',
						Value: "arg1",
					},
				},
			},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 115, 0, 0, 97, 114, 103, 49, 0, 0, 0, 0},
		},
		{
			name:     "simple address 'i' arg int value",
			message:  &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'i', Value: int(35)}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 105, 0, 0, 0, 0, 0, 35},
		},
		{
			name:     "simple address 'i' arg int32 value",
			message:  &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'i', Value: int32(35)}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 105, 0, 0, 0, 0, 0, 35},
		},
		{
			name:     "simple address 'f' arg float32 value",
			message:  &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'f', Value: float32(34.5)}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 102, 0, 0, 66, 10, 0, 0},
		},
		{
			name:     "simple address blob arg",
			message:  &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'b', Value: []byte{98, 108, 111, 98}}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 98, 0, 0, 0, 0, 0, 4, 98, 108, 111, 98},
		},
		{
			name:     "simple address blob arg with padding",
			message:  &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'b', Value: []byte{98, 108, 111}}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 98, 0, 0, 0, 0, 0, 3, 98, 108, 111, 0},
		},
		{
			name:     "simple address True arg",
			message:  &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'T', Value: true}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 84, 0, 0},
		},
		{
			name:     "simple address False arg",
			message:  &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'F', Value: false}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 70, 0, 0},
		},
		{
			name:     "simple address Infinitum arg",
			message:  &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'I', Value: math.MaxInt32}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 73, 0, 0},
		},
		{
			name:     "simple address color arg",
			message:  &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'r', Value: osc.Color{R: 20, G: 21, B: 22, A: 10}}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 114, 0, 0, 20, 21, 22, 10},
		},
		{
			name:     "simple address nil arg",
			message:  &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'N', Value: nil}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 78, 0, 0},
		},
		{
			name:     "simple address 'h' arg int64 value",
			message:  &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'h', Value: int64(281474976710655)}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 104, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255},
		},
		{
			name:     "simple address 'h' arg int32 value",
			message:  &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'h', Value: int32(10)}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 104, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10},
		},
		{
			name:     "simple address 'h' arg int value",
			message:  &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'h', Value: int(25)}}},
			expected: []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 104, 0, 0, 0, 0, 0, 0, 0, 0, 0, 25},
		},
		{
			name:    "simple address 'd' arg float64 value",
			message: &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'd', Value: float64(12.7654763)}}},
			expected: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 100, 0, 0, 0x40, 0x29, 0x87, 0xec, 0x82, 0x74, 0xb9, 0xe6,
			},
		},
		{
			name:    "simple address 'd' arg float32 value",
			message: &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'd', Value: float32(12.76)}}},
			expected: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 100, 0, 0, 0x40, 0x29, 0x85, 0x1e, 0xc0, 0x00, 0x00, 0x00,
			},
		},
		{
			name:    "simple address 'd' arg int32 value",
			message: &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'd', Value: int32(12)}}},
			expected: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 100, 0, 0, 0x40, 0x28, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			name:    "simple address 'd' arg int value",
			message: &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'd', Value: int(12)}}},
			expected: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 100, 0, 0, 0x40, 0x28, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		// TODO(jwetzell): get array args working
		// {
		// 	name: "simple address array arg",
		// 	message: OSCMessage{
		// 		Address: "/hello",
		// 		Args: []OSCArg{
		// 			[]OSCArg{
		//				{Type: 'd', Value: 12.7654763},
		// 				{Type: 'i', Value: 1000},
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
			message: &osc.Message{Address: "/oscillator/4/frequency", Args: []osc.Arg{{Type: 'f', Value: float32(440)}}},
			expected: []byte{
				47, 111, 115, 99, 105, 108, 108, 97, 116, 111, 114, 47, 52, 47, 102, 114, 101, 113, 117, 101, 110, 99, 121, 0, 44,
				102, 0, 0, 67, 220, 0, 0,
			},
		},
		{
			name: "osc 1.0 spec example 2",
			message: &osc.Message{
				Address: "/foo",
				Args: []osc.Arg{
					{Type: 'i', Value: int32(1000)},
					{Type: 'i', Value: int32(-1)},
					{Type: 's', Value: "hello"},
					// thanks IEEE 754
					{Type: 'f', Value: float32(1.2339999675750732421875)},
					{Type: 'f', Value: float32(5.677999973297119140625)},
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
			got, err := testCase.message.ToBytes()

			require.NoError(t, err)
			require.Equal(t, testCase.expected, got)
		})
	}

}

func TestBadOSCMessageEncoding(t *testing.T) {
	testCases := []struct {
		name        string
		message     *osc.Message
		errorString string
	}{
		{
			name:        "empty message",
			message:     &osc.Message{},
			errorString: "OSC Message must have an address",
		},
		{
			name:        "address does not start with /",
			message:     &osc.Message{Address: "hello"},
			errorString: "OSC Message address must start with /",
		},
		{
			name: "arg with unsupported type",
			message: &osc.Message{
				Address: "/hello",
				Args:    []osc.Arg{{Type: 'x', Value: "unsupported"}},
			},
			errorString: "unsupported OSC argument type: x",
		},
		{
			name: "string arg that is not a string",
			message: &osc.Message{
				Address: "/hello",
				Args:    []osc.Arg{{Type: 's', Value: 123}},
			},
			errorString: "OSC arg had string type but non-string value",
		},
		{
			name: "int32 arg that is not a number",
			message: &osc.Message{
				Address: "/hello",
				Args:    []osc.Arg{{Type: 'i', Value: "not an int"}},
			},
			errorString: "OSC arg had int32 type but non-int32 compatible value",
		},
		{
			name: "float32 arg that is not a number",
			message: &osc.Message{
				Address: "/hello",
				Args:    []osc.Arg{{Type: 'f', Value: "not a float"}},
			},
			errorString: "OSC arg had float32 type but non-float32 compatible value",
		},
		{
			name: "int64 arg that is not a number",
			message: &osc.Message{
				Address: "/hello",
				Args:    []osc.Arg{{Type: 'h', Value: "not an int"}},
			},
			errorString: "OSC arg had int64 type but non-int64 compatible value",
		},
		{
			name: "float64 arg that is not a number",
			message: &osc.Message{
				Address: "/hello",
				Args:    []osc.Arg{{Type: 'd', Value: "not a float"}},
			},
			errorString: "OSC arg had float64 type but non-float64 compatible value",
		},
		{
			name: "blob arg that is not a byte array",
			message: &osc.Message{
				Address: "/hello",
				Args:    []osc.Arg{{Type: 'b', Value: "not a blob"}},
			},
			errorString: "OSC arg had blob type but non-blob value",
		},
		{
			name: "color arg that is not an OSCColor",
			message: &osc.Message{
				Address: "/hello",
				Args:    []osc.Arg{{Type: 'r', Value: "not a color"}},
			},
			errorString: "OSC arg had color type but non-color value",
		},
		{
			name: "empty arg type",
			message: &osc.Message{
				Address: "/hello",
				Args:    []osc.Arg{{Type: 0, Value: "value"}},
			},
			errorString: "OSC Message argument type cannot be empty",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := testCase.message.ToBytes()

			require.Error(t, err)

			require.EqualError(t, err, testCase.errorString)
		})
	}
}

func TestGoodOSCMessageDecoding(t *testing.T) {
	testCases := []struct {
		name     string
		bytes    []byte
		expected *osc.Message
	}{
		{
			name:     "simple address no args",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 0, 0, 0},
			expected: &osc.Message{Address: "/hello", Args: []osc.Arg{}},
		},
		{
			name:     "simple address string arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 115, 0, 0, 97, 114, 103, 49, 0, 0, 0, 0},
			expected: &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 's', Value: "arg1"}}},
		},
		{
			name:     "simple address int32 arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 105, 0, 0, 0, 0, 0, 35},
			expected: &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'i', Value: int32(35)}}},
		},
		{
			name:     "simple address float arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 102, 0, 0, 66, 10, 0, 0},
			expected: &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'f', Value: float32(34.5)}}},
		},
		{
			name:     "simple address blob arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 98, 0, 0, 0, 0, 0, 4, 98, 108, 111, 98},
			expected: &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'b', Value: []byte{98, 108, 111, 98}}}},
		},
		{
			name:     "simple address True arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 84, 0, 0},
			expected: &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'T', Value: true}}},
		},
		{
			name:     "simple address False arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 70, 0, 0},
			expected: &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'F', Value: false}}},
		},
		{
			name:     "simple address color arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 114, 0, 0, 20, 21, 22, 10},
			expected: &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'r', Value: osc.Color{R: 20, G: 21, B: 22, A: 10}}}},
		},
		{
			name:     "simple address nil arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 78, 0, 0},
			expected: &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'N', Value: nil}}},
		},
		{
			name:     "simple address Inifinitum arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 73, 0, 0},
			expected: &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'I', Value: math.MaxInt32}}},
		},
		{
			name:     "simple address int64 arg",
			bytes:    []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 104, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255},
			expected: &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'h', Value: int64(281474976710655)}}},
		},
		{
			name: "simple address float64 arg",
			bytes: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 100, 0, 0, 0x40, 0x29, 0x87, 0xec, 0x82, 0x74, 0xb9, 0xe6,
			},
			expected: &osc.Message{Address: "/hello", Args: []osc.Arg{{Type: 'd', Value: float64(12.7654763)}}},
		},
		// TODO(jwetzell): support OSC array
		// {
		// 	name: "simple address array arg",
		// 	bytes: []byte{
		// 		47, 104, 101, 108, 108, 111, 0, 0, 44, 91, 100, 105, 93, 0, 0, 0, 0x40, 0x29, 0x87, 0xec, 0x82, 0x74, 0xb9, 0xe6,
		// 		0, 0, 3, 232,
		// 	},
		// 	expected: &osc.Message{
		// 		Address: "/hello",
		// 		Args: []osc.Arg{
		// 			[]osc.Arg{
		// 				{Type: 'd', Value: 12.7654763},
		// 				{Type: 'i', Value: 1000},
		// 			},
		// 		},
		// 	},
		// },
		{
			name:  "simple address no type string",
			bytes: []byte{47, 104, 101, 108, 108, 111, 0, 0},
			expected: &osc.Message{
				Address: "/hello",
				Args:    []osc.Arg{},
			},
		},
		{
			name: "osc 1.0 spec example 1",
			bytes: []byte{
				47, 111, 115, 99, 105, 108, 108, 97, 116, 111, 114, 47, 52, 47, 102, 114, 101, 113, 117, 101, 110, 99, 121, 0, 44,
				102, 0, 0, 67, 220, 0, 0,
			},
			expected: &osc.Message{Address: "/oscillator/4/frequency", Args: []osc.Arg{{Type: 'f', Value: float32(440)}}},
		},
		{
			name: "osc 1.0 spec example 2",
			bytes: []byte{
				47, 102, 111, 111, 0, 0, 0, 0, 44, 105, 105, 115, 102, 102, 0, 0, 0, 0, 3, 232, 255, 255, 255, 255, 104, 101, 108,
				108, 111, 0, 0, 0, 63, 157, 243, 182, 64, 181, 178, 45,
			},
			expected: &osc.Message{
				Address: "/foo",
				Args: []osc.Arg{
					{Type: 'i', Value: int32(1000)},
					{Type: 'i', Value: int32(-1)},
					{Type: 's', Value: "hello"},
					// thanks IEEE 754
					{Type: 'f', Value: float32(1.2339999675750732421875)},
					{Type: 'f', Value: float32(5.677999973297119140625)},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			actual, err := osc.MessageFromBytes(testCase.bytes)

			require.NoError(t, err)

			require.Equal(t, testCase.expected, actual)
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
		{
			name:        "does not start with /",
			bytes:       []byte{0, 104, 101, 108, 108, 111, 0, 0},
			errorString: "OSC Message must start with /",
		},
		{
			name:        "address string not padded",
			bytes:       []byte{47, 104, 101, 108, 108, 111, 0},
			errorString: "OSC string is not properly padded",
		},
		{
			name:        "type string not padded",
			bytes:       []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 0},
			errorString: "OSC string is not properly padded",
		},
		{
			name:        "type string does not start with ,",
			bytes:       []byte{47, 104, 101, 108, 108, 111, 0, 0, 45, 0, 0, 0},
			errorString: "type string is malformed",
		},
		{
			name: "string arg not null-terminated",
			bytes: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 115, 0, 0, 97, 114, 103, 49,
			},
			errorString: "OSC string must be null-terminated",
		},
		{
			name: "string arg not padded",
			bytes: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 115, 0, 0, 104, 105, 0,
			},
			errorString: "OSC string is not properly padded",
		},
		{
			name: "int32 arg not 4 bytes",
			bytes: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 105, 0, 0, 0,
			},
			errorString: "OSC int32 arg is not 4 bytes",
		},
		{
			name: "int64 arg not 8 bytes",
			bytes: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 104, 0, 0, 0, 0, 0, 0,
			},
			errorString: "OSC int64 arg is not 8 bytes",
		},
		{
			name: "float32 arg not 4 bytes",
			bytes: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 102, 0, 0, 66,
			},
			errorString: "OSC float32 arg is not 4 bytes",
		},
		{
			name: "float64 arg not 8 bytes",
			bytes: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 100, 0, 0, 0,
			},
			errorString: "OSC float64 arg is not 8 bytes",
		},
		{
			name: "blob arg size not valid",
			bytes: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 98, 0, 0, 0, 0, 0,
			},
			errorString: "OSC blob arg size not valid: OSC int32 arg is not 4 bytes",
		},
		{
			name: "blob arg size mismatch",
			bytes: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 98, 0, 0, 0, 0, 0, 4, 98, 108, 111,
			},
			errorString: "OSC blob arg size not valid: size specified is larger than remaining bytes",
		},
		{
			name: "color arg not 4 bytes",
			bytes: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 114, 0, 0, 20, 21,
			},
			errorString: "OSC color arg is not 4 bytes",
		},
		{
			name: "time tag arg seconds not complete",
			bytes: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 116, 0, 0, 0,
			},
			errorString: "OSC time tag seconds are not valid: OSC int32 arg is not 4 bytes",
		},
		{
			name: "time tag arg fractional seconds not complete",
			bytes: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 116, 0, 0, 0, 32, 0, 0, 0,
			},
			errorString: "OSC time tag fractional seconds are not valid: OSC int32 arg is not 4 bytes",
		},
		{
			name: "unknown arg type",
			bytes: []byte{
				47, 104, 101, 108, 108, 111, 0, 0, 44, 120, 0, 0,
			},
			errorString: "unsupported OSC argument type: x",
		},
		{
			name:        "string not padded with null bytes",
			bytes:       []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 0, 1, 2},
			errorString: "OSC string padding is not null bytes",
		},
		{
			name:        "blob not padded with null",
			bytes:       []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 98, 0, 0, 0, 0, 0, 6, 98, 108, 111, 98, 1, 1, 2, 2},
			errorString: "OSC blob is not padded with null bytes",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := osc.MessageFromBytes(testCase.bytes)

			require.Error(t, err)

			require.EqualError(t, err, testCase.errorString)
		})
	}
}

func TestOSCMessageMarshalBinary(t *testing.T) {
	message := &osc.Message{
		Address: "/hello",
		Args: []osc.Arg{
			{Type: 'i', Value: 35},
		},
	}

	got, err := message.MarshalBinary()
	require.NoError(t, err)

	expected, err := message.ToBytes()
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestGoodOSCMessageUnmarshalBinary(t *testing.T) {
	message := &osc.Message{}

	bytes := []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 105, 0, 0, 0, 0, 0, 35}

	err := message.UnmarshalBinary(bytes)
	require.NoError(t, err)

	expected, err := osc.MessageFromBytes(bytes)
	require.NoError(t, err)

	require.Equal(t, expected, message)
}

func TestBadOSCMessageUnmarshalBinary(t *testing.T) {
	message := &osc.Message{}

	bytes := []byte{}

	gotErr := message.UnmarshalBinary(bytes)
	require.Error(t, gotErr)

	_, expectedErr := osc.MessageFromBytes(bytes)

	require.EqualError(t, gotErr, expectedErr.Error())
}

func BenchmarkMessageToBytes(b *testing.B) {
	message := &osc.Message{
		Address: "/hello",
		Args: []osc.Arg{
			{Type: 'i', Value: 35},
		},
	}

	for b.Loop() {
		_, err := message.ToBytes()
		require.NoError(b, err)
	}
}

func BenchmarkMessageFromBytes(b *testing.B) {
	bytes := []byte{47, 104, 101, 108, 108, 111, 0, 0, 44, 105, 0, 0, 0, 0, 0, 35}

	for b.Loop() {
		_, err := osc.MessageFromBytes(bytes)
		require.NoError(b, err)
	}
}

func FuzzMessageFromBytes(f *testing.F) {
	seedBytes := [][]byte{
		{},
		{47, 104, 101, 108, 108, 111, 0, 0},
		{47, 104, 101, 108, 108, 111, 0, 0, 44, 105, 0, 0, 0, 0, 0, 35},
	}

	for _, seed := range seedBytes {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = osc.MessageFromBytes(data)
	})
}
