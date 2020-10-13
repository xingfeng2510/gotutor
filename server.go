package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	ln, err := net.Listen("tcp4", ":6000")
	if err != nil {
		fmt.Printf("listen error: %v\n", err)
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("accept error: %v\n", err)
			continue
		}

		go handleConnection(conn.(*net.TCPConn))
	}
}

func handleConnection(c *net.TCPConn) {
	defer c.Close()
	buf := make([]byte, 4096)

	c.SetKeepAlive(false)

	for {
		n, err := c.Read(buf)
		if err != nil || n == 0 {
			fmt.Printf("read error: %v\n", err)
			break
		}
		n, err = c.Write(buf[0:n])
		if err != nil {
			fmt.Printf("write error: %v\n", err)
			break
		}
	}
	fmt.Printf("Connection from %v closed.\n", c.RemoteAddr())
}
