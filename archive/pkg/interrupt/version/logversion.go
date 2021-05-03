package version

import (
	"path/filepath"
	"runtime"

	"github.com/p9c/log"
)

var F, E, W, I, D, T log.LevelPrinter

func init() {
	_, file,_, _ := runtime.Caller(0)
	verPath := filepath.Dir(file)+"/"
	F, E, W, I, D, T = log.GetLogPrinterSet(log.AddLoggerSubsystem(verPath))
}
