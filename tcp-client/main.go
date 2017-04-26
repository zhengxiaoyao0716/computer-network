package main

import (
	"log"
	"net"
	"os"
)

var conn net.Conn

func send(message string) {
	conn.Write([]byte(message))
	log.Println("Send message:", message)
}

var server = "localhost:5000"

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if len(os.Args) > 1 {
		server = os.Args[1]
	}

	addr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		log.Fatalln("Resolve TCP failed:", err)
	}

	conn, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatalln("TCP connect failed:", err)
	}

	send("Hello world.")
}
