package main

import (
	"github.com/p9c/pod/pkg/comm/stdconn/example/hello/hello"
	"github.com/p9c/pod/pkg/comm/stdconn/worker"
	log "github.com/p9c/pod/pkg/util/logi"
	qu "github.com/p9c/pod/pkg/util/quit"
)

func main() {
	log.L.SetLevel("trace", true, "pod")
	Info("starting up example controller")
	cmd, _ := worker.Spawn(make(qu.C), "go", "run", "hello/worker.go")
	client := hello.NewClient(cmd.StdConn)
	Info("calling Hello.Say with 'worker'")
	Info("reply:", client.Say("worker"))
	Info("calling Hello.Bye")
	Info("reply:", client.Bye())
	cmd.Kill()
}
