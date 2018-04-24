package gotransport

import (
	"io"
	"net"
)

type ConnectHandler func(transport Transport) bool
type MessageHandler func(transport Transport, packet Protocol)
type CloseHandler func(transport Transport)
type ErrorHandler func(transport Transport, err error)

type ReadWriterHookHandler func(rw io.ReadWriter) io.ReadWriter
type ConnHookHandler func(conn net.Conn) net.Conn

type Options struct {
	OnConnected ConnectHandler
	OnMessage   MessageHandler
	OnClosed    CloseHandler
	OnClosing   CloseHandler
	OnError     ErrorHandler
	Factory     ProtocolFactory
	BufferSize  int // size of transport reader buffer
	Hooks       []ReadWriterHookHandler
	ConnHooks   []ConnHookHandler
	//workerSize int  // numbers of worker go-routines
}

func MakeOptions() *Options {
	return &Options{
		Factory: PacketProtocol,
	}
}

type OptionFunc func(*Options)

func WithConnected(cb ConnectHandler) OptionFunc {
	return func(o *Options) {
		o.OnConnected = cb
	}
}

func WithMessage(cb MessageHandler) OptionFunc {
	return func(o *Options) {
		o.OnMessage = cb
	}
}

func WithClosed(cb CloseHandler) OptionFunc {
	return func(o *Options) {
		o.OnClosed = cb
	}
}

func WithClosing(cb CloseHandler) OptionFunc {
	return func(o *Options) {
		o.OnClosing = cb
	}
}

func WithError(cb ErrorHandler) OptionFunc {
	return func(o *Options) {
		o.OnError = cb
	}
}

func WithFactory(factory ProtocolFactory) OptionFunc {
	return func(o *Options) {
		o.Factory = factory
	}
}

func WithBufferSize(bufferSize int) OptionFunc {
	return func(o *Options) {
		o.BufferSize = bufferSize
	}
}

// Hook conn read and write
func WithHook(hook ReadWriterHookHandler) OptionFunc {
	return func(o *Options) {
		o.Hooks = append(o.Hooks, hook)
	}
}

func WithHookConn(hook ConnHookHandler) OptionFunc {
	return func(o *Options) {
		o.ConnHooks = append(o.ConnHooks, hook)
	}
}
