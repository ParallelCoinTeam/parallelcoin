package app

import (
	"github.com/p9c/log"
	"github.com/p9c/pod/pkg/pod"
)

// Main is the entrypoint for the pod suite
func Main() int {
	log.SetLogLevel("trace")
	var e error
	var cx *pod.State
	if cx, e = pod.GetNewContext(); F.Chk(e) {
		return 1
	}
	T.Ln("running command", cx.Config.RunningCommand.Name)
	if e = cx.Config.RunningCommand.Entrypoint(cx.Config); E.Chk(e) {
		return 1
	}
	return 0
}
