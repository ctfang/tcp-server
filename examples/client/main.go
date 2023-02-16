package main

import (
	"github.com/ctfang/tcp-server"
	"log"
)

func main() {
	server := tcp.NewClient("127.0.0.1:7777")
	server.SetEvent(&e{})
	server.Run()
}

type e struct{}

func (*e) OnConnect(connect *tcp.Connect) {
	log.Println("ok")

	connect.Write([]byte("form client message"))
}

func (*e) OnMessage(connect *tcp.Connect, msg tcp.Frame) {
	log.Println("OnMessage", string(msg.Payload))
}

func (*e) OnClose(connect *tcp.Connect) {
	log.Println("OnClose")
}

func (*e) OnError(err error) {
	log.Println("OnError", err)
}
