package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"

	osc "github.com/jwetzell/osc-go"
	"github.com/urfave/cli/v3"
)

func main() {
	var Host string
	var Port int32
	var Protocol string
	var Format string
	var Slip bool

	cmd := &cli.Command{
		Name:  "receiveosc",
		Usage: "receive OSC messages via UDP or TCP",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "ip",
				Usage:       "ip to receive OSC messages on",
				Value:       "0.0.0.0",
				Destination: &Host,
			},
			&cli.Int32Flag{
				Name:        "port",
				Usage:       "port to receive OSC messages on",
				Destination: &Port,
				Value:       8888,
			},
			&cli.StringFlag{
				Name:        "protocol",
				Usage:       "protocol to use to receive (tcp or udp)",
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
				Name:        "format",
				Usage:       "format for messages to be output in ('json')",
				Value:       "json",
				Destination: &Format,
				Validator: func(flag string) error {
					if flag != "json" {
						return fmt.Errorf("format must be 'json'")
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:        "slip",
				Value:       false,
				Usage:       "whether to slip encode the OSC Message bytes",
				Destination: &Slip,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			netAddress := fmt.Sprintf("%s:%d", Host, Port)
			switch Protocol {
			case "udp":
				listenUDP(netAddress, Format)
			case "tcp":
				if !Slip {
					return fmt.Errorf("OSC 1.0 over TCP is not supported yet")
				}
				listenTCP(netAddress, Slip, Format)
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}

func listenTCP(netAddress string, useSLIP bool, format string) {
	socket, err := net.Listen("tcp4", netAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer socket.Close()

	for {
		conn, err := socket.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleTCPConnection(conn, useSLIP, format)
	}
}

type SLIP struct {
	pendingBytes []byte
	Packets      chan osc.OSCPacket
}

func (s *SLIP) decode(bytes []byte) {
	END := byte(0xc0)
	ESC := byte(0xdb)
	ESC_END := byte(0xdc)
	ESC_ESC := byte(0xdd)

	escapeNext := false
	for _, packetByte := range bytes {

		if packetByte == ESC {
			escapeNext = true
			continue
		}

		if escapeNext {
			if packetByte == ESC_END {
				s.pendingBytes = append(s.pendingBytes, END)
			} else if packetByte == ESC_ESC {
				s.pendingBytes = append(s.pendingBytes, ESC)
			}
			escapeNext = false
		} else if packetByte == END {
			if len(s.pendingBytes) == 0 {
				// probably opening END byte, can discard
				continue
			} else {
				oscPacket, _, err := osc.PacketFromBytes(s.pendingBytes)
				if err != nil {
					panic(err)
				} else {
					s.Packets <- oscPacket
				}
			}
			s.pendingBytes = []byte{}
		} else {
			s.pendingBytes = append(s.pendingBytes, packetByte)
		}
	}
}

func handleSLIP(slip SLIP, format string) {
	for message := range slip.Packets {
		handlePacket(message, format)
	}
}

func handleTCPConnection(conn net.Conn, useSLIP bool, format string) {
	slip := SLIP{
		pendingBytes: []byte{},
		Packets:      make(chan osc.OSCPacket),
	}
	go handleSLIP(slip, format)

	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		bytesRead, err := conn.Read(buffer)

		if err != nil {
			return
		}
		if useSLIP {
			slip.decode(buffer[0:bytesRead])
		} else {
			// TODO(jwetzell): handle non-SLIP TCP messages properly
		}

	}
}

func handlePacket(message osc.OSCPacket, format string) {
	if bundle, ok := message.(*osc.OSCBundle); ok {
		handleBundle(bundle, format)
	} else if msg, ok := message.(*osc.OSCMessage); ok {
		handleMessage(msg, format)
	} else {
		fmt.Println("Received unknown OSC Packet type")
	}
}

func handleMessage(message *osc.OSCMessage, format string) {
	if format == "json" {
		jsonData, _ := json.Marshal(message)
		fmt.Println(string(jsonData))
	} else {
		fmt.Printf("%v\n", message)
	}
}

func handleBundle(bundle *osc.OSCBundle, format string) {
	for _, packet := range bundle.Contents {
		handlePacket(packet, format)
	}
}

func listenUDP(netAddress string, format string) {

	s, err := net.ResolveUDPAddr("udp4", netAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()
	buffer := make([]byte, 1024)

	for {
		bytesRead, _, err := connection.ReadFromUDP(buffer)

		if err != nil {
			panic(err)
		}

		fmt.Println("Received UDP packet")
		fmt.Println(buffer[0:bytesRead])
		oscPacket, _, err := osc.PacketFromBytes(buffer[0:bytesRead])

		if err != nil {
			panic(err)
		}
		handlePacket(oscPacket, format)
	}
}
