package main

import (
	"fmt"
	"net"
	"project/internal/pkg/protocol"
	"time"
)

var (
	address = "127.0.0.1:8011"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("connect success...")

	go send(conn)

	for {
		time.Sleep(1 * 1e9)
	}
}

func send(conn net.Conn) {
	for i := 0; i < 10; i++ {
		message := "Hi,go!哈哈哈哈!哈哈哈!"
		data, err := protocol.Pack([]byte(message))
		if err != nil {
			panic(err)
		}
		conn.Write(data)
	}
}
