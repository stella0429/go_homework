package main

import (
	"fmt"
	"net"
	"time"
)

var (
	address = "127.0.0.1:8010"
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

	fmt.Println("connect success")
	go send(conn)

	for {
		time.Sleep(1 * 1e9)
	}
}

func send(conn net.Conn) {
	for i := 0; i < 100; i++ {
		message := "Hi,go!"
		conn.Write([]byte(message))
	}
}
