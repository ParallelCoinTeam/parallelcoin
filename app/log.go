package app

import (
	"github.com/parallelcointeam/parallelcoin/pkg/log"
)

type _dtype int

var _d _dtype
var L = log.NewLogger("info")

func UseLogger(logger *log.Logger) {
	L = logger
}
