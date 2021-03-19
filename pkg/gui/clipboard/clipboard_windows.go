// Copyright 2013 @atotto. All rights reserved. Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package clipboard

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	cfUnicodetext = 13
	gmemFixed     = 0x0000
)

var (
	user32           = syscall.MustLoadDLL("user32")
	openClipboard    = user32.MustFindProc("OpenClipboard")
	closeClipboard   = user32.MustFindProc("CloseClipboard")
	emptyClipboard   = user32.MustFindProc("EmptyClipboard")
	getClipboardData = user32.MustFindProc("GetClipboardData")
	setClipboardData = user32.MustFindProc("SetClipboardData")
	
	kernel32     = syscall.NewLazyDLL("kernel32")
	globalAlloc  = kernel32.NewProc("GlobalAlloc")
	globalFree   = kernel32.NewProc("GlobalFree")
	globalLock   = kernel32.NewProc("GlobalLock")
	globalUnlock = kernel32.NewProc("GlobalUnlock")
	lstrcpy      = kernel32.NewProc("lstrcpyW")
)

func readAll() (string, error) {
	var e error
	var r uintptr
	r, _, e = openClipboard.Call(0)
	if r == 0 {
		return "", e
	}
	defer func() {
		if _, _, e = closeClipboard.Call(); E.Chk(e) {
		}
	}()
	
	var h uintptr
	h, _, e = getClipboardData.Call(cfUnicodetext)
	if r == 0 {
		return "", e
	}
	
	var l uintptr
	l, _, e = globalLock.Call(h)
	if l == 0 {
		return "", e
	}
	
	text := syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(l))[:])
	
	r, _, e = globalUnlock.Call(h)
	if r == 0 {
		return "", e
	}
	
	return text, nil
}

func writeAll(text string) (e error) {
	var r uintptr
	r, _, e = openClipboard.Call(0)
	if r == 0 {
		return e
	}
	defer func() {
		if _, _, e = closeClipboard.Call(); E.Chk(e) {
		}
	}()
	
	r, _, e = emptyClipboard.Call(0)
	if r == 0 {
		return e
	}
	
	data := syscall.StringToUTF16(text)
	
	var h uintptr
	h, _, e = globalAlloc.Call(gmemFixed, uintptr(len(data)*int(unsafe.Sizeof(data[0]))))
	if h == 0 {
		return e
	}
	
	var l uintptr
	l, _, e = globalLock.Call(h)
	if l == 0 {
		return e
	}
	
	r, _, e = lstrcpy.Call(l, uintptr(unsafe.Pointer(&data[0])))
	if r == 0 {
		return e
	}
	
	r, _, e = globalUnlock.Call(h)
	if r == 0 {
		return e
	}
	
	r, _, e = setClipboardData.Call(cfUnicodetext, h)
	if r == 0 {
		return e
	}
	return nil
}

// Start ...
func Start() {
}

func Get() string {
	str, e := readAll()
	if e != nil  {
		_, _ = fmt.Fprintln(os.Stderr, e)
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
		_, _ = fmt.Fprintln(os.Stderr, e)
	}
}
