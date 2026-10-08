package osc

// TODO(jwetzell): split things up
import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"slices"
)

func stringToOSCBytes(rawString string) []byte {
	rawStringLength := len(rawString)
	amountOfPadding := (4 - (rawStringLength+1)%4) % 4 //strings need to be 4-byte aligned

	bytes := make([]byte, rawStringLength+1+amountOfPadding) // string + null + padding
	copy(bytes, rawString)
	return bytes
}

func int32ToOSCBytes(number int32) []byte {
	bytes := make([]byte, 4)
	bytes[0] = byte((number >> 24) & 0xFF)
	bytes[1] = byte((number >> 16) & 0xFF)
	bytes[2] = byte((number >> 8) & 0xFF)
	bytes[3] = byte(number & 0xFF)
	return bytes
}

func int64ToOSCBytes(number int64) []byte {
	bytes := make([]byte, 8)
	bytes[0] = byte((number >> 56) & 0xFF)
	bytes[1] = byte((number >> 48) & 0xFF)
	bytes[2] = byte((number >> 40) & 0xFF)
	bytes[3] = byte((number >> 32) & 0xFF)
	bytes[4] = byte((number >> 24) & 0xFF)
	bytes[5] = byte((number >> 16) & 0xFF)
	bytes[6] = byte((number >> 8) & 0xFF)
	bytes[7] = byte(number & 0xFF)
	return bytes
}

func float32ToOSCBytes(number float32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, math.Float32bits(number))
	return bytes
}

func float64ToOSCBytes(number float64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, math.Float64bits(number))
	return bytes
}

func byteArrayToOSCBytes(bytes []byte) []byte {

	bytesSize := len(bytes)
	bytesSizeBytes := int32ToOSCBytes(int32(bytesSize))
	oscBytes := []byte{bytesSizeBytes[0], bytesSizeBytes[1], bytesSizeBytes[2], bytesSizeBytes[3]}
	oscBytes = append(oscBytes, bytes...)

	padLength := 4 - (bytesSize % 4)
	if padLength < 4 {
		for range padLength {
			oscBytes = append(oscBytes, 0)
		}
	}

	return oscBytes
}

func timeTagToOSCBytes(timeTag TimeTag) []byte {
	secondsBytes := int32ToOSCBytes(timeTag.Seconds)
	fractionalSecondsBytes := int32ToOSCBytes(timeTag.FractionalSeconds)
	timeTagBytes := []byte{
		secondsBytes[0],
		secondsBytes[1],
		secondsBytes[2],
		secondsBytes[3],
		fractionalSecondsBytes[0],
		fractionalSecondsBytes[1],
		fractionalSecondsBytes[2],
		fractionalSecondsBytes[3],
	}
	return timeTagBytes
}

func argsToBuffer(args []Arg) ([]byte, error) {
	var argBuffers = []byte{}

	for _, arg := range args {
		switch arg.Type {
		case "s":
			if value, ok := arg.Value.(string); ok {
				argBuffers = append(argBuffers, stringToOSCBytes(value)...)
			} else {
				return nil, errors.New("OSC arg had string type but non-string value")
			}
		case "i":
			if value, ok := arg.Value.(int); ok {
				valueBytes := int32ToOSCBytes(int32(value))
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(int32); ok {
				valueBytes := int32ToOSCBytes(int32(value))
				argBuffers = append(argBuffers, valueBytes...)
			} else {
				return nil, errors.New("OSC arg had int32 type but non-number value")
			}
		case "f":
			if value, ok := arg.Value.(float32); ok {
				valueBytes := float32ToOSCBytes(value)
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(float64); ok {
				valueBytes := float32ToOSCBytes(float32(value))
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(int); ok {
				valueBytes := float32ToOSCBytes(float32(value))
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(int32); ok {
				valueBytes := float32ToOSCBytes(float32(value))
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(int64); ok {
				valueBytes := float32ToOSCBytes(float32(value))
				argBuffers = append(argBuffers, valueBytes...)
			} else {
				return nil, errors.New("OSC arg had float32 type but non-number value")
			}
		case "b":
			if value, ok := arg.Value.([]byte); ok {
				valueBytes := byteArrayToOSCBytes(value)
				argBuffers = append(argBuffers, valueBytes...)
			} else {
				return nil, errors.New("OSC arg had blob type but non-blob value")
			}
		case "T":
			continue
		case "F":
			continue
		case "N":
			continue
		case "I":
			continue
		case "r":
			color, ok := arg.Value.(Color)
			if !ok {
				return nil, errors.New("OSC arg had color type but non-color value")
			}
			if ok {
				colorBytes := []byte{color.R, color.G, color.B, color.A}
				argBuffers = append(argBuffers, colorBytes...)
			}
		case "h":
			if value, ok := arg.Value.(int); ok {
				valueBytes := int64ToOSCBytes(int64(value))
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(int32); ok {
				valueBytes := int64ToOSCBytes(int64(value))
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(int64); ok {
				valueBytes := int64ToOSCBytes(value)
				argBuffers = append(argBuffers, valueBytes...)
			} else {
				return nil, errors.New("OSC arg had int64 type but non-number value")
			}
		case "d":
			if value, ok := arg.Value.(float32); ok {
				valueBytes := float64ToOSCBytes(float64(value))
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(float64); ok {
				valueBytes := float64ToOSCBytes(value)
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(int); ok {
				valueBytes := float64ToOSCBytes(float64(value))
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(int32); ok {
				valueBytes := float64ToOSCBytes(float64(value))
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(int64); ok {
				valueBytes := float64ToOSCBytes(float64(value))
				argBuffers = append(argBuffers, valueBytes...)
			} else {
				return nil, errors.New("OSC arg had float64 type but non-number value")
			}
		default:
			return nil, fmt.Errorf("unsupported OSC argument type: %s", arg.Type)
		}
	}
	return argBuffers, nil
}

func readOSCString(bytes []byte) (string, []byte, error) {

	nullByteIndex := slices.Index(bytes, 0)

	if nullByteIndex == -1 {
		return "", bytes, errors.New("OSC string must be null-terminated")
	}
	paddingEndIndex := nullByteIndex + (4 - (nullByteIndex)%4) - 1
	if paddingEndIndex > len(bytes)-1 {
		return "", bytes, errors.New("OSC string is not properly padded")
	}
	for i := nullByteIndex + 1; i <= paddingEndIndex; i++ {
		if bytes[i] != 0 {
			return "", bytes, errors.New("OSC string padding is not null bytes")
		}
	}
	oscString := string(bytes[:nullByteIndex])
	remainingBytes := bytes[paddingEndIndex+1:]
	return oscString, remainingBytes, nil
}

func readOSCInt32(bytes []byte) (int32, []byte, error) {
	if len(bytes) < 4 {
		return 0, bytes, errors.New("OSC int32 arg is not 4 bytes")
	}
	value := int32(bytes[0])<<24 | int32(bytes[1])<<16 | int32(bytes[2])<<8 | int32(bytes[3])
	return value, bytes[4:], nil
}

func readOSCInt64(bytes []byte) (int64, []byte, error) {
	if len(bytes) < 8 {
		return 0, bytes, errors.New("OSC int64 arg is not 8 bytes")
	}
	value := int64(bytes[0])<<56 | int64(bytes[1])<<48 | int64(bytes[2])<<40 | int64(bytes[3])<<32 | int64(bytes[4])<<24 | int64(bytes[5])<<16 | int64(bytes[6])<<8 | int64(bytes[7])
	return value, bytes[8:], nil
}

func readOSCFloat32(bytes []byte) (float32, []byte, error) {
	if len(bytes) < 4 {
		return 0, bytes, errors.New("OSC float32 arg is not 4 bytes")
	}
	bits := uint32(bytes[0])<<24 | uint32(bytes[1])<<16 | uint32(bytes[2])<<8 | uint32(bytes[3])
	return math.Float32frombits(bits), bytes[4:], nil
}

func readOSCFloat64(bytes []byte) (float64, []byte, error) {
	if len(bytes) < 8 {
		return 0, bytes, errors.New("OSC float64 arg is not 8 bytes")
	}
	bits := uint64(bytes[0])<<56 | uint64(bytes[1])<<48 | uint64(bytes[2])<<40 | uint64(bytes[3])<<32 | uint64(bytes[4])<<24 | uint64(bytes[5])<<16 | uint64(bytes[6])<<8 | uint64(bytes[7])
	return math.Float64frombits(bits), bytes[8:], nil
}

func readOSCBlob(bytes []byte) ([]byte, []byte, error) {
	blobLength, remainingBytes, err := readOSCInt32(bytes)

	if err != nil {
		return []byte{}, bytes, errors.New("OSC blob arg size not valid: " + err.Error())
	}

	if blobLength < 0 {
		return []byte{}, bytes, errors.New("OSC blob arg size not valid: size cannot be negative")
	}

	if len(remainingBytes) < int(blobLength) {
		return []byte{}, bytes, errors.New("OSC blob arg size not valid: size specified is larger than remaining bytes")
	}

	blobStartIndex := 4
	paddingAmount := (4 - (int(blobLength) % 4)) % 4

	paddingEndIndex := blobStartIndex + int(blobLength) + paddingAmount - 1

	if paddingEndIndex > len(bytes)-1 {
		return []byte{}, bytes, errors.New("OSC blob arg size not valid: size specified is larger than remaining bytes when accounting for padding")
	}

	for i := blobStartIndex + int(blobLength); i <= paddingEndIndex; i++ {
		if bytes[i] != 0 {
			return []byte{}, bytes, errors.New("OSC blob is not padded with null bytes")
		}
	}
	return bytes[blobStartIndex : blobStartIndex+int(blobLength)], bytes[paddingEndIndex+1:], nil
}

func readOSCColor(bytes []byte) (Color, []byte, error) {
	if len(bytes) < 4 {
		return Color{0, 0, 0, 0}, bytes, errors.New("OSC color arg is not 4 bytes")
	}
	oscColor := Color{
		R: bytes[0],
		G: bytes[1],
		B: bytes[2],
		A: bytes[3],
	}
	return oscColor, bytes[4:], nil
}

func readOSCTimeTag(bytes []byte) (TimeTag, []byte, error) {
	seconds, bytesAfterSeconds, err := readOSCInt32(bytes)
	if err != nil {
		return TimeTag{}, bytes, fmt.Errorf("OSC time tag seconds are not valid: %s", err)
	}
	fractionalSeconds, remainingBytes, err := readOSCInt32(bytesAfterSeconds)
	if err != nil {
		return TimeTag{}, bytes, fmt.Errorf("OSC time tag fractional seconds are not valid: %s", err)
	}

	return TimeTag{
			Seconds:           seconds,
			FractionalSeconds: fractionalSeconds,
		},
		remainingBytes,
		nil
}

func readOSCArg(bytes []byte, oscType string) (Arg, []byte, error) {
	var readArgError error

	oscArg := Arg{}
	oscArg.Type = oscType

	var remainingBytes []byte
	switch oscType {
	case "s":
		argString, bytesLeft, err := readOSCString(bytes)
		if err != nil {
			return Arg{}, bytes, err
		}
		oscArg.Value = argString
		remainingBytes = bytesLeft
	case "i":
		argInt, bytesLeft, err := readOSCInt32(bytes)
		if err != nil {
			readArgError = err
		}
		oscArg.Value = argInt
		remainingBytes = bytesLeft
	case "f":
		argFloat, bytesLeft, err := readOSCFloat32(bytes)
		if err != nil {
			readArgError = err
		}
		oscArg.Value = argFloat
		remainingBytes = bytesLeft
	case "b":
		argBytes, bytesLeft, err := readOSCBlob(bytes)
		if err != nil {
			readArgError = err
		}
		oscArg.Value = argBytes
		remainingBytes = bytesLeft
	case "T":
		oscArg.Value = true
		remainingBytes = bytes
	case "F":
		oscArg.Value = false
		remainingBytes = bytes
	case "N":
		oscArg.Value = nil
		remainingBytes = bytes
	case "I":
		oscArg.Value = math.MaxInt32
		remainingBytes = bytes
	case "r":
		argColor, bytesLeft, err := readOSCColor(bytes)
		if err != nil {
			readArgError = err
		}
		oscArg.Value = argColor
		remainingBytes = bytesLeft
	case "h":
		argInt, bytesLeft, err := readOSCInt64(bytes)
		if err != nil {
			readArgError = err
		}
		oscArg.Value = argInt
		remainingBytes = bytesLeft
	case "d":
		argFloat, bytesLeft, err := readOSCFloat64(bytes)
		if err != nil {
			readArgError = err
		}
		oscArg.Value = argFloat
		remainingBytes = bytesLeft
	case "t":
		argTimeTag, bytesLeft, err := readOSCTimeTag(bytes)
		if err != nil {
			readArgError = err
		}
		oscArg.Value = argTimeTag
		remainingBytes = bytesLeft
	default:
		return Arg{}, bytes, fmt.Errorf("unsupported OSC argument type: %s", oscType)
	}
	return oscArg, remainingBytes, readArgError
}

func PacketFromBytes(bytes []byte) (Packet, []byte, error) {
	if len(bytes) == 0 {
		return nil, bytes, errors.New("cannot create OSC Packet from empty byte array")
	}

	switch bytes[0] {
	case '#':
		bundle, remainingBytes, err := BundleFromBytes(bytes)
		if err != nil {
			return nil, bytes, err
		}
		return bundle, remainingBytes, nil
	case '/':
		message, err := MessageFromBytes(bytes)
		if err != nil {
			return nil, bytes, err
		}
		return message, []byte{}, nil
	default:
		return nil, bytes, errors.New("OSC Packet must start with # for bundle or / for message")
	}
}
