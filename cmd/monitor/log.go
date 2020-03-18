package monitor

import (
	"os"
	"runtime"
	"strings"

	log "github.com/p9c/pod/pkg/logi"
)

var (
	L = log.L
)

func init() {
	_, loc, _, _ := runtime.Caller(0)
	files := strings.Split(loc, "pod")
	var pkg string
	pkg = loc
	if len(files) > 1 {
		pkg = files[1]
	}
	split := strings.Split(pkg, string(os.PathSeparator))
	pkg = strings.Join(split[:len(split)-1], string(os.PathSeparator))
	L = log.Empty(pkg).SetLevel("trace", true, "pod")
	log.Loggers[pkg] = L
}
