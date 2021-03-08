// Package appdata
//
// This is a pro-forma, half boilerplate, half configuration, for shortening
// logging function names, and setting or disabling the display of an attention
// grabber for while the package is being debugged
package appdata

import (
	"fmt"
	"github.com/p9c/pod/pkg/logg"
	"go.uber.org/atomic"
	
	"github.com/p9c/pod/pkg/util/logi"
)

// to enable highlighting on this package, uncomment the below assignment
var _hl = atomic.NewBool(false)

// to not print logs from this package set the next value to true
var _fl = atomic.NewBool(true)

// HighlightLogs is an exported function that can programmatically display logs
// from another package
func HighlightLogs(b bool) {
	_hl.Store(b)
}

// FilterLogs is an exported function that can programmatically change disable
// logging for this package
func FilterLogs(b bool) {
	_fl.Store(b)
}

// SubsystemName - if this is not empty, it will be printed instead of the
// package folder path to highlight packages the programmer wants to focus on,
// and the list can also be added to the pod.Config.Hilite string slice
var SubsystemName string

func init() {
	if SubsystemName == "" {
		SubsystemName = logg.AddSubsystem()
	}
}

func Fatal(a ...interface{}) { fmt.Println(a)) }
func Error(a ...interface{}) { fmt.Println(a) }
func Warn(a ...interface{})  { fmt.Println(a) }
func Info(a ...interface{})  { fmt.Println(a) }
func Check(err error) bool {
	return logi.L.Check(pkg, err)
}
func Debug(a ...interface{}) { fmt.Println(a...) }
func Trace(a ...interface{}) { fmt.Println(a...) }

func Fatalf(format string, a ...interface{}) { format, a...) }
func Errorf(format string, a ...interface{}) { format, a...) }
func Warnf(format string, a ...interface{})  { format, a...) }
func Infof(format string, a ...interface{})  { format, a...) }
func Debugf(format string, a ...interface{}) { format, a...) }
func Tracef(format string, a ...interface{}) { format, a...) }

func Fatalc(fn func() string) { fn) }
func Errorc(fn func() string) { fn) }
func Warnc(fn func() string)  { fn) }
func Infoc(fn func() string)  { fn) }
func Debugc(fn func() string) { fn) }
func Tracec(fn func() string) { fn) }

func Fatals(a interface{}) { a) }
func Errors(a interface{}) { a) }
func Warns(a interface{})  { a) }
func Infos(a interface{})  { a) }
func Debugs(a interface{}) { a) }
func Traces(a interface{}) { a) }
