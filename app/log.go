package app

import (
	"github.com/p9c/pod/pkg/util/cl"
	"github.com/p9c/pod/pkg/util/pkgs"
)

// Log is the logger for node
type _dtype int

var _d _dtype
var Log = cl.NewSubSystem(pkgs.Name(_d), "info")
var log = Log.Ch
