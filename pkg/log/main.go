package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/buger/goterm"
)

var (
	colorOff     = "\033[0m"
	colorRed     = "\033[0;31m"
	colorRedB    = "\033[0;31;1m"
	colorGreen   = "\033[0;32m"
	colorGreenB  = "\033[0;32;1m"
	colorOrange  = "\033[0;33m"
	colorOrangeB = "\033[0;33;1m"
	colorBlue    = "\033[0;34m"
	colorBlueB   = "\033[0;34;1m"
	colorPurple  = "\033[0;35m"
	colorPurpleB = "\033[0;35;1m"
	colorCyan    = "\033[0;36m"
	colorCyanB   = "\033[0;36;1m"
	colorGray    = "\033[0;37m"
	colorGrayB   = "\033[0;37;1m"
	styleBold    = "\033[1m"
)

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
	Color         bool
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
	found := false
	for i := range Levels {
		if level == Levels[i] {
			found = true
			break
		}
	}
	if !found {
		level = "info"
	}
	return level
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
func (l *Logger) SetLevel(level string, color bool) {
	_, loc, line, _ := runtime.Caller(1)
	files := strings.Split(loc, "github.com/parallelcointeam/parallelcoin/")
	codeLoc := fmt.Sprint(files[1], ":", rightJustifyLineNumber(line))
	fmt.Println("setting level to", level, codeLoc)
	*l = *Empty()
	var fallen bool
	switch {
	case level == Trace || fallen:
		// fmt.Println("loading Trace printers")
		l.Trace = Println("TRACE ", color, l.LogFileHandle)
		l.Tracef = Printf("TRACE ", color, l.LogFileHandle)
		l.Tracec = Printc("TRACE ", color, l.LogFileHandle)
		fallen = true
		fallthrough
	case level == Debug || fallen:
		// fmt.Println("loading Debug printers")
		l.Debug = Println("DEBUG ", color, l.LogFileHandle)
		l.Debugf = Printf("DEBUG ", color, l.LogFileHandle)
		l.Debugc = Printc("DEBUG ", color, l.LogFileHandle)
		fallen = true
		fallthrough
	case level == Info || fallen:
		// fmt.Println("loading Info printers")
		l.Info = Println(" INFO ", color, l.LogFileHandle)
		l.Infof = Printf(" INFO ", color, l.LogFileHandle)
		l.Infoc = Printc(" INFO ", color, l.LogFileHandle)
		fallen = true
		fallthrough
	case level == Warn || fallen:
		// fmt.Println("loading Warn printers")
		l.Warn = Println(" WARN ", color, l.LogFileHandle)
		l.Warnf = Printf(" WARN ", color, l.LogFileHandle)
		l.Warnc = Printc(" WARN ", color, l.LogFileHandle)
		fallen = true
		fallthrough
	case level == Error || fallen:
		// fmt.Println("loading Error printers")
		l.Error = Println("ERROR ", color, l.LogFileHandle)
		l.Errorf = Printf("ERROR ", color, l.LogFileHandle)
		l.Errorc = Printc("ERROR ", color, l.LogFileHandle)
		fallen = true
		fallthrough
	case level == Fatal:
		// fmt.Println("loading Fatal printers")
		l.Fatal = Println("FATAL ", color, l.LogFileHandle)
		l.Fatalf = Printf("FATAL ", color, l.LogFileHandle)
		l.Fatalc = Printc("FATAL ", color, l.LogFileHandle)
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
	if s[len(s)-1] == '\n' {
		return s[:len(s)-1]
	}
	return s
}

func rightJustifyLineNumber(n int) string {
	s := fmt.Sprint(n)
	switch len(s) {
	case 1:
		s += "    "
	case 2:
		s += "   "
	case 3:
		s += "  "
	case 4:
		s += " "
	}
	return s
}

// RightJustify takes a string and right justifies it by a width or crops it
func rightJustify(s string, w int) string {
	sw := len(s)
	diff := w - sw
	if diff > 0 {
		s = strings.Repeat(" ", diff) + s
	} else if diff < 0 {
		s = s[:w]
	}
	return s
}

// Println prints a log entry like Println
func Println(level string, color bool, fh *os.File) func(a ...interface{}) {
	return func(a ...interface{}) {
		if color {
			switch strings.Trim(level, " ") {
			case "FATAL":
				level = colorPurple + level + colorOff
			case "ERROR":
				level = colorRed + level + colorOff
			case "WARN":
				level = colorOrange + level + colorOff
			case "INFO":
				level = colorGray + level + colorOff
			case "DEBUG":
				level = colorBlue + level + colorOff
			case "TRACE":
				level = colorCyan + level + colorOff
			}
		}
		terminalWidth := goterm.Width() - 3
		_, loc, line, _ := runtime.Caller(2)
		files := strings.Split(loc, "github.com/parallelcointeam/parallelcoin/")
		since := fmt.Sprint(time.Now().Sub(StartupTime) / time.
			Second * time.Second)
		var prefix string = level + " " + rightJustify(since, 9) + " "
		codeLoc := fmt.Sprint(files[1], ":", rightJustifyLineNumber(line))
		indent := strings.Repeat(" ", len(prefix))
		ellipsis := " "
		if terminalWidth > 160 {
			prefix = level + " " + rightJustify(fmt.Sprint(time.Now()), 24) + " "
			ellipsis = "."
		}
		if terminalWidth < 64 {
			prefix = ""
			indent = "  "
			ellipsis = " "
		}
		if terminalWidth < 56 {
			codeLoc = ""
		}
		text := trimReturn(fmt.Sprintln(a...))
		// wordwrap :p
		split := strings.Split(text, " ")
		out := prefix + split[0] + " "
		var final string
		cod := false
		for i := range split {
			if i > 0 {
				if len(out)+len(split[i])+1+len(
					codeLoc) > terminalWidth && !cod {
					cod = true
					final += out + strings.Repeat(ellipsis,
						terminalWidth-len(out)-len(codeLoc)) + " " +
						codeLoc + "\n"
					out = indent + split[i] + " "
				} else if len(out)+len(split[i]) > terminalWidth {
					final += out + "\n"
					out = indent + split[i] + " "
				} else {
					out += split[i] + " "
				}
			}
		}
		final += out
		if !cod {
			rem := terminalWidth - len(out) - len(codeLoc)
			if rem < 1 {
				final += strings.Repeat(" ",
					terminalWidth-len(codeLoc)) + codeLoc
			} else {
				final += strings.Repeat(ellipsis, rem) + " " + codeLoc
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
func Printf(level string, color bool, fh *os.File) func(format string,
	a ...interface{}) {
	return func(format string, a ...interface{}) {
		if color {
			switch strings.Trim(level, " ") {
			case "FATAL":
				level = colorPurple + level + colorOff
			case "ERROR":
				level = colorRed + level + colorOff
			case "WARN":
				level = colorOrange + level + colorOff
			case "INFO":
				level = colorGray + level + colorOff
			case "DEBUG":
				level = colorBlue + level + colorOff
			case "TRACE":
				level = colorCyan + level + colorOff
			}
		}
		terminalWidth := goterm.Width() - 3
		_, loc, line, _ := runtime.Caller(2)
		files := strings.Split(loc, "github.com/parallelcointeam/parallelcoin/")
		since := fmt.Sprint(time.Now().Sub(StartupTime) / time.
			Second * time.Second)
		var prefix string = level + " " + rightJustify(since, 9) + " "
		codeLoc := fmt.Sprint(files[1], ":", rightJustifyLineNumber(line))
		indent := strings.Repeat(" ", len(prefix))
		ellipsis := " "
		if terminalWidth > 160 {
			prefix = level + " " + rightJustify(fmt.Sprint(time.Now()), 24) + " "
			ellipsis = "."
		}
		if terminalWidth < 64 {
			prefix = ""
			indent = "  "
			ellipsis = " "
		}
		if terminalWidth < 56 {
			codeLoc = ""
		}
		text := trimReturn(fmt.Sprintf(format, a...))
		// wordwrap :p
		split := strings.Split(text, " ")
		out := prefix + split[0] + " "
		var final string
		cod := false
		for i := range split {
			if i > 0 {
				if len(out)+len(split[i])+1+len(
					codeLoc) > terminalWidth && !cod {
					cod = true
					final += out + strings.Repeat(ellipsis,
						terminalWidth-len(out)-len(codeLoc)) + " " +
						codeLoc + "\n"
					out = indent + split[i] + " "
				} else if len(out)+len(split[i]) > terminalWidth {
					final += out + "\n"
					out = indent + split[i] + " "
				} else {
					out += split[i] + " "
				}
			}
		}
		final += out
		if !cod {
			rem := terminalWidth - len(out) - len(codeLoc)
			if rem < 1 {
				final += strings.Repeat(" ",
					terminalWidth-len(codeLoc)) + codeLoc
			} else {
				final += strings.Repeat(ellipsis, rem) + " " + codeLoc
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
func Printc(level string, color bool, fh *os.File) func(fn func() string) {
	// level = strings.ToUpper(string(level[0]))
	return func(fn func() string) {
		if color {
			switch strings.Trim(level, " ") {
			case "FATAL":
				level = colorPurple + level + colorOff
			case "ERROR":
				level = colorRed + level + colorOff
			case "WARN":
				level = colorOrange + level + colorOff
			case "INFO":
				level = colorGray + level + colorOff
			case "DEBUG":
				level = colorBlue + level + colorOff
			case "TRACE":
				level = colorCyan + level + colorOff
			}
		}
		terminalWidth := goterm.Width() - 3
		_, loc, line, _ := runtime.Caller(2)
		files := strings.Split(loc, "github.com/parallelcointeam/parallelcoin/")
		since := fmt.Sprint(time.Now().Sub(StartupTime) / time.
			Second * time.Second)
		var prefix string = level + " " + (rightJustify(since, 9)) + " "
		codeLoc := fmt.Sprint(files[1], ":", rightJustifyLineNumber(line))
		indent := strings.Repeat(" ", len(prefix))
		ellipsis := " "
		if terminalWidth > 160 {
			prefix = level + " " + rightJustify(fmt.Sprint(time.Now()), 24) + " "
			ellipsis = "."
		}
		if terminalWidth < 64 {
			prefix = ""
			indent = "  "
			ellipsis = " "
		}
		if terminalWidth < 56 {
			codeLoc = ""
		}
		t := fn()
		text := trimReturn(t)
		split := strings.Split(text, " ")
		out := prefix + split[0] + " "
		var final string
		cod := false
		for i := range split {
			if i > 0 {
				if len(out)+len(split[i])+1+len(
					codeLoc) > terminalWidth && !cod {
					cod = true
					final += out + strings.Repeat(ellipsis,
						terminalWidth-len(out)-len(codeLoc)) + " " +
						codeLoc + "\n"
					out = indent + split[i] + " "
				} else if len(out)+len(split[i]) > terminalWidth {
					final += out + "\n"
					out = indent + split[i] + " "
				} else {
					out += split[i] + " "
				}
			}
		}
		final += out
		if !cod {
			rem := terminalWidth - len(out) - len(codeLoc)
			if rem < 1 {
				final += strings.Repeat(" ",
					terminalWidth-len(codeLoc)) + codeLoc
			} else {
				final += strings.Repeat(ellipsis, rem) + " " + codeLoc
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
