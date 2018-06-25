package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
)

type Struct struct {
	Pub string
	pri int `pri:"private"`
}

func (s *Struct) Pri() int {
	return s.pri
}

func (s *Struct) sum(o int) int {
	return s.pri + o
}

func (s *Struct) Sum(o int) int {
	s.pri = s.sum(o)
	return s.pri
}

func (s *Struct) Name(firstName, lastName string) string {
	return firstName + " " + lastName
}

func sum(a, b int) int {
	return a + b
}

type Movie struct {
	Title, Subtitle string
	Year            int
	Color           bool
	Actor           map[string]string
	Oscars          []string
	Sequel          *string
}

func main() {
	x := 1
	rx := reflect.ValueOf(&x).Elem()
	px := rx.Addr().Interface().(*int)
	*px = 2
	fmt.Println("x", x)
	rx.SetInt(4)
	fmt.Println("x", x)

	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Color:    false,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},

		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}
	Display("strangelove", strangelove)
	Display("os.Stderr", os.Stderr)

	data, err := Marshal(strangelove)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("strangelove: %s\n", data)

	s := &Struct{pri: 11}
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	fmt.Println("Type:", t)
	fmt.Println("Value:", v)
	fmt.Println("Kind:", t.Kind())

	for i := 0; i < t.Elem().NumField(); i++ {
		f := t.Elem().Field(i)
		fmt.Printf("struct field %d: %s, %s， embeded?: %v, tag: %v\n", i, f.Name, f.Type, f.Anonymous, f.Tag)
	}
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		fmt.Printf("struct method %d: %s, %s\n", i, m.Name, m.Type)
	}

	callMethod := func(s reflect.Value, methodName string, methodArgs ...reflect.Value) ([]reflect.Value, error) {
		t := s.Type()
		method, exist := t.MethodByName(methodName)
		if !exist {
			return nil, fmt.Errorf("\"%s\": is not existed for %s", methodName, t)
		}

		if regexp.MustCompile(`^[a-z]`).MatchString(method.Name) {
			return nil, fmt.Errorf("\"%s\": unexported field cannot be called", method.Name)
		}

		args := []reflect.Value{s}
		args = append(args, methodArgs...)

		return method.Func.Call(args), nil
	}

	fmt.Print("call Struct.Pri: ")
	values, _ := callMethod(v, "Pri")
	fmt.Println(values[0].Interface())

	fmt.Print("call Struct.Sum: ")
	values, _ = callMethod(v, "Sum", reflect.ValueOf(1))
	fmt.Println(values[0].Interface())

	fmt.Print("call Struct.Name: ")
	values, _ = callMethod(v, "Name", reflect.ValueOf("David"), reflect.ValueOf("Beckham"))
	fmt.Println(values[0].Interface())

	fmt.Print("call Struct.sum: ")
	fmt.Println(callMethod(v, "sum", reflect.ValueOf(1)))

	fmt.Print("call Struct.s: ")
	fmt.Println(callMethod(v, "s"))

	fn := reflect.ValueOf(sum)
	ft := fn.Type()
	for i := 0; i < ft.NumIn(); i++ {
		in := ft.In(i)
		fmt.Printf("function argument %d: %s\n", i, in)
	}
	for i := 0; i < ft.NumOut(); i++ {
		out := ft.In(i)
		fmt.Printf("function return value %d: %s\n", i, out)
	}

	i, j := 1, 3
	fmt.Printf("Call sum(%d, %d) function: %v\n",
		i, j,
		fn.Call([]reflect.Value{
			reflect.ValueOf(i),
			reflect.ValueOf(j),
		})[0].Interface(),
	)
}
