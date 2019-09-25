// +build !headless

package vue

import (
	"time"
)

type DuoVUEalert struct {
	Time      time.Time   `json:"time"`
	Title    string      `json:"title"`
	Message   interface{} `json:"message"`
	AlertType string      `json:"type"`
}

// GetMsg loads the message variable
func (dv *DuoVUE) PushDuoVUEalert(t string, m interface{}, at string) {
	a := new(DuoVUEalert)
	a.Time = time.Now()
	a.Title = t
	a.Message = m
	a.AlertType = at
	dv.Render("alert", a)
}
