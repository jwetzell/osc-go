package osc

import (
	"errors"
	"strings"
)

func MessageToBytes(message OSCMessage) []byte {
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

func MessageFromBytes(bytes []byte) (OSCMessage, error) {
	address, typeAndArgBytes := readOSCString(bytes)

	if address[0] != 47 {
		return OSCMessage{}, errors.New("OSC Message address must start with /")
	}

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
