package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var statusCodeNameMap = map[int]string{
	200: "OK",
	404: "Not Found",
}

func writeHtmlResponseSimple(conn net.Conn, statusCode int) {
	w := bufio.NewWriter(conn)
	defer w.Flush()
	if headerName, ok := statusCodeNameMap[statusCode]; ok {
		w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n", statusCode, headerName)))
	}
}

func writeHtmlResponseWithPlainBody(conn net.Conn, statusCode int, body string) {
	w := bufio.NewWriter(conn)
	defer w.Flush()
	if headerName, ok := statusCodeNameMap[statusCode]; ok {
		w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", statusCode, headerName, len(body), body)))
	}
}

func writeHtmlResponseWithFile(conn net.Conn, statusCode int, filePath string) {
	w := bufio.NewWriter(conn)
	defer w.Flush()
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		w.WriteString("file not present")
	}
	if headerName, ok := statusCodeNameMap[statusCode]; ok {
		w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n", statusCode, headerName, len(bytes))))
	}
	w.Write(bytes)
}
