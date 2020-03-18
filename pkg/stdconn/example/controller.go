package main

import (
	log "github.com/p9c/pod/pkg/logi"
	"github.com/p9c/pod/pkg/stdconn/example/hello/hello"
	"github.com/p9c/pod/pkg/stdconn/worker"
)

func main() {
	log.L.SetLevel("trace", true)
	L.Info("starting up example controller")
	cmd := worker.Spawn("go", "run", "hello/worker.go")
	client := hello.NewClient(cmd.StdConn)
	L.Info("calling Hello.Say with 'worker'")
	L.Info("reply:", client.Say("worker"))
	L.Info("calling Hello.Bye")
	L.Info("reply:", client.Bye())
}
