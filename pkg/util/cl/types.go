package cl

import "sync"

type (
	// StringClosure is a function that returns a string, used to defer execution of expensive logging operations
	StringClosure func() string
	// Value is the generic list of things processed by the log chan
	Value []interface{}
	// Fatal is a log value that indicates level and how to interpret the interface slice
	Fatal Value
	// Fatalf is a log value that indicates level and how to interpret the interface slice
	Fatalf Value
	// Ftl is a log type that is just one string
	Ftl string
	// Fatalc is for passing a closure when the log entry is expensive to compute
	Fatalc StringClosure
	// Error is a log value that indicates level and how to interpret the interface slice
	Error Value
	// Errorf is a log value that indicates level and how to interpret the interface slice
	Errorf Value
	// Err is a log type that is just one string
	Err string
	// Errorc is for passing a closure when the log entry is expensive to compute
	Errorc StringClosure
	// Warn is a log value that indicates level and how to interpret the interface slice
	Warn Value
	// Warnf is a log value that indicates level and how to interpret the interface slice
	Warnf Value
	// Wrn is a log type that is just one string
	Wrn string
	// Warnc is for passing a closure when the log entry is expensive to compute
	Warnc StringClosure
	// Info is a log value that indicates level and how to interpret the interface slice
	Info Value
	// Infof is a log value that indicates level and how to interpret the interface slice
	Infof Value
	// Inf is a log type that is just one string
	Inf string
	// Infoc is for passing a closure when the log entry is expensive to compute
	Infoc StringClosure
	// Debug is a log value that indicates level and how to interpret the interface slice
	Debug Value
	// Debugf is a log value that indicates level and how to interpret the interface slice
	Debugf Value
	// Dbg is a log type that is just one string
	Dbg string
	// Debugc is for passing a closure when the log entry is expensive to compute
	Debugc StringClosure
	// Trace is a log value that indicates level and how to interpret the interface slice
	Trace Value
	// Tracef is a log value that indicates level and how to interpret the interface slice
	Tracef Value
	// Trc is a log type that is just one string
	Trc string
	// Tracec is for passing a closure when the log entry is expensive to compute
	Tracec StringClosure
)

// A SubSystem is a logger with a specific prefix name prepended  to the entry
type SubSystem struct {
	Name        string
	Ch          chan interface{}
	Level       int
	LevelString string
	MaxLen      int
	mutex       sync.Mutex
}

// Registry is
type Registry map[string]*SubSystem
