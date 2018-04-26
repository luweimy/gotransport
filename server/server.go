package server

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"sync"
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
	mu sync.Mutex
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
	s.mu.Lock()
	if s.ln != nil {
		s.mu.Unlock()
		return ErrMultipleListenCalls
	}

	var (
		ln  net.Listener
		err error
	)
	if s.opts.ConfigTLS != nil {
		ln, err = tls.Listen(network, address, s.opts.ConfigTLS)
	} else {
		ln, err = net.Listen(network, address)
	}
	if err != nil {
		s.mu.Unlock()
		return err
	}
	s.ln = ln
	s.mu.Unlock()

	// listen loop will block the goroutine
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
	defer func() {
		ln.Close()
	}()
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
					return s.ctx.Err()
				}
				continue
			}
			return err
		}
		delay = 0

		gotransport.NewTransport(s.ctx, conn, s.opts).LoopAsync()
	}
}
