package osc_test

import (
	"testing"

	"github.com/jwetzell/osc-go"
	"github.com/stretchr/testify/require"
)

func TestGoodOSCBundleEncoding(t *testing.T) {

	testCases := []struct {
		name     string
		bundle   *osc.Bundle
		expected []byte
	}{
		{
			name: "simple contents single message",
			bundle: &osc.Bundle{
				TimeTag: osc.TimeTag{
					Seconds:           32,
					FractionalSeconds: 0,
				},
				Contents: []osc.Packet{&osc.Message{Address: "/oscillator/4/frequency", Args: []osc.Arg{{Type: "f", Value: float32(440)}}}},
			},
			expected: []byte{35, 98, 117, 110, 100, 108, 101, 0, 0, 0, 0,
				32, 0, 0, 0, 0, 0, 0, 0, 32, 47, 111,
				115, 99, 105, 108, 108, 97, 116, 111, 114, 47, 52,
				47, 102, 114, 101, 113, 117, 101, 110, 99, 121, 0,
				44, 102, 0, 0, 67, 220, 0, 0},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			got, err := testCase.bundle.ToBytes()

			require.NoError(t, err)
			require.Equal(t, testCase.expected, got)
		})
	}
}

func TestBadOSCBundleEncoding(t *testing.T) {

	testCases := []struct {
		name        string
		bundle      *osc.Bundle
		errorString string
	}{
		{
			name: "bundle contains message with bad address",
			bundle: &osc.Bundle{
				TimeTag: osc.TimeTag{
					Seconds:           32,
					FractionalSeconds: 0,
				},
				Contents: []osc.Packet{&osc.Message{Address: "hello", Args: []osc.Arg{}}},
			},
			errorString: "OSC Message address must start with /",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			_, err := testCase.bundle.ToBytes()

			require.Error(t, err)
			require.EqualError(t, err, testCase.errorString)
		})
	}
}

func TestGoodOSCBundleDecoding(t *testing.T) {
	testCases := []struct {
		name     string
		expected *osc.Bundle
		bytes    []byte
	}{
		{
			name: "simple contents single message",
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
		{
			name: "simple contents nested bundle",
			expected: &osc.Bundle{
				TimeTag: osc.TimeTag{
					Seconds:           32,
					FractionalSeconds: 0,
				},
				Contents: []osc.Packet{&osc.Bundle{
					TimeTag: osc.TimeTag{
						Seconds:           64,
						FractionalSeconds: 0,
					},
					Contents: []osc.Packet{&osc.Message{Address: "/oscillator/4/frequency", Args: []osc.Arg{{Type: "f", Value: float32(440)}}}},
				}},
			},
			bytes: []byte{35, 98, 117, 110, 100, 108, 101, 0, // #bundle
				0, 0, 0, 32, 0, 0, 0, 0, // time tag
				0, 0, 0, 52, // content size
				35, 98, 117, 110, 100, 108, 101, 0, // #bundle
				0, 0, 0, 64, 0, 0, 0, 0, // time tag
				0, 0, 0, 32, // content size
				47, 111, 115, 99, 105, 108, 108, 97, 116, 111, 114, 47, 52,
				47, 102, 114, 101, 113, 117, 101, 110, 99, 121, 0,
				44, 102, 0, 0, 67, 220, 0, 0},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, remainingBytes, err := osc.BundleFromBytes(testCase.bytes)

			require.NoError(t, err)
			require.Empty(t, remainingBytes)
			require.Equal(t, testCase.expected, actual)
		})
	}
}

func TestBadOSCBundleDecoding(t *testing.T) {
	testCases := []struct {
		name        string
		bytes       []byte
		errorString string
	}{
		{
			name:        "empty byte array",
			bytes:       []byte{},
			errorString: "OSC Bundle has to be at least 20 bytes",
		},
		{
			name:        "does not start with #",
			bytes:       []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			errorString: "OSC Bundle must start with a #",
		},
		{
			name:        "does not start with #bundle",
			bytes:       []byte{35, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			errorString: "OSC Bundle must start with #bundle string",
		},
		{
			name: "bundle header not properly null terminated",
			bytes: []byte{
				35, 98, 117, 110, 100, 108, 101,
				35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35},
			errorString: "OSC string must be null-terminated",
		},
		{
			name: "bundle contains incorrect size",
			bytes: []byte{
				35, 98, 117, 110, 100, 108, 101, 0, // #bundle
				0, 0, 0, 0, 0, 0, 0, 0, // time tag
				0, 0, 0, 100, // content size of 100 but only 10 bytes of content
				35, 35, 35, 35, 35, 35, 35, 35, 35, 35},
			errorString: "bundle doesn't have enough bytes for the content size it specifies",
		},
		{
			name: "bundle doesn't contain message or bundle",
			bytes: []byte{
				35, 98, 117, 110, 100, 108, 101, 0, // #bundle
				0, 0, 0, 0, 0, 0, 0, 0, // time tag
				0, 0, 0, 10,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			errorString: "bundle contents does not look a bundle or message",
		},
		{
			name: "bundle with bad bundle inside",
			bytes: []byte{35, 98, 117, 110, 100, 108, 101, 0, // #bundle
				0, 0, 0, 32, 0, 0, 0, 0, // time tag
				0, 0, 0, 19, // content size
				35, 98, 117, 110, 100, 108, 101, 0, // #bundle
				0, 0, 0, 64, 0, 0, 0, 0, // time tag
				0, 0, 0},
			errorString: "OSC Bundle has to be at least 20 bytes",
		},
		{
			name: "bundle with bad message inside",
			bytes: []byte{35, 98, 117, 110, 100, 108, 101, 0, // #bundle
				0, 0, 0, 32, 0, 0, 0, 0, // time tag
				0, 0, 0, 5, // content size
				47, 104, 101, 108, 108, 111},
			errorString: "OSC string must be null-terminated",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, _, err := osc.BundleFromBytes(testCase.bytes)

			require.Error(t, err)
			require.EqualError(t, err, testCase.errorString)
		})
	}
}

func TestOSCBundleMarshalBinary(t *testing.T) {
	message := &osc.Bundle{
		TimeTag: osc.TimeTag{
			Seconds:           1,
			FractionalSeconds: 2,
		},
		Contents: []osc.Packet{
			&osc.Message{
				Address: "/hello",
			},
		},
	}

	got, err := message.MarshalBinary()
	require.NoError(t, err)

	expected, err := message.ToBytes()
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestGoodOSCBundleUnmarshalBinary(t *testing.T) {
	message := &osc.Bundle{}

	bytes := []byte{35, 98, 117, 110, 100, 108, 101, 0, 0, 0, 0,
		32, 0, 0, 0, 0, 0, 0, 0, 32, 47, 111,
		115, 99, 105, 108, 108, 97, 116, 111, 114, 47, 52,
		47, 102, 114, 101, 113, 117, 101, 110, 99, 121, 0,
		44, 102, 0, 0, 67, 220, 0, 0}

	err := message.UnmarshalBinary(bytes)
	require.NoError(t, err)

	expected, _, err := osc.BundleFromBytes(bytes)
	require.NoError(t, err)

	require.Equal(t, expected, message)
}

func TestBadOSCBundleUnmarshalBinary(t *testing.T) {
	message := &osc.Bundle{}

	bytes := []byte{}

	gotErr := message.UnmarshalBinary(bytes)
	require.Error(t, gotErr)

	_, _, expectedErr := osc.BundleFromBytes(bytes)

	require.EqualError(t, gotErr, expectedErr.Error())
}

func FuzzBundleFromBytes(f *testing.F) {
	seedBytes := [][]byte{
		{},
		{35, 98, 117, 110, 100, 108, 101, 0, 0, 0, 0,
			32, 0, 0, 0, 0, 0, 0, 0, 32, 47, 111,
			115, 99, 105, 108, 108, 97, 116, 111, 114, 47, 52,
			47, 102, 114, 101, 113, 117, 101, 110, 99, 121, 0,
			44, 102, 0, 0, 67, 220, 0, 0},
		{35, 98, 117, 110, 100, 108, 101, 0, // #bundle
			0, 0, 0, 32, 0, 0, 0, 0, // time tag
			0, 0, 0, 52, // content size
			35, 98, 117, 110, 100, 108, 101, 0, // #bundle
			0, 0, 0, 64, 0, 0, 0, 0, // time tag
			0, 0, 0, 32, // content size
			47, 111, 115, 99, 105, 108, 108, 97, 116, 111, 114, 47, 52,
			47, 102, 114, 101, 113, 117, 101, 110, 99, 121, 0,
			44, 102, 0, 0, 67, 220, 0, 0},
	}

	for _, seed := range seedBytes {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		_, _, _ = osc.BundleFromBytes(data)
	})
}
