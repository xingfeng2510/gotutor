package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now().UTC().Format(time.RFC3339Nano))
	t, err := time.Parse(time.RFC3339Nano, "2016-09-26T21:58:28+08:00")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t)
	n := int64(1.474965488991146e+15) * 1000
	fmt.Println(n)
	t = time.Unix(n/1e9, n%1e9)
	fmt.Println(t.UTC().Format(time.RFC3339Nano))
}
