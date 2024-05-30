package main

import (
	"bytes"
	"net"
)

type client struct {
	conn net.Conn
}

func newClient(con net.Conn) *client {
	return &client{
		conn: con,
	}
}

func (c *client) handle(msg []byte) {
	split := bytes.Split(msg, []byte(" "))
	cmd := bytes.ToUpper(bytes.TrimSpace(split[0]))

	switch string(cmd) {
	case "PING":
		resMsg := "\rPong\r\n"
		c.conn.Write([]byte(resMsg))
	case "ECHO":
		args := bytes.Join(split[1:], []byte(" "))
		resMsg := "\r" + string(args) + "\r\n"
		c.conn.Write([]byte(resMsg))
	case "GET":
		c.get(split[1:])
	case "SET":
		c.set(split[1:])
	case "DEL":
		c.del(split[1:])
	case "QUIT":
		c.conn.Close()
	default:
		c.conn.Write([]byte("\r🧐Unknown command\r\n"))
	}
}

func (c *client) get(args [][]byte) bool {
	if len(args) != 1 {
		resMsg := "\r" + "Wrong number of argumnets passed." + "\r\n"
		c.conn.Write([]byte(resMsg))
		return false
	}
	key := string(bytes.TrimSpace(args[0]))
	val, _ := Cache.Load(key)
	if val != nil {
		resMsg := "\r" + val.(string) + "\r\n"
		c.conn.Write([]byte(resMsg))
	} else {
		resMsg := "\r-1\r\n"
		c.conn.Write([]byte(resMsg))
	}
	return true
}

func (c *client) set(args [][]byte) bool {
	if len(args) != 2 {
		resMsg := "\r" + "Wrong number of argumnets passed." + "\r\n"
		c.conn.Write([]byte(resMsg))
		return false
	}
	key := string(bytes.TrimSpace(args[0]))
	val := string(bytes.TrimSpace(args[1]))
	Cache.Store(key, val)
	resMsg := "\r" + "OK" + "\r\n"
	c.conn.Write([]byte(resMsg))
	return true
}

func (c *client) del(args [][]byte) bool {
	if len(args) != 1 {
		resMsg := "\r" + "Wrong number of argumnets passed." + "\r\n"
		c.conn.Write([]byte(resMsg))
		return false
	}
	key := string(bytes.TrimSpace(args[0]))
	_, exists := Cache.LoadAndDelete(key)
	resMsg := ""
	if exists {
		resMsg = "\r" + "Deleted" + "\r\n"
	} else {
		resMsg = "\r" + "Key doesn't exists" + "\r\n"
	}
	c.conn.Write([]byte(resMsg))
	return true
}
