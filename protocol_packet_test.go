package gotransport

import (
	"bytes"
	"testing"
)

func assert(v bool) {
	if !v {
		panic("assert")
	}
}

func assertErr(err error) {
	if err != nil {
		panic(err)
	}
}

func TestPacket_Pack(t *testing.T) {
	p := packetProtocol{}
	err := p.SetFlagOptions(byte(0x02))
	assertErr(err)
	p.SetPayload([]byte{3, 2})

	packedData, err := p.Pack()
	assertErr(err)

	assert(len(packedData) == 7)
	assert(bytes.Compare([]byte{2, 0, 0, 0, 2, 3, 2}, packedData) == 0)

	p2 := packetProtocol{}
	n, err := p2.Unpack(packedData)
	assertErr(err)
	assert(n == 7)
	assert(p2.FlagOptions().Byte() == 0x02)
	assert(bytes.Compare(p.Payload(), p2.Payload()) == 0)

	buf1 := bytes.NewBuffer(packedData)
	p3 := packetProtocol{}
	n, err = p3.ReadFrom(buf1)
	assertErr(err)
	assert(n == 7)
	assert(p3.FlagOptions().Byte() == 0x02)
	assert(bytes.Compare(p3.Payload(), p.Payload()) == 0)

	buf2 := bytes.NewBuffer([]byte{})
	n, err = p.WriteTo(buf2)
	assertErr(err)
	assert(n == 7)
	assert(bytes.Compare(buf2.Bytes(), packedData) == 0)
}
