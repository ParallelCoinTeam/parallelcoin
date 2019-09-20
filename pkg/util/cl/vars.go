package cl

import (
	"io"
	"os"
	"sync"
)

const (
	_off = iota
	_fatal
	_error
	_warn
	_info
	_debug
	_trace
)

var (
	// Levels is the map of name to level
	Levels = map[string]int{
		"off":   _off,
		"fatal": _fatal,
		"error": _error,
		"warn":  _warn,
		"info":  _info,
		"debug": _debug,
		"trace": _trace,
	}
	// Color turns on and off colouring of error type tag
	Color = true
	// ColorChan accepts a bool and flips the state accordingly
	ColorChan = make(chan bool)
	// ShuttingDown indicates if the shutdown switch has been triggered
	ShuttingDown bool
	// Writer is the place thelogs put out
	Writer = io.MultiWriter(os.Stdout)
	// Og is the root channel that processes logging messages, so, cl.Og <- Fatalf{"format string %s %d", stringy, inty} sends to the root
	Og = make(chan interface{}, 27)
	wg sync.WaitGroup
	// Quit signals the logger to stop. Invoke like this:
	//     close(clog.Quit)
	// You can call init() again to start it up again
	Quit = make(chan struct{})
	// Register is the registry of subsystems in operation
	Register = make(Registry)
	// GlobalLevel is the ostensible global level currently set
	GlobalLevel = "info"
	maxLen      int
	// LogDBC is a channel that can be set to push log messages to another goroutine
	//
	LogDBC chan string
)
