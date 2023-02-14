package tcp

import (
	"fmt"
	"net"
	"sync"
)

func NewConnect(conn net.Conn) *Connect {
	return &Connect{
		Conn: conn,
		id:   fmt.Sprintf("%v", conn.RemoteAddr()),
		ctx:  sync.Map{},
	}
}

// Connect 连接实例
type Connect struct {
	net.Conn
	ctx sync.Map

	id  string
	uid string
}

func (c *Connect) Set(key, value any) {
	c.ctx.Store(key, value)
}

func (c *Connect) Get(key any) (value any, ok bool) {
	return c.ctx.Load(key)
}

func (c *Connect) Del(key any) {
	c.ctx.Delete(key)
}

func (c *Connect) Id() string {
	return c.id
}

func (c *Connect) Write(msg []byte) error {
	_, err := c.Conn.Write(ToFrame(msg))
	return err
}

func (c *Connect) SetUid(uid string) {
	c.uid = uid
}

func (c *Connect) Uid() string {
	return c.uid
}
