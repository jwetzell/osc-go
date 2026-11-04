package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/jwetzell/osc-go"
	slip "github.com/jwetzell/slip-go"
	"github.com/urfave/cli/v3"
)

func newEncodeCommand() *cli.Command {
	return &cli.Command{
		Name:  "encode",
		Usage: "write encoded OSC message to stdout",
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
				Usage: "whether to slip encode the OSC message bytes",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			address := cmd.String("address")
			args := cmd.StringSlice("arg")
			types := []rune{}
			for _, t := range cmd.StringSlice("type") {
				if len(t) != 1 {
					return fmt.Errorf("OSC type should be a single character, got: %s", t)
				}
				types = append(types, rune(t[0]))
			}
			slip := cmd.Bool("slip")
			return encodeMsg(address, args, types, slip)
		},
	}
}

func encodeMsg(address string, args []string, types []rune, useSLIP bool) error {

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

	var outWriter io.Writer = os.Stdout
	if useSLIP {
		outWriter = slip.NewWriter(os.Stdout)

	}
	_, err = outWriter.Write(oscMessageBuffer)
	if err != nil {
		return fmt.Errorf("Error writing OSC message to stdout: %w", err)
	}
	return nil
}
