package main

import (
	"flag"
	"fmt"
	"os"
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

var addressFlag string
var slipFlag bool
var argsFlag stringSliceFlag
var typesFlag stringSliceFlag

func init() {
	flag.StringVar(&addressFlag, "address", "", "OSC address")
	flag.BoolVar(&slipFlag, "slip", false, "whether to slip encode the OSC Message bytes")
	flag.Var(&argsFlag, "arg", "OSC args")
	flag.Var(&typesFlag, "type", "OSC types")
}

func main() {

	flag.Parse()

	if addressFlag == "" {
		fmt.Println("Address is required")
		flag.PrintDefaults()
		return
	}

	types := []rune{}
	for _, t := range typesFlag {
		if len(t) != 1 {
			fmt.Println("OSC type should be a single character, got: ", t)
			return
		}
		types = append(types, rune(t[0]))
	}

	err := makeMsg(addressFlag, argsFlag, types, slipFlag)
	if err != nil {
		fmt.Printf("Error making OSC message: %v\n", err)
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

func makeMsg(address string, args []string, types []rune, slip bool) error {

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
		return fmt.Errorf("Error converting OSC message to bytes: %w", err)
	}

	if slip {
		oscMessageBuffer = slipEncode(oscMessageBuffer)
	}
	_, err = os.Stdout.Write(oscMessageBuffer)
	if err != nil {
		return fmt.Errorf("Error writing OSC message to stdout: %w", err)
	}
	return nil
}
