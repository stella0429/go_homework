package main

import (
	"bufio"
	"fmt"
	"net"
	"project/internal/pkg/protocol"
)

var (
	address = ":8011"
)

func main() {
	netListen, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer netListen.Close()

	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		fmt.Println(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		receiveStr, err := protocol.Unpack(reader)
		if err != nil {
			fmt.Println("read from connection err:", err)
			break
		}
		fmt.Println("receive data string:", receiveStr)
	}
}
