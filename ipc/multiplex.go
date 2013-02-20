package ipc

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sync"
	//"time"
)

type Multiplex struct {
	rwc    io.ReadWriteCloser
	router *Router
}

type Router struct {
	rwc          io.ReadWriter
	channels     map[int]*Channel
	channelNames map[string]int
	lsem         chan int
	channelCount int
	control      *controlChannel
}

func NewRouter(rwc io.ReadWriter) *Router {
	router := new(Router)
	router.rwc = rwc
	router.channels = make(map[int]*Channel)
	router.channelNames = make(map[string]int)
	router.lsem = make(chan int, 5)
	router.channelCount = 0

	return router
}

func (r *Router) readChannel(channel *Channel) {

	buf := make([]byte, 256) // 2 + 254
	for {
		if channel.closed {
			break
		}
		n, err := channel.readFromBuffer(buf[2:])
		//fmt.Println("readChannel", channel.Id, n)
		if err != nil {
			fmt.Errorf(err.Error())
		}
		if n > 0 {
			buf[0] = byte(channel.Id)
			buf[1] = byte(n)
			//r.mutex.Lock()
			r.lsem <- 1
			r.rwc.Write(buf[0 : n+2])
			<-r.lsem
			//r.mutex.Unlock()
			//fmt.Println(string(buf[2:]))
		} else {

			//<-channel.Blocker
		}

	}

}

func (r *Router) gateway() {

	// first 2 bytes of the buffer are used to store channel id and write length
	buf := make([]byte, 254)
	header_buf := make([]byte, 2)
	bufr := bufio.NewReader(r.rwc)

	for {

		n, err := bufr.Read(header_buf)
		fmt.Println("reading")
		if err == io.EOF {
			return
		}

		if n == 2 {
			fmt.Println("gateway ", int(header_buf[0]))
			channel, ok := r.channels[int(header_buf[0])]

			size := int(header_buf[1])
			bcap := cap(buf)
			for size > 0 {

				if size < cap(buf) {
					bcap = size
				} else {
					bcap = cap(buf)
				}

				datan, _ := bufr.Read(buf[0:bcap])
				if ok {
					channel.writeToBuffer(buf[0:datan])
				}
				size -= datan
			}

		} else if n == 1 {
			bufr.UnreadByte()
		} else {
			//<-time.After(time.Millisecond)
		}
	}

}

func (r *Router) Start() {
	// create the comms channel
	coms := NewChannel("COMS")
	r.control = newControlChannel(coms)
	r.RegisterChannel(coms)

	r.gateway()
}

func (r *Router) RegisterChannel(channel *Channel) {
	r.channels[r.channelCount] = channel
	channel.Id = r.channelCount
	fmt.Println("reigstering channel", channel.Id)
	r.channelCount += 1
	go r.readChannel(channel)
}

func NewMultiplex(rwc io.ReadWriteCloser) *Multiplex {
	router := NewRouter(rwc.(io.ReadWriter))

	mpx := &Multiplex{rwc, router}
	go router.Start()
	return mpx
}

func (m *Multiplex) RawReaderChannel(name string) *Channel {
	channel := NewChannel(name)

	m.router.RegisterChannel(channel)
	return channel
}

func (m *Multiplex) RawWriterChannel(name string) *Channel {
	channel := NewChannel(name)

	m.router.RegisterChannel(channel)
	return channel
}

type controlMessageType int

const (
	REGISTER controlMessageType = iota
	RESPONSE
)

type controlChannel struct {
	channel        *Channel
	mutex          *sync.Mutex
	responseBuffer *bytes.Buffer
	blocker        chan int
}

func newControlChannel(channel *Channel) *controlChannel {
	cc := new(controlChannel)
	cc.mutex = new(sync.Mutex)
	cc.responseBuffer = new(bytes.Buffer)
	cc.channel = channel
	cc.blocker = make(chan int)
	go cc.incomingMessages()
}

func (cc *controlChannel) formatMessage(mtype controlMessageType, data []byte) []byte {

	buf := byte(mtype)             // message type
	append(buf, []byte(len(data))) // length of data
	append(buf, data)              // data

	return buf
}

func (cc *controlChannel) registerChannel(name string, id int) {

	out := []byte(id)
	append(out, []byte(name))

	out = cc.formatMessage(REGISTER, out)

	resp := cc.sendRecieve(out)
}

func (cc *controlChannel) sendRecieve(data []byte) []byte {
	cc.mutex.Lock()

	cc.channel.Write(data)

	buf := make([]byte, 1024)

	out := make([]byte, 1)
	n, err := cc.channel.Read(buf)

	// first

	cc.mutex.Unlock()

	return out
}

func (cc *controlChannel) incomingMessages() {

}

func (cc *controlChannel) handleRequest(mtype controlMessageType, data []byte) {

}

/*func (m *Multiplex) GobChannel(name string) *Channel {
	channel := NewGobChannel(name)

	m.router.RegisterChannel(channel)
	return channel
}*/
