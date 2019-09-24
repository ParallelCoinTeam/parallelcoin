package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type PrintlnFunc func(a ...interface{})
type PrintfFunc func(format string, a ...interface{})
type Closure func(func() string)

const (
	Off   = "off"
	Fatal = "fatal"
	Error = "error"
	Warn  = "warn"
	Info  = "info"
	Debug = "debug"
	Trace = "trace"
)

type Logger struct {
	Fatal         PrintlnFunc
	Error         PrintlnFunc
	Warn          PrintlnFunc
	Info          PrintlnFunc
	Debug         PrintlnFunc
	Trace         PrintlnFunc
	Fatalf        PrintfFunc
	Errorf        PrintfFunc
	Warnf         PrintfFunc
	Infof         PrintfFunc
	Debugf        PrintfFunc
	Tracef        PrintfFunc
	Fatalc        Closure
	Errorc        Closure
	Warnc         Closure
	Infoc         Closure
	Debugc        Closure
	Tracec        Closure
	LogFileHandle *os.File
}

type Entry struct {
	Time         time.Time
	Level        string
	CodeLocation string
	Text         string
}

func Empty() *Logger {
	return &Logger{
		Fatal:  NoPrintln,
		Error:  NoPrintln,
		Warn:   NoPrintln,
		Info:   NoPrintln,
		Debug:  NoPrintln,
		Trace:  NoPrintln,
		Fatalf: NoPrintf,
		Errorf: NoPrintf,
		Warnf:  NoPrintf,
		Infof:  NoPrintf,
		Debugf: NoPrintf,
		Tracef: NoPrintf,
		Fatalc: NoClosure,
		Errorc: NoClosure,
		Warnc:  NoClosure,
		Infoc:  NoClosure,
		Debugc: NoClosure,
		Tracec: NoClosure,
	}

}

func NewLogger(level, logPath, logFileName string) (l *Logger, close func(),
	err error) {
	const timeFormat = "2006-01-02_15-04-05"
	path := filepath.Join(logFileName, logPath)
	var logFileHandle *os.File
	if FileExists(path) {
		err := os.Rename(path, filepath.Join(logPath,
			time.Now().Format(timeFormat)+".json"))
		if err != nil {
			fmt.Println("error rotating log", err)
			return
		}
	}
	logFileHandle, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("error opening log file", logFileName)
	}
	l = Empty()
	l.SetLevel(level)
	l.LogFileHandle = logFileHandle
	_, _ = fmt.Fprintln(logFileHandle, "{")

	return
}

func (l *Logger) SetLevel(level string) {
	*l = *Empty()
	var fallen bool
	switch {
	case level == Fatal:
		l.Fatal = Println(level, l.LogFileHandle)
		l.Fatalf = Printf(level, l.LogFileHandle)
		l.Fatalc = Printc(level, l.LogFileHandle)
		fallen = true
		fallthrough
	case level == Error && fallen:
		l.Error = Println(level, l.LogFileHandle)
		l.Errorf = Printf(level, l.LogFileHandle)
		l.Errorc = Printc(level, l.LogFileHandle)
		fallthrough
	case level == Warn && fallen:
		l.Warn = Println(level, l.LogFileHandle)
		l.Warnf = Printf(level, l.LogFileHandle)
		l.Warnc = Printc(level, l.LogFileHandle)
		fallthrough
	case level == Info && fallen:
		l.Info = Println(level, l.LogFileHandle)
		l.Infof = Printf(level, l.LogFileHandle)
		l.Infoc = Printc(level, l.LogFileHandle)
		fallthrough
	case level == Debug && fallen:
		l.Debug = Println(level, l.LogFileHandle)
		l.Debugf = Printf(level, l.LogFileHandle)
		l.Debugc = Printc(level, l.LogFileHandle)
		fallthrough
	case level == Trace && fallen:
		l.Trace = Println(level, l.LogFileHandle)
		l.Tracef = Printf(level, l.LogFileHandle)
		l.Tracec = Printc(level, l.LogFileHandle)
	}
}

func NoPrintln(_ ...interface{})          {}
func NoPrintf(_ string, _ ...interface{}) {}
func NoClosure(_ func() string)           {}

func Println(level string, fh *os.File) func(a ...interface{}) {
	return func(a ...interface{}) {
		_, loc, line, _ := runtime.Caller(1)
		files := strings.Split(codeLoc, "github.com/parallelcointeam/parallelcoin/")
		codeLoc := fmt.Sprint(files[1], ":", line)
		text := fmt.Sprintln(a...)
		fmt.Println(text, codeLoc)
		if fh != nil {
			out := Entry{time.Now(), level, loc, text}
			j, err := json.Marshal(out)
			if err != nil {
				fmt.Println("logging error:", err)
			}
			_, _ = fmt.Fprintln(fh, string(j))
		}
	}
}

func Printf(level string, fh *os.File) func(format string, a ...interface{}) {
	return func(format string, a ...interface{}) {
		_, loc, line, _ := runtime.Caller(1)
		files := strings.Split(loc, "github.com/parallelcointeam/parallelcoin/")
		codeLoc := fmt.Sprint(files[1], ":", line)
		text := fmt.Sprintf(format, a...)
		fmt.Printf("%s %s", text, codeLoc)
		if fh != nil {
			out := Entry{time.Now(), level, loc, text}
			j, err := json.Marshal(out)
			if err != nil {
				fmt.Println("logging error:", err)
			}
			_, _ = fmt.Fprintln(fh, j)
		}
	}
}

func Printc(level string, fh *os.File) func(fn func() string) {
	return func(fn func() string) {
		text := fn()
		_, loc, line, _ := runtime.Caller(1)
		files := strings.Split(loc, "github.com/parallelcointeam/parallelcoin/")
		codeLoc := fmt.Sprint(files[1], ":", line)
		fmt.Printf("%s %s", text, codeLoc)
		if fh != nil {
			out := Entry{time.Now(), level, loc, text}
			j, err := json.Marshal(out)
			if err != nil {
				fmt.Println("logging error:", err)
			}
			_, _ = fmt.Fprintln(fh, j)
		}
	}
}

// FileExists reports whether the named file or directory exists.
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
