package main

import (
	"log"
	"net"
	"os"
)

var conn net.Conn

func read() {
	buffer := make([]byte, 1024)
	for {

		n, err := conn.Read(buffer)
		if err != nil {
			// log.Println("Read buffer failed:", err)
			return
		}

		log.Println("Received string:", string(buffer[:n]))
	}
}

var server = "localhost:5000"

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if len(os.Args) > 1 {
		server = os.Args[1]
	}

	listener, err := net.Listen("tcp", server)
	if err != nil {
		log.Fatalln("Resolve TCP failed:", err)
	}
	defer listener.Close()

	log.Println("Listening on:", listener.Addr())
	for {
		conn, err = listener.Accept()
		if err != nil {
			continue
		}

		log.Println("Connect with:", conn.RemoteAddr())
		read()
	}
}
