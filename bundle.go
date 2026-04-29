package osc

import (
	"errors"
)

func (b *OSCBundle) ToBytes() ([]byte, error) {

	bytes := stringToOSCBytes("#bundle")

	timeTagBytes := timeTagToOSCBytes(b.TimeTag)
	bytes = append(bytes, timeTagBytes...)

	for _, packet := range b.Contents {
		packetBytes, err := packet.ToBytes()
		if err != nil {
			return nil, err
		}
		packetLength := len(packetBytes)

		packetLengthBytes := int32ToOSCBytes(int32(packetLength))
		bytes = append(bytes, packetLengthBytes...)
		bytes = append(bytes, packetBytes...)
	}

	return bytes, nil
}

func BundleFromBytes(bytes []byte) (*OSCBundle, []byte, error) {
	if len(bytes) < 20 {
		return nil, bytes, errors.New("OSC Bundle has to be at least 20 bytes")
	}

	if bytes[0] != 35 {
		return nil, bytes, errors.New("OSC Bundle must start with a #")
	}

	bundleHeader, bytesAfterBundleHeader, err := readOSCString(bytes)

	if err != nil {
		return nil, bytes, err
	}

	if bundleHeader != "#bundle" {
		return nil, bytesAfterBundleHeader, errors.New("OSC Bundle must start with #bundle string")
	}

	timeTag, bytesAfterTimeTag, err := readOSCTimeTag(bytesAfterBundleHeader)

	if err != nil {
		return nil, bytesAfterBundleHeader, err
	}

	bundleContents := []OSCPacket{}

	endOfBundle := false

	remainingBytes := bytesAfterTimeTag

	for !endOfBundle {
		contentSize, bytesAfterContentSize, err := readOSCInt32(remainingBytes)

		if err != nil {
			return nil, remainingBytes, err
		}

		remainingBytes = bytesAfterContentSize

		if len(remainingBytes) < int(contentSize) {
			return nil, remainingBytes, errors.New("bundle doesn't have enough bytes for the content size it specifies")
		}

		bundleContentBytes := remainingBytes[0:contentSize]

		if bundleContentBytes[0] == 35 {
			content, _, err := BundleFromBytes(bundleContentBytes)
			if err != nil {
				return nil, remainingBytes, err
			}
			bundleContents = append(bundleContents, content)
		} else if bundleContentBytes[0] == 47 {
			content, err := MessageFromBytes(bundleContentBytes)
			if err != nil {
				return nil, remainingBytes, err
			}
			bundleContents = append(bundleContents, content)
		} else {
			return nil, remainingBytes, errors.New("bundle contents does not look a bundle or message")
		}
		remainingBytes = bytesAfterContentSize[contentSize:]
		if len(remainingBytes) == 0 {
			endOfBundle = true
		}

	}

	return &OSCBundle{
			TimeTag:  timeTag,
			Contents: bundleContents,
		},
		remainingBytes,
		nil

}
