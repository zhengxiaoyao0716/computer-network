package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var server = "localhost:5000"
var message = "Hello world"

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if len(os.Args) > 1 {
		server = os.Args[1]
	}
	if len(os.Args) > 2 {
		message = strings.Join(os.Args[2:], " ")
	}

	// 建立连接

	addr, err := net.ResolveTCPAddr("tcp", server)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := net.Dial("ip", addr.IP.String())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// 构造数据

	buf := gopacket.NewSerializeBuffer()
	tcpLayer := layers.TCP{
		SrcPort: 50000,
		DstPort: layers.TCPPort(addr.Port),
	}
	tcpLayer.SetNetworkLayerForChecksum(&layers.IPv4{
		SrcIP:    net.ParseIP("127.0.0.1").To4(),
		DstIP:    addr.IP,
		Protocol: layers.IPProtocolTCP,
	})
	if err = gopacket.SerializeLayers(buf, gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}, &tcpLayer); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Send to %s:%d [% X]\n", addr.IP, addr.Port, buf.Bytes())

	// 写入数据

	if _, err = conn.Write(append(buf.Bytes(), message...)); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Send message:", message)
}
