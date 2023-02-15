package tcp

import (
	"net"
)

type client struct {
	*net.TCPConn
	tcpAddr *net.TCPAddr
	event   Event
}

// NewClient service := ":7777"
func NewClient(service string) *client {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		panic(err)
	}
	return &client{
		tcpAddr: tcpAddr,
	}
}

func (c *client) SetEvent(e Event) {
	c.event = e
}

func (c *client) Run() {
	var err error
	c.TCPConn, err = net.DialTCP("tcp", nil, c.tcpAddr)

	if err != nil {
		go c.event.OnError(err)
		return
	}
	defer c.Close()
	conn := NewConnect(c.TCPConn)
	go c.event.OnConnect(conn)
	defer c.event.OnClose(conn)

	p := Protocol{}
	p.Init()
	for {
		message, err := p.Read(conn.Conn)
		if err != nil {
			c.event.OnError(err)
			return
		}

		c.event.OnMessage(conn, message)
	}
}
