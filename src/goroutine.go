package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

func foo(wg *sync.WaitGroup, ch chan<- int, seq int) {
	defer wg.Done()

	if seq%2 == 0 {
		fmt.Printf("foo %d is even\n", seq)
		return
	}

	ch <- seq
	return
}

func main1() {
	b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
	var f interface{}
	err := json.Unmarshal(b, &f)
	fmt.Printf("%+v, %v", f, err)

	m := f.(map[string]interface{})
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case int:
			fmt.Println(k, "is int", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}

	sum := make(chan int, 10)
	wg := &sync.WaitGroup{}
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go foo(wg, sum, i)
	}
	wg.Wait()

	close(sum)
	xx := 0
	for part := range sum {
		xx += part
	}
	fmt.Println("sum", xx)

	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	for {
		var v map[string]interface{}
		if err := dec.Decode(&v); err != nil {
			log.Println(err)
			return
		}
		for k := range v {
			if k != "Name" {
				delete(v, k)
			}
		}
		if err := enc.Encode(&v); err != nil {
			log.Println(err)
		}
	}
}
