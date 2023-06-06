package client

import (
	"context"
	"log"
	"net"
)

type Options struct{}

type Client struct {
	conn net.Conn
}

func New(endpoint string, opts Options) (*Client, error) {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		log.Fatal(err)
	}
	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Set(ctx context.Context, key, value []byte) (any, error) {
	return nil, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
