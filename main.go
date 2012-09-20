package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

var (
	request     = []byte("<policy-file-request/>\x00")
	response    = []byte("<cross-domain-policy><allow-access-from domain=\"*\" to-ports=\"*\" /></cross-domain-policy>\x00")
	requestLen  = len(request)
	responseLen = len(response)
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:843")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer listener.Close()

	fmt.Println("=== Flash Socket Policy File Server ===")

	for {
		if conn, err := listener.Accept(); err == nil {
			go loop(conn)
		}
	}
}

func loop(conn net.Conn) {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	buff := make([]byte, requestLen)

	if n, err := io.ReadFull(conn, buff); err == nil && n == requestLen {
		conn.Write(response)
	}
}
