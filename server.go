package tcp

import (
	"github.com/sirupsen/logrus"
	"net"
	"sync"
)

// NewServer service := ":7777"
func NewServer(service string) *Server {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		panic(err)
	}
	return &Server{
		tcpAddr: tcpAddr,
	}
}

type Server struct {
	tcpAddr *net.TCPAddr
	cons    sync.Map
	event   Event
}

func (s *Server) SetEvent(e Event) {
	s.event = e
}

func (s *Server) Run() {
	listener, err := net.ListenTCP("tcp", s.tcpAddr)
	if err != nil {
		logrus.Fatalf("tpc server run err(%v)", err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.Errorf("tcp server accept err %v", err)
			continue
		}

		go s.handleClient(NewConnect(conn))
	}
}

func (s *Server) handleClient(conn *Connect) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Error("handleClient", err)
		}
	}()
	defer conn.Close()
	// 保存客户端连接
	s.cons.Store(conn.Id(), conn)
	defer s.cons.Delete(conn.Id())
	s.event.OnConnect(conn)
	defer s.event.OnClose(conn)

	p := Protocol{}
	p.Init()
	for {
		message, err := p.Read(conn.Conn)
		if err != nil {
			str := err.Error()
			if str == "EOF" {
				return
			}
			s.event.OnError(err)
			return
		}

		s.event.OnMessage(conn, message)
	}
}

// GetConn 获取客户端
func (s *Server) GetConn(id string) (*Connect, bool) {
	client, ok := s.cons.Load(id)
	if !ok {
		return nil, false
	}
	return client.(*Connect), ok
}

// Range 循环所有客户端
func (s *Server) Range(f func(key, value any) bool) {
	s.cons.Range(f)
}
