package main

import (
	"fmt"
	"net"

	osc "github.com/jwetzell/osc-go"
	"github.com/spf13/cobra"
)

func main() {
	var Host string
	var Port string
	var Protocol string

	var rootCmd = &cobra.Command{
		Use: "sendosc",
		Run: func(cmd *cobra.Command, args []string) {
			netAddress := Host + ":" + Port
			if Protocol == "udp" {
				listenUDP(netAddress)
			} else if Protocol == "tcp" {
				listenTCP(netAddress)
			}
		},
	}
	rootCmd.Flags().StringVar(&Host, "host", "127.0.0.1", "host to send OSC message to")
	rootCmd.Flags().StringVar(&Port, "port", "8888", "port to send OSC message to")
	rootCmd.Flags().StringVar(&Protocol, "protocol", "udp", "protocol to use to send (tcp or udp)")
	rootCmd.Execute()
}

func listenTCP(netAddress string) {
	socket, err := net.Listen("tcp4", netAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer socket.Close()

	fmt.Printf("listening on %s (tcp w/ SLIP)\n", netAddress)

	for {
		conn, err := socket.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

type SLIP struct {
	pendingBytes []byte
	Messages     chan osc.OSCMessage
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
			if len(s.pendingBytes) > 0 {
				message, err := osc.MessageFromBytes(s.pendingBytes)
				if err != nil {
					fmt.Println(err)
				} else {
					s.Messages <- message
				}
			}
			s.pendingBytes = []byte{}
		} else {
			s.pendingBytes = append(s.pendingBytes, packetByte)
		}
	}

}

func handleMessages(slip SLIP) {
	for message := range slip.Messages {
		fmt.Printf("%v\n", message)
	}
}

func handleConnection(conn net.Conn) {
	slip := SLIP{
		pendingBytes: []byte{},
		Messages:     make(chan osc.OSCMessage),
	}
	go handleMessages(slip)

	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		bytesRead, err := conn.Read(buffer)

		if err != nil {
			return
		}

		slip.decode(buffer[0:bytesRead])
	}
}

func listenUDP(netAddress string) {

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

	fmt.Printf("listening on %s (udp)\n", netAddress)

	defer connection.Close()
	buffer := make([]byte, 1024)

	for {
		bytesRead, _, err := connection.ReadFromUDP(buffer)

		if err != nil {
			panic(err)
		}

		oscMessage, err := osc.MessageFromBytes(buffer[0:bytesRead])

		if err != nil {
			panic(err)
		}
		fmt.Println(oscMessage)
	}
}
