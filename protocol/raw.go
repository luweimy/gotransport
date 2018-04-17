package protocol

import (
	"io"
)

type Raw struct {
	data []byte
}

func NewRaw() *Raw {
	return &Raw{}
}

func (r *Raw) Payload() []byte {
	return r.data
}

func (l *Raw) Type() byte {
	return 0
}

func (p *Raw) SetPayload(payload []byte) {
	p.data = payload
}

func (p *Raw) SetType(tp byte) {
	panic(ErrNotImplement)
}

func (p *Raw) WriteTo(w io.Writer) (int, error) {
	return w.Write(p.data)
}

func (p *Raw) ReadFrom(r io.Reader) (int, error) {
	p.data = make([]byte, 1024)
	n, err := r.Read(p.data)
	p.data = p.data[:n]
	return n, err
}

type RawFactory struct {
}

func (RawFactory) Build() Protocol {
	return NewRaw()
}
