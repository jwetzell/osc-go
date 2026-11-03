package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"
)

var version = "local"

func main() {

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{})))
	cmd := &cli.Command{
		Name:    "oscli",
		Usage:   "basic OSC utilities",
		Version: version,
		Commands: []*cli.Command{
			{
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
			},
			{
				Name:    "receive",
				Aliases: []string{"r"},
				Usage:   "receive OSC messages via UDP or TCP",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "ip",
						Usage: "ip to receive OSC messages on",
						Value: "0.0.0.0",
					},
					&cli.Int32Flag{
						Name:  "port",
						Usage: "port to receive OSC messages on",
						Value: 8888,
					},
					&cli.StringFlag{
						Name:  "protocol",
						Usage: "protocol to use to receive (tcp or udp)",
						Value: "udp",
						Validator: func(flag string) error {
							if flag != "udp" && flag != "tcp" {
								return fmt.Errorf("protocol must be either 'udp' or 'tcp'")
							}
							return nil
						},
					},
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
					ip := cmd.String("ip")
					port := cmd.Int32("port")
					protocol := cmd.String("protocol")
					format := cmd.String("format")
					slip := cmd.Bool("slip")

					netAddress := fmt.Sprintf("%s:%d", ip, port)
					switch protocol {
					case "udp":
						listenUDP(netAddress, format)
					case "tcp":
						if !slip {
							return fmt.Errorf("OSC 1.0 over TCP is not supported yet")
						}
						listenTCP(netAddress, slip, format)
					}
					return nil
				},
			},
			{
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
			},
			{
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
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Default().Error("Error running command", "err", err)
	}
}
