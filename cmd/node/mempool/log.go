package mempool

import (
	"github.com/p9c/pod/pkg/util/cl"
	"github.com/p9c/pod/pkg/util/pkgs"
)

// Log is the logger for the peer package
type _dtype int

var _d _dtype
var Log = cl.NewSubSystem(pkgs.Name(_d), "info")
var log = Log.Ch

// pickNoun returns the singular or plural form of a noun depending on the count n.
func pickNoun(n int, singular, plural string) string {

	if n == 1 {
		return singular
	}
	return plural
}
