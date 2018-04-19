package gotransport

import (
	"errors"
	"io"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

type Protocol interface {
	WriteTo(w io.Writer) (int, error)
	ReadFrom(r io.Reader) (int, error)

	Payload() []byte
	SetPayload([]byte)

	Type() byte
	SetType(byte)
}

type ProtocolFactory func() Protocol
