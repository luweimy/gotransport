package gotransport

import (
	"errors"
	"io"
)

var (
	ErrOptionsNotSupport = errors.New("protocol: options not support")
	ErrTypeNotSupport    = errors.New("protocol: type not support")
)

type ProtocolFactory func() Protocol

type Protocol interface {
	// 最终写入的长度就是返回值n
	WriteTo(w io.Writer) (n int, err error)

	// ReadFrom decode from reader(net.Conn)
	// err io.EOF must be throw out
	// io.EOF错误务必要抛出来，这样transport可以自动关闭连接。
	// 读取的Payload长度可能与返回值n不等，这是正常的，因为内部的解码会造成这种情况
	ReadFrom(r io.Reader) (n int, err error)

	// 设置和获取有效载荷数据
	SetPayload([]byte)
	Payload() []byte

	// 设置/获取一些协议支持的选项和标记等
	SetFlagOptions(value interface{}) error
	FlagOptions() Value
}

type Value interface {
	Byte() byte
	Bytes() []byte
	Int8() int8
	Int16() int16
	Int32() int32
	Int64() int64
	Uint8() uint8
	Uint16() uint16
	Uint32() uint32
	Uint64() uint64
	Float32() float32
	Float64() float64
	Raw() interface{}
}

type valueBox struct {
	value interface{}
}

func WrapValue(v interface{}) Value {
	return &valueBox{
		value: v,
	}
}

func (v valueBox) Byte() byte {
	return v.value.(byte)
}

func (v valueBox) Bytes() []byte {
	return v.value.([]byte)
}

func (v valueBox) Int8() int8 {
	return v.value.(int8)
}

func (v valueBox) Int16() int16 {
	return v.value.(int16)
}

func (v valueBox) Int32() int32 {
	return v.value.(int32)
}

func (v valueBox) Int64() int64 {
	return v.value.(int64)
}

func (v valueBox) Uint8() uint8 {
	return v.value.(uint8)
}

func (v valueBox) Uint16() uint16 {
	return v.value.(uint16)
}

func (v valueBox) Uint32() uint32 {
	return v.value.(uint32)
}

func (v valueBox) Uint64() uint64 {
	return v.value.(uint64)
}

func (v valueBox) Float32() float32 {
	return v.value.(float32)
}

func (v valueBox) Float64() float64 {
	return v.value.(float64)
}

func (v valueBox) Raw() interface{} {
	return v.value
}
