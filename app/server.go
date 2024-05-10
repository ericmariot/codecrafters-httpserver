package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	proccessConnection(conn)
}

func proccessConnection(conn net.Conn) {
	buff := make([]byte, 1024)

	n, err := conn.Read(buff)
	if err != nil {
		fmt.Printf("Failed to connect to server: %v", err)
	}

	fmt.Println(string(buff[:n]))

	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
}
