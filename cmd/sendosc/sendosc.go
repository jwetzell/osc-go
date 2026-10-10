package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"github.com/jwetzell/osc-go"

	"github.com/urfave/cli/v3"
)

func main() {

	cmd := &cli.Command{
		Name:  "sendosc",
		Usage: "send OSC messages via UDP or TCP",
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
			types := cmd.StringSlice("type")
			protocol := cmd.String("protocol")
			slip := cmd.Bool("slip")
			send(host, port, address, args, types, protocol, slip)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Printf("Error running command: %v\n", err)
	}
}

func slipEncode(bytes []byte) []byte {
	END := byte(0xc0)
	ESC := byte(0xdb)
	ESC_END := byte(0xdc)
	ESC_ESC := byte(0xdd)

	var encodedBytes = []byte{END}

	for _, byteToEncode := range bytes {
		switch byteToEncode {
		case END:
			encodedBytes = append(encodedBytes, ESC, ESC_END)
		case ESC:
			encodedBytes = append(encodedBytes, ESC, ESC_ESC)
		default:
			encodedBytes = append(encodedBytes, byteToEncode)
		}
	}

	encodedBytes = append(encodedBytes, END)
	return encodedBytes
}

func send(host string, port int32, address string, args []string, types []string, protocol string, slip bool) {

	oscMessage := osc.Message{
		Address: address,
		Args:    []osc.Arg{},
	}

	for index, arg := range args {
		oscType := "s"
		if len(types) > index {
			oscType = types[index]
		}
		typedArg, err := osc.ArgFromStringAndType(arg, oscType)
		if err != nil {
			fmt.Printf("Error converting arg to typed arg: %v", err)
			return
		}
		oscMessage.Args = append(oscMessage.Args, typedArg)

	}

	oscMessageBuffer, err := oscMessage.ToBytes()
	if err != nil {
		fmt.Printf("Error converting osc message to bytes %v", err)
		return
	}

	if slip {
		oscMessageBuffer = slipEncode(oscMessageBuffer)
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
		fmt.Printf("Dial err %v\n", err)
		return
	}
	defer conn.Close()

	if _, err = conn.Write([]byte(oscMessageBuffer)); err != nil {
		fmt.Printf("Write err %v\n", err)
		return
	}
}
