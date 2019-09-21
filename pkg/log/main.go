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

const screenWidth = 144

var StartupTime = time.Now()

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

var Levels = []string{
	Off, Fatal, Error, Warn, Info, Debug, Trace,
}

// Logger is a struct containing all the functions with nice handy names
type Logger struct {
	Name          string
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

// Entry is a log entry to be printed as json to the log file
type Entry struct {
	Time         time.Time
	Level        string
	CodeLocation string
	Text         string
}

func Empty() *Logger {
	return &Logger{
		Fatal:  NoPrintln(),
		Error:  NoPrintln(),
		Warn:   NoPrintln(),
		Info:   NoPrintln(),
		Debug:  NoPrintln(),
		Trace:  NoPrintln(),
		Fatalf: NoPrintf(),
		Errorf: NoPrintf(),
		Warnf:  NoPrintf(),
		Infof:  NoPrintf(),
		Debugf: NoPrintf(),
		Tracef: NoPrintf(),
		Fatalc: NoClosure(),
		Errorc: NoClosure(),
		Warnc:  NoClosure(),
		Infoc:  NoClosure(),
		Debugc: NoClosure(),
		Tracec: NoClosure(),
	}

}

// sanitizeLoglevel accepts a string and returns a
// default if the input is not in the Levels slice
func sanitizeLoglevel(level string) string {
	fmt.Println("sanitise")
	found := false
	for i := range Levels {
		if level == Levels[i] {
			found = true
			break
		}
	}
	fmt.Println(level, found)
	if !found {
		fmt.Println("default info")
		level = "info"
	}
	return level
}

// NewLogger creates a new logger with json entries
func NewLogger(level string) (l *Logger) {
	l = Empty()
	l.SetLevel(level)
	Register.Add(l)
	return
}

// SetLogPaths sets a file path to write logs
func (l *Logger) SetLogPaths(logPath, logFileName string) {
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
	logFileHandle, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("error opening log file", logFileName)
	}
	l.LogFileHandle = logFileHandle
	_, _ = fmt.Fprintln(logFileHandle, "{")
}

// SetLevel enables or disables the various print functions
func (l *Logger) SetLevel(level string) {
	*l = *Empty()
	var fallen bool
	switch {
	case level == Trace || fallen:
		l.Trace = Println("T", l.LogFileHandle)
		l.Tracef = Printf("T", l.LogFileHandle)
		l.Tracec = Printc("T", l.LogFileHandle)
		fallen = true
		fallthrough
	case level == Debug || fallen:
		l.Debug = Println("D", l.LogFileHandle)
		l.Debugf = Printf("D", l.LogFileHandle)
		l.Debugc = Printc("D", l.LogFileHandle)
		fallen = true
		fallthrough
	case level == Info || fallen:
		l.Info = Println("I", l.LogFileHandle)
		l.Infof = Printf("I", l.LogFileHandle)
		l.Infoc = Printc("I", l.LogFileHandle)
		fallen = true
		fallthrough
	case level == Warn || fallen:
		l.Warn = Println("W", l.LogFileHandle)
		l.Warnf = Printf("W", l.LogFileHandle)
		l.Warnc = Printc("W", l.LogFileHandle)
		fallen = true
		fallthrough
	case level == Error || fallen:
		l.Error = Println("E", l.LogFileHandle)
		l.Errorf = Printf("E", l.LogFileHandle)
		l.Errorc = Printc("E", l.LogFileHandle)
		fallen = true
		fallthrough
	case level == Fatal:
		l.Fatal = Println("F", l.LogFileHandle)
		l.Fatalf = Printf("F", l.LogFileHandle)
		l.Fatalc = Printc("F", l.LogFileHandle)
		fallen = true
	}
}

var NoPrintln = func() func(_ ...interface{}) {
	return func(_ ...interface{}) {
	}
}
var NoPrintf = func() func(_ string, _ ...interface{}) {
	return func(_ string, _ ...interface{}) {
	}
}
var NoClosure = func() func(_ func() string) {
	return func(_ func() string) {
	}
}

func trimReturn(s string) string {
	return s[:len(s)-1]
}

func rightJustify(n int) string {
	s := fmt.Sprint(n)
	switch len(s) {
	case 1:
		s += "   "
	case 2:
		s += "  "
	case 3:
		s += " "
	}
	return s
}

// Println prints a log entry like Println
func Println(level string, fh *os.File) func(a ...interface{}) {
	// level = strings.ToUpper(string(level[0]))
	return func(a ...interface{}) {
		_, loc, line, _ := runtime.Caller(1)
		files := strings.Split(loc, "github.com/parallelcointeam/parallelcoin/")
		codeLoc := fmt.Sprint(files[1], ":", rightJustify(line))
		since := fmt.Sprint(time.Now().Sub(StartupTime) / time.
			Second * time.Second)
		text := since + " " + level + " "
		indent := strings.Repeat(" ", len(text))
		text += trimReturn(fmt.Sprintln(a...))
		// wordwrap :p
		split := strings.Split(text, " ")
		out := split[0] + " "
		var final string
		cod := false
		for i := range split {
			if i > 0 {
				if len(out)+len(split[i])+1+len(codeLoc) > screenWidth && !cod {
					cod = true
					final += out + strings.Repeat(".",
						screenWidth-len(out)-len(codeLoc)) + " " +
						codeLoc + "\n"
					out = indent + split[i] + " "
				} else if len(out)+len(split[i]) > screenWidth {
					out = indent + split[i] + " "
				} else {
					out += split[i] + " "
				}
			}
		}
		final += out
		if !cod {
			rem := screenWidth - len(out) - len(codeLoc)
			if rem < 1 {
				final += "\n" + strings.Repeat(" ", screenWidth-len(codeLoc)) + codeLoc
			} else {
				final += strings.Repeat(".", rem) + " " + codeLoc
			}
		}
		fmt.Println(final)
		if fh != nil {
			out := Entry{time.Now(), level, loc, text}
			j, err := json.Marshal(out)
			if err != nil {
				fmt.Println("logging error:", err)
			}
			_, _ = fmt.Fprint(fh, string(j)+",")
		}
	}
}

// Printf prints a log entry with formatting
func Printf(level string, fh *os.File) func(format string, a ...interface{}) {
	// level = strings.ToUpper(string(level[0]))
	return func(format string, a ...interface{}) {
		_, loc, line, _ := runtime.Caller(1)
		files := strings.Split(loc, "github.com/parallelcointeam/parallelcoin/")
		codeLoc := fmt.Sprint(files[1], ":", rightJustify(line))
		since := fmt.Sprint(time.Now().Sub(StartupTime) / time.
			Second * time.Second)
		text := since + " " + level + " "
		indent := strings.Repeat(" ", len(text))
		text += trimReturn(fmt.Sprintln(a...))
		// wordwrap :p
		split := strings.Split(text, " ")
		out := split[0] + " "
		var final string
		cod := false
		for i := range split {
			if i > 0 {
				if len(out)+len(split[i])+1+len(codeLoc) > screenWidth && !cod {
					cod = true
					final += out + strings.Repeat(".",
						screenWidth-len(out)-len(codeLoc)) + " " +
						codeLoc + "\n"
					out = indent + split[i] + " "
				} else if len(out)+len(split[i]) > screenWidth {
					out = indent + split[i] + " "
				} else {
					out += split[i] + " "
				}
			}
		}
		final += out
		if !cod {
			rem := screenWidth - len(out) - len(codeLoc)
			if rem < 1 {
				final += "\n" + strings.Repeat(" ", screenWidth-len(codeLoc)) + codeLoc
			} else {
				final += strings.Repeat(".", rem) + " " + codeLoc
			}
		}
		fmt.Println(final)
		if fh != nil {
			out := Entry{time.Now(), level, loc, text}
			j, err := json.Marshal(out)
			if err != nil {
				fmt.Println("logging error:", err)
			}
			_, _ = fmt.Fprintln(fh, string(j)+",")
		}
	}
}

// Printc prints from a closure returning a string
func Printc(level string, fh *os.File) func(fn func() string) {
	// level = strings.ToUpper(string(level[0]))
	return func(fn func() string) {
		t := fn()
		_, loc, line, _ := runtime.Caller(1)
		files := strings.Split(loc, "github.com/parallelcointeam/parallelcoin/")
		codeLoc := fmt.Sprint(files[1], ":", rightJustify(line))
		since := fmt.Sprint(time.Now().Sub(StartupTime) / time.
			Second * time.Second)
		text := since + " " + level + " "
		indent := strings.Repeat(" ", len(t))
		text += trimReturn(t)
		// wordwrap :p
		split := strings.Split(text, " ")
		out := split[0] + " "
		var final string
		cod := false
		for i := range split {
			if i > 0 {
				if len(out)+len(split[i])+1+len(codeLoc) > screenWidth && !cod {
					cod = true
					final += out + strings.Repeat(".",
						screenWidth-len(out)-len(codeLoc)) + " " +
						codeLoc + "\n"
					out = indent + split[i] + " "
				} else if len(out)+len(split[i]) > screenWidth {
					out = indent + split[i] + " "
				} else {
					out += split[i] + " "
				}
			}
		}
		final += out
		if !cod {
			rem := screenWidth - len(out) - len(codeLoc)
			if rem < 1 {
				final += "\n" + strings.Repeat(" ", screenWidth-len(codeLoc)) + codeLoc
			} else {
				final += strings.Repeat(".", rem) + " " + codeLoc
			}
		}
		fmt.Println(final)
		if fh != nil {
			out := Entry{time.Now(), level, loc, text}
			j, err := json.Marshal(out)
			if err != nil {
				fmt.Println("logging error:", err)
			}
			_, _ = fmt.Fprintln(fh, string(j)+",")
		}
	}
}

// FileExists reports whether the named file or directory exists.
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// DirectionString is a helper function that returns a string that represents the direction of a connection (inbound or outbound).
func DirectionString(inbound bool) string {
	if inbound {
		return "inbound"
	}
	return "outbound"
}

// PickNoun returns the singular or plural form of a noun depending
// on the count n.
func PickNoun(n int, singular, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}
