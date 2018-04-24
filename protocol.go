package gotransport

import (
	"errors"
	"io"
)

var (
	ErrFlagsNotSupport = errors.New("protocol: flags not support")
	ErrTypeNotSupport  = errors.New("protocol: type not support")
)

type ProtocolFactory func() Protocol

type Protocol interface {
	WriteTo(w io.Writer) (n int, err error)

	// ReadFrom decode from reader(net.Conn)
	// err io.EOF must be throw out
	// io.EOF错误务必要抛出来，这样transport可以自动关闭连接。
	// 最终读取的Payload长度可能与返回值n不等，这是正常的，因为内部的解码会造成这种情况
	ReadFrom(r io.Reader) (n int, err error)

	SetPayload([]byte)
	Payload() []byte

	SetFlags(value interface{}) error
	Flags() Flags
}

type Flags interface {
	Byte() byte
	Int8() int8
	Int16() int16
	Int32() int32
	Int64() int64
	Uint8() uint8
	Uint16() uint16
	Uint32() uint32
	Uint64() uint64
	Raw() interface{}
}

type FlagsValue struct {
	value interface{}
}

func (v FlagsValue) Byte() byte {
	return v.value.(byte)
}

func (v FlagsValue) Int8() int8 {
	return v.value.(int8)
}

func (v FlagsValue) Int16() int16 {
	return v.value.(int16)
}

func (v FlagsValue) Int32() int32 {
	return v.value.(int32)
}

func (v FlagsValue) Int64() int64 {
	return v.value.(int64)
}

func (v FlagsValue) Uint8() uint8 {
	return v.value.(uint8)
}

func (v FlagsValue) Uint16() uint16 {
	return v.value.(uint16)
}

func (v FlagsValue) Uint32() uint32 {
	return v.value.(uint32)
}

func (v FlagsValue) Uint64() uint64 {
	return v.value.(uint64)
}

func (v FlagsValue) Raw() interface{} {
	return v.value
}
