package clipboard

import (
	"github.com/p9c/pod/pkg/logg"
)

var subsystem = logg.AddLoggerSubsystem()
var ftl, e, wrn, inf, dbg, trc = logg.GetLogPrinterSet(subsystem)
