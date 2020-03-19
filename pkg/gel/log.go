package gel

import (
	"runtime"

	log "github.com/p9c/pod/pkg/logi"
)

var (
	L = log.L
)

func init() {
	_, loc, _, _ := runtime.Caller(0)
	log.Register("pod", loc, L)
}
