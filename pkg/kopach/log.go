package kopach

import (
	"github.com/parallelcointeam/parallelcoin/pkg/util/cl"
	"github.com/parallelcointeam/parallelcoin/pkg/util/pkgs"
)

type _dtype int

var _d _dtype
var Log = cl.NewSubSystem(pkgs.Name(_d), "info")
var log = Log.Ch
