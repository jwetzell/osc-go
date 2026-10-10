package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jwetzell/osc-go"
	"github.com/urfave/cli/v3"
)

func main() {

	cmd := &cli.Command{
		Name:  "makeosc",
		Usage: "make osc bytes",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "address",
				Value:    "",
				Usage:    "OSC address",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:  "arg",
				Usage: "OSC args",
			},
			&cli.StringSliceFlag{
				Name:  "type",
				Usage: "OSC types",
			},
			&cli.BoolFlag{
				Name:  "slip",
				Value: false,
				Usage: "whether to slip encode the OSC Message bytes",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			address := cmd.String("address")
			args := cmd.StringSlice("arg")
			types := cmd.StringSlice("type")
			slip := cmd.Bool("slip")
			makeMsg(address, args, types, slip)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Printf("Error running command: %v\n", err)
		return
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

func makeMsg(address string, args []string, types []string, slip bool) {

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
			fmt.Printf("Error converting arg to typed arg: %v\n", err)
			return
		}
		oscMessage.Args = append(oscMessage.Args, typedArg)
	}

	oscMessageBuffer, err := oscMessage.ToBytes()
	if err != nil {
		fmt.Printf("Error converting OSC message to bytes: %v\n", err)
		return
	}

	if slip {
		oscMessageBuffer = slipEncode(oscMessageBuffer)
	}
	os.Stdout.Write(oscMessageBuffer)
}
