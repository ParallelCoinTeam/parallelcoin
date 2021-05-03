// Copyright 2013 @atotto. All rights reserved. Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// +build darwin

package clipboard

import (
	"fmt"
	"os"
	"os/exec"
)

var (
	pasteCmdArgs = "pbpaste"
	copyCmdArgs  = "pbcopy"
)

func getPasteCommand() *exec.Cmd {
	cmd := exec.Command(pasteCmdArgs)
	cmd.Env = []string{"LANG=en_US.UTF-8"}
	return cmd
}

func getCopyCommand() *exec.Cmd {
	return exec.Command(copyCmdArgs)
}

func readAll() (string, error) {
	pasteCmd := getPasteCommand()
	out, e := pasteCmd.Output()
	if e != nil  {
		return "", e
	}
	return string(out), nil
}

func writeAll(text string) (e error) {
	copyCmd := getCopyCommand()
	in, e := copyCmd.StdinPipe()
	if e != nil  {
		return e
	}

	if e := copyCmd.Start(); e != nil {
		return e
	}
	if _, e = in.Write([]byte(text)); e != nil {
		return e
	}
	if e := in.Close(); e != nil {
		return e
	}
	return copyCmd.Wait()
}

func Start() {
}

func Get() string {
	str, e := readAll()
	if e != nil  {
		fmt.Fprintln(os.Stderr, e)
		return ""
	}
	return str
}

func GetPrimary() string {
	return ""
}

func Set(text string) {
	e := writeAll(text)
	if e != nil  {
		fmt.Fprintln(os.Stderr, e)
	}
}
