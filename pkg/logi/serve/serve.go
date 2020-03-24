package serve

import (
	"errors"
	"github.com/p9c/pod/pkg/kopachctrl/job"
	"io"
	"net/rpc"
)

type Serve struct {
	*rpc.Client
}

func New(conn io.ReadWriteCloser) *Serve {
	return &Serve{}
}


// Log sends a new log message
func (c *Serve) Log(job *job.Container) (err error) {
	L.Debug("starting logger feed")
	var reply bool
	err = c.Call("Consume.Log", job, &reply)
	if err != nil {
		L.Error(err)
		return
	}
	if reply != true {
		err = errors.New("start signal not acknowledged")
	}
	return
}

type API struct {

}

// Run starts serving Log
func (a *API) Run(cmd bool, reply *bool) (err error) {
	r := true
	reply = &r
	return
}


// Pause pauses serving Log
func (a *API) Pause(cmd bool, reply *bool) (err error) {
	r := false
	reply = &r
	return
}
