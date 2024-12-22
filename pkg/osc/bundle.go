package osc

import (
	"errors"
)

func (b *OSCBundle) ToBytes() []byte {

	bytes := stringToOSCBytes("#bundle")

	bytes = append(bytes, timeTagToOSCBytes(b.TimeTag)...)

	for _, packet := range b.Contents {
		packetBytes := packet.ToBytes()
		packetLength := len(packet.ToBytes())

		bytes = append(bytes, int32ToOSCBytes(int32(packetLength))...)
		bytes = append(bytes, packetBytes...)
	}

	return bytes
}

func BundleFromBytes(bytes []byte) (OSCBundle, []byte, error) {
	if len(bytes) < 20 {
		return OSCBundle{}, bytes, errors.New("bundle has to be at least 20 bytes")
	}

	if bytes[0] != 35 {
		return OSCBundle{}, bytes, errors.New("bundle must start with a #")
	}

	bundleHeader, bytesAfterBundleHeader := readOSCString(bytes)

	if bundleHeader != "#bundle" {
		return OSCBundle{}, bytesAfterBundleHeader, errors.New("bundle must start with #bundle string")
	}

	timeTag, bytesAfterTimeTag, err := readOSCTimeTag(bytesAfterBundleHeader)

	if err != nil {
		return OSCBundle{}, bytesAfterBundleHeader, err
	}

	bundleContents := []OSCPacket{}

	endOfBundle := false

	remainingBytes := bytesAfterTimeTag

	for !endOfBundle {
		contentSize, bytesAfterContentSize, err := readOSCInt32(remainingBytes)

		if err != nil {
			return OSCBundle{}, remainingBytes, err
		}

		remainingBytes = bytesAfterContentSize

		if len(remainingBytes) < int(contentSize) {
			return OSCBundle{}, remainingBytes, errors.New("bundle doesn't have enough bytes for the content size it specifies")
		}

		bundleContentBytes := remainingBytes[0:contentSize]

		if bundleContentBytes[0] == 35 {
			content, _, err := BundleFromBytes(bundleContentBytes)
			if err != nil {
				return OSCBundle{}, remainingBytes, err
			}
			bundleContents = append(bundleContents, &content)
		} else if bundleContentBytes[0] == 47 {
			content, err := MessageFromBytes(bundleContentBytes)
			if err != nil {
				return OSCBundle{}, remainingBytes, err
			}
			bundleContents = append(bundleContents, &content)
		} else {
			return OSCBundle{}, remainingBytes, errors.New("bundle contents does not look a bundle or message")
		}
		remainingBytes = bytesAfterContentSize[contentSize:]
		if len(remainingBytes) == 0 {
			endOfBundle = true
		}

	}

	return OSCBundle{
			TimeTag:  timeTag,
			Contents: bundleContents,
		},
		remainingBytes,
		nil

}
