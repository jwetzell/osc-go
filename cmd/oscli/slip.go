package main

import (
	"fmt"

	"github.com/jwetzell/osc-go"
)

var (
	END     = byte(0xc0)
	ESC     = byte(0xdb)
	ESC_END = byte(0xdc)
	ESC_ESC = byte(0xdd)
)

type SLIP struct {
	pendingBytes []byte
	Messages     chan []byte
}

func (s *SLIP) Encode(bytes []byte) []byte {

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

func (s *SLIP) Decode(bytes []byte) {
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
				s.Messages <- s.pendingBytes
			}
			s.pendingBytes = []byte{}
		} else {
			s.pendingBytes = append(s.pendingBytes, packetByte)
		}
	}
}

func handleOSCSLIP(slip *SLIP, format string) {
	for message := range slip.Messages {
		oscPacket, _, err := osc.PacketFromBytes(message)
		if err != nil {
			fmt.Printf("Error decoding SLIP packet: %v\n", err)
			return
		}
		handlePacket(oscPacket, format)
	}
}
