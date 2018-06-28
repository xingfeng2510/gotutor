package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func xor(a ...int) int {
	result := 0
	for _, x := range a {
		result ^= x
	}
	return result
}

type auth struct {
	ak string
	sk string
}

var m map[auth]int

func main() {
	if m == nil {
		xx := auth{"1", "2"}
		r, found := m[xx]
		fmt.Println(r, found)
		println("nil map")
		m = make(map[auth]int)
	}

	a1 := auth{"ak1", "sk1"}
	m[a1] = 1

	x1, ok := m[a1]
	if ok {
		fmt.Println("find", a1, x1)
	}

	_, ok = m[auth{"ak1", "ak2"}]
	if !ok {
		fmt.Println("not found")
	}
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
