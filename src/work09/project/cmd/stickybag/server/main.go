package main

import (
	"fmt"
	"net"
)

var (
	address = ":8010"
	bytelen = 20
)

func main() {
	netListen, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer netListen.Close()

	fmt.Println("Waiting for clients...")
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
	buffer := make([]byte, bytelen)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}
		fmt.Println(conn.RemoteAddr().String(), "receive data length:", n)
		fmt.Println(conn.RemoteAddr().String(), "receive data:", buffer[:n])
		fmt.Println(conn.RemoteAddr().String(), "receive data string:", string(buffer[:n]))
	}
}
