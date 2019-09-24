package chain

import (
	"github.com/parallelcointeam/parallelcoin/pkg/util/cl"
	"github.com/parallelcointeam/parallelcoin/pkg/util/pkgs"
)

// Log is the logger for the peer package
type _dtype int

var _d _dtype
var Log = cl.NewSubSystem(pkgs.Name(_d), "info")
var log = Log.Ch
