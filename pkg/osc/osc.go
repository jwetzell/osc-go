package osc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

type OSCArg struct {
	Type  string
	Value any
}

type OSCMessage struct {
	Address string
	Args    []OSCArg
}

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

func integerToOSCBytes(number int32) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, number)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func floatToOSCBytes(number float32) []byte {
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
	oscBytes = append(oscBytes, integerToOSCBytes(int32(bytesSize))...)
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
			if value, ok := arg.Value.(int32); ok {
				argBuffers = append(argBuffers, integerToOSCBytes(value)...)
			} else {
				fmt.Println("OSC arg had integer type but non-integer value.")
			}
		case "f":
			if value, ok := arg.Value.(float32); ok {
				argBuffers = append(argBuffers, floatToOSCBytes(value)...)
			} else {
				fmt.Println("OSC arg had float type but non-float value.")
			}
		case "b":
			if value, ok := arg.Value.([]byte); ok {
				argBuffers = append(argBuffers, byteArrayToOSCBytes(value)...)
			} else {
				fmt.Println("OSC arg had blob type but non-blob value.")
			}
		default:
			fmt.Print("unhandled osc type: ")
			fmt.Printf("%s.\n", oscType)
		}
	}
	return argBuffers
}

func ToBuffer(message OSCMessage) []byte {
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
