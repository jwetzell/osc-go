package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/jwetzell/osc-go"
	"github.com/urfave/cli/v3"
)

func newReceiveCommand() *cli.Command {
	return &cli.Command{
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
	}
}

func handleTCPConnection(conn net.Conn, useSLIP bool, format string) {
	slip := &SLIP{
		pendingBytes: []byte{},
		Messages:     make(chan []byte),
	}
	go handleOSCSLIP(slip, format)

	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		bytesRead, err := conn.Read(buffer)

		if err != nil {
			return
		}
		if useSLIP {
			slip.Decode(buffer[0:bytesRead])
		} else {
			// TODO(jwetzell): handle non-SLIP TCP messages properly
		}

	}
}

func handlePacket(message osc.Packet, format string) {
	if bundle, ok := message.(*osc.Bundle); ok {
		handleBundle(bundle, format)
	} else if msg, ok := message.(*osc.Message); ok {
		handleMessage(msg, format)
	}
}

func handleMessage(message *osc.Message, format string) {
	if format == "json" {
		jsonData, _ := json.Marshal(message)
		os.Stdout.WriteString(string(jsonData))
	} else {
		os.Stdout.WriteString(fmt.Sprintf("%v\n", message))
	}
}

func handleBundle(bundle *osc.Bundle, format string) {
	for _, packet := range bundle.Contents {
		handlePacket(packet, format)
	}
}

func listenUDP(netAddress string, format string) {

	laddr, err := net.ResolveUDPAddr("udp4", netAddress)
	if err != nil {
		slog.Default().Error("Error resolving UDP address", "err", err)
		return
	}

	conn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		slog.Default().Error("Error listening on UDP", "err", err)
		return
	}

	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		bytesRead, _, err := conn.ReadFromUDP(buffer)

		if err != nil {
			slog.Default().Error("Error reading from UDP", "err", err)
			return
		}

		oscPacket, _, err := osc.PacketFromBytes(buffer[0:bytesRead])

		if err != nil {
			slog.Default().Error("Error parsing OSC packet", "err", err)
			return
		}
		handlePacket(oscPacket, format)
	}
}

func listenTCP(netAddress string, useSLIP bool, format string) {
	socket, err := net.Listen("tcp4", netAddress)
	if err != nil {
		slog.Default().Error("Error listening on TCP", "err", err)
		return
	}

	defer socket.Close()

	for {
		conn, err := socket.Accept()
		if err != nil {
			slog.Default().Error("Error accepting TCP connection", "err", err)
			continue
		}
		go handleTCPConnection(conn, useSLIP, format)
	}
}
