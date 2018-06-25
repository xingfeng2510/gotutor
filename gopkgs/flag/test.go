package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type stringArray []string

func (a *stringArray) Set(s string) error {
	*a = append(*a, s)
	return nil
}

func (a *stringArray) String() string {
	return strings.Join(*a, ",")
}

func main() {
	flagSet := flag.NewFlagSet("flagtest", flag.ExitOnError)
	addrs := stringArray{}
	flagSet.Var(&addrs, "address", "server address")
	flagSet.Parse(os.Args[1:])
	f := flagSet.Lookup("address")
	if f != nil {
		fmt.Println(f.Value.String())
	}
}
