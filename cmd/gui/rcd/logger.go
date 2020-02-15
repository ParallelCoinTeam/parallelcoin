package rcd

import (
	"github.com/p9c/pod/pkg/log"
)

func (r *RcVar) DuoUIloggerController() {
	log.L.LogChan = make(chan log.Entry)
	r.Log.LogChan = log.L.LogChan
	log.L.SetLevel(*r.cx.Config.LogLevel, false)
	go func() {
	out:
		for {
			select {
			case n := <-log.L.LogChan:
				r.Log.LogMessages = append(r.Log.LogMessages, n)
			case <-r.Log.StopLogger:
				defer func() {
					r.Log.StopLogger = make(chan struct{})
				}()
				r.Log.LogMessages = []log.Entry{}
				log.L.LogChan = nil
				break out
			}
		}
	}()
}
