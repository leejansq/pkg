package netsq

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"
)

type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func Ping(ip string, timeout int) (bool, error) {

	var (
		icmp     ICMP
		laddr    = net.IPAddr{IP: net.ParseIP("0.0.0.0")}
		raddr, _ = net.ResolveIPAddr("ip", ip)
	)

	icmp.Type = 8
	icmp.Code = 0
	icmp.Checksum = 0
	icmp.Identifier = 0
	icmp.SequenceNum = 0

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.Checksum = CheckSum(buffer.Bytes())
	buffer.Reset()
	binary.Write(&buffer, binary.BigEndian, icmp)
	//fmt.Printf("%#v", syscallutils.SysInfo())

	conn, err := net.DialIP("ip4:icmp", &laddr, raddr)

	if err != nil {
		return false, err
	}

	defer conn.Close()

	if _, err := conn.Write(buffer.Bytes()); err != nil {
		return false, err
	}

	//t_start := time.Now()

	conn.SetReadDeadline((time.Now().Add(time.Second * time.Duration(timeout))))

	recv := make([]byte, 1024)
	_, err = conn.Read(recv)

	if err != nil {
		return false, err
	}

	//t_end := time.Now()

	//dur := t_end.Sub(t_start).Nanoseconds() / 1e6

	//fmt.Printf("来自 %s 的回复: 时间 = %dms\n", raddr.String(), dur)
	//fmt.Println("connected ok!")
	//time.Sleep(5 * time.Second)

	return true, nil
}

func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
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
