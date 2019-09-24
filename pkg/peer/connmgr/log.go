package connmgr

import (
	"github.com/parallelcointeam/parallelcoin/pkg/util/cl"
	"github.com/parallelcointeam/parallelcoin/pkg/util/pkgs"
)

// Log is the logger for the connmgr package
type _dtype int

var _d _dtype
var Log = cl.NewSubSystem(pkgs.Name(_d), "info")
var log = Log.Ch

// UseLogger uses a specified Logger to output package logging info.
func UseLogger(logger *cl.SubSystem) {
	Log = logger
	log = Log.Ch
}
