package log

func Check(err error) (isNotNil bool) {
	if isNotNil = err != nil; isNotNil {
		ERROR(err)
	}
	return
}

func FATAL(a ...interface{}) {
	if L.Fatal != nil {
		f := L.Fatal
		fd := *f
		fd(a...)
	}
}

func ERROR(a ...interface{}) {
	if L.Error != nil {
		f := L.Error
		fd := *f
		fd(a...)
	}
}

func WARN(a ...interface{}) {
	if L.Warn != nil {
		f := L.Warn
		fd := *f
		fd(a...)
	}
}

func INFO(a ...interface{}) {
	if L.Info != nil {
		f := L.Info
		fd := *f
		fd(a...)
	}
}

func DEBUG(a ...interface{}) {
	if L.Debug != nil {
		f := L.Debug
		fd := *f
		fd(a...)
	}
}

func TRACE(a ...interface{}) {
	if L.Trace != nil {
		f := L.Trace
		fd := *f
		fd(a...)
	}
}

func SPEW(a interface{}) {
	if L.Traces != nil {
		f := L.Traces
		fd := *f
		fd(a)
	}
}

func FATALF(format string, a ...interface{}) {
	if L.Fatalf != nil {
		f := L.Fatalf
		fd := *f
		fd(format, a...)
	}
}

func ERRORF(format string, a ...interface{}) {
	if L.Errorf != nil {
		f := L.Errorf
		fd := *f
		fd(format, a...)
	}
}

func WARNF(format string, a ...interface{}) {
	if L.Warnf != nil {
		f := L.Warnf
		fd := *f
		fd(format, a...)
	}
}

func INFOF(format string, a ...interface{}) {
	if L.Infof != nil {
		f := L.Infof
		fd := *f
		fd(format, a...)
	}
}

func DEBUGF(format string, a ...interface{}) {
	if L.Debugf != nil {
		f := L.Debugf
		fd := *f
		fd(format, a...)
	}
}

func TRACEF(format string, a ...interface{}) {
	if L.Tracef != nil {
		f := L.Tracef
		fd := *f
		fd(format, a...)
	}
}

func FATALC(fn func() string) {
	if L.Fatalc != nil {
		f := L.Fatalc
		fd := *f
		fd(fn)
	}
}

func ERRORC(fn func() string) {
	if L.Errorc != nil {
		f := L.Errorc
		fd := *f
		fd(fn)
	}
}

func WARNC(fn func() string) {
	if L.Warnc != nil {
		f := L.Warnc
		fd := *f
		fd(fn)
	}
}

func INFOC(fn func() string) {
	if L.Infoc != nil {
		f := L.Infoc
		fd := *f
		fd(fn)
	}
}

func DEBUGC(fn func() string) {
	if L.Debugc != nil {
		f := L.Debugc
		fd := *f
		fd(fn)
	}
}

func TRACEC(fn func() string) {
	if L.Tracec != nil {
		f := L.Tracec
		fd := *f
		fd(fn)
	}
}
