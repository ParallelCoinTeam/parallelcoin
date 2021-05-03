package gel

import (
	"github.com/p9c/log"

	"github.com/p9c/gel/version"
)

var F, E, W, I, D, T = log.GetLogPrinterSet(log.AddLoggerSubsystem(version.PathBase))
