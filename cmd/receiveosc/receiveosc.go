package main

import (
	"encoding/json"
	"fmt"
	"net"

	"flag"

	"github.com/jwetzell/osc-go"
)

var ipFlag string
var portFlag int
var protocolFlag string
var formatFlag string
var slipFlag bool

func init() {
	flag.StringVar(&ipFlag, "ip", "0.0.0.0", "ip to receive OSC messages on")
	flag.IntVar(&portFlag, "port", 8888, "port to receive OSC messages on")
	flag.StringVar(&protocolFlag, "protocol", "udp", "protocol to use to receive (tcp or udp)")
	flag.StringVar(&formatFlag, "format", "json", "format for messages to be output in ('json')")
	flag.BoolVar(&slipFlag, "slip", false, "whether to slip encode the OSC Message bytes")
}

func main() {

	ip := ipFlag
	port := portFlag
	protocol := protocolFlag
	format := formatFlag
	slip := slipFlag

	netAddress := fmt.Sprintf("%s:%d", ip, port)
	switch protocol {
	case "udp":
		listenUDP(netAddress, format)
	case "tcp":
		if !slip {
			fmt.Println("OSC 1.0 over TCP is not supported yet")
			return
		}
		listenTCP(netAddress, slip, format)
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
	Packets      chan osc.Packet
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
			switch packetByte {
			case ESC_END:
				s.pendingBytes = append(s.pendingBytes, END)
			case ESC_ESC:
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
					fmt.Printf("Error decoding SLIP packet: %v\n", err)
					return
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
		Packets:      make(chan osc.Packet),
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

func handlePacket(message osc.Packet, format string) {
	if bundle, ok := message.(*osc.Bundle); ok {
		handleBundle(bundle, format)
	} else if msg, ok := message.(*osc.Message); ok {
		handleMessage(msg, format)
	} else {
		fmt.Println("Received unknown OSC Packet type")
	}
}

func handleMessage(message *osc.Message, format string) {
	if format == "json" {
		jsonData, _ := json.Marshal(message)
		fmt.Println(string(jsonData))
	} else {
		fmt.Printf("%v\n", message)
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
		fmt.Println(err)
		return
	}

	conn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		bytesRead, _, err := conn.ReadFromUDP(buffer)

		if err != nil {
			fmt.Printf("Error reading from UDP: %v\n", err)
			return
		}

		oscPacket, _, err := osc.PacketFromBytes(buffer[0:bytesRead])

		if err != nil {
			fmt.Printf("Error parsing OSC packet: %v\n", err)
			return
		}
		handlePacket(oscPacket, format)
	}
}
