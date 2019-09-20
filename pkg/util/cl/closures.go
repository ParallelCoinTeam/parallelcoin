package cl

// Ftlc appends the subsystem name to the front of a closure's output and this runs only if the log entry is called
func (s *SubSystem) Ftlc(closure StringClosure) {
	if s.Level > _off {
		s.Ch <- Fatalc(closure)
	}
}

// Errc appends the subsystem name to the front of a closure's output and this runs only if the log entry is called
func (s *SubSystem) Errc(closure StringClosure) {
	if s.Level > _fatal {
		s.Ch <- Errorc(closure)
	}
}

// Wrnc appends the subsystem name to the front of a closure's output and this runs only if the log entry is called
func (s *SubSystem) Wrnc(closure StringClosure) {
	if s.Level > _error {
		s.Ch <- Warnc(closure)
	}
}

// Infc appends the subsystem name to the front of a closure's output and this runs only if the log entry is called
func (s *SubSystem) Infc(closure StringClosure) {
	if s.Level > _warn {
		s.Ch <- Infoc(closure)
	}
}

// Dbgc appends the subsystem name to the front of a closure's output and this runs only if the log entry is called
func (s *SubSystem) Dbgc(closure StringClosure) {
	if s.Level > _info {
		s.Ch <- Debugc(closure)
	}
}

// Trcc appends the subsystem name to the front of a closure's output and this runs only if the log entry is called
func (s *SubSystem) Trcc(closure StringClosure) {
	if s.Level > _debug {
		s.Ch <- Tracec(closure)
	}
}
