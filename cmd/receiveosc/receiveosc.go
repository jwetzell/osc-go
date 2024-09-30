package main

import (
	"fmt"
	"net"

	"github.com/chabad360/go-osc/osc"
	"github.com/spf13/cobra"
)

func main() {
	var Host string
	var Port string

	var rootCmd = &cobra.Command{
		Use: "sendosc",
		Run: func(cmd *cobra.Command, args []string) {
			netAddress := Host + ":" + Port
			listen(netAddress)
		},
	}
	rootCmd.Flags().StringVar(&Host, "host", "127.0.0.1", "host to send OSC message to")
	rootCmd.Flags().StringVar(&Port, "port", "8888", "port to send OSC message to")
	rootCmd.Execute()
}

func listen(netAddress string) {

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

		oscMessage, err := osc.NewMessageFromData(buffer[0:bytesRead])

		if err != nil {
			panic(err)
		}
		fmt.Println(oscMessage)
	}
}
