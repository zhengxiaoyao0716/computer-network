package webserver

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"net/url"

	"bufio"

	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Server 服务端
type Server struct {
	Addr      string
	handleMap map[string]Handle
	staticMap map[string]string
}

// New 创建
func New(address string) *Server {
	return &Server{
		Addr:      address,
		handleMap: map[string]Handle{},
		staticMap: map[string]string{},
	}
}

// Run 开始运行 deprecated
func (s *Server) _Run() error {
	addr, err := net.ResolveTCPAddr("tcp", s.Addr)
	if err != nil {
		return err
	}

	conn, err := net.ListenIP("ip", &net.IPAddr{IP: addr.IP, Zone: addr.Zone})
	if err != nil {
		return err
	}
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		n, from, err := conn.ReadFrom(buf)
		if err != nil {
			// log.Fatalln("Read from connection failed:", err)
		}
		buf = buf[:n]
		// fmt.Printf("Received buffer: % X\n", buf)

		go (func() {
			packet := gopacket.NewPacket(buf, layers.LayerTypeTCP, gopacket.Default)
			if layer := packet.Layer(layers.LayerTypeTCP); layer != nil {
				if tcp := layer.(*layers.TCP); int(tcp.DstPort) == addr.Port {
					fmt.Printf("Read from %s:%d [% X]\n", from, tcp.SrcPort, tcp.Contents)
					fmt.Println("Read message:", string(tcp.Payload))
					fmt.Printf("%d %d %t %t\n\n", tcp.Seq, tcp.Ack, tcp.SYN, tcp.ACK)

					// buf := gopacket.NewSerializeBuffer()
					// tcp.DstPort, tcp.SrcPort = tcp.SrcPort, tcp.DstPort
					// tcp.Ack = tcp.Seq + 1
					// tcp.Seq = rand.Uint32()
					// tcp.ACK = true
					// tcp.SetNetworkLayerForChecksum(&layers.IPv4{
					// 	SrcIP:    addr.IP,
					// 	DstIP:    net.ParseIP(from.String()).To4(),
					// 	Protocol: layers.IPProtocolTCP,
					// })
					// if err = gopacket.SerializeLayers(buf, gopacket.SerializeOptions{
					// 	FixLengths:       true,
					// 	ComputeChecksums: true,
					// }, tcp); err != nil {
					// 	log.Fatalln(err)
					// }
					// if _, err := conn.WriteTo(buf.Bytes(), from); err != nil {
					// 	log.Fatalln("Weite to connection failed:", err)
					// }
				}
			}
		})()
	}

	return nil
}
func isEol(r rune) bool {
	switch r {
	case '\r', '\n':
		return true
	}
	return false
}

// Run 开始运行
func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	log.Println("Listen tcp address:", s.Addr)

	for {
		conn, _ := listener.Accept()
		go func(c net.Conn) {
			defer c.Close()

			// 读取请求

			c.SetReadDeadline(time.Time{})
			reader := bufio.NewReader(c)
			lines := []string{}
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					log.Println(err)
					return
				}
				line = strings.TrimSpace(line)
				if line == "" {
					break
				}
				lines = append(lines, line)
			}

			req := Req{
				Header: map[string][]string{},
			}

			// 首行

			if subs := strings.Fields(lines[0]); len(subs) == 3 {
				method, path, proto := subs[0], subs[1], subs[2]
				switch proto {
				case "HTTP/1.0", "HTTP/1.1", "HTTP/2":
					break
				default:
					return
				}
				req.Method = method
				if req.URL, err = url.Parse("http://" + s.Addr + path); err != nil {
					return
				}
				req.Proto = proto
			} else {
				return
			}

			// 请求头

			for _, line := range lines[1:] {
				subs := strings.Split(line, ": ")
				key, value := subs[0], subs[1]
				if _, ok := req.Header[key]; !ok {
					req.Header[key] = []string{}
				}
				req.Header[key] = append(req.Header[key], value)
			}

			// 请求体

			if v, ok := req.Header["Content-Length"]; ok {
				l, err := strconv.Atoi(v[0])
				if err != nil {
					log.Println(err)
					return
				}
				req.Content, err = reader.Peek(l)
				if err != nil {
					log.Println(err)
					return
				}
			}

			handle, ok := s.handleMap[req.URL.Path]
			if !ok {
				key := req.URL.Path
				for {
					i := strings.LastIndex(key, "/")
					key = key[:i]
					path, ok := s.staticMap[key+"/"]
					if ok {
						handle = GetStaticHandle(path + req.URL.Path[i:])
						break
					}
				}
			}
			if handle != nil {
				resp := handle(&req)
				lines := []string{
					resp.Proto + " " + strconv.Itoa(resp.StatusCode) + " " + resp.StatusText,
				}
				for key, value := range resp.Header {
					for _, line := range value {
						lines = append(lines, key+": "+line)
					}
				}
				_, err := c.Write(append([]byte(strings.Join(lines, "\r\n")+"\r\n\r\n"), resp.Content...))
				if err != nil {
					log.Println(err)
					return
				}
			}
		}(conn)
	}

	return nil
}
