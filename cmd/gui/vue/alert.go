// +build !headless

package vue

import (
	"fmt"
	"time"
)

type DuoVUEalert struct {
	Time      time.Time   `json:"time"`
	Alert     interface{} `json:"alert"`
	AlertType string      `json:"type"`
}


// GetMsg loads the message variable
func (dv *DuoVUE) PushDuoVUEalert(m interface{}, t string) {
			a := new(DuoVUEalert)
			a.Time = time.Now()
			a.Alert = m
			a.AlertType = t
			dv.Render("alert", a)
			fmt.Println("dadafddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
	fmt.Println("dadafggf da", a)
}
