package btcjson_test

import (
	"github.com/p9c/pod/pkg/logg"
)

var subsystem = logg.AddLoggerSubsystem()
var ftl, err, wrn, inf, dbg, trc = logg.GetLogPrinterSet(subsystem)
