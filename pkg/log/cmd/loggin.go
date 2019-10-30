package main

import (
	"fmt"

	"github.com/p9c/pod/pkg/log"
)

func main() {
	fmt.Println("testing level trace")
	l := log.NewLogger("trace")
	printAll(l, "log level = trace")
	fmt.Println("testing level debug")
	l.SetLevel("debug")
	printAll(l, "log level = debug")
	fmt.Println("testing level info")
	l.SetLevel("info")
	printAll(l, "log level = info")
	fmt.Println("testing level warn")
	l.SetLevel("warn")
	printAll(l, "log level = warn")
	fmt.Println("testing level error")
	l.SetLevel("error")
	printAll(l, "log level = error")
	fmt.Println("testing level fatal")
	l.SetLevel("fatal")
	printAll(l, "log level = fatal")
	fmt.Println("testing level off")
	l.SetLevel("off")
	printAll(l, "log level = off")
}

func printAll(l *log.Logger, text string) {
	l.Fatal(text, 1.0)
	l.Error(text, 1.0)
	l.Warn(text, 1.0)
	l.Info(text, 1.0)
	log.DEBUG(text, 1.0)
	l.Trace(text, 1.0)
	l.Fatalf("%s %f", text, 1.0)
	l.Errorf("%s %f", text, 1.0)
	l.Warnf("%s %f", text, 1.0)
	l.Infof("%s %f", text, 1.0)
	l.Debugf("%s %f", text, 1.0)
	l.Tracef("%s %f", text, 1.0)
	l.Fatalc(func() string {
		return fmt.Sprintln(text, 1.0)
	})
	l.Errorc(func() string {
		return fmt.Sprintln(text, 1.0)
	})
	l.Warnc(func() string {
		return fmt.Sprintln(text, 1.0)
	})
	l.Infoc(func() string {
		return fmt.Sprintln(text, 1.0)
	})
	l.Debugc(func() string {
		return fmt.Sprintln(text, 1.0)
	})
	l.Tracec(func() string {
		return fmt.Sprintln(text, 1.0)
	})
}