package main

import (
	"github.com/ctfang/tcp-server"
	"log"
	"time"
)

func main() {
	server := tcp.NewServer(":7777")
	server.SetEvent(&e{s: server})
	server.Run()
}

type e struct {
	s *tcp.Server
}

func (e *e) OnConnect(connect *tcp.Connect) {
	go func() {
		for true {
			e.s.Range(func(key, value any) bool {
				id := key.(string)
				connect.Write([]byte("client = " + id))
				return true
			})
			time.Sleep(1 * time.Second)
		}
	}()
}

func (*e) OnMessage(connect *tcp.Connect, frame tcp.Frame) {
	log.Println(string(frame.Payload))
}

func (*e) OnClose(connect *tcp.Connect) {
	log.Println("OnClose")
}

func (*e) OnError(err error) {
	log.Println("OnError", err)
}
