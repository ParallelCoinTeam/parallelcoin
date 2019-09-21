package log

import (
	"fmt"
	"sync"
)

// Register is the registry of subsystems in operation
var Register = make(Registry)
var globalLevel string
var mutex sync.Mutex

// Registry is the centralised store that links all the loggers so they  can
// be accessed programmatically
type Registry map[string]*Logger

func // Add appends a new subsystem to its map for access and introspection
(r *Registry) Add(s *Logger) {
	mutex.Lock()
	_, ok := (*r)[s.Name]
	if !ok {
		(*r)[s.Name] = s
	}
	mutex.Unlock()
}

func // List returns a string slice containing all the registered loggers
(r *Registry) List() (out []string) {
	mutex.Lock()
	for i := range *r {
		fmt.Println("List", (*r)[i].Name)
		out = append(out, (*r)[i].Name)
	}
	mutex.Unlock()
	return
}

func // Get returns the subsystem.
// This could then be used to close or set its level eg `*r.Get("subsystem").
// SetLevel("debug")`
(r *Registry) Get(name string) (out *Logger) {
	var ok bool
	mutex.Lock()
	if out, ok = (*r)[name]; ok {
		mutex.Unlock()
	}
	return
}

func // GetGlobalLevel returns the global level
(r *Registry) GetGlobalLevel() string {
	mutex.Lock()
	out := globalLevel
	mutex.Lock()
	return out
}

func // SetAllLevels sets the level in all registered loggers
(r *Registry) SetAllLevels(level string) {
	fmt.Println("levels", level)
	level = sanitizeLoglevel(level)
	fmt.Println(level)
	mutex.Lock()
	globalLevel = level
	fmt.Println(globalLevel)
	mutex.Unlock()
	loggers := r.List()
	for _, x := range loggers {
		fmt.Println(x)
		r.Get(x).SetLevel(level)
	}
}
