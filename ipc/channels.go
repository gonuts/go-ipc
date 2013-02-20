package ipc

import (
	"bytes"
	//"fmt"
	"io"
)

type Channel struct {
	Name       string
	in         *bytes.Buffer
	out        *bytes.Buffer
	Id         int
	inBlocker  chan int
	outBlocker chan int
	closed     bool
}

func NewChannel(name string) *Channel {
	c := new(Channel)
	c.Name = name

	c.init()

	return c
}

func (c *Channel) init() {
	c.closed = false
	c.in = new(bytes.Buffer)
	c.out = new(bytes.Buffer)
	c.Id = 0
	c.inBlocker = make(chan int)
	c.outBlocker = make(chan int)
}

func (c *Channel) Write(data []byte) (int, error) {

	a, b := c.out.Write(data)

	select {
	case c.outBlocker <- 1:
	default:
	}
	return a, b
}

func (c *Channel) writeToBuffer(data []byte) (int, error) {
	a, b := c.in.Write(data)

	select {
	case c.inBlocker <- 1:
	default:
	}
	return a, b
}

func (c *Channel) readFromBuffer(buf []byte) (int, error) {
	if c.out.Len() == 0 {
		<-c.outBlocker
	}

	n, err := c.out.Read(buf)

	if err == io.EOF {
		return n, nil
	}
	return n, err
}

func (c *Channel) Read(buf []byte) (int, error) {

	if c.in.Len() == 0 {
		<-c.inBlocker
	}

	n, err := c.in.Read(buf)

	if err == io.EOF {
		return n, nil
	}
	return n, err
}

func (c *Channel) Close() error {
	c.closed = true
	c.inBlocker <- 1
	return nil
}

func (c *Channel) Closed() bool {
	return c.closed
}
