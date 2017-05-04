package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/zhengxiaoyao0716/computer-network/webserver"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	signal := flag.String("s", "run", "信号 run | scan")
	host := flag.String("host", "127.0.0.1", "指定主机IP或域名")
	port := flag.Int("port", -1, "指定端口号（-1表示自动检索）")
	flag.Parse()

	switch *signal {
	case "scan":
		fmt.Println("正在扫描本机可用IPv4地址")
		scanIPv4()
		fmt.Println("完成")
		os.Exit(0)
	case "run":
		address := *host + ":"
		if *port == -1 {
			*port = 4000
			for true {
				if checkAddress(address + strconv.Itoa(*port)) {
					break
				}
				*port++
			}
		}
		address += strconv.Itoa(*port)
		if !checkAddress(address) {
			log.Println("检查到连接不可用")
			os.Exit(0)
		}

		run(address)
	default:
		flag.PrintDefaults()
	}
}
func scanIPv4() {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatalln(err)
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatalln(err)
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil && ipnet.IP.IsGlobalUnicast() {
					fmt.Println(iface.Name, ":", ipnet.IP)
				}
			}
		}
	}
}
func checkAddress(address string) bool {
	_, err := http.Get("http://" + address)
	if err == nil || !strings.Contains(err.Error(), "refused") {
		return false
	}
	return true
}

func run(address string) {
	s := webserver.New(address)

	s.Get("/", func(req *webserver.Req) *webserver.Resp {
		log.Println("/", req)
		return nil
	})
	s.Get("/index.html", func(req *webserver.Req) *webserver.Resp {
		log.Println("/index.html", req)
		return nil
	})

	if err := s.Run(); err != nil {
		log.Fatalln("Run server failed:", err)
	}
}
