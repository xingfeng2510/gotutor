package main

import (
	"container/ring"
	"fmt"
)

func main() {
	r := ring.New(10)
	p := r
	for i := 0; i < r.Len(); i++ {
		p.Value = i
		p = p.Next()
	}

	r.Do(func(v interface{}) {
		fmt.Println(v)
	})
	fmt.Println()

	p = p.Move(6)
	fmt.Println(p.Value)
	fmt.Println()

	s := r.Unlink(18)
	s.Do(func(v interface{}) {
		fmt.Println(v)
	})
	fmt.Println()

	r.Do(func(v interface{}) {
		fmt.Println(v)
	})
	fmt.Println()

	r2 := r.Link(s)
	fmt.Println(r2.Value)
	fmt.Println()

	r.Do(func(v interface{}) {
		fmt.Println(v)
	})
}
