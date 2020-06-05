package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

var _ TelnetClient=(*client)(nil)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

type client struct {
	address string
	timeout time.Duration
	in io.ReadCloser
	out io.Writer
	connection net.Conn
	inbuf []byte
	outbuf []byte
}

func (c *client) Connect() (e error) {
	c.connection, e= net.DialTimeout("tcp", c.address,c.timeout)
	return
}

func (c *client) Close() error {
	if c.connection!=nil{
		return c.connection.Close()
	}
	return fmt.Errorf("connection is nil")
}

func (c *client) Send() (e error) {
	_,e=io.Copy(c.connection,c.in)
	return

	/*scanner := bufio.NewScanner(c.in)
//OUTER:
	for {
		select {
		//case <-ctx.Done():
		//	break OUTER
		default:
			if !scanner.Scan() {
				return fmt.Errorf("cannot scan in")
			}
			text := scanner.Text()
			_,e=fmt.Fprintln(c.connection,text)
			if e!=nil{
				fmt.Println(e)
				return
			}
		}
	}
	fmt.Println(e)
	return
	/*_,e:=c.in.Read(c.inbuf)
	if e!=nil{
		return e
	}
	_,e=c.connection.Write(c.inbuf)
	if e!=nil{
		return e
	}
	return nil*/
}

func (c *client) Receive() (e error) {
	_,e=io.Copy(c.out,c.connection)
	return

/*scanner := bufio.NewScanner(c.connection)
	OUTER:
		for {
			select {
			case <-ctx.Done():
				break OUTER
			default:
				if !scanner.Scan() {
					return fmt.Errorf("cannot scan connection")
				}
				text := scanner.Text()
				_,e:=fmt.Fprint(c.out,text)
				if e!=nil{
					return e
				}
				log.Printf("From server: %s", text)
			}
		}
		log.Printf("Finished readRoutine")*/
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	c:=client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
	return &c
}

// Place your code here
// P.S. Author's solution takes no more than 50 lines
