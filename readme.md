# tcp 服务实现

#### 数据帧（Data Framing）
````go
type Frame struct {
	// < 126 这个数用来表示传输数据的长度
	// == 126 2个字节表示的是一个16进制无符号数，这个数用来表示传输数据的长度
	// == 127 8个字节表示的一个64位无符合数，这个数用来表示传输数据的长度
	PayloadLen int
	Payload    []byte
}
````

#### 如何使用
server.go
````go
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
	log.Println("OnError")
}

````
client.go
````go
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
	log.Println("OnError")
}

````