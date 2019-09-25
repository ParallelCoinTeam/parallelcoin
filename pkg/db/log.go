package database

import (
	"github.com/p9c/pod/pkg/util/cl"
	"github.com/p9c/pod/pkg/util/pkgs"
)

type _dtype int

// Log is the logger for the peer package
var (
	_d  _dtype
	Log = cl.NewSubSystem(pkgs.Name(_d), "info")
	log = Log.Ch
)

// UseLogger uses a specified Logger to output package logging info.
func UseLogger(logger *cl.SubSystem) {
	Log = logger
	log = Log.Ch
}
