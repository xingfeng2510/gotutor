package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/rpc"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type IArith interface {
	Multiply(*Args, *int) error
	Divide(*Args, *Quotient) error
}

type ArithImpl struct {
	pimpl IArith
}

func (i *ArithImpl) Multiply(args *Args, reply *int) error {
	return i.pimpl.Multiply(args, reply)
}

func (i *ArithImpl) Divide(args *Args, quo *Quotient) error {
	return i.pimpl.Divide(args, quo)
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func main() {
	arith := new(Arith)
	var ai IArith = &ArithImpl{arith}
	rpc.Register(ai)
	rpc.HandleHTTP()

	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
