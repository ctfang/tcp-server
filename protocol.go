package tcp

import (
	"encoding/binary"
	"errors"
	"net"
)

type Protocol struct {
	// 本地缓冲区
	cacheByte []byte
	// 缓冲长度
	cacheCount int
}

func (w *Protocol) Init() {

}

func (w *Protocol) Read(conn net.Conn) (Frame, error) {
	got := Frame{}
	payloadLenByte, err := w.readConnOrCache(conn, 1)
	if len(payloadLenByte) == 0 {
		return got, errors.New("EOF")
	}
	got.PayloadLen = int(payloadLenByte[0] & 0x7F) // 有效负载, 这里只读取7位
	// TODO 预留一个第8位, 如果有效负载长度大于127字节，则第8位为1
	// mask := payloadLenByte[0] >> 7

	switch got.PayloadLen {
	case 126: // 两个字节表示的是一个16进制无符号数，这个数用来表示传输数据的长度
		temLen, _ := w.readConnOrCache(conn, 2)
		got.PayloadLen = int(binary.BigEndian.Uint16(temLen))
	case 127: // 8个字节表示的一个64位无符合数，这个数用来表示传输数据的长度
		temLen, _ := w.readConnOrCache(conn, 8)
		got.PayloadLen = int(binary.BigEndian.Uint64(temLen))
	}

	got.Payload, err = w.readConnOrCache(conn, got.PayloadLen)
	return got, err
}

func ToFrame(msg []byte) []byte {
	length := len(msg)
	sendByte := make([]byte, 0)

	var payLenByte byte
	switch {
	case length <= 125:
		payLenByte = byte(0x80) | byte(length)
		sendByte = append(sendByte, payLenByte)
	case length <= 65536:
		payLenByte = byte(0x80) | byte(0x7e)
		sendByte = append(sendByte, payLenByte)
		// 随后的两个字节表示的是一个16进制无符号数，用来表示传输数据的长度
		payLenByte2 := make([]byte, 2)
		binary.BigEndian.PutUint16(payLenByte2, uint16(length))
		sendByte = append(sendByte, payLenByte2...)
	default:
		payLenByte = byte(0x80) | byte(0x7f)
		sendByte = append(sendByte, payLenByte)
		// 随后的是8个字节表示的一个64位无符合数，这个数用来表示传输数据的长度
		payLenByte8 := make([]byte, 8)
		binary.BigEndian.PutUint64(payLenByte8, uint64(length))
		sendByte = append(sendByte, payLenByte8...)
	}
	sendByte = append(sendByte, msg...)
	return sendByte
}

// 读取指定长度数据
func (w *Protocol) readConnOrCache(conn net.Conn, count int) ([]byte, error) {
	if w.cacheCount > 0 {
		// 拥有缓冲数据
		if count <= w.cacheCount {
			// 缓冲数据比需要的还要大，直接拿取
			msg := w.cacheByte[:count]
			w.cacheCount = w.cacheCount - count
			w.cacheByte = w.cacheByte[count:]
			return msg, nil
		} else {
			// 缓冲数据不足，剩余需要的位数，多读取一点，可以优化速度
			data := make([]byte, count+512)
			cacheCount, err := conn.Read(data)
			if err != nil {
				return nil, errors.New("读取数据失败")
			}
			w.cacheCount = w.cacheCount + cacheCount
			w.cacheByte = append(w.cacheByte, data[:cacheCount]...)
			return w.readConnOrCache(conn, count)
		}
	} else {
		// 缓冲是空的
		data := make([]byte, 1024)
		cacheCount, err := conn.Read(data)
		if err != nil {
			return nil, errors.New("读取数据失败")
		}
		w.cacheCount = cacheCount
		w.cacheByte = append(w.cacheByte, data[:cacheCount]...)
		return w.readConnOrCache(conn, count)
	}
}
