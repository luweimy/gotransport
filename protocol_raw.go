package gotransport

import (
	"io"
)

type rawProtocol struct {
	data []byte
}

func RawProtocol() Protocol {
	return &rawProtocol{}
}

func (r *rawProtocol) Payload() []byte {
	return r.data
}

func (l *rawProtocol) Type() byte {
	return 0
}

func (p *rawProtocol) SetPayload(payload []byte) {
	p.data = payload
}

func (p *rawProtocol) SetType(tp byte) {
	panic(ErrNotImplemented)
}

func (p *rawProtocol) WriteTo(w io.Writer) (int, error) {
	return w.Write(p.data)
}

func (p *rawProtocol) ReadFrom(r io.Reader) (int, error) {
	p.data = make([]byte, 1024)
	n, err := r.Read(p.data)
	p.data = p.data[:n]
	return n, err
}
