package cl

// Add appends a new subsystem to its map for access and introspeection
func (r *Registry) Add(s *SubSystem) {
	_, ok := (*r)[s.Name]
	if ok {
		Og <- Error{s.Name, "subsystem already registered", Ine3()}
	} else {
		(*r)[s.Name] = s
	}
}

// List returns a string slice containing all the available subsystems registered with clog
func (r *Registry) List() (out []string) {
	for _, x := range *r {
		out = append(out, x.Name)
	}
	return
}

// Get returns the subsystem. This could then be used to close or set its level eg `*r.Get("subsystem").SetLevel("debug")`
func (r *Registry) Get(name string) (out *SubSystem) {
	var ok bool
	if out, ok = (*r)[name]; ok {
		return
	}
	return
}

// SetAllLevels is
func (r *Registry) SetAllLevels(level string) {
	level = sanitizeLoglevel(level)
	GlobalLevel = level
	loggers := r.List()
	for _, x := range loggers {
		r.Get(x).SetLevel(level)
	}
}
