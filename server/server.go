package server

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/luweimy/gotransport"
)

var (
	ErrMultipleListenCalls = errors.New("server multiple listen calls")
)

type Server struct {
	opts *gotransport.Options
	ctx  context.Context

	ln net.Listener
}

func New(opts ...gotransport.OptionFunc) *Server {
	s := &Server{
		opts: gotransport.MakeOptions(),
		ctx:  context.Background(),
	}
	s.Options(opts...)
	return s
}

func (s *Server) Options(opts ...gotransport.OptionFunc) {
	for _, o := range opts {
		o(s.opts)
	}
}

// Listen announces on the local network address.
//
// The network must be "tcp", "tcp4", "tcp6", "unix" or "unixpacket".
func (s *Server) Listen(network, address string) error {
	if s.ln != nil {
		return ErrMultipleListenCalls
	}
	ln, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	s.ln = ln
	return s.listenLoop(ln)
}

// Addr returns the listener's network address, a *TCPAddr.
// The Addr returned is shared by all invocations of Addr, so
// do not modify it.
func (s *Server) Addr() net.Addr {
	return s.ln.Addr()
}

// Close stops listening on the TCP address.
// Already Accepted connections are not closed.
func (s *Server) Close() error {
	return s.ln.Close()
}

func (s *Server) listenLoop(ln net.Listener) error {
	var delay time.Duration
	for {
		conn, err := ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if delay == 0 {
					delay = 5 * time.Millisecond
				} else {
					delay *= 2
				}
				if delay >= time.Second {
					delay = time.Second
				}
				select {
				case <-time.After(delay):
				case <-s.ctx.Done():
				}
				continue
			}
			return err
		}
		delay = 0

		gotransport.NewTransport(conn, s.opts).AsyncLoop()
	}
}
