package rcd

import (
	"github.com/p9c/pod/cmd/gui/models"
	"time"
)

var ALERT = models.DuoUIalert{}


// GetMsg loads the message variable
func (r *RcVar) PushDuoUIalert(t string, m interface{}, at string) (d *models.DuoUIalert) {
	a := new(models.DuoUIalert)
	a.Time = time.Now()
	a.Title = t
	a.Message = m
	a.AlertType = at
	//d.Render("alert", a)
	return
}
