package osc

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

func StringArg(value string) Arg {
	return Arg{Type: "s", Value: value}
}

func Int32Arg(value int32) Arg {
	return Arg{Type: "i", Value: value}
}

func FloatArg(value float32) Arg {
	return Arg{Type: "f", Value: value}
}

func BlobArg(value []byte) Arg {
	return Arg{Type: "b", Value: value}
}

func TrueArg() Arg {
	return Arg{Type: "T", Value: true}
}

func FalseArg() Arg {
	return Arg{Type: "F", Value: false}
}

func NilArg() Arg {
	return Arg{Type: "N", Value: nil}
}

func ColorArg(red, green, blue, alpha byte) Arg {
	return Arg{Type: "r", Value: Color{R: red, G: green, B: blue, A: alpha}}
}

func Int64Arg(value int64) Arg {
	return Arg{Type: "h", Value: value}
}

func DoubleArg(value float64) Arg {
	return Arg{Type: "d", Value: value}
}

func TimeTagArg(seconds, fractional int32) Arg {
	return Arg{Type: "t", Value: TimeTag{Seconds: seconds, FractionalSeconds: fractional}}
}

func ArgFromStringAndType(rawArg string, oscType string) (Arg, error) {
	switch oscType {
	case "s":
		return Arg{
			Value: rawArg,
			Type:  "s",
		}, nil
	case "i":
		number, err := strconv.ParseInt(rawArg, 10, 32)
		if err != nil {
			// ... handle error
			return Arg{}, err
		}
		return Arg{
			Value: int32(number),
			Type:  "i",
		}, nil
	case "f":
		number, err := strconv.ParseFloat(rawArg, 32)
		if err != nil {
			return Arg{}, err
		}
		return Arg{
			Value: float32(number),
			Type:  "f",
		}, nil
	case "b":
		data, err := hex.DecodeString(rawArg)
		if err != nil {
			return Arg{}, err
		}
		return Arg{
			Value: data,
			Type:  "b",
		}, nil
	case "h":
		number, err := strconv.ParseInt(rawArg, 10, 64)
		if err != nil {
			return Arg{}, err
		}
		return Arg{
			Value: int64(number),
			Type:  "h",
		}, nil
	case "d":
		number, err := strconv.ParseFloat(rawArg, 64)
		if err != nil {
			return Arg{}, err
		}
		return Arg{
			Value: float64(number),
			Type:  "d",
		}, nil
	case "T":
		return Arg{
			Value: true,
			Type:  "T",
		}, nil
	case "F":
		return Arg{
			Value: false,
			Type:  "F",
		}, nil
	case "N":
		return Arg{
			Value: nil,
			Type:  "N",
		}, nil
	default:
		return Arg{}, fmt.Errorf("unsupported OSC arg type: %s", oscType)
	}
}
