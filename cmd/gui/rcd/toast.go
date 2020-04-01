package rcd

import (
	"github.com/p9c/pod/cmd/gui/model"
)

func (r *RcVar) toastAdd(t, m string) {
	r.Toasts = append(r.Toasts, model.DuoUItoast{
		Title:   t,
		Message: m,
	})
}
