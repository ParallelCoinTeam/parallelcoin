package main

import (
	"fmt"
	qu "github.com/p9c/pod/pkg/util/quit"
	"net/rpc"
	"os"

	"github.com/p9c/pod/pkg/comm/stdconn"
)

type Hello struct {
	Quit qu.C
}

func NewHello() *Hello {
	return &Hello{qu.T()}
}

func (h *Hello) Say(name string, reply *string) (err error) {
	r := "hello " + name
	*reply = r
	return
}

func (h *Hello) Bye(_ int, reply *string) (err error) {
	r := "i hear and obey *dies*"
	*reply = r
	h.Quit.Q()
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
	hello.Quit.Wait()
	printlnE("i am dead! x_X")
}

func printlnE(a ...interface{}) {
	out := append([]interface{}{"[Hello]"}, a...)
	_, _ = fmt.Fprintln(os.Stderr, out...)
}
