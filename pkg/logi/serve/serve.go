package consume

import (
	"errors"
	"github.com/p9c/pod/pkg/kopachctrl/job"
	"io"
	"net/rpc"
)

type Logs struct {
	*rpc.Client
}

func New(conn io.ReadWriteCloser) *Logs {
	return &Logs{}
}


// Log receives a new log message
func (c *Logs) Log(job *job.Container) (err error) {
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
