package protocol

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
	p := Packet{}
	p.SetType(0x02)
	p.SetPayload([]byte{3, 2})

	packedData, err := p.Pack()
	assertErr(err)

	assert(len(packedData) == 7)
	assert(bytes.Compare([]byte{2, 0, 0, 0, 2, 3, 2}, packedData) == 0)

	p2 := Packet{}
	n, err := p2.Unpack(packedData)
	assertErr(err)
	assert(n == 7)
	assert(p2.Type() == 0x02)
	assert(bytes.Compare(p.Payload(), p2.Payload()) == 0)

	buf1 := bytes.NewBuffer(packedData)
	p3 := Packet{}
	n, err = p3.ReadFrom(buf1)
	assertErr(err)
	assert(n == 7)
	assert(p3.Type() == 0x02)
	assert(bytes.Compare(p3.Payload(), p.Payload()) == 0)

	buf2 := bytes.NewBuffer([]byte{})
	n, err = p.WriteTo(buf2)
	assertErr(err)
	assert(n == 7)
	assert(bytes.Compare(buf2.Bytes(), packedData) == 0)
}
