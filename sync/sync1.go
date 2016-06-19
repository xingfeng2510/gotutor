package main

import (
	"bytes"
	"log"
	"os"
	"sync"
	"time"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func main0() {
	b := bufPool.Get().(*bytes.Buffer)
	b.Reset()
	b.WriteString(time.Now().UTC().Format(time.RFC3339))
	b.WriteByte(' ')
	b.WriteString("path=/search?q=flowers")
	b.WriteByte('\n')
	os.Stdout.Write(b.Bytes())
	bufPool.Put(b)
}

func testCondSignal() {
	var m sync.Mutex
	c := sync.NewCond(&m)
	n := 20
	running := make(chan bool, n)
	awake := make(chan bool, n)
	for i := 0; i < n; i++ {
		go func(g int) {
			m.Lock()
			running <- true
			c.Wait()
			awake <- true
			m.Unlock()
		}(i)
	}
	for i := 0; i < n; i++ {
		<-running
	}
	for n > 0 {
		select {
		case <-awake:
			log.Fatal("goroutine not asleep")
		default:
		}
		m.Lock()
		c.Signal()
		m.Unlock()
		<-awake
		select {
		case <-awake:
			log.Fatal("too many goroutines awake")
		default:
		}
		n--
	}
	c.Signal()
}

func testCondBroadcast() {
	var m sync.Mutex
	c := sync.NewCond(&m)
	n := 200
	running := make(chan int, n)
	awake := make(chan int, n)
	exit := false
	for i := 0; i < n; i++ {
		go func(g int) {
			m.Lock()
			for !exit {
				running <- g
				c.Wait()
				awake <- g
			}
			m.Unlock()
		}(i)
	}
	for i := 0; i < n; i++ {
		for i := 0; i < n; i++ {
			<-running
		}
		if i == n-1 {
			m.Lock()
			exit = true
			m.Unlock()
		}
		select {
		case <-awake:
			log.Fatal("goroutine not asleep")
		default:
		}
		m.Lock()
		c.Broadcast()
		m.Unlock()
		seen := make([]bool, n)
		for i := 0; i < n; i++ {
			g := <-awake
			if seen[g] {
				log.Fatal("goroutine woke up twice")
			}
			seen[g] = true
		}
	}
}

func main() {
	testCondSignal()
	testCondBroadcast()
}
