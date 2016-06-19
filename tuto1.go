package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/context"
)

/*
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

func m() {
	fmt.Println("before f()")
	f()
	fmt.Println("after f()")
}

func f() {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("runtime panic:", x)
		}
	}()
	g(3)
	fmt.Println("after g(3)")
}

func g(i int) {
	if i > 2 {
		panic(i)
	}
}

func logPanics(handler func(int, int)) func(int, int) {
	return func(a int, b int) {
		defer func() {
			if x := recover(); x != nil {
				fmt.Println(x)
			}
		}()
		handler(a, b)
	}
}

type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	f(w, req)
}

func ArgServer(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, os.Args)
}

func ParallelWrite(data []byte) chan error {
	res := make(chan error, 2)
	f1, err := os.Create("file1")
	if err != nil {
		res <- err
	} else {
		go func() {
			_, err := f1.Write(data)
			res <- err
			f1.Close()
		}()
	}
	f2, err := os.Create("file2")
	if err != nil {
		res <- err
	} else {
		go func() {
			_, err := f2.Write(data)
			res <- err
			f2.Close()
		}()
	}
	return res
}

func main0() {
	c := make(chan bool)
	m := make(map[string]string)
	go func() {
		fmt.Println("aaa")
		m["1"] = "a" // First conflicting access.
		c <- true
	}()
	fmt.Println("bbb")
	m["2"] = "b" // Second conflicting access.
	<-c
	for k, v := range m {
		fmt.Println(k, v)
	}
}

type Iface interface {
	F2()
}

type Iface2 interface {
	F1()
}

type Base struct {
	i Iface2
	a int
}

func (b *Base) F1() {
	fmt.Println("base.a", b.a)
}

func (b *Base) F2() {
	fmt.Println("base.F2")
	b.i.F1()
}

type Derive struct {
	*Base
	b int
}

func (d *Derive) F1() {
	fmt.Println("derive.b", d.b)
}

type AAlias A
type APtrAlias *A

type Point struct{ X, Y float64 }

func (p Point) Distance(q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

type Rocket struct{}

func (r *Rocket) Launch() { fmt.Println("rocket launching...") }

func main2() {
	flag.Parse()

	cancel := make(chan struct{})
	var tick <-chan time.Time
	if flag.Arg(0) == "verbose" {
		tick = time.Tick(time.Second * 2)
	}
	go func() {
		os.Stdin.Read(make([]byte, 1))
		cancel <- struct{}{}
	}()

	select {
	case <-tick:
		fmt.Println("tick...")
	case <-cancel:
		fmt.Println("cancel...")
	}

	filepath.Walk("/tmp/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		if info.IsDir() {
			fmt.Println("dir", path, "skip")
			return filepath.SkipDir
		}
		fmt.Println("filename", path)
		return nil
	})

	var r Rocket
	fmt.Printf("type: %T\n", r.Launch)
	time.AfterFunc(time.Second, r.Launch)

	time.Sleep(time.Second * 2)

	p1 := Point{1, 2}
	q1 := Point{4, 6}

	distanceFromP := p1.Distance
	fmt.Println(distanceFromP(q1))

	var distance func(p, q Point) float64
	distance = Point.Distance
	distance(p1, q1)

	var jsonBlob = []byte(`[
		{"Name": "Platypus", "Order": "Monotremata", "Age": 1234567},
		{"Name": "Quoll",    "Order": "Dasyuromorphia", "Age": 12.34}
	]`)
	type Animal struct {
		Name  string
		Order string
		Age   interface{}
	}
	var animals []Animal
	err := json.Unmarshal(jsonBlob, &animals)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", animals)
	fmt.Printf("type %T, %T\n", animals[0].Age, animals[1].Age)

	dec := json.NewDecoder(bytes.NewReader([]byte(`{"x": 123456789}`)))
	dec.UseNumber()
	var x interface{}
	dec.Decode(&x)
	fmt.Println(x.(map[string]interface{})["x"].(json.Number).Int64())

	mm := map[string]interface{}{
		"user":  "zhangsan",
		"age":   20,
		"score": 89.12,
		"birth": time.Now().UTC().Format(time.RFC3339Nano),
		"profile": map[string]interface{}{
			"addr": "henan",
			"uid":  71461,
		},
	}
	data, _ := json.Marshal(mm)
	s := string(data)
	fmt.Printf("data: %s\n", s)

	var mx map[string]interface{}
	fmt.Println("err", json.Unmarshal(data, &mx))
	fmt.Println("mx", mx)

	ch2 := make(chan int, 0)
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(time.Second * 1)
		timeout <- true
		fmt.Println("timeouting...")
	}()

	go func() {
		time.Sleep(time.Millisecond * 200)
		ch2 <- 11
	}()

	select {
	case v := <-ch2:
		fmt.Println("received value", v)
	case <-timeout:
		fmt.Println("timedout")
	}

	time.Sleep(time.Second * 5)

	pa2 := &A{x: 1, y: 2}
	var pa3 APtrAlias = pa2
	fmt.Println(pa3)

	var pd *Derive
	pd = &Derive{
		Base: &Base{a: 1},
		b:    2,
	}
	pd.Base.i = pd
	var iface Iface = pd
	iface.F2()

	return

	pa := &A{2, 11}
	v1 := reflect.ValueOf(pa)
	v2 := reflect.Indirect(v1)
	fmt.Printf("%v : %v\n", v1.Interface(), v2.Interface())

	a := [][]int{{1, 2, 3}, {4, 5, 6}}

RowLoop:
	for _, row := range a {
		for _, data := range row {
			if data == 2 {
				continue RowLoop
			}
			println(data)
		}
	}

	ch := make(chan int, 10)
	go func(ch chan int) {
		for i := 0; i < 100; i++ {
			ch <- i
		}
	}(ch)
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case i := <-ch:
			fmt.Println(i)
		case t := <-ticker.C:
			fmt.Println(t)
		}
		time.Sleep(200 * time.Millisecond)
	}

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			fmt.Println("xxx", r)
		}
	}()

	var p *int
	*p = 0

	t := time.Tick(1 * time.Second)
	fmt.Printf("%T\n", t)
	for now := range t {
		fmt.Printf("%v\n", now)
	}

	c := ParallelWrite([]byte("abcd"))
	fmt.Println(<-c)
	fmt.Println(<-c)

	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(j int) {
			fmt.Println(j) // Not the 'i' you are looking for.
			wg.Done()
		}(i)
	}
	wg.Wait()
}*/

func query(i int) int {
	fmt.Println("query start", i)
	//time.Sleep(time.Second * time.Duration(i))
	fmt.Println("query end ", i)
	return i
}

type simpleSleeper struct {
	seconds int
}

func (s *simpleSleeper) sleepFor(retry int) {
	if retry <= 1 {
		s.seconds = 1
	} else {
		s.seconds *= 2
	}
	if s.seconds > 32 {
		s.seconds = 32
	}
	fmt.Println("seconds ", s.seconds)
	//time.Sleep(time.Duration(s.seconds) * time.Second)
}

func walkFiles(done chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)
	go func() {
		defer close(paths)
		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			select {
			case paths <- path:
			case <-done:
				return errors.New("walk canceled")
			}
			return nil
		})
	}()
	return paths, errc
}

func ParseTimeRange(expr string) (string, error) {
	slice := strings.Split(expr, "-")
	if len(slice) != 2 {
		return "", errors.New(`time range should be seperated by "-"`)
	}
	if strings.TrimSpace(slice[0]) != `$(now)` {
		return "", errors.New(`time range should be prefixed with "$(now)"`)
	}
	d, err := time.ParseDuration(strings.TrimSpace(slice[1]))
	if err != nil {
		return "", fmt.Errorf("invalid time range: %v", err)
	}
	from := time.Now().Add(-d).Format("2006-01-02 15")
	to := time.Now().Format("2006-01-02 15")
	return fmt.Sprintf("%s ~ %s", from, to), nil
}

var bytePool = sync.Pool{
	New: newPool,
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func newPool() interface{} {
	b := make([]byte, 1024)
	return &b
}

func BenchmarkAlloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		_ = buf
	}
}

func BenchmarkPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()
		_ = buf
		bufPool.Put(buf)
	}
}

func printx(ctx context.Context) {
	list := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50}
	for i, num := range list {
		select {
		case <-ctx.Done():
			fmt.Println(i, "xxx--", ctx.Err())
			for _, num := range list[i:] {
				fmt.Println("yyy--", num)
			}
			fmt.Println("exit")
			return
		default:
		}
		fmt.Println(num, "---------")
		time.Sleep(time.Millisecond * 10)
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

var targetURL *url.URL

func handler(w http.ResponseWriter, r *http.Request) {
	o := new(http.Request)

	*o = *r

	o.Host = targetURL.Host
	o.URL.Scheme = targetURL.Scheme
	o.URL.Host = targetURL.Host
	o.URL.Path = singleJoiningSlash(targetURL.Path, o.URL.Path)

	if q := o.URL.RawQuery; q != "" {
		o.URL.RawPath = o.URL.Path + "?" + q
	} else {
		o.URL.RawPath = o.URL.Path
	}

	o.URL.RawQuery = targetURL.RawQuery

	o.Proto = "HTTP/1.1"
	o.ProtoMajor = 1
	o.ProtoMinor = 1
	o.Close = false

	transport := http.DefaultTransport

	res, err := transport.RoundTrip(o)

	if err != nil {
		log.Printf("http: proxy error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hdr := w.Header()

	for k, vv := range res.Header {
		for _, v := range vv {
			hdr.Add(k, v)
		}
	}

	w.WriteHeader(res.StatusCode)

	if res.Body != nil {
		io.Copy(w, res.Body)
	}
}

// var encoding = flag.Bool("encoding", false, "set jvm encoding")

func err1(n int) (int, error) {
	if n == 1 {
		return 1, errors.New("error 1")
	}
	return n, nil
}

func err2() (int, error) {
	return 2, errors.New("error 2")
}

func foo2(i int) (n int, err error) {
	defer func() {
		if err != nil {
			fmt.Println(n, err)
		}
	}()

	j, err := err1(2)
	if err != nil {
		return 0, err
	}
	_ = j

	for i > 0 {
		_, err := err2()
		if err != nil {
			return 5, err
		}
	}
	return n, err
}

type foo1 struct {
	bar string
}

type fooHandler struct {}

func (h fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func main() {
	http.Handle("/foo", fooHandler{})
	http.ListenAndServe(":8080", nil)

	//cmd := exec.Command("sleep", "60")
	//err := cmd.Start()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("Waiting for command to finish...")
	//err = cmd.Wait()
	//log.Printf("Command finished with error: %v", err)

	//tr := &http.Transport{
	//	ResponseHeaderTimeout: 5 * time.Minute,
	//}
	//var client = &http.Client{
	//	Transport: tr,
	//}
	//start := time.Now()
	//buf := bytes.NewBufferString("hello world")
	//log.Println("start...")
	//resp, err := client.Post("http://localhost:8080", "application/json", buf)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//b, _ := ioutil.ReadAll(resp.Body)
	//log.Printf("body: %s, elapsed: %v\n", b, time.Since(start))

	// s1 := stat{}
	// s1.disk.good = 2
	// s1.disk.bad = 1
	// a := stats{
	// 	s1,
	// }
	// fmt.Printf("%+v\n", a)

	// for i := 0; i < len(a); i++ {
	// 	a[i].disk.good = 3
	// 	a[i].disk.bad = 0
	// }
	// fmt.Printf("%+v\n", a)

	// conn, err := net.Dial("tcp", "cs21:40041")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer conn.Close()

	// host, port, _ := net.SplitHostPort(conn.LocalAddr().String())
	// fmt.Printf("local ip: %s, port: %s", host, port)

	// time.Sleep(time.Second * 30)
	// fmt.Println("close")

	// pfoo := &foo1{
	// 	bar: "bar1",
	// }
	// pfoo2 := &(*pfoo)
	// pfoo2.bar = "bar2"
	// fmt.Printf("pfoo.bar: %s, %p\n", pfoo.bar, pfoo)
	// fmt.Printf("pfoo2.bar: %s, %p\n", pfoo2.bar, pfoo2)

	// fmt.Printf("uid=%d, euid=%d, gid=%d, egid=%d\n",
	// 	syscall.Getuid(), syscall.Geteuid(), syscall.Getgid(), syscall.Getegid())

	// f, err := os.Open("/tmp/abc")
	// if err != nil {
	// 	log.Fatal("open fail: ", err)
	// }
	// b := make([]byte, 20)
	// n, err := f.Read(b)
	// fmt.Printf("read: %d, %s, %v\n", n, b, err)

	// time.Sleep(time.Hour)

	// m := gomail.NewMessage()

	// m.SetAddressHeader("From", "xingfeng25100@163.com", "xingfeng2510")
	// m.SetHeader("To", "xingmengbang@qiniu.com")
	// m.SetHeader("Cc", "danghexuan@qiniu.com", "xingmengbang@qiniu.com")
	// m.SetHeader("Subject", "Hello!")
	// m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")

	// d := gomail.NewDialer("smtp.163.com", 25, "xingfeng25100@163.com", "xingbang728")

	// Send the email to Bob, Cora and Dan.
	// if err := d.DialAndSend(m); err != nil {
	// 	panic(err)
	// }

	// fmt.Println("pid", syscall.Getpid())
	// fmt.Println("ppid", syscall.Getppid())
	// gpid, _ := syscall.Getpgid(syscall.Getpid())
	// fmt.Println("gpid", gpid)
	// sid, _ := syscall.Getsid(syscall.Getpid())
	// fmt.Println("sid", sid)
	// fmt.Println("uid", syscall.Getuid())

	// gen := func(ctx context.Context) <-chan int {
	// 	dst := make(chan int)
	// 	n := 1
	// 	go func() {
	// 		for {
	// 			select {
	// 			case <-ctx.Done():
	// 				fmt.Println(ctx.Err())
	// 				return // returning not to leak the goroutine
	// 			case dst <- n:
	// 				n++
	// 			}
	// 		}
	// 	}()
	// 	return dst
	// }

	// ctx, cancel := context.WithCancel(context.Background())

	// for n := range gen(ctx) {
	// 	fmt.Println(n)
	// 	if n == 5 {
	// 		break
	// 	}
	// }

	// cancel() // cancel when we are finished consuming integers
	// time.Sleep(10 * time.Second)

	// fmt.Println(pa.lerr)
	// fmt.Println(pa.get(), pa.lerr)
	// url, err := url.Parse("https://www.qiniu.com")
	// if err != nil {
	// 	panic("bad url")
	// }
	// targetURL = url
	// http.HandleFunc("/", handler)
	// http.ListenAndServe(":1234", nil)

	// ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	// defer cancel()

	// printx(ctx)

	// for i := 0; i < 30; i++ {
	// 	select {
	// 	case <-ctx.Done():
	// 		fmt.Println("done")
	// 	default:
	// 	}
	// 	fmt.Println(i, ctx.Err())
	// }

	// l := new(sync.Mutex)
	// c := sync.NewCond(l)
	// fire := false
	// count := 1
	// go func() {
	// 	c.L.Lock()
	// 	for !fire {
	// 		fmt.Println("before fire", count)
	// 		c.Wait()
	// 		count++
	// 	}
	// 	c.L.Unlock()
	// 	fmt.Println("firing", count)
	// }()

	// for i := 0; i < 100; i++ {
	// 	if i == 9 {
	// 		c.L.Lock()
	// 		fire = true
	// 		c.Signal()
	// 		c.L.Unlock()
	// 		break
	// 	}
	// 	time.Sleep(time.Millisecond * 50)
	// 	c.Signal()
	// }

	// time.Sleep(time.Second * 20)

	// b := bufPool.Get().(*bytes.Buffer)
	// b.Reset()
	// b.WriteString("hello world")
	// fmt.Println(b.String())
	// bufPool.Put(b)

	// c := make(chan int, 1)
	// for i := 1; i <= 3; i++ {
	// 	go func(i int) {
	// 		select {
	// 		case c <- query(i):
	// 			fmt.Println("send ", i)
	// 		default:
	// 			fmt.Println("nonblock ", i)
	// 		}
	// 	}(i)
	// }
	// time.Sleep(time.Second)
	// fmt.Println("receive: ", <-c)

	// time.Sleep(time.Second * 10)

	// for skip := 0; ; skip++ {
	// 	pc, file, line, ok := runtime.Caller(skip)
	// 	if !ok {
	// 		break
	// 	}
	// 	fmt.Printf("skip = %v, pc = %v, file = %v, line = %v\n", skip, pc, file, line)
	// }

	// flag.Parse()

	// ch := make(chan bool, 0)
	// go func() {
	// 	fmt.Println("before")
	// 	<-ch
	// 	fmt.Println("after")
	// }()

	// time.Sleep(time.Second * 1)
	// close(ch)
	// time.Sleep(time.Second * 1)

	// myType := &MyType{22, "helo"}
	// mtt := reflect.TypeOf(myType)
	// nm := mtt.NumMethod()
	// for i := 0; i < nm; i++ {
	// 	fmt.Printf("method %d: %s %s\n", i, mtt.Method(i).Name, mtt.Method(i).Type)
	// }

	// var x float64 = 3.4
	// fmt.Println("type: ", reflect.TypeOf(x))

	// f := func(r rune) bool {
	// 	return !unicode.IsNumber(r) && !unicode.IsLetter(r)
	// }
	// fmt.Printf("Fields are %v", strings.FieldsFunc("  foo1;bar2,baz3...", f))

	// f2 := fib()
	// fmt.Println(f2(), f2(), f2())

	// fmt.Println(filter([]int{1, 2, 5, 7, 8, 9}, func(a int) bool {
	// 	return a%2 == 0
	// }))

	// str := "hello"
	// p := (*struct {
	// 	str uintptr
	// 	len int
	// })(unsafe.Pointer(&str))

	// fmt.Printf("%+v\n", p)

	// var slice []int32 = make([]int32, 5, 10)
	// p2 := (*struct {
	// 	array uintptr
	// 	len   int
	// 	cap   int
	// })(unsafe.Pointer(&slice))

	// fmt.Printf("output: %+v\n", p2)

	// var iface interface{} = "Hello World!"
	// p3 := (*struct {
	// 	tab  uintptr
	// 	data uintptr
	// })(unsafe.Pointer(&iface))

	// fmt.Printf("%+v\n", p3)

	// var ia int = 1
	// sa := []int{1, 2, 4}
	// fmt.Println(unsafe.Sizeof(ia), unsafe.Sizeof(str), unsafe.Sizeof(sa))

	// var d Datas
	// fmt.Println(unsafe.Offsetof(d.c0)) // 0
	// fmt.Println(unsafe.Offsetof(d.c1)) // 8
	// fmt.Println(unsafe.Offsetof(d.c2)) // 16
	// fmt.Println(unsafe.Offsetof(d.c3)) // 32

	// fmt.Println(unsafe.Alignof(d.c0))
	// fmt.Println(unsafe.Alignof(d.c1))
	// fmt.Println(unsafe.Alignof(d.c2))
	// fmt.Println(unsafe.Alignof(d.c3))

	// pa := &A{}
	// pp := unsafe.Pointer(pa)
	// offset := unsafe.Offsetof(pa.y)
	// var px *int32 = (*int32)(pp)
	// *px = 32
	// var py *int64 = (*int64)(unsafe.Pointer(uintptr(pp) + offset))
	// *py = 64
	// fmt.Println(pa.x, pa.y)

	// dec := json.NewDecoder(os.Stdin)
	// enc := json.NewEncoder(os.Stdout)
	// for {
	// 	var v map[string]interface{}
	// 	if err := dec.Decode(&v); err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// 	for k := range v {
	// 		if k != "Name" {
	// 			delete(v, k)
	// 		}
	// 	}
	// 	if err := enc.Encode(&v); err != nil {
	// 		log.Println(err)
	// 	}
	// }
}
