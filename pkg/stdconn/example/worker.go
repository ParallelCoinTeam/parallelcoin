package main

import (
	"fmt"
	"github.com/p9c/pod/pkg/stdconn"
	"net/rpc"
	"os"
)

type Hello struct {
	Quit chan struct{}
}

func NewHello() *Hello {
	return &Hello{make(chan struct{})}
}

func (h *Hello) Say(name string, reply *string) (err error) {
	r := "hello " + name
	*reply = r
	return
}

func (h *Hello) Bye(_ int, reply *string) (err error) {
	r := "i hear and obey *dies*"
	*reply = r
	close(h.Quit)
	return
}

func main() {
	printlnE("starting up example worker")
	hello := NewHello()
	stdConn := stdconn.New(os.Stdin, os.Stdout, hello.Quit)
	err := rpc.Register(hello)
	if err != nil {
		printlnE(err)
		return
	}
	go rpc.ServeConn(stdConn)
	<-hello.Quit
	printlnE("i am dead! x_X")
}

func printlnE(a ...interface{}) {
	out := append([]interface{}{"[Hello]"}, a...)
	_, _ = fmt.Fprintln(os.Stderr, out...)
}
