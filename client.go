package main

import (
	"net"
)

type Client struct {
	conn net.Conn
}

func NewClient(con net.Conn) *Client {
	return &Client{
		conn: con,
	}
}

func (c *Client) SendResponse(response string) {
	response += "\n"
	c.conn.Write([]byte(response))
}

func (c *Client) CloseConnection() error {
	return c.conn.Close()
}
