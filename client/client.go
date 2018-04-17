package client

import (
	"net"
)

type Client struct {
	conn net.Conn
}

func New() *Client {
	return &Client{}
}

func (c *Client) Connect(network, address string) error {
	conn, err := net.Dial(network, address)
	if err != nil {
		return err
	}
	c.conn = conn

	//c.conn.Write()

	return nil
}

//func (c *Client) Write(b []byte) (n int, err error) {
//
//}
