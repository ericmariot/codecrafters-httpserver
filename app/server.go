package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const crlf = "\r\n"

type HTTPRequest struct {
	body    string
	method  string
	path    string
	version string
	agent   string
}
type HTTPResponse struct {
	Version       string
	StatusCode    int
	StatusMessage string
}

func (r *HTTPResponse) Build() string {
	return fmt.Sprintf("%s %d %s \r\n\r\n", r.Version, r.StatusCode, r.StatusMessage)
}

func MakeHTTPResponse(statusCode int, statusMessage string) *HTTPResponse {
	return &HTTPResponse{
		Version:       "HTTP/1.1",
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
	}
}

func MakeHTTPRequest(conn *net.Conn) (*HTTPRequest, error) {
	buf := make([]byte, 1024)
	_, err := (*conn).Read(buf)
	if err != nil {
		return nil, err
	}
	str := string(buf)
	parts := strings.Split(str, crlf)
	pathElements := strings.Split(parts[0], " ")
	req := new(HTTPRequest)
	req.body = strings.Replace(parts[len(parts)-1], "\x00", "", -1)
	req.method = pathElements[0]
	req.path = pathElements[1]
	req.version = pathElements[2]
	req.agent = parts[2]
	return req, nil

}

func proccessConnection(conn net.Conn) {
	req, err := MakeHTTPRequest(&conn)
	if err != nil {
		log.Fatal(err)
	}

	if req.path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		conn.Close()
	}

	if strings.HasPrefix(req.path, "/echo/") {
		echo := strings.TrimPrefix(req.path, "/echo/")
		res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(echo), echo)

		conn.Write([]byte(res))
		conn.Close()
	}

	if strings.HasPrefix(req.path, "/user-agent") {
		user := strings.TrimPrefix(req.agent, "User-Agent: ")
		res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(user), user)

		conn.Write([]byte(res))
		conn.Close()
	}

	if req.method == "POST" && strings.HasPrefix(req.path, "/files/") {
		filename := strings.TrimPrefix(req.path, "/files/")
		file := dir + filename
		err := os.WriteFile(file, []byte(req.body), 0644)
		if err != nil {
			fmt.Println("Error writing file: ", err.Error())
			os.Exit(1)
		}

		conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
		conn.Close()
	}

	if strings.HasPrefix(req.path, "/files/") {
		path := fmt.Sprint(dir, "/", strings.Join(strings.Split(req.path, "/files/")[1:], ""))

		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			fmt.Fprint(conn, "HTTP/1.1 404 NOT FOUND\r\n\r\n")
			conn.Close()
			return
		}

		content, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: " + strconv.Itoa(len(content)) + "\r\n\r\n" + string(content)))
		conn.Close()
	}

	conn.Write([]byte(MakeHTTPResponse(404, "Not Found").Build()))
	conn.Close()
}

var dir string

func main() {
	fmt.Println("Logs from your program will appear here!")
	directory := flag.String("directory", ".", "Specify the directory")
	flag.Parse()
	dir = *directory

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go proccessConnection(conn)
	}
}
