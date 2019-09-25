package addrmgr

import (
	"github.com/p9c/pod/pkg/util/cl"
	"github.com/p9c/pod/pkg/util/pkgs"
)

// Log is the logger for the addrmgr package
type _dtype int

var _d _dtype
var Log = cl.NewSubSystem(pkgs.Name(_d), "info")
var log = Log.Ch
