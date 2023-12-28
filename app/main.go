package main

import (
	"fmt"
	"net"
	"os"
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
	handleConn(conn)
}

func handleConn(conn net.Conn) {
	req := make([]byte, 1000)
	n, err := conn.Read(req)
	if err != nil {
		fmt.Printf("failed to read from connection. Err: %v\n", err)
		os.Exit(1)
	}
	req = req[:n]
	lines := strings.Split(string(req), "\r\n")
	tokens := strings.Split(lines[0], " ")
	path := tokens[1]
	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

	conn.Close()
}
