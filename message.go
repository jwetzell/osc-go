package osc

import (
	"errors"
)

func (m *Message) ToBytes() ([]byte, error) {

	if len(m.Address) == 0 {
		return nil, errors.New("OSC Message must have an address")
	}

	if m.Address[0] != '/' {
		return nil, errors.New("OSC Message address must start with /")
	}

	typeString := ","

	for _, arg := range m.Args {
		if arg.Type == 0 {
			return nil, errors.New("OSC Message argument type cannot be empty")
		}
		typeString += string(arg.Type)
	}
	addressBuffer := stringToOSCBytes(m.Address)

	oscBuffer := addressBuffer

	if m.Args != nil {
		typeStringBuffer := stringToOSCBytes(typeString)

		argsBuffer, err := argsToBuffer(m.Args)
		if err != nil {
			return nil, err
		}
		oscBuffer = append(oscBuffer, typeStringBuffer...)
		oscBuffer = append(oscBuffer, argsBuffer...)
	}

	return oscBuffer, nil
}

func (m Message) MarshalBinary() ([]byte, error) {
	return m.ToBytes()
}

func (m *Message) UnmarshalBinary(data []byte) error {
	msg, err := MessageFromBytes(data)
	if err != nil {
		return err
	}
	m.Address = msg.Address
	m.Args = msg.Args
	return nil
}

func MessageFromBytes(bytes []byte) (*Message, error) {
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

	oscMessage := Message{
		Address: address,
		Args:    nil,
	}

	if len(typeAndArgBytes) == 0 {
		// NOTE(jwetzell): no type string return early.
		return &oscMessage, nil
	}

	if typeAndArgBytes[0] != ',' {
		return nil, errors.New("type string is malformed")
	}

	oscMessage.Args = []Arg{}

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
			oscArg, remainingBytes, err := readOSCArg(argBytes, oscType)
			if err != nil {
				return nil, err
			}
			argBytes = remainingBytes
			oscMessage.Args = append(oscMessage.Args, oscArg)
		}
	}

	if len(argBytes) > 0 {
		return nil, errors.New("extra bytes found after parsing all arguments")
	}

	return &oscMessage, nil
}
