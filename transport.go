package gotransport

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/luweimy/gotransport/protocol"
)

const (
	ReaderBufferSize = 1024
)

// 代表一个连接
type Transport interface {
	Write(b []byte) (n int, err error)
	WriteString(s string) (n int, err error)
	WritePacket(packet protocol.Protocol) (n int, err error)
	//SetWriteDeadline(t time.Time) error

	Close() error   // 关闭连接
	Peer() net.Addr // 获取对方信息
	Host() net.Addr // 获取本机信息
	//Error() <-chan struct{}
}

type transport struct {
	conn net.Conn
	opts Options
}

func NewTransport(conn net.Conn, opts Options) *transport {
	return &transport{
		conn: conn,
		opts: opts,
	}
}

func (t *transport) ProtocolFactory() protocol.Factory {
	return t.opts.Factory
}

func (t *transport) Write(b []byte) (n int, err error) {
	packet := t.ProtocolFactory().Build()
	packet.SetPayload(b)
	return packet.WriteTo(t.conn)
}

func (t *transport) WriteString(s string) (n int, err error) {
	return t.Write([]byte(s))
}

func (t *transport) WritePacket(packet protocol.Protocol) (n int, err error) {
	return packet.WriteTo(t.conn)
}

func (t *transport) Close() error {
	if err := t.conn.Close(); err != nil {
		return err
	}
	if t.opts.OnClose != nil {
		t.opts.OnClose(t)
	}
	return nil
}

func (t *transport) Peer() net.Addr {
	return t.conn.RemoteAddr()
}

func (t *transport) Host() net.Addr {
	return t.conn.LocalAddr()
}

func (t *transport) ReadLoop() {
	defer func() {
		if err := recover(); err != nil {
			t.opts.OnError(t, errorWrap(err))
		}
		if err := t.Close(); err != nil {
			t.opts.OnError(t, err)
			return
		}
	}()

	reader := bufio.NewReaderSize(t.conn, ReaderBufferSize)
	for {
		packet := t.opts.Factory.Build()
		_, err := packet.ReadFrom(reader)
		if err != nil {
			if err == io.EOF {
				return
			}
			t.opts.OnError(t, err)
			return
		}
		t.opts.OnMessage(t, packet)
	}
}

func errorWrap(v interface{}) error {
	return fmt.Errorf("%v", v)
}
