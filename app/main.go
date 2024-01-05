package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	var dirFlag = flag.String("directory", ".", "directory to serve files from")
	flag.Parse()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConn(conn, *dirFlag)
	}
}

func handleConn(conn net.Conn, dir string) {
	defer conn.Close()
	req := make([]byte, 1000)
	n, err := conn.Read(req)
	if err != nil {
		fmt.Printf("failed to read from connection. Err: %v\n", err)
		return
	}
	req = req[:n]
	lines := strings.Split(string(req), "\r\n")
	tokens := strings.Split(lines[0], " ")
	method := tokens[0]
	path := tokens[1]
	if path == "/" {
		writeHtmlResponseSimple(conn, 200)
	} else if strings.Contains(path, "/echo/") {
		resp := path[6:]
		writeHtmlResponseWithPlainBody(conn, 200, resp)
	} else if strings.Contains(path, "/user-agent") {
		for _, header := range lines {
			if strings.Contains(header, "User-Agent: ") {
				headerValue := header[12:]
				writeHtmlResponseWithPlainBody(conn, 200, headerValue)
			}
		}
	} else if strings.Contains(path, "/files/") {
		fileName := path[7:]
		filePath := filepath.Join(dir, fileName)
		if method == "GET" {
			_, err := os.Stat(filePath)
			if err != nil {
				writeHtmlResponseSimple(conn, 404)
				return
			}
			writeHtmlResponseWithFile(conn, 200, filePath)
		} else if method == "POST" {
			file, err := os.Create(filePath)
			if err != nil {
				writeHtmlResponseSimple(conn, 500)
			}
			defer file.Close()
			file.WriteString(lines[len(lines)-1])
			writeHtmlResponseSimple(conn, 201)
		}
	} else {
		writeHtmlResponseSimple(conn, 404)
	}
}
