package interrupt

import (

	"git.parallelcoin.io/dev/pod/pkg/util/cl"
	"git.parallelcoin.io/dev/pod/pkg/util/pkgs"
)

// Log is the logger for node
//nolint
type _dtype int

var _d _dtype
var Log = cl.NewSubSystem(pkgs.Name(_d), "info")
var log = Log.Ch

// UseLogger uses a specified Logger to output package logging info. This
// should be used in preference to SetLogWriter if the caller is also using log.
func UseLogger(logger *cl.SubSystem,
) {
	Log = logger
	log = Log.Ch
}
