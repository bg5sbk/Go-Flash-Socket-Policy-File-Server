package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"
)

var (
	addr     = flag.String("addr", "0.0.0.0:843", "address")
	file     = flag.String("file", "", "the socket policy file")
	request  = []byte("<policy-file-request/>\x00")
	response = []byte("<cross-domain-policy><allow-access-from domain=\"*\" to-ports=\"*\" /></cross-domain-policy>\x00")
)

func main() {
	flag.Parse()

	var err error
	var lsn net.Listener

	if *file != "" {
		response, err = ioutil.ReadFile(*file)
		if err != nil {
			log.Fatal(err)
		}
		response = append(response, 0)
	}

	lsn, err = net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Flash Socket Policy File Server ===")

	for {
		if conn, err := lsn.Accept(); err == nil {
			go loop(conn)
		}
	}
}

func loop(conn net.Conn) {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	buff := make([]byte, len(request))

	if n, err := io.ReadFull(conn, buff); err == nil && n == len(request) {
		conn.Write(response)
	}
}
