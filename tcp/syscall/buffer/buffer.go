package buffer

import (
	"bytes"
	"encoding/binary"
	"net"
	"syscall"
	"unsafe"
)

// PsdHeader 伪首部
type PsdHeader struct {
	SrcAddr   uint32
	DstAddr   uint32
	Zero      uint8
	ProtoType uint8
	TCPLength uint16
}

// TCPHeader 首部
type TCPHeader struct {
	SrcPort   uint16
	DstPort   uint16
	SeqNum    uint32
	AckNum    uint32
	Offset    uint8
	Flag      uint8
	Window    uint16
	Checksum  uint16
	UrgentPtr uint16
}

func inetAddr(addr string) uint32 {
	ip := net.ParseIP(addr).To4()
	return uint32(ip[3])<<24 + uint32(ip[2])<<16 + uint32(ip[1])<<8 + uint32(ip[0])
}
func checkSum(data []byte) uint16 {
	var (
		sum    uint32
		index  int
		length = len(data)
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}

// Create 创建一个buffer对象
func Create(dst string, len int) (buf bytes.Buffer) {
	psdheader := PsdHeader{
		SrcAddr:   inetAddr("127.0.0.1"),
		DstAddr:   inetAddr(dst),
		Zero:      0,
		ProtoType: syscall.IPPROTO_TCP,
		TCPLength: uint16(unsafe.Sizeof(TCPHeader{})) + uint16(len),
	}

	tcpheader := TCPHeader{
		SrcPort:  50000,
		DstPort:  5000,
		SeqNum:   0,
		AckNum:   0,
		Offset:   uint8(uint16(unsafe.Sizeof(TCPHeader{}))/4) << 4,
		Flag:     2,
		Window:   60000,
		Checksum: 0,
	}

	binary.Write(&buf, binary.BigEndian, psdheader)
	binary.Write(&buf, binary.BigEndian, tcpheader)
	tcpheader.Checksum = checkSum(buf.Bytes())
	buf.Reset()
	binary.Write(&buf, binary.BigEndian, tcpheader)

	return
}
