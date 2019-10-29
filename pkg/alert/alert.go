package alert

import (
	"time"
)

type DuOSalert struct {
	Time      time.Time   `json:"time"`
	Title     string      `json:"title"`
	Message   interface{} `json:"message"`
	AlertType string      `json:"type"`
}

// GetMsg loads the message variable
func (d *DuOSalert) PushDuOSalert(t string, m interface{}, at string) {
	a := new(DuOSalert)
	a.Time = time.Now()
	a.Title = t
	a.Message = m
	a.AlertType = at
	//d.Render("alert", a)
}
