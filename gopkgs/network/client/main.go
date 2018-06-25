package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func main() {
	sip := net.ParseIP("127.0.0.1")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: sip, Port: 9981}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		log.Fatal("error: ", err)
	}
	defer conn.Close()
	s := strings.Repeat("a", 1743)
	conn.Write([]byte(s))
	fmt.Printf("<%s>\n", conn.RemoteAddr())

	time.Sleep(time.Second * 10)
}
