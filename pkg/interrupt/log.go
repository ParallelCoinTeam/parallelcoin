package interrupt

import (
	"github.com/p9c/log"
	"github.com/p9c/interrupt/version"
)

var F, E, W, I, D, T = log.GetLogPrinterSet(log.AddLoggerSubsystem(version.PathBase))
