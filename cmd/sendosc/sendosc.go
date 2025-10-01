package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"

	osc "github.com/jwetzell/osc-go"

	"github.com/urfave/cli/v3"
)

func main() {
	var Host string
	var Port int32
	var Address string
	var Protocol string
	var Args []string
	var Types []string
	var Slip bool

	cmd := &cli.Command{
		Name:  "sendosc",
		Usage: "send OSC messages via UDP or TCP",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Usage:       "host to send OSC message to",
				Destination: &Host,
				Required:    true,
			},
			&cli.Int32Flag{
				Name:        "port",
				Usage:       "port to send OSC message to",
				Destination: &Port,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "protocol",
				Usage:       "protocol to use to send (tcp or udp)",
				Value:       "udp",
				Destination: &Protocol,
				Validator: func(flag string) error {
					if flag != "udp" && flag != "tcp" {
						return fmt.Errorf("protocol must be either 'udp' or 'tcp'")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:        "address",
				Usage:       "OSC address",
				Destination: &Address,
				Required:    true,
			},
			&cli.StringSliceFlag{
				Name:        "arg",
				Usage:       "OSC args",
				Value:       []string{},
				Destination: &Args,
			},
			&cli.StringSliceFlag{
				Name:        "type",
				Usage:       "OSC types",
				Value:       []string{},
				Destination: &Types,
			},
			&cli.BoolFlag{
				Name:        "slip",
				Value:       false,
				Usage:       "whether to slip encode the OSC Message bytes",
				Destination: &Slip,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			send(Host, Port, Address, Args, Types, Protocol, Slip)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}

func argToTypedArg(rawArg string, oscType string) osc.OSCArg {

	switch oscType {
	case "s":
		return osc.OSCArg{
			Value: rawArg,
			Type:  "s",
		}
	case "i":
		number, err := strconv.ParseInt(rawArg, 10, 32)
		if err != nil {
			// ... handle error
			panic(err)
		}
		return osc.OSCArg{
			Value: int32(number),
			Type:  "i",
		}
	case "f":
		number, err := strconv.ParseFloat(rawArg, 32)
		if err != nil {
			// ... handle error
			panic(err)
		}
		return osc.OSCArg{
			Value: float32(number),
			Type:  "f",
		}
	case "b":
		data, err := hex.DecodeString(rawArg)
		if err != nil {
			// ... handle error
			panic(err)
		}
		return osc.OSCArg{
			Value: data,
			Type:  "b",
		}
	case "h":
		number, err := strconv.ParseInt(rawArg, 10, 64)
		if err != nil {
			// ... handle error
			panic(err)
		}
		return osc.OSCArg{
			Value: int64(number),
			Type:  "h",
		}
	case "d":
		number, err := strconv.ParseFloat(rawArg, 64)
		if err != nil {
			// ... handle error
			panic(err)
		}
		return osc.OSCArg{
			Value: float64(number),
			Type:  "d",
		}
	case "T":
		return osc.OSCArg{
			Value: true,
			Type:  "T",
		}
	case "F":
		return osc.OSCArg{
			Value: false,
			Type:  "F",
		}
	case "N":
		return osc.OSCArg{
			Value: nil,
			Type:  "N",
		}
	default:
		fmt.Print("unhandled osc type: ")
		fmt.Printf("%s.\n", oscType)
		// TODO(jwetzell): something better than this like actual nil, err thing
		return osc.OSCArg{}
	}
}

func slipEncode(bytes []byte) []byte {
	END := byte(0xc0)
	ESC := byte(0xdb)
	ESC_END := byte(0xdc)
	ESC_ESC := byte(0xdd)

	var encodedBytes = []byte{END}

	for _, byteToEncode := range bytes {
		if byteToEncode == END {
			encodedBytes = append(encodedBytes, ESC_END)
		} else if byteToEncode == ESC {
			encodedBytes = append(encodedBytes, ESC_ESC)
		} else {
			encodedBytes = append(encodedBytes, byteToEncode)
		}
	}

	encodedBytes = append(encodedBytes, END)
	return encodedBytes
}

func send(host string, port int32, address string, args []string, types []string, protocol string, slip bool) {

	oscMessage := osc.OSCMessage{
		Address: address,
		Args:    []osc.OSCArg{},
	}

	for index, arg := range args {
		oscType := "s"
		if len(types) > index {
			oscType = types[index]
		}
		arg := argToTypedArg(arg, oscType)

		oscMessage.Args = append(oscMessage.Args, arg)

	}

	oscMessageBuffer := oscMessage.ToBytes()

	if slip {
		oscMessageBuffer = slipEncode(oscMessageBuffer)
	}

	netAddress := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial(protocol, netAddress)
	if err != nil {
		fmt.Printf("Dial err %v", err)
		panic(err)
	}
	defer conn.Close()

	if _, err = conn.Write([]byte(oscMessageBuffer)); err != nil {
		fmt.Printf("Write err %v", err)
		panic(err)
	}
}
