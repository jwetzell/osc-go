package main

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/jwetzell/osc-go"
)

func send(host string, port int32, address string, args []string, types []rune, protocol string, slip bool) error {

	oscMessage := osc.Message{
		Address: address,
		Args:    []osc.Arg{},
	}

	for index, arg := range args {
		oscType := 's'
		if len(types) > index {
			oscType = types[index]
		}
		typedArg, err := osc.ArgFromStringAndType(arg, oscType)
		if err != nil {

			return fmt.Errorf("Error converting arg to typed arg: %w", err)
		}
		oscMessage.Args = append(oscMessage.Args, typedArg)

	}

	oscMessageBuffer, err := oscMessage.ToBytes()
	if err != nil {
		return fmt.Errorf("Error converting osc message to bytes: %w", err)
	}

	if slip {
		s := SLIP{}

		oscMessageBuffer = s.Encode(oscMessageBuffer)
	} else if protocol == "tcp" {
		// OSC 1.0 prepends a 4 byte size header for non-SLIP TCP messages
		size := uint32(len(oscMessageBuffer))
		sizeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(sizeBytes, size)
		oscMessageBuffer = append(sizeBytes, oscMessageBuffer...)
	}

	netAddress := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.Dial(protocol, netAddress)
	if err != nil {
		return fmt.Errorf("Dial err: %w", err)
	}
	defer conn.Close()

	if _, err = conn.Write([]byte(oscMessageBuffer)); err != nil {
		return fmt.Errorf("Write err: %w", err)
	}
	return nil
}
