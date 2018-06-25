package main

import (
	"log"
	"os"
)

func main() {
	file, err := os.Create("/tmp/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString("This is first line\n")
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("file size after first write: %d\n", fi.Size())

	file.WriteString("This is second line\n")
	fi, err = file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("file size after second write: %d\n", fi.Size())
}
