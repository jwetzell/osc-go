package osc

import (
	"errors"
	"strings"
)

func (m *OSCMessage) ToBytes() ([]byte, error) {
	//TODO(jwetzell): add error handling
	oscBuffer := []byte{}

	oscBuffer = append(oscBuffer, stringToOSCBytes(m.Address)...)

	var sb strings.Builder

	sb.WriteString(",")

	for _, arg := range m.Args {
		sb.WriteString(arg.Type)
	}

	oscBuffer = append(oscBuffer, stringToOSCBytes(sb.String())...)
	argsBuffer, err := argsToBuffer(m.Args)
	if err != nil {
		return nil, err
	}
	oscBuffer = append(oscBuffer, argsBuffer...)

	return oscBuffer, nil
}

func MessageFromBytes(bytes []byte) (*OSCMessage, error) {
	if len(bytes) == 0 {
		return nil, errors.New("cannot create OSC Message from empty byte array")
	}
	if bytes[0] != 47 {
		return nil, errors.New("OSC Message must start with /")
	}

	address, typeAndArgBytes, err := readOSCString(bytes)

	if err != nil {
		return nil, err
	}

	oscMessage := OSCMessage{
		Address: address,
		Args:    []OSCArg{},
	}

	if len(typeAndArgBytes) == 0 {
		// NOTE(jwetzell): no type string return early.
		return &oscMessage, nil
	}

	typeString, argBytes, err := readOSCString(typeAndArgBytes)

	if err != nil {
		return nil, err
	}

	for index, oscType := range typeString {
		if index == 0 {
			if oscType != ',' {
				return nil, errors.New("type string is malformed")
			}
		} else {
			oscArg, remainingBytes, err := readOSCArg(argBytes, string(oscType))
			if err != nil {
				return nil, err
			}
			argBytes = remainingBytes
			oscMessage.Args = append(oscMessage.Args, oscArg)
		}
	}

	return &oscMessage, nil
}
