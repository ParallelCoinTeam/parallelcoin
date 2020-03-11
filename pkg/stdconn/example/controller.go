package main

import (
	log "github.com/p9c/logi"
	"github.com/p9c/pod/pkg/stdconn/example/hello/hello"
	"github.com/p9c/pod/pkg/stdconn/worker"
)

func main() {
	log.L.SetLevel("trace", true)
	log.L.Info("starting up example controller")
	cmd := worker.Spawn("go", "run", "hello/worker.go")
	client := hello.NewClient(cmd.StdConn)
	log.L.Info("calling Hello.Say with 'worker'")
	log.L.Info("reply:", client.Say("worker"))
	log.L.Info("calling Hello.Bye")
	log.L.Info("reply:", client.Bye())
}
