package osc

type Packet interface {
	ToBytes() ([]byte, error)
}

type Bundle struct {
	Contents []Packet `json:"contents"`
	TimeTag  TimeTag  `json:"timeTag"`
}

type Arg struct {
	Value any    `json:"value"`
	Type  string `json:"type"`
}

type Message struct {
	Address string `json:"address"`
	Args    []Arg  `json:"args"`
}

type Color struct {
	r uint8
	g uint8
	b uint8
	a uint8
}

type TimeTag struct {
	seconds           int32
	fractionalSeconds int32
}
