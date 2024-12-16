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

func int32ToOSCBytes(number int32) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, number)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func int64ToOSCBytes(number int64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, number)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func float32ToOSCBytes(number float32) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, number)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func float64ToOSCBytes(number float64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, number)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func byteArrayToOSCBytes(bytes []byte) []byte {
	oscBytes := []byte{}

	bytesSize := len(bytes)
	oscBytes = append(oscBytes, int32ToOSCBytes(int32(bytesSize))...)
	oscBytes = append(oscBytes, bytes...)

	padLength := 4 - (bytesSize % 4)
	if padLength < 4 {
		for i := 0; i < padLength; i++ {
			oscBytes = append(oscBytes, 0)
		}
	}

	return oscBytes
}

func argsToBuffer(args []OSCArg) []byte {
	//TODO(jwetzell): add error handling
	var argBuffers = []byte{}

	for _, arg := range args {
		switch oscType := arg.Type; oscType {
		case "s":
			if value, ok := arg.Value.(string); ok {
				argBuffers = append(argBuffers, stringToOSCBytes(value)...)
			} else {
				fmt.Println("OSC arg had string type but non-string value.")
			}
		case "i":
			if value, ok := arg.Value.(int); ok {
				argBuffers = append(argBuffers, int32ToOSCBytes(int32(value))...)
			} else if value, ok := arg.Value.(int32); ok {
				argBuffers = append(argBuffers, int32ToOSCBytes(value)...)
			} else {
				fmt.Println("OSC arg had integer type but non-integer value.")
			}
		case "f":
			if value, ok := arg.Value.(float32); ok {
				argBuffers = append(argBuffers, float32ToOSCBytes(float32(value))...)
			} else if value, ok := arg.Value.(float64); ok {
				argBuffers = append(argBuffers, float32ToOSCBytes(float32(value))...)
			} else if value, ok := arg.Value.(int); ok {
				argBuffers = append(argBuffers, float32ToOSCBytes(float32(value))...)
			} else if value, ok := arg.Value.(int32); ok {
				argBuffers = append(argBuffers, float32ToOSCBytes(float32(value))...)
			} else if value, ok := arg.Value.(int64); ok {
				argBuffers = append(argBuffers, float32ToOSCBytes(float32(value))...)
			} else {
				fmt.Println("OSC arg had float type but non-float value.")
			}
		case "b":
			if value, ok := arg.Value.([]byte); ok {
				argBuffers = append(argBuffers, byteArrayToOSCBytes(value)...)
			} else {
				fmt.Println("OSC arg had blob type but non-blob value.")
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
			if ok {
				colorBytes := []byte{color.r, color.g, color.b, color.a}
				argBuffers = append(argBuffers, colorBytes...)
			}
		case "h":
			if value, ok := arg.Value.(int); ok {
				argBuffers = append(argBuffers, int64ToOSCBytes(int64(value))...)
			} else if value, ok := arg.Value.(int32); ok {
				argBuffers = append(argBuffers, int64ToOSCBytes(int64(value))...)
			} else if value, ok := arg.Value.(int64); ok {
				argBuffers = append(argBuffers, int64ToOSCBytes(value)...)
			} else {
				fmt.Println("OSC arg had integer type but non-integer value.")
			}
		case "d":
			if value, ok := arg.Value.(float32); ok {
				argBuffers = append(argBuffers, float64ToOSCBytes(float64(value))...)
			} else if value, ok := arg.Value.(float64); ok {
				argBuffers = append(argBuffers, float64ToOSCBytes(float64(value))...)
			} else if value, ok := arg.Value.(int); ok {
				argBuffers = append(argBuffers, float64ToOSCBytes(float64(value))...)
			} else if value, ok := arg.Value.(int32); ok {
				argBuffers = append(argBuffers, float64ToOSCBytes(float64(value))...)
			} else if value, ok := arg.Value.(int64); ok {
				argBuffers = append(argBuffers, float64ToOSCBytes(float64(value))...)
			} else {
				fmt.Println("OSC arg had float type but non-float value.")
			}
		default:
			fmt.Printf("unhandled osc type: %s.\n", oscType)
		}
	}
	return argBuffers
}

func ToBytes(message OSCMessage) []byte {
	//TODO(jwetzell): add error handling
	oscBuffer := []byte{}

	oscBuffer = append(oscBuffer, stringToOSCBytes(message.Address)...)

	var sb strings.Builder

	sb.WriteString(",")

	for _, arg := range message.Args {
		sb.WriteString(arg.Type)
	}

	oscBuffer = append(oscBuffer, stringToOSCBytes(sb.String())...)
	oscBuffer = append(oscBuffer, argsToBuffer(message.Args)...)

	return oscBuffer
}

func readOSCString(bytes []byte) (string, []byte) {
	//TODO(jwetzell): add error handling
	oscString := ""
	stringFinished := false
	stringEndIndex := 0
	remainingBytes := []byte{}

	for index, byteIn := range bytes {
		if !stringFinished {
			if byteIn == 0 {
				oscString = string(bytes[0:index])
				stringEndIndex = index + 1
				break
			}
		}
	}

	stringPadding := 4 - (stringEndIndex % 4)

	if stringPadding < 4 {
		stringEndIndex = stringEndIndex + stringPadding
	}

	remainingBytes = bytes[stringEndIndex:]

	return oscString, remainingBytes
}

func readOSCInt32(bytes []byte) (int32, []byte, error) {
	if len(bytes) < 4 {
		return 0, bytes, errors.New("int data must be at least 4 bytes large")
	}
	bits := binary.BigEndian.Uint32(bytes[0:4])
	return int32(bits), bytes[4:], nil
}

func readOSCInt64(bytes []byte) (int64, []byte, error) {
	if len(bytes) < 8 {
		return 0, bytes, errors.New("int data must be at least 4 bytes large")
	}
	bits := binary.BigEndian.Uint64(bytes[0:8])
	return int64(bits), bytes[8:], nil
}

func readOSCFloat32(bytes []byte) (float32, []byte, error) {
	if len(bytes) < 4 {
		return 0, bytes, errors.New("float data must be at least 4 bytes large")
	}
	bits := binary.BigEndian.Uint32(bytes[0:4])
	return math.Float32frombits(bits), bytes[4:], nil
}

func readOSCFloat64(bytes []byte) (float64, []byte, error) {
	if len(bytes) < 4 {
		return 0, bytes, errors.New("float data must be at least 4 bytes large")
	}
	bits := binary.BigEndian.Uint64(bytes[0:8])
	return math.Float64frombits(bits), bytes[8:], nil
}

func readOSCBlob(bytes []byte) ([]byte, []byte, error) {
	blobLength, remainingBytes, err := readOSCInt32(bytes)

	if err != nil {
		return []byte{}, bytes, errors.New("problem reading blob data size")
	}

	if len(remainingBytes) < int(blobLength) {
		return []byte{}, bytes, errors.New("blob data specified a size larger than the remaining message data")
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
		return OSCColor{0, 0, 0, 0}, bytes, errors.New("color data must be at least 4 bytes large")
	}
	oscColor := OSCColor{
		r: bytes[0],
		g: bytes[1],
		b: bytes[2],
		a: bytes[3],
	}
	return oscColor, bytes[4:], nil
}

func readOSCArg(bytes []byte, oscType string) (OSCArg, []byte, error) {
	var readArgError error

	oscArg := OSCArg{}
	oscArg.Type = oscType

	remainingBytes := []byte{}
	//TODO(jwetzell): add error handling
	switch oscType {
	case "s":
		argString, bytesLeft := readOSCString(bytes)
		oscArg.Value = argString
		remainingBytes = bytesLeft
	case "i":
		argInt, bytesLeft, error := readOSCInt32(bytes)
		if error != nil {
			readArgError = error
		}
		oscArg.Value = argInt
		remainingBytes = bytesLeft
	case "f":
		argFloat, bytesLeft, error := readOSCFloat32(bytes)
		if error != nil {
			readArgError = error
		}
		oscArg.Value = argFloat
		remainingBytes = bytesLeft
	case "b":
		argBytes, bytesLeft, error := readOSCBlob(bytes)
		if error != nil {
			readArgError = error
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
		argColor, bytesLeft, error := readOSCColor(bytes)
		if error != nil {
			readArgError = error
		}
		oscArg.Value = argColor
		remainingBytes = bytesLeft
	case "h":
		argInt, bytesLeft, error := readOSCInt64(bytes)
		if error != nil {
			readArgError = error
		}
		oscArg.Value = argInt
		remainingBytes = bytesLeft
	case "d":
		argFloat, bytesLeft, error := readOSCFloat64(bytes)
		if error != nil {
			readArgError = error
		}
		oscArg.Value = argFloat
		remainingBytes = bytesLeft
	default:
		fmt.Printf("unsupported osc type: %s\n", oscType)
		readArgError = errors.New("unsupported osc type: " + oscType)
	}
	return oscArg, remainingBytes, readArgError
}

func FromBytes(bytes []byte) (OSCMessage, error) {
	//TODO(jwetzell): add Message and Bundle support
	address, typeAndArgBytes := readOSCString(bytes)

	oscMessage := OSCMessage{
		Address: address,
		Args:    []OSCArg{},
	}

	typeString, argBytes := readOSCString(typeAndArgBytes)

	for index, oscType := range typeString {
		if index == 0 {
			if oscType != ',' {
				return OSCMessage{}, errors.New("type string is malformed")
			}
		} else {
			oscArg, remainingBytes, error := readOSCArg(argBytes, string(oscType))
			if error != nil {
				return oscMessage, error
			}
			argBytes = remainingBytes
			oscMessage.Args = append(oscMessage.Args, oscArg)
		}
	}

	return oscMessage, nil
}
