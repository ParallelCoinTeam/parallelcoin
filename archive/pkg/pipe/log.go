package pipe

import (
	"github.com/p9c/log"
	"github.com/p9c/pod/version"
)

var subsystem = log.AddLoggerSubsystem(version.PathBase)
var F, E, W, I, D, T log.LevelPrinter = log.GetLogPrinterSet(subsystem)
