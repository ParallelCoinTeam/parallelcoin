package logg

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var started = time.Now()

// RepoPath is the part of the filesystem path that contains the application
// repository.
//
// todo: can be injected using stroy
var RepoPath = "/home/loki/src/github.com/p9c/pod/"

// Sep is just a convenient shortcut for this very longwinded expression
var Sep = string(os.PathSeparator)

var currentLevel = _info

// writer can be swapped out for any io.*writer* that you want to use instead of
// stdout. Part of the goals of this logging library is to form a small and
// incompressible part of a framework/pattern that takes advantage of the
// universality of []byte and io.Reader/writer - by changing the writer to a
// goroutine that is reading from it to deliver or process the data in some way,
// one can arbitrarily change the logger to send to another server, pass each
// chunk across a channel, or run a closure on it
//
// Obviously, this should be set at startup and not changed except during
// reload/restart cycles.
var writer = &os.Stderr

// SetWriter atomically changes the log writer writer interface
func SetWriter(wr *io.Writer) {
	w := unsafe.Pointer(writer)
	c := unsafe.Pointer(wr)
	atomic.SwapPointer(&w, c)
}

// AllSubsystems stores all of the package subsystem names found in the current
// application
var AllSubsystems []string

// SortAllSubsystems sorts the list of subsystems, to keep the data read-only,
// call this function right at the top of the main, which runs after
// declarations and main/init. Really this is just here to alert the reader.
func SortAllSubsystems() {
	sort.Strings(AllSubsystems)
}

// AddSubsystem adds a subsystem to the list of known subsystems and returns the
// string so it is nice and neat in the package logg.go file
func AddSubsystem() (subsystem string) {
	var pkgPath string
	var split []string
	_, pkgPath, _, _ = runtime.Caller(3)
	fromRoot := strings.Split(pkgPath, RepoPath)[1]
	split = strings.Split(fromRoot, Sep)
	subsystem = strings.Join(split[:len(split)-1], "/")
	AllSubsystems = append(AllSubsystems, subsystem)
	return
}

// highlighted is a text that helps visually distinguish a log entry by category
var highlighted = make(map[string]struct{})
var highlightMx sync.Mutex

// logFilter specifies a set of packages that will not print logs
var logFilter = make(map[string]struct{})
var logFilterMx sync.Mutex

const (
	_off = iota
	_fatal
	_error
	_check
	_info
	_debug
	_trace
)

type levelSpec struct {
	id        int
	name      string
	colorizer func(format string, a ...interface{}) string
}

var levels = []levelSpec{
	{_off, "off", nil},
	{_fatal, "fatal", color.RedString},
	{_error, "error", color.YellowString},
	{_check, "check", color.GreenString},
	{_info, "info", color.WhiteString},
	{_debug, "debug", color.BlueString},
	{_trace, "trace", color.MagentaString},
}

// SetLevel sets the log level via a string, which can be truncated as the first
// letter is unique
func SetLevel(l string) {
	lvl := _info
	for i := range levels {
		if levels[i].name[:len(l)] == l {
			lvl = levels[i].id
		}
	}
	currentLevel = lvl
}

func getLoc(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("%s:%d", file, line)
}

// IsHighlighted returns true if the subsystem is in the list to have attention
// getters added to them
func IsHighlighted(subsystem string) (found bool) {
	highlightMx.Lock()
	_, found = highlighted[subsystem]
	highlightMx.Unlock()
	return
}

// AddHighlighted adds a new subsystem name to the highlighted list
func AddHighlighted(hl string) {
	highlightMx.Lock()
	highlighted[hl] = struct{}{}
	highlightMx.Unlock()
}

// SetHighlighted sets the list of subsystems to highlight
func SetHighlighted(highlights []string) (found bool) {
	highlightMx.Lock()
	highlighted = make(map[string]struct{}, len(highlights))
	for i := range highlights {
		highlighted[highlights[i]] = struct{}{}
	}
	highlightMx.Unlock()
	return
}

// GetHighlighted returns a copy of the map of highlighted subsystems
func GetHighlighted() (o []string) {
	highlightMx.Lock()
	o = make([]string, len(logFilter))
	var counter int
	for i := range logFilter {
		o[counter] = i
		counter++
	}
	highlightMx.Unlock()
	sort.Strings(o)
	return
}

// SetFiltered sets the list of subsystems to filter
func SetFiltered(filter []string) {
	logFilterMx.Lock()
	logFilter = make(map[string]struct{}, len(filter))
	for i := range filter {
		logFilter[filter[i]] = struct{}{}
	}
	logFilterMx.Unlock()
}

// IsFiltered returns true if the subsystem should not print logs
func IsFiltered(subsystem string) (found bool) {
	logFilterMx.Lock()
	_, found = logFilter[subsystem]
	logFilterMx.Unlock()
	return
}

// GetFiltered returns a copy of the map of filtered subsystems
func GetFiltered() (o []string) {
	logFilterMx.Lock()
	o = make([]string, len(logFilter))
	var counter int
	for i := range logFilter {
		o[counter] = i
		counter++
	}
	logFilterMx.Unlock()
	sort.Strings(o)
	return
}

// join constructs a string from an slice of interface same as Println but
// without the terminal newline
func join(sep string, a ...interface{}) (o string) {
	for i := range a {
		o += fmt.Sprint(a[i])
		if i < len(a)-1 {
			o += sep
		}
	}
	return
}

// PrintLn output joins a list of various kinds of typed variables
func PrintLn(level int, subsystem string, a ...interface{}) {
	if level >= currentLevel || !IsFiltered(subsystem) {
		if IsHighlighted(subsystem) {
			subsystem = " (((" + strings.ToUpper(subsystem) + ")))"
		} else {
			subsystem = ""
		}
		// set the runtime.Caller depth to appropriate for the caller
		depth := 3
		if level == _check {
			depth = 4
		}
		fmt.Fprintln(
			*writer,
			levels[level].colorizer(
				"%v %s%s [%s] %s",
				time.Now().Truncate(time.Millisecond).Sub(started),
				levels[level].name,
				subsystem,
				color.WhiteString(join(" ", a)),
				getLoc(depth),
			),
		)
	}
}

// PrintF is a printf type function, uses Sprintf so format is the same
func PrintF(level int, subsystem string, format string, a ...interface{}) {
	if level >= currentLevel || !IsFiltered(subsystem) {
		if IsHighlighted(subsystem) {
			subsystem = " (((" + strings.ToUpper(subsystem) + ")))"
		} else {
			subsystem = ""
		}
		fmt.Fprintf(
			*writer, format,
			levels[level].colorizer(
				"%v %s%s [%s] %s",
				time.Now().Truncate(time.Millisecond).Sub(started),
				levels[level].name,
				subsystem,
				color.WhiteString(join(" ", a)),
				getLoc(3),
			),
		)
	}
}

// PrintS spews the given list of variables
func PrintS(level int, subsystem string, a ...interface{}) {
	if level >= currentLevel || !IsFiltered(subsystem) {
		if IsHighlighted(subsystem) {
			subsystem = " (((" + strings.ToUpper(subsystem) + ")))"
		} else {
			subsystem = ""
		}
		fmt.Fprintln(
			*writer,
			levels[level].colorizer(
				"%v %s%s\n%s\n%s",
				time.Now().Truncate(time.Millisecond).Sub(started),
				levels[level].name,
				subsystem,
				color.WhiteString(join("\n", a)),
				getLoc(3),
			),
		)
	}
}

// PrintC runs a closure to generate a log string so as to shift this processing
// out of a goroutine for such as heavy tracing output
func PrintC(level int, subsystem string, closure func() string) {
	if level >= currentLevel || !IsFiltered(subsystem) {
		if IsHighlighted(subsystem) {
			subsystem = " (((" + strings.ToUpper(subsystem) + ")))"
		} else {
			subsystem = ""
		}
		PrintLn(level, subsystem, false, closure())
	}
}

// Check evaluates an error value and returns true if it is not nil
func Check(err error, subsystem string, a ...interface{}) bool {
	if err != nil {
		if IsHighlighted(subsystem) {
			subsystem = " (((" + strings.ToUpper(subsystem) + ")))"
		} else {
			subsystem = ""
		}
		if !IsFiltered(subsystem) {
			PrintLn(_check, subsystem, a)
		}
		return true
	}
	return false
}
