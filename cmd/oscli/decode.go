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

func newDecodeCommand() *cli.Command {
	return &cli.Command{
		Name:    "decode",
		Aliases: []string{"r"},
		Usage:   "receive OSC messages via Stdin",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "format",
				Usage: "format for messages to be output in ('json')",
				Value: "json",
				Validator: func(flag string) error {
					if flag != "json" {
						return fmt.Errorf("format must be 'json'")
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "slip",
				Value: false,
				Usage: "whether to slip encode the OSC Message bytes",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			format := cmd.String("format")
			slip := cmd.Bool("slip")

			listenStdin(slip, format)
			return nil
		},
	}
}

func listenStdin(useSLIP bool, format string) {
	var buffer = make([]byte, 1024)
	message := &osc.Message{}
	var reader io.Reader = os.Stdin
	if useSLIP {
		reader = slip.NewReader(os.Stdin)
	}
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			os.Stderr.WriteString(err.Error())
			break
		}
		if n > 0 {
			err := message.UnmarshalBinary(buffer[:n])
			if err != nil {
				os.Stderr.WriteString(err.Error())
				break
			}
			handleMessage(message, format)
		}
	}
}
