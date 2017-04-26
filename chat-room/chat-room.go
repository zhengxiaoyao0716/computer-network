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
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	signal := flag.String("s", "run", "信号 run | scan")
	host := flag.String("host", "localhost", "指定主机IP或域名")
	port := flag.Int("port", -1, "指定端口号（-1表示自动检索）")
	room := flag.String("room", "http://localhost:4001", "指定聊天室通讯地址\n"+
		"    \t * 当聊天室地址与本机服务地址相同时\n"+
		"    \t * (room == http://address:port)\n"+
		"    \t * 应用将作为聊天室启动\n"+
		"    \t\b",
	)
	flag.Parse()

	switch *signal {
	case "scan":
		fmt.Println("正在扫描本机可用IPv4地址")
		scanIPv4()
		fmt.Println("完成")
		os.Exit(0)
	case "run":
		address := "http://" + *host + ":"
		if *port == -1 {
			*port = 4001
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
		if !run(address, *room) {
			log.Println("初始化失败")
			os.Exit(0)
		}
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
	_, err := http.Get(address)
	if err == nil || !strings.Contains(err.Error(), "refused") {
		return false
	}
	return true
}

func run(address string, room string) bool {
	if address != room {
		log.Println("Join room：", room)
	} else {
		log.Println("Create room：", room)
	}
	return true
}
