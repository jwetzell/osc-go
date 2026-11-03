package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/jwetzell/osc-go"
)

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
