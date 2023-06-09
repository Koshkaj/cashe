package client

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Koshkaj/cashe/core"
)

type Options struct{}

type Client struct {
	conn net.Conn
}

func NewFromConn(conn net.Conn) *Client {
	return &Client{
		conn: conn,
	}
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

func (c *Client) Get(ctx context.Context, key []byte) ([]byte, error) {
	cmd := &core.CommandGet{
		Key: key,
	}
	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, err
	}

	resp, err := core.ParseGetResponse(c.conn)
	if err != nil {
		return nil, err
	}
	if resp.Status == core.StatusKeyNotFound {
		return nil, fmt.Errorf("could not find key (%s)", resp.Status)
	}
	if resp.Status != core.StatusOK {
		return nil, fmt.Errorf("server responsed with error %s", resp.Status)
	}
	return resp.Value, nil
}

func (c *Client) Delete(ctx context.Context, key []byte) error {
	cmd := &core.CommandDel{
		Key: key,
	}
	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Set(ctx context.Context, key []byte, value []byte, ttl int) error {
	cmd := &core.CommandSet{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}
	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return err
	}

	resp, err := core.ParseSetResponse(c.conn)
	if err != nil {
		return err
	}
	if resp.Status != core.StatusOK {
		return fmt.Errorf("server responsed with error %s", resp.Status)
	}
	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
