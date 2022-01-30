package main

import (
	"bufio"
	"fmt"
	week9 "leiax00.cn/week9"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("连接服务端出错,错误原因: ", err)
		return
	}
	defer conn.Close()
	fmt.Println("与服务端连接建立成功...")
	fmt.Printf("请输入指令: [0,10]")
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		_, _ = conn.Write([]byte(text))
		readResp(conn)
	}
}

func readResp(conn net.Conn) {
	//var bytes = make([]byte, 1024)
	//_, err := conn.Read(bytes)
	//fmt.Println(err)
	//fmt.Println(string(bytes))
	rr := bufio.NewReader(conn)
	decode := week9.GoImD.WithStream(rr).Decode()
	result, err := decode.Result()
	if err != nil {
		fmt.Printf("err happen: %v", err)
	}
	fmt.Println(result.Pretty())
}
