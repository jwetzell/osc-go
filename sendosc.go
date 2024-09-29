package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type OSCArg struct {
	Type  string
	Value any
}

type OSCMessage struct {
	Address string
	Args    []OSCArg
}

func stringToOSCBytes(rawString string) []byte {
	var sb strings.Builder

	sb.WriteString(rawString)
	sb.WriteString("\u0000")

	padLength := 4 - (len(sb.String()) % 4)
	if padLength < 4 {
		for i := 0; i < padLength; i++ {
			sb.WriteString("\u0000")
		}
	}

	return []byte(sb.String())
}

func integerToOSCBytes(number int32) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, number)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func floatToOSCBytes(number float32) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, number)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func byteArrayToOSCBytes(bytes []byte) []byte {
	oscBytes := []byte{}

	bytesSize := len(bytes)
	oscBytes = append(oscBytes, integerToOSCBytes(int32(bytesSize))...)
	oscBytes = append(oscBytes, bytes...)

	padLength := 4 - (bytesSize % 4)
	if padLength < 4 {
		for i := 0; i < padLength; i++ {
			oscBytes = append(oscBytes, 0)
		}
	}

	return oscBytes
}

func argsToBuffer(args []OSCArg) []byte {
	var argBuffers = []byte{}

	for _, arg := range args {
		switch oscType := arg.Type; oscType {
		case "s":
			if value, ok := arg.Value.(string); ok {
				argBuffers = append(argBuffers, stringToOSCBytes(value)...)
			} else {
				fmt.Println("OSC arg had string type but non-string value.")
			}
		case "i":
			if value, ok := arg.Value.(int32); ok {
				argBuffers = append(argBuffers, integerToOSCBytes(value)...)
			} else {
				fmt.Println("OSC arg had integer type but non-integer value.")
			}
		case "f":
			if value, ok := arg.Value.(float32); ok {
				argBuffers = append(argBuffers, floatToOSCBytes(value)...)
			} else {
				fmt.Println("OSC arg had float type but non-float value.")
			}
		case "b":
			if value, ok := arg.Value.([]byte); ok {
				argBuffers = append(argBuffers, byteArrayToOSCBytes(value)...)
			} else {
				fmt.Println("OSC arg had blob type but non-blob value.")
			}
		default:
			fmt.Print("unhandled osc type: ")
			fmt.Printf("%s.\n", oscType)
		}
	}
	return argBuffers
}

func messageToBuffer(message OSCMessage) []byte {
	oscBuffer := []byte{}

	oscBuffer = append(oscBuffer, stringToOSCBytes(message.Address)...)

	var sb strings.Builder

	sb.WriteString(",")

	for _, arg := range message.Args {
		sb.WriteString(arg.Type)
	}

	oscBuffer = append(oscBuffer, stringToOSCBytes(sb.String())...)
	oscBuffer = append(oscBuffer, argsToBuffer(message.Args)...)

	return oscBuffer
}

func main() {
	var Host string
	var Port int32
	var Address string
	var Args []string
	var Types []string

	var rootCmd = &cobra.Command{
		Use: "sendosc",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Host)
			fmt.Println(Port)
			fmt.Println(Address)
			fmt.Println(Args)
			send(Host, Port, Address, Args, Types)
		},
	}
	rootCmd.Flags().StringVar(&Host, "host", "", "host to send OSC message to")
	rootCmd.Flags().Int32Var(&Port, "port", 9999, "port to send OSC message to")
	rootCmd.Flags().StringVar(&Address, "address", "", "OSC address")
	rootCmd.Flags().StringArrayVar(&Args, "arg", []string{}, "OSC args")
	rootCmd.Flags().StringArrayVar(&Types, "type", []string{}, "OSC types")
	rootCmd.MarkFlagRequired("host")
	rootCmd.MarkFlagRequired("port")
	rootCmd.MarkFlagRequired("address")
	rootCmd.Execute()
}

func argToTypedArg(rawArg string, oscType string) OSCArg {
	switch oscType {
	case "s":
		return OSCArg{
			Type:  "s",
			Value: rawArg,
		}
	case "i":
		number, err := strconv.ParseInt(rawArg, 10, 32)
		if err != nil {
			// ... handle error
			panic(err)
		}
		return OSCArg{
			Type:  "i",
			Value: int32(number),
		}
	case "f":
		number, err := strconv.ParseFloat(rawArg, 32)
		if err != nil {
			// ... handle error
			panic(err)
		}
		return OSCArg{
			Type:  "f",
			Value: float32(number),
		}
	case "b":
		data, err := hex.DecodeString(rawArg)
		if err != nil {
			// ... handle error
			panic(err)
		}
		return OSCArg{
			Type:  "b",
			Value: data,
		}
	default:
		fmt.Print("unhandled osc type: ")
		fmt.Printf("%s.\n", oscType)
		return OSCArg{}
	}
}

func send(host string, port int32, address string, args []string, types []string) {

	oscArgs := []OSCArg{}

	for index, arg := range args {
		oscType := "s"
		if len(types) > index {
			oscType = types[index]
		}

		oscArgs = append(oscArgs, argToTypedArg(arg, oscType))
	}
	fmt.Println(oscArgs)

	oscMessage := OSCMessage{
		Address: address,
		Args:    oscArgs,
	}

	oscMessageBuffer := messageToBuffer(oscMessage)

	netAddress := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("udp", netAddress)
	if err != nil {
		fmt.Printf("Dial err %v", err)
		os.Exit(-1)
	}
	defer conn.Close()

	if _, err = conn.Write([]byte(oscMessageBuffer)); err != nil {
		fmt.Printf("Write err %v", err)
		os.Exit(-1)
	}
}
