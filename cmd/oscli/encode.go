package main

import (
	"fmt"
	"os"

	"github.com/jwetzell/osc-go"
)

func encodeMsg(address string, args []string, types []rune, slip bool) error {

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
		return fmt.Errorf("Error converting OSC message to bytes: %w", err)
	}

	if slip {
		s := &SLIP{}
		oscMessageBuffer = s.Encode(oscMessageBuffer)
	}
	_, err = os.Stdout.Write(oscMessageBuffer)
	if err != nil {
		return fmt.Errorf("Error writing OSC message to stdout: %w", err)
	}
	return nil
}
