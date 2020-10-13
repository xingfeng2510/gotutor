package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", ":6000")
	if err != nil {
		log.Printf("dail error: %v\n", err)
		os.Exit(1)
	}
	conn.(*net.TCPConn).SetKeepAlive(false)

	buf := make([]byte, 4096)

	fmt.Fprintf(conn, "hello")
	if _, err = bufio.NewReader(conn).Read(buf); err != nil {
		log.Printf("read error: %v\n", err)
	} else {
		log.Printf("read: %s\n", buf)
	}

	go func() {
		for {
			fmt.Fprintf(conn, "hello2")
			if _, err = bufio.NewReader(conn).Read(buf); err != nil {
				log.Printf("read error: %v\n", err)
				break
			} else {
				log.Printf("read: %s\n", buf)
			}

			time.Sleep(10 * time.Second)
		}
	}()

	time.Sleep(2 * time.Minute)

	fmt.Fprintf(conn, "hello3")
	if _, err = bufio.NewReader(conn).Read(buf); err != nil {
		log.Printf("read error: %v\n", err)
	} else {
		log.Printf("read: %s\n", buf)
	}

	time.Sleep(2 * time.Minute)
}
