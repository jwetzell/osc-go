package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jwetzell/osc-go"
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

func listenStdin(slip bool, format string) {
	var buffer = make([]byte, 1024)
	message := &osc.Message{}
	s := &SLIP{
		pendingBytes: []byte{},
		Messages:     make(chan []byte),
	}

	if slip {
		go handleOSCSLIP(s, format)
	}
	for {
		n, err := os.Stdin.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			os.Stderr.WriteString(err.Error())
			break
		}
		if n > 0 {
			if slip {
				s.Decode(buffer[:n])
			} else {
				err := message.UnmarshalBinary(buffer[:n])
				if err != nil {
					os.Stderr.WriteString(err.Error())
					break
				}
				handleMessage(message, format)
			}
		}
	}
}
