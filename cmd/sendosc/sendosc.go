package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"

	"github.com/hypebeast/go-osc/osc"

	"github.com/spf13/cobra"
)

func main() {
	var Host string
	var Port int32
	var Address string
	var Protocol string
	var Args []string
	var Types []string
	var Slip bool

	var rootCmd = &cobra.Command{
		Use: "sendosc",
		Run: func(cmd *cobra.Command, args []string) {
			send(Host, Port, Address, Args, Types, Protocol, Slip)
		},
	}
	rootCmd.Flags().StringVar(&Host, "host", "", "host to send OSC message to")
	rootCmd.Flags().Int32Var(&Port, "port", 9999, "port to send OSC message to")
	rootCmd.Flags().StringVar(&Address, "address", "", "OSC address")
	rootCmd.Flags().StringVar(&Protocol, "protocol", "udp", "protocol to use to send (tcp or udp)")
	rootCmd.Flags().StringArrayVar(&Args, "arg", []string{}, "OSC args")
	rootCmd.Flags().StringArrayVar(&Types, "type", []string{}, "OSC types")
	rootCmd.Flags().BoolVar(&Slip, "slip", false, "whether to slip encode the OSC Message bytes")
	rootCmd.MarkFlagRequired("host")
	rootCmd.MarkFlagRequired("port")
	rootCmd.MarkFlagRequired("address")
	rootCmd.Execute()
}

func argToTypedArg(rawArg string, oscType string) any {
	switch oscType {
	case "s":
		return rawArg
	case "i":
		number, err := strconv.ParseInt(rawArg, 10, 32)
		if err != nil {
			// ... handle error
			panic(err)
		}
		return int32(number)
	case "f":
		number, err := strconv.ParseFloat(rawArg, 32)
		if err != nil {
			// ... handle error
			panic(err)
		}
		return float32(number)
	case "b":
		data, err := hex.DecodeString(rawArg)
		if err != nil {
			// ... handle error
			panic(err)
		}
		return data
	default:
		fmt.Print("unhandled osc type: ")
		fmt.Printf("%s.\n", oscType)
		return rawArg
	}
}

func slipEncode(bytes []byte) []byte {
	END := byte(0xc0)
	ESC := byte(0xdb)
	ESC_END := byte(0xdc)
	ESC_ESC := byte(0xdd)

	var encodedBytes = []byte{}

	for _, byteToEncode := range bytes {
		if byteToEncode == END {
			encodedBytes = append(encodedBytes, ESC_END)
		} else if byteToEncode == ESC {
			encodedBytes = append(encodedBytes, ESC_ESC)
		} else {
			encodedBytes = append(encodedBytes, byteToEncode)
		}
	}

	encodedBytes = append(encodedBytes, END)
	return encodedBytes
}

func send(host string, port int32, address string, args []string, types []string, protocol string, slip bool) {

	oscMessage := osc.NewMessage(address)

	for index, arg := range args {
		oscType := "s"
		if len(types) > index {
			oscType = types[index]
		}

		oscMessage.Append(argToTypedArg(arg, oscType))
	}

	oscMessageBuffer, err := oscMessage.MarshalBinary()

	if err != nil {
		panic(err)
	}

	if slip {
		oscMessageBuffer = slipEncode(oscMessageBuffer)
	}

	netAddress := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial(protocol, netAddress)
	if err != nil {
		fmt.Printf("Dial err %v", err)
		panic(err)
	}
	defer conn.Close()

	if _, err = conn.Write([]byte(oscMessageBuffer)); err != nil {
		fmt.Printf("Write err %v", err)
		panic(err)
	}
}
