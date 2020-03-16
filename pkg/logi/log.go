package logi

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/p9c/goterm"
)

const (
	Off   = "off"
	Fatal = "fatal"
	Error = "error"
	Warn  = "warn"
	Info  = "info"
	Check = "check"
	Debug = "debug"
	Trace = "trace"
)

var (
	// NoClosure is a noop for a closure print function
	NoClosure = func() PrintcFunc {
		f := func(_ func() string) {}
		return f
	}
	// NoPrintf is a noop for a closure printf function
	NoPrintf = func() PrintfFunc {
		f := func(_ string, _ ...interface{}) {}
		return f
	}
	// NoPrintln is a noop for a println function
	NoPrintln = func() PrintlnFunc {
		f := func(_ ...interface{}) {}
		return f
	}
	// NoPrintln is a noop for a println function
	NoCheck = func() CheckFunc {
		f := func(_ error) bool {
			return true
		}
		return f
	}
	// NoSpew is a noop for a spew function
	NoSpew = func() SpewFunc {
		f := func(_ interface{}) {}
		return f
	}
	// StartupTime allows a shorter log prefix as time since start
	StartupTime    = time.Now()
	BackgroundGrey = "\u001b[48;5;240m"
	ColorBlue      = "\u001b[38;5;33m"
	ColorBold      = "\u001b[1m"
	ColorBrown     = "\u001b[38;5;130m"
	ColorCyan      = "\u001b[36m"
	ColorFaint     = "\u001b[2m"
	ColorGreen     = "\u001b[38;5;40m"
	ColorItalic    = "\u001b[3m"
	ColorOff       = "\u001b[0m"
	ColorOrange    = "\u001b[38;5;208m"
	ColorPurple    = "\u001b[38;5;99m"
	ColorRed       = "\u001b[38;5;196m"
	ColorUnderline = "\u001b[4m"
	ColorViolet    = "\u001b[38;5;201m"
	ColorYellow    = "\u001b[38;5;226m"

	L = Empty()

	Levels = []string{
		Off, Fatal, Error, Warn, Info, Check, Debug, Trace,
	}
)

var wr LogWriter

// Entry is a log entry to be printed as json to the log file
type Entry struct {
	Time         time.Time
	Level        string
	CodeLocation string
	Text         string
}

type Logger struct {
	Fatal         PrintlnFunc
	Error         PrintlnFunc
	Warn          PrintlnFunc
	Info          PrintlnFunc
	Check         CheckFunc
	Debug         PrintlnFunc
	Trace         PrintlnFunc
	Fatalf        PrintfFunc
	Errorf        PrintfFunc
	Warnf         PrintfFunc
	Infof         PrintfFunc
	Checkf        CheckFunc
	Debugf        PrintfFunc
	Tracef        PrintfFunc
	Fatalc        PrintcFunc
	Errorc        PrintcFunc
	Warnc         PrintcFunc
	Infoc         PrintcFunc
	Checkc        CheckFunc
	Debugc        PrintcFunc
	Tracec        PrintcFunc
	Fatals        SpewFunc
	Errors        SpewFunc
	Warns         SpewFunc
	Infos         SpewFunc
	Debugs        SpewFunc
	Traces        SpewFunc
	LogFileHandle *os.File
	Writer        LogWriter
	Color         bool
	Split         string
	// If this channel is loaded log entries are composed and sent to it
	LogChan chan Entry
}

// SetLevel enables or disables the various print functions
func (l *Logger) SetLevel(level string, color bool, split string) *Logger {
	l.Split = split + string(os.PathSeparator)
	level = sanitizeLoglevel(level)
	var fallen bool
	switch {
	case level == Trace || fallen:
		l.Trace = printlnFunc("TRC", color, l.LogFileHandle, nil, l.Split)
		l.Tracef = printfFunc("TRC", color, l.LogFileHandle, nil, l.Split)
		l.Tracec = printcFunc("TRC", color, l.LogFileHandle, nil, l.Split)
		l.Traces = ps("TRC", color, l.LogFileHandle, l.Split)
		fallen = true
		fallthrough
	case level == Debug || fallen:
		l.Debug = printlnFunc("DBG", color, l.LogFileHandle, l.LogChan, l.Split)
		l.Debugf = printfFunc("DBG", color, l.LogFileHandle, l.LogChan, l.Split)
		l.Debugc = printcFunc("DBG", color, l.LogFileHandle, l.LogChan, l.Split)
		l.Debugs = ps("DBG", color, l.LogFileHandle, l.Split)
		fallen = true
		fallthrough
	case level == Check || fallen:
		l.Check = checkFunc(color, l.LogFileHandle, l.LogChan, l.Split)
		l.Checkf = checkFunc(color, l.LogFileHandle, l.LogChan, l.Split)
		l.Checkc = checkFunc(color, l.LogFileHandle, l.LogChan, l.Split)
		fallen = true
		fallthrough
	case level == Info || fallen:
		l.Info = printlnFunc("INF", color, l.LogFileHandle, l.LogChan, l.Split)
		l.Infof = printfFunc("INF", color, l.LogFileHandle, l.LogChan, l.Split)
		l.Infoc = printcFunc("INF", color, l.LogFileHandle, l.LogChan, l.Split)
		l.Infos = ps("INF", color, l.LogFileHandle, l.Split)
		fallen = true
		fallthrough
	case level == Warn || fallen:
		l.Warn = printlnFunc("WRN", color, l.LogFileHandle, l.LogChan, l.Split)
		l.Warnf = printfFunc("WRN", color, l.LogFileHandle, l.LogChan, l.Split)
		l.Warnc = printcFunc("WRN", color, l.LogFileHandle, l.LogChan, l.Split)
		l.Warns = ps("WRN", color, l.LogFileHandle, l.Split)
		fallen = true
		fallthrough
	case level == Error || fallen:
		l.Error = printlnFunc("ERR", color, l.LogFileHandle, l.LogChan, l.Split)
		l.Errorf = printfFunc("ERR", color, l.LogFileHandle, l.LogChan, l.Split)
		l.Errorc = printcFunc("ERR", color, l.LogFileHandle, l.LogChan, l.Split)
		l.Errors = ps("ERR", color, l.LogFileHandle, l.Split)
		fallen = true
		fallthrough
	case level == Fatal:
		l.Fatal = printlnFunc("FTL", color, l.LogFileHandle, l.LogChan, l.Split)
		l.Fatalf = printfFunc("FTL", color, l.LogFileHandle, l.LogChan, l.Split)
		l.Fatalc = printcFunc("FTL", color, l.LogFileHandle, l.LogChan, l.Split)
		l.Fatals = ps("FTL", color, l.LogFileHandle, l.Split)
		fallen = true
	}
	return l
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
			wr.Println("error rotating log", err)
			return
		}
	}
	logFileHandle, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		wr.Println("error opening log file", logFileName)
	}
	l.LogFileHandle = logFileHandle
	_, _ = fmt.Fprintln(logFileHandle, "{")
}

type LogWriter struct {
	io.Writer
}

func (w *LogWriter) Print(a ...interface{}) {
	_, _ = fmt.Fprint(wr, a...)
}

func (w *LogWriter) Printf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(wr, format, a...)
}

func (w *LogWriter) Println(a ...interface{}) {
	_, _ = fmt.Fprintln(wr, a...)
}

type PrintcFunc func(func() string)
type PrintfFunc func(format string, a ...interface{})
type PrintlnFunc func(a ...interface{})
type CheckFunc func(err error) bool
type SpewFunc func(interface{})

func Composite(text, level string, color bool, split string) string {
	dots := "."
	terminalWidth := goterm.Width()
	if terminalWidth <= 80 {
		terminalWidth = 80
	}
	skip := 2
	if level == Check {
		skip = 3
	}
	_, loc, iline, _ := runtime.Caller(skip)
	line := fmt.Sprint(iline)
	files := strings.Split(loc, split)
	var file, since string
	file = loc
	if len(files) > 1 {
		file = files[1]
	}
	switch {
	case terminalWidth <= 60:
		since = ""
		file = ""
		line = ""
		dots = " "
	case terminalWidth <= 80:
		dots = " "
		if len(file) > 30 {
			file = ""
			line = ""
		}
		since = fmt.Sprintf("%v", time.Now().Sub(StartupTime)/time.Second*time.Second)
	case terminalWidth < 120:
		if len(file) > 40 {
			file = ""
			line = ""
			dots = " "
		}
		since = fmt.Sprintf("%v", time.Now().Sub(StartupTime)/time.Millisecond*time.Millisecond)
	case terminalWidth < 160:
		if len(file) > 60 {
			file = ""
			line = ""
			dots = " "
		}
		since = fmt.Sprint(time.Now())[:19]
	case terminalWidth >= 200:
		since = fmt.Sprint(time.Now())[:39]
	default:
		since = fmt.Sprint(time.Now())[:19]
	}
	levelLen := len(level) + 1
	sinceLen := len(since) + 1
	textLen := len(text) + 1
	fileLen := len(file) + 1
	lineLen := len(line) + 1
	if file != "" {
		file += ":"
	}
	if color {
		switch level {
		case "FTL":
			level = ColorBold + ColorRed + level + ColorOff
			since = ColorRed + since + ColorOff
			file = ColorItalic + ColorBlue + file
			line = line + ColorOff
		case "ERR":
			level = ColorBold + ColorOrange + level + ColorOff
			since = ColorOrange + since + ColorOff
			file = ColorItalic + ColorBlue + file
			line = line + ColorOff
		case "WRN":
			level = ColorBold + ColorYellow + level + ColorOff
			since = ColorYellow + since + ColorOff
			file = ColorItalic + ColorBlue + file
			line = line + ColorOff
		case "INF":
			level = ColorBold + ColorGreen + level + ColorOff
			since = ColorGreen + since + ColorOff
			file = ColorItalic + ColorBlue + file
			line = line + ColorOff
		case "CHK":
			level = ColorBold + ColorCyan + level + ColorOff
			since = since
			file = ColorItalic + ColorBlue + file
			line = line + ColorOff
		case "DBG":
			level = ColorBold + ColorBlue + level + ColorOff
			since = ColorBlue + since + ColorOff
			file = ColorItalic + ColorBlue + file
			line = line + ColorOff
		case "TRC":
			level = ColorBold + ColorViolet + level + ColorOff
			since = ColorViolet + since + ColorOff
			file = ColorItalic + ColorBlue + file
			line = line + ColorOff
		}
	}
	final := ""
	if levelLen+sinceLen+textLen+fileLen+lineLen > terminalWidth {
		lines := strings.Split(text, "\n")
		// log text is multiline
		line1len := terminalWidth - levelLen - sinceLen - fileLen - lineLen
		restLen := terminalWidth - levelLen - sinceLen
		if len(lines) > 1 {
			final = fmt.Sprintf("%s %s %s %s%s", level, since,
				strings.Repeat(" ",
					terminalWidth-levelLen-sinceLen-fileLen-lineLen),
				file, line)
			final += text[:len(text)-1]
		} else {
			// log text is a long line
			spaced := strings.Split(text, " ")
			var rest bool
			curLineLen := 0
			final += fmt.Sprintf("%s %s ", level, since)
			var i int
			for i = range spaced {
				if i > 0 {
					curLineLen += len(spaced[i-1]) + 1
					if !rest {
						if curLineLen >= line1len {
							rest = true
							spacers := terminalWidth - levelLen - sinceLen -
								fileLen - lineLen - curLineLen + len(spaced[i-1]) + 1
							if spacers < 1 {
								spacers = 1
							}
							final += strings.Repeat(dots, spacers)
							final += fmt.Sprintf(" %s%s\n",
								file, line)
							final += strings.Repeat(" ", levelLen+sinceLen)
							final += spaced[i-1] + " "
							curLineLen = len(spaced[i-1]) + 1
						} else {
							final += spaced[i-1] + " "
						}
					} else {
						if curLineLen >= restLen-1 {
							final += "\n" + strings.Repeat(" ",
								levelLen+sinceLen)
							final += spaced[i-1] + dots
							curLineLen = len(spaced[i-1]) + 1
						} else {
							final += spaced[i-1] + " "
						}
					}
				}
			}
			curLineLen += len(spaced[i])
			if !rest {
				if curLineLen >= line1len {
					final += fmt.Sprintf("%s %s%s\n",
						strings.Repeat(dots,
							len(spaced[i])+line1len-curLineLen),
						file, line)
					final += strings.Repeat(" ", levelLen+sinceLen)
					final += spaced[i] // + "\n"
				} else {
					final += fmt.Sprintf("%s %s %s%s\n",
						spaced[i],
						strings.Repeat(dots,
							terminalWidth-curLineLen-fileLen-lineLen),
						file, line)
				}
			} else {
				if curLineLen >= restLen {
					final += "\n" + strings.Repeat(" ", levelLen+sinceLen)
				}
				final += spaced[i]
			}
		}
	} else {
		final = fmt.Sprintf("%s %s %s %s %s%s", level, since, text,
			strings.Repeat(dots,
				terminalWidth-levelLen-sinceLen-textLen-fileLen-lineLen),
			file, line)
	}
	return final
}

// printlnFunc prints a log entry like Println
func DirectionString(inbound bool) string {
	if inbound {
		return "inbound"
	}
	return "outbound"
}

// PickNoun returns the singular or plural form of a noun depending
// on the count n.
func Empty() *Logger {
	return &Logger{
		Fatal:  NoPrintln(),
		Error:  NoPrintln(),
		Warn:   NoPrintln(),
		Info:   NoPrintln(),
		Check:  NoCheck(),
		Debug:  NoPrintln(),
		Trace:  NoPrintln(),
		Fatalf: NoPrintf(),
		Errorf: NoPrintf(),
		Warnf:  NoPrintf(),
		Infof:  NoPrintf(),
		Checkf: NoCheck(),
		Debugf: NoPrintf(),
		Tracef: NoPrintf(),
		Fatalc: NoClosure(),
		Errorc: NoClosure(),
		Warnc:  NoClosure(),
		Infoc:  NoClosure(),
		Checkc: NoCheck(),
		Debugc: NoClosure(),
		Tracec: NoClosure(),
		Fatals: NoSpew(),
		Errors: NoSpew(),
		Warns:  NoSpew(),
		Infos:  NoSpew(),
		Debugs: NoSpew(),
		Traces: NoSpew(),
		Writer: wr,
	}

}

// sanitizeLoglevel accepts a string and returns a
// default if the input is not in the Levels slice
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// DirectionString is a helper function that returns a string that represents the direction of a connection (inbound or outbound).
func init() {
	SetLogWriter(os.Stderr)
	L.SetLevel("info", true, "log/")
	L.Trace("starting up logger")
}

func PickNoun(n int, singular, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}
func Print(a ...interface{}) {
	wr.Print(a...)
}

func printcFunc(level string, color bool, fh *os.File, ch chan Entry, split string) PrintcFunc {
	f := func(fn func() string) {
		t := fn()
		text := trimReturn(t)
		wr.Println(Composite(text, level, color, split))
		if fh != nil || ch != nil {
			_, loc, line, _ := runtime.Caller(2)
			out := Entry{time.Now(), level, fmt.Sprint(loc, ":", line), text}
			if fh != nil {
				j, err := json.Marshal(out)
				if err != nil {
					wr.Println("logging error:", err)
				}
				_, _ = fmt.Fprint(fh, string(j)+",")
			}
			if ch != nil {
				ch <- out
			}
		}
	}
	return f
}

// ps spews a variable
func Printf(format string, a ...interface{}) {
	wr.Printf(format, a...)
}

// Logger is a struct containing all the functions with nice handy names
func printfFunc(level string, color bool, fh *os.File, ch chan Entry, split string) PrintfFunc {
	f := func(format string, a ...interface{}) {
		text := fmt.Sprintf(format, a...)
		wr.Println(Composite(text, level, color, split))
		if fh != nil || ch != nil {
			_, loc, line, _ := runtime.Caller(2)
			out := Entry{time.Now(), level, fmt.Sprint(loc, ":", line), text}
			if fh != nil {
				j, err := json.Marshal(out)
				if err != nil {
					wr.Println("logging error:", err)
				}
				_, _ = fmt.Fprint(fh, string(j)+",")
			}
			if ch != nil {
				ch <- out
			}
		}
	}
	return f
}

// printcFunc prints from a closure returning a string
func Println(a ...interface{}) {
	wr.Println(a...)
}

func printlnFunc(level string, color bool, fh *os.File, ch chan Entry, split string) PrintlnFunc {
	f := func(a ...interface{}) {
		text := trimReturn(fmt.Sprintln(a...))
		wr.Println(Composite(text, level, color, split))
		if fh != nil || ch != nil {
			_, loc, line, _ := runtime.Caller(2)
			out := Entry{time.Now(), level, fmt.Sprint(loc, ":", line), text}
			if fh != nil {
				j, err := json.Marshal(out)
				if err != nil {
					wr.Println("logging error:", err)
				}
				_, _ = fmt.Fprint(fh, string(j)+",")
			}
			if ch != nil {
				ch <- out
			}
		}
	}
	return f
}

func checkFunc(color bool, fh *os.File, ch chan Entry, split string) CheckFunc {
	f := func(err error) bool {
		n := err == nil
		if n {
			return false
		}
		text := err.Error()
		wr.Println(Composite(text, "CHK", color, split))
		if fh != nil || ch != nil {
			_, loc, line, _ := runtime.Caller(3)
			out := Entry{time.Now(), "CHK", fmt.Sprint(loc, ":", line), text}
			if fh != nil {
				j, err := json.Marshal(out)
				if err != nil {
					wr.Println("logging error:", err)
				}
				_, _ = fmt.Fprint(fh, string(j)+",")
			}
			if ch != nil {
				ch <- out
			}
		}
		return true
	}
	return f
}

// printfFunc prints a log entry with formatting
func ps(level string, color bool, fh *os.File, split string) SpewFunc {
	f := func(a interface{}) {
		text := trimReturn(spew.Sdump(a))
		o := "" + Composite("spew:", level, color, split)
		o += "\n" + text + "\n"
		wr.Print(o)
		if fh != nil {
			_, loc, line, _ := runtime.Caller(2)
			out := Entry{time.Now(), level, fmt.Sprint(loc, ":", line), text}
			j, err := json.Marshal(out)
			if err != nil {
				wr.Println("logging error:", err)
			}
			_, _ = fmt.Fprint(fh, string(j)+",")
		}
	}
	return f
}

// FileExists reports whether the named file or directory exists.
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

func SetLogWriter(w io.Writer) {
	wr.Writer = w
}

func trimReturn(s string) string {
	if s[len(s)-1] == '\n' {
		return s[:len(s)-1]
	}
	return s
}
