// SPDX-License-Identifier: Unlicense OR MIT

// Package widget implements state tracking and event handling of
// common user interface controls. To draw widgets, use a theme
// packages such as package gioui.org/widget/material.
package controller

import "time"

type Command struct {
	Com      interface{}
	ComID    string
	Category string
	Out      func()
	Time     time.Time
}
