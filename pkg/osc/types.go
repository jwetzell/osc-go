package osc

type OSCPacket interface {
	ToBytes() []byte
}


type OSCArg struct {
	Type  string
	Value any
}

type OSCMessage struct {
	Address string
	Args    []OSCArg
}

type OSCColor struct {
	r uint8
	g uint8
	b uint8
	a uint8
}