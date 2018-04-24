package gotransport

import (
	"bufio"
	"context"
	"io"
	"net"
)

const (
	DefaultBufferSize = 1024
)

// Transport representing a network connection
type Transport interface {
	// Write writes data to the connection.
	// Write can be made to time out and return an Error with Timeout() == true
	// after a fixed time limit; see SetDeadline and SetWriteDeadline.
	Write(b []byte) (n int, err error)
	WriteString(s string) (n int, err error)
	WritePacket(packet Protocol) (n int, err error)

	// Close closes the connection.
	// Any blocked Read or Write operations will be unblocked and return errors.
	Close() error
	IsClosed() bool

	// Peer returns the remote network address.
	Peer() net.Addr

	// Host returns the local network address.
	Host() net.Addr

	ProtocolMake() Protocol

	//Error() <-chan error
	//Done() <-chan struct{}
}

type TransportHijacker interface {
	Hijack() net.Conn
}

type transport struct {
	ctx  context.Context
	conn net.Conn
	rw   io.ReadWriter // 方便hook conn的读写
	opts *Options
}

func NewTransport(ctx context.Context, conn net.Conn, opts *Options) *transport {
	if opts == nil {
		opts = MakeOptions()
	}
	t := &transport{
		ctx:  ctx,
		conn: conn,
		rw:   conn,
		opts: opts,
	}
	for _, hook := range opts.ConnHooks {
		t.conn = hook(t.conn)
	}
	t.rw = t.conn
	for _, hook := range opts.Hooks {
		t.rw = hook(t.rw)
	}
	return t
}

func (t *transport) Hijack() net.Conn {
	return t.conn
}

func (t *transport) ProtocolMake() Protocol {
	return t.opts.Factory()
}

func (t *transport) Write(b []byte) (n int, err error) {
	packet := t.ProtocolMake()
	packet.SetPayload(b)
	return packet.WriteTo(t.rw)
}

func (t *transport) WriteString(s string) (n int, err error) {
	return t.Write([]byte(s))
}

func (t *transport) WritePacket(packet Protocol) (n int, err error) {
	return packet.WriteTo(t.rw)
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
	_, err := t.rw.Read([]byte{})
	return IsClosedConnError(err)
}

func (t *transport) Peer() net.Addr {
	return t.conn.RemoteAddr()
}

func (t *transport) Host() net.Addr {
	return t.conn.LocalAddr()
}

func (t *transport) LoopSync() *transport {
	if t.opts.OnConnected != nil && !t.opts.OnConnected(t) {
		t.conn.Close()
		return t
	}
	t.readLoop()
	return t
}

func (t *transport) LoopAsync() *transport {
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
	reader := bufio.NewReaderSize(t.rw, bufferSize)
	for {
		// TODO: timeout 不准
		select {
		case <-t.ctx.Done():
			return
		default:
		}
		packet := t.ProtocolMake()
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
