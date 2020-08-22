package old

import (
	"github.com/stalker-loki/pod/cmd/gui/model"
)

func (r *RcVar) toastAdd(t, m string) {
	r.Toasts = append(r.Toasts, model.DuoUItoast{
		Title:   t,
		Message: m,
	})
}
