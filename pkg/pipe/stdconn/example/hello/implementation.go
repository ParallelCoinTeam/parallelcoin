package main

import (
	"fmt"
	"net/rpc"
	"os"
	
	"github.com/p9c/qu"
	
	"github.com/p9c/pod/pkg/pipe/stdconn"
)

type Hello struct {
	Quit qu.C
}

func NewHello() *Hello {
	return &Hello{qu.T()}
}

func (h *Hello) Say(name string, reply *string) (e error) {
	r := "hello " + name
	*reply = r
	return
}

func (h *Hello) Bye(_ int, reply *string) (e error) {
	r := "i hear and obey *dies*"
	*reply = r
	h.Quit.Q()
	return
}

func main() {
	printlnE("starting up example worker")
	hello := NewHello()
	stdConn := stdconn.New(os.Stdin, os.Stdout, hello.Quit)
	e := rpc.Register(hello)
	if e != nil  {
		printlnE(e)
		return
	}
	go rpc.ServeConn(stdConn)
	hello.Quit.Wait()
	printlnE("i am dead! x_X")
}

func printlnE(a ...interface{}) {
	out := append([]interface{}{"[Hello]"}, a...)
	_, _ = fmt.Fprintln(os.Stderr, out...)
}
