package log

var L = Empty()

func FATAL(a ...interface{}) {
	if L.Fatal != nil {
		(*L).Fatal(a...)
	}
}
func ERROR(a ...interface{}) {
	if L.Error != nil {
		(*L).Error(a...)
	}
}
func WARN(a ...interface{}) {
	if L.Warn != nil {
		(*L).Warn(a...)
	}
}
func INFO(a ...interface{}) {
	if L.Info != nil {
		(*L).Info(a...)
	}
}
func DEBUG(a ...interface{}) {
	if L.Debug != nil {
		(*L).Debug(a...)
	}
}
func TRACE(a ...interface{}) {
	if L.Trace != nil {
		(*L).Trace(a...)
	}
}
func FATALF(format string, a ...interface{}) {
	if L.Fatalf != nil {
		(*L).Fatalf(format, a...)
	}
}
func ERRORF(format string, a ...interface{}) {
	if L.Errorf != nil {
		(*L).Errorf(format, a...)
	}
}
func WARNF(format string, a ...interface{}) {
	if L.Warnf != nil {
		(*L).Warnf(format, a...)
	}
}
func INFOF(format string, a ...interface{}) {
	if L.Infof != nil {
		(*L).Infof(format, a...)
	}
}
func DEBUGF(format string, a ...interface{}) {
	if L.Debugf != nil {
		(*L).Debugf(format, a...)
	}
}
func TRACEF(format string, a ...interface{}) {
	if L.Tracef != nil {
		(*L).Tracef(format, a...)
	}
}
func FATALC(f func() string) {
	if L.Fatalc != nil {
		(*L).Fatalc(f)
	}
}
func ERRORC(f func() string) {
	if L.Errorc != nil {
		(*L).Errorc(f)
	}
}
func WARNC(f func() string) {
	if L.Warnc != nil {
		(*L).Warnc(f)
	}
}
func INFOC(f func() string) {
	if L.Infoc != nil {
		(*L).Infoc(f)
	}
}
func DEBUGC(f func() string) {
	if L.Debugc != nil {
		(*L).Debugc(f)
	}
}
func TRACEC(f func() string) {
	if L.Tracec != nil {
		(*L).Tracec(f)
	}
}
