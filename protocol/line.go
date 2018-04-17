package protocol

import (
	"bufio"
	"bytes"
	"io"
)

type Line struct {
	data []byte
}

func NewLine() *Line {
	return &Line{}
}

func (l *Line) Payload() []byte {
	return l.data
}

func (l *Line) Type() byte {
	return 0
}

func (p *Line) SetPayload(payload []byte) {
	p.data = payload
}

func (p *Line) SetType(tp byte) {
	panic(ErrNotImplement)
}

func (p *Line) WriteTo(w io.Writer) (int, error) {
	// windows:\r\n，unix:\n，mac:\n
	if p.data[len(p.data)-1] != '\n' {
		p.data = append(p.data, '\n')
	}
	return w.Write(p.data)
}

func (p *Line) ReadFrom(r io.Reader) (int, error) {
	var (
		reader = bufio.NewReader(r)
		buffer = bytes.Buffer{}

		part   []byte
		prefix bool
		err    error
	)
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			return buffer.Len(), err
		}
		buffer.Write(part)
		if !prefix {
			break
		}
	}

	p.data = buffer.Bytes()
	return len(p.data), nil
}

type LineFactory struct {
}

func (LineFactory) Build() Protocol {
	return NewLine()
}
