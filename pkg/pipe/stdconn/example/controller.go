package main

import (
	"github.com/p9c/log"
	"github.com/p9c/pod/pkg/pipe/stdconn/example/hello/hello"
	"github.com/p9c/pod/pkg/pipe/stdconn/worker"
	"github.com/p9c/qu"
)

func main() {
	log.SetLogLevel("trace")
	I.Ln("starting up example controller")
	cmd, _ := worker.Spawn(qu.T(), "go", "run", "hello/worker.go")
	client := hello.NewClient(cmd.StdConn)
	I.Ln("calling Hello.Say with 'worker'")
	I.Ln("reply:", client.Say("worker"))
	I.Ln("calling Hello.Bye")
	I.Ln("reply:", client.Bye())
	if e := cmd.Kill(); E.Chk(e) {
	}
}
