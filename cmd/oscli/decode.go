package main

import (
	"os"

	"github.com/jwetzell/osc-go"
)

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
