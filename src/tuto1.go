package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unsafe"
)

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func fib() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

func filter(s []int, f func(int) bool) []int {
	var p []int
	for _, v := range s {
		if f(v) {
			p = append(p, v)
		}
	}
	return p
}

type Datas struct {
	c0 byte
	c1 int
	c2 string
	c3 int
}

type A struct {
	x int32
	y int64
}

type MyType struct {
	i    int
	name string
}

func (mt *MyType) SetI(i int) {
	mt.i = i
}

func (mt *MyType) SetName(name string) {
	mt.name = name
}

func (mt *MyType) String() string {
	return fmt.Sprintf("%p", mt) + "--name:" + mt.name + " i:" + strconv.Itoa(mt.i)
}

func main() {
	ch := make(chan bool, 0)
	go func() {
		fmt.Println("before")
		<-ch
		fmt.Println("after")
	}()

	time.Sleep(time.Second * 1)
	close(ch)
	time.Sleep(time.Second * 1)

	myType := &MyType{22, "helo"}
	mtt := reflect.TypeOf(myType)
	nm := mtt.NumMethod()
	for i := 0; i < nm; i++ {
		fmt.Printf("method %d: %s %s\n", i, mtt.Method(i).Name, mtt.Method(i).Type)
	}

	var x float64 = 3.4
	fmt.Println("type: ", reflect.TypeOf(x))

	f := func(r rune) bool {
		return !unicode.IsNumber(r) && !unicode.IsLetter(r)
	}
	fmt.Printf("Fields are %v", strings.FieldsFunc("  foo1;bar2,baz3...", f))

	f2 := fib()
	fmt.Println(f2(), f2(), f2())

	fmt.Println(filter([]int{1, 2, 5, 7, 8, 9}, func(a int) bool {
		return a%2 == 0
	}))

	var str string = "hello"
	p := (*struct {
		str uintptr
		len int
	})(unsafe.Pointer(&str))

	fmt.Printf("%+v\n", p)

	var slice []int32 = make([]int32, 5, 10)
	p2 := (*struct {
		array uintptr
		len   int
		cap   int
	})(unsafe.Pointer(&slice))

	fmt.Printf("output: %+v\n", p2)

	var iface interface{} = "Hello World!"
	p3 := (*struct {
		tab  uintptr
		data uintptr
	})(unsafe.Pointer(&iface))

	fmt.Printf("%+v\n", p3)

	var ia int = 1
	sa := []int{1, 2, 4}
	fmt.Println(unsafe.Sizeof(ia), unsafe.Sizeof(str), unsafe.Sizeof(sa))

	var d Datas
	fmt.Println(unsafe.Offsetof(d.c0)) // 0
	fmt.Println(unsafe.Offsetof(d.c1)) // 8
	fmt.Println(unsafe.Offsetof(d.c2)) // 16
	fmt.Println(unsafe.Offsetof(d.c3)) // 32

	fmt.Println(unsafe.Alignof(d.c0))
	fmt.Println(unsafe.Alignof(d.c1))
	fmt.Println(unsafe.Alignof(d.c2))
	fmt.Println(unsafe.Alignof(d.c3))

	pa := &A{}
	pp := unsafe.Pointer(pa)
	offset := unsafe.Offsetof(pa.y)
	var px *int32 = (*int32)(pp)
	*px = 32
	var py *int64 = (*int64)(unsafe.Pointer(uintptr(pp) + offset))
	*py = 64
	fmt.Println(pa.x, pa.y)
}
