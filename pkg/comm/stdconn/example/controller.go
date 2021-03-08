package main

import (
	"github.com/p9c/pod/pkg/comm/stdconn/example/hello/hello"
	"github.com/p9c/pod/pkg/comm/stdconn/worker"
	log "github.com/p9c/pod/pkg/util/logi"
	qu "github.com/p9c/pod/pkg/util/quit"
)

func main() {
	log.L.SetLevel("trace", true, "pod")
	inf.Ln("starting up example controller")
	cmd, _ := worker.Spawn(qu.T(), "go", "run", "hello/worker.go")
	client := hello.NewClient(cmd.StdConn)
	inf.Ln("calling Hello.Say with 'worker'")
	inf.Ln("reply:", client.Say("worker"))
	inf.Ln("calling Hello.Bye")
	inf.Ln("reply:", client.Bye())
	if e := cmd.Kill(); dbg.Ok(err) {
	}
}
