package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var server = "127.0.0.1:5000"

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if len(os.Args) > 1 {
		server = os.Args[1]
	}

	// 建立连接
	addr, err := net.ResolveTCPAddr("tcp", server)
	if err != nil {
		log.Fatalln("Resolve IP address failed:", err)
	}

	conn, err := net.ListenIP("ip", &net.IPAddr{IP: addr.IP, Zone: addr.Zone})
	if err != nil {
		log.Fatalln("Listen IP failed:", err)
	}
	defer conn.Close()

	// 读取数据

	for {
		buf := make([]byte, 1024)
		n, from, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatalln(err)
		}
		buf = buf[:n]

		// 解析数据

		go (func() {
			packet := gopacket.NewPacket(buf, layers.LayerTypeTCP, gopacket.Default)
			if layer := packet.Layer(layers.LayerTypeTCP); layer != nil {
				if tcp := layer.(*layers.TCP); int(tcp.DstPort) == addr.Port {
					fmt.Printf("Read from %s:%d [% X]\n", from, tcp.SrcPort, tcp.Contents)
					fmt.Println("Read message:", string(tcp.Payload))
				}
			}
		})()
	}
}
