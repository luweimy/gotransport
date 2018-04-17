package server

import (
	"context"
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
	opts gotransport.Options
	ctx  context.Context

	ln net.Listener
	mu sync.Mutex
}

func New(opts ...gotransport.OptionFunc) *Server {
	s := &Server{
		ctx:  context.Background(),
		opts: gotransport.MakeOptions(),
	}
	for _, o := range opts {
		o(&s.opts)
	}
	return s
}

func (s *Server) Listen(network, address string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
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

func (s *Server) Host() net.Addr {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ln == nil {
		return nil
	}
	return s.ln.Addr()
}

func (s *Server) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ln == nil {
		return nil
	}
	ln := s.ln
	s.ln = nil
	return ln.Close()
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

		transport := gotransport.NewTransport(conn, s.opts)
		if s.opts.OnConnect != nil && !s.opts.OnConnect(transport) {
			conn.Close()
			continue
		}
		go transport.ReadLoop()
	}
}
