package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

var _ TelnetClient = (*client)(nil)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

type client struct {
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	connection net.Conn
}

//Connect connects to tcp server with address and timeout.
func (c *client) Connect() (e error) {
	c.connection, e = net.DialTimeout("tcp", c.address, c.timeout)
	return
}

//Close closes connect.
func (c *client) Close() error {
	if c.connection == nil {
		return fmt.Errorf("tcp connection is nil")
	}
	return c.connection.Close()
}

//Send reads messages from Reader and sent to connection.
func (c *client) Send() (e error) {
	if c.connection == nil {
		return fmt.Errorf("tcp connection is nil")
	}
	_, e = io.Copy(c.connection, c.in)
	return
}

//Receive reads messages from connection and write to Writer.
func (c *client) Receive() (e error) {
	if c.connection == nil {
		return fmt.Errorf("tcp connection is nil")
	}
	_, e = io.Copy(c.out, c.connection)
	return
}

//NewTelnetClient returns client with fields address,timeout,Reader and Writer.
func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{address: address, timeout: timeout, in: in, out: out}
}
