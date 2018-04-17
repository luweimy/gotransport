package protocol

import (
	"errors"
	"io"
)

var (
	ErrNotImplement = errors.New("not implement")
)

type Protocol interface {
	WriteTo(w io.Writer) (int, error)
	ReadFrom(r io.Reader) (int, error)

	Payload() []byte
	Type() byte
	SetPayload([]byte)
	SetType(byte)
}

type Factory interface {
	Build() Protocol
}
