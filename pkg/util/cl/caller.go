package cl

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/mitchellh/colorstring"
)

// Ine (cl.Ine) returns caller location in source code
var Ine = func() string {
	_, file, line, _ := runtime.Caller(1)
	files := strings.Split(file, "github.com/parallelcointeam/parallelcoin/")
	file = "./" + files[1]
	return colorstring.Color(fmt.Sprintf(" [dim]%s:%d", file, line))
}

// Ine3 (cl.Ine) returns caller location in source code
var Ine3 = func() string {
	_, file, line, _ := runtime.Caller(3)
	return colorstring.Color(fmt.Sprintf(" [dim]%s:%d", file, line))
}
