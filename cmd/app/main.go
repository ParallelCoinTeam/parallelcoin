package app

import (
	"github.com/p9c/log"
	"github.com/p9c/pod/pkg/pod"
)

func Main() int {
	log.SetLogLevel("trace")
	var e error
	var cx *pod.State
	if cx, e = pod.GetNewContext(); F.Chk(e) {
		return 1
	}
	_ = cx
	return 0
}
