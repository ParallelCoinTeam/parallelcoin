package rpcclient

import (
	"github.com/p9c/pod/pkg/util/cl"
	"github.com/p9c/pod/pkg/util/pkgs"
)

// Log is the logger for the rpcclient package
type _dtype int

var _d _dtype
var Log = cl.NewSubSystem(pkgs.Name(_d), "info")
var log = Log.Ch

// UseLogger uses a specified Logger to output package logging info.
func UseLogger(logger *cl.SubSystem) {
	Log = logger
	log = Log.Ch
}

// // LogClosure is a closure that can be printed with %v to be used to generate expensive-to-create data for a detailed log level and avoid doing the work if the data isn't printed.
// type logClosure func() string

// // String invokes the log closure and returns the results string.
// func (c logClosure) String() string {
// 	return c()
// }

// // newLogClosure returns a new closure over the passed function which allows it to be used as a parameter in a logging function that is only invoked when the logging level is such that the message will actually be logged.
// func newLogClosure(// 	c func() string) logClosure {
// 	return logClosure(c)
// }
