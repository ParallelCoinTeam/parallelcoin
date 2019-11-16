package main

import (
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/worker"
	"io"
	"net/rpc"
)

func main() {
	log.L.SetLevel("trace", true)
	log.INFO("starting up example controller")
	cmd := worker.SpawnWorker("go", "run", "worker.go")
	client := NewHelloClient(cmd.StdConn)
	log.INFO("calling Hello.Say with 'worker'")
	log.INFO("reply:", client.Say("worker"))
	log.INFO("calling Hello.Bye")
	log.INFO("reply:", client.Bye())
	if err := cmd.KillWorker(); err != nil {
		log.ERROR(err)
	}
}

type HelloClient struct {
	*rpc.Client
}

func NewHelloClient(conn io.ReadWriteCloser) *HelloClient {
	return &HelloClient{rpc.NewClient(conn)}

}

func (h *HelloClient) Say(name string) (reply string) {
	err := h.Call("Hello.Say", "worker", &reply)
	if err != nil {
		log.ERROR(err)
		return "error: " + err.Error()
	}
	return
}

func (h *HelloClient) Bye() (reply string) {
	err := h.Call("Hello.Bye", 1, &reply)
	if err != nil {
		log.ERROR(err)
		return "error: " + err.Error()
	}
	return
}
