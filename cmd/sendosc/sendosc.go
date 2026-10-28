package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/jwetzell/osc-go"
)

type stringSliceFlag []string

func (i *stringSliceFlag) String() string {
	return strings.Join(*i, ", ")
}

func (i *stringSliceFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var hostFlag string
var portFlag int
var protocolFlag string
var addressFlag string
var slipFlag bool
var argsFlag stringSliceFlag
var typesFlag stringSliceFlag

func init() {
	flag.StringVar(&hostFlag, "host", "", "host to send OSC message to")
	flag.IntVar(&portFlag, "port", 0, "port to send OSC message to")
	flag.StringVar(&protocolFlag, "protocol", "udp", "protocol to use to send (tcp or udp)")
	flag.StringVar(&addressFlag, "address", "", "OSC address")
	flag.BoolVar(&slipFlag, "slip", false, "whether to slip encode the OSC Message bytes")
	flag.Var(&argsFlag, "arg", "OSC args")
	flag.Var(&typesFlag, "type", "OSC types")
}

func main() {
	flag.Parse()

	if hostFlag == "" {
		slog.Error("host is required")
		return
	}

	if portFlag == 0 {
		slog.Error("port is required")
		return
	}

	if addressFlag == "" {
		slog.Error("address is required")
		return
	}

	switch protocolFlag {
	case "":
		protocolFlag = "udp"
	case "udp", "tcp":
		// do nothing
	default:
		slog.Error("invalid protocol", "protocol", protocolFlag)
		return
	}

	types := []rune{}
	for _, t := range typesFlag {
		if len(t) != 1 {
			slog.Error("invalid OSC type", "type", t)
			return
		}
		types = append(types, rune(t[0]))
	}

	err := send(hostFlag, portFlag, addressFlag, argsFlag, types, protocolFlag, slipFlag)
	if err != nil {
		slog.Error("error sending OSC message", "err", err)
		return
	}
}

func slipEncode(bytes []byte) []byte {
	END := byte(0xc0)
	ESC := byte(0xdb)
	ESC_END := byte(0xdc)
	ESC_ESC := byte(0xdd)

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

func send(host string, port int, address string, args []string, types []rune, protocol string, slip bool) error {

	oscMessage := osc.Message{
		Address: address,
		Args:    []osc.Arg{},
	}

	for index, arg := range args {
		oscType := 's'
		if len(types) > index {
			oscType = types[index]
		}
		typedArg, err := osc.ArgFromStringAndType(arg, oscType)
		if err != nil {

			return fmt.Errorf("Error converting arg to typed arg: %w", err)
		}
		oscMessage.Args = append(oscMessage.Args, typedArg)

	}

	oscMessageBuffer, err := oscMessage.ToBytes()
	if err != nil {
		return fmt.Errorf("Error converting osc message to bytes: %w", err)
	}

	if slip {
		oscMessageBuffer = slipEncode(oscMessageBuffer)
	} else if protocol == "tcp" {
		// OSC 1.0 prepends a 4 byte size header for non-SLIP TCP messages
		size := uint32(len(oscMessageBuffer))
		sizeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(sizeBytes, size)
		oscMessageBuffer = append(sizeBytes, oscMessageBuffer...)
	}

	netAddress := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.Dial(protocol, netAddress)
	if err != nil {
		return fmt.Errorf("Dial err: %w", err)
	}
	defer conn.Close()

	if _, err = conn.Write([]byte(oscMessageBuffer)); err != nil {
		return fmt.Errorf("Write err: %w", err)
	}
	return nil
}
