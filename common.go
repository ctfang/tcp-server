package tcp

type Event interface {
	OnConnect(connect *Connect)
	OnMessage(connect *Connect, frame Frame)
	OnClose(connect *Connect)
	OnError(err error)
}

type Frame struct {
	// < 126 这个数用来表示传输数据的长度
	// == 126 2个字节表示的是一个16进制无符号数，这个数用来表示传输数据的长度
	// == 127 8个字节表示的一个64位无符合数，这个数用来表示传输数据的长度
	PayloadLen int
	// 简单的校验位，例如将Payload中所有数据的二进制值累加，并且将和对255取模，这将产生一个校验和值
	Checksum byte
	Payload  []byte
}
