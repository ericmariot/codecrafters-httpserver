package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
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
	defer conn.Close()

	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
	}

	request := string(data[:n])
	lines := strings.Split(request, "\r\n")
	firstLine := strings.Fields(lines[0])

	HTTP_HEADER_OK := "HTTP/1.1 200 OK\r\n"
	HTTP_HEADER_NOT_FOUND := "HTTP/1.1 404 Not Found\r\n"
	HTTP_CONTENT_TYPE_TEXT := "Content-Type: text/plain\r\n"
	HTTP_CONTENT_LENGTH := "Content-Length: "
	HTTP_SEPARATE_LINE := "\r\n\r\n"
	HTTP_END := "\r\n"

	if strings.Contains(firstLine[1], "/echo/") {
		response := strings.TrimPrefix(firstLine[1], "/echo/")
		contentLength := strconv.Itoa(len((response)))
		conn.Write([]byte(
			HTTP_HEADER_OK +
				HTTP_CONTENT_TYPE_TEXT +
				HTTP_CONTENT_LENGTH +
				contentLength +
				HTTP_SEPARATE_LINE +
				response))
	} else if firstLine[1] == "/" {
		conn.Write([]byte(HTTP_HEADER_OK + HTTP_END))
	} else {
		conn.Write([]byte(HTTP_HEADER_NOT_FOUND + HTTP_END))
	}
}
