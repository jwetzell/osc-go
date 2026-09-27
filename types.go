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
	R uint8
	G uint8
	B uint8
	A uint8
}

type TimeTag struct {
	Seconds           int32
	FractionalSeconds int32
}
