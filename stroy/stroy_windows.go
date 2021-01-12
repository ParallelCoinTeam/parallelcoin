package main

import (
	"os/exec"
	"strings"
	"syscall"
)

func WindowsExec(split []string) (out *exec.Cmd) {
	out = exec.Command(split[0])
	out.SysProcAttr = &syscall.SysProcAttr{}
	out.SysProcAttr.CmdLine = strings.Join(split, " ")
	return
}
