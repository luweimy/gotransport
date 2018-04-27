package gotransport

import (
	"bufio"
	"context"
	"net"
)

const (
	BufferSize = 1024
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

	// Notify close, the error is the reason of transport close
	Done() <-chan error
}

// TransportHijacker hijack transport net.Conn
type TransportHijacker interface {
	Hijack() net.Conn
}

type transport struct {
	ctx  context.Context
	opts *Options
	conn net.Conn

	// Done chan
	doneCh chan error
}

func NewTransport(ctx context.Context, conn net.Conn, opts *Options) *transport {
	if opts == nil {
		opts = MakeOptions()
	}
	t := &transport{
		ctx:  ctx,
		opts: opts,
		conn: conn,

		doneCh: make(chan error, 1),
	}
	for _, hook := range opts.Hooks {
		t.conn = hook(t.conn)
	}
	return t
}

func (t *transport) Write(b []byte) (n int, err error) {
	packet := t.ProtocolMake()
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
	return t.close(nil)
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

func (t *transport) Done() <-chan error {
	return t.doneCh
}

func (t *transport) Hijack() net.Conn {
	return t.conn
}

func (t *transport) ProtocolMake() Protocol {
	return t.opts.Factory()
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
	var readErr error
	defer func() {
		if err := recover(); err != nil {
			readErr = errorWrap(err)
		}
		t.close(readErr)
	}()

	buffSize := BufferSize
	if t.opts.BufferSize > 0 {
		buffSize = t.opts.BufferSize
	}
	reader := bufio.NewReaderSize(t.conn, buffSize)
	for {
		select {
		case <-t.ctx.Done():
			readErr = t.ctx.Err()
			return
		default:
		}
		packet := t.ProtocolMake()
		_, readErr = packet.ReadFrom(reader)
		if readErr != nil {
			return
		}
		t.notify(packet)
	}
}

func (t *transport) close(err error) error {
	if t.opts.OnClosing != nil {
		t.opts.OnClosing(t, err)
	}
	t.doneCh <- err
	// close the conn
	closeErr := t.conn.Close()
	if t.opts.OnClosed != nil {
		t.opts.OnClosed(t, closeErr)
	}
	return closeErr
}

func (t *transport) notify(packet Protocol) {
	if t.opts.OnMessage != nil {
		t.opts.OnMessage(t, packet)
	}
}
