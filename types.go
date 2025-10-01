package osc

type OSCPacket interface {
	ToBytes() []byte
}

type OSCBundle struct {
	Contents []OSCPacket `json:"contents"`
	TimeTag  OSCTimeTag  `json:"timeTag"`
}

type OSCArg struct {
	Value any    `json:"value"`
	Type  string `json:"type"`
}

type OSCMessage struct {
	Address string   `json:"address"`
	Args    []OSCArg `json:"args"`
}

type OSCColor struct {
	r uint8
	g uint8
	b uint8
	a uint8
}

type OSCTimeTag struct {
	seconds           int32
	fractionalSeconds int32
}
