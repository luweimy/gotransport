package gotransport

import "github.com/luweimy/gotransport/protocol"

type OnConnectFunc func(Transport) bool
type OnMessageFunc func(Transport, protocol.Protocol)
type OnCloseFunc func(Transport)
type OnErrorFunc func(Transport, error)

type Options struct {
	OnConnect OnConnectFunc
	OnMessage OnMessageFunc
	OnClose   OnCloseFunc
	OnError   OnErrorFunc
	Factory   protocol.Factory
	//workerSize int  // numbers of worker go-routines
	//bufferSize int  // size of buffered channel
}

func MakeOptions() Options {
	return Options{
		OnConnect: func(Transport) bool { return true },
		OnMessage: func(Transport, protocol.Protocol) {},
		OnClose:   func(Transport) {},
		OnError:   func(Transport, error) {},
		Factory:   protocol.PacketFactory{},
	}
}

type OptionFunc func(*Options)

func WithConnect(cb OnConnectFunc) OptionFunc {
	return func(o *Options) {
		o.OnConnect = cb
	}
}

func WithMessage(cb OnMessageFunc) OptionFunc {
	return func(o *Options) {
		o.OnMessage = cb
	}
}

func WithClose(cb OnCloseFunc) OptionFunc {
	return func(o *Options) {
		o.OnClose = cb
	}
}

func WithError(cb OnErrorFunc) OptionFunc {
	return func(o *Options) {
		o.OnError = cb
	}
}

func WithFactory(factory protocol.Factory) OptionFunc {
	return func(o *Options) {
		o.Factory = factory
	}
}
