package osc

import "encoding/binary"

func (b *OSCBundle) ToBytes() []byte {

	bytes := stringToOSCBytes("#bundle")

	bytes = binary.BigEndian.AppendUint64(bytes, b.TimeTag)

	for _, packet := range b.Contents {
		packetBytes := packet.ToBytes()
		packetLength := len(packet.ToBytes())

		bytes = append(bytes, int32ToOSCBytes(int32(packetLength))...)
		bytes = append(bytes, packetBytes...)
	}

	return bytes
}
