package cfgutil

import (
	"runtime"
)

var (
	L = log.L
)

func init() {
	_, loc, _, _ := runtime.Caller(0)
	log.Register("pod", loc, L)
}
