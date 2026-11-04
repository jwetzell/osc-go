package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/jwetzell/osc-go"
	slip "github.com/jwetzell/slip-go"
	"github.com/urfave/cli/v3"
)

func newSendCommand() *cli.Command {
	return &cli.Command{
		Name:    "send",
		Aliases: []string{"s"},
		Usage:   "send OSC messages via UDP or TCP",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "host",
				Usage:    "host to send OSC message to",
				Required: true,
			},
			&cli.Int32Flag{
				Name:     "port",
				Usage:    "port to send OSC message to",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "protocol",
				Usage: "protocol to use to send (tcp or udp)",
				Value: "udp",
				Validator: func(flag string) error {
					if flag != "udp" && flag != "tcp" {
						return fmt.Errorf("protocol must be either 'udp' or 'tcp'")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "address",
				Usage:    "OSC address",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:  "arg",
				Usage: "OSC args",
				Value: []string{},
			},
			&cli.StringSliceFlag{
				Name:  "type",
				Usage: "OSC types",
				Value: []string{},
			},
			&cli.BoolFlag{
				Name:  "slip",
				Value: false,
				Usage: "whether to slip encode the OSC Message bytes",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			host := cmd.String("host")
			port := cmd.Int32("port")
			address := cmd.String("address")
			args := cmd.StringSlice("arg")
			types := []rune{}
			for _, t := range cmd.StringSlice("type") {
				if len(t) != 1 {
					return fmt.Errorf("OSC type should be a single character, got: %s", t)
				}
				types = append(types, rune(t[0]))
			}
			protocol := cmd.String("protocol")
			slip := cmd.Bool("slip")
			return send(host, port, address, args, types, protocol, slip)
		},
	}
}

func send(host string, port int32, address string, args []string, types []rune, protocol string, useSLIP bool) error {

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

	outBuffer := bytes.NewBuffer(nil)

	if useSLIP {
		s := slip.NewWriter(outBuffer)

		_, err := s.Write(oscMessageBuffer)
		if err != nil {
			return fmt.Errorf("Error writing SLIP encoded OSC message: %w", err)
		}
	} else {
		if protocol == "tcp" {
			// OSC 1.0 prepends a 4 byte size header for non-SLIP TCP messages
			size := uint32(len(oscMessageBuffer))
			sizeBytes := make([]byte, 4)
			binary.BigEndian.PutUint32(sizeBytes, size)
			oscMessageBuffer = append(sizeBytes, oscMessageBuffer...)
		}
		_, err := outBuffer.Write(oscMessageBuffer)
		if err != nil {
			return fmt.Errorf("Error writing OSC message to buffer: %w", err)
		}
	}

	netAddress := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.Dial(protocol, netAddress)
	if err != nil {
		return fmt.Errorf("Dial err: %w", err)
	}
	defer conn.Close()

	if _, err = conn.Write(outBuffer.Bytes()); err != nil {
		return fmt.Errorf("Write err: %w", err)
	}
	return nil
}
