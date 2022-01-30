package main

import (
	"fmt"
	"github.com/Terry-Mao/goim/api/protocol"
	"github.com/Terry-Mao/goim/pkg/bytes"
	"net"
	"strings"
	"time"
)

func main() {
	socket, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("开启监听失败,错误原因: ", err)
		return
	}
	defer socket.Close()
	fmt.Println("开启监听...")

	for {
		conn, err := socket.Accept()
		if err != nil {
			fmt.Println("建立链接失败,错误原因: ", err)
			return
		}
		defer conn.Close()
		fmt.Println("建立链接成功,客户端地址是: ", conn.RemoteAddr())
		logic(conn)

	}
}

func logic(conn net.Conn) {
	for {
		tmp := make([]byte, 3)
		_, _ = conn.Read(tmp) // 演示, 忽略异常
		fmt.Println(string(tmp))
		cmd := strings.TrimSpace(string(tmp))
		var writer *bytes.Writer
		switch cmd {
		case "1":
			writer = reply(1, 0, 1, "hand shake....")
		case "2":
			writer = reply(1, 1, 2, "hand shake reply")
		default:
			writer = reply(1, 4, 4, "send message")
		}
		_, _ = conn.Write(writer.Buffer())
		time.Sleep(time.Second * 2)
	}
}

func reply(ver, op, seq int32, body string) *bytes.Writer {
	bs := []byte(body)
	writer := bytes.NewWriterSize(len(body) + 64)
	proto := &protocol.Proto{
		Ver:  ver,
		Op:   op,
		Seq:  seq,
		Body: bs,
	}
	proto.WriteTo(writer)
	return writer
}
