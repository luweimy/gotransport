package gotransport

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

var ErrTooLarge = errors.New("packet: too large")

const (
	MaxPacketSize = 1024 * 1024 * 100 // 100MB
	HeaderSize    = 5                 // type(1-byte) + length(4-byte)
)

// packetProtocol
// message format:
//  [00000000][00000000][00000000][00000000][00000000][00000000][00000000]...
//  | (uint8)||               (uint32)               ||     (binary)
//  |  1-byte||                4-byte                ||      N-byte
//  ----------------------------------------------------------------------...
//      type                   length                        value
//       \-----------------------/
//            header(5-byte)
type packetProtocol struct {
	tag   byte
	value []byte
}

func PacketProtocol() Protocol {
	return &packetProtocol{}
}

func (p *packetProtocol) Payload() []byte {
	return p.value
}

func (p *packetProtocol) SetPayload(payload []byte) {
	p.value = payload
}

func (p *packetProtocol) SetFlags(value interface{}) error {
	if tag, ok := value.(byte); ok {
		p.tag = tag
		return nil
	}
	return ErrTypeNotSupport
}

func (p *packetProtocol) Flags() Flags {
	return FlagsValue{p.tag}
}

func (p *packetProtocol) WriteTo(w io.Writer) (int, error) {
	if len(p.value)+HeaderSize > MaxPacketSize {
		return 0, ErrTooLarge
	}
	var (
		total int
	)

	if err := binary.Write(w, binary.BigEndian, p.tag); err != nil {
		return total, err
	}
	total += 1

	if err := binary.Write(w, binary.BigEndian, uint32(len(p.value))); err != nil {
		return total, err
	}
	total += 4

	n, err := w.Write(p.value)
	total += n
	if err != nil {
		return total, err
	}
	return total, nil
}

func (p *packetProtocol) ReadFrom(r io.Reader) (int, error) {
	var (
		total int
	)

	if err := binary.Read(r, binary.BigEndian, &p.tag); err != nil {
		return total, err
	}
	total += 1

	var length uint32 = 0
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return total, err
	}
	total += 4

	if length+HeaderSize > MaxPacketSize {
		return total, ErrTooLarge
	}

	var value = make([]byte, length)
	n, err := io.ReadFull(r, value)
	p.value = value[:n]
	total += n
	if err != nil {
		return total, err
	}
	return total, nil
}

func (p *packetProtocol) Pack() ([]byte, error) {
	buf := &bytes.Buffer{}
	if _, err := p.WriteTo(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *packetProtocol) Unpack(data []byte) (int, error) {
	return p.ReadFrom(bytes.NewBuffer(data))
}
