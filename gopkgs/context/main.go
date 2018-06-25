package main

import (
	"fmt"
	"runtime"

	"golang.org/x/net/context"
)

func doSomething(ctx context.Context) {
	select {
	case <-ctx.Done():

	}
}

func main() {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	// defer cancel()
	fmt.Println(runtime.GOMAXPROCS(8))
	fmt.Println(runtime.GOMAXPROCS(2))
	fmt.Println(runtime.GOMAXPROCS(4))
	fmt.Println(runtime.NumCPU())
}
