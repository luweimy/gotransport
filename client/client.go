package client

import (
	"context"
	"errors"
	"net"
	"sync"

	"github.com/luweimy/gotransport"
)

var (
	ErrMultipleConnectCalls = errors.New("server multiple connect calls")
)

type Client struct {
	opts *gotransport.Options
	ctx  context.Context

	gotransport.Transport
	mu sync.Mutex
}

func New(opts ...gotransport.OptionFunc) *Client {
	c := &Client{
		opts: gotransport.MakeOptions(),
		ctx:  context.Background(),
	}
	c.Options(opts...)
	return c
}

func (c *Client) Options(opts ...gotransport.OptionFunc) {
	for _, o := range opts {
		o(c.opts)
	}
}

func (c *Client) Connect(network, address string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Transport != nil && !c.Transport.IsClosed() {
		return ErrMultipleConnectCalls
	}
	conn, err := net.Dial(network, address)
	if err != nil {
		return err
	}
	c.Transport = gotransport.NewTransport(c.ctx, conn, c.opts).LoopAsync()
	return nil
}
