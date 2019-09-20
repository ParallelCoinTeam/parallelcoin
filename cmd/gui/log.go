//+build !nogui


package gui

import (
	"git.parallelcoin.io/dev/pod/pkg/util/cl"
	"git.parallelcoin.io/dev/pod/pkg/util/pkgs"
)

type _dtype int

var _d _dtype
var Log = cl.NewSubSystem(pkgs.Name(_d), "info")
var log = Log.Ch
