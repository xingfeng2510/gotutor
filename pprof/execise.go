package main

import (
	"container/list"
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	"net/http"
	_ "net/http/pprof"
)

func Append(slice []int, elements ...int) []int {
	n := len(slice)
	total := len(slice) + len(elements)
	if total > cap(slice) {
		newSlice := make([]int, total, 2*total+1)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[:total]
	copy(slice[n:], elements)
	return slice
}

func Insert(slice []int, index int, value int) []int {
	slice = slice[0 : len(slice)+1]
	copy(slice[index+1:], slice[index:])
	slice[index] = value
	return slice
}

type List2 struct {
	*list.List
	ref     *Ref2
	payload []byte
}

func New(ref *Ref2) *List2 {
	return &List2{
		List:    list.New(),
		ref:     ref,
		payload: make([]byte, 100*1024*1024),
	}
}

type Ref2 struct {
	a int
}

func foo(ref *Ref2) {
	l := New(ref)
	e4 := l.PushBack(4)
	e1 := l.PushFront(1)
	l.InsertBefore(3, e4)
	l.InsertAfter(2, e1)

	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Print(e.Value)
	}
	fmt.Println(l.ref.a)
}

func main() {
	flag.Parse()

	go func() {
		log.Println(http.ListenAndServe("localhost:6666", nil))
	}()

	runtime.GOMAXPROCS(runtime.NumCPU())

	ref := &Ref2{a: 11}
	for i := 0; i < 1000; i++ {
		foo(ref)
		time.Sleep(time.Millisecond * 50)
	}

	time.Sleep(time.Second * 10)

	// slice1 := []int{0, 1, 2, 3, 4}
	// slice2 := []int{55, 66, 77}
	// fmt.Println(slice1)
	// slice1 = Append(slice1, slice2...)
	// fmt.Println(slice1)

	// slice := make([]int, 10, 20)
	// for i := range slice {
	// 	slice[i] = i
	// }
	// fmt.Println(slice)

	// Insert(slice, 5, 99)
	// fmt.Println(slice)
}
