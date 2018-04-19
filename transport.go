package gotransport

import (
	"bufio"
	"io"
	"net"
)

const (
	DefaultBufferSize = 1024
)

// 代表一个连接
type Transport interface {
	Write(b []byte) (n int, err error)
	WriteString(s string) (n int, err error)
	WritePacket(packet Protocol) (n int, err error)

	//SetWriteDeadline(t time.Time) error

	Close() error   // 关闭连接
	IsClosed() bool // 是否已经关闭
	Peer() net.Addr // 获取对方信息
	Host() net.Addr // 获取本机信息

	//Error() <-chan error
	//Done() <-chan struct{}

	Protocol() Protocol
}

type transport struct {
	opts *Options
	conn net.Conn
}

func NewTransport(conn net.Conn, opts *Options) *transport {
	if opts == nil {
		opts = MakeOptions()
	}
	t := &transport{
		opts: opts,
		conn: conn,
	}
	return t
}

func (t *transport) Protocol() Protocol {
	return t.opts.Factory()
}

func (t *transport) Write(b []byte) (n int, err error) {
	packet := t.Protocol()
	packet.SetPayload(b)
	return packet.WriteTo(t.conn)
}

func (t *transport) WriteString(s string) (n int, err error) {
	return t.Write([]byte(s))
}

func (t *transport) WritePacket(packet Protocol) (n int, err error) {
	return packet.WriteTo(t.conn)
}

func (t *transport) Close() error {
	if t.opts.OnClosing != nil {
		t.opts.OnClosing(t)
	}
	if err := t.conn.Close(); err != nil {
		return err
	}
	if t.opts.OnClosed != nil {
		t.opts.OnClosed(t)
	}
	return nil
}

func (t *transport) IsClosed() bool {
	if t.conn == nil {
		return true
	}
	_, err := t.conn.Read([]byte{})
	return IsClosedConnError(err)
}

func (t *transport) Peer() net.Addr {
	return t.conn.RemoteAddr()
}

func (t *transport) Host() net.Addr {
	return t.conn.LocalAddr()
}

func (t *transport) SyncLoop() *transport {
	if t.opts.OnConnected != nil && !t.opts.OnConnected(t) {
		t.conn.Close()
		return t
	}
	t.readLoop()
	return t
}

func (t *transport) AsyncLoop() *transport {
	if t.opts.OnConnected != nil && !t.opts.OnConnected(t) {
		t.conn.Close()
		return t
	}
	go t.readLoop()
	return t
}

func (t *transport) readLoop() {
	defer func() {
		if err := recover(); err != nil {
			t.callErr(errorWrap(err))
		}
		if err := t.Close(); err != nil {
			t.callErr(err)
			return
		}
	}()

	bufferSize := DefaultBufferSize
	if t.opts.BufferSize > 0 {
		bufferSize = t.opts.BufferSize
	}
	reader := bufio.NewReaderSize(t.conn, bufferSize)
	for {
		packet := t.Protocol()
		_, err := packet.ReadFrom(reader)
		if err != nil {
			if err == io.EOF {
				return
			}
			t.callErr(err)
			return
		}
		if t.opts.OnMessage != nil {
			t.opts.OnMessage(t, packet)
		}
	}
}

func (t *transport) callErr(err error) {
	if t.opts.OnError != nil {
		t.opts.OnError(t, err)
	}
}
