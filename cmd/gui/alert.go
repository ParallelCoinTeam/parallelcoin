//+build !headless

package gui

import (
	"time"
)

var ALERT = DuOSalert{}

type DuOSalert struct {
	Time      time.Time   `json:"time"`
	Title     string      `json:"title"`
	Message   interface{} `json:"message"`
	AlertType string      `json:"type"`
}

// GetMsg loads the message variable
func (r *rcvar) PushDuOSalert(t string, m interface{}, at string) (d *DuOSalert) {
	a := new(DuOSalert)
	a.Time = time.Now()
	a.Title = t
	a.Message = m
	a.AlertType = at
	//d.Render("alert", a)
	return
}
