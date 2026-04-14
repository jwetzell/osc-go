package osc

// TODO(jwetzell): split things up
import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strings"
)

func stringToOSCBytes(rawString string) []byte {
	var sb strings.Builder

	sb.WriteString(rawString)
	sb.WriteString("\u0000")

	padLength := 4 - (len(sb.String()) % 4)
	if padLength < 4 {
		for i := 0; i < padLength; i++ {
			sb.WriteString("\u0000")
		}
	}

	return []byte(sb.String())
}

func int32ToOSCBytes(number int32) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, number)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func int64ToOSCBytes(number int64) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, number)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func float32ToOSCBytes(number float32) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, number)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func float64ToOSCBytes(number float64) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, number)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func byteArrayToOSCBytes(bytes []byte) ([]byte, error) {
	oscBytes := []byte{}

	bytesSize := len(bytes)
	bytesSizeBytes, err := int32ToOSCBytes(int32(bytesSize))
	if err != nil {
		return nil, err
	}
	oscBytes = append(oscBytes, bytesSizeBytes...)
	oscBytes = append(oscBytes, bytes...)

	padLength := 4 - (bytesSize % 4)
	if padLength < 4 {
		for i := 0; i < padLength; i++ {
			oscBytes = append(oscBytes, 0)
		}
	}

	return oscBytes, nil
}

func timeTagToOSCBytes(timeTag OSCTimeTag) ([]byte, error) {
	timeTagBytes, err := int32ToOSCBytes(timeTag.seconds)
	if err != nil {
		return nil, err
	}
	fractionalSecondsBytes, err := int32ToOSCBytes(timeTag.fractionalSeconds)
	if err != nil {
		return nil, err
	}
	timeTagBytes = append(timeTagBytes, fractionalSecondsBytes...)

	return timeTagBytes, nil
}

func argsToBuffer(args []OSCArg) ([]byte, error) {
	//TODO(jwetzell): add error handling
	var argBuffers = []byte{}

	for _, arg := range args {
		switch oscType := arg.Type; oscType {
		case "s":
			if value, ok := arg.Value.(string); ok {
				argBuffers = append(argBuffers, stringToOSCBytes(value)...)
			} else {
				return nil, errors.New("OSC arg had string type but non-string value")
			}
		case "i":
			if value, ok := arg.Value.(int); ok {
				valueBytes, err := int32ToOSCBytes(int32(value))
				if err != nil {
					return nil, err
				}
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(int32); ok {
				valueBytes, err := int32ToOSCBytes(int32(value))
				if err != nil {
					return nil, err
				}
				argBuffers = append(argBuffers, valueBytes...)
			} else {
				return nil, errors.New("OSC arg had int32 type but non-number value")
			}
		case "f":
			if value, ok := arg.Value.(float32); ok {
				valueBytes, err := float32ToOSCBytes(value)
				if err != nil {
					return nil, err
				}
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(float64); ok {
				valueBytes, err := float32ToOSCBytes(float32(value))
				if err != nil {
					return nil, err
				}
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(int); ok {
				valueBytes, err := float32ToOSCBytes(float32(value))
				if err != nil {
					return nil, err
				}
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(int32); ok {
				valueBytes, err := float32ToOSCBytes(float32(value))
				if err != nil {
					return nil, err
				}
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(int64); ok {
				valueBytes, err := float32ToOSCBytes(float32(value))
				if err != nil {
					return nil, err
				}
				argBuffers = append(argBuffers, valueBytes...)
			} else {
				return nil, errors.New("OSC arg had float32 type but non-number value")
			}
		case "b":
			if value, ok := arg.Value.([]byte); ok {
				valueBytes, err := byteArrayToOSCBytes(value)
				if err != nil {
					return nil, err
				}
				argBuffers = append(argBuffers, valueBytes...)
			} else {
				return nil, errors.New("OSC arg had blob type but non-blob value")
			}
		case "T":
			argBuffers = append(argBuffers, make([]byte, 0)...)
		case "F":
			argBuffers = append(argBuffers, make([]byte, 0)...)
		case "N":
			argBuffers = append(argBuffers, make([]byte, 0)...)
		case "I":
			argBuffers = append(argBuffers, make([]byte, 0)...)
		case "r":
			color, ok := arg.Value.(OSCColor)
			if !ok {
				return nil, errors.New("OSC arg had color type but non-color value")
			}
			if ok {
				colorBytes := []byte{color.r, color.g, color.b, color.a}
				argBuffers = append(argBuffers, colorBytes...)
			}
		case "h":
			if value, ok := arg.Value.(int); ok {
				valueBytes, err := int64ToOSCBytes(int64(value))
				if err != nil {
					return nil, err
				}
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(int32); ok {
				valueBytes, err := int64ToOSCBytes(int64(value))
				if err != nil {
					return nil, err
				}
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(int64); ok {
				valueBytes, err := int64ToOSCBytes(value)
				if err != nil {
					return nil, err
				}
				argBuffers = append(argBuffers, valueBytes...)
			} else {
				return nil, errors.New("OSC arg had int64 type but non-number value")
			}
		case "d":
			if value, ok := arg.Value.(float32); ok {
				valueBytes, err := float64ToOSCBytes(float64(value))
				if err != nil {
					return nil, err
				}
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(float64); ok {
				valueBytes, err := float64ToOSCBytes(value)
				if err != nil {
					return nil, err
				}
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(int); ok {
				valueBytes, err := float64ToOSCBytes(float64(value))
				if err != nil {
					return nil, err
				}
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(int32); ok {
				valueBytes, err := float64ToOSCBytes(float64(value))
				if err != nil {
					return nil, err
				}
				argBuffers = append(argBuffers, valueBytes...)
			} else if value, ok := arg.Value.(int64); ok {
				valueBytes, err := float64ToOSCBytes(float64(value))
				if err != nil {
					return nil, err
				}
				argBuffers = append(argBuffers, valueBytes...)
			} else {
				return nil, errors.New("OSC arg had float64 type but non-number value")
			}
		default:
			return nil, fmt.Errorf("unsupported OSC argument type: %s", oscType)
		}
	}
	return argBuffers, nil
}

func readOSCString(bytes []byte) (string, []byte, error) {
	oscString := ""
	stringEndIndex := 0

	nullByteFound := false
	for index, byteIn := range bytes {
		if byteIn == 0 {
			nullByteFound = true
			oscString = string(bytes[0:index])
			stringEndIndex = index + 1
			break
		}
	}

	if !nullByteFound {
		return "", bytes, errors.New("OSC string must be null-terminated")
	}

	stringPadding := 4 - (stringEndIndex % 4)

	if stringPadding < 4 {
		stringEndIndex = stringEndIndex + stringPadding
	}

	if stringEndIndex > len(bytes) {
		return "", bytes, errors.New("OSC string is not properly padded")
	}

	remainingBytes := bytes[stringEndIndex:]

	return oscString, remainingBytes, nil
}

func readOSCInt32(bytes []byte) (int32, []byte, error) {
	if len(bytes) < 4 {
		return 0, bytes, errors.New("OSC int32 arg is not 4 bytes")
	}
	bits := binary.BigEndian.Uint32(bytes[0:4])
	return int32(bits), bytes[4:], nil
}

func readOSCInt64(bytes []byte) (int64, []byte, error) {
	if len(bytes) < 8 {
		return 0, bytes, errors.New("OSC int64 arg is not 8 bytes")
	}
	bits := binary.BigEndian.Uint64(bytes[0:8])
	return int64(bits), bytes[8:], nil
}

func readOSCFloat32(bytes []byte) (float32, []byte, error) {
	if len(bytes) < 4 {
		return 0, bytes, errors.New("OSC float32 arg is not 4 bytes")
	}
	bits := binary.BigEndian.Uint32(bytes[0:4])
	return math.Float32frombits(bits), bytes[4:], nil
}

func readOSCFloat64(bytes []byte) (float64, []byte, error) {
	if len(bytes) < 4 {
		return 0, bytes, errors.New("OSC float64 arg is not 8 bytes")
	}
	bits := binary.BigEndian.Uint64(bytes[0:8])
	return math.Float64frombits(bits), bytes[8:], nil
}

func readOSCBlob(bytes []byte) ([]byte, []byte, error) {
	blobLength, remainingBytes, err := readOSCInt32(bytes)

	if err != nil {
		return []byte{}, bytes, errors.New("OSC blob arg size not valid: " + err.Error())
	}

	if len(remainingBytes) < int(blobLength) {
		return []byte{}, bytes, errors.New("OSC blob arg size not valid: size specified is larger than remaining bytes")
	}

	blobLengthPadding := 4 - (blobLength % 4)
	blobEnd := 4 + blobLength

	if blobLengthPadding < 4 {
		blobEnd = blobEnd + blobLengthPadding
	}
	return bytes[4 : 4+blobLength], bytes[blobEnd:], nil
}

func readOSCColor(bytes []byte) (OSCColor, []byte, error) {
	if len(bytes) < 4 {
		return OSCColor{0, 0, 0, 0}, bytes, errors.New("OSC color arg is not 4 bytes")
	}
	oscColor := OSCColor{
		r: bytes[0],
		g: bytes[1],
		b: bytes[2],
		a: bytes[3],
	}
	return oscColor, bytes[4:], nil
}

func readOSCTimeTag(bytes []byte) (OSCTimeTag, []byte, error) {
	seconds, bytesAfterSeconds, err := readOSCInt32(bytes)
	if err != nil {
		return OSCTimeTag{}, bytes, fmt.Errorf("OSC time tag seconds are not valid: %s", err)
	}
	fractionalSeconds, remainingBytes, err := readOSCInt32(bytesAfterSeconds)
	if err != nil {
		return OSCTimeTag{}, bytes, fmt.Errorf("OSC time tag fractional seconds are not valid: %s", err)
	}

	return OSCTimeTag{
			seconds:           seconds,
			fractionalSeconds: fractionalSeconds,
		},
		remainingBytes,
		nil
}

func readOSCArg(bytes []byte, oscType string) (OSCArg, []byte, error) {
	var readArgError error

	oscArg := OSCArg{}
	oscArg.Type = oscType

	remainingBytes := []byte{}
	//TODO(jwetzell): add error handling
	switch oscType {
	case "s":
		argString, bytesLeft, err := readOSCString(bytes)
		if err != nil {
			return OSCArg{}, bytes, err
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
		fmt.Printf("unsupported osc type: %s\n", oscType)
		readArgError = errors.New("unsupported OSC argument type: " + oscType)
	}
	return oscArg, remainingBytes, readArgError
}

func PacketFromBytes(bytes []byte) (OSCPacket, []byte, error) {
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
