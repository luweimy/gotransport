package gotransport

import (
	"bufio"
	"bytes"
	"io"
)

type lineProtocol struct {
	data []byte
}

func LineProtocol() Protocol {
	return &lineProtocol{}
}

func (l *lineProtocol) Payload() []byte {
	return l.data
}

func (l *lineProtocol) Type() byte {
	return 0
}

func (p *lineProtocol) SetPayload(payload []byte) {
	p.data = payload
}

func (p *lineProtocol) SetType(tp byte) {
	panic(ErrNotImplemented)
}

func (p *lineProtocol) WriteTo(w io.Writer) (int, error) {
	if p.data[len(p.data)-1] != '\n' {
		p.data = append(p.data, '\n')
	}
	return w.Write(p.data)
}

func (p *lineProtocol) ReadFrom(r io.Reader) (int, error) {
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
