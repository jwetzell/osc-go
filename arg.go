package osc

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
	return Arg{Type: "r", Value: Color{r: red, g: green, b: blue, a: alpha}}
}

func Int64Arg(value int64) Arg {
	return Arg{Type: "h", Value: value}
}

func DoubleArg(value float64) Arg {
	return Arg{Type: "d", Value: value}
}
