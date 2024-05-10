package main

import (
	"fmt"
	"net"
	"testing"
)

func TestServerResponse(t *testing.T) {
	go main()

	conn, err := net.Dial("tcp", "localhost:4221")
	if err != nil {
		t.Fatal(err)
	}

	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}

	fmt.Printf("N: %v", n)
	fmt.Println("---")
	fmt.Printf("buff: %v", string(buff[:n]))

	expected := "HTTP/1.1 200 OK\r\n\r\n"
	if string(buff[:n]) != expected {
		t.Errorf("Expected '%s' got '%s'", expected, buff[:n])
	}
}
