package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"syscall"

	"github.com/zhengxiaoyao0716/computer-network/tcp/syscall/buffer"
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

	if runtime.GOOS == "windows" {
		var data syscall.WSAData
		if err := syscall.WSAStartup((2<<8)+2, &data); err != nil {
			log.Fatalln("WSAStartup failed:", err)
		}
		defer syscall.WSACleanup()
	}

	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_IP)
	if err != nil {
		log.Fatalln("Create socket failed:", err)
	}
	defer syscall.Closesocket(sock)

	// 构造数据

	buf := buffer.Create(addr.IP.String(), len(message))

	fmt.Printf("Send to %s:%d [% X]\n", addr.IP, addr.Port, buf.Bytes())
	_, err = buf.WriteString(message)
	if err != nil {
		log.Fatalln(err)
	}

	var sockAddr syscall.SockaddrInet4
	sockAddr.Port = addr.Port
	copy(sockAddr.Addr[:], addr.IP.To4())

	// 发送数据

	switch runtime.GOOS {
	case "windows":
		if err := syscall.Connect(sock, &sockAddr); err != nil {
			log.Fatalln("Connect socket failed:", err)
		}
		var (
			sent       uint32
			overlapped syscall.Overlapped
		)
		err = syscall.WSASend(sock, &syscall.WSABuf{
			Len: uint32(buf.Len()),
			Buf: &buf.Bytes()[0],
		}, 1, &sent, 0, &overlapped, nil)
	default:
		err = syscall.Sendto(sock, buf.Bytes(), 0, &sockAddr)
	}
	if err != nil {
		log.Fatalln("Send message failed: ", err)
	}
	fmt.Println("Send message:", message)
}
